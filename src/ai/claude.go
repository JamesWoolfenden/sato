package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/vertex"
	"github.com/rs/zerolog/log"
)

const callTimeout = 60 * time.Second

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
	once   sync.Once
	client anthropic.Client
	model  string
}

// NewClaude returns a lazily-initialised client. If ANTHROPIC_VERTEX_PROJECT_ID
// is set it routes via Google Vertex AI (using CLOUD_ML_REGION and ADC);
// otherwise it reads ANTHROPIC_API_KEY for the direct Anthropic API.
// ANTHROPIC_MODEL overrides the default model in either mode.
func NewClaude() *Claude {
	return &Claude{}
}

func (c *Claude) init(ctx context.Context) {
	c.model = defaultModel
	if m := os.Getenv("ANTHROPIC_MODEL"); m != "" {
		c.model = m
	}

	var opts []option.RequestOption
	if project := os.Getenv("ANTHROPIC_VERTEX_PROJECT_ID"); project != "" {
		region := os.Getenv("CLOUD_ML_REGION")
		log.Info().Msgf("ai: using Vertex AI (project=%s region=%s model=%s)", project, region, c.model)
		opts = append(opts, vertex.WithGoogleAuth(ctx, region, project))
	} else {
		log.Info().Msgf("ai: using Anthropic API (model=%s)", c.model)
	}

	c.client = anthropic.NewClient(opts...)
}

// Convert implements Converter.
func (c *Claude) Convert(ctx context.Context, req Request) (*Result, error) {
	ctx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	c.once.Do(func() { c.init(ctx) })

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
