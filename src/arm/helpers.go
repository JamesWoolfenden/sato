package arm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// IsLocal looks up whether a variable is a local or not
func IsLocal(target string, result map[string]interface{}) bool {
	locals := result["locals"].(map[string]interface{})
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

	myType := myItem["type"].(string)
	switch myType {
	case "object":
		{
			var types string
			var result string

			for name, item := range myItem["defaultValue"].(map[string]interface{}) {

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
			myItem["type"] = "number"
		}
	case "map[string]interface{}":
		{
			myItem["type"] = "string"
		}
	case "string", "securestring":
		{
			myItem["default"] = escapeQuote(myItem["default"])
		}
	case "list(string)", "number":
		{
			//do nothing
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
		return strings.Replace(item.(string), "\"", "\\\"", -1)
	}
	return ""
}

func arrayToString(defaultValue []interface{}) string {
	var newValue = "["
	for count, value := range defaultValue {
		if count == len(defaultValue)-1 {
			newValue = newValue + "\"" + value.(string) + "\""
		} else {
			newValue = newValue + "\"" + value.(string) + "\"" + ","
		}
	}
	return newValue + "]"
}
