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

func syncExecute(_ *cobra.Command, _ []string) error {
	if err := sync(); err != nil {
		return fmt.Errorf("executing sync command: %w", err)
	}
	return nil
}

// sync does the heavy lifting to ensure that the local directory tree(s) match
// the desired state as defined in the specfile.
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

SpecLoop:
	for _, spec := range specs {
		// process stored vdm metafile so we know what operations to actually
		// perform for existing directories
		vdmMeta, err := spec.GetVDMMeta()
		if err != nil {
			return fmt.Errorf("getting vdm metadata file for sync: %w", err)
		}

		if vdmMeta == (vdmspec.Spec{}) {
			logrus.Infof("%s not found at local path '%s' -- will be created", vdmspec.MetaFileName, filepath.Join(spec.LocalPath))
		} else {
			if vdmMeta.Version != spec.Version && vdmMeta.Remote != spec.Remote {
				logrus.Infof("Will change '%s' from current local version spec '%s' to '%s'...", spec.Remote, vdmMeta.Version, spec.Version)
				panic("not implemented")
			} else {
				logrus.Infof("Version unchanged (%s) in spec file for '%s' --> '%s', skipping", spec.Version, spec.Remote, spec.LocalPath)
				continue SpecLoop
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
