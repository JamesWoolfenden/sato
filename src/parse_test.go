package sato

import (
	"reflect"
	"testing"
	tftemplate "text/template"

	"github.com/awslabs/goformation/v7/cloudformation"
)

func TestParse(t *testing.T) {
	type args struct {
		file        string
		destination string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Parse(tt.args.file, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringToMap(t *testing.T) {
	type args struct {
		param      cloudformation.Parameter
		myVariable Variable
	}
	tests := []struct {
		name string
		args args
		want Variable
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToMap(tt.args.param, tt.args.myVariable); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseResources(t *testing.T) {
	type args struct {
		resources   cloudformation.Resources
		funcMap     tftemplate.FuncMap
		destination string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseResources(tt.args.resources, tt.args.funcMap, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("ParseResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrite(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Write(tt.args.output, tt.args.location, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestToTFName(t *testing.T) {
	type args struct {
		CFN string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToTFName(tt.args.CFN); got != tt.want {
				t.Errorf("ToTFName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceVariables(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceVariables(tt.args.str1); got != tt.want {
				t.Errorf("ReplaceVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseVariables(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func TestGetVariableType(t *testing.T) {
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetVariableType(tt.args.param, tt.args.myVariable, tt.args.DataResources, tt.args.m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVariableType() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("GetVariableType() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetVariableDefault(t *testing.T) {
	type args struct {
		param      cloudformation.Parameter
		myVariable Variable
	}
	tests := []struct {
		name string
		args args
		want Variable
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVariableDefault(tt.args.param, tt.args.myVariable); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVariableDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReplaceDependant(t *testing.T) {
	type args struct {
		str1 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceDependant(tt.args.str1); got != tt.want {
				t.Errorf("ReplaceDependant() = %v, want %v", got, tt.want)
			}
		})
	}
}
