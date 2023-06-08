package arm

import (
	"errors"
	"fmt"
	"regexp"
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

			_, ok := myItem["defaultValue"].(map[string]interface{})
			if ok {
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
								result += "," + name + "= [" + temp + "]"
								types += "," + name + "= list(object(" + temptTypes + "))"
							} else {
								result += name + "= [" + temp + "]"
								types += name + "= list(object(" + temptTypes + "))"
							}
						}
					case map[string]interface{}:
						{
							log.Print(item)
						}
					case string:
						{
							if result == "" {
								result = name + " = " + "\"" + escapeQuote(item) + "\""
								types = name + " = " + "string"
							} else {
								result += ",\n\t" + name + " = " + "\"" + escapeQuote(item) + "\""
								types += ",\n\t" + name + " = " + "string"
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
					case int:
						{
							if result == "" {
								result = name + " = " + strconv.Itoa(item.(int))
								types = name + " = " + "number"
							} else {
								temp := result
								result = temp + ",\n\t" + name + " = " + strconv.Itoa(item.(int))
								types = types + ",\n\t" + name + " = " + "number"
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
			} else {
				return nil, fmt.Errorf("default value is not an object %v", myItem["defaultValue"])
			}
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
	case typeListString, typeNumber, "bool":
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
		if _, ok := name.(string); ok {
			tagged += "\t\"" + item + "\"" + " = " + "\"" + name.(string) + "\"\n"
		} else {
			tagged += "\t\"" + item + "\"" + " = " + "\"OBJECT\"\n"
		}
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

func enabled(status string) bool {
	if strings.ToLower(status) == "enabled" {
		return true
	}
	return false
}

func loseSQBrackets(newAttribute string) string {
	re := regexp.MustCompile(`^\[(.*)\]`) // format('{0}/{1}',
	Matched := re.FindStringSubmatch(newAttribute)
	if len(Matched) > 1 {
		return Matched[1]
	}
	return newAttribute
}

func ditch(Attribute string, name string) string {

	leftBrackets := strings.SplitAfter(Attribute, "(")

	if len(leftBrackets) == 0 {
		return Attribute
	}

	var brackets []string

	for _, item := range leftBrackets {
		rbrackets := strings.SplitAfter(item, ")")
		brackets = append(brackets, rbrackets...)
	}

	brackets = deleteEmpty(brackets)

	y := 100
	var raw []string
	for x, item := range brackets {
		if strings.Contains(item, name) {
			y = len(brackets) - 2 + x
			raw = append(raw, strings.Replace(item, name+"(", "", 1))
		}

		if y != x && !strings.Contains(item, name) {
			raw = append(raw, item)
		}

		if y == x {
			raw = append(raw, strings.Replace(item, ")", "", 1))
		}
	}
	return strings.Join(raw, "")
}

func uuid(count int) string {
	var i int
	var uuids string
	for i = 0; i < (count); i++ {
		uuids += "resource \"random_uuid\" \"sato" + strconv.Itoa(i) + "\" {}\n"
	}
	return uuids
}

func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
