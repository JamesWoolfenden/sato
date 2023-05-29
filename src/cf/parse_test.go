package cf_test

import (
	"reflect"
	sato "sato/src/cf"
	"testing"
	tftemplate "text/template"

	"github.com/awslabs/goformation/v7/cloudformation"
)

//goland:noinspection GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter
func TestParse(t *testing.T) {
	t.Parallel()

	type args struct {
		file        string
		destination string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		//{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := sato.Parse(tt.args.file, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//goland:noinspection GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter
func TestStringToMap(t *testing.T) {
	t.Parallel()

	type args struct {
		param      cloudformation.Parameter
		myVariable sato.Variable
	}

	tests := []struct {
		name string
		args args
		want sato.Variable
	}{
		//{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.StringToMap(tt.args.param, tt.args.myVariable); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter
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
		//{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := sato.Write(tt.args.output, tt.args.location, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//goland:noinspection GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter
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
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.ToTFName(tt.args.CFN); got != tt.want {
				t.Errorf("ToTFName() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter,GoLinter,GoLinter,GoLinter,GoLinter
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
			if got := sato.ReplaceVariables(tt.args.str1); got != tt.want {
				t.Errorf("ReplaceVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter
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
		want    []sato.Variable
		wantErr bool
	}{
		//{},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := sato.ParseVariables(tt.args.template, tt.args.funcMap, tt.args.destination)
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

//goland:noinspection GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter
func TestGetVariableType(t *testing.T) {
	t.Parallel()

	type args struct {
		param         cloudformation.Parameter
		myVariable    sato.Variable
		DataResources []string
		m             map[string]bool
	}

	tests := []struct {
		name  string
		args  args
		want  []string
		want1 sato.Variable
		want2 map[string]bool
	}{
		//{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, got2 := sato.GetVariableType(tt.args.param, tt.args.myVariable, tt.args.DataResources, tt.args.m)
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

//goland:noinspection GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter
func TestGetVariableDefault(t *testing.T) {
	t.Parallel()

	type args struct {
		param      cloudformation.Parameter
		myVariable sato.Variable
	}

	tests := []struct {
		name string
		args args
		want sato.Variable
	}{
		//{},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.GetVariableDefault(tt.args.param, tt.args.myVariable); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVariableDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

//goland:noinspection GoLinter,GoLinter,GoLinter,GoLinter,GoLinter,GoLinter
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
		{},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.ReplaceDependant(tt.args.str1); got != tt.want {
				t.Errorf("ReplaceDependant() = %v, want %v", got, tt.want)
			}
		})
	}
}
