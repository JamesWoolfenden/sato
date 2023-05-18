package arm

import (
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
		if strings.Contains(str, v) {
			return &v, true
		}
	}

	return nil, false
}

func translate(target string) (string, error) {
	if strings.Contains(target, "reference") {
		log.Printf("skipped %s", target)
	} else {
		s, err := handleResource(target)
		if err != nil {
			return s, err
		}
		return s, nil
	}

	return target, nil
}

func fixType(myItem map[string]interface{}) map[string]interface{} {
	if myItem["type"] == nil {
		return nil
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
						log.Print(item)
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
	case "list(string)", "string", "number", "securestring":
		{
			//do nothing- should move this all to preprocess
		}
	default:
		{
			log.Warn().Msgf("missed type %s", myType)
		}
	}
	return myItem
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
