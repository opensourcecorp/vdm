package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckAvailable(t *testing.T) {
	ctx := context.Background()

	t.Run("git", func(t *testing.T) {
		// Host of this test better have git available lol
		gitAvailable := checkGitAvailable(ctx)
		assert.NoError(t, gitAvailable)

		t.Setenv("PATH", "")
		gitAvailable = checkGitAvailable(ctx)
		assert.Error(t, gitAvailable)
	})
}
