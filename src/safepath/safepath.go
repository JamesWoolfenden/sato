// Package safepath keeps names taken from untrusted templates inside the directory they are written to.
package safepath

import (
	"path/filepath"
	"strings"

	"sato/src/satoerrors"
)

// Join returns dir/name, or an error when name could place the result outside dir.
// name is a single file name, never a path: separators, NUL bytes and "." / ".." are rejected,
// and the cleaned result must still resolve under dir.
func Join(dir string, name string) (string, error) {
	if name == "" || name == "." || name == ".." ||
		strings.ContainsAny(name, "/\\"+"\x00") {
		return "", &satoerrors.UnsafePathError{Name: name}
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", &satoerrors.FilepathError{Path: dir, Err: err}
	}

	joined := filepath.Join(absDir, name)

	rel, err := filepath.Rel(absDir, joined)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", &satoerrors.UnsafePathError{Name: name}
	}

	return joined, nil
}
