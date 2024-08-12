// Package vars houses constants etc. for working with command-line flag values
// across packages. These helpers are pushed down to their own package in order
// to avoid import cycles.
package vars

import (
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

var (
	// VDMHome stores the location of vdm's own home directory
	VDMHome string
)

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic("unable to determine home directory")
	}

	vdmHomeOverride, ok := os.LookupEnv(VDMHomeEnvVarName)
	if ok {
		VDMHome = filepath.Join(vdmHomeOverride, ".vdm")
	} else {
		VDMHome = filepath.Join(homedir, ".vdm")
	}
	err = os.Setenv(VDMHomeEnvVarName, VDMHome)
	if err != nil {
		panic("unable to set VDM home directory")
	}
}

// GetVDMCacheDir returns the determined path to vdm's cache directory.
func GetVDMCacheDir() string {
	return filepath.Join(os.Getenv(VDMHomeEnvVarName), "cache")
}
