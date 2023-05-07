package main

import (
	"os"
)

func main() {
	checkRootUsage(os.Args)
	cmd, ok := subcommands[os.Args[1]] // length-guarded already by checkRootUsage() above
	if !ok {
		errLogger.Fatalf("Unrecognized vdm subcommand '%s'", os.Args[1])
	}
	registerCommonFlags()
	cmd.Parse(os.Args[2:])

	ctx := registerContextKeys()

	err := checkGitAvailable(ctx)
	if err != nil {
		os.Exit(1)
	}

	specs := getSpecsFromFile(ctx, specFilePath)

	for _, spec := range specs {
		err := spec.Validate(ctx)
		if err != nil {
			errLogger.Fatalf("Your vdm spec file is malformed: %v", err)
		}
	}

	sync(ctx, specs)

	happyLogger.Print("All done!")
}
