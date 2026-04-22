package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

const defaultModel = "claude-sonnet-4-6"

const systemPrompt = `You convert a single CloudFormation or ARM resource into Terraform.

Reply with exactly three sections, no prose:

===TFTYPE===
<terraform resource type, e.g. aws_sqs_queue>
===HCL===
<one resource block for the given instance; valid HCL only>
===TEMPLATE===
<a Go text/template that generalises the HCL over .resource and .item>

Template rules:
- The block label is "{{.item}}".
- Source fields are under .resource using the upstream casing (e.g. .resource.BucketName for CFN, .resource.properties.x for ARM).
- Wrap optional fields in {{- if .resource.X}} ... {{- end}}.
- Pipe string attributes through |Quote so var./local. references stay unquoted.
- Helpers available: Quote, Nild, Boolean, Marshal, Tags, Snake, Replace.

Example template:
resource "aws_sns_topic" "{{.item}}" {
  name = {{Nild .resource.TopicName (.item)|Quote}}
{{- if .resource.KmsMasterKeyId}}
  kms_master_key_id = {{.resource.KmsMasterKeyId|Quote}}
{{- end}}
}
`

// Claude implements Converter using the Anthropic Messages API.
type Claude struct {
	client anthropic.Client
	model  string
}

// NewClaude builds a client that reads ANTHROPIC_API_KEY from the environment.
func NewClaude() *Claude {
	return &Claude{client: anthropic.NewClient(), model: defaultModel}
}

// Convert implements Converter.
func (c *Claude) Convert(ctx context.Context, req Request) (*Result, error) {
	msg, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     c.model,
		MaxTokens: 4096,
		System:    []anthropic.TextBlockParam{{Text: systemPrompt}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(buildPrompt(req))),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("anthropic: %w", err)
	}

	var sb strings.Builder
	for _, block := range msg.Content {
		if block.Type == "text" {
			sb.WriteString(block.Text)
		}
	}

	return parseResponse(sb.String())
}

func buildPrompt(req Request) string {
	body, err := json.MarshalIndent(req.Resource, "", "  ")
	if err != nil {
		body = fmt.Appendf(nil, "%v", req.Resource)
	}

	hint := req.TFType
	if hint == "" {
		hint = "(infer from source type)"
	}

	return fmt.Sprintf(
		"Provider: %s\nSource type: %s\nSuggested terraform type: %s\nBlock label: %s\n\nResource body:\n%s\n",
		req.Provider, req.SourceType, hint, req.Name, body,
	)
}

var sectionRE = regexp.MustCompile(`(?s)===TFTYPE===\s*(.*?)\s*===HCL===\s*(.*?)\s*===TEMPLATE===\s*(.*)`)

func parseResponse(text string) (*Result, error) {
	m := sectionRE.FindStringSubmatch(text)
	if m == nil {
		return nil, fmt.Errorf("ai: response missing section markers")
	}

	return &Result{
		TFType:   strings.TrimSpace(stripFences(m[1])),
		HCL:      strings.TrimSpace(stripFences(m[2])),
		Template: strings.TrimSpace(stripFences(m[3])) + "\n",
	}, nil
}

func stripFences(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```hcl")
	s = strings.TrimPrefix(s, "```terraform")
	s = strings.TrimPrefix(s, "```gotemplate")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}
