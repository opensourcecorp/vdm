package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	// ANSI color codes for log messages
	colorErr   = "\033[31m" // red
	colorDebug = "\033[33m" // yellow
	colorHappy = "\033[32m" // green
	colorReset = "\033[0m"
)

var (
	// CLI args
	debug        bool
	specFilePath string

	// Loggers, which include embedded ANSI color codes
	infoLogger  = log.New(os.Stderr, fmt.Sprintf("%s[vdm] ", colorReset), 0)
	errLogger   = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm]%s ", colorReset, colorErr, colorReset), 0)
	debugLogger = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm]%s ", colorReset, colorDebug, colorReset), 0)
	happyLogger = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm]%s ", colorReset, colorHappy, colorReset), 0)
)

// In case I need to pass these around, so we're not relying on globals
type runFlags struct {
	SpecFilePath string
	Debug        bool
}

func main() {
	flag.StringVar(&specFilePath, "spec-file", "./.vdm", "vdm dependency spec file")
	flag.BoolVar(&debug, "debug", false, "Print debug logs")
	flag.Parse()

	runFlags := runFlags{
		SpecFilePath: specFilePath,
		Debug:        debug,
	}

	err := checkGitAvailable()
	if err != nil {
		os.Exit(1)
	}

	specs := getSpecsFromFile(specFilePath, runFlags)

	for _, spec := range specs {
		err := spec.Validate()
		if err != nil {
			errLogger.Fatalf("your vdm spec file is malformed: %v", err)
		}
	}

	sync(specs)

	happyLogger.Print("All done!")
}
