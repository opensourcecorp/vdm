package sumdb

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/opensourcecorp/vdm/cmd/vars"
	"github.com/opensourcecorp/vdm/internal/archive"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
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
	src := "https://github.com/org/user"
	want := "aHR0cHM6Ly9naXRodWIuY29tL29yZy91c2Vy"
	got := StringAsBase64(src)

	assert.Equal(t, want, got)
}

func TestCacheRemote(t *testing.T) {
	t.Setenv(vars.VDMHomeEnvVarName, filepath.Join(os.TempDir(), "vdmhome"))

	f, err := os.Open("../../../testdata/sumdb/sha256test.txt")
	require.NoError(t, err)

	remote := vdmspec.Remote{
		Source:  "https://github.com/org/user",
		Version: "v1.0.0",
	}
	cachedPath := filepath.Join(os.Getenv(vars.VDMHomeEnvVarName), "cache", StringAsBase64(remote.Source))

	err = CacheRemote(remote, f)
	assert.NoError(t, err)

	want := fmt.Sprintf("%s %s %s", remote.Source, remote.Version, "25fce0ea957324f2fdab37fa2c35df8dc1c62703b1970a744a723603407b630b")
	got, err := os.ReadFile(cachedPath)
	require.NoError(t, err)

	assert.Equal(t, want, string(got))
}
