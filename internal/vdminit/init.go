package vdminit

import (
	"fmt"
	"os"

	"github.com/opensourcecorp/vdm/cmd/vars"
	"github.com/opensourcecorp/vdm/internal/archive/cache"
)

func Paths() error {
	vdmCacheDir, err := vars.GetVDMCacheDir()
	if err != nil {
		return fmt.Errorf("determining vdm cache directory while running init: %w", err)
	}

	err = os.MkdirAll(vdmCacheDir, 0755)
	if err != nil {
		return fmt.Errorf("creating vdm cache directory %q during init: %w", vdmCacheDir, err)
	}

	err = cache.CreateSumDB()
	if err != nil {
		return fmt.Errorf("creating sumdb during init: %w", err)
	}

	return nil
}
