package vdminit

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/opensourcecorp/vdm/cmd/vars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaths(t *testing.T) {
	// SetupVDMForTest itself calls Paths(), handily
	_, cleanup := SetupVDMForTest(t)

	t.Cleanup(cleanup)

	t.Run("env var is set right", func(t *testing.T) {
		want := filepath.Join(os.TempDir(), "vdm-tmp")
		got := os.Getenv(vars.VDMHomeEnvVarName)
		assert.Equal(t, want, got)
	})

	t.Run("vdm cache is then returned as being under the test homedir", func(t *testing.T) {
		vdmHome, err := vars.GetVDMHomeDir()
		require.NoError(t, err)
		want := filepath.Join(vdmHome, "cache")

		got, err := vars.GetVDMCacheDir()
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("sumdb exists and has sums table in it", func(t *testing.T) {
		vdmCache, err := vars.GetVDMCacheDir()
		require.NoError(t, err)

		sumDBPath := filepath.Join(vdmCache, "sum.db")
		_, err = os.Stat(sumDBPath)
		assert.NoError(t, err)

		db, err := sql.Open("sqlite", sumDBPath)
		require.NoError(t, err)

		want := 1
		var got int
		err = db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'sums';`).Scan(&got)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
