// Package utils provides integration-test helpers for validating sato output
// with OpenTofu.
package utils

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TofuPath returns the path to the tofu binary, or skips the test if it is
// not installed.
func TofuPath(t *testing.T) string {
	t.Helper()

	p, err := exec.LookPath("tofu")
	if err != nil {
		t.Skip("tofu not found on PATH; skipping validation")
	}

	return p
}

func run(t *testing.T, tofu, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command(tofu, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "TF_IN_AUTOMATION=1")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("tofu %v in %s failed: %v\n%s", args, dir, err, out.String())
	}
}

// Init runs `tofu init -backend=false` in dir. It downloads providers, so a
// plugin cache is configured under the user cache dir to keep repeat runs fast.
func Init(t *testing.T, tofu, dir string) {
	t.Helper()

	if cache, err := os.UserCacheDir(); err == nil {
		cacheDir := filepath.Join(cache, "sato-tofu-plugins")
		_ = os.MkdirAll(cacheDir, 0o750)
		t.Setenv("TF_PLUGIN_CACHE_DIR", cacheDir)
	}

	run(t, tofu, dir, "init", "-backend=false", "-input=false", "-no-color")
}

// Validate runs `tofu validate` in dir. Init must have been called first.
func Validate(t *testing.T, tofu, dir string) {
	t.Helper()
	run(t, tofu, dir, "validate", "-no-color")
}
