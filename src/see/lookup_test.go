package see

import (
	"reflect"
	"testing"
)

func TestLookup(t *testing.T) {
	t.Parallel()

	type args struct {
		resource string
		reverse  bool
	}

	result := "aws_appautoscaling_target"
	myServiceBus := "azurerm_servicebus_namespace"
	reverse := "aws::efs::filesystem"

	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{name: "Pass", args: args{"AWS::ApplicationAutoScaling::ScalableTarget", false}, want: &result, wantErr: false},
		{name: "Pass", args: args{"Microsoft.ServiceBus/namespaces/", false}, want: &myServiceBus, wantErr: false},
		{name: "Fail", args: args{"AWS::Guff::Guffing", false}, want: nil, wantErr: true},
		{name: "Reverse", args: args{resource: "aws_efs_file_system", reverse: true}, want: &reverse, wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Lookup(tt.args.resource, tt.args.reverse)
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
