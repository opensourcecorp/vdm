package cachetest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/opensourcecorp/vdm/cmd/vars"
	"github.com/opensourcecorp/vdm/internal/archive/cache"
	"github.com/opensourcecorp/vdm/internal/remotes"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCacheRemote(t *testing.T) {
	vdmHomePath := filepath.Join(os.TempDir(), "vdmhome")
	t.Setenv(vars.VDMHomeEnvVarName, vdmHomePath)
	t.Cleanup(func() {
		err := os.RemoveAll(vdmHomePath)
		require.NoError(t, err)
	})

	f, err := os.Open("../../../../testdata/sumdb/sha256test.txt")
	require.NoError(t, err)
	t.Cleanup(func() {
		err := f.Close()
		require.NoError(t, err)
	})

	remote := remotes.Git{
		RemoteTemplate: vdmspec.RemoteTemplate{
			Source:  "https://github.com/org/user",
			Version: "v1.0.0",
		},
	}
	cachedPath := filepath.Join(os.Getenv(vars.VDMHomeEnvVarName), "cache", cache.StringToBase64(remote.Source))

	err = cache.AddRemote(remote, "")
	assert.NoError(t, err)

	want := fmt.Sprintf("%s %s %s", remote.Source, remote.Version, "25fce0ea957324f2fdab37fa2c35df8dc1c62703b1970a744a723603407b630b")
	got, err := os.ReadFile(cachedPath)
	require.NoError(t, err)

	assert.Equal(t, want, string(got))
}
