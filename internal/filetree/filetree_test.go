package filetree

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFilePathsInDirectory(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		root, err := filepath.Abs("../../testdata/filetree")
		require.NoError(t, err)

		want := []string{
			filepath.Join(root, "top-file"),
			filepath.Join(root, "subdir", "subdir-file"),
			filepath.Join(root, "subdir", "subdir-2", "subdir-2-file"),
		}
		got, err := GetFilePathsInDirectory(root)
		require.NoError(t, err)

		assert.ElementsMatch(t, want, got)
	})

	t.Run("fails in some way, such as when given nonexistent root", func(t *testing.T) {
		badRoot := "dir-that-doesnt-exist"
		_, err := GetFilePathsInDirectory(badRoot)
		assert.Error(t, err)
	})
}
