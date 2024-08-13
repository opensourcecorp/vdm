package cache

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/vdm/cmd/vars"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
)

// AddRemote uses the provided [vdmspec.Remoter] information along with a file
// handle for a target archive to actually write the archive data.
func AddRemote(remote vdmspec.Remoter, archiveFileHandle io.Reader) (err error) {
	err = os.MkdirAll(vars.GetVDMCacheDir(), 0755)
	if err != nil {
		return fmt.Errorf("creating vdm cache directory %s: %w", vars.GetVDMCacheDir(), err)
	}

	cacheFileName := StringToBase64(remote.GetSource())
	f, err := os.Create(filepath.Join(vars.GetVDMCacheDir(), cacheFileName))
	if err != nil {
		return fmt.Errorf("creating file %s: %w", cacheFileName, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing cache file %s: %w", cacheFileName, closeErr))
		}
	}()

	sum, err := CalculateSHASum(archiveFileHandle)
	if err != nil {
		return fmt.Errorf("calculating sum when caching remote %s: %w", remote.GetSource(), err)
	}

	contents := fmt.Sprintf("%s %s %s", remote.GetSource(), remote.GetVersion(), sum)

	err = os.WriteFile(f.Name(), []byte(contents), 0644)
	if err != nil {
		return fmt.Errorf("writing cache file %s: %w", f.Name(), err)
	}

	return err
}
