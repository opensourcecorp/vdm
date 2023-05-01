package main

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSync(t *testing.T) {
	const testVDMRoot = "./testdata"
	specFilePath := filepath.Join(testVDMRoot, ".vdm")

	specs := getSpecsFromFile(specFilePath, runFlags{SpecFilePath: specFilePath})

	sync(specs)

	t.Run("spec[0] used a tag", func(t *testing.T) {
		vdmMeta := specs[0].getVDMMeta()
		assert.Equal(t, vdmMeta.Version, "v0.2.0")
	})

	t.Run("spec[1] used 'latest'", func(t *testing.T) {
		vdmMeta := specs[1].getVDMMeta()
		assert.Equal(t, vdmMeta.Version, "latest")
	})

	t.Run("spec[2] used a branch", func(t *testing.T) {
		vdmMeta := specs[2].getVDMMeta()
		assert.Equal(t, vdmMeta.Version, "main")
	})

	t.Run("spec[4] used a hash", func(t *testing.T) {
		vdmMeta := specs[3].getVDMMeta()
		assert.Equal(t, vdmMeta.Version, "2e6657f5ac013296167c4dd92fbb46f0e3dbdc5f")
	})
}
