package arm

import (
	"bytes"
	"sato/src/cf"
	tftemplate "text/template"
)

// ParseData writes out to data.tf.
func ParseData(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) error {
	if result["data"] == nil {
		return nil
	}

	data := result["data"]

	var output bytes.Buffer

	tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(dataFile))
	if err != nil {
		return &templateNewError{err: err}
	}

	err = tmpl.Execute(&output, m{
		"data": data,
	})

	if err != nil {
		return &templateExecuteError{err: err}
	}

	err = cf.Write(output.String(), destination, "data")
	if err != nil {
		return &writeFileError{destination: destination, err: err}
	}

	return nil
}
