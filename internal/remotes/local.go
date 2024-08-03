package remotes

import (
	"github.com/opensourcecorp/vdm/internal/vdmspec"
)

// SyncLocal is the root of the sync operations for "local" remote types.
func SyncLocal(remote vdmspec.Remote) error {
	return nil
}

func retrieveLocal(remote vdmspec.Remote) (err error) {
	return nil
}
