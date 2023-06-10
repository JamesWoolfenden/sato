package arm_test

import (
	"sato/src/arm"
	"testing"
	"text/template"
)

func Test_parseData(t *testing.T) {
	t.Parallel()

	type args struct {
		result      map[string]interface{}
		funcMap     template.FuncMap
		destination string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		//		{},
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
