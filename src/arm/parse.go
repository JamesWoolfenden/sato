package arm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	sato "sato/src"
	"sato/src/see"
	"strings"
	tftemplate "text/template"
)

type m map[string]interface{}

// Parse turn CFN into Terraform
func Parse(file string, destination string) error {
	fileAbs, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	jsonFile, err := os.Open(fileAbs)

	if err != nil {
		fmt.Println(err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result map[string]interface{}
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		return err
	}

	funcMap := tftemplate.FuncMap{
		"Array":        sato.Array,
		"ArrayReplace": sato.ArrayReplace,
		"Contains":     sato.Contains,
		"Sprint":       sato.Sprint,
		"Decode64":     sato.Decode64,
		"Boolean":      sato.Boolean,
		"Dequote":      sato.Dequote,
		"Quote":        sato.Quote,
		"Demap":        sato.Demap,
		"ToUpper":      strings.ToUpper,
		"ToLower":      sato.Lower,
		"Deref":        func(str *string) string { return *str },
		"Nil":          sato.Nill,
		"Nild":         sato.Nild,
		"Marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"Split":        sato.Split,
		"Replace":      sato.Replace,
		"Tags":         sato.Tags,
		"RandomString": sato.RandomString,
		"Map":          sato.Map,
		"Snake":        sato.Snake,
		"Kebab":        sato.Kebab,
		"ZipFile":      sato.Zipfile,
	}

	_, err = ParseVariables(result, funcMap, destination)
	if err != nil {
		return err
	}

	Resources, err := ParseResources(result)
	log.Print(Resources)
	if err != nil {
		return err
	}

	return nil
}

// ParseVariables convert ARM Parameters into terraform variables
func ParseVariables(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) ([]sato.Variable, error) {
	parameters := result["parameters"].(map[string]interface{})
	variables := result["variables"].(map[string]interface{})

	var All string
	var locals string

	locals, All, myVariables, err := parseParameters(parameters, locals, funcMap, All)

	if err != nil {
		return myVariables, err
	}

	for name, value := range variables {
		var myVar sato.Variable

		switch value.(type) {
		case string:
			myVar.Type = "string"
			myVar.Name = name
			if value != nil {
				myVar.Default = value.(string)
			} else {
				myVar.Default = ""
			}
		case map[string]interface{}:
			myVar.Type = "string"
			myVar.Name = name
			myVar.Default = ""
		default:
			log.Print(value)
		}

		var output bytes.Buffer
		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(VariableFile))
		if err != nil {
			return nil, err
		}
		_ = tmpl.Execute(&output, m{
			"variable": myVar,
			"item":     name,
		})
		All = All + output.String()
		myVariables = append(myVariables, myVar)
	}

	err = sato.Write(All, destination, "variables")
	if err != nil {
		return nil, err
	}

	locals = "locals {\n" + locals + "}\n"
	err = sato.Write(locals, destination, "locals")
	if err != nil {
		return nil, err
	}

	return myVariables, nil
}

func parseParameters(parameters map[string]interface{}, locals string, funcMap tftemplate.FuncMap, All string) (string, string, []sato.Variable, error) {
	var myVariables []sato.Variable
	for name, item := range parameters {
		var myVariable sato.Variable
		myVariable.Name = name
		myItem := item.(map[string]interface{})

		if myItem["defaultValue"] != nil {
			var local string
			//contains function
			if strings.Contains(myItem["defaultValue"].(string), "()") {
				switch myItem["defaultValue"].(string) {
				case "[resourceGroup().location]":
					local = "\t" + name + " = " + "data.azurerm_resource_group.sato.id" + "\n"
				default:
					local = "\t" + name + " = " + myItem["defaultValue"].(string) + "\n"
				}
				locals = locals + local
				continue
			} else {
				myVariable.Default = myItem["defaultValue"].(string)
			}
		}
		myVariable.Type = myItem["type"].(string)
		myDescription := myItem["metadata"].(map[string]interface{})
		if myDescription != nil {
			myVariable.Description = myDescription["description"].(string)
		}

		var output bytes.Buffer
		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(VariableFile))
		if err != nil {
			return "", "", nil, err
		}
		_ = tmpl.Execute(&output, m{
			"variable": myVariable,
			"item":     name,
		})
		All = All + output.String()
		myVariables = append(myVariables, myVariable)
	}
	return locals, All, myVariables, nil
}

// ParseResources handles resources in ARM conversion
func ParseResources(result map[string]interface{}) (map[string]interface{}, error) {
	resources := result["resources"].([]interface{})
	newResources, err := parseList(resources)
	if err != nil {
		return nil, err
	}
	result["resources"] = newResources
	return result, nil
}

func parseList(resources []interface{}) ([]interface{}, error) {
	var newResources []interface{}
	matches := []string{"parameters", "variables"}
	for _, resource := range resources {
		myResource := resource.(map[string]interface{})
		myResource, err := parseMap(myResource, matches)
		if err != nil {
			return nil, err
		}
		newResources = append(newResources, myResource)
	}
	return newResources, nil
}

func parseMap(myResource map[string]interface{}, matches []string) (map[string]interface{}, error) {
	for name, attribute := range myResource {

		switch attribute.(type) {
		case string:
			myResource[name] = *parseString(matches, attribute.(string))
		case map[string]interface{}:
			var err error
			myResource[name], err = parseMap(attribute.(map[string]interface{}), matches)
			if err != nil {
				return nil, err
			}
		case []interface{}:
			myArray := attribute.([]interface{})
			for index, resource := range myArray {
				switch resource.(type) {
				case string:
					myArray[index] = parseString(matches, resource.(string))
				case map[string]interface{}:
					myArray[index], _ = parseMap(resource.(map[string]interface{}), matches)
				case []interface{}:
					log.Print(resource)
				}

			}
			myResource[name] = myArray
		case float64, bool:
			//it's fine

		default:
			log.Print(attribute)
		}
	}
	return myResource, nil
}

func parseString(matches []string, newAttribute string) *string {
	newAttribute, err := translate(newAttribute)
	if err != nil {
		return nil
	}
	if contains(matches, newAttribute) {
		if strings.Contains(newAttribute, "parameters") {
			newAttribute = strings.Replace(newAttribute, "parameters('", "var.", -1)
			newAttribute = strings.Replace(newAttribute, "')", "", -1)
		}
		if strings.Contains(newAttribute, "variables") {
			newAttribute = strings.Replace(newAttribute, "variables('", "var.", -1)
			newAttribute = strings.Replace(newAttribute, "')", "", -1)
		}
		newAttribute = strings.Replace(newAttribute, "[", "", 1)
		newAttribute = strings.Replace(newAttribute, "]", "", 1)
	}
	return &newAttribute
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(str, v) {
			return true
		}
	}

	return false
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

func handleResource(target string) (string, error) {
	if strings.Contains(target, "resourceId") {
		var re = regexp.MustCompile(`\[resourceId\((.*)\)\]`)
		matches := re.FindStringSubmatch(target)
		splitten := strings.Split(matches[1], ",")
		resourceName := splitten[1]
		var resourceMatch []string
		if strings.Contains(splitten[1], "parameters") {
			var myRe = regexp.MustCompile(`parameters\('(.*)'\)`)
			resourceMatch = myRe.FindStringSubmatch(resourceName)
		}

		if strings.Contains(splitten[1], "variables") {
			var myRe = regexp.MustCompile(`variables\('(.*)'\)`)
			resourceMatch = myRe.FindStringSubmatch(resourceName)
		}
		resource, err := see.Lookup(strings.Replace(splitten[0], "'", "", 2))
		if err != nil {
			return "", err
		}
		return *resource + "." + resourceMatch[1], nil
	}
	return "", nil
}
