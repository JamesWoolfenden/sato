package arm

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sato/src/cf"
	"sato/src/see"
	"strconv"
	"strings"
	tftemplate "text/template"
	"unicode"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

type m map[string]interface{}

// Parse turn CFN into Terraform.
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
		"Array":        cf.Array,
		"ArrayReplace": cf.ArrayReplace,
		"Contains":     cf.Contains,
		"Enabled":      enabled,
		"Sprint":       cf.Sprint,
		"Decode64":     cf.Decode64,
		"Boolean":      cf.Boolean,
		"Dequote":      cf.Dequote,
		"Quote":        cf.Quote,
		"Demap":        cf.Demap,
		"Tags":         tags,
		"ToUpper":      strings.ToUpper,
		"ToLower":      cf.Lower,
		"Deref":        func(str *string) string { return *str },
		"Nil":          cf.Nill,
		"Nild":         cf.Nild,
		"Marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)

			return string(a)
		},
		"Set":          arrayToString,
		"Split":        cf.Split,
		"SplitOn":      cf.SplitOn,
		"Replace":      cf.Replace,
		"RandomString": cf.RandomString,
		"Map":          cf.Map,
		"NotNil":       notNil,
		"Snake":        cf.Snake,
		"Kebab":        cf.Kebab,
		"ZipFile":      cf.Zipfile,
		"Uuid":         uuid,
	}

	result = preprocess(result)
	result, err = parseVariables(result, funcMap, destination)

	if err != nil {
		return err
	}

	result, err = parseResources(result, funcMap, destination)
	if err != nil {
		return err
	}

	err = parseOutputs(result, funcMap, destination)
	if err != nil {
		return err
	}

	err = parseData(result, funcMap, destination)
	if err != nil {
		return err
	}
	return nil
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
			var temp string
			temp, result = parseString(attribute.(string), result)
			myResource[name] = temp
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
					var temp string
					temp, result = parseString(resource.(string), result)
					myArray[index] = temp
				case map[string]interface{}:
					myArray[index], _ = parseMap(resource.(map[string]interface{}), result)
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

func parseString(attribute string, result map[string]interface{}) (string, map[string]interface{}) {
	matches := []string{
		"parameters", "variables", "toLower", "resourceGroup().location", "resourceGroup().id",
		"substring", "uniqueString", "reference", "resourceId", "listKeys", "format('", "SubscriptionResourceId",
		"concat", "subscription().tenantId", "uuid",
	}

	if what, found := contains(matches, attribute); found {
		attribute, result = replace(matches, attribute, what, result)
		attribute = loseSQBrackets(attribute)
	}
	return attribute, result
}

func replace(matches []string, newAttribute string, what *string, result map[string]interface{}) (string, map[string]interface{}) {
	var Attribute string
	switch *what {
	case "concat":
		{
			Attribute = loseSQBrackets(newAttribute)
			ditched := ditch(Attribute, "concat")

			raw := strings.Split(ditched, ",")
			var after string

			for item, value := range raw {
				if value == "/" {
					value = "_"
				}

				if strings.Contains(value, "Microsoft") {
					var err error
					value, err = resourceToName(value, result)
					if err != nil {
						log.Debug().Msgf("Concat failed: %v", err)
					}
				}

				s := []string{"var.", "azurerm_", "local.", "substr"}
				for _, v := range s {
					if strings.Contains(strings.ToLower(value), strings.ToLower(v)) {
						if v == "substr" {
							after = "${" + strings.TrimSpace(strings.Join(raw[1:], ",")) + "}"
							continue
						} else {
							raw[item] = fmt.Sprintf("${%s}", strings.ReplaceAll(strings.TrimSpace(value), "'", ""))
						}
					}

				}
				raw[item] = strings.ReplaceAll(strings.TrimSpace(raw[item]), "'", "")
			}

			if after == "" {
				Attribute = strings.Join(raw, "")
			} else {
				Attribute = raw[0] + after
			}
		}
	case "reference":
		{
			Attribute = ditch(loseSQBrackets(newAttribute), "reference")
		}
	case "resourceId":
		{
			Attribute = loseSQBrackets(newAttribute)

			var err error
			Attribute, err = replaceResourceID(Attribute, result)
			if err != nil {
				log.Warn().Msgf("failed to parse %s", newAttribute)
			}
		}
	case "uniqueString", "uniquestring":
		{
			target := "uniquestring\\((.*?)\\)"
			re := regexp.MustCompile(target) // format('{0}/{1}',
			Match := re.ReplaceAllString(strings.ToLower(newAttribute), "substr(uuid(), 0, 8)")
			Attribute = Match
		}
	case "subscriptionResourceId":
		{
			Attribute = ditch(newAttribute, "subscriptionResourceId")
		}
	case "format('":
		{
			re := regexp.MustCompile(`{.}`) // format('{0}/{1}',
			Match := re.ReplaceAllString(newAttribute, "%s")
			Match = strings.Replace(Match, "'", "\"", -1)
			Match = strings.Replace(Match, "/", "-", -1)
			Attribute = loseSQBrackets(Match)
		}
	case "listKeys":
		{
			Attribute = loseSQBrackets(newAttribute)
			re := regexp.MustCompile(`listKeys\((.*)\)`) // format('{0}/{1}',
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
			re := regexp.MustCompile(`parameters\('(.*?)\'\)`)
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
			re := regexp.MustCompile(`variables\('(.*?)\'\)`)
			Match := re.FindStringSubmatch(newAttribute)
			if (Match) != nil {
				var myTemp string
				if IsLocal(Match[1], result) {
					myTemp = "local." + Match[1]
				} else {
					myTemp = "var." + Match[1]
				}

				Attribute = strings.Replace(newAttribute, Match[0], myTemp, -1)
			} else {
				log.Printf("not found %s", newAttribute)
			}
		}
	case "toLower", "tolower":
		{
			re := regexp.MustCompile(`(?i)tolower`)
			Attribute = re.ReplaceAllString(newAttribute, "lower")
		}
	case "resourceGroup().location":
		{
			Attribute = strings.Replace(newAttribute, "resourceGroup().location",
				"data.azurerm_resource_group.sato.location", -1)
			data := make(map[string]interface{})
			if result["data"] == nil {
				data["resource_group"] = true
			} else {
				data = result["data"].(map[string]interface{})
				data["resource_group"] = true
			}
			result["data"] = data
		}
	case "uuid":
		{
			data := make(map[string]interface{})
			if result["data"] == nil {
				data["uuid"] = 0
			} else {
				data = result["data"].(map[string]interface{})
				//temp := data["uuid"].(int) + 1
				if data["uuid"] != nil {
					data["uuid"] = data["uuid"].(int) + 1
				} else {
					data["uuid"] = 0
				}
			}

			result["data"] = data
			replacement := "random_uuid.sato" + strconv.Itoa(data["uuid"].(int)) + ".result"
			Attribute = strings.Replace(newAttribute, "uuid()", replacement, 1)
		}
	case "subscription().tenantId":
		{
			Attribute = strings.Replace(newAttribute, "subscription().tenantId",
				"data.azurerm_client_config.sato.tenant_id", -1)
			data := make(map[string]interface{})
			if result["data"] == nil {
				data["client_config"] = true
			} else {
				data = result["data"].(map[string]interface{})
				data["client_config"] = true
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
		// allow failure
		if Attribute != newAttribute {
			Attribute, result = replace(matches, Attribute, again, result)
		} else {
			log.Printf("having a problem parsing %s", newAttribute)
		}
	}
	return Attribute, result
}

func replaceResourceID(Match string, result map[string]interface{}) (string, error) {
	if strings.Contains(Match, "extensionResourceId") {
		Match = strings.Replace(Match, "extensionResourceId", "resourceId", 1)
	}

	if strings.Contains(Match, "subscriptionResourceId") {
		Match = strings.Replace(Match, "subscriptionResourceId", "resourceId", 1)
	}

	re := regexp.MustCompile(`resourceId\((.*?)\)`)
	Attribute := re.FindStringSubmatch(Match)

	re2 := regexp.MustCompile(`resourceId\((.*?\))\),`)
	Attribute2 := re2.FindStringSubmatch(Match)

	if len(Attribute2) > 1 {
		Attribute = Attribute2
	}

	if len(Attribute) <= 1 {
		re3 := regexp.MustCompile(`,(![^[]*\])`)
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

		resourceName, err = see.Lookup(arm)

		if err != nil {
			return "", err
		}

		switch strings.ToLower(arm) {
		case "microsoft.containerregistry/registries":
			{
				temp := "azurerm_container_registry"
				resourceName = &temp
				name, err = findResourceName(result, name)
				if err != nil {
					log.Warn().Msgf("no match found %s", arm)
				}
			}
		case "microsoft.network/virtualnetworks/subnets":
			{
				temp := "tolist(azurerm_virtual_network"
				resourceName = &temp
				splutters := strings.Split(name, ", ")
				for item, splutter := range splutters {
					splutters[item], _ = findResourceName(result, splutter)
				}
				name = strings.Split(splutters[0], ".")[0] + ".subnet)[0].id"
			}
		case "microsoft.authorization/roledefinitions":
			{
				if unicode.IsDigit(rune(name[0])) {
					name = "_" + name
				}
			}
		case "microsoft.network/privateendpoints/privatednszonegroups":
			{
				// this isn't a separate terraform resource just part of
				// private_dns_zone_group in a azurerm_private_endpoint
				name, err = resourceToName(Match, result)
				if err != nil {
					return "", err
				}

				return *resourceName + "." + name + ".private_dns_zone_group.id", nil
			}
		case "microsoft.network/applicationgateways/httplisteners":
			{
				name, err = resourceToName(Match, result)

				if err != nil {
					return "", err
				}
				return *resourceName + "." + name + ".http_listener.id", nil
			}
		case "microsoft.network/applicationgateways/frontendipconfigurations":
			{
				name, err = resourceToName(Match, result)

				if err != nil {
					return "", err
				}
				return *resourceName + "." + name + ".frontend_ip_configuration.id", nil
			}
		case "microsoft.network/applicationgateways/frontendports":
			{
				name, err = resourceToName(Match, result)

				if err != nil {
					return "", err
				}
				return *resourceName + "." + name + ".frontend_port", nil
			}
		case "microsoft.network/applicationgateways/backendaddresspools":
			{
				name, err = resourceToName(Match, result)

				if err != nil {
					return "", err
				}
				return *resourceName + "." + name + ".frontend_ip_configuration.id", nil
			}
		case "microsoft.network/applicationgateways/backendhttpsettingscollection":
			{
				name, err = resourceToName(Match, result)

				if err != nil {
					return "", err
				}
				return *resourceName + "." + name + ".backend_http_settings", nil
			}
		default:
			{
				if strings.Contains(name, "local") {
					name, err = findResourceName(result, name)
					if err != nil {
						log.Warn().Msgf("no match found %s", arm)
					}
				}
				resourceName, err = see.Lookup(cf.Dequote(arm))
				if err != nil {
					resourceName = toUnknownPointer()
					log.Warn().Msgf("no match found %s", arm)
				}
			}
		}
	}

	if resourceName != nil {
		temp := *resourceName + "." + name
		return strings.Replace(Match, Attribute[0], temp, -1), nil
	}

	return "", err
}

func resourceToName(Match string, result map[string]interface{}) (string, error) {
	re := regexp.MustCompile(`^resourceId\((.*)\)`)
	splitter := re.FindStringSubmatch(Match)
	if len(splitter) > 1 {
		//ditch type
		_, found, _ := strings.Cut(splitter[1], ",")
		myResourceName, _, _ := strings.Cut(found, ",")
		name, err := findResourceName(result, strings.TrimSpace(myResourceName))

		if err != nil {
			log.Warn().Msgf("no match found %s", Match)
		} else {
			return name, err
		}
	}

	if len(splitter) == 0 {
		name, err := findResourceName(result, strings.TrimSpace(Match))

		if err != nil {
			log.Warn().Msgf("no match found %s", Match)
		} else {
			return name, err
		}
	}

	return "", fmt.Errorf("failed to split resource %s", Match)
}

func toUnknownPointer() *string {
	temp := "UNKNOWN"
	return &temp
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
			re := regexp.MustCompile(`'(.*?)'`)
			newAttribute := re.FindStringSubmatch(splitsy[1])
			if len(newAttribute) <= 1 {
				arm = cf.Dequote(splitsy[0])
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
	default:
		{
			// more than 2
		}
	}
	return arm, name, nil
}

func findResourceName(result map[string]interface{}, name string) (string, error) {
	if strings.HasPrefix(name, "format") {
		return name, fmt.Errorf("uses inline format function %s", name)
	}

	name = cf.Dequote(name)
	var err error

	if result["resources"] == nil {
		return "", fmt.Errorf("resources are empty")
	}

	resources, ok := result["resources"].([]interface{})

	if !ok {
		return name, fmt.Errorf("no resources found")
	}

	for _, myResource := range resources {
		test, ok := myResource.(map[string]interface{})
		if !ok {
			log.Print("resource is not a map")
			continue
		}
		temp := loseSQBrackets(test["name"].(string))

		if name == temp {
			return test["resource"].(string), nil
		}

		trimName := strings.Replace(name, "var.", "", 1)
		trimTemp := strings.Replace(ditch(temp, "variables"), "'", "", 2)

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
			name = strings.Trim(name, "var.")

			if name == strings.ReplaceAll(splits[1], "'", "") {
				return test["resource"].(string), nil
			}
		}
	}

	// not simple name lookup
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
		if len(rawNames) != 2 {
			return name, fmt.Errorf("failed to match value %s", name)
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

func findResourceType(result map[string]interface{}, name string) bool {
	if result["resources"] == nil {
		return false
	}

	resources := result["resources"].([]interface{})
	for _, myResource := range resources {
		test := myResource.(map[string]interface{})
		if name == test["type"].(string) {
			return true
		}
	}
	return false
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
			switch strings.ToLower(myType) {
			case "string", "securestring":
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
					// todo
					// myResult["default"] = myResult["defaultValue"]
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
			case "bool":
				{
					myResult["default"] = fmt.Sprintf("%v", myResult["defaultValue"])
					newParams[item] = myResult
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
