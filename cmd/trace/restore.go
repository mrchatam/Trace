package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/store"
)

func cmdRestore(root string, args []string) int {
	fs := flag.NewFlagSet("restore", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	force := fs.Bool("force", false, "overwrite existing .trace/trace.db")
	args = flagsFirst(args, map[string]bool{"force": false})
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: trace restore <backup.db> [--force]\n")
		return exitUsage
	}
	src := fs.Arg(0)
	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "restore: %v\n", err)
		return exitFail
	}
	fmt.Fprintf(os.Stderr, "restore: installing into %s/.trace/trace.db\n", abs)
	if err := store.Restore(abs, src, *force); err != nil {
		fmt.Fprintf(os.Stderr, "restore: %v\n", err)
		return exitFail
	}
	fmt.Println(abs)
	return exitOK
}
