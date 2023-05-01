package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	ctx := context.Background()

	t.Run("passes", func(t *testing.T) {
		spec := vdmSpec{
			Remote:    "https://some-remote",
			Version:   "v1.0.0",
			LocalPath: "./deps/some-remote",
		}
		err := spec.Validate(ctx)
		assert.NoError(t, err)
	})

	t.Run("fails on zero-length remote", func(t *testing.T) {
		spec := vdmSpec{
			Remote:    "",
			Version:   "v1.0.0",
			LocalPath: "./deps/some-remote",
		}
		err := spec.Validate(ctx)
		assert.Error(t, err)
	})

	t.Run("fails on remote without valid protocol", func(t *testing.T) {
		spec := vdmSpec{
			Remote:    "some-remote",
			Version:   "v1.0.0",
			LocalPath: "./deps/some-remote",
		}
		err := spec.Validate(ctx)
		assert.Error(t, err)
	})

	t.Run("fails on zero-length version", func(t *testing.T) {
		spec := vdmSpec{
			Remote:    "https://some-remote",
			Version:   "",
			LocalPath: "./deps/some-remote",
		}
		err := spec.Validate(ctx)
		assert.Error(t, err)
	})

	t.Run("fails on zero-length local path", func(t *testing.T) {
		spec := vdmSpec{
			Remote:    "https://some-remote",
			Version:   "v1.0.0",
			LocalPath: "",
		}
		err := spec.Validate(ctx)
		assert.Error(t, err)
	})
}
