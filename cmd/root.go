package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:              "vdm",
	Short:            "vdm -- a Versioned-Dependency Manager",
	Long:             "vdm is used to manage arbitrary remote dependencies",
	TraverseChildren: true,
}

type rootFlags struct {
	SpecFilePath string
	Debug        bool
}

// RootFlagValues contains an initalized [rootFlags] struct with populated
// values.
var RootFlagValues rootFlags

func init() {
	rootCmd.PersistentFlags().StringVar(&RootFlagValues.SpecFilePath, "specfile-path", "./vdm.yaml", "Path to vdm specfile")
	rootCmd.PersistentFlags().BoolVar(&RootFlagValues.Debug, "debug", false, "Show debug messages during runtime")

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
