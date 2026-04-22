package cf

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sato/src/ai"
)

type stubConverter struct {
	got *ai.Request
}

func (s *stubConverter) Convert(_ context.Context, req ai.Request) (*ai.Result, error) {
	s.got = &req

	return &ai.Result{
		TFType:   "aws_stub_thing",
		HCL:      "resource \"aws_stub_thing\" \"" + req.Name + "\" {}\n",
		Template: "resource \"aws_stub_thing\" \"{{.item}}\" {}\n",
	}, nil
}

func Test_aiFallback(t *testing.T) {
	t.Parallel()

	dest := t.TempDir()
	stub := &stubConverter{}

	err := aiFallback(stub, "AWS::Made::Up", "MyThing", map[string]any{"Foo": "bar"}, dest)
	if err != nil {
		t.Fatalf("aiFallback: %v", err)
	}

	if stub.got == nil || stub.got.Provider != "aws" || stub.got.SourceType != "AWS::Made::Up" {
		t.Fatalf("request not populated: %+v", stub.got)
	}

	hcl, err := os.ReadFile(filepath.Join(dest, "aws_stub_thing.mything.tf"))
	if err != nil {
		t.Fatalf("read hcl: %v", err)
	}

	if !strings.HasPrefix(string(hcl), ai.Header) {
		t.Errorf("missing AI header in %q", hcl)
	}

	if _, err := os.Stat(filepath.Join(dest, "_drafts", "aws_stub_thing.template")); err != nil {
		t.Errorf("draft not written: %v", err)
	}
}
