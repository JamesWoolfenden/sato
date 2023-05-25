package arm

import (
	"reflect"
	"testing"
)

func Test_lookup(t *testing.T) {
	type args struct {
		myType string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lookup(tt.args.myType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lookup() = %v, want %v", got, tt.want)
			}
		})
	}
}
