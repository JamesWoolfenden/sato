package arm

import (
	"bytes"
	"sato/src/cf"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"
)

// parseData writes out to data.tf
func parseData(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) error {
	if result["data"] == nil {
		return nil
	}

	data := result["data"]

	var output bytes.Buffer

	tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(dataFile))

	if err != nil {
		return err
	}

	err = tmpl.Execute(&output, m{
		"data": data,
	})

	if err != nil {
		log.Print(err)

		return err
	}

	err = cf.Write(output.String(), destination, "data")
	if err != nil {
		log.Print(err)

		return err
	}

	return nil
}
