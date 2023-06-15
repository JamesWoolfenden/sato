package arm_test

import (
	"reflect"
	"testing"
	"text/template"

	"sato/src/arm"
)

func Test_parseVariables(t *testing.T) {
	t.Parallel()

	type args struct {
		result      map[string]interface{}
		funcMap     template.FuncMap
		destination string
	}

	empty := make(map[string]interface{})
	emptyOutputs := make(map[string]interface{})
	emptyOutputs["outputs"] = make(map[string]interface{})

	results := make(map[string]interface{})
	outputs := make(map[string]interface{})
	entry := map[string]interface{}{
		"type": "string",
		"value": "[reference(resourceId('Microsoft.Network/publicIPAddresses'," +
			" parameters('publicIpName')), '2022-05-01').dnsSettings.fqdn]",
	}

	outputs["hostname"] = entry
	results["outputs"] = outputs
	wants := make(map[string]interface{})
	// emptySlice = append(emptySlice, empty)

	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{"empty", args{empty, funcMap, "test-output"}, wants, false},

		{"emptyOutputs", args{emptyOutputs, funcMap, "test-output"}, nil, false},
		{"Outputs", args{results, funcMap, "test-output"}, nil, false},
		{"Bad funk", args{results, badFunk, "test-output"}, nil, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := arm.ParseVariables(tt.args.result, tt.args.funcMap, tt.args.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVariables() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseVariables() got = %v, want %v", got, tt.want)
			}
		})
	}
}
