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

func Test_quote(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test", args{"something"}, "\"something\""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := quote(tt.args.target); got != tt.want {
				t.Errorf("quote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_snake(t *testing.T) {
	type args struct {
		Camel string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Camel", args{"ThisIsSnakeCasing"}, "this_is_snake_casing"},
		{"snake", args{"This-Is-Snake-Casing"}, "this_is_snake_casing"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := snake(tt.args.Camel); got != tt.want {
				t.Errorf("snake() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kebab(t *testing.T) {
	type args struct {
		Camel string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Camel", args{"ThisIsCamelCasing"}, "this-is-camel-casing"},
		{"snake", args{"This_Is_Camel_Casing"}, "this-is-camel-casing"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := kebab(tt.args.Camel); got != tt.want {
				t.Errorf("kebab() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lower(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"basic", args{"SHOUTING"}, "shouting"},
		{"with special chars", args{"SHOUTING!"}, "shouting!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lower(tt.args.target); got != tt.want {
				t.Errorf("lower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nill(t *testing.T) {
	type args struct {
		str *string
	}
	exists := "here"
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Nil", args{nil}, ""},
		{"Not Nil", args{&exists}, exists},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nill(tt.args.str); got != tt.want {
				t.Errorf("nill() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nild(t *testing.T) {
	type args struct {
		str       *string
		myDefault string
	}
	exists := "here"
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Nil", args{nil, "hello"}, "hello"},
		{name: "NotNil", args: args{str: &exists, myDefault: "hello"}, want: "here"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nild(tt.args.str, tt.args.myDefault); got != tt.want {
				t.Errorf("nild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_array(t *testing.T) {
	type args struct {
		mySlice []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "replace",
			args: args{mySlice: []string{"some", "thing", "punt"}},
			want: "[\n\t\"some\",\t\"thing\",\t\"punt\"\n\t]\n"},
		{name: "nil",
			args: args{mySlice: nil},
			want: "[]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := array(tt.args.mySlice); got != tt.want {
				t.Errorf("array() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayReplace(t *testing.T) {
	type args struct {
		mySlice     []string
		target      string
		replacement string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "replace",
			args: args{mySlice: []string{"some", "thing", "punt"}, target: "pu", replacement: "ca"},
			want: "[\n\t\"some\",\t\"thing\",\t\"cant\"\n\t]\n"},
		{name: "nil",
			args: args{mySlice: nil, target: "pu", replacement: "ca"},
			want: "[]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arrayReplace(tt.args.mySlice, tt.args.target, tt.args.replacement); got != tt.want {
				t.Errorf("arrayReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		target    string
		substring string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"found", args{"scunthorpe", "thor"}, true},
		{"not found", args{"typhoo", "thor"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.target, tt.args.substring); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_zipfile(t *testing.T) {
	type args struct {
		code     string
		filename string
		runtime  string
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
			if got := zipfile(tt.args.code, tt.args.filename, tt.args.runtime); got != tt.want {
				t.Errorf("zipfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
