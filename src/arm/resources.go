package arm

import (
	"bytes"
	"fmt"
	"sato/src/see"
	"sato/src/tfgen"
	"strings"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"
)

// ParseResources handles resources in ARM conversion.
func ParseResources(
	result map[string]interface{}, funcMap tftemplate.FuncMap, destination string) (map[string]interface{}, error) {
	resources, ok := result["resources"].([]interface{})

	if !ok {
		return result, &castError{"[]interface{}"}
	}

	newResources, err := ParseList(resources, result)
	if err != nil {
		return nil, err
	}

	result["resources"] = newResources

	for _, resource := range newResources {
		var output bytes.Buffer

		var name *string

		myType, ok := resource.(map[string]interface{})

		if !ok {
			log.Warn().Msg("resource is not map[string]interface{}, skipping")

			continue
		}

		resourceType, ok := myType["type"].(string)

		if !ok {
			log.Warn().Msg("resource type is not a string, skipping")

			continue
		}

		myContent := lookup(resourceType)

		first, err := see.Lookup(resourceType, false)
		if err != nil {
			log.Warn().Err(err).Msgf("no terraform mapping for %s, skipping", resourceType)

			continue
		}

		temp, ok := myType["resource"].(string)

		if !ok {
			log.Warn().Msgf("resource name for %s is not a string, skipping", resourceType)

			continue
		}

		name = &temp

		// needs to pivot on policy template from resource
		tmpl, err := tftemplate.New("sato").Funcs(funcMap).Parse(string(myContent))
		if err != nil {
			log.Warn().Msgf("failed at %s  for %s %s", err, *first, *name)

			continue
		}

		_ = tmpl.Execute(&output, tfgen.M{
			"resource": resource,
			"item":     name,
		})

		err = tfgen.Write(output.String(), destination, *first+"."+strings.Replace(*name, "var.", "", 1))
		if err != nil {
			return nil, fmt.Errorf("write failure %w", err)
		}
	}

	return result, nil
}
