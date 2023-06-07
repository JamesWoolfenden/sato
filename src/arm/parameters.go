package arm

import (
	"bytes"
	"fmt"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"
)

func parseParameters(result map[string]interface{}, funcMap tftemplate.FuncMap, All string) (string, []interface{}, error) {
	parameters, ok := result["parameters"].(map[string]interface{})
	if !ok {
		return "", nil, fmt.Errorf("failed to cast to map")
	}
	myVariables := make([]interface{}, 0)
	var err error
	for name, item := range parameters {

		myItem := item.(map[string]interface{})

		myItem, err = fixType(myItem)
		if err != nil {
			log.Print(err)
		}

		var output bytes.Buffer
		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(variableFile))
		if err != nil {
			return "", nil, err
		}
		_ = tmpl.Execute(&output, m{
			"variable": myItem,
			"item":     name,
		})
		All = All + output.String()
		myVariables = append(myVariables, myItem)
	}
	return All, myVariables, nil
}
