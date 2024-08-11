package sumdb

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/vdm/cmd/vars"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
)

// CalculateSHASum takes an arbitrary [io.Reader] (such as an open file handle)
// and calculates the SHA256 checksum for it.
func CalculateSHASum(reader io.Reader) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, reader); err != nil {
		return "", fmt.Errorf("writing reader to hasher: %w", err)
	}
	sum := fmt.Sprintf("%x", hasher.Sum(nil))
	return sum, nil
}

func CacheRemote(remote vdmspec.Remote, archiveFileHandle io.Reader) (err error) {
	err = os.MkdirAll(vars.GetVDMCacheDir(), 0755)
	if err != nil {
		return fmt.Errorf("creating vdm cache directory %s: %w", vars.GetVDMCacheDir(), err)
	}

	cacheFileName := StringAsBase64(remote.Source)
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
		return fmt.Errorf("calculating sum when caching remote %s: %w", remote.Source, err)
	}

	contents := fmt.Sprintf("%s %s %s", remote.Source, remote.Version, sum)

	err = os.WriteFile(f.Name(), []byte(contents), 0644)
	if err != nil {
		return fmt.Errorf("writing cache file %s: %w", f.Name(), err)
	}

	return err
}

func StringAsBase64(s string) string {
	// We use base64-encoded names for caches, because remote sources have
	// slashes and that can break FS pathing
	data := []byte(s)
	encodedBytes := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encodedBytes, data)
	return string(encodedBytes)
}
