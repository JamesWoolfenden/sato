package cf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	tftemplate "text/template"

	"github.com/awslabs/goformation/v7"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/rs/zerolog/log"
)

const typeNumber = "number"
const typeListString = "list(string)"
const typeString = "string"

// M is a wrapper object to help pass in multiple object to template.
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

// Parse turn CFN into Terraform.
func Parse(file string, destination string) error {
	// Open a template from file (can be JSON or YAML)
	fileAbs, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	template, err := goformation.Open(fileAbs)
	if err != nil {
		log.Fatal().Err(err).Msg("Parse failure")
	}

	funcMap := tftemplate.FuncMap{
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
		"Marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)

			return string(a)
		},
		"Split":        Split,
		"SplitOn":      SplitOn,
		"Replace":      Replace,
		"Tags":         Tags,
		"RandomString": RandomString,
		"Map":          Map,
		"Snake":        Snake,
		"Kebab":        Kebab,
		"ZipFile":      Zipfile,
	}
	_, err = ParseVariables(template, funcMap, destination)

	if err != nil {
		return err
	}

	err = parseResources(template.Resources, funcMap, destination)
	if err != nil {
		return err
	}

	return nil
}

// ParseVariables convert CFN Parameters into terraform variables
func ParseVariables(
	template *cloudformation.Template,
	funcMap tftemplate.FuncMap,
	destination string,
) ([]Variable, error) {
	var All string

	var (
		myMap         = make(map[string]bool)
		DataResources []string
	)

	var myVariables []Variable

	for Name, param := range template.Parameters {
		var myVariable Variable

		DataResources, myVariable, myMap = GetVariableType(param, myVariable, DataResources, myMap)
		myVariable.Description = strings.ReplaceAll(*param.Description, "${", "$${")
		myVariable.Name = Name
		myVariable = GetVariableDefault(param, myVariable)

		var output bytes.Buffer

		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(variableFile))
		if err != nil {
			return nil, err
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
		return nil, err
	}

	err = Write(strings.Join(DataResources, "\n"), destination, "data")
	if err != nil {
		return nil, err
	}

	return myVariables, nil
}

// GetVariableType determines variable types.
func GetVariableType(
	param cloudformation.Parameter,
	myVariable Variable,
	dataResources []string,
	myMap map[string]bool) ([]string, Variable, map[string]bool) {
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
					myVariable = StringToMap(param, myVariable)
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
func StringToMap(param cloudformation.Parameter, myVariable Variable) Variable {
	temp := strings.Split(param.Default.(string), "=")

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

// Write out Terraform
func Write(output string, location string, name string) error {
	if output != "" {
		newPath, _ := filepath.Abs(location)
		err := os.MkdirAll(newPath, os.ModePerm)

		if err != nil {
			return err
		}

		d1 := []byte(output)

		destination, _ := filepath.Abs(fmt.Sprint(location, "/", name, ".tf"))
		err = os.WriteFile(destination, d1, 0o644)
		log.Info().Msgf("Created %s", destination)

		if err != nil {
			return err
		}
	}

	return nil
}

// ToTFName creates a Terraform resource name from a CFN type (approximates)
func ToTFName(CFN string) string {
	return strings.ToLower(strings.ReplaceAll(CFN, "::", "_"))
}

// ReplaceVariables looks to see if u can translate CFN vars into terraform
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
