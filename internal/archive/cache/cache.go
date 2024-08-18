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
	b64 := StringToBase64(remote.GetSourceVersion())
	message.Debugf("remote %q base64'd to %q", remote.GetSourceVersion(), b64)

	cacheFileName := b64 + ".tar.gz"
	cacheTargetPath := filepath.Join(vars.GetVDMCacheDir(), cacheFileName)

	remoteAlreadyInSumDB, err := CheckIfRemoteInSumDB(remote)
	if err != nil {
		return "", fmt.Errorf("checking if remote %q already in sumdb: %w", remote.GetSourceVersion(), err)
	}

	if !remoteAlreadyInSumDB {
		message.Infof("Remote %q already found in local cache, will retrieve from there")
		err = archive.CreateArchive(cacheRoot, cacheTargetPath)
		if err != nil {
			return "", fmt.Errorf("creating archive while adding remote %q: %w", remote.GetSourceVersion(), err)
		}

		err = AddToSumDB(remote, cacheTargetPath)
		if err != nil {
			return "", fmt.Errorf("adding %q details to sumdb: %w", cacheTargetPath, err)
		}
	}

	return cacheTargetPath, nil
}

func GetTempCachePath(remote vdmspec.Remoter) string {
	// tmpCachePath is where the actual retrieval is targeted, which should then
	// be later archived to the persistent cache
	randID := 100000000000 + rand.Intn(999999999999)
	tmpCacheRoot := filepath.Join(os.TempDir(), fmt.Sprintf("vdm-tmp-%d", randID))
	tmpCachePath := filepath.Join(tmpCacheRoot, filepath.Base(remote.GetSource()))

	return tmpCachePath
}
