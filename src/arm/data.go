package arm

import (
	"bytes"
	"fmt"
	tftemplate "text/template"

	"sato/src/cf"
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
		return fmt.Errorf("failed to build parser %w", err)
	}

	err = tmpl.Execute(&output, m{
		"data": data,
	})

	if err != nil {
		return fmt.Errorf("failed to execute parser %w", err)
	}

	err = cf.Write(output.String(), destination, "data")
	if err != nil {
		return fmt.Errorf("failed to write data.tf %w", err)
	}

	return nil
}
