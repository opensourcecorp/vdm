package remotes

import (
	"os"
	"testing"

	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestGitSpec() vdmspec.Remote {
	specLocalPath := "./deps/go-common"
	return vdmspec.Remote{
		Type:        "git",
		Source:      "https://github.com/opensourcecorp/go-common",
		Version:     "v0.2.0",
		Destination: specLocalPath,
	}
}

func TestSyncGit(t *testing.T) {
	spec := getTestGitSpec()
	err := SyncGit(spec)
	require.NoError(t, err)

	defer t.Cleanup(func() {
		if cleanupErr := os.RemoveAll(spec.Destination); cleanupErr != nil {
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
	spec := getTestGitSpec()
	cloneErr := gitClone(spec)

	defer t.Cleanup(func() {
		if cleanupErr := os.RemoveAll(spec.Destination); cleanupErr != nil {
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
