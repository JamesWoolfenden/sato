package main

import (
	"bytes"
	_ "embed" //required for embed
	"fmt"
	"github.com/awslabs/goformation/v7"
	"log"
	tftemplate "text/template"
)

//go:embed aws_sns_topic.policy.template
var awsSNSTopic []byte

//go:embed aws_iam_role.template
var awsIamRole []byte

type M map[string]interface{}

func main() {

	// Open a template from file (can be JSON or YAML)
	template, err := goformation.Open("template.yaml")
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
		fmt.Print(output.String())
	}

}
