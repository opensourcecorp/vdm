package cmd

import (
	"errors"
	"fmt"

	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// !!! DO NOT TOUCH, the version-bumper script handles updating this !!!
const vdmVersion string = "v0.3.0"

// rootFlags defines the CLI flags for the root command.
type rootFlags struct {
	SpecFilePath string
	Debug        bool
}

// rootFlagValues contains an initalized [rootFlags] struct with populated
// values.
var rootFlagValues rootFlags

// Flag name keys
const (
	specFilePathFlagKey string = "specfile-path"
	debugFlagKey        string = "debug"
)

func newRootCommand() *cobra.Command {
	var err error
	cmd := &cobra.Command{
		Use:              "vdm",
		Short:            "vdm -- a Versioned-Dependency Manager",
		Long:             "vdm is used to manage retrieval of arbitrary remote dependencies",
		TraverseChildren: true,
		Version:          vdmVersion,
		SilenceUsage:     true,
		SilenceErrors:    true,
		RunE:             executeRootCommand,
	}

	cmd.PersistentFlags().StringVar(&rootFlagValues.SpecFilePath, specFilePathFlagKey, "./vdm.yaml", "Path to vdm specfile")
	err = viper.BindPFlag(specFilePathFlagKey, cmd.PersistentFlags().Lookup(specFilePathFlagKey))
	if err != nil {
		message.Fatalf("internal error: unable to bind state of flag --%s: %v", specFilePathFlagKey, err)
	}

	cmd.PersistentFlags().BoolVar(&rootFlagValues.Debug, debugFlagKey, false, "Show debug messages during runtime")
	err = viper.BindPFlag(debugFlagKey, cmd.PersistentFlags().Lookup(debugFlagKey))
	if err != nil {
		message.Fatalf("internal error: unable to bind state of flag --%s: %v", debugFlagKey, err)
	}

	cmd.AddCommand(newSyncCommand())

	return cmd
}

func executeRootCommand(cmd *cobra.Command, args []string) error {
	maybeSetDebug()
	if len(args) == 0 {
		err := cmd.Help()
		if err != nil {
			return errors.New("failed to print help message, somehow")
		}
	}

	return errors.New("You must provide a subcommand to vdm")
}

// Execute wraps the primary execution logic for vdm's root command, and returns
// any errors encountered to the caller.
func Execute() error {
	cmd := newRootCommand()
	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("executing root command: %w", err)
	}

	return nil
}
