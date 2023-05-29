package arm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

const typeNumber = "number"
const typeString = "string"
const typeListString = "list(string)"

// IsLocal looks up whether a variable is a local or not.
func IsLocal(target string, result map[string]interface{}) bool {
	locals, ok := result["locals"].(map[string]interface{})

	if !ok {
		return false
	}

	for name := range locals {
		if name == target {
			return true
		}
	}

	return false
}

func contains(s []string, str string) (*string, bool) {
	for _, v := range s {
		if strings.Contains(strings.ToLower(str), strings.ToLower(v)) {
			return &v, true
		}
	}

	return nil, false
}

func fixType(myItem map[string]interface{}) (map[string]interface{}, error) {
	if myItem["type"] == nil {
		return myItem, fmt.Errorf("object type is nil %s", myItem)
	}

	myType, ok := myItem["type"].(string)

	if !ok {
		return myItem, errors.New("type is not string")
	}

	switch myType {
	case "object":
		{
			var types string
			var result string

			for name, item := range myItem["defaultValue"].(map[string]interface{}) {
				//goland:noinspection GoLinter
				switch item.(type) {
				case []interface{}:
					{
						var temp string
						var temptTypes string
						var myType string
						for _, y := range item.([]interface{}) {
							for name, value := range y.(map[string]interface{}) {
								temp = temp + "\t   " + name + " = \"" + value.(string) + "\"\n"
								myType = myType + "\t   " + name + " = string\n"
							}
							temp = "{\n" + strings.TrimSuffix(temp, "\n") + "}"
							temptTypes = "{\n" + strings.TrimSuffix(myType, "\n") + "}"
						}
						if result != "" {
							result = result + "," + name + "= [" + temp + "]"
							types = types + "," + name + "= list(object(" + temptTypes + "))"
						} else {
							result = result + name + "= [" + temp + "]"
							types = types + name + "= list(object(" + temptTypes + "))"
						}
					}
				case map[string]interface{}:
					{
						log.Print(item)
					}
				case string:
					{
						if result == "" {
							result = name + " = " + escapeQuote(item)
							types = name + " = " + "string"
						} else {
							temp := result
							result = temp + ",\n\t" + name + " = " + escapeQuote(item)
							types = types + ",\n\t" + name + " = " + "string"
						}
					}
				case bool:
					{
						if result == "" {
							result = name + " = " + strconv.FormatBool(item.(bool))
							types = name + " = " + "bool"
						} else {
							temp := result
							result = temp + ",\n\t" + name + " = " + strconv.FormatBool(item.(bool))
							types = types + ",\n\t" + name + " = " + "bool"
						}
					}
				default:
					{
						log.Print(item)
					}
				}
			}
			myItem["default"] = "{\n\t" + result + "}"
			myItem["type"] = "object({\n\t" + types + "})"
		}
	case "int", "float":
		{
			myItem["type"] = typeNumber
		}
	case "map[string]interface{}":
		{
			myItem["type"] = typeString
		}
	case typeString, "securestring":
		{
			myItem["default"] = escapeQuote(myItem["default"])
		}
	case typeListString, typeNumber:
		{
			// do nothing
		}
	default:
		{
			log.Warn().Msgf("missed type %s", myType)
		}
	}

	return myItem, nil
}

func escapeQuote(item interface{}) string {
	if item != nil {
		return strings.ReplaceAll(item.(string), "\"", "\\\"")
	}

	return ""
}

func arrayToString(defaultValue []interface{}) string {
	newValue := "["

	for count, value := range defaultValue {
		if count == len(defaultValue)-1 {
			newValue = newValue + "\"" + value.(string) + "\""
		} else {
			newValue = newValue + "\"" + value.(string) + "\"" + ","
		}
	}

	return newValue + "]"
}

func tags(tags map[string]interface{}) string {
	tagged := "{\n"
	for item, name := range tags {
		tagged += "\t\"" + item + "\"" + " = " + "\"" + name.(string) + "\"\n"
	}

	tagged += "\t}"

	return tagged
}

func notNil(unknown interface{}) bool {
	if unknown == nil {

		return false
	}

	return true
}
