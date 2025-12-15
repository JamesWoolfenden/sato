package cf

import (
	"reflect"
	"testing"
	tftemplate "text/template"

	"github.com/awslabs/goformation/v7/cloudformation"
)

//goland:noinspection GoLinter
func TestParse(t *testing.T) {
	// Note: Not using t.Parallel() due to concurrent map writes bug in goformation library
	type args struct {
		file        string
		destination string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"minimal", args{"./testdata/min.cg.yaml", "./testdata/.terraform/output"}, false},
		{"minimal with params", args{"./testdata/minwithparams.cg.yaml", "./testdata/.terraform/output"}, false},
		{"not a file", args{"./testdata/nothing.cg.yaml", "./testdata/.terraform/output"}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Not using t.Parallel() here because goformation library is not thread-safe
			if err := Parse(tt.args.file, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//goland:noinspection GoLinter
func TestStringToMap(t *testing.T) {
	t.Parallel()

	type args struct {
		param cloudformation.Parameter
	}

	//Description := "Enter t2.micro, m1.small, or m1.large. Default is t2.micro."
	//result := Variable{
	//	Description: "",
	//	Type:        "map(string)",
	//	Default:     "{ \"t2.micro\" = }",
	//	Name:        "",
	//}

	tests := []struct {
		name string
		args args
		want Variable
	}{
		//{name: "pass",
		//	args: args{
		//		param: cloudformation.Parameter{
		//			Type:          "String",
		//			Description:   &Description,
		//			Default:       "{ pike=params }",
		//			AllowedValues: []interface{}{"t2.micro", "m1.small", "m1.large"},
		//		}},
		//	want: result},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := StringToMap(tt.args.param); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter
func TestWrite(t *testing.T) {
	t.Parallel()

	type args struct {
		output   string
		location string
		name     string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Pass", args{"gibberine", ".", "test"}, false},
		//{"Fail", args{"gibberine", "/usr/bin/local", "test"}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := Write(tt.args.output, tt.args.location, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//goland:noinspection GoLinter
func TestToTFName(t *testing.T) {
	t.Parallel()

	type args struct {
		CFN string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"pass", args{"AWS::SNS::Topic"}, "aws_sns_topic"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ToTFName(tt.args.CFN); got != tt.want {
				t.Errorf("ToTFName() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter
func TestReplaceVariables(t *testing.T) {
	t.Parallel()

	type args struct {
		str1 string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"pass", args{"Im ${good}"}, "Im ${var.good}"},
		{"do nothing", args{"Im ${good::aws}"}, "Im ${good::aws}"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ReplaceVariables(tt.args.str1); got != tt.want {
				t.Errorf("ReplaceVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter
func TestParseVariables(t *testing.T) {
	t.Parallel()

	type args struct {
		template    *cloudformation.Template
		funcMap     tftemplate.FuncMap
		destination string
	}

	tests := []struct {
		name    string
		args    args
		want    []Variable
		wantErr bool
	}{
		//{},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseVariables(tt.args.template, tt.args.funcMap, tt.args.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVariables() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter
func TestGetVariableType(t *testing.T) {
	t.Parallel()

	type args struct {
		param         cloudformation.Parameter
		myVariable    Variable
		DataResources []string
		m             map[string]bool
	}

	tests := []struct {
		name  string
		args  args
		want  []string
		want1 Variable
		want2 map[string]bool
	}{
		{
			name: "String type",
			args: args{
				param:         cloudformation.Parameter{Type: "String", Default: "test"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{provider},
			want1: Variable{Type: "string"},
			want2: map[string]bool{provider: true},
		},
		{
			name: "String type with bool default",
			args: args{
				param:         cloudformation.Parameter{Type: "String", Default: "true"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{provider},
			want1: Variable{Type: "bool"},
			want2: map[string]bool{provider: true},
		},
		{
			name: "Number type",
			args: args{
				param:         cloudformation.Parameter{Type: "Number"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{provider},
			want1: Variable{Type: "number"},
			want2: map[string]bool{provider: true},
		},
		{
			name: "CommaDelimitedList type",
			args: args{
				param:         cloudformation.Parameter{Type: "CommaDelimitedList"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{provider},
			want1: Variable{Type: "list(string)"},
			want2: map[string]bool{provider: true},
		},
		{
			name: "AWS::EC2::Subnet::Id type",
			args: args{
				param:         cloudformation.Parameter{Type: "AWS::EC2::Subnet::Id"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{dataSubnet, provider},
			want1: Variable{Type: "string"},
			want2: map[string]bool{dataSubnet: true, provider: true},
		},
		{
			name: "AWS::EC2::KeyPair::KeyName type",
			args: args{
				param:         cloudformation.Parameter{Type: "AWS::EC2::KeyPair::KeyName"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{dataKeyPair, provider},
			want1: Variable{Type: "string"},
			want2: map[string]bool{dataKeyPair: true, provider: true},
		},
		{
			name: "AWS::EC2::VPC::Id type",
			args: args{
				param:         cloudformation.Parameter{Type: "AWS::EC2::VPC::Id"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{dataVpc, provider},
			want1: Variable{Type: "string"},
			want2: map[string]bool{dataVpc: true, provider: true},
		},
		{
			name: "AWS::Region type",
			args: args{
				param:         cloudformation.Parameter{Type: "AWS::Region"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{dataRegion, provider},
			want1: Variable{Type: "string"},
			want2: map[string]bool{dataRegion: true, provider: true},
		},
		{
			name: "List<AWS::EC2::Subnet::Id> type",
			args: args{
				param:         cloudformation.Parameter{Type: "List<AWS::EC2::Subnet::Id>"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{provider},
			want1: Variable{Type: "list(string)"},
			want2: map[string]bool{provider: true},
		},
		{
			name: "List<AWS::EC2::AvailabilityZone::Name> type",
			args: args{
				param:         cloudformation.Parameter{Type: "List<AWS::EC2::AvailabilityZone::Name>"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{dataAvailabilityZone, provider},
			want1: Variable{Type: "list(string)"},
			want2: map[string]bool{dataAvailabilityZone: true, provider: true},
		},
		{
			name: "AWS::EC2::SecurityGroup::Id type",
			args: args{
				param:         cloudformation.Parameter{Type: "AWS::EC2::SecurityGroup::Id"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{dataSecurityGroup, provider},
			want1: Variable{Type: "string"},
			want2: map[string]bool{dataSecurityGroup: true, provider: true},
		},
		{
			name: "AWS::EC2::Image::Id type",
			args: args{
				param:         cloudformation.Parameter{Type: "AWS::EC2::Image::Id"},
				myVariable:    Variable{},
				DataResources: []string{},
				m:             make(map[string]bool),
			},
			want:  []string{provider},
			want1: Variable{Type: "string"},
			want2: map[string]bool{provider: true},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, got2 := GetVariableType(tt.args.param, tt.args.myVariable, tt.args.DataResources, tt.args.m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVariableType() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetVariableType() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GetVariableType() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

//goland:noinspection GoLinter
func TestGetVariableDefault(t *testing.T) {
	t.Parallel()

	type args struct {
		param      cloudformation.Parameter
		myVariable Variable
	}

	tests := []struct {
		name string
		args args
		want Variable
	}{
		{
			name: "String default value",
			args: args{
				param:      cloudformation.Parameter{Default: "my-default-value"},
				myVariable: Variable{Type: "string"},
			},
			want: Variable{Type: "string", Default: "\"my-default-value\""},
		},
		{
			name: "Boolean string default",
			args: args{
				param:      cloudformation.Parameter{Default: "true"},
				myVariable: Variable{Type: "bool"},
			},
			want: Variable{Type: "bool", Default: "true"},
		},
		{
			name: "Boolean value default",
			args: args{
				param:      cloudformation.Parameter{Default: true},
				myVariable: Variable{Type: "bool"},
			},
			want: Variable{Type: "bool", Default: "true"},
		},
		{
			name: "Numeric string default",
			args: args{
				param:      cloudformation.Parameter{Default: "42"},
				myVariable: Variable{},
			},
			want: Variable{Type: "number", Default: "42"},
		},
		{
			name: "Float default",
			args: args{
				param:      cloudformation.Parameter{Default: 3.14},
				myVariable: Variable{},
			},
			want: Variable{Type: "number", Default: "3.14"},
		},
		{
			name: "Empty list default",
			args: args{
				param:      cloudformation.Parameter{Default: []interface{}{}},
				myVariable: Variable{Type: "list(string)"},
			},
			want: Variable{Type: "list(string)", Default: "[]"},
		},
		{
			name: "Map string with equals",
			args: args{
				param:      cloudformation.Parameter{Default: "Name=MyApp"},
				myVariable: Variable{},
			},
			want: Variable{Type: "map(string)", Default: "{ \"Name\" = \"MyApp\"}"},
		},
		{
			name: "False boolean default",
			args: args{
				param:      cloudformation.Parameter{Default: false},
				myVariable: Variable{Type: "bool"},
			},
			want: Variable{Type: "bool", Default: "false"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := GetVariableDefault(tt.args.param, tt.args.myVariable); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVariableDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter
func TestReplaceDependant(t *testing.T) {
	t.Parallel()

	type args struct {
		str1 string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"pass", args{"AWS::Region::Something"}, "data.aws_region.current.name::Something"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ReplaceDependant(tt.args.str1); got != tt.want {
				t.Errorf("ReplaceDependant() = %v, want %v", got, tt.want)
			}
		})
	}
}
