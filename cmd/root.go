package cmd

import (
	"fmt"
	"os"

	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// !!! DO NOT TOUCH, the version-bumper script handles updating this !!!
const vdmVersion string = "v0.3.0"

var rootCmd = cobra.Command{
	Use:              "vdm",
	Short:            "vdm -- a Versioned-Dependency Manager",
	Long:             "vdm is used to manage retrieval of arbitrary remote dependencies",
	TraverseChildren: true,
	Version:          vdmVersion,
	Run: func(cmd *cobra.Command, args []string) {
		maybeSetDebug()
		if len(args) == 0 {
			message.Errorf("You must provide a subcommand to vdm")
			err := cmd.Help()
			if err != nil {
				message.Fatalf("failed to print help message, somehow")
			}
			os.Exit(1)
		}
	},
}

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

func init() {
	var err error

	rootCmd.PersistentFlags().StringVar(&rootFlagValues.SpecFilePath, specFilePathFlagKey, "./vdm.yaml", "Path to vdm specfile")
	err = viper.BindPFlag(specFilePathFlagKey, rootCmd.PersistentFlags().Lookup(specFilePathFlagKey))
	if err != nil {
		message.Fatalf("internal error: unable to bind state of flag --%s: %v", specFilePathFlagKey, err)
	}

	rootCmd.PersistentFlags().BoolVar(&rootFlagValues.Debug, debugFlagKey, false, "Show debug messages during runtime")
	err = viper.BindPFlag(debugFlagKey, rootCmd.PersistentFlags().Lookup(debugFlagKey))
	if err != nil {
		message.Fatalf("internal error: unable to bind state of flag --%s: %v", debugFlagKey, err)
	}

	rootCmd.AddCommand(syncCmd)
}

// Execute wraps the primary execution logic for vdm's root command, and returns
// any errors encountered to the caller.
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("executing root command: %w", err)
	}

	return nil
}
