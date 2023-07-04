package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testVDMRoot = "../testdata"

var (
	testSpecFilePath = filepath.Join(testVDMRoot, "vdm.json")
)

func TestSync(t *testing.T) {
	specs, err := vdmspec.GetSpecsFromFile(testSpecFilePath)
	require.NoError(t, err)

	// Need to override for test
	RootFlagValues.SpecFilePath = testSpecFilePath
	err = sync()
	require.NoError(t, err)

	t.Run("spec[0] used a tag", func(t *testing.T) {
		vdmMeta, err := specs[0].GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, "v0.2.0", vdmMeta.Version)
	})

	t.Run("spec[1] used 'latest'", func(t *testing.T) {
		vdmMeta, err := specs[1].GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, "latest", vdmMeta.Version)
	})

	t.Run("spec[2] used a branch", func(t *testing.T) {
		vdmMeta, err := specs[2].GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, "main", vdmMeta.Version)
	})

	t.Run("spec[4] used a hash", func(t *testing.T) {
		vdmMeta, err := specs[3].GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, "2e6657f5ac013296167c4dd92fbb46f0e3dbdc5f", vdmMeta.Version)
	})

	t.Cleanup(func() {
		for _, spec := range specs {
			err := os.RemoveAll(spec.LocalPath)
			assert.NoError(t, err)
		}
	})
}
