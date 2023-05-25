package arm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	sato "sato/src"
	"sato/src/see"
	"strconv"
	"strings"
	tftemplate "text/template"
	"unicode"

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
	//goland:noinspection GoUnhandledErrorResult
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
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
	result, err = ParseVariables(result, funcMap, destination)
	if err != nil {
		return err
	}

	result, err = ParseResources(result, funcMap, destination)
	if err != nil {
		return err
	}

	err = ParseOutputs(result, funcMap, destination)
	if err != nil {
		return err
	}

	err = ParseData(result, funcMap, destination)
	if err != nil {
		return err
	}
	return nil
}

// ParseVariables convert ARM Parameters into terraform variables
func ParseVariables(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) (map[string]interface{}, error) {

	variables := make(map[string]interface{})
	if result["variables"] != nil {
		variables = result["variables"].(map[string]interface{})
	}

	var All string

	All, myVariables, err := parseParameters(result, funcMap, All)

	if err != nil {
		return result, err
	}

	locals, result, err := parseLocals(result)

	if err != nil {
		return result, err
	}

	for name, value := range variables {
		myItem := make(map[string]interface{})

		if value != nil {
			var local string

			if reflect.TypeOf(value).String() == "string" {
				if strings.Contains(value.(string), "()") ||
					strings.Contains(value.(string), "[") {
					value, result = parseString(value.(string), result)
					local = "\t" + name + " = " + value.(string) + "\n"
					locals = locals + local
					continue
				}
				myItem["default"] = value
			}

			if reflect.TypeOf(value).String() == "map[string]interface {}" {
				blob, _ := json.Marshal(value)
				myItem["default"] = string(blob)
			}

			myItem["name"] = name
			myItem["type"] = "string"
		}

		myItem = fixType(myItem)

		var output bytes.Buffer
		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(VariableFile))
		if err != nil {
			return result, err
		}
		_ = tmpl.Execute(&output, m{
			"variable": myItem,
			"item":     name,
		})
		All = All + output.String()
		myVariables = append(myVariables, myItem)
	}

	err = sato.Write(All, destination, "variables")
	if err != nil {
		return result, err
	}

	locals = "locals {\n" + locals + "}\n"
	err = sato.Write(locals, destination, "locals")
	if err != nil {
		return result, err
	}

	return result, nil
}

func parseParameters(result map[string]interface{}, funcMap tftemplate.FuncMap, All string) (string, []interface{}, error) {
	parameters := result["parameters"].(map[string]interface{})
	myVariables := make([]interface{}, 0)

	for name, item := range parameters {

		myItem := item.(map[string]interface{})

		myItem = fixType(myItem)

		var output bytes.Buffer
		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(VariableFile))
		if err != nil {
			return "", nil, err
		}
		_ = tmpl.Execute(&output, m{
			"variable": myItem,
			"item":     name,
		})
		All = All + output.String()
		myVariables = append(myVariables, myItem)
	}
	return All, myVariables, nil
}

func parseLocals(result map[string]interface{}) (string, map[string]interface{}, error) {
	var locals string
	myLocals := result["locals"].(map[string]interface{})

	for x, y := range myLocals {
		var temp *string
		temp, result = parseString(y.(string), result)
		local := "\t" + x + " = " + *temp + "\n"
		locals = locals + local
	}

	return locals, result, nil
}

// ParseResources handles resources in ARM conversion
func ParseResources(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) (map[string]interface{}, error) {
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
		first, err := see.Lookup(myType["type"].(string))
		if err != nil {
			log.Warn().Err(err)
			continue
		}

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
			var temp *string
			temp, result = parseString(attribute.(string), result)
			myResource[name] = *temp
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
					var temp *string
					temp, result = parseString(resource.(string), result)
					myArray[index] = *temp
				case map[string]interface{}:
					myArray[index], _ = parseMap(resource.(map[string]interface{}), result)
				case []interface{}:
					log.Print(resource)
				}

			}
			myResource[name] = myArray
		case float64, bool:
			//it's fine
			myResource[name] = attribute
		default:
			log.Print(attribute)
		}
	}
	return myResource, nil
}

func parseString(newAttribute string, result map[string]interface{}) (*string, map[string]interface{}) {
	var matches = []string{"parameters", "variables", "toLower", "resourceGroup().location", "resourceGroup().id",
		"substring", "uniqueString", "reference", "resourceId", "listKeys", "format('"}

	if what, found := contains(matches, newAttribute); found {
		newAttribute, result = replace(matches, newAttribute, what, result)
		newAttribute = loseSQBrackets(newAttribute)
	}
	return &newAttribute, result
}

func loseSQBrackets(newAttribute string) string {
	var re = regexp.MustCompile(`^\[(.*)\]`) //format('{0}/{1}',
	Matched := re.FindStringSubmatch(newAttribute)
	if len(Matched) > 1 {
		return Matched[1]
	}
	return newAttribute
}

func replace(matches []string, newAttribute string, what *string, result map[string]interface{}) (string, map[string]interface{}) {
	var Attribute string
	switch *what {
	case "reference", "resourceId":
		{
			if *what == "reference" {
				Attribute = loseSQBrackets(newAttribute)
				var re = regexp.MustCompile(`reference\((.*?)\)\.`) //format('{0}/{1}',
				Match := re.FindStringSubmatch(Attribute)
				if (Match) != nil {
					var rem = regexp.MustCompile(`\)\.(.*)`) //format('{0}/{1}',
					remains := rem.FindStringSubmatch(Attribute)
					if len(remains) <= 1 {
						log.Print("no attribute")
					}
					var re2 = regexp.MustCompile(`resourceId\((.*?)\)`) //format('{0}/{1}',
					Match2 := re2.FindStringSubmatch(Attribute)
					Attribute = Match2[0] + "." + remains[1]
				} else {
					log.Info().Msgf("ignoring %s", newAttribute)
				}
			} else {
				var err error
				Attribute, err = replaceResourceID(loseSQBrackets(newAttribute), result)
				if err != nil {
					log.Warn().Msgf("failed to parse %s", newAttribute)
				}
			}
		}
	case "uniqueString", "uniquestring":
		{
			target := *what + "\\((.*?)\\)"
			var re = regexp.MustCompile(target) //format('{0}/{1}',
			Match := re.ReplaceAllString(newAttribute, "substr(uuid(), 0, 8)")
			Attribute = Match
		}
	case "format('":
		{
			var re = regexp.MustCompile(`{.}`) //format('{0}/{1}',
			Match := re.ReplaceAllString(newAttribute, "%%s")
			Match = strings.Replace(Match, "'", "\"", -1)
			Attribute = loseSQBrackets(Match)
		}
	case "listKeys":
		{
			Attribute = loseSQBrackets(newAttribute)
			var re = regexp.MustCompile(`listKeys\((.*)\)`) //format('{0}/{1}',
			Match := re.FindStringSubmatch(Attribute)
			if len(Match) > 1 {
				resource := strings.Split(Match[1], ",")[0]
				Attribute = strings.Replace(Attribute, Match[0], resource, -1)
			} else {
				log.Warn().Msgf("failed to parse list keys")
			}
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
				Attribute = loseSQBrackets(strings.Replace(newAttribute, Match[0], temp, -1))
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
			data := make(map[string]bool)
			if result["data"] == nil {
				data["resource_group"] = true
			} else {
				data = result["data"].(map[string]bool)
				data["resource_group"] = true
			}
			result["data"] = data
		}
	case "resourceGroup().id":
		{
			Attribute = loseSQBrackets(strings.Replace(newAttribute, "resourceGroup().id",
				"data.azurerm_resource_group.sato.id", -1))
		}
	case "substring":
		{
			Attribute = strings.Replace(newAttribute, "substring",
				"substr", -1)
			Attribute = strings.Replace(Attribute, "'", "\"", -1)
		}
	}

	if again, still := contains(matches, Attribute); still {
		//allow failure
		if Attribute != newAttribute {
			Attribute, result = replace(matches, Attribute, again, result)
		} else {
			log.Printf("having a problem parsing %s", newAttribute)
		}
	}
	return Attribute, result
}

func isCompound(newAttribute string) bool {
	var compound bool
	if strings.Contains(newAttribute, "\"") {
		compound = true
	}
	return compound
}

func replaceResourceID(Match string, result map[string]interface{}) (string, error) {

	if strings.Contains(Match, "extensionResourceId") {
		var re = regexp.MustCompile(`extensionResourceId\((.*?)\)`)
		Match = re.ReplaceAllString(Match, "NIL")
	}

	var re = regexp.MustCompile(`resourceId\((.*?)\)`)
	Attribute := re.FindStringSubmatch(Match)

	var re2 = regexp.MustCompile(`resourceId\((.*?\))\),`)
	Attribute2 := re2.FindStringSubmatch(Match)

	if len(Attribute2) > 1 {
		Attribute = Attribute2
	}

	if len(Attribute) <= 1 {
		var re3 = regexp.MustCompile(`,(![^[]*\])`)
		Attribute = re3.FindStringSubmatch(Match)
	}
	arm, name, err := splitResourceName(Attribute[1])

	if err != nil {
		log.Warn().Msgf("failed to parse %s", Attribute[1])
	}

	name, err = findResourceName(result, name)
	if err != nil {
		log.Print(err)
	}
	var resourceName *string
	if findResourceType(result, arm) {
		resourceName, err = see.Lookup(arm)
		if err != nil {
			log.Printf("no match found %s", arm)
		}
	} else {
		if strings.Contains(arm, " ") {
			deSplitter := strings.Split(arm, " ")
			for _, name := range deSplitter {
				if strings.Contains(name, "Microsoft") {
					arm = name
				}
			}
		}

		switch arm {
		case "Microsoft.ContainerRegistry/registries":
			{
				temp := "azurerm_container_registry"
				resourceName = &temp
				name, err = findResourceName(result, name)
				if err != nil {
					log.Warn().Msgf("no match found %s", arm)
				}
			}
		case "Microsoft.Network/virtualNetworks/subnets":
			{
				temp := "azurerm_virtual_network"
				resourceName = &temp
				splutters := strings.Split(name, ", ")
				for item, splutter := range splutters {
					splutters[item], _ = findResourceName(result, splutter)
				}
				name = strings.Join(splutters, ".") + ".subnet"
			}
		case "Microsoft.Authorization/roleDefinitions":
			{
				temp := "azurerm_role_definition"
				resourceName = &temp
				if unicode.IsDigit(rune(name[0])) {
					name = "_" + name
				}
			}
		default:
			{
				if strings.Contains(name, "local") {
					name, err = findResourceName(result, name)
					if err != nil {
						log.Warn().Msgf("no match found %s", arm)
					}
				}
				resourceName, err = see.Lookup(sato.Dequote(arm))
				if err != nil {
					log.Warn().Msgf("no match found %s", arm)
				}
			}
		}
	}
	temp := *resourceName + "." + name
	return strings.Replace(Match, Attribute[0], temp, -1), nil
}

func splitResourceName(Attribute string) (string, string, error) {
	splitsy := strings.SplitN(Attribute, ",", 2)
	var name string
	var arm string

	switch len(splitsy) {
	case 1:
		{
			return "", "", fmt.Errorf("failed to parse resource name %s", Attribute)
		}
	case 2:
		{
			var re = regexp.MustCompile(`'(.*?)'`)
			newAttribute := re.FindStringSubmatch(splitsy[1])
			if len(newAttribute) <= 1 {
				arm = sato.Dequote(splitsy[0])
				name = strings.TrimSpace(splitsy[1])
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
	default:
		{
			//more than 2
		}
	}
	return arm, name, nil
}

func findResourceName(result map[string]interface{}, name string) (string, error) {
	if strings.HasPrefix(name, "format") {
		return name, fmt.Errorf("uses inline format function %s", name)
	}

	name = sato.Dequote(name)
	var err error
	resources := result["resources"].([]interface{})
	for _, myResource := range resources {
		test := myResource.(map[string]interface{})
		if name == test["name"].(string) {
			return test["resource"].(string), nil
		}
	}

	//not simple name lookup
	if strings.Contains(name, ",") {
		Lots := strings.Split(name, ",")
		var newName []string
		for _, lot := range Lots {
			var part string
			if part, err = findResourceName(result, strings.TrimSpace(lot)); err != nil {
				part, err = getNameValue(result, strings.TrimSpace(lot))

				if err != nil {
					return "", err
				}
			}

			newName = append(newName, part)
		}
		return strings.Join(newName, "."), nil
	}

	name, err = getNameValue(result, name)
	if err != nil {
		return "", err
	}

	return name, nil
}

func getNameValue(result map[string]interface{}, name string) (string, error) {
	if strings.Contains(name, ".") {
		rawNames := strings.Split(name, ".")
		if len(rawNames) <= 1 {
			return name, fmt.Errorf("failed to match value %s", name)
		}
		rawName := rawNames[1]
		variables := result["variables"].(map[string]interface{})
		for myVariable, value := range variables {
			if rawName == myVariable {
				return value.(string), nil
			}
		}
	}
	return name, nil
}

func findResourceType(result map[string]interface{}, name string) bool {
	resources := result["resources"].([]interface{})
	for _, myResource := range resources {
		test := myResource.(map[string]interface{})
		if name == test["type"].(string) {
			return true
		}
	}
	return false
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
	if result["outputs"] == nil {
		return nil
	}

	outputs := result["outputs"].(map[string]interface{})

	var All string
	for name, value := range outputs {
		var myVar sato.Output
		var someString *string
		myVar.Type = "string"
		myVar.Name = name
		temp := value.(map[string]interface{})
		someString, result = parseString(temp["value"].(string), result)
		myVar.Value = *someString
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

// ParseData writes out to data.tf
func ParseData(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) error {
	if result["data"] == nil {
		return nil
	}

	data := result["data"]

	var output bytes.Buffer
	tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(DataFile))
	if err != nil {
		return err
	}
	err = tmpl.Execute(&output, m{
		"data": data,
	})

	if err != nil {
		log.Print(err)
		return err
	}

	err = sato.Write(output.String(), destination, "data")
	if err != nil {
		log.Print(err)
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
	results["resources"] = setResourceNames(results)
	locals := make(map[string]interface{})

	if results["variables"] != nil {
		paraVariables := results["variables"].(map[string]interface{})

		newVariables := make(map[string]interface{})
		for item, result := range paraVariables {
			switch result.(type) {
			case string:
				{
					if strings.Contains(result.(string), "[") {
						locals[item] = result
					} else {
						newVariables[item] = result
					}
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
		_, ok := myResult["defaultValue"]
		if ok {
			myType := myResult["type"].(string)
			switch myType {
			case "string":
				{
					defaultValue := myResult["defaultValue"].(string)
					if strings.Contains(defaultValue, "[") {
						newLocals[item] = defaultValue
						break
					}
					myResult["default"] = myResult["defaultValue"].(string)
					newParams[item] = myResult
				}
			case "int":
				{
					myResult["type"] = "number"
					myResult["default"] = fmt.Sprintf("%v", myResult["defaultValue"].(float64))
					newParams[item] = myResult
				}
			case "object", "list(string)":
				{
					//todo
					//myResult["default"] = myResult["defaultValue"]
					myResult["default"] = ""
					newParams[item] = myResult
				}
			case "array":
				{
					myResult["type"] = "list(string)"
					myResult["default"] = arrayToString(myResult["defaultValue"].([]interface{}))
					newParams[item] = myResult
				}
			case "map[string]interface{}":
				{
					log.Printf("handled %s", myType)
				}
			default:
				{
					log.Printf("unhandled %s", myType)
				}
			}
		} else {
			myResult["default"] = ""
			newParams[item] = myResult
		}
	}
	results["parameters"] = newParams
	maps.Copy(locals, newLocals)
	results["locals"] = locals
	return results
}

func setResourceNames(results map[string]interface{}) []interface{} {
	resources := results["resources"].([]interface{})
	var newResults []interface{}
	for item, result := range resources {
		inside := result.(map[string]interface{})
		counter := map[string]interface{}{"resource": fmt.Sprintf("sato" + strconv.Itoa(item))}
		maps.Copy(inside, counter)
		newResults = append(newResults, inside)
	}
	return newResults
}
