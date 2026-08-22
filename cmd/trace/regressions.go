package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/domain"
)

func cmdRegressions(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace regressions list\n")
		return exitUsage
	}
	switch args[0] {
	case "list":
		return cmdRegressionsList(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown regressions subcommand: %s\n", args[0])
		return exitUsage
	}
}

func cmdRegressionsList(root string, args []string) int {
	fs := flag.NewFlagSet("regressions list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "optional task UUID filter")
	changeID := fs.String("change", "", "optional change UUID filter")
	limit := fs.Int("limit", 0, "max rows (default 32, cap 64)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true, "change": true, "limit": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace regressions list [--task <id>] [--change <id>] [--limit N]\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	if code := failCLIDenied(svc, "regressions", "regressions list"); code != exitOK {
		return code
	}

	rows, err := svc.ListRegressions(context.Background(), domain.EvidenceQueryOpts{
		TaskID: *taskID, ChangeID: *changeID, Limit: *limit,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "regressions list: %v\n", err)
		return exitFail
	}
	if rows == nil {
		rows = []domain.RegressionQueryRow{}
	}
	if err := json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "regressions": rows, "count": len(rows),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "regressions list: %v\n", err)
		return exitFail
	}
	return exitOK
}
