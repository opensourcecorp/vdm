package main

import (
	"errors"
	"os/exec"
)

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
