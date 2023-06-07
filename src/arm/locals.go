package arm

import (
	"strings"
)

func parseLocals(result map[string]interface{}) (string, map[string]interface{}, error) {
	var locals string
	myLocals := result["locals"].(map[string]interface{})

	for x, value := range myLocals {
		var theValue string
		var local string
		theValue, result = parseString(value.(string), result)

		myLocals[x] = theValue
		if strings.Contains(theValue, "${") {
			local = "\t" + x + " = \"" + theValue + "\" #" + value.(string) + "\n"
		} else {
			local = "\t" + x + " = " + theValue + " #" + value.(string) + "\n"
		}

		locals += strings.ReplaceAll(local, "'", "\"")
	}

	result["locals"] = myLocals

	return locals, result, nil
}
