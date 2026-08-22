package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrchatam/Trace/internal/eval"
)

func cmdEval(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace eval rules | trace eval results --task <id>\n")
		return exitUsage
	}
	switch args[0] {
	case "rules":
		return cmdEvalRules(root, args[1:])
	case "results":
		return cmdEvalResults(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown eval subcommand: %s\n", args[0])
		return exitUsage
	}
}

func cmdEvalRules(root string, args []string) int {
	fs := flag.NewFlagSet("eval rules", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace eval rules\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	if code := failCLIDenied(svc, "eval", "eval rules"); code != exitOK {
		return code
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "eval rules: %v\n", err)
		return exitFail
	}

	load, err := eval.LoadRules(context.Background(), abs, st)
	if err != nil {
		fmt.Fprintf(os.Stderr, "eval rules: %v\n", err)
		return exitFail
	}
	payload := map[string]any{
		"path":       load.Path,
		"loaded":     load.Loaded,
		"mechanisms": load.Mechanisms,
		"invariants": load.Invariants,
		"cached_at":  load.CachedAt,
	}
	if err := json.NewEncoder(os.Stdout).Encode(payload); err != nil {
		fmt.Fprintf(os.Stderr, "eval rules: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdEvalResults(root string, args []string) int {
	fs := flag.NewFlagSet("eval results", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "task id")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace eval results --task <id>\n")
		return exitUsage
	}
	if strings.TrimSpace(*taskID) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace eval results --task <id>\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	if code := failCLIDenied(svc, "eval", "eval results"); code != exitOK {
		return code
	}

	rows, err := eval.ListResults(context.Background(), svc, *taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "eval results: %v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(rows); err != nil {
		fmt.Fprintf(os.Stderr, "eval results: %v\n", err)
		return exitFail
	}
	return exitOK
}
