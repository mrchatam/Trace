package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

type patternsRefreshResponse struct {
	OK              bool `json:"ok"`
	PatternsUpdated int  `json:"patterns_updated"`
}

func cmdPatterns(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace patterns refresh|list\n")
		return exitUsage
	}
	switch args[0] {
	case "refresh":
		return cmdPatternsRefresh(root, args[1:])
	case "list":
		return cmdPatternsList(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown patterns subcommand: %s\n", args[0])
		return exitUsage
	}
}

func cmdPatternsRefresh(root string, args []string) int {
	fs := flag.NewFlagSet("patterns refresh", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace patterns refresh\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "patterns refresh: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "patterns refresh: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "patterns", "patterns refresh"); code != exitOK {
		return code
	}

	n, err := svc.RefreshChangePatterns(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "patterns refresh: %v\n", err)
		return exitFail
	}
	resp := patternsRefreshResponse{OK: true, PatternsUpdated: n}
	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		fmt.Fprintf(os.Stderr, "patterns refresh: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdPatternsList(root string, args []string) int {
	fs := flag.NewFlagSet("patterns list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	limit := fs.Int("limit", 0, "max rows (default 32, cap 64)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"limit": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace patterns list [--limit N]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "patterns list: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "patterns list: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "patterns", "patterns list"); code != exitOK {
		return code
	}

	rows, err := svc.ListChangePatterns(context.Background(), *limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "patterns list: %v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(rows); err != nil {
		fmt.Fprintf(os.Stderr, "patterns list: %v\n", err)
		return exitFail
	}
	return exitOK
}
