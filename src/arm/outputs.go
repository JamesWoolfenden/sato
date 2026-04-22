package arm

import (
	"bytes"
	"sato/src/tfgen"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"
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
		var myVar tfgen.Output

		var someString string

		myVar.Type = "string"
		myVar.Name = name
		temp, ok := value.(map[string]interface{})

		if !ok {
			log.Warn().Msgf("output %q is not an object", name)
			continue
		}

		raw, ok := temp["value"].(string)
		if !ok {
			log.Warn().Msgf("output %q has non-string value", name)
			continue
		}

		someString, result = ParseString(raw, result)
		myVar.Value = someString

		var output bytes.Buffer

		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(outputFile))
		if err != nil {
			return &templateNewError{Err: err}
		}

		err = tmpl.Execute(&output, m{
			"variable": myVar,
			"item":     name,
		})

		if err != nil {
			return &templateExecuteError{Err: err}
		}

		All += output.String()
	}

	err := tfgen.Write(All, destination, "outputs")
	if err != nil {
		return &writeFileError{Destination: destination, Err: err}
	}

	return nil
}
