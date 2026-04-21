package cf

import (
	"strings"

	"github.com/awslabs/goformation/v7/cloudformation/tags"
)

// Tags renders a slice of CloudFormation tags as Terraform map entries.
func Tags(v []tags.Tag) string {
	var builder strings.Builder

	for _, item := range v {
		if item.Key != "" {
			builder.WriteString("\t\"")
			builder.WriteString(item.Key)
			builder.WriteString("\"=\"")
			builder.WriteString(item.Value)
			builder.WriteString("\"\n")
		}
	}

	return builder.String()
}
