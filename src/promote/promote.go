// Package promote installs an AI-drafted template into the source tree so it
// becomes a built-in conversion. It is a maintainer tool: it edits Go source
// in-place and only works when run from a sato checkout.
package promote

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type target struct {
	resourcesDir string
	embedsFile   string
	mapFile      string
	mapClose     string
	keyTransform func(string) string
}

var targets = map[string]target{
	"aws": {
		resourcesDir: "src/cf/resources",
		embedsFile:   "src/cf/resources.go",
		mapFile:      "src/cf/resource_mapping.go",
		mapClose:     "\n}\n",
		keyTransform: func(s string) string { return s },
	},
	"azurerm": {
		resourcesDir: "src/arm/resources",
		embedsFile:   "src/arm/resource.go",
		mapFile:      "src/arm/lookup.go",
		mapClose:     "\n\t}\n",
		keyTransform: strings.ToLower,
	},
}

// Run copies draft into the source tree and registers it in the embed and
// lookup tables. sourceType is the upstream type (e.g. AWS::Scheduler::Schedule
// or Microsoft.Compute/disks). root is the repo root.
func Run(draft, sourceType, root string) error {
	tfType := strings.TrimSuffix(filepath.Base(draft), ".template")

	provider, _, ok := strings.Cut(tfType, "_")
	if !ok {
		return fmt.Errorf("cannot infer provider from %q", tfType)
	}

	tgt, ok := targets[provider]
	if !ok {
		return fmt.Errorf("unsupported provider prefix %q", provider)
	}

	body, err := os.ReadFile(draft) // #nosec G304
	if err != nil {
		return err
	}

	dest := filepath.Join(root, tgt.resourcesDir, tfType+".template")
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("%s already exists; edit it directly", dest)
	}

	if err := os.WriteFile(dest, body, 0o600); err != nil { // #nosec G703 -- maintainer tool writes into checkout
		return err
	}

	varName := camel(tfType)

	if err := appendEmbed(filepath.Join(root, tgt.embedsFile), tfType, varName); err != nil {
		return err
	}

	entry := fmt.Sprintf("\t%q: %s,", tgt.keyTransform(sourceType), varName)
	if err := insertBefore(filepath.Join(root, tgt.mapFile), tgt.mapClose, entry); err != nil {
		return err
	}

	return updateSee(filepath.Join(root, "src/see/resource_mapping.go"), sourceType, tfType)
}

func camel(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}

	return strings.Join(parts, "")
}

func appendEmbed(path, tfType, varName string) error {
	body, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return err
	}

	if bytes.Contains(body, []byte("resources/"+tfType+".template")) {
		return nil
	}

	add := fmt.Sprintf("\n//go:embed resources/%s.template\nvar %s []byte\n", tfType, varName)

	return os.WriteFile(path, append(body, []byte(add)...), 0o600) // #nosec G703
}

func insertBefore(path, marker, line string) error {
	body, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return err
	}

	if bytes.Contains(body, []byte(line)) {
		return nil
	}

	idx := bytes.LastIndex(body, []byte(marker))
	if idx == -1 {
		return fmt.Errorf("marker not found in %s", path)
	}

	out := append([]byte{}, body[:idx]...)
	out = append(out, '\n')
	out = append(out, []byte(line)...)
	out = append(out, body[idx:]...)

	return os.WriteFile(path, out, 0o600) // #nosec G703
}

func updateSee(path, sourceType, tfType string) error {
	body, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%q:", strings.ToLower(sourceType))
	if i := bytes.Index(body, []byte(key)); i != -1 {
		end := bytes.IndexByte(body[i:], '\n')
		line := body[i : i+end]
		if bytes.Contains(line, []byte("none")) {
			newLine := fmt.Sprintf("%s %q,", key, tfType)
			out := bytes.Replace(body, line, []byte(newLine), 1)

			return os.WriteFile(path, out, 0o600) // #nosec G703
		}

		return nil
	}

	entry := fmt.Sprintf("\t%s %q,", key, tfType)

	return insertBefore(path, "\n}\n", entry)
}

// OpenPR branches, commits the promotion edits, pushes, and opens a GitHub PR
// via the gh CLI. Requires git and gh on PATH and an authenticated gh session.
func OpenPR(root, sourceType, tfType string) error {
	branch := "promote/" + tfType
	title := fmt.Sprintf("promote: add %s template (%s)", tfType, sourceType)

	steps := [][]string{
		{"git", "checkout", "-b", branch},
		{"git", "add", "src/"},
		{"git", "commit", "-m", title},
		{"git", "push", "-u", "origin", branch},
		{"gh", "pr", "create", "--title", title, "--body",
			"Promotes AI-drafted template for `" + sourceType + "` into the built-in lookup tables.", "--head", branch},
	}

	for _, s := range steps {
		cmd := exec.Command(s[0], s[1:]...) // #nosec G204
		cmd.Dir = root
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s: %w", strings.Join(s, " "), err)
		}
	}

	return nil
}
