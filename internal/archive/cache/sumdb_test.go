package cache

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/opensourcecorp/vdm/internal/archive"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateSHASum(t *testing.T) {
	t.Run("works with an open file", func(t *testing.T) {
		f, err := os.Open("../../../testdata/sumdb/sha256test.txt")
		require.NoError(t, err)
		t.Cleanup(func() {
			closeErr := f.Close()
			require.NoError(t, closeErr)
		})

		want := `25fce0ea957324f2fdab37fa2c35df8dc1c62703b1970a744a723603407b630b`
		got, err := CalculateSHASum(f)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("works on a created archive file", func(t *testing.T) {
		rootDir := "../../../testdata/filetree"
		archivePath := filepath.Join(os.TempDir(), "archive-to-hash.tar.gz")
		err := archive.CreateArchive(rootDir, archivePath)
		require.NoError(t, err)

		f, err := os.Open(archivePath)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := os.RemoveAll(archivePath)
			require.NoError(t, err)
		})

		want := `c1c5e8e5cd54819ea0243db2b203e581aae7308937fb575204914fc2e41ef2a7`
		got, err := CalculateSHASum(f)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestStringAsBase64(t *testing.T) {
	t.Run("works on string with no resulting padding", func(t *testing.T) {
		src := "https://github.com/org/user"
		want := "aHR0cHM6Ly9naXRodWIuY29tL29yZy91c2Vy"
		got := StringToBase64(src)
		assert.Equal(t, want, got)
	})

	t.Run("works on string WITH resulting padding", func(t *testing.T) {
		src := "https://github.com/org/user12"
		want := "aHR0cHM6Ly9naXRodWIuY29tL29yZy91c2VyMTI_EQ"
		got := StringToBase64(src)
		assert.Equal(t, want, got)
	})
}

func TestStringFromBase64(t *testing.T) {
	t.Run("works on string with no resulting padding", func(t *testing.T) {
		src := "aHR0cHM6Ly9naXRodWIuY29tL29yZy91c2Vy"
		want := "https://github.com/org/user"
		got, err := StringFromBase64(src)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("works on string WITH resulting padding", func(t *testing.T) {
		src := "aHR0cHM6Ly9naXRodWIuY29tL29yZy91c2VyMTI_EQ"
		want := "https://github.com/org/user12"
		got, err := StringFromBase64(src)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestNonAlphaBase64(t *testing.T) {
	t.Run("replacer works", func(t *testing.T) {
		want := "aHR0cHM6Ly9naXRodWIuY29tL29yZy91c2Vy_EQ"
		got := replaceNonAlphaBase64Characters("aHR0cHM6Ly9naXRodWIuY29tL29yZy91c2Vy=")
		assert.Equal(t, want, got)
	})

	t.Run("restorer works", func(t *testing.T) {
		want := "aHR0cHM6Ly9naXRodWIuY29tL29yZy91c2Vy="
		got := restoreNonAlphaBase64Characters("aHR0cHM6Ly9naXRodWIuY29tL29yZy91c2Vy_EQ")
		assert.Equal(t, want, got)
	})
}
