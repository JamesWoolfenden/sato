package arm

import (
	"errors"
	"testing"
)

func TestSplitResourceError(t *testing.T) {
	tests := []struct {
		name     string
		match    string
		expected string
	}{
		{
			name:     "basic resource",
			match:    "resource/123",
			expected: "failed to split resource resource/123",
		},
		{
			name:     "empty resource",
			match:    "",
			expected: "failed to split resource ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := splitResourceError{match: tt.match}
			if got := err.Error(); got != tt.expected {
				t.Errorf("splitResourceError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFilepathError(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "valid path",
			path:     "/path/to/file",
			expected: "not implemented /path/to/file",
		},
		{
			name:     "empty path",
			path:     "",
			expected: "not implemented ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &filepathError{Path: tt.path}
			if got := err.Error(); got != tt.expected {
				t.Errorf("filepathError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseListError(t *testing.T) {
	err := &parseListError{}
	expected := "parseListError"
	if got := err.Error(); got != expected {
		t.Errorf("parseListError.Error() = %v, want %v", got, expected)
	}
}

func TestParseMapError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "with inner error",
			err:      errors.New("inner error"),
			expected: "parseMapError inner error",
		},
		{
			name:     "nil error",
			err:      nil,
			expected: "parseMapError <nil>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &parseMapError{Err: tt.err}
			if got := err.Error(); got != tt.expected {
				t.Errorf("parseMapError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestEmptyResourceError(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		expected string
	}{
		{
			name:     "resource name",
			resource: "myResource",
			expected: "myResource is empty",
		},
		{
			name:     "empty name",
			resource: "",
			expected: " is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &emptyResourceError{Name: tt.resource}
			if got := err.Error(); got != tt.expected {
				t.Errorf("emptyResourceError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseResourceError(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		expected string
	}{
		{
			name:     "resource name",
			resource: "myResource",
			expected: "failed to parse resource name myResource",
		},
		{
			name:     "empty name",
			resource: "",
			expected: "failed to parse resource name ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &parseResourceError{Name: tt.resource}
			if got := err.Error(); got != tt.expected {
				t.Errorf("parseResourceError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInlineFormatError(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		expected string
	}{
		{
			name:     "format name",
			resource: "myFormat",
			expected: "uses inline format function myFormat",
		},
		{
			name:     "empty format",
			resource: "",
			expected: "uses inline format function ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &inlineFormatError{Name: tt.resource}
			if got := err.Error(); got != tt.expected {
				t.Errorf("inlineFormatError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestMatchValueError(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "value name",
			value:    "myValue",
			expected: "failed to match value myValue",
		},
		{
			name:     "empty value",
			value:    "",
			expected: "failed to match value ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &matchValueError{Name: tt.value}
			if got := err.Error(); got != tt.expected {
				t.Errorf("matchValueError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}
