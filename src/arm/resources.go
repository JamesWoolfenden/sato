package arm

import (
	"bytes"
	"context"
	"fmt"
	"sato/src/ai"
	"sato/src/see"
	"sato/src/tfgen"
	"strings"
	tftemplate "text/template"

	"github.com/rs/zerolog/log"
)

// ParseResources handles resources in ARM conversion.
func ParseResources(
	result map[string]interface{}, funcMap tftemplate.FuncMap, destination string, opts ...Option,
) (map[string]interface{}, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

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

		first, seeErr := see.Lookup(resourceType, false)

		temp, ok := myType["resource"].(string)
		if !ok {
			log.Warn().Msgf("resource name for %s is not a string, skipping", resourceType)

			continue
		}
		name = &temp

		if seeErr != nil || len(myContent) == 0 {
			if o.fallback != nil {
				if err := aiFallback(o.fallback, resourceType, first, *name, myType, destination); err != nil {
					log.Warn().Err(err).Msgf("ai fallback failed for %s", resourceType)
				}
			} else {
				log.Warn().Msgf("no terraform mapping for %s, skipping", resourceType)
			}

			continue
		}

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

func aiFallback(conv ai.Converter, armType string, tfHint *string, name string, resource map[string]interface{}, destination string) error {
	hint := ""
	if tfHint != nil {
		hint = *tfHint
	}

	res, err := conv.Convert(context.Background(), ai.Request{
		SourceType: armType,
		TFType:     hint,
		Provider:   "azurerm",
		Name:       name,
		Resource:   resource,
	})
	if err != nil {
		return err
	}

	log.Info().Msgf("AI-converted %s -> %s", armType, res.TFType)

	if err := tfgen.Write(ai.Header+res.HCL, destination, res.TFType+"."+name); err != nil {
		return err
	}

	return ai.WriteDraft(destination, res.TFType, res.Template)
}
