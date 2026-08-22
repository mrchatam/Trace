package main

import (
	"fmt"
	"os"
)

func cmdTests(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace tests verifying\n")
		return exitUsage
	}
	switch args[0] {
	case "verifying":
		return cmdTestVerifying(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown tests subcommand: %s\n", args[0])
		return exitUsage
	}
}
