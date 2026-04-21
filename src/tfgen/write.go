package tfgen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"sato/src/satoerrors"

	"github.com/rs/zerolog/log"
)

// Write renders a Terraform file to disk and formats it with `tofu fmt`.
func Write(output string, location string, name string) error {
	if output == "" {
		return nil
	}

	newPath, err := filepath.Abs(location)
	if err != nil {
		return &satoerrors.FilepathError{Path: location, Err: err}
	}

	if err := os.MkdirAll(newPath, 0o750); err != nil {
		return &satoerrors.MakeDirError{Err: err}
	}

	destination, err := filepath.Abs(fmt.Sprint(location, "/", name, ".tf"))
	if err != nil {
		return &satoerrors.FilepathError{Path: fmt.Sprint(location, "/", name, ".tf"), Err: err}
	}

	if err := os.WriteFile(destination, []byte(output), 0o600); err != nil {
		return &satoerrors.WriteFileError{Destination: destination, Err: err}
	}

	log.Info().Msgf("Created %s", destination)

	cmd := exec.Command("tofu", "fmt", destination) // #nosec G204 -- destination is a validated filepath
	if err := cmd.Run(); err != nil {
		log.Warn().Msgf("Could not format %s with tofu fmt: %v", destination, err)
	} else {
		log.Info().Msgf("Formatted %s with tofu fmt", destination)
	}

	return nil
}
