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

	testRemote = Remote{
		Remote:    "https://some-remote",
		Version:   "v1.0.0",
		LocalPath: testVDMRoot,
	}

	testSpecFilePath = filepath.Join(testVDMRoot, "vdm.yaml")

	testVDMMetaContents = fmt.Sprintf(
		`{"remote": "https://some-remote", "version": "v1.0.0", "local_path": "%s"}`,
		testVDMRoot,
	)
)

func TestVDMMeta(t *testing.T) {
	t.Run("GetVDMMeta", func(t *testing.T) {
		err := os.WriteFile(testVDMMetaFilePath, []byte(testVDMMetaContents), 0644)
		require.NoError(t, err)

		got, err := testRemote.GetVDMMeta()
		require.NoError(t, err)
		assert.Equal(t, testRemote, got)

		t.Cleanup(func() {
			err := os.RemoveAll(testVDMMetaFilePath)
			require.NoError(t, err)
		})
	})

	t.Run("WriteVDMMeta", func(t *testing.T) {
		// Needs to have parent dir(s) exist for write to work
		err := os.MkdirAll(testRemote.LocalPath, 0644)
		require.NoError(t, err)

		err = testRemote.WriteVDMMeta()
		require.NoError(t, err)

		got, err := testRemote.GetVDMMeta()
		require.NoError(t, err)
		assert.Equal(t, testRemote, got)

		t.Cleanup(func() {
			err := os.RemoveAll(testVDMMetaFilePath)
			require.NoError(t, err)
		})
	})

	t.Run("GetSpecsFromFile", func(t *testing.T) {
		spec, err := GetSpecFromFile(testSpecFilePath)
		require.NoError(t, err)
		assert.Equal(t, 5, len(spec.Remotes))

		t.Cleanup(func() {
			err := os.RemoveAll(testVDMMetaFilePath)
			require.NoError(t, err)
		})
	})
}
