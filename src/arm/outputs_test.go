package arm_test

import (
	"encoding/json"
	"strings"
	"testing"
	"text/template"

	"sato/src/arm"
	"sato/src/tfgen"
)

var badFunk template.FuncMap

var funcMap = template.FuncMap{ //nolint:gochecknoglobals
	"Array":        tfgen.Array,
	"ArrayReplace": tfgen.ArrayReplace,
	"Contains":     tfgen.Contains,
	"Sprint":       tfgen.Sprint,
	"Decode64":     tfgen.Decode64,
	"Boolean":      tfgen.Boolean,
	"Dequote":      tfgen.Dequote,
	"Quote":        tfgen.Quote,
	"Demap":        tfgen.Demap,
	"ToUpper":      strings.ToUpper,
	"ToLower":      tfgen.Lower,
	"Deref":        func(str *string) string { return *str },
	"Nil":          tfgen.Nill,
	"Nild":         tfgen.Nild,
	"Marshal": func(v interface{}) string {
		a, _ := json.Marshal(v) //nolint:errchkjson

		return string(a)
	},
	"Split":        tfgen.Split,
	"SplitOn":      tfgen.SplitOn,
	"Replace":      tfgen.Replace,
	"Tags":         arm.Tags,
	"RandomString": tfgen.RandomString,
	"Map":          tfgen.Map,
	"Snake":        tfgen.Snake,
	"Kebab":        tfgen.Kebab,
	"ZipFile":      tfgen.Zipfile,
	"Uuid":         arm.UUID,
}

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
