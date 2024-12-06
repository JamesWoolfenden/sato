package arm

import "fmt"

type splitResourceError struct {
	match string
}

func (e splitResourceError) Error() string {
	return fmt.Sprintf("failed to split resource %s", e.match)
}

type filepathError struct {
	Path string
}

func (m *filepathError) Error() string {
	return fmt.Sprintf("not implemented %s", m.Path)
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

type templateNewError struct {
	err error
}

func (e templateNewError) Error() string {
	return fmt.Sprintf("failed to build parser %v", e.err)
}

type templateExecuteError struct {
	err error
}

func (e templateExecuteError) Error() string {
	return fmt.Sprintf("failed to execute parser %v", e.err)
}

type writeFileError struct {
	destination string
	err         error
}

func (e writeFileError) Error() string {
	return fmt.Sprintf("failed to write %s %v", e.destination, e.err)
}
