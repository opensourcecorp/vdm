package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:              "vdm",
	Short:            "vdm -- a Versioned-Dependency Manager",
	Long:             "vdm is used to manage arbitrary remote dependencies",
	TraverseChildren: true,
}

type RootFlags struct {
	SpecFilePath string
	Debug        bool
}

var RootFlagValues RootFlags

func init() {
	rootCmd.PersistentFlags().StringVar(&RootFlagValues.SpecFilePath, "specfile-path", "./vdm.yaml", "Path to vdm specfile")
	rootCmd.PersistentFlags().BoolVar(&RootFlagValues.Debug, "debug", false, "Show debug logs")

	if RootFlagValues.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	rootCmd.AddCommand(syncCmd)
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("executing root command: %v", err)
	}

	return nil
}
