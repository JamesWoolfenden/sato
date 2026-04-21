// Package tfgen provides provider-agnostic types and template helpers
// used by both the cf (CloudFormation) and arm (Azure) converters.
package tfgen

// M is a wrapper object to help pass multiple values into a text/template.
type M map[string]interface{}

// Variable describes a Terraform variable.
type Variable struct {
	Description string
	Type        string
	Default     string
	Name        string
}

// Output describes a Terraform output.
type Output struct {
	Description string
	Type        string
	Value       string
	Name        string
}
