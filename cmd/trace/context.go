package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func cmdContext(root string, args []string) int {
	fs := flag.NewFlagSet("context", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	depth := fs.Int("depth", 1, "expand depth (1=TaskContext, 2=ExpandContext; max 2)")
	maxLayer := fs.Int("max-layer", 1, "progressive layer ceiling 1=L0–L1 (default), 2|3=opt-in deeper layers")
	format := fs.String("format", "json", "json|markdown|both")
	includeWhy := fs.Bool("include-why", false, "include why_trace")
	query := fs.String("query", "", "optional agent query merged into context (task moat preserved)")
	// Allow flags before or after the positional task-id (stdlib flag stops at first non-flag).
	if err := fs.Parse(flagsFirst(args, map[string]bool{
		"depth": true, "max-layer": true, "format": true, "include-why": false, "query": true,
	})); err != nil {
		return exitUsage
	}
	pos := fs.Args()
	if len(pos) != 1 {
		fmt.Fprintf(os.Stderr, "usage: trace context <task-id> [--depth 1|2] [--max-layer 1|2|3] [--format json|markdown|both] [--include-why] [--query <text>]\n")
		return exitUsage
	}
	taskID := pos[0]
	if *depth < 1 || *depth > 2 {
		fmt.Fprintf(os.Stderr, "context: --depth must be 1 or 2\n")
		return exitUsage
	}
	if *maxLayer < 1 || *maxLayer > 3 {
		fmt.Fprintf(os.Stderr, "context: --max-layer must be 1, 2, or 3\n")
		return exitUsage
	}
	switch *format {
	case "json", "markdown", "both":
	default:
		fmt.Fprintf(os.Stderr, "context: --format must be json|markdown|both\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "context: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "context: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "context", "context"); code != exitOK {
		return code
	}

	eng := retrieval.New(st)
	if repo, rerr := tryOpenGit(abs, st); rerr == nil {
		defer repo.Close()
		eng = eng.WithVCS(repo)
	}
	comp := compiler.New(st).WithRetrieval(eng)
	opts := compiler.ContextOptions{
		MaxLayer:        *maxLayer,
		IncludeWhy:      *includeWhy,
		IncludeMarkdown: *format == "markdown" || *format == "both",
		Query:           *query,
	}

	ctx := context.Background()
	var pkt compiler.Packet
	if *depth == 1 {
		pkt, err = comp.TaskContext(ctx, taskID, opts)
	} else {
		pkt, err = comp.ExpandContext(ctx, taskID, *depth, opts)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "context: %v\n", err)
		return exitFail
	}

	switch *format {
	case "json":
		b, err := pkt.JSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "context: %v\n", err)
			return exitFail
		}
		fmt.Println(string(b))
	case "markdown":
		md := pkt.Markdown()
		fmt.Print(md)
		if len(md) == 0 || md[len(md)-1] != '\n' {
			fmt.Println()
		}
	case "both":
		b, err := pkt.JSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "context: %v\n", err)
			return exitFail
		}
		fmt.Println(string(b))
		fmt.Println("---")
		md := pkt.Markdown()
		fmt.Print(md)
		if len(md) == 0 || md[len(md)-1] != '\n' {
			fmt.Println()
		}
	}
	return exitOK
}
