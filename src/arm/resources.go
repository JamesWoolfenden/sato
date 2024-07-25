package arm

import (
	"bytes"
	"fmt"
	"sato/src/cf"
	"sato/src/see"
	"strings"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"
)

// ParseResources handles resources in ARM conversion
func ParseResources(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) (map[string]interface{}, error) {
	resources := result["resources"].([]interface{})

	newResources, err := ParseList(resources, result)
	if err != nil {
		return nil, err
	}

	result["resources"] = newResources

	for _, resource := range newResources {
		var output bytes.Buffer

		var name *string

		myType := resource.(map[string]interface{})

		myContent := lookup(myType["type"].(string))

		first, err := see.Lookup(myType["type"].(string), false)
		if err != nil {
			log.Warn().Err(err)

			continue
		}

		temp, ok := myType["resource"].(string)

		if !ok {
			log.Warn().Msg("myType[\"resource\"] is not string")
		}

		name = &temp

		// needs to pivot on policy template from resource
		tmpl, err := tftemplate.New("sato").Funcs(funcMap).Parse(string(myContent))
		if err != nil {
			log.Warn().Msgf("failed at %s  for %s %s", err, *first, *name)

			continue
		}

		_ = tmpl.Execute(&output, cf.M{
			"resource": resource,
			"item":     name,
		})

		err = cf.Write(output.String(), destination, *first+"."+strings.Replace(*name, "var.", "", 1))
		if err != nil {
			return nil, fmt.Errorf("write failure %w", err)
		}
	}

	return result, nil
}
