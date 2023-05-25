package arm

import (
	"reflect"
	"testing"
)

func TestIsLocal(t *testing.T) {
	result := make(map[string]interface{})
	locals := map[string]interface{}{
		"securityProfileJson": "example",
		"storageAccountName":  "dont care",
		"location":            "anywhere",
	}
	result["locals"] = locals

	type args struct {
		target string
		result map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Pass", args: args{target: "var.sato", result: result}, want: false},
		{name: "Found", args: args{target: "location", result: result}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLocal(tt.args.target, tt.args.result); got != tt.want {
				t.Errorf("IsLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayToString(t *testing.T) {
	type args struct {
		defaultValue []interface{}
	}

	data := []interface{}{"I", "AM", "HERE"}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"Pass", args{defaultValue: data}, "[\"I\",\"AM\",\"HERE\"]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arrayToString(tt.args.defaultValue); got != tt.want {
				t.Errorf("arrayToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		s   []string
		str string
	}

	result := "dave"
	tests := []struct {
		name  string
		args  args
		want  *string
		want1 bool
	}{
		{"Not Found", args{[]string{"dave", "annie"}, "fred"}, nil, false},
		{"Found", args{[]string{"dave", "annie"}, "dave"}, &result, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := contains(tt.args.s, tt.args.str)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("contains() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("contains() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_escapeQuote(t *testing.T) {
	type args struct {
		item interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"pass", args{"\""}, "\\\""},
		{"multiple", args{"\"this\""}, "\\\"this\\\""},
		{"do nothing", args{"/"}, "/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeQuote(tt.args.item); got != tt.want {
				t.Errorf("escapeQuote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fixType(t *testing.T) {
	type args struct {
		myItem map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fixType(tt.args.myItem); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fixType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_translate(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := translate(tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("translate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("translate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
