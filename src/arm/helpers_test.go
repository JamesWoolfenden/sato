package arm

import (
	"reflect"
	"testing"
)

func TestIsLocal(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := IsLocal(tt.args.target, tt.args.result); got != tt.want {
				t.Errorf("IsLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arrayToString(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := arrayToString(tt.args.defaultValue); got != tt.want {
				t.Errorf("arrayToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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
	t.Parallel()

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
		{"nil", args{nil}, ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := escapeQuote(tt.args.item); got != tt.want {
				t.Errorf("escapeQuote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fixType(t *testing.T) {
	t.Parallel()

	type args struct {
		myItem map[string]interface{}
	}

	myItem := map[string]interface{}{
		"type":         "number",
		"defaultValue": 1,
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	myNullItem := map[string]interface{}{
		"type":         nil,
		"defaultValue": 1,
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	myNullItemReturns := map[string]interface{}{
		"type":         nil,
		"defaultValue": 1,
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	myObject := map[string]interface{}{
		"type":         "object",
		"defaultValue": 1,
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	//myObject2 := map[string]interface{}{
	//	"type": "object",
	//	"defaultValue": map[string]interface{}{
	//		"firewallRules": []interface{}{
	//			map[string]interface{}{
	//				"firewallRuleName": "AllowFromAll",
	//				"rangeStart":       "0.0.0.0",
	//				"rangeEnd":         "255.255.255.255",
	//			},
	//		},
	//	},
	//	"maxValue": 25,
	//	"minValue": 0,
	//	"metadata": map[string]interface{}{
	//		"description": "Minimum number of replicas that will be deployed",
	//	},
	//	"default": "1",
	//}
	//
	//myObjectReturned := map[string]interface{}{
	//	"type": "object({\n\tfirewallRules= list(object({\n\t   firewallRuleName = string\n\t   rangeStart = string\n\t   rangeEnd = string}))})",
	//	"defaultValue": map[string]interface{}{
	//		"firewallRules": []interface{}{
	//			map[string]interface{}{
	//				"firewallRuleName": "AllowFromAll",
	//				"rangeStart":       "0.0.0.0",
	//				"rangeEnd":         "255.255.255.255",
	//			},
	//		},
	//	},
	//	"maxValue": 25,
	//	"minValue": 0,
	//	"metadata": map[string]interface{}{
	//		"description": "Minimum number of replicas that will be deployed",
	//	},
	//	"default": "{\n\tfirewallRules= [{\n\t   firewallRuleName = \"AllowFromAll\"\n\t   rangeStart = \"0.0.0.0\"\n\t   rangeEnd = \"255.255.255.255\"}]}",
	//}

	newObject := myObject
	newObject["type"] = "string"
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{"Nil by mouth", args{myNullItem}, myNullItemReturns, true},
		{"Do Nothing", args{myItem}, myItem, false},
		{"Object", args{myObject}, newObject, false},
		//{"Convert Object", args{myObject2}, myObjectReturned, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := fixType(tt.args.myItem)
			if (err != nil) != tt.wantErr {
				t.Errorf("fixType() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fixType() got = %v, want %v", got, tt.want)
			}
		})
	}
}
