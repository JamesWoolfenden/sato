package arm_test

import (
	"reflect"
	"sato/src/arm"
	"testing"
)

func TestIsLocal(t *testing.T) {
	t.Parallel()

	empty := make(map[string]interface{})
	result := make(map[string]interface{})
	locals := map[string]interface{}{
		"securityProfileJson": "example",
		"storageAccountName":  "dont care",
		"location":            "anywhere",
	}

	result["locals"] = locals
	empty["locals"] = "test"

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
		{name: "Empty", args: args{target: "var.sato", result: empty}, want: false},
		{name: "Found", args: args{target: "location", result: result}, want: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := arm.IsLocal(tt.args.target, tt.args.result); got != tt.want {
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
			if got := arm.ArrayToString(tt.args.defaultValue); got != tt.want {
				t.Errorf("ArrayToString() = %v, want %v", got, tt.want)
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

			got, got1 := arm.Contains(tt.args.s, tt.args.str)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Contains() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Contains() got1 = %v, want %v", got1, tt.want1)
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
			if got := arm.EscapeQuote(tt.args.item); got != tt.want {
				t.Errorf("EscapeQuote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fixType(t *testing.T) {
	t.Parallel()

	type args struct {
		myItem map[string]interface{}
	}

	myString := map[string]interface{}{
		"type":         "string",
		"defaultValue": "2",
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	myStringResult := map[string]interface{}{
		"type":         "string",
		"defaultValue": "2",
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	myInt := map[string]interface{}{
		"type":         "int",
		"defaultValue": 2,
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	myIntResult := map[string]interface{}{
		"type":         "number",
		"defaultValue": 2,
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	mySlice := map[string]interface{}{
		"type":         "[]interface{}",
		"defaultValue": []interface{}{1, 2},
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	mySliceResult := map[string]interface{}{
		"type":         "[]interface{}",
		"defaultValue": []interface{}{1, 2},
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
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
		"defaultValue": map[string]interface{}{"name": 1},
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	newObject := map[string]interface{}{
		"type": `object({
	name = number})`,
		"defaultValue": map[string]interface{}{"name": 1},
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": `{
	name = 1}`,
	}

	notObject := map[string]interface{}{
		"type":         43,
		"defaultValue": 1,
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	returnNotObject := map[string]interface{}{
		"type":         43,
		"defaultValue": 1,
		"maxValue":     25,
		"minValue":     0,
		"metadata": map[string]interface{}{
			"description": "Minimum number of replicas that will be deployed",
		},
		"default": "1",
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{"Nil by mouth", args{myNullItem}, myNullItemReturns, true},
		{"Do Nothing", args{myItem}, myItem, false},
		{"Object", args{myObject}, newObject, false},
		{"Not string", args{notObject}, returnNotObject, true},
		{"Slice", args{mySlice}, mySliceResult, false},
		{"Int", args{myInt}, myIntResult, false},
		{"String", args{myString}, myStringResult, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := arm.FixType(tt.args.myItem)
			if (err != nil) != tt.wantErr {
				t.Errorf("FixType() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FixType() got =\n %v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_ditch(t *testing.T) {
	t.Parallel()

	type args struct {
		Attribute string
		name      string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"none", args{"bastionPublicIpAddressName", "variables"}, "bastionPublicIpAddressName"},
		{"variables", args{"variables('bastionPublicIpAddressName')", "variables"}, "'bastionPublicIpAddressName'"},
		{"variables not", args{"parameters('vmName')", "variables"}, "parameters('vmName')"},
		{"mixed",
			args{
				"concat(parameters('vmName'),'/', variables('omsAgentForLinuxName'))",
				"variables"},
			"concat(parameters('vmName'),'/', 'omsAgentForLinuxName'))"},
		{"mixed 2",
			args{
				"concat(variables('blobPrivateDnsZoneName'), '/link_to_', toLower(parameters('virtualNetworkName')))",
				"variables"},
			"concat('blobPrivateDnsZoneName'), '/link_to_', toLower(parameters('virtualNetworkName'))"},
		{"concat", args{
			"uri(local._artifactsLocation, concat('artifacts/vm1.default.htm', var._artifactsLocationSasToken))",
			"concat"},
			"uri(local._artifactsLocation, 'artifacts/vm1.default.htm', var._artifactsLocationSasToken)"},
		{"works",
			args{
				`[uri(parameters("_artifactsLocation"), concat("artifacts/vm2.default.htm", parameters("_artifactsLocationSasToken")))]`,
				"uri"},
			"[parameters(\"_artifactsLocation\"), concat(\"artifacts/vm2.default.htm\", parameters(\"_artifactsLocationSasToken\"))]"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := arm.Ditch(tt.args.Attribute, tt.args.name); got != tt.want {
				t.Errorf("Ditch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_enabled(t *testing.T) {
	t.Parallel()

	type args struct {
		status string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Enabled", args{"Enabled"}, true},
		{"Disabled", args{"Disabled"}, false},
		{"Guff", args{"guff"}, false},
		{"Nil", args{""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := arm.Enabled(tt.args.status); got != tt.want {
				t.Errorf("Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notNil(t *testing.T) {
	t.Parallel()

	type args struct {
		unknown interface{}
	}

	var (
		test  interface{}
		test2 interface{}
		test3 interface{}
	)

	test = true
	test3 = false
	test4 := ""

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"True", args{test}, true},
		{"Nil", args{test2}, false},
		{"False", args{test3}, true},
		{"String", args{""}, true},
		{"Pointer", args{&test4}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := arm.NotNil(tt.args.unknown); got != tt.want {
				t.Errorf("NotNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tags(t *testing.T) {
	t.Parallel()

	myTags := map[string]interface{}{
		"test":    "something",
		"another": "stuff",
	}

	myBodge := map[string]interface{}{
		"test":    "something",
		"another": false,
	}

	type args struct {
		tags map[string]interface{}
	}

	result := "{\n\t\"test\" = \"something\"\n\t\"another\" = \"stuff\"\n\t}"
	bodge := "{\n\t\"test\" = \"something\"\n\t\"another\" = \"OBJECT\"\n\t}"

	tests := []struct {
		name string
		args args
		want string
	}{
		{"pass", args{myTags}, result},
		{"bodge", args{myBodge}, bodge},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := arm.Tags(tt.args.tags); got != tt.want {
				t.Errorf("Tags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uuid(t *testing.T) {
	t.Parallel()

	type args struct {
		count int
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"Pass", args{1}, "resource \"random_uuid\" \"sato0\" {}\n"},
		{"Two", args{2}, "resource \"random_uuid\" \"sato0\" {}\nresource \"random_uuid\" \"sato1\" {}\n"},
		{"Zero", args{0}, ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := arm.UUID(tt.args.count); got != tt.want {
				t.Errorf("uuid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loseSQBrackets(t *testing.T) {
	t.Parallel()

	type args struct {
		newAttribute string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"pass", args{"[pass]"}, "pass"},
		{"no pass", args{"[pass"}, "[pass"},
		{"leave", args{"stuff[pass]"}, "stuff[pass]"},
		{"leave with outside", args{"[stuff[pass]]"}, "stuff[pass]"},
		{"just", args{"[]"}, ""},
		{"nil", args{""}, ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := arm.LoseSQBrackets(tt.args.newAttribute); got != tt.want {
				t.Errorf("LoseSQBrackets() = %v, want %v", got, tt.want)
			}
		})
	}
}
