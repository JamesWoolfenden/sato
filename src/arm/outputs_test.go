package arm_test

import (
	"testing"
	"text/template"

	"sato/src/arm"
	"sato/src/tfgen"
)

var badFunk template.FuncMap

var funcMap = tfgen.FuncMap(template.FuncMap{
	"Enabled": arm.Enabled,
	"NotNil":  arm.NotNil,
	"Set":     arm.ArrayToString,
	"Tags":    arm.Tags,
	"Uuid":    arm.UUID,
})

func Test_parseOutputs(t *testing.T) {
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

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{empty, funcMap, "test-output"}, false},
		{"emptyOutputs", args{emptyOutputs, funcMap, "test-output"}, false},
		{"Outputs", args{results, funcMap, "test-output"}, false},
		{"Bad funk", args{results, badFunk, "test-output"}, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := arm.ParseOutputs(tt.args.result, tt.args.funcMap, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
