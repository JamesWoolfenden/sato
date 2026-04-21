// Package see provides resource type mapping and lookup functionality for converting
// cloud provider resource types (AWS CloudFormation and Azure ARM) to Terraform resource types.
package see

import (
	"fmt"
	"strings"
)

const none string = "none"

var reverseMapping = reverseMap(lookupMapping)

// missingResourceError represents a resource lookup failure
type missingResourceError struct {
	Resource string
}

func (e *missingResourceError) Error() string {
	return fmt.Sprintf("resource %s not found", e.Resource)
}

// Lookup converts from cloudformation/ARM to terraform resource name.
func Lookup(resource string, reverse bool) (*string, error) {
	if resource == "" {
		return nil, &missingResourceError{
			Resource: resource,
		}
	}

	var result string

	if reverse {
		result = reverseMapping[resource]
	} else {
		result = lookupMapping[strings.TrimSuffix(strings.ToLower(resource), "/")]
	}

	if result == "" {
		return nil, &missingResourceError{
			Resource: resource,
		}
	}

	return &result, nil
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		if v == none {
			continue
		}
		n[v] = k
	}

	return n
}
