package main

import (
	"os"
)

func main() {
	checkRootUsage(os.Args)
	cmd, ok := subcommands[os.Args[1]] // length-guarded already by checkRootUsage() above
	if !ok {
		showRootUsage()
		errLogger.Fatalf("Unrecognized vdm subcommand '%s'", os.Args[1])
	}
	registerFlags()
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

	switch cmd.Name() {
	case syncCmd.Name():
		sync(ctx, specs)
	default: // should never get here since we check above, but still
		showRootUsage()
		errLogger.Fatalf("Unrecognized vdm subcommand '%s'", cmd.Name())
	}

	happyLogger.Print("All done!")
}
