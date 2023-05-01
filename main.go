package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
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

	// Loggers, which include embedded ANSI color codes
	infoLogger  = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm]%s ", colorReset, colorInfo, colorReset), 0)
	errLogger   = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm]%s ", colorReset, colorErr, colorReset), 0)
	debugLogger = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm]%s ", colorReset, colorDebug, colorReset), 0)
	happyLogger = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm]%s ", colorReset, colorHappy, colorReset), 0)
)

func setCommonFlags() {
	for _, cmd := range subcommands {
		cmd.StringVar(&specFilePath, "spec-file", "./.vdm", "vdm dependency spec file")
		cmd.BoolVar(&debug, "debug", false, "Print debug logs")
	}
}

func isDebug(ctx context.Context) bool {
	debugVal := ctx.Value("debug")
	if debugVal == nil {
		panic("somehow the debug context key ended up as <nil>")
	}

	return debugVal.(bool)
}

func main() {
	if len(os.Args) == 1 {
		errLogger.Fatal("You must provide a command to vdm")
	}
	cmd, ok := subcommands[os.Args[1]]
	if !ok {
		errLogger.Fatalf("Unrecognized vmd subcommand '%s'", os.Args[1])
	}
	setCommonFlags()
	cmd.Parse(os.Args[2:])

	ctx := context.WithValue(context.Background(), "debug", debug)
	ctx = context.WithValue(ctx, "specFilePath", specFilePath)

	err := checkGitAvailable(ctx)
	if err != nil {
		os.Exit(1)
	}

	specs := getSpecsFromFile(ctx, specFilePath)

	for _, spec := range specs {
		err := spec.Validate(ctx)
		if err != nil {
			errLogger.Fatalf("Your vdm spec file is malformed: %v", err)
		}
	}

	sync(ctx, specs)

	happyLogger.Print("All done!")
}
