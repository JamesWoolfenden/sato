package arm

import (
	"fmt"
	"strings"
)

// ParseLocals parses fields into locals.tf.
func ParseLocals(result map[string]interface{}) (string, map[string]interface{}, error) {
	var locals string

	myLocals, ok := result["locals"].(map[string]interface{})

	if !ok || myLocals == nil {
		return "", result, fmt.Errorf("locals is empty")
	}

	for item, value := range myLocals {
		var (
			theValue string
			local    string
		)

		theValue, result = ParseString(value.(string), result)

		myLocals[item] = theValue

		if strings.Contains(theValue, "${") {
			local = "\t" + item + " = \"" + theValue + "\" #" + value.(string) + "\n"
		} else {
			local = "\t" + item + " = " + theValue + " #" + value.(string) + "\n"
		}

		locals += strings.ReplaceAll(local, "'", "\"")
	}

	result["locals"] = myLocals

	return locals, result, nil
}
