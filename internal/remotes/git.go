package remotes

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/opensourcecorp/vdm/internal/archive/cache"
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
	// tmpCachePath is where the actual retrieval is targeted, which is then
	// later archived to the persistent cache
	tmpCacheRoot := filepath.Join(os.TempDir(), "vdm-tmp")
	tmpCachePath := filepath.Join(tmpCacheRoot, filepath.Base(remote.Source))
	message.Debugf("tmpCachePath: %s", tmpCachePath)
	defer func() {
		if rmErr := os.RemoveAll(tmpCacheRoot); rmErr != nil {
			err = errors.Join(err, fmt.Errorf("removing temporary cache path '%s': %w", tmpCacheRoot, rmErr))
		}
	}()

	if err := os.RemoveAll(tmpCacheRoot); err != nil {
		return fmt.Errorf("trying to clean up possibly-duplicate old temp cache data at '%s': %w", tmpCachePath, err)
	}

	message.Infof("%s: Retrieving...", remote.OpMsg())
	err = gitClone(remote.Source, tmpCachePath)
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

	err = cache.AddRemote(remote, tmpCachePath)
	if err != nil {
		return fmt.Errorf("caching git remote %s: %w", remote.Source, err)
	}

	return err
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

func gitClone(src string, dest string) error {
	err := checkGitAvailable()
	if err != nil {
		return fmt.Errorf("remote '%s' is a git type, but git may not installed/available on PATH: %w", src, err)
	}

	cloneCmdArgs := []string{"clone", src, dest}
	message.Debugf("git args: %v", cloneCmdArgs)

	cloneCmd := exec.Command("git", cloneCmdArgs...)
	cloneOutput, err := cloneCmd.CombinedOutput()
	message.Debugf("git clone command output: %s", string(cloneOutput))
	if err != nil {
		return fmt.Errorf("cloning remote: exec error '%w', with output: %s", err, string(cloneOutput))
	}

	return nil
}

var _ vdmspec.Remoter = Git{}
