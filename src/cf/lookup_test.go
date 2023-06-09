package cf

import (
	"reflect"
	"testing"
	tftemplate "text/template"

	"github.com/awslabs/goformation/v7/cloudformation"
)

func TestParseResources(t *testing.T) {
	t.Parallel()

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
		// {},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := parseResources(tt.args.resources, tt.args.funcMap, tt.args.destination); (err != nil) != tt.wantErr {
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
