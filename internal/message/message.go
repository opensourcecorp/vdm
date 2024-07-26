// Package message controls message printing. THis isn't a "logging" package per
// se, but adds some niceties for log-like needs.
package message

import (
	"fmt"
	"os"
)

func Debugf(format string, args ...any) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf("DEBUG: "+format+"\n", args...)
	}
}

func Infof(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

func Warnf(format string, args ...any) {
	fmt.Printf("WARNING: "+format+"\n", args...)
}

func Errorf(format string, args ...any) {
	fmt.Printf("ERROR: "+format+"\n", args...)
}

func Fatalf(format string, args ...any) {
	fmt.Printf("ERROR: "+format+"\n", args...)
	os.Exit(1)
}
