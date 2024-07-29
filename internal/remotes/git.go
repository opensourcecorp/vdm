package remotes

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
)

// SyncGit is the root of the sync operations for "git" remote types.
func SyncGit(remote vdmspec.Remote, keepGitDir bool) error {
	err := gitClone(remote)
	if err != nil {
		return fmt.Errorf("cloing remote: %w", err)
	}

	if remote.Version != "latest" {
		message.Infof("%s: Setting specified version...", remote.OpMsg())
		checkoutCmd := exec.Command("git", "-C", remote.LocalPath, "checkout", remote.Version)
		checkoutOutput, err := checkoutCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error checking out specified revision: exec error '%w', with output: %s", err, string(checkoutOutput))
		}
	}

	if !keepGitDir {
		message.Debugf("removing .git dir for local path '%s'", remote.LocalPath)
		dotGitPath := filepath.Join(remote.LocalPath, ".git")
		err = os.RemoveAll(dotGitPath)
		if err != nil {
			return fmt.Errorf("removing directory %s: %w", dotGitPath, err)
		}
	}

	return nil
}

func checkGitAvailable() error {
	cmd := exec.Command("git", "--version")
	sysOutput, err := cmd.CombinedOutput()
	if err != nil {
		message.Debugf("%s: %s", err.Error(), string(sysOutput))
		return errors.New("git does not seem to be available on your PATH, so cannot continue")
	}
	message.Debugf("git was found on PATH")
	return nil
}

func gitClone(remote vdmspec.Remote) error {
	err := checkGitAvailable()
	if err != nil {
		return fmt.Errorf("remote '%s' is a git type, but git may not installed/available on PATH: %w", remote.Remote, err)
	}

	// If users want "latest", then we can just do a depth-one clone and
	// skip the checkout operation. But if they want non-latest, we need the
	// full history to be able to find a specified revision
	var cloneCmdArgs []string
	if remote.Version == "latest" {
		message.Debugf("%s: version specified as 'latest', so making shallow clone and skipping separate checkout operation", remote.OpMsg())
		cloneCmdArgs = []string{"clone", "--depth=1", remote.Remote, remote.LocalPath}
	} else {
		message.Debugf("%s: version specified as NOT latest, so making regular clone and will make separate checkout operation", remote.OpMsg())
		cloneCmdArgs = []string{"clone", remote.Remote, remote.LocalPath}
	}

	message.Infof("%s: Retrieving...", remote.OpMsg())
	cloneCmd := exec.Command("git", cloneCmdArgs...)
	cloneOutput, err := cloneCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cloning remote: exec error '%w', with output: %s", err, string(cloneOutput))
	}

	return nil
}
