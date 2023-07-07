package vdmspec

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testVDMRoot = "../../testdata"

var (
	testVDMMetaFilePath = filepath.Join(testVDMRoot, MetaFileName)

	testSpec = VDMSpec{
		Remote:    "https://some-remote",
		Version:   "v1.0.0",
		LocalPath: testVDMRoot,
	}

	testSpecFilePath = filepath.Join(testVDMRoot, "vdm.json")

	testVDMMetaContents = fmt.Sprintf(
		`{"remote": "https://some-remote", "version": "v1.0.0", "local_path": "%s"}`,
		testVDMRoot,
	)
)

func TestVDMMeta(t *testing.T) {
	t.Run("GetVDMMeta", func(t *testing.T) {
		err := os.WriteFile(testVDMMetaFilePath, []byte(testVDMMetaContents), 0644)
		require.NoError(t, err)

		got, err := testSpec.GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, testSpec, got)

		t.Cleanup(func() {
			err := os.RemoveAll(testVDMMetaFilePath)
			assert.NoError(t, err)
		})
	})

	t.Run("WriteVDMMeta", func(t *testing.T) {
		// Needs to have parent dir(s) exist for write to work
		err := os.MkdirAll(testSpec.LocalPath, 0644)
		require.NoError(t, err)

		err = testSpec.WriteVDMMeta()
		require.NoError(t, err)

		got, err := testSpec.GetVDMMeta()
		assert.NoError(t, err)
		assert.Equal(t, testSpec, got)

		t.Cleanup(func() {
			err := os.RemoveAll(testVDMMetaFilePath)
			assert.NoError(t, err)
		})
	})

	t.Run("GetSpecsFromFile", func(t *testing.T) {
		specs, err := GetSpecsFromFile(testSpecFilePath)
		assert.NoError(t, err)
		assert.Equal(t, 5, len(specs))

		t.Cleanup(func() {
			err := os.RemoveAll(testVDMMetaFilePath)
			assert.NoError(t, err)
		})
	})
}
