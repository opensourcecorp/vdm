package remote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckGitAvailable(t *testing.T) {
	t.Run("git", func(t *testing.T) {
		// Host of this test better have git available lol
		gitAvailable := checkGitAvailable()
		assert.NoError(t, gitAvailable)

		t.Setenv("PATH", "")
		gitAvailable = checkGitAvailable()
		assert.Error(t, gitAvailable)
	})
}
