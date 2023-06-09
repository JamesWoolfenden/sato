package arm

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log" //nolint:depguard
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

// Contains looks into slice for string.
func Contains(s []string, str string) (*string, bool) {
	for _, v := range s {
		if strings.Contains(strings.ToLower(str), strings.ToLower(v)) {
			return &v, true
		}
	}

	return nil, false
}

// FixType converts from go to terraform type names.
func FixType(myItem map[string]interface{}) (map[string]interface{}, error) {
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

			defaultValue, ok := myItem["defaultValue"].(map[string]interface{})
			if ok {
				for name, item := range defaultValue {
					switch item := item.(type) {
					case []interface{}:
						{
							var temp string
							var temptTypes string
							var myType string
							for _, y := range item {
								y, ok := y.(map[string]interface{})
								if !ok {
									return myItem, errors.New("failed to assert  (map[string]interface{}")
								}

								for name, value := range y {
									temp += "\t   " + name + " = \"" + value.(string) + "\"\n"
									myType += "\t   " + name + " = string\n"
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
								result = name + " = " + "\"" + EscapeQuote(item) + "\""
								types = name + " = " + "string"
							} else {
								result += ",\n\t" + name + " = " + "\"" + EscapeQuote(item) + "\""
								types += ",\n\t" + name + " = " + "string"
							}
						}
					case bool:
						{
							if result == "" {
								result = name + " = " + strconv.FormatBool(item)
								types = name + " = " + "bool"
							} else {
								temp := result
								result = temp + ",\n\t" + name + " = " + strconv.FormatBool(item)
								types = types + ",\n\t" + name + " = " + "bool"
							}
						}
					case int:
						{
							if result == "" {
								result = name + " = " + strconv.Itoa(item)
								types = name + " = " + "number"
							} else {
								temp := result
								result = temp + ",\n\t" + name + " = " + strconv.Itoa(item)
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
			myItem["default"] = EscapeQuote(myItem["default"])
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

// EscapeQuote them quotes.
func EscapeQuote(item interface{}) string {
	if item != nil {
		return strings.ReplaceAll(item.(string), "\"", "\\\"")
	}

	return ""
}

// ArrayToString squashes slice into string.
func ArrayToString(defaultValue []interface{}) string {
	newValue := "["

	for count, value := range defaultValue {
		if count == len(defaultValue)-1 {
			newValue += "\"" + value.(string) + "\""
		} else {
			newValue += "\"" + value.(string) + "\"" + ","
		}
	}

	return newValue + "]"
}

// Tags take map into a string for tags
func Tags(tags map[string]interface{}) string {
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

// NotNil handles interfaces that are bools and their value being set
func NotNil(unknown interface{}) bool {
	return unknown != nil
}

// Enabled cast from string to bool
func Enabled(status string) bool {
	return strings.ToLower(status) == "Enabled"
}

// LoseSQBrackets ditches square brackets
func LoseSQBrackets(newAttribute string) string {
	re := regexp.MustCompile(`^\[(.*)\]`) // format('{0}/{1}',
	Matched := re.FindStringSubmatch(newAttribute)

	if len(Matched) > 1 {
		return Matched[1]
	}

	return newAttribute
}

// Ditch helps to drop functions for arm
func Ditch(attribute string, ditch string) string {

	leftBrackets := strings.SplitAfter(attribute, "(")

	if len(leftBrackets) == 0 {
		return attribute
	}

	var brackets []string

	for _, item := range leftBrackets {
		rbrackets := strings.SplitAfter(item, ")")
		brackets = append(brackets, rbrackets...)
	}

	brackets = deleteEmpty(brackets)

	var opposite = 100

	var raw []string

	for x, item := range brackets {
		if strings.Contains(item, ditch) {
			opposite = len(brackets) - 2 + x
			item = strings.Replace(item, ditch+"(", "", 1)
		}

		if opposite != x {
			if !strings.Contains(item, ditch) {
				raw = append(raw, item)
			} else {
				raw = append(raw, strings.Replace(item, ")", "", 1))
			}
		}
	}

	return strings.Join(raw, "")
}

// UUID replaces.
func UUID(count int) string {
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
