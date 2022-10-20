package sato

import (
	"reflect"
	"testing"
)

func Test_replace(t *testing.T) {
	type args struct {
		input string
		from  string
		to    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"pass", args{"hankey", "k", "l"}, "hanley"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replace(tt.args.input, tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_add(t *testing.T) {
	type args struct {
		s string
		a []string
		m map[string]bool
	}
	var added []string
	m := map[string]bool{"test": false}
	after := []string{"test"}
	afterMap := map[string]bool{"test": true}
	exists := []string{"test"}

	tests := []struct {
		name  string
		args  args
		want  []string
		want1 map[string]bool
	}{
		{"add-empty", args{"test", added, m}, after, afterMap},
		{"no-change", args{"test", exists, afterMap}, after, afterMap},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := add(tt.args.s, tt.args.a, tt.args.m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("add() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("add() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_split(t *testing.T) {
	type args struct {
		source    string
		separator string
	}
	result := []string{"test", "here"}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"pass", args{"test:here", ":"}, result},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := split(tt.args.source, tt.args.separator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dequote(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"dequote", args{"This is a \"pain\""}, "This is a pain"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dequote(tt.args.target); got != tt.want {
				t.Errorf("dequote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boolean(t *testing.T) {
	type args struct {
		test *bool
	}
	test := true
	var undefined bool
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"true", args{&test}, true},
		{"false", args{nil}, false},
		{"undefined", args{&undefined}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := boolean(tt.args.test); got != tt.want {
				t.Errorf("boolean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decode64(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"decode", args{"QSB0ZXN0"}, "A test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decode64(tt.args.str); got != tt.want {
				t.Errorf("decode64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sprint(t *testing.T) {
	type args struct {
		unknown interface{}
	}

	var empty interface{}
	var sliceEmpty []interface{}
	var sliceNotEmpty []interface{}

	sliceNotEmpty = append(sliceNotEmpty, "SliceNotEmpty")
	notEmpty := "stuff"

	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{empty}, "\"\""},
		{"notEmpty", args{notEmpty}, "stuff"},
		{"sliceEmpty", args{sliceEmpty}, ""},
		{"sliceNotEmpty", args{sliceNotEmpty}, "SliceNotEmpty"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sprint(tt.args.unknown); got != tt.want {
				t.Errorf("sprint() = %v, want %v", got, tt.want)
			}
		})
	}
}
