package arm

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log" //nolint:depguard
)

const (
	typeNumber     = "number"
	typeString     = "string"
	typeListString = "list(string)"
)

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
			if !ok {
				return nil, fmt.Errorf("default value is not an object %v", myItem["defaultValue"])
			}

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

// Tags take map into a string for tags.
func Tags(raw interface{}) string {
	tags, ok := raw.(map[string]interface{})
	if !ok {
		return raw.(string)
	}

	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	tagged := "{\n"

	for _, k := range keys {
		if value, ok := tags[k].(string); ok {
			tagged += "\t\"" + k + "\"" + " = " + "\"" + value + "\"\n"
		} else {
			tagged += "\t\"" + k + "\"" + " = " + "\"OBJECT\"\n"
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
	return strings.ToLower(status) == "enabled"
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

// Ditch helps to drop functions for arm.
func Ditch(attribute string, ditch string) string {
	splitter := strings.SplitN(attribute, ditch+"(", 2)

	if len(splitter) != 2 {
		return attribute
	}

	remains := strings.Replace(splitter[0], ditch+"(", "", 1)
	next := splitter[1]

	cards := 1

	var newString string

	leftBrackets := "("
	rightBrackets := ")"

	found := false

	for _, y := range next {
		char := string(y)
		switch char {
		case leftBrackets:
			{
				if !found {
					cards++
				}

				newString += char
			}
		case rightBrackets:
			{
				if !found {
					if cards != 1 {
						newString += string(y)
					} else {
						found = true
					}

					cards--
				} else {
					newString += string(y)
				}
			}
		default:
			{
				newString += char
			}
		}
	}

	newString = remains + newString

	return newString
}

// UUID replaces.
func UUID(count int) string {
	var (
		i     int
		uuids string
	)

	for i = 0; i < (count); i++ {
		uuids += "resource \"random_uuid\" \"sato" + strconv.Itoa(i) + "\" {}\n"
	}

	return uuids
}
