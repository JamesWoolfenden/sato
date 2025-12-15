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

//goland:noinspection GoLinter
func TestParseSub(t *testing.T) {
	t.Parallel()

	type args struct {
		input      string
		parameters map[string]cloudformation.Parameter
	}

	params := map[string]cloudformation.Parameter{
		"BucketName":   {Type: "String"},
		"Environment":  {Type: "String"},
		"DatabaseName": {Type: "String"},
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Simple parameter substitution",
			args: args{
				input:      "arn:aws:s3:::${BucketName}",
				parameters: params,
			},
			want: "arn:aws:s3:::${var.bucketname}",
		},
		{
			name: "AWS::Region pseudo-parameter",
			args: args{
				input:      "arn:aws:lambda:${AWS::Region}:123456789012:function:MyFunction",
				parameters: params,
			},
			want: "arn:aws:lambda:${data.aws_region.current.name}:123456789012:function:MyFunction",
		},
		{
			name: "AWS::AccountId pseudo-parameter",
			args: args{
				input:      "arn:aws:iam::${AWS::AccountId}:role/MyRole",
				parameters: params,
			},
			want: "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/MyRole",
		},
		{
			name: "Multiple substitutions",
			args: args{
				input:      "${Environment}-${DatabaseName}-db",
				parameters: params,
			},
			want: "${var.environment}-${var.databasename}-db",
		},
		{
			name: "Resource attribute reference",
			args: args{
				input:      "https://${MyBucket.DomainName}/index.html",
				parameters: make(map[string]cloudformation.Parameter),
			},
			want: "https://${mybucket.domainname}/index.html",
		},
		{
			name: "Mix of parameter and pseudo-parameter",
			args: args{
				input:      "${BucketName}-${AWS::Region}",
				parameters: params,
			},
			want: "${var.bucketname}-${data.aws_region.current.name}",
		},
		{
			name: "Plain text without substitution",
			args: args{
				input:      "no variables here",
				parameters: params,
			},
			want: "no variables here",
		},
		{
			name: "AWS::StackName pseudo-parameter",
			args: args{
				input:      "${AWS::StackName}-resources",
				parameters: params,
			},
			want: "${var.stack_name}-resources",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ParseSub(tt.args.input, tt.args.parameters); got != tt.want {
				t.Errorf("ParseSub() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter
func TestParseJoin(t *testing.T) {
	t.Parallel()

	type args struct {
		delimiter string
		values    []string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Join with hyphen delimiter",
			args: args{
				delimiter: "-",
				values:    []string{"web", "server", "01"},
			},
			want: `join("-", ["web", "server", "01"])`,
		},
		{
			name: "Join with empty delimiter",
			args: args{
				delimiter: "",
				values:    []string{"a", "b", "c"},
			},
			want: `join("", ["a", "b", "c"])`,
		},
		{
			name: "Join single value",
			args: args{
				delimiter: ",",
				values:    []string{"single"},
			},
			want: `join(",", ["single"])`,
		},
		{
			name: "Join with special characters",
			args: args{
				delimiter: "::",
				values:    []string{"arn", "aws", "s3"},
			},
			want: `join("::", ["arn", "aws", "s3"])`,
		},
		{
			name: "Join empty array",
			args: args{
				delimiter: "-",
				values:    []string{},
			},
			want: `join("-", [])`,
		},
		{
			name: "Join with values containing quotes",
			args: args{
				delimiter: ",",
				values:    []string{`value"with"quotes`, "normal"},
			},
			want: `join(",", ["value\"with\"quotes", "normal"])`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ParseJoin(tt.args.delimiter, tt.args.values); got != tt.want {
				t.Errorf("ParseJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter
func TestParseSplit(t *testing.T) {
	t.Parallel()

	type args struct {
		delimiter    string
		sourceString string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Split with pipe delimiter",
			args: args{
				delimiter:    "|",
				sourceString: "a|b|c",
			},
			want: `split("|", "a|b|c")`,
		},
		{
			name: "Split with comma delimiter",
			args: args{
				delimiter:    ",",
				sourceString: "item1,item2,item3",
			},
			want: `split(",", "item1,item2,item3")`,
		},
		{
			name: "Split with colon delimiter",
			args: args{
				delimiter:    ":",
				sourceString: "arn:aws:s3:::my-bucket",
			},
			want: `split(":", "arn:aws:s3:::my-bucket")`,
		},
		{
			name: "Split with multi-char delimiter",
			args: args{
				delimiter:    "::",
				sourceString: "part1::part2::part3",
			},
			want: `split("::", "part1::part2::part3")`,
		},
		{
			name: "Split empty string",
			args: args{
				delimiter:    ",",
				sourceString: "",
			},
			want: `split(",", "")`,
		},
		{
			name: "Split with quotes in source",
			args: args{
				delimiter:    ",",
				sourceString: `value"with"quotes,normal`,
			},
			want: `split(",", "value\"with\"quotes,normal")`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ParseSplit(tt.args.delimiter, tt.args.sourceString); got != tt.want {
				t.Errorf("ParseSplit() = %v, want %v", got, tt.want)
			}
		})
	}
}
