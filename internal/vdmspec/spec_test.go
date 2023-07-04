package vdmspec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func TestVDMMeta(t *testing.T) {
	const testVDMRoot = "./testdata"
	testVDMMetaFilePath := filepath.Join(testVDMRoot, MetaFileName)

	t.Run("GetVDMMeta", func(t *testing.T) {
		spec := VDMSpec{
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

		got, err := spec.GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, spec, got)

		t.Cleanup(func() {
			err := os.RemoveAll(testVDMMetaFilePath)
			assert.NoError(t, err)
		})
	})

	t.Run("WriteVDMMeta", func(t *testing.T) {
		spec := VDMSpec{
			Remote:    "https://some-remote",
			Version:   "v1.0.0",
			LocalPath: "./testdata",
		}
		err := spec.WriteVDMMeta()
		require.NoError(t, err)

		got, err := spec.GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, spec, got)

		t.Cleanup(func() {
			err := os.RemoveAll(testVDMMetaFilePath)
			assert.NoError(t, err)
		})
	})

	t.Run("GetSpecsFromFile", func(t *testing.T) {
		specFilePath := "./testdata/vdm.json"

		specs, err := GetSpecsFromFile(specFilePath)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(specs))

		t.Cleanup(func() {
			err := os.RemoveAll(testVDMMetaFilePath)
			assert.NoError(t, err)
		})
	})
}
