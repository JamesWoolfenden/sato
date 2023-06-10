package arm

import (
	"bytes"
	"sato/src/cf"
	"sato/src/see"
	"strings"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"
)

// ParseResources handles resources in ARM conversion
func ParseResources(result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) (map[string]interface{}, error) {
	resources := result["resources"].([]interface{})

	newResources, err := parseList(resources, result)

	if err != nil {
		return nil, err
	}

	result["resources"] = newResources

	for _, resource := range newResources {
		var output bytes.Buffer
		var name *string
		myType := resource.(map[string]interface{})
		myContent := lookup(myType["type"].(string))
		first, err := see.Lookup(myType["type"].(string))

		if err != nil {
			log.Warn().Err(err)
			continue
		}

		temp := myType["resource"].(string)
		name = &temp

		// needs to pivot on policy template from resource
		tmpl, err := tftemplate.New("sato").Funcs(funcMap).Parse(string(myContent))
		if err != nil {
			log.Printf("failed at %s  for %s %s", err, *first, *name)
			continue
		}

		_ = tmpl.Execute(&output, cf.M{
			"resource": resource,
			"item":     name,
		})

		err = cf.Write(output.String(), destination, *first+"."+strings.Replace(*name, "var.", "", 1))
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
