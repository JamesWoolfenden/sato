package arm_test

import (
	"path/filepath"
	"reflect"
	"sato/src/arm"
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
		{
			name: "Pass",
			args: args{
				file:        filepath.FromSlash("../../examples/arm/microsoft.compute/vm-simple-windows/azuredeploy.json"),
				destination: ".sato",
			},
		},
		{
			name: "Fail",
			args: args{
				file:        filepath.FromSlash("../../examples/arm/microsoft.compute/vm-simple-windows/nofile.json"),
				destination: ".sato",
			},
			wantErr: true,
		},
		{
			name: "Empty",
			args: args{
				file:        filepath.FromSlash("../../examples/arm/smallest.json"),
				destination: ".sato",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := arm.Parse(tt.args.file, tt.args.destination); (err != nil) != tt.wantErr {
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
		{"Not found", args{target, "MyDNSZone"}, "MyDNSZone", false},
		{"format", args{target, "format('PrivateDNSZone')"}, "format('PrivateDNSZone')", true},
		{"Empty", args{empty, "PrivateDNSZone"}, "PrivateDNSZone", false},
		{"Very Empty", args{veryEmpty, "PrivateDNSZone"}, "", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := arm.FindResourceName(tt.args.result, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindResourceName() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("FindResourceName() got = %v, want %v", got, tt.want)
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
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := arm.FindResourceType(tt.args.result, tt.args.name); got != tt.want {
				t.Errorf("FindResourceType() = %v, want %v", got, tt.want)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := arm.GetNameValue(tt.args.result, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNameValue() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("GetNameValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseList(t *testing.T) {
	t.Parallel()

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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := arm.ParseList(tt.args.resources, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseList() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseList() got = %v, want %v", got, tt.want)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := arm.ParseMap(tt.args.myResource, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMap() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMap() got = %v, want %v", got, tt.want)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := arm.Preprocess(tt.args.results); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Preprocess() = %v, want %v", got, tt.want)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := arm.ReplaceResourceID(tt.args.Match, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceResourceID() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("ReplaceResourceID() got = %v, want %v", got, tt.want)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := arm.SetResourceNames(tt.args.results); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetResourceNames() = %v, want %v", got, tt.want)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, err := arm.SplitResourceName(tt.args.Attribute)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitResourceName() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("SplitResourceName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitResourceName() got1 = %v, want %v", got1, tt.want1)
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1, err := arm.ParseLocals(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLocals() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("ParseLocals()\n got:\n %v\n want:\n %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseLocals() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_replace(t *testing.T) {
	t.Parallel()
	//matches := []string{
	//	"parameters", "variables", "toLower", "resourceGroup().location", "resourceGroup().id",
	//	"substring", "uniqueString", "reference", "resourceId", "listKeys", "format('", "SubscriptionResourceId",
	//	"concat", "subscription().tenantId", "uuid",
	//}
	//
	//target := map[string]interface{}{
	//	"locals": map[string]interface{}{
	//		"location": "[resourceGroup().location]",
	//	},
	//	"resources": []interface{}{
	//		map[string]interface{}{
	//			"name":     "[variables('appGwPublicIpName')]",
	//			"resource": "sato0",
	//			"type":     "Microsoft.Network/publicIPAddresses",
	//		},
	//	},
	//}
	//result := target

	type args struct {
		matches      []string
		newAttribute string
		what         *string
		result       map[string]interface{}
	}

	// vars := "variables"
	// resourceId := "resourceId"
	// params := "parameters"

	tests := []struct {
		name  string
		args  args
		want  string
		want1 map[string]interface{}
	}{
		//{"Pass",
		//	args{matches, "[resourceId('Microsoft.Network/publicIPAddresses',variables('appGwPublicIpName'))]", &vars, target},
		//	"azurerm_public_ip.sato0", result},
		//{"resourceId",
		//	args{matches, "[resourceId('Microsoft.Network/publicIPAddresses',local.appGwPublicIpName)]", &resourceId, target},
		//	"azurerm_public_ip.sato0", result},
		//{"parameters", args{matches, "[parameters('wafMode')]", &params, target}, "local.wafMode", result},
		//{"bodge",
		//	args{
		//		matches, "[resourceId('Microsoft.Network/virtualNetworks/subnets', variables('virtualNetworkName'), variables('appGatewaySubnetName'))]",
		//		&vars,
		//		target},
		//	"tolist(azurerm_virtual_network.sato6.subnet)[0].id",
		//	result},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, got1 := arm.Replace(tt.args.matches, tt.args.newAttribute, tt.args.what, tt.args.result)
			if got != tt.want {
				t.Errorf("Replace() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Replace() got1 = %v, want %v", got1, tt.want1)
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

	//want := "${var._artifactsLocation}artifacts/vm2.default.htm"
	//target := map[string]interface{}{}
	//
	//result := map[string]interface{}{}

	tests := []struct {
		name  string
		args  args
		want  *string
		want1 map[string]interface{}
	}{
		//{"Pass", args{"vm2DefaultHtmFullPath", target}, &want, result},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, got1 := arm.ParseString(tt.args.newAttribute, tt.args.result)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseString() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ParseString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
