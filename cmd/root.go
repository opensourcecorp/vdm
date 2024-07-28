package cmd

import (
	"fmt"

	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = cobra.Command{
	Use:              "vdm",
	Short:            "vdm -- a Versioned-Dependency Manager",
	Long:             "vdm is used to manage arbitrary remote dependencies",
	TraverseChildren: true,
	Run: func(_ *cobra.Command, _ []string) {
		MaybeSetDebug()
	},
}

type rootFlags struct {
	SpecFilePath string
	Debug        bool
}

// RootFlagValues contains an initalized [rootFlags] struct with populated
// values.
var RootFlagValues rootFlags

// Flag name keys
const (
	specFilePathFlagKey string = "specfile-path"
	debugFlagKey        string = "debug"
)

func init() {
	var err error

	rootCmd.PersistentFlags().StringVar(&RootFlagValues.SpecFilePath, specFilePathFlagKey, "./vdm.yaml", "Path to vdm specfile")
	err = viper.BindPFlag(specFilePathFlagKey, rootCmd.PersistentFlags().Lookup(specFilePathFlagKey))
	if err != nil {
		message.Fatalf("internal error: unable to bind state of flag --%s", specFilePathFlagKey)
	}

	rootCmd.PersistentFlags().BoolVar(&RootFlagValues.Debug, debugFlagKey, false, "Show debug messages during runtime")
	err = viper.BindPFlag(debugFlagKey, rootCmd.PersistentFlags().Lookup(debugFlagKey))
	if err != nil {
		message.Fatalf("internal error: unable to bind state of flag --%s", debugFlagKey)
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
