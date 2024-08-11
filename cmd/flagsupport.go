package cmd

import (
	"os"

	"github.com/opensourcecorp/vdm/cmd/vars"
	"github.com/opensourcecorp/vdm/internal/message"
	"github.com/spf13/viper"
)

// maybeSetDebug sets the DEBUG environment variable if it was set as a flag by
// the caller.
func maybeSetDebug() {
	if viper.GetBool(debugFlagKey) {
		err := os.Setenv(vars.Debug, "true")
		if err != nil {
			message.Fatalf("internal error: unable to set environment variable %s", vars.Debug)
		}
	}
}

// maybeTryLocalSources sets the TRY_LOCAL_SOURCES environment variable if it
// was set as a flag by the caller.
func maybeTryLocalSources() {
	if viper.GetBool(tryLocalSourcesFlagKey) {
		err := os.Setenv(vars.TryLocalSources, "true")
		if err != nil {
			message.Fatalf("internal error: unable to set environment variable %s", vars.TryLocalSources)
		}
	}
}
