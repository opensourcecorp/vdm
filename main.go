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
	"regexp"
)

const (
	colorErr   = "\033[31m" // red
	colorDebug = "\033[33m" // yellow
	colorReset = "\033[0m"
)

var (
	debug  bool
	logger = log.New(os.Stderr, "[vdm] ", 0)
)

type depSpec struct {
	Remote    string `json:"remote"`
	Version   string `json:"version"`
	LocalPath string `json:"local_path"`
}

func main() {
	specFilePath := flag.String("spec-file", "./.vdm", "vdm dependency spec file. Defaults to `.vdm` in the current directory")
	flag.BoolVar(&debug, "debug", false, "Print debug logs")
	flag.Parse()

	err := checkGitAvailable()
	if err != nil {
		os.Exit(1)
	}
	if debug {
		logger.Printf("%sgit was found on PATH%s", colorDebug, colorReset)
	}

	depsFile, err := os.ReadFile(*specFilePath)
	if err != nil {
		if debug {
			logger.Printf("%serror reading depsFile from disk: %v%s", colorDebug, err, colorReset)
		}
		logger.Fatalf("%sThere was a problem reading your vdm file from '%s' -- does it not exist?%s", colorErr, *specFilePath, colorReset)
	}
	if debug {
		logger.Printf("%sdepsFile contents read:\n%s%s", colorDebug, string(depsFile), colorReset)
	}

	var specs []depSpec
	err = json.Unmarshal(depsFile, &specs)
	if err != nil {
		if debug {
			logger.Printf("%serror during depsFile unmarshal: %v%s", colorDebug, err, colorReset)
		}
		logger.Fatalf("%sThere was a problem reading the contents of your vdm spec file%s", colorErr, colorReset)
	}
	if debug {
		logger.Printf("%sdepSpecs unmarshalled: %+v%s", colorDebug, specs, colorReset)
	}

	for _, spec := range specs {
		err = spec.Validate()
		if err != nil {
			logger.Fatalf("%syour vdm spec file is malformed: %v%s", colorErr, err, colorReset)
		}
	}

	for _, spec := range specs {
		if debug {
			logger.Printf("%sremoving any old data for '%s'%s", colorDebug, spec.LocalPath, colorReset)
		}
		os.RemoveAll(spec.LocalPath)

		operationMsg := fmt.Sprintf("%s@%s --> %s", spec.Remote, spec.Version, spec.LocalPath)

		// If users want "latest", then we can just do a depth-one clone and
		// skip the checkout operation. But if they want non-latest, we need the
		// full history to be able to find a specified revision
		var cloneCmdArgs []string
		if spec.Version == "latest" {
			if debug {
				logger.Printf("%s%s version specified as 'latest', so making shallow clone and skipping separate checkout operation%s", colorDebug, operationMsg, colorReset)
			}
			cloneCmdArgs = []string{"clone", "--depth=1", spec.Remote, spec.LocalPath}
		} else {
			if debug {
				logger.Printf("%s%s version specified as NOT latest, so making regular clone and will make separate checkout operation%s", colorDebug, operationMsg, colorReset)
			}
			cloneCmdArgs = []string{"clone", spec.Remote, spec.LocalPath}
		}

		logger.Printf("Retrieving %s ...", operationMsg)
		cloneCmd := exec.Command("git", cloneCmdArgs...)
		cloneOutput, err := cloneCmd.CombinedOutput()
		if err != nil {
			logger.Fatalf("%serror cloning remote: exec error '%v', with output: %s%s", colorErr, err, string(cloneOutput), colorReset)
		}

		if spec.Version != "latest" {
			logger.Printf("Using specified version for %s", operationMsg)
			checkoutCmd := exec.Command("git", "-C", spec.LocalPath, "checkout", spec.Version)
			checkoutOutput, err := checkoutCmd.CombinedOutput()
			if err != nil {
				logger.Fatalf("%serror checking out specified revision: exec error '%v', with output: %s%s", colorErr, err, string(checkoutOutput), colorReset)
			}
		}

		if debug {
			logger.Printf("%sremoving .git dir for local path '%s'%s", colorDebug, spec.LocalPath, colorReset)
		}
		os.RemoveAll(path.Join(spec.LocalPath, ".git"))

		logger.Printf("Done with %s", operationMsg)
	}

	logger.Print("All done!")
}

func checkGitAvailable() error {
	cmd := exec.Command("git", "--version")
	sysOutput, err := cmd.CombinedOutput()
	if err != nil {
		if debug {
			logger.Printf("%s%v: %v%s", colorDebug, err.Error(), string(sysOutput), colorReset)
		}
		return errors.New("git does not seem to be available on your PATH, so cannot continue")
	}
	return nil
}

func (spec depSpec) Validate() error {
	var allErrors []error

	if debug {
		logger.Print(colorDebug, "validating field 'Remote'", colorReset)
	}
	if len(spec.Remote) == 0 {
		allErrors = append(allErrors, errors.New("all 'remote' fields must be non-zero length"))
	}
	protocolRegex := regexp.MustCompile(`(https://|git://git@)`)
	if !protocolRegex.MatchString(spec.Remote) {
		allErrors = append(
			allErrors,
			fmt.Errorf("remote provided as '%s', but all 'remote' fields must begin with a protocol specifier or other valid prefix (e.g. 'https://', 'git@', etc.)", spec.Remote),
		)
	}

	if debug {
		logger.Print(colorDebug, "validating field 'Version'", colorReset)
	}
	if len(spec.Version) == 0 {
		allErrors = append(allErrors, errors.New("all 'version' fields must be non-zero length. If you don't care about the version (even though you should), then use 'latest'"))
	}

	if debug {
		logger.Print(colorDebug, "validating field 'LocalPath'", colorReset)
	}
	if len(spec.LocalPath) == 0 {
		allErrors = append(allErrors, errors.New("all 'local_path' fields must be non-zero length"))
	}

	if len(allErrors) > 0 {
		for _, err := range allErrors {
			logger.Printf("validation failure: %s", err.Error())
		}
		return fmt.Errorf("%s%d validation failure(s) found in your vdm spec file%s", colorErr, len(allErrors), colorReset)
	}

	return nil
}
