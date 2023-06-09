package arm_test

import (
	"reflect"
	"sato/src/arm"
	"testing"
	"text/template"
)

func Test_parseVariables(t *testing.T) {
	t.Parallel()

	type args struct {
		result      map[string]interface{}
		funcMap     template.FuncMap
		destination string
	}

	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		//{},
	}
	for _, tt := range tests {
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
