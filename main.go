package main

import (
	"github.com/opensourcecorp/vdm/cmd"
	"github.com/opensourcecorp/vdm/internal/message"
)

func main() {
	if err := cmd.Execute(); err != nil {
		message.Fatalf("running vdm: %v", err)
	}
}
