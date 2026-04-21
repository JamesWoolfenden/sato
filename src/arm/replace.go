package arm

import (
	"fmt"
	"regexp"
	"sato/src/see"
	"sato/src/tfgen"
	"strconv"
	"strings"
	"unicode"

	"github.com/rs/zerolog/log"
)

// subResourceSuffix maps ARM sub-resource types that translate to a block on
// a parent terraform resource rather than a resource of their own. The value
// is the attribute path appended after "<tf_type>.<name>.".
var subResourceSuffix = map[string]string{
	"microsoft.network/privateendpoints/privatednszonegroups":             "private_dns_zone_group.id",
	"microsoft.network/applicationgateways/httplisteners":                 "http_listener.id",
	"microsoft.network/applicationgateways/frontendipconfigurations":      "frontend_ip_configuration.id",
	"microsoft.network/applicationgateways/frontendports":                 "frontend_port",
	"microsoft.network/applicationgateways/backendaddresspools":           "frontend_ip_configuration.id",
	"microsoft.network/applicationgateways/backendhttpsettingscollection": "backend_http_settings",
	"microsoft.network/applicationgateways/authenticationcertificates":    "authentication_certificates",
	"microsoft.network/applicationgateways/sslcertificates":               "ssl_certificates",
}

// replaceReference rewrites reference(expr [, apiVersion [, 'Full']]) to just
// expr, preserving any prefix/suffix around the call. ARM's reference() returns
// a runtime object whose attributes map onto the terraform resource directly,
// so the wrapper and its extra args are dropped.
func replaceReference(attr string) string {
	idx := strings.Index(attr, "reference(")
	if idx < 0 {
		return attr
	}

	prefix := attr[:idx]
	rest := attr[idx+len("reference("):]

	depth := 1
	argEnd := -1

	var i int
	for i = 0; i < len(rest); i++ {
		switch rest[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				if argEnd < 0 {
					argEnd = i
				}
				return prefix + rest[:argEnd] + rest[i+1:]
			}
		case ',':
			if depth == 1 && argEnd < 0 {
				argEnd = i
			}
		}
	}

	return attr
}

// replaceUUID handles the uuid ARM function conversion.
func replaceUUID(newAttribute string, result map[string]interface{}) (string, map[string]interface{}) {
	if !strings.Contains(newAttribute, "uuid()") {
		return newAttribute, result
	}

	data := ensureDataMap(result)

	if data["uuid"] != nil {
		data["uuid"] = data["uuid"].(int) + 1
	} else {
		data["uuid"] = 0
	}

	result["data"] = data
	replacement := "random_uuid.sato" + strconv.Itoa(data["uuid"].(int)) + ".result"
	return strings.Replace(newAttribute, "uuid()", replacement, 1), result
}

// replaceConcat handles the concat ARM function conversion.
func replaceConcat(newAttribute string, result map[string]interface{}) string {
	Attribute := LoseSQBrackets(newAttribute)
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
		return strings.Join(raw, "")
	}

	return raw[0] + after
}

// Replace convert ARM functions to tf.
func Replace(
	matches []string,
	newAttribute string,
	what *string,
	result map[string]interface{}) (string, map[string]interface{}) {

	var Attribute string

	switch *what {
	case "uri(":
		{
			Attribute = Ditch(LoseSQBrackets(newAttribute), "uri")
		}
	case "concat":
		{
			Attribute = replaceConcat(newAttribute, result)
		}
	case "reference":
		{
			Attribute = replaceReference(LoseSQBrackets(newAttribute))
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
			re := regexp.MustCompile(`(?i)uniquestring\((.*?)\)`)
			Attribute = re.ReplaceAllString(newAttribute, "substr(uuid(), 0, 8)")
		}
	case "SubscriptionResourceId", "subscriptionResourceId":
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
				log.Warn().Msgf("no match found %s", newAttribute)
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
				log.Warn().Msgf("not found %s", newAttribute)
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
			data := ensureDataMap(result)
			data["resource_group"] = true
			result["data"] = data
		}
	case "uuid(":
		{
			Attribute, result = replaceUUID(newAttribute, result)
		}
	case "subscription().tenantId":
		{
			Attribute = strings.ReplaceAll(newAttribute, "subscription().tenantId",
				"data.azurerm_client_config.sato.tenant_id")
			data := ensureDataMap(result)
			data["client_config"] = true
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
			log.Warn().Msgf("having a problem parsing %s", newAttribute)
		}
	}

	return Attribute, result
}

// ReplaceResourceID ditches rssourceID.
func ReplaceResourceID(match string, result map[string]interface{}) (string, error) {
	match = strings.Replace(match, "extensionResourceId", "resourceId", 1)
	match = strings.Replace(match, "subscriptionResourceId", "resourceId", 1)

	match = strings.ReplaceAll(match, "resourceID", "resourceId")

	re := regexp.MustCompile(`resourceId\((.*?)\)`)
	Attribute := re.FindStringSubmatch(match)

	re2 := regexp.MustCompile(`resourceId\((.*?\))\),`)
	Attribute2 := re2.FindStringSubmatch(match)

	if len(Attribute2) > 1 {
		Attribute = Attribute2
	}

	if len(Attribute) <= 1 {
		re3 := regexp.MustCompile(`,(![^[]*\])`)
		Attribute = re3.FindStringSubmatch(match)
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
			log.Warn().Msgf("no match found %s", arm)
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

		if suffix, ok := subResourceSuffix[strings.ToLower(arm)]; ok {
			name, err = resourceToName(match, result)
			if err != nil {
				return "", err
			}

			return *resourceName + "." + name + "." + suffix, nil
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
		default:
			{
				if strings.Contains(name, "local") {
					name, err = FindResourceName(result, name)
					if err != nil {
						log.Warn().Msgf("no match found %s", arm)
					}
				}

				resourceName, err = see.Lookup(tfgen.Dequote(arm), false)

				if err != nil {
					resourceName = toUnknownPointer()

					log.Warn().Msgf("no match found %s", arm)
				}
			}
		}
	}

	if resourceName != nil {
		temp := *resourceName + "." + name
		return strings.ReplaceAll(match, Attribute[0], temp), nil
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
