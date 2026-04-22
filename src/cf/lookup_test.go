package cf

import (
	"reflect"
	"testing"
	tftemplate "text/template"

	"github.com/awslabs/goformation/v7"
	"github.com/awslabs/goformation/v7/cloudformation"
)

func TestParseResources(t *testing.T) {
	t.Parallel()

	type args struct {
		resources   cloudformation.Resources
		funcMap     tftemplate.FuncMap
		destination string
	}

	cloudFormation, _ := goformation.Open("../../examples/template.yaml")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Parsed", args{
			resources:   cloudFormation.Resources,
			funcMap:     funcMap,
			destination: ".",
		}, false},
		{"empty function map", args{
			resources:   cloudFormation.Resources,
			funcMap:     tftemplate.FuncMap{},
			destination: "",
		}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := parseResources(
				tt.args.resources, tt.args.funcMap, tt.args.destination, options{}); (err != nil) != tt.wantErr {
				t.Errorf("parseResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLookup(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := lookup(tt.args.myType)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lookup() = %v, want %v", got, tt.want)
			}
		})
	}
}
