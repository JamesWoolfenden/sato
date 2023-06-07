package arm

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sato/src/cf"
	"strings"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"
)

// parseVariables convert ARM Parameters into terraform variables.
func parseVariables(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) (map[string]interface{}, error) {
	variables := make(map[string]interface{})
	if result["variables"] != nil {
		variables = result["variables"].(map[string]interface{})
	}

	var All string

	All, myVariables, err := parseParameters(result, funcMap, All)
	if err != nil {
		return result, err
	}

	var locals string
	locals, result, err = parseLocals(result)

	if err != nil {
		return result, err
	}

	for name, value := range variables {
		myItem := make(map[string]interface{})

		if value != nil {
			var local string

			if reflect.TypeOf(value).String() == typeString {
				if strings.Contains(value.(string), "()") ||
					strings.Contains(value.(string), "[") {
					value, result = parseString(value.(string), result)

					local = "\t" + name + " = " + value.(string) + "\n"
					locals += local
					continue
				}

				myItem["default"] = value
			}

			if reflect.TypeOf(value).String() == "map[string]interface {}" {
				blob, err := json.Marshal(value)
				if err != nil {
					log.Warn().Msgf("fail to marshal %s", value)
				}

				myItem["default"] = string(blob)
			}

			myItem["name"] = name
			myItem["type"] = typeString
		}

		myItem, err = fixType(myItem)
		if err != nil {
			log.Print(err)
		}

		var output bytes.Buffer

		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(variableFile))

		if err != nil {
			return result, err
		}
		_ = tmpl.Execute(&output, m{
			"variable": myItem,
			"item":     name,
		})
		All += output.String()

		myVariables = append(myVariables, myItem)
	}

	err = cf.Write(All, destination, "variables")
	if err != nil {
		return result, err
	}

	locals = "locals {\n" + locals + "}\n"
	err = cf.Write(locals, destination, "locals")

	if err != nil {
		return result, err
	}

	return result, nil
}
