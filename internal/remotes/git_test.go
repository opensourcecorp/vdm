package remotes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestGitSpec() vdmspec.Remote {
	specLocalPath := "./deps/go-common"
	return vdmspec.Remote{
		Type:      "git",
		Remote:    "https://github.com/opensourcecorp/go-common",
		Version:   "v0.2.0",
		LocalPath: specLocalPath,
	}
}

func TestSyncGit(t *testing.T) {
	t.Run("with keepGitDir=false", func(t *testing.T) {
		spec := getTestGitSpec()
		err := SyncGit(spec, false)
		require.NoError(t, err)
		defer t.Cleanup(func() {
			if cleanupErr := os.RemoveAll(spec.LocalPath); cleanupErr != nil {
				t.Fatalf("removing specLocalPath: %v", cleanupErr)
			}
		})

		_, err = os.Stat(filepath.Join(spec.LocalPath, ".git"))
		assert.ErrorIs(t, err, os.ErrNotExist, ".git directory should be removed")
	})

	t.Run("with keepGitDir=true", func(t *testing.T) {
		spec := getTestGitSpec()
		err := SyncGit(spec, true)
		require.NoError(t, err)
		defer t.Cleanup(func() {
			if cleanupErr := os.RemoveAll(spec.LocalPath); cleanupErr != nil {
				t.Fatalf("removing specLocalPath: %v", cleanupErr)
			}
		})

		_, err = os.Stat(filepath.Join(spec.LocalPath, ".git"))
		assert.NoError(t, err, ".git directory should exist")
		assert.DirExists(t, filepath.Join(spec.LocalPath, ".git"), ".git should be a directory")
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
	spec := getTestGitSpec()
	cloneErr := gitClone(spec)

	defer t.Cleanup(func() {
		if cleanupErr := os.RemoveAll(spec.LocalPath); cleanupErr != nil {
			t.Fatalf("removing specLocalPath: %v", cleanupErr)
		}
	})

	t.Run("no error on success", func(t *testing.T) {
		require.NoError(t, cloneErr)
	})

	t.Run("LocalPath is a directory, not a file", func(t *testing.T) {
		outDir, err := os.Stat("./deps/go-common")
		require.NoError(t, err)
		assert.True(t, outDir.IsDir())
	})

	t.Run("a known file in the remote exists, and is a file", func(t *testing.T) {
		sampleFile, err := os.Stat("./deps/go-common/go.mod")
		require.NoError(t, err)
		assert.False(t, sampleFile.IsDir())
	})
}
