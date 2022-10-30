package sato

import (
	"reflect"
	"testing"
	tftemplate "text/template"

	"github.com/awslabs/goformation/v7/cloudformation"
)

func TestParseResources(t *testing.T) {
	type args struct {
		resources   cloudformation.Resources
		funcMap     tftemplate.FuncMap
		destination string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseResources(tt.args.resources, tt.args.funcMap, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("ParseResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLookup(t *testing.T) {
	type args struct {
		myType string
	}

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"found", args{"AWS::SNS::Topic"}, awsSNSTopic},
		{"not found", args{"AWS::SNS::Balderdash"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lookup(tt.args.myType)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lookup() = %v, want %v", got, tt.want)
			}
		})
	}
}
