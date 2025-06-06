package arm

import (
	"bytes"
	"log"
	"sato/src/cf"
	tftemplate "text/template"
)

// ParseOutputs writes out to outputs.tf.
func ParseOutputs(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) error {
	if result["outputs"] == nil {
		return nil
	}

	outputs, ok := result["outputs"].(map[string]interface{})

	if !ok {
		return &castError{"map[string]interface{}"}
	}

	var All string

	for name, value := range outputs {
		var myVar cf.Output

		var someString string

		myVar.Type = "string"
		myVar.Name = name
		temp, ok := value.(map[string]interface{})

		if !ok {
			log.Printf("fail to assert value to map[string]interface{}")
		}

		someString, result = ParseString(temp["value"].(string), result)
		myVar.Value = someString

		var output bytes.Buffer

		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(outputFile))
		if err != nil {
			return templateNewError{err}
		}

		err = tmpl.Execute(&output, m{
			"variable": myVar,
			"item":     name,
		})

		if err != nil {
			return templateExecuteError{err: err}
		}

		All += output.String()
	}

	err := cf.Write(All, destination, "outputs")
	if err != nil {
		return writeFileError{destination: destination, err: err}
	}

	return nil
}
