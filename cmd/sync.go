package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/opensourcecorp/vdm/internal/remotes"
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
	spec, err := vdmspec.GetSpecFromFile(RootFlagValues.SpecFilePath)
	if err != nil {
		return fmt.Errorf("getting specs from spec file: %w", err)
	}

	err = spec.Validate()
	if err != nil {
		return fmt.Errorf("your vdm spec file is malformed: %w", err)
	}

SpecLoop:
	for _, remote := range spec.Remotes {
		// process stored vdm metafile so we know what operations to actually
		// perform for existing directories
		vdmMeta, err := remote.GetVDMMeta()
		if err != nil {
			return fmt.Errorf("getting vdm metadata file for sync: %w", err)
		}

		if vdmMeta == (vdmspec.Remote{}) {
			logrus.Infof("%s not found at local path '%s' -- will be created", vdmspec.MetaFileName, filepath.Join(remote.LocalPath))
		} else {
			if vdmMeta.Version != remote.Version && vdmMeta.Remote != remote.Remote {
				logrus.Infof("Will change '%s' from current local version spec '%s' to '%s'...", remote.Remote, vdmMeta.Version, remote.Version)
				panic("not implemented")
			} else {
				logrus.Infof("Version unchanged (%s) in spec file for '%s' --> '%s', skipping", remote.Version, remote.Remote, remote.LocalPath)
				continue SpecLoop
			}
		}

		switch remote.Type {
		case "git", "":
			remotes.SyncGit(remote)
		case "file":
			remotes.SyncFile(remote)
		default:
			return fmt.Errorf("unrecognized remote type '%s'", remote.Type)
		}

		err = remote.WriteVDMMeta()
		if err != nil {
			return fmt.Errorf("could not write %s file to disk: %w", vdmspec.MetaFileName, err)
		}

		logrus.Infof("%s -- Done.", remote.OpMsg())
	}

	logrus.Info("All done!")
	return nil
}
