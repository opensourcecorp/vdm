package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/stretchr/testify/assert"
)

func TestSync(t *testing.T) {
	const testVDMRoot = "./testdata"

	specFilePath := filepath.Join(testVDMRoot, "vdm.json")

	specs, err := vdmspec.GetSpecsFromFile(specFilePath)
	assert.NoError(t, err)

	sync()

	t.Run("spec[0] used a tag", func(t *testing.T) {
		vdmMeta, err := specs[0].GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, vdmMeta.Version, "v0.2.0")
	})

	t.Run("spec[1] used 'latest'", func(t *testing.T) {
		vdmMeta, err := specs[1].GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, vdmMeta.Version, "latest")
	})

	t.Run("spec[2] used a branch", func(t *testing.T) {
		vdmMeta, err := specs[2].GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, vdmMeta.Version, "main")
	})

	t.Run("spec[4] used a hash", func(t *testing.T) {
		vdmMeta, err := specs[3].GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, vdmMeta.Version, "2e6657f5ac013296167c4dd92fbb46f0e3dbdc5f")
	})

	t.Cleanup(func() {
		for _, spec := range specs {
			err := os.RemoveAll(spec.LocalPath)
			assert.NoError(t, err)
		}
	})
}
