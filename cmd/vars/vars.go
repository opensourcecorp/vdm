// Package vars houses constants etc. for working with command-line flag values
// across packages. These helpers are pushed down to their own package in order
// to avoid import cycles.
package vars

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// Debug anchors to the DEBUG env var
	Debug = "DEBUG"
	// TryLocalSources anchors to the TRY_LOCAL_SOURCES env var
	TryLocalSources = "TRY_LOCAL_SOURCES"
	// VDMHomeEnvVarName anchors to the VDM_HOME env var
	VDMHomeEnvVarName = "VDM_HOME"
)

func GetVDMHomeDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("determining home directory: %w", err)
	}

	var vdmHome string
	vdmHomeOverride, ok := os.LookupEnv(VDMHomeEnvVarName)
	if ok {
		vdmHome = filepath.Join(vdmHomeOverride, ".vdm")
	} else {
		vdmHome = filepath.Join(homedir, ".vdm")
	}

	return vdmHome, nil
}

// GetVDMCacheDir returns the determined path to vdm's cache directory.
func GetVDMCacheDir() (string, error) {
	vdmHome, err := GetVDMHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting vdm home directory: %w", err)
	}

	return filepath.Join(vdmHome, "cache"), nil
}
