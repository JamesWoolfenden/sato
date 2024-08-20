package cf

import (
	"encoding/json"
	"html/template"
	"reflect"
	"strings"
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

	var funcMap = template.FuncMap{ //nolint:gochecknoglobals
		"ToUpper": strings.ToUpper,
		"Deref":   func(str *string) string { return *str },
		"Marshal": func(v interface{}) string {
			a, _ := json.Marshal(v) //nolint:errchkjson

			return string(a)
		},
		// Nild is a template function.
		"Nild": func(str *string, myDefault string) string {
			if str == nil || *str == "" {
				return myDefault
			}

			return *str
		},
		"Quote": func(target string) string {
			if (strings.Contains(target, "var.") || strings.Contains(target, "local.") ||
				(strings.Contains(target, "_") && strings.Contains(target, "."))) && (!strings.Contains(target, "${")) {
				return target
			}

			return "\"" + target + "\""
		},
	}

	//cut down template for a slimmed function map
	cloudFormation, _ := goformation.Open("../../examples/template.yaml")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Parsed", args{
			resources:   cloudFormation.Resources,
			funcMap:     funcMap,
			destination: "",
		}, false},
		{"empty function map", args{
			resources:   cloudFormation.Resources,
			funcMap:     template.FuncMap{},
			destination: "",
		}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := parseResources(
				tt.args.resources, tt.args.funcMap, tt.args.destination); (err != nil) != tt.wantErr {
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
