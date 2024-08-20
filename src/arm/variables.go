package arm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sato/src/cf"
	"strings"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"
)

// ParseVariables convert ARM Parameters into terraform variables.
func ParseVariables(
	result map[string]interface{},
	funcMap tftemplate.FuncMap,
	destination string) (map[string]interface{}, error) {

	variables, ok := result["variables"].(map[string]interface{})

	if !ok {
		return result, &castError{"map[string]interface{}"}
	}

	var All string

	All, myVariables, err := parseParameters(result, funcMap, All)
	if err != nil {
		log.Warn().Msgf("parseParameters failed %s", err)
	}

	var locals string
	locals, result, err = ParseLocals(result)

	if err != nil {
		log.Warn().Msgf("locals failed %s", err)
	}

	if locals != "" {
		locals = "locals {\n" + locals + "}\n"
		err = cf.Write(locals, destination, "locals")

		if err != nil {
			return result, fmt.Errorf("failed to write locals %w", err)
		}
	}

	for name, value := range variables {
		myItem := make(map[string]interface{})

		if value != nil {
			var local string

			if reflect.TypeOf(value).String() == typeString {
				value, ok := value.(string)
				if !ok {
					log.Printf("type assertion failed")

					continue
				}

				if strings.Contains(value, "()") ||
					strings.Contains(value, "[") {
					value, result = ParseString(value, result)

					local = "\t" + name + " = " + value + "\n"
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

		myItem, err = FixType(myItem)
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
		return result, fmt.Errorf("variable write fail %w", err)
	}

	return result, nil
}
