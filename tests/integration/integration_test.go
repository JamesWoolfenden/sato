package integration_test

import (
	"os"
	"path"
	"testing"

	sato "sato/src/cf"
	"sato/tests/utils"

	"github.com/stretchr/testify/assert"
)

func TestLocalModules(t *testing.T) {
	t.Parallel()
	t.Run("Test sato conversion", func(t *testing.T) {
		t.Parallel()

		// Use local example file instead of cloning non-existent remote repo
		target := "../../examples/aws-vpc.template.yaml"
		destination := t.TempDir()

		err := sato.Parse(target, destination)
		assert.NoError(t, err, "Failed to parse")

		// Count generated files
		files, err := os.ReadDir(destination)
		assert.NoError(t, err, "Failed to read destination directory")

		// Should generate at least some .tf files
		assert.Greater(t, len(files), 0, "Should generate at least one file")

		// Verify at least one .tf file was created
		hasTfFile := false
		for _, file := range files {
			if path.Ext(file.Name()) == ".tf" {
				hasTfFile = true
				break
			}
		}
		assert.True(t, hasTfFile, "Should generate at least one .tf file")

		err = utils.TfInit(destination)
		assert.NoError(t, err, "Failed to tf init output")
	})
}
