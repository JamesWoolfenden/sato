package arm

import (
	"bytes"
	"fmt"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"
)

type castError struct {
	Type string
}

func (m *castError) Error() string {
	return fmt.Sprintf("failed to cast to %s", m.Type)
}

func parseParameters(
	result map[string]interface{},
	funcMap tftemplate.FuncMap, all string,
) (string, []interface{}, error) {
	parameters, ok := result["parameters"].(map[string]interface{})
	if !ok {
		return "", nil, &castError{Type: "map[string]interface{}"}
	}

	myVariables := make([]interface{}, 0)

	var err error

	for name, item := range parameters {
		myItem, ok := item.(map[string]interface{})

		if !ok {
			return "", nil, &castError{Type: "map[string]interface{}"}
		}

		myItem, err = FixType(myItem)
		if err != nil {
			log.Info().Err(err)
		}

		var output bytes.Buffer

		tmpl, err := tftemplate.New("test").Funcs(funcMap).Parse(string(variableFile))
		if err != nil {
			return "", nil, templateNewError{err: err}
		}

		err = tmpl.Execute(&output, m{
			"variable": myItem,
			"item":     name,
		})
		if err != nil {
			return "", nil, templateExecuteError{err: err}
		}

		all += output.String()

		myVariables = append(myVariables, myItem)
	}

	return all, myVariables, nil
}
