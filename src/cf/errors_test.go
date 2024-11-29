package cf

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
			name:     "basic error",
			path:     "/test/path",
			err:      errors.New("test error"),
			expected: "not implemented /test/path as raised test error",
		},
		{
			name:     "empty path",
			path:     "",
			err:      errors.New("no path"),
			expected: "not implemented  as raised no path",
		},
		{
			name:     "nil error",
			path:     "/some/path",
			err:      nil,
			expected: "not implemented /some/path as raised <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fe := &filepathError{
				Path: tt.path,
				err:  tt.err,
			}
			if got := fe.Error(); got != tt.expected {
				t.Errorf("filepathError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMissingResourceInputError(t *testing.T) {
	err := &missingResourceInputError{}
	expected := "invalid input parameters"
	if got := err.Error(); got != expected {
		t.Errorf("missingResourceInputError.Error() = %v, want %v", got, expected)
	}
}

func TestGoformationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "with error",
			err:      errors.New("formation failed"),
			expected: "goformation parse failure formation failed",
		},
		{
			name:     "nil error",
			err:      nil,
			expected: "goformation parse failure <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ge := &goformationError{err: tt.err}
			if got := ge.Error(); got != tt.expected {
				t.Errorf("goformationError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseVariablesError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "with error",
			err:      errors.New("var parse failed"),
			expected: "parse varriables failure var parse failed",
		},
		{
			name:     "nil error",
			err:      nil,
			expected: "parse varriables failure <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pe := &parseVariablesError{err: tt.err}
			if got := pe.Error(); got != tt.expected {
				t.Errorf("parseVariablesError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseResourcesError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "with error",
			err:      errors.New("resource parse failed"),
			expected: "parse resources failure resource parse failed",
		},
		{
			name:     "nil error",
			err:      nil,
			expected: "parse resources failure <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pe := &parseResourcesError{err: tt.err}
			if got := pe.Error(); got != tt.expected {
				t.Errorf("parseResourcesError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}
