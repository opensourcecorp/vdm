package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSync(t *testing.T) {
	ctx := context.Background()

	const testVDMRoot = "./testdata"
	specFilePath := filepath.Join(testVDMRoot, "vdm.json")

	specs := getSpecsFromFile(ctx, specFilePath)

	sync(ctx, specs)

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

	t.Cleanup(func() {
		for _, spec := range specs {
			os.RemoveAll(spec.LocalPath)
		}
	})
}

func TestShouldKeepGitDir(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, keepGitDirKey{}, true)

	const testVDMRoot = "./testdata"
	specFilePath := filepath.Join(testVDMRoot, "vdm.json")

	specs := getSpecsFromFile(ctx, specFilePath)

	sync(ctx, specs)

	for _, spec := range specs {
		assert.DirExists(t, filepath.Join(spec.LocalPath, ".git"))
	}

	t.Cleanup(func() {
		for _, spec := range specs {
			os.RemoveAll(spec.LocalPath)
		}
	})
}
