package arm_test

import (
	"reflect"
	"sato/src/arm"
	"testing"
	"text/template"
)

func Test_parseResources(t *testing.T) {
	t.Parallel()

	type args struct {
		result      map[string]interface{}
		funcMap     template.FuncMap
		destination string
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		//{ "Pass", args{}, ResultPass, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := arm.ParseResources(tt.args.result, tt.args.funcMap, tt.args.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseResources() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseResources() got = %v, want %v", got, tt.want)
			}
		})
	}
}
