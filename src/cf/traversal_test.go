package cf

import (
	"os"
	"path/filepath"
	"testing"
)

// A logical ID containing path segments must not write files outside the destination (GHSA-prgx-fhgv-v466).
func TestParse_rejectsTraversingLogicalID(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "out", "sub")
	template := filepath.Join(root, "evil.yaml")

	evil := `AWSTemplateFormatVersion: "2010-09-09"
Resources:
  "../../conftest":
    Type: AWS::Lambda::Function
    Properties:
      Runtime: python3.9
      Handler: index.handler
      Role: x
      Code:
        ZipFile: |
          print("pwned")
`
	if err := os.WriteFile(template, []byte(evil), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := Parse(template, dest); err == nil {
		t.Error("Parse() accepted a logical ID containing path separators")
	}

	for _, escaped := range []string{
		filepath.Join(root, "conftest.py"),
		filepath.Join(root, "conftest.zip"),
		filepath.Join(root, "out", "conftest.py"),
	} {
		if _, err := os.Stat(escaped); err == nil {
			t.Errorf("file written outside the destination: %s", escaped)
		}
	}
}
