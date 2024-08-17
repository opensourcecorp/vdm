package remotes

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/vdm/internal/archive"
	"github.com/opensourcecorp/vdm/internal/archive/cache"
	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
)

// File defines the file remote type
type File struct {
	vdmspec.RemoteTemplate
}

// Cache provides the [vdmspec.Remoter.Cache] operations for "file" remote types.
func (remote File) Cache() (cachePath string, err error) {
	tmpCachePath := cache.GetTempCachePath(remote)
	message.Debugf("tmpCachePath: %q", tmpCachePath)
	defer func() {
		if rmErr := os.RemoveAll(filepath.Dir(tmpCachePath)); rmErr != nil {
			err = errors.Join(err, fmt.Errorf("removing temporary cache path %q: %w", tmpCachePath, rmErr))
		}
	}()

	err = ensureParentDirs(tmpCachePath)
	if err != nil {
		return "", fmt.Errorf("creating parent temp cache directories for file remote %q: %w", tmpCachePath, err)
	}

	remote.OpMsg("Retrieving...")
	err = retrieveFile(remote, tmpCachePath)
	if err != nil {
		return "", fmt.Errorf("retrieving file: %w", err)
	}

	cachePath, err = cache.AddRemote(remote, tmpCachePath)
	if err != nil {
		return "", fmt.Errorf("caching file remote %q: %w", remote.Source, err)
	}

	remote.OpMsg("Done.")
	return cachePath, err
}

// Sync provides the [vdmspec.Remoter.Sync] operations for "file" remote types.
func (remote File) Sync(src, dest string) error {
	// We want to make sure the parent directories exist for the real destination, not the temp cache
	err := ensureParentDirs(remote.Destination)
	if err != nil {
		return fmt.Errorf("creating parent directories for file %q: %w", dest, err)
	}

	err = archive.ExtractArchiveToDestination(src, dest)
	if err != nil {
		return fmt.Errorf("syncing file cache for remote %q: %w", remote.GetSource(), err)
	}
	return nil
}

// GetSource returns the Source field.
func (remote File) GetSource() string {
	return remote.Source
}

// GetVersion returns the Version field.
func (remote File) GetVersion() string {
	return remote.Version
}

func checkFileExists(remote File) (bool, error) {
	fullPath, err := filepath.Abs(remote.Destination)
	if err != nil {
		return false, fmt.Errorf("determining abspath for file %q: %w", remote.Destination, err)
	}

	_, err = os.Stat(remote.Destination)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("couldn't check if %q exists at %q: %w", remote.Destination, fullPath, err)
	}

	return true, nil
}

func retrieveFile(remote File, dest string) (err error) {
	resp, err := http.Get(remote.Source)
	if err != nil {
		return fmt.Errorf("retrieving remote file %q: %w", remote.Source, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing response body after remote file %q retrieval: %w", remote.Source, err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccessful status code '%d' from server when retrieving remote file %q", resp.StatusCode, remote.Source)
	}

	// Note: I would normally use os.WriteFile() using the returned bytes
	// directly, but the internet says this os.Create()/io.Copy() approach
	// appears to be idiomatic
	message.Debugf("landing file to be created at %q", dest)
	outFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("creating landing file %q for remote file: %w", dest, err)
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing local file %q after remote file %q retrieval: %w", dest, remote.Source, err))
		}
	}()

	bytesWritten, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("copying HTTP response to disk: ")
	}
	message.Debugf("wrote %d bytes to %q", bytesWritten, dest)

	return nil
}

func ensureParentDirs(path string) error {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("determining abspath for file %q: %w", path, err)
	}
	message.Debugf("absolute filepath for %q determined to be %q", path, fullPath)
	dir := filepath.Dir(fullPath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("making directories: %w", err)
	}
	message.Debugf("created director(ies): %q", dir)

	return nil
}

var _ vdmspec.Remoter = File{}
