package cmd

import (
	"path/filepath"
	"testing"

	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testVDMRoot = "../testdata"

var (
	testSpecFilePath = filepath.Join(testVDMRoot, "vdm.yaml")
)

func TestSync(t *testing.T) {
	spec, err := vdmspec.GetSpecFromFile(testSpecFilePath)
	require.NoError(t, err)

	// Need to override for test
	RootFlagValues.SpecFilePath = testSpecFilePath
	err = sync()
	require.NoError(t, err)

	// defer t.Cleanup(func() {
	// 	for _, remote := range spec.Remotes {
	// 		err := os.RemoveAll(remote.LocalPath)
	// 		require.NoError(t, err)
	// 	}
	// })

	t.Run("SyncGit", func(t *testing.T) {
		t.Run("spec[0] used a tag", func(t *testing.T) {
			vdmMeta, err := spec.Remotes[0].GetVDMMeta()
			require.NoError(t, err)
			assert.Equal(t, "v0.2.0", vdmMeta.Version)
		})

		t.Run("spec[1] used 'latest'", func(t *testing.T) {
			vdmMeta, err := spec.Remotes[1].GetVDMMeta()
			require.NoError(t, err)
			assert.Equal(t, "latest", vdmMeta.Version)
		})

		t.Run("spec[2] used a branch", func(t *testing.T) {
			vdmMeta, err := spec.Remotes[2].GetVDMMeta()
			require.NoError(t, err)
			assert.Equal(t, "main", vdmMeta.Version)
		})

		t.Run("spec[3] used a hash", func(t *testing.T) {
			vdmMeta, err := spec.Remotes[3].GetVDMMeta()
			require.NoError(t, err)
			assert.Equal(t, "2e6657f5ac013296167c4dd92fbb46f0e3dbdc5f", vdmMeta.Version)
		})
	})

	t.Run("SyncFile", func(t *testing.T) {
		t.Run("spec[4] had an implicit version", func(t *testing.T) {
			vdmMeta, err := spec.Remotes[4].GetVDMMeta()
			require.NoError(t, err)
			assert.Equal(t, "", vdmMeta.Version)
		})
	})
}
