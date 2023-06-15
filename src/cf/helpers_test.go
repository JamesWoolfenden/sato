package cf_test

import (
	"reflect"
	"testing"

	sato "sato/src/cf"
)

func Test_replace(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Replace(tt.args.input, tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_add(t *testing.T) {
	t.Parallel()

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
		{"Add-empty", args{"test", added, m}, after, afterMap},
		{"no-change", args{"test", exists, afterMap}, after, afterMap},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1 := sato.Add(tt.args.s, tt.args.a, tt.args.m)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Add() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_split(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Split(tt.args.source, tt.args.separator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dequote(t *testing.T) {
	t.Parallel()

	type args struct {
		target string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"Dequote", args{"This is a \"pain\""}, "This is a pain"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Dequote(tt.args.target); got != tt.want {
				t.Errorf("Dequote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_boolean(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Boolean(tt.args.test); got != tt.want {
				t.Errorf("Boolean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decode64(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Decode64(tt.args.str); got != tt.want {
				t.Errorf("Decode64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sprint(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Sprint(tt.args.unknown); got != tt.want {
				t.Errorf("Sprint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_quote(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Quote(tt.args.target); got != tt.want {
				t.Errorf("Quote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_snake(t *testing.T) {
	t.Parallel()

	type args struct {
		Camel string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"Camel", args{"ThisIsSnakeCasing"}, "this_is_snake_casing"},
		{"Snake", args{"This-Is-Snake-Casing"}, "this_is_snake_casing"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Snake(tt.args.Camel); got != tt.want {
				t.Errorf("Snake() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_kebab(t *testing.T) {
	t.Parallel()

	type args struct {
		Camel string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Camel", args{"ThisIsCamelCasing"}, "this-is-camel-casing"},
		{"Snake", args{"This_Is_Camel_Casing"}, "this-is-camel-casing"},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Kebab(tt.args.Camel); got != tt.want {
				t.Errorf("Kebab() = %v, want %v", got, tt.want)
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
			if got := sato.Lower(tt.args.target); got != tt.want {
				t.Errorf("Lower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nill(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Nill(tt.args.str); got != tt.want {
				t.Errorf("Nill() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nild(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Nild(tt.args.str, tt.args.myDefault); got != tt.want {
				t.Errorf("Nild() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_array(t *testing.T) {
	t.Parallel()

	type args struct {
		mySlice []string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Replace",
			args: args{mySlice: []string{"some", "thing", "punt"}},
			want: "[\n\t\"some\",\t\"thing\",\t\"punt\"\n\t]\n",
		},
		{
			name: "nil",
			args: args{mySlice: nil},
			want: "[]",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Array(tt.args.mySlice); got != tt.want {
				t.Errorf("Array() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayReplace(t *testing.T) {
	t.Parallel()

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
		{
			name: "Replace",
			args: args{mySlice: []string{"some", "thing", "punt"}, target: "pu", replacement: "ca"},
			want: "[\n\t\"some\",\t\"thing\",\t\"cant\"\n\t]\n",
		},
		{
			name: "nil",
			args: args{mySlice: nil, target: "pu", replacement: "ca"},
			want: "[]",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.ArrayReplace(tt.args.mySlice, tt.args.target, tt.args.replacement); got != tt.want {
				t.Errorf("ArrayReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Contains(tt.args.target, tt.args.substring); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_zipfile(t *testing.T) {
	t.Parallel()

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
		//{},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := sato.Zipfile(tt.args.code, tt.args.filename, tt.args.runtime); got != tt.want {
				t.Errorf("Zipfile() = %v, want %v", got, tt.want)
			}
		})
	}
}
