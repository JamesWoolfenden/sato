package cf

import (
	"bytes"
	"fmt"
	"strings"
	tftemplate "text/template"

	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/rs/zerolog/log"
)

// parseResources converts resource to Terraform.
func parseResources(resources cloudformation.Resources, funcMap tftemplate.FuncMap, destination string) error {
	if resources == nil || funcMap == nil || destination == "" {
		return &missingResourceInputError{}
	}

	for item, resource := range resources {
		var output bytes.Buffer

		myType := resources[item].AWSCloudFormationType()

		myContent := lookup(myType)

		// needs to pivot on policy template from resource
		tmpl, err := tftemplate.New("sato").Funcs(funcMap).Parse(string(myContent))
		if err != nil {
			return fmt.Errorf("failed to template %w", err)
		}

		_ = tmpl.Execute(&output, M{
			"resource": resource,
			"item":     item,
		})

		err = Write(
			ReplaceDependant(
				ReplaceVariables(output.String())), destination, fmt.Sprint(ToTFName(myType), ".", strings.ToLower(item)))
		if err != nil {
			return err
		}
	}

	return nil
}

//goland:noinspection GoLinter
//nolint:funlen
func lookup(myType string) []byte {

	var myContent []byte

	var ok bool

	if tfLookup[myType] != nil {
		myContent, ok = tfLookup[myType].([]byte)
		if !ok {
			log.Warn().Msg("failed to cast to []byte")
		}
	} else {
		// we don't want to half the parsing so log it.
		log.Warn().Msgf("%s not found", myType)
	}

	return myContent
}
