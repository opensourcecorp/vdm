// Package message controls message printing. THis isn't a "logging" package per
// se, but adds some niceties for log-like needs.
package message

import (
	"fmt"
	"os"

	"github.com/opensourcecorp/vdm/cmd/flagvars"
)

// Debugf prints out debug-level information messages with a formatting
// directive.
func Debugf(format string, args ...any) {
	if os.Getenv(flagvars.Debug) != "" {
		fmt.Printf("DEBUG: "+format+"\n", args...)
	}
}

// Infof prints out debug-level information messages with a formatting
// directive.
func Infof(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

// Warnf prints out debug-level information messages with a formatting
// directive.
func Warnf(format string, args ...any) {
	fmt.Printf("WARNING: "+format+"\n", args...)
}

// Errorf prints out debug-level information messages with a formatting
// directive.
func Errorf(format string, args ...any) {
	fmt.Printf("ERROR: "+format+"\n", args...)
}

// Fatalf prints out debug-level information messages with a formatting
// directive, and then exits with code 1.
func Fatalf(format string, args ...any) {
	fmt.Printf("ERROR: "+format+"\n", args...)
	os.Exit(1)
}
