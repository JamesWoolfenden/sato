package arm_test

import (
	"testing"
	"text/template"

	"sato/src/arm"
)

func Test_parseData(t *testing.T) {
	t.Parallel()

	type args struct {
		result      map[string]interface{}
		funcMap     template.FuncMap
		destination string
	}

	var badFunk template.FuncMap

	empty := make(map[string]interface{})
	emptyData := make(map[string]interface{})
	emptyData["data"] = make(map[string]interface{})

	results := make(map[string]interface{})
	data := make(map[string]interface{})

	data["resource_group"] = true
	results["data"] = data

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{empty, funcMap, "test-output"}, false},
		{"emptyData", args{emptyData, funcMap, "test-output"}, false},
		{"data", args{results, funcMap, "test-output"}, false},
		{"bad funk", args{results, badFunk, "test-output"}, true},
		// {"bad destination", args{results, funcMap, "/usr/bin/nowhere"}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := arm.ParseData(tt.args.result, tt.args.funcMap, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("ParseData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
