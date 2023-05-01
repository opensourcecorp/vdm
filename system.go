package main

import (
	"context"
	"errors"
	"os/exec"
)

func checkGitAvailable(ctx context.Context) error {
	cmd := exec.Command("git", "--version")
	sysOutput, err := cmd.CombinedOutput()
	if err != nil {
		if isDebug(ctx) {
			debugLogger.Printf("%s: %s", err.Error(), string(sysOutput))
		}
		return errors.New("git does not seem to be available on your PATH, so cannot continue")
	}
	if isDebug(ctx) {
		debugLogger.Print("git was found on PATH")
	}
	return nil
}
