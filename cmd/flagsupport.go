package cmd

import (
	"os"

	"github.com/opensourcecorp/vdm/cmd/flagvars"
	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/spf13/viper"
)

// maybeSetDebug sets the DEBUG environment variable if it was set as a flag by
// the caller.
func maybeSetDebug() {
	if viper.GetBool(debugFlagKey) {
		err := os.Setenv(flagvars.Debug, "true")
		if err != nil {
			message.Fatalf("internal error: unable to set environment variable %s", flagvars.Debug)
		}
	}
}

// maybeTryLocalSources sets the TRY_LOCAL_SOURCES environment variable if it
// was set as a flag by the caller.
func maybeTryLocalSources() {
	if viper.GetBool(tryLocalSourcesFlagKey) {
		err := os.Setenv(flagvars.TryLocalSources, "true")
		if err != nil {
			message.Fatalf("internal error: unable to set environment variable %s", flagvars.TryLocalSources)
		}
	}
}
