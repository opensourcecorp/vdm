package remote

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/opensourcecorp/vdm/internal/vdmspec"
	"github.com/sirupsen/logrus"
)

// SyncGit is the root of the sync operations for "git" remote types.
func SyncGit(spec vdmspec.Spec) error {
	err := gitClone(spec)
	if err != nil {
		return fmt.Errorf("cloing remote: %w", err)
	}

	if spec.Version != "latest" {
		logrus.Infof("%s -- Setting specified version...", spec.OpMsg())
		checkoutCmd := exec.Command("git", "-C", spec.LocalPath, "checkout", spec.Version)
		checkoutOutput, err := checkoutCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error checking out specified revision: exec error '%w', with output: %s", err, string(checkoutOutput))
		}
	}

	logrus.Debugf("removing .git dir for local path '%s'", spec.LocalPath)
	dotGitPath := filepath.Join(spec.LocalPath, ".git")
	err = os.RemoveAll(dotGitPath)
	if err != nil {
		return fmt.Errorf("removing directory %s: %w", dotGitPath, err)
	}

	return nil
}

func checkGitAvailable() error {
	cmd := exec.Command("git", "--version")
	sysOutput, err := cmd.CombinedOutput()
	if err != nil {
		logrus.Debugf("%s: %s", err.Error(), string(sysOutput))
		return errors.New("git does not seem to be available on your PATH, so cannot continue")
	}
	logrus.Debug("git was found on PATH")
	return nil
}

func gitClone(spec vdmspec.Spec) error {
	err := checkGitAvailable()
	if err != nil {
		return fmt.Errorf("remote '%s' is a git type, but git may not installed/available on PATH: %w", spec.Remote, err)
	}

	// If users want "latest", then we can just do a depth-one clone and
	// skip the checkout operation. But if they want non-latest, we need the
	// full history to be able to find a specified revision
	var cloneCmdArgs []string
	if spec.Version == "latest" {
		logrus.Debugf("%s -- version specified as 'latest', so making shallow clone and skipping separate checkout operation", spec.OpMsg())
		cloneCmdArgs = []string{"clone", "--depth=1", spec.Remote, spec.LocalPath}
	} else {
		logrus.Debugf("%s -- version specified as NOT latest, so making regular clone and will make separate checkout operation", spec.OpMsg())
		cloneCmdArgs = []string{"clone", spec.Remote, spec.LocalPath}
	}

	logrus.Infof("%s -- Retrieving...", spec.OpMsg())
	cloneCmd := exec.Command("git", cloneCmdArgs...)
	cloneOutput, err := cloneCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cloning remote: exec error '%w', with output: %s", err, string(cloneOutput))
	}

	return nil
}
