package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

func cmdMigrate(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace migrate status\n")
		return exitUsage
	}
	switch args[0] {
	case "status":
		if len(args) > 1 {
			fmt.Fprintf(os.Stderr, "usage: trace migrate status\n")
			return exitUsage
		}
		return cmdMigrateStatus(root)
	default:
		fmt.Fprintf(os.Stderr, "unknown migrate subcommand: %s\n", args[0])
		fmt.Fprintf(os.Stderr, "usage: trace migrate status\n")
		return exitUsage
	}
}

func cmdMigrateStatus(root string) int {
	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate status: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate status: %v\n", err)
		return exitFail
	}
	defer st.Close()

	ms, err := st.MigrationStatus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "migrate status: %v\n", err)
		return exitFail
	}
	applied := make([]string, len(ms.AppliedVersions))
	for i, v := range ms.AppliedVersions {
		applied[i] = strconv.Itoa(v)
	}
	fmt.Printf("max_applied=%d\n", ms.MaxApplied)
	fmt.Printf("embed_expected=%d\n", ms.EmbedExpected)
	fmt.Printf("applied=%s\n", strings.Join(applied, ","))
	fmt.Printf("pending=%d\n", ms.PendingCount)
	return exitOK
}
