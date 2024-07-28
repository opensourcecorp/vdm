package cmd

import (
	"os"

	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/spf13/viper"
)

// MaybeSetDebug sets the DEBUG environment variable if it was set as a flag by
// the caller.
func MaybeSetDebug() {
	if viper.GetBool(debugFlagKey) {
		err := os.Setenv("DEBUG", "true")
		if err != nil {
			message.Fatalf("internal error: unable to set environment variable DEBUG")
		}
	}
}
