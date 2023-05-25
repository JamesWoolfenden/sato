package arm

import (
	"reflect"
	sato "sato/src"
	"testing"
	"text/template"
)

func TestGetValue(t *testing.T) {
	type args struct {
		item      string
		variables []sato.Variable
	}
	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetValue(tt.args.item, tt.args.variables)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestParseData(t *testing.T) {
	type args struct {
		result      map[string]interface{}
		funcMap     template.FuncMap
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
			if err := ParseData(tt.args.result, tt.args.funcMap, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("ParseData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseOutputs(t *testing.T) {
	type args struct {
		result      map[string]interface{}
		funcMap     template.FuncMap
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
			if err := ParseOutputs(tt.args.result, tt.args.funcMap, tt.args.destination); (err != nil) != tt.wantErr {
				t.Errorf("ParseOutputs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseResources(t *testing.T) {
	type args struct {
		result      map[string]interface{}
		funcMap     template.FuncMap
		destination string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseResources(tt.args.result, tt.args.funcMap, tt.args.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseResources() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseVariables(t *testing.T) {
	type args struct {
		result      map[string]interface{}
		funcMap     template.FuncMap
		destination string
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVariables(tt.args.result, tt.args.funcMap, tt.args.destination)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVariables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseVariables() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findResourceName(t *testing.T) {
	type args struct {
		result map[string]interface{}
		name   string
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
	type args struct {
		result map[string]interface{}
		name   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findResourceType(tt.args.result, tt.args.name); got != tt.want {
				t.Errorf("findResourceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNameValue(t *testing.T) {
	type args struct {
		result map[string]interface{}
		name   string
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

func Test_handleResource(t *testing.T) {
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
			got, err := handleResource(tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handleResource() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isCompound(t *testing.T) {
	type args struct {
		newAttribute string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCompound(tt.args.newAttribute); got != tt.want {
				t.Errorf("isCompound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loseSQBrackets(t *testing.T) {
	type args struct {
		newAttribute string
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
			if got := loseSQBrackets(tt.args.newAttribute); got != tt.want {
				t.Errorf("loseSQBrackets() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

func Test_parseParameters(t *testing.T) {
	type args struct {
		result  map[string]interface{}
		funcMap template.FuncMap
		All     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseParameters(tt.args.result, tt.args.funcMap, tt.args.All)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseParameters() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseParameters() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_preprocess(t *testing.T) {
	type args struct {
		results map[string]interface{}
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
			if got := preprocess(tt.args.results); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("preprocess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replaceResourceID(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	type args struct {
		results map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setResourceNames(tt.args.results); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setResourceNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitResourceName(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
	type args struct {
		result map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseLocals(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLocals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseLocals() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseLocals() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_replace(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
