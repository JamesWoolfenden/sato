package satoerrors

import (
	"errors"
	"testing"
)

func TestFilepathError(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		err      error
		expected string
	}{
		{
			name:     "with_error",
			path:     "/test/path",
			err:      errors.New("permission denied"),
			expected: "filepath error for /test/path: permission denied",
		},
		{
			name:     "nil_error",
			path:     "/some/path",
			err:      nil,
			expected: "filepath error for /some/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := &FilepathError{Path: tt.path, Err: tt.err}
			if got := fe.Error(); got != tt.expected {
				t.Errorf("FilepathError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFilepathErrorUnwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	fe := &FilepathError{Path: "/test", Err: innerErr}

	unwrapped := fe.Unwrap()
	if !errors.Is(unwrapped, innerErr) {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, innerErr)
	}
}

func TestWriteFileError(t *testing.T) {
	tests := []struct {
		name        string
		destination string
		err         error
		expected    string
	}{
		{
			name:        "with_error",
			destination: "/tmp/test.tf",
			err:         errors.New("disk full"),
			expected:    "failed to write file /tmp/test.tf: disk full",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			we := &WriteFileError{Destination: tt.destination, Err: tt.err}
			if got := we.Error(); got != tt.expected {
				t.Errorf("WriteFileError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestWriteFileErrorUnwrap(t *testing.T) {
	innerErr := errors.New("inner error")
	we := &WriteFileError{Destination: "/test", Err: innerErr}

	unwrapped := we.Unwrap()
	if !errors.Is(unwrapped, innerErr) {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, innerErr)
	}
}

func TestMakeDirError(t *testing.T) {
	innerErr := errors.New("permission denied")
	me := &MakeDirError{Err: innerErr}

	expected := "failed to create directory: permission denied"
	if got := me.Error(); got != expected {
		t.Errorf("MakeDirError.Error() = %v, want %v", got, expected)
	}
}

func TestTemplateNewError(t *testing.T) {
	innerErr := errors.New("syntax error")
	te := &TemplateNewError{Err: innerErr}

	expected := "failed to create template: syntax error"
	if got := te.Error(); got != expected {
		t.Errorf("TemplateNewError.Error() = %v, want %v", got, expected)
	}
}

func TestTemplateExecuteError(t *testing.T) {
	innerErr := errors.New("missing variable")
	tee := &TemplateExecuteError{Err: innerErr}

	expected := "failed to execute template: missing variable"
	if got := tee.Error(); got != expected {
		t.Errorf("TemplateExecuteError.Error() = %v, want %v", got, expected)
	}
}
