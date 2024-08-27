package vdminit

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/opensourcecorp/vdm/cmd/vars"
)

// SetupVDMForTest is a helper function for setting up vdm's tests. Notably, vdm
// requires itself to be initialized with its own home directory etc. before
// first use, and tests need a clean way to do that, and an equally-clean way to
// tear down their mess. This function returns both the string path to vdm's
// test home directory, as well as a cleanup function to be called by the test.
func SetupVDMForTest(t *testing.T) (vdmHome string, cleanup func()) {
	t.Helper()

	randID := 100000000000 + rand.Intn(999999999999)
	vdmHome = filepath.Join(os.TempDir(), fmt.Sprintf("vdm-tmp-%d", randID))
	t.Setenv(vars.VDMHomeEnvVarName, vdmHome)

	err := Paths()
	if err != nil {
		t.Errorf("instantiating vdm paths for test: %v", err)
	}

	cleanup = func() {
		err := os.RemoveAll(vdmHome)
		if err != nil {
			t.Errorf("failed to clean up after test: %v", err)
		}
	}

	return vdmHome, cleanup
}
