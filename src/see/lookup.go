package see

import (
	"fmt"
	"strings"
)

const none string = "none"

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
		Reversed := reverseMap(lookupMapping)
		result = Reversed[resource]
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

// reverseMap creates a new map with keys and values swapped.
// It assumes unique values in the input map to avoid conflicts.
func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		n[v] = k
	}

	return n
}
