package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/store"
)

func cmdBackup(root string, args []string) int {
	fs := flag.NewFlagSet("backup", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	out := fs.String("o", "", "output path for trace.db snapshot")
	outLong := fs.String("output", "", "output path for trace.db snapshot")
	includeToken := fs.Bool("include-token", false, "also write access.token beside the snapshot")
	args = flagsFirst(args, map[string]bool{"o": true, "output": true, "include-token": false})
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	dest := *out
	if dest == "" {
		dest = *outLong
	}
	if dest == "" || fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace backup -o <path> [--include-token]\n")
		return exitUsage
	}
	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "backup: %v\n", err)
		return exitFail
	}
	fmt.Fprintf(os.Stderr, "backup: writing %s\n", dest)
	if err := store.Backup(abs, dest, store.BackupOptions{IncludeToken: *includeToken}); err != nil {
		fmt.Fprintf(os.Stderr, "backup: %v\n", err)
		return exitFail
	}
	fmt.Println(dest)
	return exitOK
}
