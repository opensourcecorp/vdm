package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/opensourcecorp/vdm/internal/remotes"
	"github.com/opensourcecorp/vdm/internal/vdminit"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/spf13/cobra"
)

// syncFlags defines the CLI flags for the sync subcommand.
type syncFlags struct {
	TryLocalSources bool
}

// syncFlagValues contains an initalized [syncFlags] struct with populated
// values.
var syncFlagValues syncFlags

func newSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync remotes based on specfile",
		RunE:  executeSyncSubCommand,
	}

	return cmd
}

func executeSyncSubCommand(_ *cobra.Command, _ []string) error {
	maybeSetDebug()

	err := vdminit.Paths()
	if err != nil {
		return fmt.Errorf("initializing vdm: %w", err)
	}

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
		var determinedRemote vdmspec.Remoter
		switch remote.Type {
		case vdmspec.GitType, "":
			determinedRemote = remotes.Git{RemoteTemplate: remote}
		case vdmspec.ArchiveType:
			return errors.New("cannot process 'archive' remote types, as they are not yet fully implemented")
		case vdmspec.LocalType:
			return errors.New("cannot process 'local' remote types, as they are not yet fully implemented")
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

		remote.OpMsg("Done.")
	}

	message.Infof("All done!")
	return nil
}
