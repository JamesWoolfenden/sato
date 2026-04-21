package arm

import (
	"fmt"
	"sato/src/satoerrors"
)

// filepathError is an alias for the satoerrors error type.
type filepathError = satoerrors.FilepathError

// writeFileError is an alias for the satoerrors error type.
type writeFileError = satoerrors.WriteFileError

// templateNewError is an alias for the satoerrors error type.
type templateNewError = satoerrors.TemplateNewError

// templateExecuteError is an alias for the satoerrors error type.
type templateExecuteError = satoerrors.TemplateExecuteError

// parseVariablesError is an alias for the satoerrors error type.
type parseVariablesError = satoerrors.ParseVariablesError

// parseResourcesError is an alias for the satoerrors error type.
type parseResourcesError = satoerrors.ParseResourcesError

type splitResourceError struct {
	match string
}

func (e splitResourceError) Error() string {
	return fmt.Sprintf("failed to split resource %s", e.match)
}

type parseListError struct{}

func (m *parseListError) Error() string {
	return "parseListError"
}

type parseMapError struct {
	Err error
}

func (m *parseMapError) Error() string {
	return fmt.Sprintf("parseMapError %v", m.Err)
}

func (m *parseMapError) Unwrap() error {
	return m.Err
}

type emptyResourceError struct {
	Name string
}

func (m *emptyResourceError) Error() string {
	return fmt.Sprintf("%s is empty", m.Name)
}

type parseResourceError struct {
	Name string
}

func (m *parseResourceError) Error() string {
	return fmt.Sprintf("failed to parse resource name %s", m.Name)
}

type inlineFormatError struct {
	Name string
}

func (m *inlineFormatError) Error() string {
	return fmt.Sprintf("uses inline format function %s", m.Name)
}

type matchValueError struct {
	Name string
}

func (m *matchValueError) Error() string {
	return fmt.Sprintf("failed to match value %s", m.Name)
}
