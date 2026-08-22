package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func cmdExplore(root string, args []string) int {
	fs := flag.NewFlagSet("explore", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	query := fs.String("query", "", "optional agent query merged into task context (task moat preserved)")
	limit := fs.Int("limit", 0, "FTS max hits (default 32, cap 64)")
	maxNodes := fs.Int("max-nodes", 0, "neighborhood max_nodes (default 100)")
	depth := fs.Int("depth", 0, "neighborhood depth (default 2)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{
		"query": true, "limit": true, "max-nodes": true, "depth": true,
	})); err != nil {
		return exitUsage
	}
	pos := fs.Args()
	if len(pos) != 1 {
		fmt.Fprintf(os.Stderr, "usage: trace explore <task-id> [--query <text>] [--limit N] [--max-nodes N] [--depth N]\n")
		return exitUsage
	}
	taskID := strings.TrimSpace(pos[0])
	if taskID == "" {
		fmt.Fprintf(os.Stderr, "explore: task-id is required\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "explore: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "explore: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "explore", "explore"); code != exitOK {
		return code
	}

	eng := retrieval.New(st)
	if repo, rerr := tryOpenGit(abs, st); rerr == nil {
		defer repo.Close()
		eng = eng.WithVCS(repo)
	}
	comp := compiler.New(st).WithRetrieval(eng)

	out, err := compiler.Explore(context.Background(), comp, eng, compiler.ExploreOpts{
		TaskID:   taskID,
		Query:    *query,
		Limit:    *limit,
		MaxNodes: *maxNodes,
		Depth:    *depth,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "explore: %v\n", err)
		return exitFail
	}

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "explore: %v\n", err)
		return exitFail
	}
	fmt.Println(string(b))
	return exitOK
}
