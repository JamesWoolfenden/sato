package tfgen

import (
	"os"
	"path/filepath"
	"testing"
)

//goland:noinspection GoLinter
func TestWrite(t *testing.T) {
	t.Parallel()

	type args struct {
		output   string
		location string
		name     string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Pass", args{"gibberine", ".", "test"}, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := Write(tt.args.output, tt.args.location, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrite_rejectsTraversal(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	dest := filepath.Join(root, "out")

	for _, name := range []string{"../escape", "../../escape", "a/b", `..\escape`} {
		if err := Write("content", dest, name); err == nil {
			t.Errorf("Write(name=%q) = nil, want error", name)
		}
	}

	if _, err := os.Stat(filepath.Join(root, "escape.tf")); err == nil {
		t.Error("Write escaped the destination directory")
	}
}
