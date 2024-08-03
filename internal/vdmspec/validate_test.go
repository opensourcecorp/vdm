package vdmspec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Run("passes", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Source:      "https://some-remote",
				Version:     "v1.0.0",
				Destination: "./deps/some-remote",
			}},
		}
		err := spec.Validate()
		require.NoError(t, err)
	})

	t.Run("fails on zero-length remote", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Source:      "",
				Version:     "v1.0.0",
				Destination: "./deps/some-remote",
			}},
		}
		err := spec.Validate()
		assert.Error(t, err)
	})

	t.Run("fails on remote without valid protocol", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Source:      "some-remote",
				Version:     "v1.0.0",
				Destination: "./deps/some-remote",
			}},
		}
		err := spec.Validate()
		assert.Error(t, err)
	})

	t.Run("fails on zero-length version for git remote type", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Source:      "https://some-remote",
				Version:     "",
				Destination: "./deps/some-remote",
				Type:        GitType,
			}},
		}
		err := spec.Validate()
		assert.Error(t, err)
	})

	t.Run("fails on unrecognized remote type", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Source:      "https://some-remote",
				Version:     "",
				Destination: "./deps/some-remote",
				Type:        "bad",
			}},
		}
		err := spec.Validate()
		assert.Error(t, err)
	})

	t.Run("fails on zero-length local path", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Source:      "https://some-remote",
				Version:     "v1.0.0",
				Destination: "",
			}},
		}
		err := spec.Validate()
		assert.Error(t, err)
	})
}
