package vdmspec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("passes", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Remote:    "https://some-remote",
				Version:   "v1.0.0",
				LocalPath: "./deps/some-remote",
			}},
		}
		err := spec.Validate()
		assert.NoError(t, err)
	})

	t.Run("fails on zero-length remote", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Remote:    "",
				Version:   "v1.0.0",
				LocalPath: "./deps/some-remote",
			}},
		}
		err := spec.Validate()
		assert.Error(t, err)
	})

	t.Run("fails on remote without valid protocol", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Remote:    "some-remote",
				Version:   "v1.0.0",
				LocalPath: "./deps/some-remote",
			}},
		}
		err := spec.Validate()
		assert.Error(t, err)
	})

	t.Run("fails on zero-length version for git remote type", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Remote:    "https://some-remote",
				Version:   "",
				LocalPath: "./deps/some-remote",
				Type:      "git",
			}},
		}
		err := spec.Validate()
		assert.Error(t, err)
	})

	t.Run("fails on unrecognized remote type", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Remote:    "https://some-remote",
				Version:   "",
				LocalPath: "./deps/some-remote",
				Type:      "bad",
			}},
		}
		err := spec.Validate()
		assert.Error(t, err)
	})

	t.Run("fails on zero-length local path", func(t *testing.T) {
		spec := Spec{
			Remotes: []Remote{{
				Remote:    "https://some-remote",
				Version:   "v1.0.0",
				LocalPath: "",
			}},
		}
		err := spec.Validate()
		assert.Error(t, err)
	})
}
