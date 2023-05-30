package arm

import (
	"reflect"
	"testing"
)

func Test_lookup(t *testing.T) {
	t.Parallel()

	type args struct {
		myType string
	}

	test := azurermTemplateDeployment
	var empty []byte

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"Find", args{"Microsoft.Resources/deployments"}, test},
		{"Dont Find", args{"garbage"}, empty},
		{"Nil", args{""}, empty},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := lookup(tt.args.myType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lookup() = %v, want %v", got, tt.want)
			}
		})
	}
}
