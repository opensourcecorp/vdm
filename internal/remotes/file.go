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

// SyncFile is the root of the sync operations for "file" remote types.
func SyncFile(remote vdmspec.Remote) error {
	fileExists, err := checkFileExists(remote)
	if err != nil {
		return fmt.Errorf("checking if file exists locally: %w", err)
	}

	if !fileExists {
		message.Infof("File '%s' does not exist locally, retrieving", remote.LocalPath)
		err = retrieveFile(remote)
		if err != nil {
			return fmt.Errorf("retrieving file: %w", err)
		}
	} else {
		message.Infof("File '%s' already exists locally, skipping", remote.LocalPath)
	}

	return nil
}

func checkFileExists(remote vdmspec.Remote) (bool, error) {
	fullPath, err := filepath.Abs(remote.LocalPath)
	if err != nil {
		return false, fmt.Errorf("determining abspath for file '%s': %w", remote.LocalPath, err)
	}

	_, err = os.Stat(remote.LocalPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("couldn't check if %s exists at '%s': %w", remote.LocalPath, fullPath, err)
	}

	return true, nil
}

func retrieveFile(remote vdmspec.Remote) error {
	resp, err := http.Get(remote.Remote)
	if err != nil {
		return fmt.Errorf("retrieving remote file '%s': %w", remote.Remote, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			message.Errorf("closing response body after remote file '%s' retrieval: %v", remote.Remote, err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccessful status code '%d' from server when retrieving remote file '%s'", resp.StatusCode, remote.Remote)
	}

	// Note: I would normally use os.WriteFile() using the returned bytes
	// directly, but the internet says this os.Create()/io.Copy() approach
	// appears to be idiomatic
	outFile, err := os.Create(remote.LocalPath)
	if err != nil {
		return fmt.Errorf("creating landing file '%s' for remote file: %w", remote.LocalPath, err)
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil {
			message.Errorf("closing local file '%s' after remote file '%s' retrieval: %v", remote.LocalPath, remote.Remote, err)
		}
	}()

	bytesWritten, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("copying HTTP response to disk: ")
	}
	message.Debugf("wrote %d bytes to '%s'", bytesWritten, remote.LocalPath)

	return nil
}
