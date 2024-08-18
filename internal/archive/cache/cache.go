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
// handle for a target archive to actually write the archive data. TODO fix this
func AddRemote(remote vdmspec.Remoter, cacheRoot string) (cachePath string, err error) {
	err = os.MkdirAll(vars.GetVDMCacheDir(), 0755)
	if err != nil {
		return "", fmt.Errorf("creating vdm cache directory %q: %w", vars.GetVDMCacheDir(), err)
	}

	b64 := StringToBase64(remote.GetSourceVersion())
	if err != nil {
		return "", fmt.Errorf("calculating sum when caching remote %q: %w", remote.GetSourceVersion(), err)
	}
	message.Debugf("remote %q base64'd to %q", remote.GetSourceVersion(), b64)

	cacheFileName := b64 + ".tar.gz"
	cacheTargetPath := filepath.Join(vars.GetVDMCacheDir(), cacheFileName)

	cacheTarget, err := archive.CreateArchive(cacheRoot, cacheTargetPath)
	if err != nil {
		return "", fmt.Errorf("creating archive while adding remote %q: %w", remote.GetSourceVersion(), err)
	}
	defer cacheTarget.Close()

	err = AddToSumDB(remote, cacheTarget)
	if err != nil {
		return "", fmt.Errorf("adding %q details to sumdb: %w", cacheTargetPath, err)
	}

	return cacheTargetPath, err
}

func GetTempCachePath(remote vdmspec.Remoter) string {
	// tmpCachePath is where the actual retrieval is targeted, which should then
	// be later archived to the persistent cache
	randID := 100000000000 + rand.Intn(999999999999)
	tmpCacheRoot := filepath.Join(os.TempDir(), fmt.Sprintf("vdm-tmp-%d", randID))
	tmpCachePath := filepath.Join(tmpCacheRoot, filepath.Base(remote.GetSource()))

	return tmpCachePath
}
