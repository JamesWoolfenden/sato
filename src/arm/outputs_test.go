package arm

import (
	"testing"
	"text/template"
)

func Test_parseOutputs(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := parseOutputs(tt.args.result, tt.args.funcMap, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("parseOutputs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
