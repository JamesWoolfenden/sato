package arm

import "fmt"

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
