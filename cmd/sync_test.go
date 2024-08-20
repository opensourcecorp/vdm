package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/opensourcecorp/vdm/internal/vdminit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testVDMRoot = "../testdata"

var (
	testSpecFilePath = filepath.Join(testVDMRoot, "vdm.yaml")
)

func TestSync(t *testing.T) {
	_, cleanup := vdminit.SetupVDMForTest(t)

	// This runs as part of the outer test container, because we want to inspect
	// the filesystem state as we go
	t.Cleanup(cleanup)

	cmd := newRootCommand()
	cmd.SetArgs([]string{
		"--debug",
		"--specfile-path", testSpecFilePath,
		"sync",
	})

	t.Run("sync works without throwing any errors", func(t *testing.T) {
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("filesytem state is as expected", func(t *testing.T) {
		expectedGitDirs := map[string]string{
			"git-tag":    "vdm",
			"git-branch": "osc-infra",
			"git-hash":   "go-common",
		}
		for topDir, secondDir := range expectedGitDirs {
			sourceRoot := filepath.Join("deps", topDir, secondDir)
			t.Run("source directory exists at its destination", func(t *testing.T) {
				gitTagSource, err := os.Stat(sourceRoot)
				require.NoError(t, err)
				assert.True(t, gitTagSource.IsDir())
			})

			t.Run(".git directory was removed", func(t *testing.T) {
				dotGitDir := filepath.Join(sourceRoot, ".git")
				_, err := os.Stat(dotGitDir)
				assert.ErrorIs(t, err, os.ErrNotExist)
			})

			t.Run("a known file in the remote exists, and is a file", func(t *testing.T) {
				readmePath := filepath.Join(sourceRoot, "README.md")
				_, err := os.Stat(readmePath)
				assert.NoError(t, err)
			})
		}
	})
}
