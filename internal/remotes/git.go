package remotes

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/opensourcecorp/vdm/internal/archive"
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
func (remote Git) Cache() (cachePath string, err error) {
	tmpCachePath := cache.GetTempCachePath(remote)
	message.Debugf("tmpCachePath: %q", tmpCachePath)
	defer func() {
		if rmErr := os.RemoveAll(filepath.Dir(tmpCachePath)); rmErr != nil {
			err = errors.Join(err, fmt.Errorf("removing temporary cache path %q: %w", tmpCachePath, rmErr))
		}
	}()

	remoteAlreadyInSumDB, err := cache.CheckIfRemoteInSumDB(remote)
	if err != nil {
		return "", fmt.Errorf("checking if remote %q already in sumdb: %w", remote.GetSumDBKey(), err)
	}

	if !remoteAlreadyInSumDB {
		remote.OpMsg("Retrieving...")
		err = gitClone(remote.Source, tmpCachePath)
		if err != nil {
			return "", fmt.Errorf("cloning git repository: %w", err)
		}
		defer func() {
			if rmErr := os.RemoveAll(tmpCachePath); rmErr != nil {
				err = errors.Join(err, fmt.Errorf("removing temporary cache directory for %q: %w", remote.GetSource(), rmErr))
			}
		}()

		remote.OpMsg("Setting specified version...")
		checkoutCmd := exec.Command("git", "-C", tmpCachePath, "checkout", remote.Version)
		checkoutOutput, err := checkoutCmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("error checking out specified revision: exec error '%w', with output: %s", err, string(checkoutOutput))
		}

		message.Debugf("removing .git dir for local path %q", tmpCachePath)
		dotGitPath := filepath.Join(tmpCachePath, ".git")
		err = os.RemoveAll(dotGitPath)
		if err != nil {
			return "", fmt.Errorf("removing directory %q: %w", dotGitPath, err)
		}

		cachePath, err = cache.AddRemote(remote, tmpCachePath)
		if err != nil {
			return "", fmt.Errorf("caching git remote %q: %w", remote.Source, err)
		}
	} else {
		remote.OpMsg("Remote found in local cache; restoring...")
		cachePath, err = cache.GetPersistentCacheFilePath(remote)
		if err != nil {
			return "", fmt.Errorf("getting existing cache path for remote %q: %w", remote.GetSumDBKey(), err)
		}
	}

	return cachePath, err
}

// Sync provides the [vdmspec.Remoter.Sync] operations for "git" remote types.
func (remote Git) Sync(src, dest string) error {
	err := archive.ExtractTGZArchive(src, dest)
	if err != nil {
		return fmt.Errorf("syncing git cache for remote %q: %w", remote.GetSource(), err)
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

// GetSumDBKey returns the Source & Version fields, concatenated with an '@'.
func (remote Git) GetSumDBKey() string {
	return fmt.Sprintf("%s@%s", remote.Source, remote.Version)
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
		return fmt.Errorf("remote %q is a git type, but git may not installed/available on PATH: %w", src, err)
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
