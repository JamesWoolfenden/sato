// Package arm provides functionality for parsing and converting Azure Resource Manager (ARM)
// templates to Terraform configuration files.
package arm

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sato/src/tfgen"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"
)

type m map[string]interface{}

var funcMap = tfgen.FuncMap(tftemplate.FuncMap{
	"Enabled": Enabled,
	"NotNil":  NotNil,
	"Set":     ArrayToString,
	"Tags":    Tags,
	"Uuid":    UUID,
})

// Parse turn ARM into Terraform.
func Parse(file string, destination string, opts ...Option) error {
	fileAbs, err := filepath.Abs(file)
	if err != nil {
		return &filepathError{Path: file}
	}

	jsonFile, err := os.Open(fileAbs) // #nosec G304 -- User-provided file path is expected for CLI tool
	if err != nil {
		return &openFileError{path: file, err: err}
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer func() {
		if closeErr := jsonFile.Close(); closeErr != nil {
			log.Warn().Err(closeErr).Msg("failed to close file")
		}
	}()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return &readFileError{path: file, err: err}
	}

	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)

	if err != nil {
		return &unmarshalError{err: err}
	}

	result = Preprocess(result)
	result, err = ParseVariables(result, funcMap, destination)

	if err != nil {
		return &parseVariablesError{Err: err}
	}

	result, err = ParseResources(result, funcMap, destination, opts...)
	if err != nil {
		return err
	}

	err = ParseOutputs(result, funcMap, destination)
	if err != nil {
		return &parseResourcesError{Err: err}
	}

	err = ParseData(result, funcMap, destination)
	if err != nil {
		return &parseDataError{err: err}
	}

	return nil
}

// ParseList parses List objects.
func ParseList(resources []interface{}, result map[string]interface{}) ([]interface{}, error) {
	var newResources []interface{}

	for _, resource := range resources {
		myResource, ok := resource.(map[string]interface{})
		if ok {
			myResource, err := ParseMap(myResource, result)
			if err != nil {
				return nil, &parseMapError{err}
			}

			newResources = append(newResources, myResource)
		} else {
			return nil, &parseListError{}
		}
	}

	return newResources, nil
}

// ParseMap parses a map type.
func ParseMap(myResource map[string]interface{}, result map[string]interface{}) (map[string]interface{}, error) {
	for name, attribute := range myResource {
		switch mapType := attribute.(type) {
		case string:
			var temp string
			temp, result = ParseString(mapType, result)
			myResource[name] = temp
		case map[string]interface{}:
			var err error
			myResource[name], err = ParseMap(mapType, result)

			if err != nil {
				return nil, err
			}
		case []interface{}:
			myArray := mapType
			for index, resource := range myArray {
				switch resource := resource.(type) {
				case string:
					var temp string
					temp, result = ParseString(resource, result)
					myArray[index] = temp
				case map[string]interface{}:
					myArray[index], _ = ParseMap(resource, result)
				case []interface{}:
					log.Print(resource)
				}
			}

			myResource[name] = myArray
		case float64, bool:
			// it's fine
			myResource[name] = attribute
		default:
			log.Print(attribute)
		}
	}

	return myResource, nil
}

// ParseString does just that.
func ParseString(attribute string, result map[string]interface{}) (string, map[string]interface{}) {
	matches := []string{
		"parameters", "variables", "toLower", "resourceGroup().location", "resourceGroup().id",
		"substring", "uniqueString", "reference", "resourceId", "listKeys", "format('", "SubscriptionResourceId",
		"concat", "subscription().tenantId", "uuid(", "uri(",
	}

	if what, found := Contains(matches, attribute); found {
		attribute, result = Replace(matches, attribute, what, result)
		attribute = LoseSQBrackets(attribute)
	}

	return attribute, result
}

// ensureDataMap ensures the data map exists and returns it.
func ensureDataMap(result map[string]interface{}) map[string]interface{} {
	if result["data"] == nil {
		return make(map[string]interface{})
	}

	data, ok := result["data"].(map[string]interface{})
	if !ok {
		log.Warn().Msgf("assertion failure %s", result["data"])
		return make(map[string]interface{})
	}

	return data
}
