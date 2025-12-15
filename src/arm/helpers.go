package arm

import (
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
		return myItem, &castError{Type: "string"}
	}

	switch myType {
	case "object":
		{
			var typesBuilder strings.Builder
			var resultBuilder strings.Builder

			defaultValue, ok := myItem["defaultValue"].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("default value is not an object %v", myItem["defaultValue"])
			}

			for name, item := range defaultValue {
				switch item := item.(type) {
				case []interface{}:
					{
						var tempBuilder strings.Builder
						var tempTypesBuilder strings.Builder

						for _, y := range item {
							y, ok := y.(map[string]interface{})
							if !ok {
								return myItem, &castError{Type: "map[string]interface{}"}
							}

							for name, value := range y {
								tempBuilder.WriteString("\t   ")
								tempBuilder.WriteString(name)
								tempBuilder.WriteString(" = \"")
								tempBuilder.WriteString(value.(string))
								tempBuilder.WriteString("\"\n")

								tempTypesBuilder.WriteString("\t   ")
								tempTypesBuilder.WriteString(name)
								tempTypesBuilder.WriteString(" = string\n")
							}

							temp := "{\n" + strings.TrimSuffix(tempBuilder.String(), "\n") + "}"
							temptTypes := "{\n" + strings.TrimSuffix(tempTypesBuilder.String(), "\n") + "}"

							if resultBuilder.Len() > 0 {
								resultBuilder.WriteString(",")
							}
							resultBuilder.WriteString(name)
							resultBuilder.WriteString("= [")
							resultBuilder.WriteString(temp)
							resultBuilder.WriteString("]")

							if typesBuilder.Len() > 0 {
								typesBuilder.WriteString(",")
							}
							typesBuilder.WriteString(name)
							typesBuilder.WriteString("= list(object(")
							typesBuilder.WriteString(temptTypes)
							typesBuilder.WriteString("))")
						}
					}
				case map[string]interface{}:
					{
						log.Print(item)
					}
				case string:
					{
						if resultBuilder.Len() > 0 {
							resultBuilder.WriteString(",\n\t")
							typesBuilder.WriteString(",\n\t")
						}
						resultBuilder.WriteString(name)
						resultBuilder.WriteString(" = \"")
						resultBuilder.WriteString(EscapeQuote(item))
						resultBuilder.WriteString("\"")

						typesBuilder.WriteString(name)
						typesBuilder.WriteString(" = string")
					}
				case bool:
					{
						if resultBuilder.Len() > 0 {
							resultBuilder.WriteString(",\n\t")
							typesBuilder.WriteString(",\n\t")
						}
						resultBuilder.WriteString(name)
						resultBuilder.WriteString(" = ")
						resultBuilder.WriteString(strconv.FormatBool(item))

						typesBuilder.WriteString(name)
						typesBuilder.WriteString(" = bool")
					}
				case int:
					{
						if resultBuilder.Len() > 0 {
							resultBuilder.WriteString(",\n\t")
							typesBuilder.WriteString(",\n\t")
						}
						resultBuilder.WriteString(name)
						resultBuilder.WriteString(" = ")
						resultBuilder.WriteString(strconv.Itoa(item))

						typesBuilder.WriteString(name)
						typesBuilder.WriteString(" = number")
					}
				default:
					{
						log.Print(item)
					}
				}
			}

			myItem["default"] = "{\n\t" + resultBuilder.String() + "}"
			myItem["type"] = "object({\n\t" + typesBuilder.String() + "})"
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
			newValue += "\"" + EscapeQuote(value.(string)) + "\""
		} else {
			newValue += "\"" + EscapeQuote(value.(string)) + "\"" + ","
		}
	}

	return newValue + "]"
}

// Tags take a map into a string for tags.
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

// Ditch helps to drop functions for ARM.
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
