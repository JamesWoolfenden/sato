package arm

import (
	"path/filepath"
	"reflect"
	"testing"
)

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
		{"Pass", args{filepath.FromSlash("../../examples/arm/microsoft.compute/vm-simple-windows/azuredeploy.json"), ".sato"}, false},
		{"Fail", args{filepath.FromSlash("../../examples/arm/microsoft.compute/vm-simple-windows/nofile.json"), ".sato"}, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := Parse(tt.args.file, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_findResourceName(t *testing.T) {
	t.Parallel()

	type args struct {
		result map[string]interface{}
		name   string
	}

	target := map[string]interface{}{
		"resources": []interface{}{
			map[string]interface{}{
				"name":     "[variables('PrivateDNSZone')]",
				"resource": "sato0",
			},
		},
	}

	empty := map[string]interface{}{
		"resources": []interface{}{
			nil,
		},
	}

	veryEmpty := map[string]interface{}{
		"resources": nil,
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Pass", args{target, "PrivateDNSZone"}, "sato0", false},
		{"Fail", args{target, "MyDNSZone"}, "", true},
		{"format", args{target, "format('PrivateDNSZone')"}, "format('PrivateDNSZone')", true},
		{"Empty", args{empty, "PrivateDNSZone"}, "PrivateDNSZone", false},
		{"Very Empty", args{veryEmpty, "PrivateDNSZone"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findResourceName(tt.args.result, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("findResourceName() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("findResourceName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findResourceType(t *testing.T) {
	t.Parallel()

	type args struct {
		result map[string]interface{}
		name   string
	}

	target := map[string]interface{}{
		"resources": []interface{}{
			map[string]interface{}{
				"name":     "[variables('PrivateDNSZone')]",
				"resource": "sato0",
				"type":     "fred",
			},
		},
	}

	empty := map[string]interface{}{
		"resources": []interface{}{
			nil,
		},
	}

	veryEmpty := map[string]interface{}{
		"resources": nil,
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Pass", args{target, "fred"}, true},
		{"Empty", args{empty, ""}, false},
		{"Very Empty", args{veryEmpty, ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := findResourceType(tt.args.result, tt.args.name); got != tt.want {
				t.Errorf("findResourceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNameValue(t *testing.T) {
	t.Parallel()

	type args struct {
		result map[string]interface{}
		name   string
	}

	myEmpty := map[string]interface{}{}
	myResult := map[string]interface{}{
		"variables": map[string]interface{}{
			"finder": "found",
		},
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Found", args{myResult, "var.finder"}, "found", false},
		{"Double dots", args{myResult, "var.test.arse"}, "var.test.arse", true},
		{"Not Found", args{myResult, "var.missing"}, "var.missing", false},
		{"Empty Result Map", args{myEmpty, "var.missing"}, "var.missing", false},
		{"Nil", args{myResult, "guff"}, "guff", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNameValue(tt.args.result, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNameValue() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("getNameValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseList(t *testing.T) {
	type args struct {
		resources []interface{}
		result    map[string]interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		//{"pass", args{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseList(tt.args.resources, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseList() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseMap(t *testing.T) {
	t.Parallel()

	type args struct {
		myResource map[string]interface{}
		result     map[string]interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseMap(tt.args.myResource, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMap() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_preprocess(t *testing.T) {
	t.Parallel()

	type args struct {
		results map[string]interface{}
	}

	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		//{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := preprocess(tt.args.results); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("preprocess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replaceResourceID(t *testing.T) {
	t.Parallel()

	type args struct {
		Match  string
		result map[string]interface{}
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		//{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := replaceResourceID(tt.args.Match, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("replaceResourceID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("replaceResourceID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setResourceNames(t *testing.T) {
	t.Parallel()

	type args struct {
		results map[string]interface{}
	}

	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		//{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := setResourceNames(tt.args.results); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setResourceNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitResourceName(t *testing.T) {
	t.Parallel()

	type args struct {
		Attribute string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		//{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, err := splitResourceName(tt.args.Attribute)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitResourceName() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("splitResourceName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitResourceName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_parseLocals(t *testing.T) {
	t.Parallel()

	type args struct {
		result map[string]interface{}
	}

	parsed := "\tvmSubnetNsgId = azurerm_network_security_group.sato8 #[resourceId(\"Microsoft.Network/networkSecurityGroups\",variables(\"vmSubnetNsgName\"))]\n"
	bodge := make(map[string]interface{})

	startLocals := map[string]interface{}{
		"vmSubnetNsgId": "[resourceId('Microsoft.Network/networkSecurityGroups',variables('vmSubnetNsgName'))]",
	}

	locals := map[string]interface{}{
		"vmSubnetNsgId": "azurerm_network_security_group.sato8",
	}
	var resources []interface{}
	resource := map[string]interface{}{
		"type":     "Microsoft.Network/networkSecurityGroups",
		"name":     "[variables('vmSubnetNsgName')]",
		"location": "[parameters('location')]",
		"resource": "sato8",
	}
	resources = append(resources, resource)

	result := make(map[string]interface{})

	result["locals"] = locals
	result["resources"] = resources

	bodge["locals"] = startLocals
	bodge["resources"] = resources

	tests := []struct {
		name    string
		args    args
		want    string
		want1   map[string]interface{}
		wantErr bool
	}{
		{"test", args{bodge}, parsed, result, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, err := parseLocals(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLocals() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("parseLocals()\n got:\n %v\n want:\n %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseLocals() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_replace(t *testing.T) {
	t.Parallel()

	type args struct {
		matches      []string
		newAttribute string
		what         *string
		result       map[string]interface{}
	}

	tests := []struct {
		name  string
		args  args
		want  string
		want1 map[string]interface{}
	}{
		//{},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, got1 := replace(tt.args.matches, tt.args.newAttribute, tt.args.what, tt.args.result)
			if got != tt.want {
				t.Errorf("replace() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("replace() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_parseString(t *testing.T) {
	t.Parallel()

	type args struct {
		newAttribute string
		result       map[string]interface{}
	}

	tests := []struct {
		name  string
		args  args
		want  *string
		want1 map[string]interface{}
	}{
		//{},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1 := parseString(tt.args.newAttribute, tt.args.result)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseString() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
