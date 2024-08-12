package remotes

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/opensourcecorp/vdm/cmd/vars"
	"github.com/opensourcecorp/vdm/internal/archive/sumdb"
	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/opensourcecorp/vdm/internal/vdmspec"
)

// Git defines the git remote type, and implements the [vdmspec.Remoter]
// interface.
type Git struct {
	vdmspec.RemoteTemplate
}

// Cache provides the [vdmspec.Remoter.Cache] operations for "git" remote types.
func (remote Git) Cache() (err error) {
	tmpCachePath := filepath.Join(os.TempDir(), "vdm-tmp", filepath.Base(remote.Destination))
	message.Debugf("tmpCachePath: %s", tmpCachePath)

	if err := os.RemoveAll(tmpCachePath); err != nil {
		return fmt.Errorf("trying to clean up possibly-duplicate old cache data at %s: %w", tmpCachePath, err)
	}

	err = gitClone(remote, tmpCachePath)
	if err != nil {
		return fmt.Errorf("cloning git repository: %w", err)
	}
	defer func() {
		if rmErr := os.RemoveAll(tmpCachePath); rmErr != nil {
			err = errors.Join(err, fmt.Errorf("removing temporary cache directory for '%s': %w", remote.GetSource(), rmErr))
		}
	}()

	message.Infof("%s: Setting specified version...", remote.OpMsg())
	checkoutCmd := exec.Command("git", "-C", tmpCachePath, "checkout", remote.Version)
	checkoutOutput, err := checkoutCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error checking out specified revision: exec error '%w', with output: %s", err, string(checkoutOutput))
	}

	message.Debugf("removing .git dir for local path '%s'", tmpCachePath)
	dotGitPath := filepath.Join(tmpCachePath, ".git")
	err = os.RemoveAll(dotGitPath)
	if err != nil {
		return fmt.Errorf("removing directory %s: %w", dotGitPath, err)
	}

	// TODO-NOW: we need this to create two paths: one for the actual
	// gzipped-tar cache, and one for the sumdb file
	cachePath, err := os.Create(filepath.Join(vars.GetVDMCacheDir()))
	sumdb.CacheRemote(remote)

	// return nil

	return errors.New("not implemented")
}

// Sync provides the [vdmspec.Remoter.Sync] operations for "git" remote types.
func (remote Git) Sync() error {
	return errors.New("not implemented")
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

func gitClone(remote Git, dest string) error {
	err := checkGitAvailable()
	if err != nil {
		return fmt.Errorf("remote '%s' is a git type, but git may not installed/available on PATH: %w", remote.Source, err)
	}

	cloneCmdArgs := []string{"clone", remote.Source, dest}
	message.Debugf("git args: %v", cloneCmdArgs)

	message.Infof("%s: Retrieving...", remote.OpMsg())
	cloneCmd := exec.Command("git", cloneCmdArgs...)
	cloneOutput, err := cloneCmd.CombinedOutput()
	message.Debugf("git clone command output: %s", string(cloneOutput))
	if err != nil {
		return fmt.Errorf("cloning remote: exec error '%w', with output: %s", err, string(cloneOutput))
	}

	return nil
}

var _ vdmspec.Remoter = Git{}
