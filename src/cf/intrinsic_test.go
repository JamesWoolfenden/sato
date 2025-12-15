package cf

import (
	"testing"

	"github.com/awslabs/goformation/v7/cloudformation"
)

//goland:noinspection GoLinter
func TestParseRef(t *testing.T) {
	t.Parallel()

	type args struct {
		input      string
		parameters map[string]cloudformation.Parameter
		resources  cloudformation.Resources
	}

	// Create test parameters
	params := map[string]cloudformation.Parameter{
		"VpcId":    {Type: "String"},
		"SubnetId": {Type: "String"},
		"KeyName":  {Type: "String"},
	}

	// Empty resources for simple testing
	resources := cloudformation.Resources{}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Parameter reference",
			args: args{
				input:      "${VpcId}",
				parameters: params,
				resources:  resources,
			},
			want: "var.vpcid",
		},
		{
			name: "Ref:: prefix with parameter",
			args: args{
				input:      "${Ref::SubnetId}",
				parameters: params,
				resources:  resources,
			},
			want: "var.subnetid",
		},
		{
			name: "Multiple references",
			args: args{
				input:      "vpc=${VpcId}, subnet=${SubnetId}",
				parameters: params,
				resources:  resources,
			},
			want: "vpc=var.vpcid, subnet=var.subnetid",
		},
		{
			name: "Unknown reference defaults to variable",
			args: args{
				input:      "${UnknownRef}",
				parameters: make(map[string]cloudformation.Parameter),
				resources:  cloudformation.Resources{},
			},
			want: "var.unknownref",
		},
		{
			name: "Plain text without references",
			args: args{
				input:      "no references here",
				parameters: params,
				resources:  resources,
			},
			want: "no references here",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ParseRef(tt.args.input, tt.args.parameters, tt.args.resources); got != tt.want {
				t.Errorf("ParseRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter
func TestParseGetAtt(t *testing.T) {
	t.Parallel()

	type args struct {
		input string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ARN attribute",
			args: args{input: "${MyFunction.Arn}"},
			want: "${myfunction.arn}",
		},
		{
			name: "DNS name attribute",
			args: args{input: "${MyLoadBalancer.DNSName}"},
			want: "${myloadbalancer.dns_name}",
		},
		{
			name: "Queue URL attribute",
			args: args{input: "${MyQueue.QueueUrl}"},
			want: "${myqueue.url}",
		},
		{
			name: "Multiple GetAtt in string",
			args: args{input: "arn=${MyRole.Arn}, name=${MyRole.Name}"},
			want: "arn=${myrole.arn}, name=${myrole.name}",
		},
		{
			name: "CamelCase attribute conversion",
			args: args{input: "${MyResource.HostedZoneId}"},
			want: "${myresource.hosted_zone_id}",
		},
		{
			name: "Simple ID attribute",
			args: args{input: "${MyResource.Id}"},
			want: "${myresource.id}",
		},
		{
			name: "Plain text without GetAtt",
			args: args{input: "no getatt here"},
			want: "no getatt here",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ParseGetAtt(tt.args.input); got != tt.want {
				t.Errorf("ParseGetAtt() = %v, want %v", got, tt.want)
			}
		})
	}
}
