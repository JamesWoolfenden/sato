package cf

import (
	"fmt"
	"sato/src/satoerrors"
)

type (
	filepathError       = satoerrors.FilepathError
	parseVariablesError = satoerrors.ParseVariablesError
	parseResourcesError = satoerrors.ParseResourcesError
)

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

func (m *goformationError) Unwrap() error {
	return m.err
}

type writeError struct {
	destination string
	err         error
}

func (e *writeError) Error() string {
	return fmt.Sprintf("write failed %s %v", e.destination, e.err)
}

func (e *writeError) Unwrap() error {
	return e.err
}

type emptyPathsError struct{}

func (e emptyPathsError) Error() string {
	return "paths cannot be empty"
}
