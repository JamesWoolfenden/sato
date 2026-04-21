package arm

import (
	"fmt"
	"regexp"
	"sato/src/tfgen"
	"strings"

	"github.com/rs/zerolog/log"
)

// SplitResourceName splits and converts arm names into terraform.
func SplitResourceName(attribute string) (string, string, error) {
	const even = 2

	splitsy := strings.SplitN(attribute, ",", even)

	var (
		name string
		arm  string
	)

	switch len(splitsy) {
	case 1:
		{
			return "", "", &parseResourceError{attribute}
		}
	case 2:
		{
			re := regexp.MustCompile(`'(.*?)'`)
			newAttribute := re.FindStringSubmatch(splitsy[1])

			if len(newAttribute) <= 1 {
				arm = tfgen.Dequote(splitsy[0])
				// check it's not an array
				name = strings.TrimSpace(splitsy[1])
				if strings.Contains(name, ",") {
					name = strings.Split(name, ",")[0]
				}
			} else {
				name = newAttribute[1]
				newArm := re.FindStringSubmatch(splitsy[0])

				if len(newArm) > 1 {
					arm = newArm[1]
				} else {
					arm = splitsy[0]
				}
			}
		}
	}

	return arm, name, nil
}

// FindResourceName looks for resource names.
func FindResourceName(result map[string]interface{}, name string) (string, error) {
	if strings.HasPrefix(name, "format") {
		return name, &inlineFormatError{Name: name}
	}

	name = tfgen.Dequote(name)

	var err error

	if result["resources"] == nil {
		return "", &emptyResourceError{"resources"}
	}

	resources, ok := result["resources"].([]interface{})

	if !ok {
		return name, &emptyResourceError{"resources"}
	}

	for _, myResource := range resources {
		test, ok := myResource.(map[string]interface{})
		if !ok {
			log.Print("resource is not a map")

			continue
		}

		temp := LoseSQBrackets(test["name"].(string))

		if name == temp {
			return test["resource"].(string), nil
		}

		trimName := strings.Replace(name, "var.", "", 1)
		trimTemp := strings.Replace(Ditch(temp, "variables"), "'", "", 2)

		if trimTemp == trimName {
			return test["resource"].(string), nil
		}

		if strings.Contains(name, "local") {
			resourceName := strings.Split(name, ".")[1]
			if strings.Contains(temp, resourceName) {
				retrieved := test["resource"].(string)

				return retrieved, nil
			}
		}

		re := regexp.MustCompile(`\((.*?)\)`)
		splits := re.FindStringSubmatch(temp)

		if len(splits) > 1 {
			if trimName == strings.ReplaceAll(splits[1], "'", "") {
				return test["resource"].(string), nil
			}
		}
	}

	// not a simple name lookup
	if strings.Contains(name, ",") {
		Lots := strings.Split(name, ",")

		var newName []string

		for _, lot := range Lots {
			var part string
			if part, err = FindResourceName(result, strings.TrimSpace(lot)); err != nil {
				part, err = GetNameValue(result, strings.TrimSpace(lot))

				if err != nil {
					return "", err
				}
			}

			newName = append(newName, part)
		}

		return strings.Join(newName, "."), nil
	}

	name, err = GetNameValue(result, name)
	if err != nil {
		return "", fmt.Errorf("get Name value failed: %w", err)
	}

	return name, nil
}

// GetNameValue does just that.
func GetNameValue(result map[string]interface{}, name string) (string, error) {
	if strings.Contains(name, ".") {
		rawNames := strings.Split(name, ".")
		if len(rawNames) != 2 {
			return name, &matchValueError{name}
		}

		rawName := rawNames[1]

		if result["variables"] != nil {
			variables := result["variables"].(map[string]interface{})

			for myVariable, value := range variables {
				if rawName == myVariable {
					return value.(string), nil
				}
			}
		}
	}

	return name, nil
}

// FindResourceType get resource types.
func FindResourceType(result map[string]interface{}, name string) bool {
	if result["resources"] == nil {
		return false
	}

	resources, ok := result["resources"].([]interface{})
	if ok {
		for _, myResource := range resources {
			test, ok := myResource.(map[string]interface{})
			if ok {
				if name == test["type"].(string) {
					return true
				}
			}

		}
	}

	return false
}
