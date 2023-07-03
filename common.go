// Common and/or initialization consts, vars, and functions
package main

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
)

var (
	// Subcommands
	syncCmd = flag.NewFlagSet("sync", flag.ExitOnError)

	subcommands = map[string]*flag.FlagSet{
		syncCmd.Name(): syncCmd,
	}

	// CLI args common to each subcommand
	debug        bool
	specFilePath string

	// sync CLI flags
	keepGitDir bool
)

// registerFlags assigns values to flags that should belong to each and/or all
// command(s)
func registerFlags() {
	// common
	for _, cmd := range subcommands {
		cmd.StringVar(&specFilePath, "spec-file", "./vdm.json", "Path to vdm spec file")
		cmd.BoolVar(&debug, "debug", false, "Print debug logs")
	}

	// sync
	syncCmd.BoolVar(&keepGitDir, "keep-git-dir", false, "should vdm keep the .git directory within git-sourced directories? Most useful if you're using vdm to initialize groups of actual repositories you intend to work in")
}

// rootUsage has help text for the root command, so that users don't get an
// unhelpful error when forgetting to specify a subcommand
func showRootUsage() {
	fmt.Printf(`vdm declaratively manages remote dependencies as local directories.

Subcommands:
	sync	sync local paths based on your vdm spec file

`)
}

// checkRootUsage prints usage information if a user doesn't specify a subcommand
func checkRootUsage(args []string) {
	helpFlagRegex := regexp.MustCompile(`\-?h(elp)?`)
	if len(args) == 1 || (len(args) == 2 && helpFlagRegex.MatchString(args[1])) {
		showRootUsage()
		logrus.Fatal("You must provide a command to vdm")
	}
}
