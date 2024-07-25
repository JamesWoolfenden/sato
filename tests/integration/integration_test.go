package integration_test

import (
	"os"
	"path"
	sato "sato/src/cf"
	"sato/tests/utils"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalModules(t *testing.T) {
	t.Parallel()
	t.Run("Test sato conversion", func(t *testing.T) {
		t.Parallel()

		localTagExampleRepo := "https://github.com/JamesWoolfenden/aws-cloudformation-templates"
		repoPath := utils.CloneRepo(localTagExampleRepo, "6758465103b4431e9de4a93d30faff7912204847")

		defer func() {
			_ = os.RemoveAll(repoPath)
		}()

		target := path.Join(repoPath, "aws/solutions/ADConnector/templates/Adconnector.yaml")
		destination := path.Join(repoPath, "aws/solutions/ADConnector/templates/Adconnector/")

		err := sato.Parse(target, destination)
		if err != nil {
			assert.Fail(t, "Failed to parse")
		}

		// Count
		files, _ := os.ReadDir(destination)
		expect := 14
		assert.Equal(
			t, expect, len(files),
			"The number of files found "+strconv.Itoa(len(files))+" should match expectation "+strconv.Itoa(expect))

		err = utils.TfInit(destination)
		assert.Nil(t, err, "Failed to tf init output")
	})
}
