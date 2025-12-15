package cf

import (
	"fmt"
	"sato/src/satoerrors"
)

// filepathError is an alias for the satoerrors error type.
type filepathError = satoerrors.FilepathError

// writeFileError is an alias for the satoerrors error type.
type writeFileError = satoerrors.WriteFileError

// makeDirError is an alias for the satoerrors error type.
type makeDirError = satoerrors.MakeDirError

// missingResourceError represents a resource lookup failure.
type missingResourceInputError struct{}

func (e *missingResourceInputError) Error() string {
	return "invalid input parameters"
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
	return fmt.Sprintf("parse variables failure %v", m.err)
}

type parseResourcesError struct {
	err error
}

func (m *parseResourcesError) Error() string {
	return fmt.Sprintf("parse resources failure %v", m.err)
}

type writeError struct {
	destination string
	err         error
}

func (e *writeError) Error() string {
	return fmt.Sprintf("write failed %s %v", e.destination, e.err)
}

type emptyPathsError struct{}

func (e emptyPathsError) Error() string {
	return "paths cannot be empty"
}
