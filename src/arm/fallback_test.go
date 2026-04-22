package arm

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"sato/src/ai"
)

type stubConverter struct{}

func (stubConverter) Convert(_ context.Context, req ai.Request) (*ai.Result, error) {
	return &ai.Result{
		TFType:   "azurerm_stub_thing",
		HCL:      "resource \"azurerm_stub_thing\" \"" + req.Name + "\" {}\n",
		Template: "resource \"azurerm_stub_thing\" \"{{.item}}\" {}\n",
	}, nil
}

func TestParseResources_aiFallback(t *testing.T) {
	t.Parallel()

	dest := t.TempDir()

	result := map[string]interface{}{
		"resources": []interface{}{
			map[string]interface{}{
				"type":     "Microsoft.Made/up",
				"name":     "mything",
				"resource": "mything",
			},
		},
	}

	_, err := ParseResources(result, funcMap, dest, WithAIFallback(stubConverter{}))
	if err != nil {
		t.Fatalf("ParseResources: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dest, "azurerm_stub_thing.mything.tf")); err != nil {
		t.Fatalf("hcl not written: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dest, "_drafts", "azurerm_stub_thing.template")); err != nil {
		t.Errorf("draft not written: %v", err)
	}
}
