package cache

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/vdm/cmd/vars"
	"github.com/opensourcecorp/vdm/internal/archive"
	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
)

// AddRemote uses the provided [vdmspec.Remoter] information along with a file
// handle for a target archive to actually write the archive data.
func AddRemote(remote vdmspec.Remoter, cacheRoot string) (string, error) {
	cachePath, err := GetPersistentCacheFilePath(remote)
	if err != nil {
		return "", fmt.Errorf("getting cache location for remote %q: %w", remote.GetSumDBKey(), err)
	}
	message.Debugf("persistent cache file path: %q", cachePath)

	err = archive.CreateArchive(cacheRoot, cachePath)
	if err != nil {
		return "", fmt.Errorf("creating archive while adding remote %q: %w", remote.GetSumDBKey(), err)
	}

	err = AddToSumDB(remote, cachePath)
	if err != nil {
		return "", fmt.Errorf("adding %q details to sumdb: %w", cachePath, err)
	}

	return cachePath, nil
}

func GetPersistentCacheFilePath(remote vdmspec.Remoter) (string, error) {
	b64 := StringToBase64(remote.GetSumDBKey())
	message.Debugf("remote %q base64'd to %q", remote.GetSumDBKey(), b64)

	cacheFileName := b64 + ".tar.gz"

	vdmCacheDir, err := vars.GetVDMCacheDir()
	if err != nil {
		return "", fmt.Errorf("determining vdm cache directory while adding remote: %w", err)
	}

	cachePath := filepath.Join(vdmCacheDir, cacheFileName)
	return cachePath, nil
}

func GetTempCachePath(remote vdmspec.Remoter) string {
	// tmpCachePath is where the actual retrieval is targeted, which should then
	// be later archived to the persistent cache
	randID := 100000000000 + rand.Intn(999999999999)
	tmpCacheRoot := filepath.Join(os.TempDir(), fmt.Sprintf("vdm-tmp-%d", randID))
	tmpCachePath := filepath.Join(tmpCacheRoot, filepath.Base(remote.GetSource()))

	return tmpCachePath
}
