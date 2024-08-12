package remotes

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
)

// File defines the file remote type
type File struct {
	vdmspec.Remote
}

// Cache provides the [vdmspec.Remoter.Cache] operations for "file" remote types.
func (remote File) Cache() error {
	return errors.New("not implemented")
}

// Sync provides the [vdmspec.Remoter.Sync] operations for "file" remote types.
func (remote File) Sync() error {
	fileExists, err := checkFileExists(remote)
	if err != nil {
		return fmt.Errorf("checking if file exists locally: %w", err)
	}

	if !fileExists {
		message.Infof("File '%s' does not exist locally, retrieving", remote.Destination)
		err = retrieveFile(remote)
		if err != nil {
			return fmt.Errorf("retrieving file: %w", err)
		}
	} else {
		message.Infof("File '%s' already exists locally, skipping", remote.Destination)
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
		return false, fmt.Errorf("determining abspath for file '%s': %w", remote.Destination, err)
	}

	_, err = os.Stat(remote.Destination)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("couldn't check if %s exists at '%s': %w", remote.Destination, fullPath, err)
	}

	return true, nil
}

func retrieveFile(remote File) (err error) {
	resp, err := http.Get(remote.Source)
	if err != nil {
		return fmt.Errorf("retrieving remote file '%s': %w", remote.Source, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing response body after remote file '%s' retrieval: %w", remote.Source, err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccessful status code '%d' from server when retrieving remote file '%s'", resp.StatusCode, remote.Source)
	}

	err = ensureParentDirs(remote.Destination)
	if err != nil {
		return fmt.Errorf("creating parent directories for file: %w", err)
	}

	// Note: I would normally use os.WriteFile() using the returned bytes
	// directly, but the internet says this os.Create()/io.Copy() approach
	// appears to be idiomatic
	outFile, err := os.Create(remote.Destination)
	if err != nil {
		return fmt.Errorf("creating landing file '%s' for remote file: %w", remote.Destination, err)
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing local file '%s' after remote file '%s' retrieval: %w", remote.Destination, remote.Source, err))
		}
	}()

	bytesWritten, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("copying HTTP response to disk: ")
	}
	message.Debugf("wrote %d bytes to '%s'", bytesWritten, remote.Destination)

	return nil
}

func ensureParentDirs(path string) error {
	fullPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("determining abspath for file '%s': %w", path, err)
	}
	message.Debugf("absolute filepath for '%s' determined to be '%s'", path, fullPath)
	dir := filepath.Dir(fullPath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("making directories: %w", err)
	}
	message.Debugf("created director(ies): %s", dir)

	return nil
}

var _ vdmspec.Remoter = File{}
