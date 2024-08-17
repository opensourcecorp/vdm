package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/opensourcecorp/vdm/internal/remotes"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

func newSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync remotes based on specfile",
		RunE:  executeSyncSubCommand,
	}

	cmd.Flags().BoolVar(&syncFlagValues.TryLocalSources, tryLocalSourcesFlagKey, false, "Whether to try & process local copies of sources before retrieving their remote copies")
	err := viper.BindPFlag(tryLocalSourcesFlagKey, cmd.Flags().Lookup(tryLocalSourcesFlagKey))
	if err != nil {
		message.Fatalf("internal error: unable to bind state of flag --%s: %v", tryLocalSourcesFlagKey, err)
	}

	return cmd
}

func executeSyncSubCommand(_ *cobra.Command, _ []string) error {
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
		// TODO: add this back, but it's unused right now
		// // process stored vdm metafile so we know what operations to actually
		// // perform for existing directories
		// vdmMeta, err := remote.GetVDMMeta()
		// if err != nil {
		// 	return fmt.Errorf("getting vdm metadata file for sync: %w", err)
		// }

		// if vdmMeta == (vdmspec.RemoteTemplate{}) {
		// 	message.Infof("%s: %s not found at local path, will be created", remote.OpMsg(), vdmspec.MetaFileName)
		// } else {
		// 	if vdmMeta.Version != remote.Version && vdmMeta.Source != remote.Source {
		// 		message.Infof("%s: Will change %q from current local version spec %q to %q...", remote.OpMsg(), remote.Source, vdmMeta.Version, remote.Version)
		// 		panic("jk not implemented")
		// 	}
		// 	message.Infof("%s: version unchanged in spec file, skipping", remote.OpMsg())
		// 	continue
		// }

		var determinedRemote vdmspec.Remoter
		switch remote.Type {
		case vdmspec.GitType, "":
			determinedRemote = remotes.Git{RemoteTemplate: remote}
		case vdmspec.FileType:
			return errors.New("cannot process 'file' remote types, as they are not yet fully implemented")
		default:
			return fmt.Errorf("unrecognized remote type %q", remote.Type)
		}

		cachePath, err := determinedRemote.Cache()
		if err != nil {
			return fmt.Errorf("caching %q remote: %w", remote.Type, err)
		}

		absDestination, err := filepath.Abs(remote.Destination)
		if err != nil {
			return fmt.Errorf("determining abspath of remote's destination %q: %w", remote.Destination, err)
		}

		err = determinedRemote.Sync(cachePath, absDestination)
		if err != nil {
			return fmt.Errorf("syncing %q remote: %w", remote.Type, err)
		}

		// TODO: add back, as above
		// err = remote.WriteVDMMeta()
		// if err != nil {
		// 	return fmt.Errorf("could not write %s file to disk: %w", vdmspec.MetaFileName, err)
		// }

		remote.OpMsg("Done.")
	}

	message.Infof("All done!")
	return nil
}
