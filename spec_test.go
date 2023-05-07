package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecGetVDMMeta(t *testing.T) {
	const testVDMRoot = "./testdata"
	testVDMMetaFilePath := filepath.Join(testVDMRoot, "VDMMETA")

	t.Run("getVDMMeta", func(t *testing.T) {
		spec := vdmSpec{
			Remote:    "https://some-remote",
			Version:   "v1.0.0",
			LocalPath: "./testdata",
		}
		vdmMetaContents := `
		{
			"remote": "https://some-remote",
			"version": "v1.0.0",
			"local_path": "./testdata"
		}`
		err := os.WriteFile(testVDMMetaFilePath, []byte(vdmMetaContents), 0644)
		if err != nil {
			t.Fatal(err)
		}

		got := spec.getVDMMeta()
		assert.Equal(t, spec, got)

		t.Cleanup(func() {
			os.RemoveAll(testVDMMetaFilePath)
		})
	})

	t.Run("writeVDMMeta", func(t *testing.T) {
		spec := vdmSpec{
			Remote:    "https://some-remote",
			Version:   "v1.0.0",
			LocalPath: "./testdata",
		}
		err := spec.writeVDMMeta()
		assert.NoError(t, err)

		got := spec.getVDMMeta()
		assert.Equal(t, spec, got)

		t.Cleanup(func() {
			os.RemoveAll(testVDMMetaFilePath)
		})
	})

	t.Run("getSpecsFromFile", func(t *testing.T) {
		specFilePath := "./testdata/vdm.json"

		specs := getSpecsFromFile(context.Background(), specFilePath)
		assert.Equal(t, 4, len(specs))

		t.Cleanup(func() {
			os.RemoveAll(testVDMMetaFilePath)
		})
	})
}
