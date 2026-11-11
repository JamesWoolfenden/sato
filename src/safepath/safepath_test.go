package safepath

import (
	"errors"
	"path/filepath"
	"testing"

	"sato/src/satoerrors"
)

func TestJoin(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	for _, name := range []string{"a", "aws_lambda_function.my_fn", "x.y.z", "..foo", "foo.."} {
		got, err := Join(dir, name)
		if err != nil {
			t.Errorf("Join(%q) unexpected error: %v", name, err)
		}

		if got != filepath.Join(dir, name) {
			t.Errorf("Join(%q) = %q", name, got)
		}
	}

	for _, name := range []string{"", ".", "..", "../x", "../../conftest", "a/b", `a\b`, "/etc/passwd", `C:\x`, "a\x00b"} {
		_, err := Join(dir, name)

		var unsafe *satoerrors.UnsafePathError
		if err == nil || !errors.As(err, &unsafe) && name != "" {
			t.Errorf("Join(%q) = %v, want UnsafePathError", name, err)
		}
	}
}
