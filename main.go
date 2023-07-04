package main

import (
	"github.com/opensourcecorp/vdm/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logrus.Fatalf("running vdm: %v", err)
	}
}
