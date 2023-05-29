package see

import (
	"reflect"
	"testing"
)

func TestLookup(t *testing.T) {
	t.Parallel()

	type args struct {
		resource string
	}

	result := "aws_appautoscaling_target"

	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{name: "Pass", args: args{"AWS::ApplicationAutoScaling::ScalableTarget"}, want: &result, wantErr: false},
		{name: "Fail", args: args{"AWS::Guff::Guffing"}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Lookup(tt.args.resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lookup() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lookup() got = %v, want %v", got, tt.want)
			}
		})
	}
}
