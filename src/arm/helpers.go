package arm

import (
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
			myItem["type"] = "object()"
		}
	case "int", "float":
		{
			myItem["type"] = "number"
		}
	case "string", "map[string]interface{}":
		{
			myItem["type"] = "string"
		}
	case "securestring":
		{
			{
				myItem["type"] = "securestring"
			}
		}
	default:
		{
			log.Warn().Msgf("missed type %s", myType)
		}
	}
	return myItem
}
