package utils

import (
	"context"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/rs/zerolog/log"
)

// CloneRepo test utility
func CloneRepo(repoPath string, commitHash string) string {
	dir, err := os.MkdirTemp("", "temp-repo")
	if err != nil {
		log.Fatal().Err(err)
	}

	// Clones the repository into the given dir, just as a normal git clone does
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: repoPath,
	})

	if err != nil {
		log.Fatal().Err(err)
	}

	if commitHash != "" {
		wt, _ := repo.Worktree()

		commitRef := plumbing.NewHash(commitHash)
		_ = wt.Checkout(&git.CheckoutOptions{Hash: commitRef})
	}

	if err != nil {
		log.Fatal().Err(err)
	}

	return dir
}

// TfInit supports checking the results of Sato parse
func TfInit(workingDir string) error {
	installer := &releases.LatestVersion{
		Product: product.Terraform,
	}
	execPath, err := installer.Install(context.Background())

	if err != nil {
		return err
	}

	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		return err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return err
	}

	return nil
}
