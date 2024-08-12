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

// Git defines the git remote type
type Git struct {
	vdmspec.Remote
}

// Cache provides the [vdmspec.Remoter.Cache] operations for "git" remote types.
func (remote Git) Cache() error {
	return errors.New("not implemented")
}

// Sync provides the [vdmspec.Remoter.Sync] operations for "git" remote types.
func (remote Git) Sync() error {
	err := gitClone(remote)
	if err != nil {
		return fmt.Errorf("cloning git repository: %w", err)
	}

	message.Infof("%s: Setting specified version...", remote.OpMsg())
	checkoutCmd := exec.Command("git", "-C", remote.Destination, "checkout", remote.Version)
	checkoutOutput, err := checkoutCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error checking out specified revision: exec error '%w', with output: %s", err, string(checkoutOutput))
	}

	message.Debugf("removing .git dir for local path '%s'", remote.Destination)
	dotGitPath := filepath.Join(remote.Destination, ".git")
	err = os.RemoveAll(dotGitPath)
	if err != nil {
		return fmt.Errorf("removing directory %s: %w", dotGitPath, err)
	}

	return nil
}

// GetSource returns the Source field.
func (remote Git) GetSource() string {
	return remote.Source
}

// GetVersion returns the Version field.
func (remote Git) GetVersion() string {
	return remote.Version
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

func gitClone(remote Git) error {
	err := checkGitAvailable()
	if err != nil {
		return fmt.Errorf("remote '%s' is a git type, but git may not installed/available on PATH: %w", remote.Source, err)
	}

	cloneCmdArgs := []string{"clone", remote.Source, remote.Destination}

	message.Infof("%s: Retrieving...", remote.OpMsg())
	cloneCmd := exec.Command("git", cloneCmdArgs...)
	cloneOutput, err := cloneCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cloning remote: exec error '%w', with output: %s", err, string(cloneOutput))
	}

	return nil
}

var _ vdmspec.Remoter = Git{}
