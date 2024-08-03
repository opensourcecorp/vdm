package cmd

import (
	"fmt"

	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/opensourcecorp/vdm/internal/remotes"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync remotes based on specfile",
	RunE:  syncExecute,
}

// syncFlags defines the CLI flags for the sync subcommand.
type syncFlags struct {
	TryLocalSources bool
}

// syncFlagValues contains an initalized [syncFlags] struct with populated
// values.
var syncFlagValues syncFlags

// Flag name keys
const (
	tryLocalSourcesFlagKey string = "try-local-sources"
)

func init() {
	var err error

	syncCmd.Flags().BoolVar(&syncFlagValues.TryLocalSources, tryLocalSourcesFlagKey, false, "Whether to try & process local copies of sources before retrieving their remote copies")
	err = viper.BindPFlag(tryLocalSourcesFlagKey, syncCmd.Flags().Lookup(tryLocalSourcesFlagKey))
	if err != nil {
		message.Fatalf("internal error: unable to bind state of flag --%s: %v", tryLocalSourcesFlagKey, err)
	}
}

func syncExecute(_ *cobra.Command, _ []string) error {
	maybeSetDebug()
	maybeTryLocalSources()
	if err := sync(); err != nil {
		return fmt.Errorf("executing sync command: %w", err)
	}
	return nil
}

// sync does the heavy lifting to ensure that the local directory tree(s) match
// the desired state as defined in the specfile.
func sync() error {
	spec, err := vdmspec.GetSpecFromFile(rootFlagValues.SpecFilePath)
	if err != nil {
		return fmt.Errorf("getting specs from spec file: %w", err)
	}

	err = spec.Validate()
	if err != nil {
		return fmt.Errorf("your vdm spec file is malformed: %w", err)
	}

	for _, remote := range spec.Remotes {
		// process stored vdm metafile so we know what operations to actually
		// perform for existing directories
		vdmMeta, err := remote.GetVDMMeta()
		if err != nil {
			return fmt.Errorf("getting vdm metadata file for sync: %w", err)
		}

		if vdmMeta == (vdmspec.Remote{}) {
			message.Infof("%s: %s not found at local path, will be created", remote.OpMsg(), vdmspec.MetaFileName)
		} else {
			if vdmMeta.Version != remote.Version && vdmMeta.Source != remote.Source {
				message.Infof("%s: Will change '%s' from current local version spec '%s' to '%s'...", remote.OpMsg(), remote.Source, vdmMeta.Version, remote.Version)
				panic("not implemented")
			}
			message.Infof("%s: version unchanged in spec file, skipping", remote.OpMsg())
			continue
		}

		switch remote.Type {
		case vdmspec.GitType, "":
			if err := remotes.SyncGit(remote); err != nil {
				return fmt.Errorf("syncing git remote: %w", err)
			}
		case vdmspec.FileType:
			if err := remotes.SyncFile(remote); err != nil {
				return fmt.Errorf("syncing file remote: %w", err)
			}
		default:
			return fmt.Errorf("unrecognized remote type '%s'", remote.Type)
		}

		err = remote.WriteVDMMeta()
		if err != nil {
			return fmt.Errorf("could not write %s file to disk: %w", vdmspec.MetaFileName, err)
		}

		message.Infof("%s: Done.", remote.OpMsg())
	}

	message.Infof("All done!")
	return nil
}
