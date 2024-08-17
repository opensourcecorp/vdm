package cache

import (
	"errors"
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

	remoteWithVersion := fmt.Sprintf("%s@%s", remote.GetSource(), remote.GetVersion())
	b64 := StringToBase64(remoteWithVersion)
	if err != nil {
		return "", fmt.Errorf("calculating sum when caching remote %q: %w", remoteWithVersion, err)
	}
	message.Debugf("remote %q base64'd to %q", remoteWithVersion, b64)

	cacheFileName := b64 + ".tar.gz"
	cacheTargetPath := filepath.Join(vars.GetVDMCacheDir(), cacheFileName)

	cacheTarget, err := archive.CreateArchive(cacheRoot, cacheTargetPath)
	if err != nil {
		return "", fmt.Errorf("creating archive while adding remote %q: %w", remoteWithVersion, err)
	}

	sumDBFile, err := GetOrCreateSumDBFile()
	if err != nil {
		return "", fmt.Errorf("creating/opening sumdb file %q: %w", sumDBFile.Name(), err)
	}
	defer func() {
		if closeErr := sumDBFile.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing sumdb file %q: %w", sumDBFile.Name(), closeErr))
		}
	}()

	sum, err := CalculateSHASum(cacheTarget)
	if err != nil {
		return "", fmt.Errorf("calculating checksum for writing: %w", err)
	}
	message.Debugf("sum calculated for path %q was %q", cacheTargetPath, sum)

	sumDBContents := fmt.Sprintf("%s %s %s\n", remote.GetSource(), remote.GetVersion(), sum)
	_, err = fmt.Fprint(sumDBFile, sumDBContents)
	if err != nil {
		return "", fmt.Errorf("writing to cache file %q: %w", cacheTargetPath, err)
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
