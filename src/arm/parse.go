package arm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	sato "sato/src"
	"sato/src/see"
	"strconv"
	"strings"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"

	"golang.org/x/exp/maps"
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
		"SplitOn":      sato.SplitOn,
		"Replace":      sato.Replace,
		"Tags":         sato.Tags,
		"RandomString": sato.RandomString,
		"Map":          sato.Map,
		"Snake":        sato.Snake,
		"Kebab":        sato.Kebab,
		"ZipFile":      sato.Zipfile,
	}

	result = preprocess(result)
	variables, err := ParseVariables(result, funcMap, destination)
	if err != nil {
		return err
	}

	result, err = ParseResources(result, funcMap, destination, variables)
	if err != nil {
		return err
	}
	err = ParseOutputs(result, funcMap, destination)

	if err != nil {
		return err
	}

	return nil
}

// ParseVariables convert ARM Parameters into terraform variables
func ParseVariables(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) ([]sato.Variable, error) {

	variables := result["variables"].(map[string]interface{})

	var All string
	var locals string

	locals, All, myVariables, err := parseParameters(result, locals, funcMap, All)

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
				if strings.Contains(value.(string), "()") ||
					strings.Contains(value.(string), "[") {

					local := "\t" + name + " = " + *parseString(value.(string), result) + "\n"
					locals = locals + local
					continue
				}
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

func parseParameters(result map[string]interface{}, locals string, funcMap tftemplate.FuncMap, All string) (string, string, []sato.Variable, error) {
	parameters := result["parameters"].(map[string]interface{})
	var myVariables []sato.Variable
	for name, item := range parameters {
		var myVariable sato.Variable
		myVariable.Name = name
		myItem := item.(map[string]interface{})

		if myItem["defaultValue"] != nil {
			var local string
			//contains function
			if strings.Contains(myItem["defaultValue"].(string), "()") ||
				strings.Contains(myItem["defaultValue"].(string), "[") {

				local = "\t" + name + " = " + *parseString(myItem["defaultValue"].(string), result) + "\n"
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
func ParseResources(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string, variables []sato.Variable) (map[string]interface{}, error) {
	resources := result["resources"].([]interface{})
	newResources, err := parseList(resources, result)
	if err != nil {
		return nil, err
	}
	result["resources"] = newResources

	for _, resource := range newResources {
		var output bytes.Buffer
		var name *string
		myType := resource.(map[string]interface{})
		myContent := lookup(myType["type"].(string))
		//item := strings.Replace(myType["name"].(string), "var.", "", 1)
		first, err := see.Lookup(myType["type"].(string))
		if err != nil {
			return nil, err
		}

		//name, err = GetValue(item, variables)
		//if err != nil || name == nil || *name == "" {
		//
		temp := myType["resource"].(string)
		name = &temp

		//needs to pivot on policy template from resource
		tmpl, err := tftemplate.New("sato").Funcs(funcMap).Parse(string(myContent))

		if err != nil {
			log.Printf("failed at %s  for %s %s", err, *first, *name)
			continue
		}

		_ = tmpl.Execute(&output, sato.M{
			"resource": resource,
			"item":     name,
		})

		err = sato.Write(output.String(), destination, *first+"."+strings.Replace(*name, "var.", "", 1))
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func parseList(resources []interface{}, result map[string]interface{}) ([]interface{}, error) {
	var newResources []interface{}
	for _, resource := range resources {
		myResource := resource.(map[string]interface{})
		myResource, err := parseMap(myResource, result)
		if err != nil {
			return nil, err
		}
		newResources = append(newResources, myResource)
	}
	return newResources, nil
}

func parseMap(myResource map[string]interface{}, result map[string]interface{}) (map[string]interface{}, error) {
	for name, attribute := range myResource {

		switch attribute.(type) {
		case string:
			myResource[name] = *parseString(attribute.(string), result)
		case map[string]interface{}:
			var err error
			myResource[name], err = parseMap(attribute.(map[string]interface{}), result)
			if err != nil {
				return nil, err
			}
		case []interface{}:
			myArray := attribute.([]interface{})
			for index, resource := range myArray {
				switch resource.(type) {
				case string:
					myArray[index] = parseString(resource.(string), result)
				case map[string]interface{}:
					myArray[index], _ = parseMap(resource.(map[string]interface{}), result)
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

func parseString(newAttribute string, result map[string]interface{}) *string {
	var matches = []string{"parameters", "variables", "toLower", "resourceGroup().location", "resourceGroup().id",
		"substring", "format('", "uniqueString", "reference", "resourceId"}

	if what, found := contains(matches, newAttribute); found {
		newAttribute = replace(matches, newAttribute, what, result)
		newAttribute = loseSQBrackets(newAttribute)
	}
	return &newAttribute
}

func loseSQBrackets(newAttribute string) string {
	newAttribute = strings.Replace(newAttribute, "[", "", 1)
	newAttribute = strings.Replace(newAttribute, "]", "", 1)
	return newAttribute
}

func replace(matches []string, newAttribute string, what *string, result map[string]interface{}) string {
	var Attribute string
	switch *what {
	case "reference", "resourceId":
		{
			if *what == "reference" {
				var re = regexp.MustCompile(`reference\((.*?)\)`) //format('{0}/{1}',
				Match := re.FindStringSubmatch(newAttribute)
				if (Match) != nil {
					splitten := strings.Split(newAttribute, "')")
					if len(splitten) > 1 {
						remainder := strings.Replace(splitten[len(splitten)-1], "]", "", 1)
						name, resourceName := replaceResourceID(Match[1], result)
						Attribute = *resourceName + "." + name + remainder
					} else {
						log.Printf("no match found %s", splitten[0])
					}
				} else {
					log.Printf("no match found %s", newAttribute)
				}
			} else {

				name, resourceName := replaceResourceID(loseSQBrackets(newAttribute), result)
				Attribute = *resourceName + "." + name
			}

		}
	case "uniqueString":
		{
			var re = regexp.MustCompile(`uniqueString\((.*?)\)`) //format('{0}/{1}',
			Match := re.ReplaceAllString(newAttribute, "substr(uuid(), 0, 8)")
			Attribute = Match
		}
	case "format('":
		{
			var re = regexp.MustCompile(`{.}`) //format('{0}/{1}',
			Match := re.ReplaceAllString(newAttribute, "%%s")
			Match = strings.Replace(Match, "'", "\"", -1)
			Attribute = Match
		}
	case "parameters":
		{
			var re = regexp.MustCompile(`parameters\('(.*?)\'\)`)
			Match := re.FindStringSubmatch(newAttribute)
			if (Match) != nil {
				var temp string
				if IsLocal(Match[1], result) {
					temp = "local." + Match[1]
				} else {
					temp = "var." + Match[1]
				}
				Attribute = strings.Replace(newAttribute, Match[0], temp, -1)
			} else {
				log.Printf("no match found %s", newAttribute)
			}
		}
	case "variables":
		{
			var re = regexp.MustCompile(`variables\('(.*?)\'\)`)
			Match := re.FindStringSubmatch(newAttribute)
			if (Match) != nil {
				var temp string
				if IsLocal(Match[1], result) {
					temp = "local." + Match[1]
				} else {
					temp = "var." + Match[1]
				}
				Attribute = strings.Replace(newAttribute, Match[0], temp, -1)
			} else {
				log.Printf("not found %s", newAttribute)
			}
		}
	case "toLower":
		{
			Attribute = strings.Replace(newAttribute, "toLower", "lower", -1)
		}
	case "resourceGroup().location":
		{
			Attribute = strings.Replace(newAttribute, "resourceGroup().location",
				"data.azurerm_resource_group.sato.location", -1)
		}
	case "resourceGroup().id":
		{
			Attribute = strings.Replace(newAttribute, "resourceGroup().id",
				"data.azurerm_resource_group.sato.id", -1)
		}
	case "substring":
		{
			Attribute = strings.Replace(newAttribute, "substring",
				"substr", -1)
			Attribute = strings.Replace(Attribute, "'", "\"", -1)
		}
	}
	if again, still := contains(matches, Attribute); still {
		Attribute = replace(matches, Attribute, again, result)
	}
	return Attribute
}

func replaceResourceID(Match string, result map[string]interface{}) (string, *string) {
	resource := strings.Replace(Match, "resourceId(", "", 1)
	splitter := strings.Split(resource, ",")
	arm := strings.Replace(splitter[0], "'", "", 2)
	name := strings.Replace(strings.Replace(splitter[1], ")", "", 1), " ", "", -1)
	name, err := findResourceName(result, name)
	if err != nil {
		log.Print(err)
	}
	resourceName, err := see.Lookup(arm)
	if err != nil {
		log.Printf("no match found %s", arm)
	}
	return name, resourceName
}

func findResourceName(result map[string]interface{}, name string) (string, error) {
	resources := result["resources"].([]interface{})
	for _, myResource := range resources {
		test := myResource.(map[string]interface{})
		if name == test["name"].(string) {
			return test["resource"].(string), nil
		}
	}
	return "", errors.New("failed to find {name}")
}

func handleResource(target string) (string, error) {
	var attribute string
	if strings.Contains(target, "[") {
		var re = regexp.MustCompile(`\[(.*)\]`)
		Match := re.FindStringSubmatch(target)
		target = Match[1]
	}

	if strings.Contains(target, "reference") {
		var re = regexp.MustCompile(`reference\((.*)\)`)
		Match := re.ReplaceAllString(target, "")
		attribute = Match
	}

	if strings.Contains(target, "resourceId") {
		var re = regexp.MustCompile(`resourceId\((.*)\)`)
		myMatches := re.FindStringSubmatch(target)
		if myMatches == nil {
			return target, errors.New("[resourceId] not found")
		}
		splitten := strings.Split(myMatches[1], ",")
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
		return *resource + "." + resourceMatch[1] + attribute, nil
	}
	return target, nil
}

// ParseOutputs writes out to outputs.tf
func ParseOutputs(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) error {
	outputs := result["outputs"].(map[string]interface{})

	var All string
	for name, value := range outputs {
		var myVar sato.Output
		myVar.Type = "string"
		myVar.Name = name
		temp := value.(map[string]interface{})
		myVar.Value, _ = handleResource(temp["value"].(string))

		var output bytes.Buffer
		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(OutputFile))
		if err != nil {
			return err
		}
		_ = tmpl.Execute(&output, m{
			"variable": myVar,
			"item":     name,
		})
		All = All + output.String()
	}

	err := sato.Write(All, destination, "outputs")
	if err != nil {
		return err
	}

	return nil
}

// GetValue gets from variables
func GetValue(item string, variables []sato.Variable) (*string, error) {
	for _, x := range variables {
		if x.Name == item {
			return &x.Default, nil
		}
	}
	something := fmt.Errorf("%s Not found", item)
	return nil, something
}

func preprocess(results map[string]interface{}) map[string]interface{} {
	resources := results["resources"].([]interface{})
	var newResults []interface{}
	for item, result := range resources {
		inside := result.(map[string]interface{})
		counter := map[string]interface{}{"resource": fmt.Sprintf("sato" + strconv.Itoa(item))}
		maps.Copy(inside, counter)
		newResults = append(newResults, inside)
	}
	results["resources"] = newResults

	paraVariables := results["variables"].(map[string]interface{})
	locals := make(map[string]interface{})
	for item, result := range paraVariables {
		switch result.(type) {
		case string:
			{
				if strings.Contains(result.(string), "[") {
					locals[item] = result
				}
			}
		default:
			jasoned, _ := json.Marshal(result)
			if strings.Contains(string(jasoned), "[") {
				locals[item] = string(jasoned)
			}

		}

		//needs a data principal for data sources assumed
		//need to find
	}

	paraParameters := results["parameters"].(map[string]interface{})
	newLocals := make(map[string]interface{})
	for item, result := range paraParameters {
		myResult := result.(map[string]interface{})
		_, ok := myResult["defaultValue"]
		if ok {
			defaultValue := myResult["defaultValue"].(string)
			if strings.Contains(defaultValue, "[") {
				newLocals[item] = defaultValue
			}
		} else {
			continue
		}

	}
	maps.Copy(locals, newLocals)
	results["locals"] = locals
	return results
}
