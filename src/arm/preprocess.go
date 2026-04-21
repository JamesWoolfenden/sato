package arm

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sato/src/tfgen"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

// processParameterType processes a single parameter based on its type.
func processParameterType(item string, myResult map[string]interface{}, newLocals, newParams map[string]interface{}) {
	_, ok := myResult["defaultValue"]

	if !ok {
		myResult["default"] = ""
		newParams[item] = myResult
		return
	}

	myType := myResult["type"].(string)
	switch strings.ToLower(myType) {
	case "string", "securestring":
		defaultValue := myResult["defaultValue"].(string)
		if strings.Contains(defaultValue, "[") {
			newLocals[item] = defaultValue
			return
		}

		myResult["default"] = myResult["defaultValue"].(string)
		newParams[item] = myResult
	case "int":
		myResult["type"] = "number"
		myResult["default"] = fmt.Sprintf("%v", myResult["defaultValue"].(float64))
		newParams[item] = myResult
	case "object", "list(string)":
		// todo
		// myResult["default"] = myResult["defaultValue"]
		myResult["default"] = ""
		newParams[item] = myResult
	case "array":
		myResult["type"] = "list(string)"
		myResult["default"] = ArrayToString(myResult["defaultValue"].([]interface{}))
		newParams[item] = myResult
	case "map[string]interface{}":
		log.Debug().Msgf("handled %s", myType)
	case "bool":
		myResult["default"] = fmt.Sprintf("%v", myResult["defaultValue"])
		newParams[item] = myResult
	default:
		log.Warn().Msgf("unhandled type %s", myType)
	}
}

// Preprocess examines raw ARM loads.
func Preprocess(results map[string]interface{}) map[string]interface{} {
	results["resources"] = SetResourceNames(results)
	locals := make(map[string]interface{})

	// only satisfied if empty
	_, ok := results["variables"].(map[string]interface{})

	if !ok {
		paraVariables := results["variables"].(map[string]interface{})

		newVariables := make(map[string]interface{})

		for item, result := range paraVariables {
			switch result := result.(type) {
			case string:
				if strings.Contains(result, "[") {
					locals[item] = result
				} else {
					newVariables[item] = result
				}
			default:
				jasoned, _ := json.Marshal(result)
				if strings.Contains(string(jasoned), "[") {
					locals[item] = string(jasoned)
				} else {
					newVariables[item] = result
				}
			}
		}

		results["variables"] = newVariables
	}

	paraParameters := results["parameters"].(map[string]interface{})

	newLocals := make(map[string]interface{})
	newParams := make(map[string]interface{})

	for item, result := range paraParameters {
		myResult := result.(map[string]interface{})
		processParameterType(item, myResult, newLocals, newParams)
	}

	results["parameters"] = newParams

	maps.Copy(locals, newLocals)
	results["locals"] = locals

	return results
}

// generateResourceName creates a meaningful resource name from type and resource name.
func generateResourceName(resourceType string, resourceName string, index int) string {
	// Extract the last part of the resource type
	// e.g., "Microsoft.Network/virtualNetworks" -> "virtualNetworks"
	parts := strings.Split(resourceType, "/")
	typeName := parts[len(parts)-1]

	// Convert to snake case
	typeName = tfgen.Snake(typeName)

	// Clean up the resource name if it's a simple string
	cleanName := strings.ToLower(resourceName)
	cleanName = strings.ReplaceAll(cleanName, "-", "_")
	cleanName = strings.ReplaceAll(cleanName, " ", "_")
	cleanName = regexp.MustCompile(`[^\w_]`).ReplaceAllString(cleanName, "")

	// If the name contains template expressions like [parameters or [variables, just use type with index
	if strings.Contains(resourceName, "[") || cleanName == "" {
		return fmt.Sprintf("%s_%d", typeName, index)
	}

	// Truncate if too long
	const maxLength = 50
	combined := fmt.Sprintf("%s_%s", typeName, cleanName)
	if len(combined) > maxLength {
		combined = combined[:maxLength]
	}

	return combined
}

// SetResourceNames gets resource names for results.
func SetResourceNames(results map[string]interface{}) []interface{} {
	resources := results["resources"].([]interface{})

	newResults := make([]interface{}, 0, len(resources))
	nameCounter := make(map[string]int)

	for item, result := range resources {
		inside := result.(map[string]interface{})

		// Get resource type and name
		resourceType, typeOk := inside["type"].(string)
		resourceName, nameOk := inside["name"].(string)

		var generatedName string
		if typeOk && nameOk {
			baseName := generateResourceName(resourceType, resourceName, item)

			// Handle duplicate names by adding a suffix
			if count, exists := nameCounter[baseName]; exists {
				nameCounter[baseName] = count + 1
				generatedName = fmt.Sprintf("%s_%d", baseName, count+1)
			} else {
				nameCounter[baseName] = 0
				generatedName = baseName
			}
		} else {
			// Fallback to old naming scheme if type or name is missing
			generatedName = fmt.Sprintf("sato%d", item)
		}

		counter := map[string]interface{}{"resource": generatedName}
		maps.Copy(inside, counter)
		newResults = append(newResults, inside)
	}

	return newResults
}
