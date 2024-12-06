package cf

import "fmt"

type filepathError struct {
	Path string
	err  error
}

// missingResourceError represents a resource lookup failure
type missingResourceInputError struct {
}

func (e *missingResourceInputError) Error() string {
	return "invalid input parameters"
}

func (m *filepathError) Error() string {
	return fmt.Sprintf("not implemented %s as raised %v", m.Path, m.err)
}

type goformationError struct {
	err error
}

func (m *goformationError) Error() string {
	return fmt.Sprintf("goformation parse failure %v", m.err)
}

type parseVariablesError struct {
	err error
}

func (m *parseVariablesError) Error() string {
	return fmt.Sprintf("parse varriables failure %v", m.err)
}

type parseResourcesError struct {
	err error
}

func (m *parseResourcesError) Error() string {
	return fmt.Sprintf("parse resources failure %v", m.err)
}

type makeDirError struct {
	err error
}

func (e *makeDirError) Error() string {
	return fmt.Sprintf("mkdir failed %v", e.err)
}

type writeError struct {
	destination string
	err         error
}

func (e *writeError) Error() string {
	return fmt.Sprintf("write failed %s %v", e.destination, e.err)
}

type writeFileError struct {
	destination string
	err         error
}

type emptyPathsError struct{}

func (e emptyPathsError) Error() string {
	return "paths cannot be empty"
}

func (e *writeFileError) Error() string {
	return fmt.Sprintf("write file failed %s %v", e.destination, e.err)
}
