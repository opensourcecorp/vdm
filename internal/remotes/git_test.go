package remotes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestGitRemote(t *testing.T) (Git, string) {
	t.Helper()
	specLocalPath := "./deps/go-common"
	remote := Git{
		RemoteTemplate: vdmspec.RemoteTemplate{
			Type:        "git",
			Source:      "https://github.com/opensourcecorp/go-common",
			Version:     "v0.2.0",
			Destination: specLocalPath,
		},
	}
	dest := filepath.Join(os.TempDir(), "vdm-test", filepath.Base(remote.Source))
	return remote, dest
}

func TestSyncGit(t *testing.T) {
	remote, dest := getTestGitRemote(t)
	err := remote.Sync()
	require.NoError(t, err)

	defer t.Cleanup(func() {
		if cleanupErr := os.RemoveAll(dest); cleanupErr != nil {
			t.Fatalf("removing specLocalPath: %v", cleanupErr)
		}
	})

	t.Run(".git directory was removed", func(t *testing.T) {
		_, err := os.Stat("./deps/go-common-tag/.git")
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
}

func TestCheckGitAvailable(t *testing.T) {
	t.Run("checkGitAvailable", func(t *testing.T) {
		t.Run("no error when git is available", func(t *testing.T) {
			// Host of this test better have git available lol
			gitAvailable := checkGitAvailable()
			require.NoError(t, gitAvailable)
		})

		t.Run("error when git is NOT available", func(t *testing.T) {
			t.Setenv("PATH", "")
			gitAvailable := checkGitAvailable()
			assert.Error(t, gitAvailable)
		})
	})
}

func TestGitClone(t *testing.T) {
	remote, dest := getTestGitRemote(t)
	cloneErr := gitClone(remote.Source, dest)

	defer t.Cleanup(func() {
		if cleanupErr := os.RemoveAll(dest); cleanupErr != nil {
			t.Fatalf("removing specLocalPath: %v", cleanupErr)
		}
	})

	t.Run("no error on success", func(t *testing.T) {
		require.NoError(t, cloneErr)
	})

	t.Run("LocalPath is a directory, not a file", func(t *testing.T) {
		outDir, err := os.Stat(dest)
		require.NoError(t, err)
		assert.True(t, outDir.IsDir())
	})

	t.Run("a known file in the remote exists, and is a file", func(t *testing.T) {
		sampleFile, err := os.Stat(filepath.Join(dest, "go.mod"))
		require.NoError(t, err)
		assert.False(t, sampleFile.IsDir())
	})
}
