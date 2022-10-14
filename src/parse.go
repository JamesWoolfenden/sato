package sato

import (
	"bytes"
	"fmt"
	"github.com/awslabs/goformation/v7"
	"log"
	"os"
	"path/filepath"
	"strings"
	tftemplate "text/template"
)

type M map[string]interface{}

func Parse(file string, destination string) {
	// Open a template from file (can be JSON or YAML)
	fileAbs, err := filepath.Abs(file)
	if err != nil {
		log.Fatal(err)
	}

	template, err := goformation.Open(fileAbs)
	if err != nil {
		log.Fatalf("There was an error processing the template: %s", err)
	}

	resources := template.Resources
	for item, resource := range resources {
		var output bytes.Buffer

		myType := resources[item].AWSCloudFormationType()

		TFLookup := map[string]interface{}{
			"AWS::SNS::Topic": awsSNSTopic,
			"AWS::IAM::Role":  awsIamRole,
		}

		var myContent []byte
		if TFLookup[myType] != nil {
			myContent = TFLookup[myType].([]byte)
		} else {
			log.Printf("%s not found", myType)
			continue
		}

		//needs to pivot on policy template from resource
		tmpl, err := tftemplate.New("test").Parse(string(myContent))
		if err != nil {
			panic(err)
		}

		_ = tmpl.Execute(&output, M{
			"resource": resource,
			"item":     item,
		})
		Write(output.String(), destination, fmt.Sprint(ToTFName(myType), ".", strings.ToLower(item)))
	}
}

// Write out terraform
func Write(output string, location string, name string) error {

	newPath, _ := filepath.Abs(location)
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		return err
	}

	d1 := []byte(output)

	destination, _ := filepath.Abs(fmt.Sprint(location, "/", name, ".tf"))
	err = os.WriteFile(destination, d1, 0644)

	if err != nil {
		return err
	}

	return nil
}

func ToTFName(CFN string) string {
	return strings.ToLower(strings.ReplaceAll(CFN, "::", "_"))
}
