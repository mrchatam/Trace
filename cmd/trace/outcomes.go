package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func cmdOutcomes(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace outcomes compare|improvements|failed|worked\n")
		return exitUsage
	}
	switch args[0] {
	case "compare":
		return cmdOutcomesCompare(root, args[1:])
	case "improvements":
		return cmdOutcomesImprovements(root, args[1:])
	case "failed":
		return cmdOutcomesFailed(root, args[1:])
	case "worked":
		return cmdOutcomesWorked(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown outcomes subcommand: %s\n", args[0])
		return exitUsage
	}
}

func cmdOutcomesCompare(root string, args []string) int {
	fs := flag.NewFlagSet("outcomes compare", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "task UUID")
	kind := fs.String("kind", "", "test or evaluation")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true, "kind": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 || strings.TrimSpace(*taskID) == "" || strings.TrimSpace(*kind) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace outcomes compare --task <id> --kind test|evaluation\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	if code := failCLIDenied(svc, "outcomes", "outcomes compare"); code != exitOK {
		return code
	}

	result, err := svc.CompareIterationOutcomes(context.Background(), *taskID, *kind)
	if err != nil {
		fmt.Fprintf(os.Stderr, "outcomes compare: %v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "outcomes compare: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdOutcomesImprovements(root string, args []string) int {
	fs := flag.NewFlagSet("outcomes improvements", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	changeID := fs.String("change", "", "change UUID")
	taskID := fs.String("task", "", "task UUID")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"change": true, "task": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace outcomes improvements --change <id> | --task <id>\n")
		return exitUsage
	}
	hasChange := strings.TrimSpace(*changeID) != ""
	hasTask := strings.TrimSpace(*taskID) != ""
	if hasChange == hasTask {
		fmt.Fprintf(os.Stderr, "usage: trace outcomes improvements --change <id> | --task <id>\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	if code := failCLIDenied(svc, "outcomes", "outcomes improvements"); code != exitOK {
		return code
	}

	ctx := context.Background()
	var rows []store.Improvement
	var err error
	if hasChange {
		rows, err = svc.ListImprovementsByChangeID(ctx, *changeID)
	} else {
		rows, err = svc.ListImprovementsByTaskID(ctx, *taskID)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "outcomes improvements: %v\n", err)
		return exitFail
	}
	payload := map[string]any{"ok": true, "improvements": rows, "count": len(rows)}
	if err := json.NewEncoder(os.Stdout).Encode(payload); err != nil {
		fmt.Fprintf(os.Stderr, "outcomes improvements: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdOutcomesFailed(root string, args []string) int {
	fs := flag.NewFlagSet("outcomes failed", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "optional task UUID filter")
	limit := fs.Int("limit", 0, "max rows (default 32, cap 64)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true, "limit": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace outcomes failed [--task <id>] [--limit N]\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	if code := failCLIDenied(svc, "outcomes", "outcomes failed"); code != exitOK {
		return code
	}

	rows, err := svc.ListFailedOutcomes(context.Background(), domain.EvidenceQueryOpts{
		TaskID: *taskID, Limit: *limit,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "outcomes failed: %v\n", err)
		return exitFail
	}
	if rows == nil {
		rows = []domain.FailedOutcomeRow{}
	}
	if err := json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "failed": rows, "count": len(rows),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "outcomes failed: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdOutcomesWorked(root string, args []string) int {
	fs := flag.NewFlagSet("outcomes worked", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "optional task UUID filter")
	limit := fs.Int("limit", 0, "max rows (default 32, cap 64)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true, "limit": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace outcomes worked [--task <id>] [--limit N]\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	if code := failCLIDenied(svc, "outcomes", "outcomes worked"); code != exitOK {
		return code
	}

	rows, err := svc.ListWorkedApproaches(context.Background(), domain.EvidenceQueryOpts{
		TaskID: *taskID, Limit: *limit,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "outcomes worked: %v\n", err)
		return exitFail
	}
	if rows == nil {
		rows = []domain.WorkedApproachRow{}
	}
	if err := json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "worked": rows, "count": len(rows),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "outcomes worked: %v\n", err)
		return exitFail
	}
	return exitOK
}
