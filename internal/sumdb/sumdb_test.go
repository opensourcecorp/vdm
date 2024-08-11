package sumdb

import (
	"bytes"
	"compress/gzip"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateSHASum(t *testing.T) {
	t.Run("works with an open file", func(t *testing.T) {
		f, err := os.Open("../../testdata/sumdb/sha256test.txt")
		require.NoError(t, err)
		t.Cleanup(func() {
			closeErr := f.Close()
			require.NoError(t, closeErr)
		})

		want := `25fce0ea957324f2fdab37fa2c35df8dc1c62703b1970a744a723603407b630b`
		got, err := CalculateSHASum(f)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("works with an in-memory gzipped directory", func(t *testing.T) {
		var b bytes.Buffer
		zw := gzip.NewWriter(&b)

		f, err := os.Open("../../testdata/sumdb/sha256test.txt")
		require.NoError(t, err)
		t.Cleanup(func() {
			closeErr := f.Close()
			require.NoError(t, closeErr)
		})

		want := `25fce0ea957324f2fdab37fa2c35df8dc1c62703b1970a744a723603407b630b`
		got, err := CalculateSHASum(f)

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}
