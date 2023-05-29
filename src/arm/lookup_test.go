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

	tests := []struct {
		name string
		args args
		want []byte
	}{
		{},
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
