package sato

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	tftemplate "text/template"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/awslabs/goformation/v7"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/tags"
)

// M is a wrapper object to help pass in multiple object to template
type M map[string]interface{}

// Variable describes a terraform variable
type Variable struct {
	Description string
	Type        string
	Default     string
	Name        string
}

// Parse turn CFN into Terraform
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
		"Array":        array,
		"ArrayReplace": arrayReplace,
		"Contains":     contains,
		"Sprint":       sprint,
		"Decode64":     decode64,
		"Boolean":      boolean,
		"Dequote":      dequote,
		"Quote":        quote,
		"Demap": func(str string) []string {
			str = strings.Replace(str, "{", "", -1)
			str = strings.Replace(str, "}", "", -1)
			str = strings.Replace(str, "\"", "", -1)
			str = strings.Replace(str, ":", "", -1)
			str = strings.Replace(str, " ", "", -1)
			return strings.Split(str, ",")
		},
		"ToUpper": strings.ToUpper,
		"ToLower": lower,
		"Deref":   func(str *string) string { return *str },
		"Nil":     nill,
		"Nild":    nild,
		"Marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"Split":   split,
		"Replace": replace,
		"Tags": func(v []tags.Tag) string {
			var temp string
			for _, item := range v {
				if item.Key != "" {
					temp = temp + "\t\"" + item.Key + "\"" + "=" + "\"" + item.Value + "\"" + "\n"
				}
			}
			return temp
		},
		"RandomString": func(n int) string {
			rand.Seed(time.Now().UnixNano())
			var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			b := make([]rune, n)
			for i := range b {
				b[i] = letters[rand.Intn(len(letters))]
			}
			return string(b)
		},
		"Map": func(myMap map[string]string) string {
			result := "{ \n"
			for item, stuff := range myMap {
				result = result + "\t\"" + item + "\"" + "=" + "\"" + stuff + "\"\n"
			}
			result = result + " }"

			return result
		},
		"Snake":   snake,
		"Kebab":   kebab,
		"ZipFile": zipfile,
	}
	_, err = ParseVariables(template, funcMap, destination)
	if err != nil {
		return err
	}

	err = ParseResources(template.Resources, funcMap, destination)
	if err != nil {
		return err
	}

	return nil
}

// ParseVariables convert CFN Parameters into terraform variables
func ParseVariables(template *cloudformation.Template, funcMap tftemplate.FuncMap, destination string) ([]Variable, error) {
	var All string

	var (
		m             = make(map[string]bool)
		DataResources []string
	)

	var myVariables []Variable
	for Name, param := range template.Parameters {
		var myVariable Variable

		DataResources, myVariable, m = GetVariableType(param, myVariable, DataResources, m)
		myVariable.Description = strings.Replace(param.Description, "${", "$${", -1)
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
		All = All + output.String()
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

// GetVariableType determines variable types
func GetVariableType(param cloudformation.Parameter, myVariable Variable, DataResources []string, m map[string]bool) ([]string, Variable, map[string]bool) {
	switch param.Type {
	case "String":
		if param.Default == "false" || param.Default == "true" || param.Default == true || param.Default == false {
			myVariable.Type = "bool"
		} else {
			myVariable.Type = strings.ToLower(param.Type)
		}
	case "CommaDelimitedList":
		myVariable.Type = "list(string)"
	case "List<AWS::EC2::AvailabilityZone::Name>":
		myVariable.Type = "list(string)"
		DataResources, m = add(dataAvailabilityZone, DataResources, m)
	case "AWS::EC2::Subnet::Id":
		myVariable.Type = "string"
		DataResources, m = add(dataSubnet, DataResources, m)
	case "AWS::EC2::KeyPair::KeyName":
		myVariable.Type = "string"
		DataResources, m = add(dataKeyPair, DataResources, m)
	case "AWS::EC2::VPC::Id", "List<AWS::EC2::VPC::Id>":
		myVariable.Type = "string"
		DataResources, m = add(dataVpc, DataResources, m)
	case "AWS::EC2::SecurityGroup::Id":
		myVariable.Type = "string"
		DataResources, m = add(dataSecurityGroup, DataResources, m)
	case "AWS::EC2::Image::Id", "AWS::SSM::Parameter::Value<AWS::EC2::Image::Id>":
		myVariable.Type = "string"
	case "AWS::Region":
		myVariable.Type = "string"
		DataResources, m = add(dataRegion, DataResources, m)
	case "List<AWS::EC2::Subnet::Id>":
		myVariable.Type = "list(string)"
	case "Number":
		myVariable.Type = "number"
	default:
		log.Info().Msgf("Variable %s", param.Type)
	}

	DataResources, m = add(provider, DataResources, m)
	return DataResources, myVariable, m
}

// GetVariableDefault determines a variables default value
func GetVariableDefault(param cloudformation.Parameter, myVariable Variable) Variable {
	switch param.Default.(type) {
	case string:
		_, err := strconv.Atoi(param.Default.(string))
		if err == nil {
			myVariable.Type = "number"
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
		myVariable.Type = "number"
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

// StringToMap converts maps in strings(for tags)
func StringToMap(param cloudformation.Parameter, myVariable Variable) Variable {
	temp := strings.Split(param.Default.(string), "=")
	var myMap string
	for n := 0; n < len(temp); n++ {
		if n == 0 {
			myMap = myMap + "{ "
		}
		if n%2 == 0 {
			myMap = myMap + "\"" + temp[n] + "\" = "
		} else {
			myMap = myMap + "\"" + temp[n] + "\""
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
		err = os.WriteFile(destination, d1, 0644)
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
		if strings.Contains(target, "::") {
		} else {
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
