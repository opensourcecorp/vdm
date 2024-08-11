package archive

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateArchive(t *testing.T) {
	archiveRoot := "../../testdata/filetree"
	archivePath := filepath.Join(os.TempDir(), "archive-test.tar.gz")
	t.Cleanup(func() {
		err := os.RemoveAll(archivePath)
		require.NoError(t, err)
	})

	err := CreateArchive(archiveRoot, archivePath)
	require.NoError(t, err)

	// TODO: add test for inspecting contents once we implement an archive
	// extractor -- as of now, I'm just checking on the CLI if the expected tree
	// matches
}

func TestMaybeGetTopLevelDir(t *testing.T) {
	t.Run("works when rootDir is a directory", func(t *testing.T) {
		rootDir := "../../testdata/filetree"
		wantTopLevelDir := "filetree"
		gotTopLevelDir, err := maybeGetTopLevelDir(rootDir)

		assert.NoError(t, err)
		assert.Equal(t, wantTopLevelDir, gotTopLevelDir)
	})

	t.Run("works when rootDir is a file", func(t *testing.T) {
		rootDir := "../../testdata/filetree/top-file"
		wantTopLevelDir := ""
		gotTopLevelDir, err := maybeGetTopLevelDir(rootDir)

		assert.NoError(t, err)
		assert.Equal(t, wantTopLevelDir, gotTopLevelDir)
	})
}
