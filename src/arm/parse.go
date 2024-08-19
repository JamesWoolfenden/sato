package arm

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	tftemplate "text/template"
	"unicode"

	"sato/src/cf"
	"sato/src/see"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

type splitResourceError struct {
	match string
}

func (e splitResourceError) Error() string {
	return fmt.Sprintf("failed to split resource %s", e.match)
}

type filepathError struct {
	Path string
}

func (m *filepathError) Error() string {
	return fmt.Sprintf("not implemented %s", m.Path)
}

type parseListError struct{}

func (m *parseListError) Error() string {
	return "parseListError"
}

type parseMapError struct {
	Err error
}

func (m *parseMapError) Error() string {
	return fmt.Sprintf("parseMapError %s", m.Err)
}

type m map[string]interface{}

// Parse turn CFN into Terraform.
func Parse(file string, destination string) error {
	fileAbs, err := filepath.Abs(file)
	if err != nil {
		return &filepathError{Path: file}
	}

	jsonFile, err := os.Open(fileAbs)
	if err != nil {
		return fmt.Errorf("file open failure %w", err)
	}

	// defer the closing of our jsonFile so that we can parse it later on
	//goland:noinspection GoUnhandledErrorResult
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)

	if err != nil {
		return fmt.Errorf("unmarshall failure %w", err)
	}

	funcMap := tftemplate.FuncMap{
		"Array":        cf.Array,
		"ArrayReplace": cf.ArrayReplace,
		"Contains":     cf.Contains,
		"Enabled":      Enabled,
		"Sprint":       cf.Sprint,
		"Decode64":     cf.Decode64,
		"Boolean":      cf.Boolean,
		"Dequote":      cf.Dequote,
		"Quote":        cf.Quote,
		"Demap":        cf.Demap,
		"Tags":         Tags,
		"ToUpper":      strings.ToUpper,
		"ToLower":      cf.Lower,
		"Deref":        func(str *string) string { return *str },
		"Nil":          cf.Nill,
		"Nild":         cf.Nild,
		"Marshal": func(v interface{}) string {
			a, err := json.Marshal(v)
			if err != nil {
				log.Printf("marshal failure")
			}

			return string(a)
		},
		"Set":          ArrayToString,
		"Split":        cf.Split,
		"SplitOn":      cf.SplitOn,
		"Replace":      cf.Replace,
		"RandomString": cf.RandomString,
		"Map":          cf.Map,
		"NotNil":       NotNil,
		"Snake":        cf.Snake,
		"Kebab":        cf.Kebab,
		"ZipFile":      cf.Zipfile,
		"Uuid":         UUID,
	}

	result = Preprocess(result)
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
		return fmt.Errorf("parse data %w", err)
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
		"concat", "subscription().tenantId", "UUID", "uri(",
	}

	if what, found := Contains(matches, attribute); found {
		attribute, result = Replace(matches, attribute, what, result)
		attribute = LoseSQBrackets(attribute)
	}

	return attribute, result
}

// Replace convert ARM functions to tf
func Replace(matches []string, newAttribute string, what *string, result map[string]interface{}) (string, map[string]interface{}) {
	var Attribute string

	switch *what {
	case "uri(":
		{
			Attribute = Ditch(LoseSQBrackets(newAttribute), "uri")
		}
	case "concat":
		{
			Attribute = LoseSQBrackets(newAttribute)
			ditched := Ditch(Attribute, "concat")

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
						}

						raw[item] = fmt.Sprintf("${%s}", strings.ReplaceAll(strings.TrimSpace(value), "'", ""))
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
			Attribute = Ditch(LoseSQBrackets(newAttribute), "reference")
		}
	case "resourceId":
		{
			Attribute = LoseSQBrackets(newAttribute)

			var err error
			Attribute, err = ReplaceResourceID(Attribute, result)

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
			Attribute = Ditch(newAttribute, "subscriptionResourceId")
		}
	case "format('":
		{
			re := regexp.MustCompile(`{.}`) // format('{0}/{1}',
			Match := re.ReplaceAllString(newAttribute, "%s")
			Match = strings.ReplaceAll(Match, "'", "\"")
			Match = strings.ReplaceAll(Match, "/", "-")
			Attribute = LoseSQBrackets(Match)
		}
	case "listKeys":
		{
			Attribute = LoseSQBrackets(newAttribute)
			re := regexp.MustCompile(`listKeys\((.*)\)`) // format('{0}/{1}',
			Match := re.FindStringSubmatch(Attribute)

			if len(Match) > 1 {
				resource := strings.Split(Match[1], ",")[0]
				Attribute = strings.ReplaceAll(Attribute, Match[0], resource)
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

				Attribute = LoseSQBrackets(strings.ReplaceAll(newAttribute, Match[0], temp))
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

				Attribute = strings.ReplaceAll(newAttribute, Match[0], myTemp)
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
			Attribute = strings.ReplaceAll(newAttribute, "resourceGroup().location",
				"data.azurerm_resource_group.sato.location")
			data := make(map[string]interface{})

			if result["data"] == nil {
				data["resource_group"] = true
			} else {
				var ok bool

				data, ok = result["data"].(map[string]interface{})
				if !ok {
					log.Printf("assertion failure %s", result["data"])

					break
				}

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
				data, ok := result["data"].(map[string]interface{})

				if !ok {
					log.Printf("assertion failure %s", result["data"])
				}

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
				data, ok := result["data"].(map[string]interface{})
				if !ok {
					log.Printf("assertion failure %s", result["data"])
				}

				data["client_config"] = true
			}

			result["data"] = data
		}
	case "resourceGroup().id":
		{
			Attribute = LoseSQBrackets(strings.ReplaceAll(
				newAttribute,
				"resourceGroup().id",
				"data.azurerm_resource_group.sato.id",
			))
		}
	case "substring":
		{
			Attribute = strings.ReplaceAll(newAttribute, "substring", "substr")
			Attribute = strings.ReplaceAll(Attribute, "'", "\"")
		}
	}

	if again, still := Contains(matches, Attribute); still {
		// allow failure
		if Attribute != newAttribute {
			Attribute, result = Replace(matches, Attribute, again, result)
		} else {
			log.Printf("having a problem parsing %s", newAttribute)
		}
	}

	return Attribute, result
}

// ReplaceResourceID ditches rssourceID
func ReplaceResourceID(Match string, result map[string]interface{}) (string, error) {
	Match = strings.Replace(Match, "extensionResourceId", "resourceId", 1)
	Match = strings.Replace(Match, "subscriptionResourceId", "resourceId", 1)

	Match = strings.ReplaceAll(Match, "resourceID", "resourceId")

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

	arm, name, err := SplitResourceName(Attribute[1])
	if err != nil {
		log.Warn().Msgf("failed to parse %s", Attribute[1])
	}

	name, err = FindResourceName(result, name)
	if err != nil {
		log.Print(err)
	}

	var resourceName *string
	if FindResourceType(result, arm) {
		resourceName, err = see.Lookup(arm, false)
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

		resourceName, err = see.Lookup(arm, false)

		if err != nil {
			return "", err
		}

		switch strings.ToLower(arm) {
		case "microsoft.containerregistry/registries":
			{
				temp := "azurerm_container_registry"
				resourceName = &temp
				name, err = FindResourceName(result, name)

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
					splutters[item], _ = FindResourceName(result, splutter)
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
		case "microsoft.network/applicationgateways/authenticationcertificates":
			{
				name, err = resourceToName(Match, result)

				if err != nil {
					return "", err
				}

				return *resourceName + "." + name + ".authentication_certificates", nil
			}
		case "microsoft.network/applicationgateways/sslcertificates":
			{
				name, err = resourceToName(Match, result)

				if err != nil {
					return "", err
				}

				return *resourceName + "." + name + ".ssl_certificates", nil
			}
		default:
			{
				if strings.Contains(name, "local") {
					name, err = FindResourceName(result, name)
					if err != nil {
						log.Warn().Msgf("no match found %s", arm)
					}
				}

				resourceName, err = see.Lookup(cf.Dequote(arm), false)

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

func resourceToName(match string, result map[string]interface{}) (string, error) {
	re := regexp.MustCompile(`^resourceId\((.*)\)`)
	splitter := re.FindStringSubmatch(match)

	if len(splitter) > 1 {
		// Ditch type
		_, found, _ := strings.Cut(splitter[1], ",")
		myResourceName, _, _ := strings.Cut(found, ",")
		name, err := FindResourceName(result, strings.TrimSpace(myResourceName))

		if err != nil {
			log.Warn().Msgf("no match found %s", match)
		} else {
			return name, err
		}
	}

	if len(splitter) == 0 {
		name, err := FindResourceName(result, strings.TrimSpace(match))

		if err != nil {
			log.Warn().Msgf("no match found %s", match)
		} else {
			return name, err
		}
	}

	return "", &splitResourceError{match}
}

func toUnknownPointer() *string {
	temp := "UNKNOWN"

	return &temp
}

// SplitResourceName splits and converts arm names into terraform
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

// FindResourceName looks for resource names
func FindResourceName(result map[string]interface{}, name string) (string, error) {
	if strings.HasPrefix(name, "format") {
		return name, &inlineFormatError{Name: name}
	}

	name = cf.Dequote(name)

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

// FindResourceType get resource types
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

// Preprocess examines raw ARM loads.
func Preprocess(results map[string]interface{}) map[string]interface{} {
	results["resources"] = SetResourceNames(results)
	locals := make(map[string]interface{})

	if results["variables"] != nil {
		paraVariables := results["variables"].(map[string]interface{})

		newVariables := make(map[string]interface{})

		for item, result := range paraVariables {
			switch result := result.(type) {
			case string:
				{
					if strings.Contains(result, "[") {
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
					myResult["default"] = ArrayToString(myResult["defaultValue"].([]interface{}))
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

// SetResourceNames gets resource names for results
func SetResourceNames(results map[string]interface{}) []interface{} {
	resources := results["resources"].([]interface{})

	var newResults []interface{}

	for item, result := range resources {
		inside := result.(map[string]interface{})
		counter := map[string]interface{}{"resource": fmt.Sprintf("sato%d", item)}
		maps.Copy(inside, counter)
		newResults = append(newResults, inside)
	}

	return newResults
}
