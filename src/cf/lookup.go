package cf

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	tftemplate "text/template"

	"sato/src/ai"
	"sato/src/see"
	"sato/src/tfgen"

	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/rs/zerolog/log"
)

type templateNewError struct {
	Err error
}

func (e *templateNewError) Error() string {
	return fmt.Sprintf("failed to create template %v", e.Err)
}

// parseResources converts resource to Terraform.
func parseResources(resources cloudformation.Resources, funcMap tftemplate.FuncMap, destination string, o options) error {
	if resources == nil || funcMap == nil || destination == "" {
		return &missingResourceInputError{}
	}

	for item, resource := range resources {
		var output bytes.Buffer

		myType := resources[item].AWSCloudFormationType()

		myContent := lookup(myType)

		if len(myContent) == 0 {
			if o.fallback != nil {
				if err := aiFallback(o.fallback, myType, item, resource, destination); err != nil {
					log.Warn().Err(err).Msgf("ai fallback failed for %s", myType)
				}
			}
			continue
		}

		// needs to pivot on policy template from resource
		tmpl, err := tftemplate.New("sato").Funcs(funcMap).Parse(string(myContent))
		if err != nil {
			return &templateNewError{Err: err}
		}

		_ = tmpl.Execute(&output, M{
			"resource":    resource,
			"item":        item,
			"destination": destination,
		})

		result := ReplaceDependant(ReplaceVariables(output.String()))

		// Special processing for Step Functions state machines to convert resource references
		if myType == "AWS::StepFunctions::StateMachine" {
			result = ReplaceStepFunctionsReferences(result, resources)
		}

		err = tfgen.Write(result, destination, fmt.Sprint(ToTFName(myType), ".", strings.ToLower(item)))
		if err != nil {
			return &writeError{destination: destination, err: err}
		}
	}

	return nil
}

func aiFallback(conv ai.Converter, cfnType, item string, resource any, destination string) error {
	tfType := ""
	if t, err := see.Lookup(cfnType, false); err == nil && t != nil {
		tfType = *t
	}

	res, err := conv.Convert(context.Background(), ai.Request{
		SourceType: cfnType,
		TFType:     tfType,
		Provider:   "aws",
		Name:       strings.ToLower(item),
		Resource:   resource,
	})
	if err != nil {
		return err
	}

	log.Info().Msgf("AI-converted %s -> %s", cfnType, res.TFType)

	hcl := ai.Header + ReplaceDependant(ReplaceVariables(res.HCL))
	if err := tfgen.Write(hcl, destination, res.TFType+"."+strings.ToLower(item)); err != nil {
		return err
	}

	return ai.WriteDraft(destination, res.TFType, res.Template)
}

//goland:noinspection GoLinter
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
