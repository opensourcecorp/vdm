package main

import (
	"context"
	"errors"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func checkGitAvailable(ctx context.Context) error {
	cmd := exec.Command("git", "--version")
	sysOutput, err := cmd.CombinedOutput()
	if err != nil {
		if isDebug(ctx) {
			logrus.Debugf("%s: %s", err.Error(), string(sysOutput))
		}
		return errors.New("git does not seem to be available on your PATH, so cannot continue")
	}
	if isDebug(ctx) {
		logrus.Debug("git was found on PATH")
	}
	return nil
}
