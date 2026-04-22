package ai

import (
	"strings"
	"testing"
)

func Test_parseResponse(t *testing.T) {
	t.Parallel()

	raw := `Here you go.
===TFTYPE===
aws_sqs_queue
===HCL===
` + "```hcl" + `
resource "aws_sqs_queue" "x" {
  name = "x"
}
` + "```" + `
===TEMPLATE===
resource "aws_sqs_queue" "{{.item}}" {
  name = {{.resource.QueueName|Quote}}
}
`

	got, err := parseResponse(raw)
	if err != nil {
		t.Fatalf("parseResponse: %v", err)
	}
	if got.TFType != "aws_sqs_queue" {
		t.Errorf("TFType = %q", got.TFType)
	}
	if !strings.HasPrefix(got.HCL, `resource "aws_sqs_queue" "x"`) {
		t.Errorf("HCL not stripped of fences: %q", got.HCL)
	}
	if !strings.Contains(got.Template, "{{.item}}") {
		t.Errorf("Template missing .item: %q", got.Template)
	}
}

func Test_parseResponse_missingMarkers(t *testing.T) {
	t.Parallel()

	if _, err := parseResponse("nope"); err == nil {
		t.Fatal("expected error for malformed response")
	}
}

func Test_buildPrompt(t *testing.T) {
	t.Parallel()

	p := buildPrompt(Request{
		SourceType: "AWS::SQS::Queue",
		TFType:     "",
		Provider:   "aws",
		Name:       "myqueue",
		Resource:   map[string]any{"QueueName": "q"},
	})

	for _, want := range []string{"AWS::SQS::Queue", "(infer from source type)", "myqueue", "QueueName"} {
		if !strings.Contains(p, want) {
			t.Errorf("prompt missing %q:\n%s", want, p)
		}
	}
}

func Test_WriteDraft(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	if err := WriteDraft(dir, "aws_sqs_queue", "resource ...\n"); err != nil {
		t.Fatalf("WriteDraft: %v", err)
	}
	if err := WriteDraft(dir, "", "x"); err != nil {
		t.Errorf("WriteDraft with empty tfType should noop, got %v", err)
	}
}
