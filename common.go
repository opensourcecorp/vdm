// Common and/or initialization consts, vars, and functions
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

const (
	// ANSI color codes for log messages
	colorErr   = "\033[31m" // red
	colorDebug = "\033[33m" // yellow
	colorInfo  = "\033[36m" // cyan
	colorHappy = "\033[32m" // green
	colorReset = "\033[0m"
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

	// Loggers, which include embedded ANSI color codes
	infoLogger  = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm]%s ", colorReset, colorInfo, colorReset), 0)
	errLogger   = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm]%s ", colorReset, colorErr, colorReset), 0)
	debugLogger = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm]%s ", colorReset, colorDebug, colorReset), 0)
	happyLogger = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm]%s ", colorReset, colorHappy, colorReset), 0)
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

// Linter is mad about using string keys for context.Context, so define empty
// struct types for each usable key here
type debugKey struct{}
type specFilePathKey struct{}
type keepGitDirKey struct{}

// registerContextKeys assigns common values to the context that is passed
// around, such as CLI flags
func registerContextKeys() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, debugKey{}, debug)
	ctx = context.WithValue(ctx, specFilePathKey{}, specFilePath)
	ctx = context.WithValue(ctx, keepGitDirKey{}, keepGitDir)

	return ctx
}

// isDebug checks against the passed context to determine if the debug CLI flag
// was set by the user
func isDebug(ctx context.Context) bool {
	debugVal := ctx.Value(debugKey{})
	if debugVal == nil {
		return false
	}

	return debugVal.(bool)
}

// shouldKeepGitDir checks against the passed context to determine if the
// keepGitDir CLI flag was set by the user
func shouldKeepGitDir(ctx context.Context) bool {
	keepGitDirVal := ctx.Value(keepGitDirKey{})
	if keepGitDirVal == nil {
		return false
	}

	return keepGitDirVal.(bool)
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
		errLogger.Fatal("You must provide a command to vdm")
	}
}
