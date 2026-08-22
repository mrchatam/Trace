package main

import (
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/store"
)

func cmdAuth(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace auth set <token> | clear | status\n")
		return exitUsage
	}
	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "auth: %v\n", err)
		return exitFail
	}
	switch args[0] {
	case "set":
		if len(args) != 2 {
			fmt.Fprintf(os.Stderr, "usage: trace auth set <token>\n")
			return exitUsage
		}
		if err := store.SetAccessToken(abs, args[1]); err != nil {
			fmt.Fprintf(os.Stderr, "auth set: %v\n", err)
			return exitFail
		}
		fmt.Println("enabled")
		return exitOK
	case "clear":
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "usage: trace auth clear\n")
			return exitUsage
		}
		if err := store.ClearAccessToken(abs); err != nil {
			fmt.Fprintf(os.Stderr, "auth clear: %v\n", err)
			return exitFail
		}
		fmt.Println("disabled")
		return exitOK
	case "status":
		if len(args) != 1 {
			fmt.Fprintf(os.Stderr, "usage: trace auth status\n")
			return exitUsage
		}
		enabled, err := store.AccessTokenEnabled(abs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "auth status: %v\n", err)
			return exitFail
		}
		if enabled {
			fmt.Println("enabled")
		} else {
			fmt.Println("disabled")
		}
		return exitOK
	default:
		fmt.Fprintf(os.Stderr, "unknown auth subcommand: %s\n", args[0])
		fmt.Fprintf(os.Stderr, "usage: trace auth set <token> | clear | status\n")
		return exitUsage
	}
}
