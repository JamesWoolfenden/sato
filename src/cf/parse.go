// Package cf provides functionality for parsing and converting AWS CloudFormation
// templates to Terraform configuration files.
package cf

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/awslabs/goformation/v7"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/rs/zerolog/log"
)

const (
	typeNumber     = "number"
	typeListString = "list(string)"
	typeString     = "string"
)

// M is a wrapper object to help pass in multiple objects to template.
type M map[string]interface{}

// Variable describes a terraform variable.
type Variable struct {
	Description string
	Type        string
	Default     string
	Name        string
}

// Output describes Tf output type.
type Output struct {
	Description string
	Type        string
	Value       string
	Name        string
}

var funcMap = template.FuncMap{
	"Array":        Array,
	"ArrayReplace": ArrayReplace,
	"Contains":     Contains,
	"Sprint":       Sprint,
	"Decode64":     Decode64,
	"Boolean":      Boolean,
	"Dequote":      Dequote,
	"Quote":        Quote,
	"Demap":        Demap,
	"ToUpper":      strings.ToUpper,
	"ToLower":      Lower,
	"Deref":        func(str *string) string { return *str },
	"Nil":          Nill,
	"Nild":         Nild,
	"Marshal":      Marshal,
	"Split":        Split,
	"SplitOn":      SplitOn,
	"Replace":      Replace,
	"Tags":         Tags,
	"RandomString": RandomString,
	"Map":          Map,
	"Snake":        Snake,
	"Kebab":        Kebab,
	"ZipFile":      Zipfile,
	"TFJoin":       ParseJoin,  // CloudFormation !Join → Terraform join()
	"TFSplit":      ParseSplit, // CloudFormation !Split → Terraform split()
}

// Parse turn CFN into Terraform.
func Parse(file string, destination string) error {
	if file == "" || destination == "" {
		return &emptyPathsError{}
	}

	// Open a cloudFormation from file (can be JSON or YAML)
	fileAbs, err := filepath.Abs(file)
	if err != nil {
		return &filepathError{Path: file, Err: err}
	}

	cloudFormation, err := goformation.Open(fileAbs)
	if err != nil {
		return &goformationError{err: err}
	}

	_, err = ParseVariables(cloudFormation, funcMap, destination)

	if err != nil {
		return &parseVariablesError{err: err}
	}

	err = parseResources(cloudFormation.Resources, funcMap, destination)
	if err != nil {
		return &parseResourcesError{err: err}
	}

	return nil
}

// ParseVariables convert CFN Parameters into terraform variables.
func ParseVariables(
	cloudFormation *cloudformation.Template,
	funcMap template.FuncMap,
	destination string,
) ([]Variable, error) {
	var All string

	var (
		myMap         = make(map[string]bool)
		DataResources []string
		myVariables   []Variable
	)

	for Name, param := range cloudFormation.Parameters {
		var myVariable Variable

		DataResources, myVariable, myMap = GetVariableType(param, myVariable, DataResources, myMap)
		myVariable.Description = strings.ReplaceAll(*param.Description, "${", "$${")
		myVariable.Name = Name
		myVariable = GetVariableDefault(param, myVariable)

		var output bytes.Buffer

		tmpl, err := template.New("test").Funcs(funcMap).Parse(string(variableFile))
		if err != nil {
			return nil, &templateNewError{Err: err}
		}

		_ = tmpl.Execute(&output, M{
			"variable": myVariable,
			"item":     Name,
		})

		All += output.String()

		myVariables = append(myVariables, myVariable)
	}

	err := Write(All, destination, "variables")
	if err != nil {
		return nil, &writeError{destination: destination, err: err}
	}

	err = Write(strings.Join(DataResources, "\n"), destination, "data")
	if err != nil {
		return nil, &writeError{destination: destination, err: err}
	}

	return myVariables, nil
}

// GetVariableType determines variable types.
func GetVariableType(
	param cloudformation.Parameter,
	myVariable Variable,
	dataResources []string,
	myMap map[string]bool) (
	[]string,
	Variable,
	map[string]bool,
) {
	switch param.Type {
	case "String":
		if param.Default == "false" || param.Default == "true" || param.Default == true || param.Default == false {
			myVariable.Type = "bool"
		} else {
			myVariable.Type = strings.ToLower(param.Type)
		}
	case "CommaDelimitedList":
		myVariable.Type = typeListString
	case "List<AWS::EC2::AvailabilityZone::Name>":
		myVariable.Type = typeListString
		dataResources, myMap = Add(dataAvailabilityZone, dataResources, myMap)
	case "AWS::EC2::Subnet::Id":
		myVariable.Type = typeString
		dataResources, myMap = Add(dataSubnet, dataResources, myMap)
	case "AWS::EC2::KeyPair::KeyName":
		myVariable.Type = typeString
		dataResources, myMap = Add(dataKeyPair, dataResources, myMap)
	case "AWS::EC2::VPC::Id", "List<AWS::EC2::VPC::Id>":
		myVariable.Type = typeString
		dataResources, myMap = Add(dataVpc, dataResources, myMap)
	case "AWS::EC2::SecurityGroup::Id":
		myVariable.Type = typeString
		dataResources, myMap = Add(dataSecurityGroup, dataResources, myMap)
	case "AWS::EC2::Image::Id", "AWS::SSM::Parameter::Value<AWS::EC2::Image::Id>":
		myVariable.Type = typeString
	case "AWS::Region":
		myVariable.Type = typeString
		dataResources, myMap = Add(dataRegion, dataResources, myMap)
	case "List<AWS::EC2::Subnet::Id>":
		myVariable.Type = typeListString
	case "Number":
		myVariable.Type = typeNumber
	default:
		log.Info().Msgf("Variable %s", param.Type)
	}

	dataResources, myMap = Add(provider, dataResources, myMap)

	return dataResources, myVariable, myMap
}

// GetVariableDefault determines a variables default value.
func GetVariableDefault(param cloudformation.Parameter, myVariable Variable) Variable {
	//goland:noinspection GoLinter
	switch param.Default.(type) {
	case string:
		_, err := strconv.Atoi(param.Default.(string))

		if err == nil {
			myVariable.Type = typeNumber

			myVariable.Default = param.Default.(string)
		} else {
			if myVariable.Type == "bool" {
				myVariable.Default = param.Default.(string)
			} else {
				if strings.Contains(param.Default.(string), "=") {
					myVariable = StringToMap(param)
				} else {
					myVariable.Default = "\"" + param.Default.(string) + "\""
				}
			}
		}
	case float64:
		myVariable.Type = typeNumber
		myVariable.Default = fmt.Sprintf("%v", param.Default.(float64))
	case bool:
		myVariable.Default = strconv.FormatBool(param.Default.(bool))
	case interface{}:
		myVariable.Default = "[]"
	default:
		myVariable.Default = "null"
	}

	return myVariable
}

// StringToMap converts maps in strings(for tags).
func StringToMap(param cloudformation.Parameter) Variable {
	temp := strings.Split(param.Default.(string), "=")

	var myVariable Variable

	var myMap string

	for item := 0; item < len(temp); item++ {
		if item == 0 {
			myMap += "{ "
		}

		if item%2 == 0 {
			myMap += "\"" + temp[item] + "\" = "
		} else {
			myMap += "\"" + temp[item] + "\""
		}
	}

	myVariable.Default = myMap + "}"
	myVariable.Type = "map(string)"

	return myVariable
}

// Write out Terraform.
func Write(output string, location string, name string) error {
	if output != "" {
		newPath, err := filepath.Abs(location)
		if err != nil {
			return &filepathError{Path: location, Err: err}
		}

		err = os.MkdirAll(newPath, 0o750)
		if err != nil {
			return &makeDirError{Err: err}
		}

		d1 := []byte(output)

		destination, err := filepath.Abs(fmt.Sprint(location, "/", name, ".tf"))
		if err != nil {
			return &filepathError{Path: fmt.Sprint(location, "/", name, ".tf"), Err: err}
		}

		err = os.WriteFile(destination, d1, 0o600)
		log.Info().Msgf("Created %s", destination)

		if err != nil {
			return &writeFileError{Destination: destination, Err: err}
		}

		// Format the generated Terraform file with tofu fmt
		cmd := exec.Command("tofu", "fmt", destination) // #nosec G204 -- destination is validated filepath from Write function
		if err := cmd.Run(); err != nil {
			// If tofu is not available, log a warning but don't fail
			log.Warn().Msgf("Could not format %s with tofu fmt: %v", destination, err)
		} else {
			log.Info().Msgf("Formatted %s with tofu fmt", destination)
		}
	}

	return nil
}

// ToTFName creates a Terraform resource name from a CFN type (approximates).
func ToTFName(cloudformation string) string {
	return strings.ToLower(strings.ReplaceAll(cloudformation, "::", "_"))
}

// ReplaceVariables looks to see if u can translate CFN vars into terraform.
func ReplaceVariables(str1 string) string {
	re := regexp.MustCompile(`\${.*?}`)
	submatch := re.FindAllString(str1, -1)

	for _, target := range submatch {
		if !strings.Contains(target, "::") {
			brReplace := strings.Replace(target, "${", "${var.", 1)
			str1 = strings.Replace(str1, target, brReplace, 1)
		}
	}

	return str1
}

// ReplaceDependant is fancy!
func ReplaceDependant(str1 string) string {
	replacer := strings.NewReplacer(
		"AWS::Region", "data.aws_region.current.name")

	return replacer.Replace(str1)
}

// ParseRef converts CloudFormation Ref intrinsic function to Terraform references.
// Handles both ${Ref::ResourceName} and plain Ref patterns.
func ParseRef(input string, parameters map[string]cloudformation.Parameter, resources cloudformation.Resources) string {
	// Pattern: ${Ref::ResourceName} or ${ResourceName}
	refPattern := regexp.MustCompile(`\$\{(?:Ref::)?([^}]+)\}`)

	return refPattern.ReplaceAllStringFunc(input, func(match string) string {
		// Extract the resource/parameter name
		name := refPattern.FindStringSubmatch(match)[1]

		// Check if it's a parameter
		if _, isParam := parameters[name]; isParam {
			return fmt.Sprintf("var.%s", strings.ToLower(name))
		}

		// Check if it's a resource
		if resource, isResource := resources[name]; isResource {
			resourceType := ToTFName(resource.AWSCloudFormationType())
			return fmt.Sprintf("%s.%s.id", resourceType, strings.ToLower(name))
		}

		// Return as variable if unknown
		return fmt.Sprintf("var.%s", strings.ToLower(name))
	})
}

// ParseGetAtt converts CloudFormation GetAtt intrinsic function to Terraform attribute references.
// Handles patterns like ${ResourceName.AttributeName} or Fn::GetAtt syntax.
func ParseGetAtt(input string) string {
	// Pattern: ${ResourceName.AttributeName}
	getAttPattern := regexp.MustCompile(`\$\{([^.]+)\.([^}]+)\}`)

	return getAttPattern.ReplaceAllStringFunc(input, func(match string) string {
		parts := getAttPattern.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}

		resourceName := parts[1]
		attributeName := parts[2]

		// Convert attribute name to Terraform format (CamelCase to snake_case)
		// Handle transitions from lowercase/digit to uppercase
		tfAttribute := regexp.MustCompile(`([a-z0-9])([A-Z])`).ReplaceAllString(attributeName, "${1}_${2}")
		// Handle transitions from multiple uppercase to lowercase (e.g., "DNSName" -> "DNS_Name")
		tfAttribute = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`).ReplaceAllString(tfAttribute, "${1}_${2}")
		tfAttribute = strings.ToLower(tfAttribute)

		// Common CloudFormation to Terraform attribute mappings
		attributeMap := map[string]string{
			"arn":                  "arn",
			"id":                   "id",
			"name":                 "name",
			"dns_name":             "dns_name",
			"hosted_zone_id":       "hosted_zone_id",
			"regional_domain_name": "regional_domain_name",
			"queue_url":            "url",
			"topic_arn":            "arn",
			"function_arn":         "arn",
			"role_arn":             "arn",
		}

		if mappedAttr, ok := attributeMap[tfAttribute]; ok {
			tfAttribute = mappedAttr
		}

		// For unknown resources, we can't determine the type, so use a placeholder
		return fmt.Sprintf("${%s.%s}", strings.ToLower(resourceName), tfAttribute)
	})
}

// ParseSub converts CloudFormation !Sub intrinsic function to Terraform string interpolation.
// Handles variable substitution patterns like ${VariableName} and pseudo-parameters like ${AWS::Region}.
func ParseSub(input string, parameters map[string]cloudformation.Parameter) string {
	// Pattern for ${Variable} substitution
	subPattern := regexp.MustCompile(`\$\{([^}]+)\}`)

	return subPattern.ReplaceAllStringFunc(input, func(match string) string {
		varName := subPattern.FindStringSubmatch(match)[1]

		// Handle pseudo-parameters
		pseudoParams := map[string]string{
			"AWS::Region":           "data.aws_region.current.name",
			"AWS::AccountId":        "data.aws_caller_identity.current.account_id",
			"AWS::StackName":        "var.stack_name",
			"AWS::StackId":          "var.stack_id",
			"AWS::Partition":        "data.aws_partition.current.partition",
			"AWS::URLSuffix":        "data.aws_partition.current.dns_suffix",
			"AWS::NotificationARNs": "var.notification_arns",
		}

		if tfValue, isPseudo := pseudoParams[varName]; isPseudo {
			return fmt.Sprintf("${%s}", tfValue)
		}

		// Check if it's a parameter
		if _, isParam := parameters[varName]; isParam {
			return fmt.Sprintf("${var.%s}", strings.ToLower(varName))
		}

		// Assume it's a resource reference (GetAtt pattern: Resource.Attribute)
		if strings.Contains(varName, ".") {
			parts := strings.SplitN(varName, ".", 2)
			resourceName := strings.ToLower(parts[0])
			attribute := strings.ToLower(parts[1])
			return fmt.Sprintf("${%s.%s}", resourceName, attribute)
		}

		// Plain variable reference
		return fmt.Sprintf("${var.%s}", strings.ToLower(varName))
	})
}

// ParseJoin converts CloudFormation !Join intrinsic function to Terraform join().
// Example: !Join ["-", ["a", "b", "c"]] becomes join("-", ["a", "b", "c"]).
func ParseJoin(delimiter string, values []string) string {
	// Escape any quotes in values
	escapedValues := make([]string, len(values))
	for i, v := range values {
		escapedValues[i] = fmt.Sprintf("\"%s\"", strings.ReplaceAll(v, "\"", "\\\""))
	}

	return fmt.Sprintf("join(\"%s\", [%s])", delimiter, strings.Join(escapedValues, ", "))
}

// ParseSplit converts CloudFormation !Split intrinsic function to Terraform split().
// Example: !Split ["|", "a|b|c"] becomes split("|", "a|b|c").
func ParseSplit(delimiter string, sourceString string) string {
	// Escape quotes in delimiter and source
	escapedDelimiter := strings.ReplaceAll(delimiter, "\"", "\\\"")
	escapedSource := strings.ReplaceAll(sourceString, "\"", "\\\"")

	return fmt.Sprintf("split(\"%s\", \"%s\")", escapedDelimiter, escapedSource)
}

// ReplaceStepFunctionsReferences converts variable references in Step Functions definitions
// from ${var.ResourceNameArn} to ${aws_lambda_function.ResourceName.arn}.
func ReplaceStepFunctionsReferences(definition string, resources cloudformation.Resources) string {
	// Pattern to match ${var.XXX} references
	varPattern := regexp.MustCompile(`\$\{var\.([^}]+)\}`)

	return varPattern.ReplaceAllStringFunc(definition, func(match string) string {
		// Extract the variable name (e.g., "TerminateEC2Arn" from "${var.TerminateEC2Arn}")
		varName := varPattern.FindStringSubmatch(match)[1]

		// Try to find a matching resource by checking if varName matches or is a prefix
		for resourceName, resource := range resources {
			// Check if varName starts with resourceName (e.g., "TerminateEC2Arn" starts with "TerminateEC2")
			// or if varName equals resourceName
			if strings.EqualFold(varName, resourceName) ||
				strings.HasPrefix(strings.ToLower(varName), strings.ToLower(resourceName)) {

				resourceType := ToTFName(resource.AWSCloudFormationType())

				// Determine the attribute based on the suffix
				lowerVarName := strings.ToLower(varName)
				var attribute string
				switch {
				case strings.HasSuffix(lowerVarName, "topicarn"):
					attribute = "arn"
				case strings.HasSuffix(lowerVarName, "arn"):
					attribute = "arn"
				case strings.HasSuffix(lowerVarName, "url"):
					attribute = "url"
				default:
					attribute = "id"
				}

				// Preserve the original resource name case
				return fmt.Sprintf("${%s.%s.%s}", resourceType, resourceName, attribute)
			}
		}

		// If no matching resource found, leave as variable reference
		return match
	})
}
