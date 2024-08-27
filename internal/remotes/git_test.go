package remotes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckGitAvailable(t *testing.T) {
	t.Run("checkGitAvailable", func(t *testing.T) {
		t.Run("no error when git is available", func(t *testing.T) {
			// Host of this test better have git available lol
			gitAvailable := checkGitAvailable()
			require.NoError(t, gitAvailable)
		})

		t.Run("error when git is NOT available", func(t *testing.T) {
			t.Setenv("PATH", "")
			gitAvailable := checkGitAvailable()
			assert.Error(t, gitAvailable)
		})
	})
}
