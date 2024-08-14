package cache

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/vdm/cmd/vars"
	"github.com/opensourcecorp/vdm/internal/archive"
	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
)

// AddRemote uses the provided [vdmspec.Remoter] information along with a file
// handle for a target archive to actually write the archive data. TODO fix this
func AddRemote(remote vdmspec.Remoter, cacheRoot string) (err error) {
	err = os.MkdirAll(vars.GetVDMCacheDir(), 0755)
	if err != nil {
		return fmt.Errorf("creating vdm cache directory %s: %w", vars.GetVDMCacheDir(), err)
	}

	b64 := StringToBase64(remote.GetSource())
	if err != nil {
		return fmt.Errorf("calculating sum when caching remote %s: %w", remote.GetSource(), err)
	}
	message.Debugf("remote '%s' base64'd to '%s'", remote.GetSource(), b64)

	cacheFileName := b64 + ".tar.gz"
	cacheTargetPath := filepath.Join(vars.GetVDMCacheDir(), cacheFileName)

	cacheTarget, err := archive.CreateArchive(cacheRoot, cacheTargetPath)
	if err != nil {
		return fmt.Errorf("creating archive while adding remote '%s': %w", remote.GetSource(), err)
	}

	sumDBFile, err := GetOrCreateSumDBFile()
	if err != nil {
		return fmt.Errorf("creating/opening sumdb file '%s': %w", sumDBFile.Name(), err)
	}
	defer func() {
		if closeErr := sumDBFile.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing sumdb file '%s': %w", sumDBFile.Name(), closeErr))
		}
	}()

	sum, err := CalculateSHASum(cacheTarget)
	if err != nil {
		return fmt.Errorf("calculating checksum for writing: %w", err)
	}
	message.Debugf("sum calculated for path '%s' was '%s'", cacheTargetPath)

	sumDBContents := fmt.Sprintf("%s %s %s\n", remote.GetSource(), remote.GetVersion(), sum)
	_, err = fmt.Fprint(sumDBFile, sumDBContents)
	if err != nil {
		return fmt.Errorf("writing to cache file '%s': %w", cacheTargetPath, err)
	}

	return err
}
