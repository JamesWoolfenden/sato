package ai

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteDraft_rejectsTraversal(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	dest := filepath.Join(root, "out")

	if err := WriteDraft(dest, "../../evil", "x"); err == nil {
		t.Error("WriteDraft accepted a traversing type name")
	}

	if _, err := os.Stat(filepath.Join(root, "evil.template")); err == nil {
		t.Error("WriteDraft escaped the destination")
	}

	if err := WriteDraft(dest, "aws_s3_bucket", "x"); err != nil {
		t.Errorf("WriteDraft(valid) = %v", err)
	}
}
