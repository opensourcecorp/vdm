package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// sync ensures that the only local dependencies are ones defined in the specfile
func sync(ctx context.Context, specs []vdmSpec) {
	for _, spec := range specs {
		// Common log line prefix
		opMsg := fmt.Sprintf("%s@%s --> %s", spec.Remote, spec.Version, spec.LocalPath)

		// process stored VDMMETA so we know what operations to actually perform for existing directories
		vdmMeta := spec.getVDMMeta()
		if vdmMeta == (vdmSpec{}) {
			logrus.Infof("VDMMETA not found under local path '%s' -- will be created", spec.LocalPath)
		} else {
			if vdmMeta.Version != spec.Version {
				logrus.Infof("Changing '%s' from current local version spec '%s' to '%s'...", spec.Remote, vdmMeta.Version, spec.Version)
			} else {
				if isDebug(ctx) {
					logrus.Debugf("Version unchanged (%s) in spec file for '%s' --> '%s'", spec.Version, spec.Remote, spec.LocalPath)
				}
			}
		}

		switch spec.Type {
		case "git", "":
			syncGitRemote(ctx, spec, opMsg)
		default:
			logrus.Fatalf("")
		}

		err := spec.writeVDMMeta()
		if err != nil {
			logrus.Fatalf("Could not write VDMMETA file to disk: %v", err)
		}

		logrus.Infof("%s -- Done.", opMsg)
	}
}

func syncGitRemote(ctx context.Context, spec vdmSpec, operationMsg string) {
	// TODO: pull this up so that it only runs if the version changed or the user requested a wipe
	if !shouldKeepGitDir(ctx) {
		if isDebug(ctx) {
			logrus.Debugf("removing any old data for '%s'", spec.LocalPath)
		}
		os.RemoveAll(spec.LocalPath)
	}

	gitClone(ctx, spec, operationMsg)

	if spec.Version != "latest" {
		logrus.Infof("%s -- Setting specified version...", operationMsg)
		checkoutCmd := exec.Command("git", "-C", spec.LocalPath, "checkout", spec.Version)
		checkoutOutput, err := checkoutCmd.CombinedOutput()
		if err != nil {
			logrus.Fatalf("error checking out specified revision: exec error '%v', with output: %s", err, string(checkoutOutput))
		}
	}

	if !shouldKeepGitDir(ctx) {
		if isDebug(ctx) {
			logrus.Debugf("removing .git dir for local path '%s'", spec.LocalPath)
		}
		os.RemoveAll(filepath.Join(spec.LocalPath, ".git"))
	}
}

func gitClone(ctx context.Context, spec vdmSpec, operationMsg string) {
	// If users want "latest", then we can just do a depth-one clone and
	// skip the checkout operation. But if they want non-latest, we need the
	// full history to be able to find a specified revision
	var cloneCmdArgs []string
	if spec.Version == "latest" {
		if isDebug(ctx) {
			logrus.Debugf("%s -- version specified as 'latest', so making shallow clone and skipping separate checkout operation", operationMsg)
		}
		cloneCmdArgs = []string{"clone", "--depth=1", spec.Remote, spec.LocalPath}
	} else {
		if isDebug(ctx) {
			logrus.Debugf("%s -- version specified as NOT latest, so making regular clone and will make separate checkout operation", operationMsg)
		}
		cloneCmdArgs = []string{"clone", spec.Remote, spec.LocalPath}
	}

	logrus.Infof("%s -- Retrieving...", operationMsg)
	cloneCmd := exec.Command("git", cloneCmdArgs...)
	cloneOutput, err := cloneCmd.CombinedOutput()
	if err != nil {
		logrus.Fatalf("error cloning remote: exec error '%v', with output: %s", err, string(cloneOutput))
	}
}
