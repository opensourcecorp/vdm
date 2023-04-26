package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
)

const (
	colorErr   = "\033[31m" // red
	colorDebug = "\033[33m" // yellow
	colorHappy = "\033[32m" // green
	colorReset = "\033[0m"
)

var (
	debug        bool
	specFilePath string

	infoLogger  = log.New(os.Stderr, fmt.Sprintf("%s[vdm] ", colorReset), 0)
	errLogger   = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm] ", colorReset, colorErr), 0)
	debugLogger = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm] ", colorReset, colorDebug), 0)
	happyLogger = log.New(os.Stderr, fmt.Sprintf("%s%s[vdm] ", colorReset, colorHappy), 0)
)

// In case I need to pass these around, so we're not relying on globals
type runFlags struct {
	SpecFilePath string
	Debug        bool
}

type depSpec struct {
	Remote    string `json:"remote"`
	Version   string `json:"version"`
	LocalPath string `json:"local_path"`
}

func main() {
	flag.StringVar(&specFilePath, "spec-file", "./.vdm", "vdm dependency spec file. Defaults to `.vdm` in the current directory")
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

	depsFile, err := os.ReadFile(specFilePath)
	if err != nil {
		if debug {
			debugLogger.Printf("error reading depsFile from disk: %v", err)
		}
		errLogger.Fatalf("There was a problem reading your vdm file from '%s' -- does it not exist?", specFilePath)
	}
	if debug {
		debugLogger.Printf("depsFile contents read:\n%s", string(depsFile))
	}

	var specs []depSpec
	err = json.Unmarshal(depsFile, &specs)
	if err != nil {
		if debug {
			debugLogger.Printf("error during depsFile unmarshal: %v", err)
		}
		errLogger.Fatal("There was a problem reading the contents of your vdm spec file")
	}
	if debug {
		debugLogger.Printf("depSpecs unmarshalled: %+v", specs)
	}

	for _, spec := range specs {
		err = spec.Validate(runFlags)
		if err != nil {
			errLogger.Fatalf("your vdm spec file is malformed: %v", err)
		}
	}

	for _, spec := range specs {
		if debug {
			debugLogger.Printf("removing any old data for '%s'", spec.LocalPath)
		}
		os.RemoveAll(spec.LocalPath)

		operationMsg := fmt.Sprintf("%s@%s --> %s", spec.Remote, spec.Version, spec.LocalPath)

		// If users want "latest", then we can just do a depth-one clone and
		// skip the checkout operation. But if they want non-latest, we need the
		// full history to be able to find a specified revision
		var cloneCmdArgs []string
		if spec.Version == "latest" {
			if debug {
				debugLogger.Printf("%s -- version specified as 'latest', so making shallow clone and skipping separate checkout operation", operationMsg)
			}
			cloneCmdArgs = []string{"clone", "--depth=1", spec.Remote, spec.LocalPath}
		} else {
			if debug {
				debugLogger.Printf("%s -- version specified as NOT latest, so making regular clone and will make separate checkout operation", operationMsg)
			}
			cloneCmdArgs = []string{"clone", spec.Remote, spec.LocalPath}
		}

		infoLogger.Printf("%s -- Retrieving...", operationMsg)
		cloneCmd := exec.Command("git", cloneCmdArgs...)
		cloneOutput, err := cloneCmd.CombinedOutput()
		if err != nil {
			errLogger.Fatalf("error cloning remote: exec error '%v', with output: %s", err, string(cloneOutput))
		}

		if spec.Version != "latest" {
			infoLogger.Printf("%s -- Setting specified version...", operationMsg)
			checkoutCmd := exec.Command("git", "-C", spec.LocalPath, "checkout", spec.Version)
			checkoutOutput, err := checkoutCmd.CombinedOutput()
			if err != nil {
				errLogger.Fatalf("error checking out specified revision: exec error '%v', with output: %s", err, string(checkoutOutput))
			}
		}

		if debug {
			debugLogger.Printf("removing .git dir for local path '%s'", spec.LocalPath)
		}
		os.RemoveAll(path.Join(spec.LocalPath, ".git"))

		infoLogger.Printf("%s -- Done.", operationMsg)
	}

	happyLogger.Print("All done!")
}

func checkGitAvailable() error {
	cmd := exec.Command("git", "--version")
	sysOutput, err := cmd.CombinedOutput()
	if err != nil {
		if debug {
			debugLogger.Printf("%s: %s", err.Error(), string(sysOutput))
		}
		return errors.New("git does not seem to be available on your PATH, so cannot continue")
	}
	if debug {
		debugLogger.Print("git was found on PATH")
	}
	return nil
}
