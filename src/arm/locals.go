package arm

import (
	"strings"
)

// ParseLocals parses fields into locals.tf.
func ParseLocals(result map[string]interface{}) (string, map[string]interface{}, error) {
	var locals string

	myLocals, ok := result["locals"].(map[string]interface{})

	if !ok || myLocals == nil {
		return "", result, &emptyResourceError{}
	}

	for item, value := range myLocals {
		original, ok := value.(string)
		if !ok {
			continue
		}

		var (
			theValue string
			local    string
		)

		theValue, result = ParseString(original, result)

		myLocals[item] = theValue

		if strings.Contains(theValue, "${") {
			local = "\t" + item + " = \"" + theValue + "\" #" + original + "\n"
		} else {
			local = "\t" + item + " = " + theValue + " #" + original + "\n"
		}

		locals += strings.ReplaceAll(local, "'", "\"")
	}

	result["locals"] = myLocals

	return locals, result, nil
}
