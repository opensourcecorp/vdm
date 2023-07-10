package remote

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/sirupsen/logrus"
)

// SyncFile is the root of the sync operations for "file" remote types.
func SyncFile(spec vdmspec.Spec) error {
	fileExists, err := checkFileExists(spec)
	if err != nil {
		return fmt.Errorf("checking if file exists locally: %w", err)
	}

	if !fileExists {
		logrus.Infof("File '%s' does not exist locally, retrieving", spec.LocalPath)
		err = retrieveFile(spec)
		if err != nil {
			return fmt.Errorf("retrieving file: %w", err)
		}
	} else {
		logrus.Infof("File '%s' already exists locally, skipping", spec.LocalPath)
	}

	return nil
}

func checkFileExists(spec vdmspec.Spec) (bool, error) {
	fullPath, err := filepath.Abs(spec.LocalPath)
	if err != nil {
		return false, fmt.Errorf("determining abspath for file '%s': %w", spec.LocalPath, err)
	}

	_, err = os.Stat(spec.LocalPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("couldn't check if %s exists at '%s': %w", spec.LocalPath, fullPath, err)
	}

	return true, nil
}

func retrieveFile(spec vdmspec.Spec) error {
	resp, err := http.Get(spec.Remote)
	if err != nil {
		return fmt.Errorf("retrieving remote file '%s': %w", spec.Remote, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logrus.Errorf("closing response body after remote file '%s' retrieval: %v", spec.Remote, err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccessful status code '%d' from server when retrieving remote file '%s'", resp.StatusCode, spec.Remote)
	}

	// Note: I would normally use os.WriteFile() using the returned bytes
	// directly, but the internet says this os.Create()/io.Copy() approach
	// appears to be idiomatic
	outFile, err := os.Create(spec.LocalPath)
	if err != nil {
		return fmt.Errorf("creating landing file '%s' for remote file: %w", spec.LocalPath, err)
	}
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil {
			logrus.Errorf("closing local file '%s' after remote file '%s' retrieval: %v", spec.LocalPath, spec.Remote, err)
		}
	}()

	bytesWritten, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("copying HTTP response to disk: ")
	}
	logrus.Debugf("wrote %d bytes to '%s'", bytesWritten, spec.LocalPath)

	return nil
}
