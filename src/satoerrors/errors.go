// Package satoerrors provides shared error types used across sato packages.
package satoerrors

import "fmt"

// FilepathError represents an error resolving a file path.
type FilepathError struct {
	Path string
	Err  error
}

func (f *FilepathError) Error() string {
	if f.Err != nil {
		return fmt.Sprintf("filepath error for %s: %v", f.Path, f.Err)
	}
	return fmt.Sprintf("filepath error for %s", f.Path)
}

func (f *FilepathError) Unwrap() error {
	return f.Err
}

// WriteFileError represents an error writing a file.
type WriteFileError struct {
	Destination string
	Err         error
}

func (w *WriteFileError) Error() string {
	return fmt.Sprintf("failed to write file %s: %v", w.Destination, w.Err)
}

func (w *WriteFileError) Unwrap() error {
	return w.Err
}

// MakeDirError represents an error creating a directory.
type MakeDirError struct {
	Err error
}

func (m *MakeDirError) Error() string {
	return fmt.Sprintf("failed to create directory: %v", m.Err)
}

func (m *MakeDirError) Unwrap() error {
	return m.Err
}

// TemplateNewError represents an error creating a template.
type TemplateNewError struct {
	Err error
}

func (t *TemplateNewError) Error() string {
	return fmt.Sprintf("failed to create template: %v", t.Err)
}

func (t *TemplateNewError) Unwrap() error {
	return t.Err
}

// TemplateExecuteError represents an error executing a template.
type TemplateExecuteError struct {
	Err error
}

func (t *TemplateExecuteError) Error() string {
	return fmt.Sprintf("failed to execute template: %v", t.Err)
}

func (t *TemplateExecuteError) Unwrap() error {
	return t.Err
}

// ParseVariablesError represents a failure converting template parameters to Terraform variables.
type ParseVariablesError struct {
	Err error
}

func (p *ParseVariablesError) Error() string {
	return fmt.Sprintf("parse variables failure: %v", p.Err)
}

func (p *ParseVariablesError) Unwrap() error {
	return p.Err
}

// ParseResourcesError represents a failure converting template resources to Terraform resources.
type ParseResourcesError struct {
	Err error
}

func (p *ParseResourcesError) Error() string {
	return fmt.Sprintf("parse resources failure: %v", p.Err)
}

func (p *ParseResourcesError) Unwrap() error {
	return p.Err
}
