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

type readFileError struct {
	path string
	err  error
}

func (r *readFileError) Error() string {
	return fmt.Sprintf("failed to read file: %s %v", r.path, r.err)
}

func (r *readFileError) Unwrap() error {
	return r.err
}

type openFileError struct {
	path string
	err  error
}

func (o *openFileError) Error() string {
	return fmt.Sprintf("failed to open file: %s %v", o.path, o.err)
}

func (o *openFileError) Unwrap() error {
	return o.err
}

type unmarshalError struct {
	err error
}

func (u *unmarshalError) Error() string {
	return fmt.Sprintf("unmarshal failure %v", u.err)
}

func (u *unmarshalError) Unwrap() error {
	return u.err
}

type parseDataError struct {
	err error
}

func (p *parseDataError) Error() string {
	return fmt.Sprintf("parse data failure %v", p.err)
}

func (p *parseDataError) Unwrap() error {
	return p.err
}
