// Package flagvars houses constants etc. for working with command-line flag
// values across packages. These helpers are pushed down to their own package in
// order to avoid import cycles.
package flagvars

const (
	// Debug anchors to the DEBUG env var
	Debug = "DEBUG"
	// TryLocalSources anchors to the TRY_LOCAL_SOURCES env var
	TryLocalSources = "TRY_LOCAL_SOURCES"
)
