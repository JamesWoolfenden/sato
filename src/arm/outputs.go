package arm

import (
	"bytes"
	"sato/src/cf"
	tftemplate "text/template"
)

// parseOutputs writes out to outputs.tf
func parseOutputs(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) error {
	if result["outputs"] == nil {
		return nil
	}

	outputs := result["outputs"].(map[string]interface{})

	var All string
	for name, value := range outputs {
		var myVar cf.Output
		var someString string
		myVar.Type = "string"
		myVar.Name = name
		temp := value.(map[string]interface{})
		someString, result = parseString(temp["value"].(string), result)
		myVar.Value = someString
		var output bytes.Buffer
		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(outputFile))
		if err != nil {
			return err
		}
		_ = tmpl.Execute(&output, m{
			"variable": myVar,
			"item":     name,
		})
		All = All + output.String()
	}

	err := cf.Write(All, destination, "outputs")
	if err != nil {
		return err
	}

	return nil
}
