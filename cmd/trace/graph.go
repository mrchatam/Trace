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

func cmdGraph(root string, args []string) int {
	if len(args) < 1 {
		printGraphHelp(os.Stderr)
		return exitUsage
	}
	switch args[0] {
	case "infer":
		return cmdGraphInfer(root, args[1:])
	case "help", "-h", "--help":
		printGraphHelp(os.Stdout)
		return exitOK
	default:
		fmt.Fprintf(os.Stderr, "graph: unknown subcommand %q\n", args[0])
		printGraphHelp(os.Stderr)
		return exitUsage
	}
}

func printGraphHelp(w *os.File) {
	fmt.Fprint(w, `trace graph — opt-in graph utilities

Subcommands:
  infer [--dry-run]
        Run MVP scope-membership inference (R-PATH, R-TITLE, R-PLAN).
        Writes scope_member edges with source_type=INFERRED (confidence 0.4).
        Opt-in only — does not run on trace index / watch / daemon.
        --dry-run: propose candidates with zero writes; JSON summary on stdout.
`)
}

func cmdGraphInfer(root string, args []string) int {
	fs := flag.NewFlagSet("graph infer", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	dryRun := fs.Bool("dry-run", false, "propose candidates; write nothing")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "graph infer: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "graph infer: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "graph", "graph infer"); code != exitOK {
		return code
	}

	rep, err := svc.InferScopes(context.Background(), domain.InferOptions{DryRun: *dryRun})
	if err != nil {
		fmt.Fprintf(os.Stderr, "graph infer: %v\n", err)
		return exitFail
	}
	out := map[string]any{
		"ok":               true,
		"dry_run":          *dryRun,
		"inserted":         rep.Inserted,
		"skipped_existing": rep.SkippedExisting,
		"skipped_conflict": rep.SkippedConflict,
		"skipped_explicit": rep.SkippedExplicit,
		"candidates":       rep.Candidates,
		"skips":            rep.Skips,
	}
	if err := json.NewEncoder(os.Stdout).Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "graph infer: %v\n", err)
		return exitFail
	}
	return exitOK
}
