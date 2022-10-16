package sato

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/awslabs/goformation/v7"
	"github.com/awslabs/goformation/v7/cloudformation/tags"
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
			"AWS::SNS::Topic":                       awsSNSTopic,
			"AWS::IAM::Role":                        awsIamRole,
			"AWS::EC2::Route":                       awsRoute,
			"AWS::EC2::RouteTable":                  awsRouteTable,
			"AWS::EC2::NatGateway":                  awsNatGateway,
			"AWS::EC2::VPCGatewayAttachment":        awsVpnGatewayAttachment,
			"AWS::EC2::NetworkAclEntry":             awsNetworkAclRule,
			"AWS::EC2::NetworkAcl":                  awsNetworkAcl,
			"AWS::EC2::EIP":                         awsEIP,
			"AWS::EC2::SubnetRouteTableAssociation": awsRouteTableAssociation,
			"AWS::EC2::Subnet":                      awsSubnet,
			"AWS::Logs::LogGroup":                   awsCloudwatchLogGroup,
			"AWS::EC2::VPCDHCPOptionsAssociation":   awsVpcDhcpOptionsAssociation,
			"AWS::EC2::DHCPOptions":                 awsVpcDhcpOptions,
			"AWS::EC2::SubnetNetworkAclAssociation": awsNetworkAclAssociation,
			"AWS::EC2::FlowLog":                     awsFlowLog,
			"AWS::EC2::VPCEndpoint":                 awsVpcEndpoint,
			"AWS::EC2::InternetGateway":             awsInternetGateway,
			"AWS::EC2::VPC":                         awsVpc,
		}

		var myContent []byte
		if TFLookup[myType] != nil {
			myContent = TFLookup[myType].([]byte)
		} else {
			log.Printf("%s not found", myType)
			continue
		}

		//needs to pivot on policy template from resource
		tmpl, err := tftemplate.New("test").Funcs(tftemplate.FuncMap{
			"Deref": func(str *string) string { return *str },
			"Marshal": func(v interface{}) string {
				a, _ := json.Marshal(v)
				return string(a)
			},
			"Tags": func(v []tags.Tag) string {
				var temp string
				for _, item := range v {
					if item.Key != "" {
						temp = temp + "\"" + item.Key + "\"" + "=" + "\"" + item.Value + "\"" + "\n"
					}
				}
				return temp
			},
		}).Parse(string(myContent))

		if err != nil {
			panic(err)
		}

		_ = tmpl.Execute(&output, M{
			"resource": resource,
			"item":     item,
		})
		err = Write(output.String(), destination, fmt.Sprint(ToTFName(myType), ".", strings.ToLower(item)))
		if err != nil {
			log.Fatal(err)
		}
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

func Test() {

	var fudge []tags.Tag
	var temp tags.Tag
	temp.Value = "timmy"

	fudge = append(fudge, temp)

	log.Print(fudge)
}
