package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/opensourcecorp/vdm/internal/remote"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync remotes based on specfile",
	RunE:  syncExecute,
}

type SyncFlags struct {
	KeepGitDir bool
}

var SyncFlagValues SyncFlags

func syncExecute(_ *cobra.Command, _ []string) error {
	if err := sync(); err != nil {
		return fmt.Errorf("executing sync command: %w", err)
	}
	return nil
}

// sync ensures that the only local dependencies are ones defined in the specfile
func sync() error {
	specs, err := vdmspec.GetSpecsFromFile(RootFlagValues.SpecFilePath)
	if err != nil {
		return fmt.Errorf("getting specs from spec file: %w", err)
	}

	for _, spec := range specs {
		err := spec.Validate()
		if err != nil {
			return fmt.Errorf("your vdm spec file is malformed: %w", err)
		}
	}

	for _, spec := range specs {
		// process stored vdm metafile so we know what operations to actually
		// perform for existing directories
		vdmMeta, err := spec.GetVDMMeta()
		if err != nil {
			return fmt.Errorf("getting vdm metadata file for sync: %w", err)
		}

		if vdmMeta == (vdmspec.VDMSpec{}) {
			logrus.Infof("%s not found at local path '%s' -- will be created", vdmspec.MetaFileName, filepath.Join(spec.LocalPath))
		} else {
			if vdmMeta.Version != spec.Version {
				logrus.Infof("Changing '%s' from current local version spec '%s' to '%s'...", spec.Remote, vdmMeta.Version, spec.Version)
			} else {
				logrus.Debugf("Version unchanged (%s) in spec file for '%s' --> '%s'", spec.Version, spec.Remote, spec.LocalPath)
			}
		}

		switch spec.Type {
		case "git", "":
			remote.SyncGit(spec)
		case "file":
			remote.SyncFile(spec)
		default:
			return fmt.Errorf("unrecognized remote type '%s'", spec.Type)
		}

		err = spec.WriteVDMMeta()
		if err != nil {
			return fmt.Errorf("could not write %s file to disk: %w", vdmspec.MetaFileName, err)
		}

		logrus.Infof("%s -- Done.", spec.OpMsg())
	}

	logrus.Info("All done!")
	return nil
}
