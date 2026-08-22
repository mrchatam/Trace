package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

type searchResponse struct {
	OK    bool            `json:"ok"`
	Hits  []retrieval.Hit `json:"hits"`
	Count int             `json:"count"`
}

func cmdSearch(root string, args []string) int {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	limit := fs.Int("limit", 0, "max hits (default 32, cap 64)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"limit": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: trace search <query> [--limit N]\n")
		return exitUsage
	}
	query := strings.TrimSpace(fs.Arg(0))

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "search: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "search: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "search", "search"); code != exitOK {
		return code
	}

	resp := searchResponse{OK: true, Hits: []retrieval.Hit{}, Count: 0}
	if query != "" {
		hits, err := retrieval.New(st).Search(context.Background(), query, retrieval.SearchOptions{Limit: *limit})
		if err != nil {
			fmt.Fprintf(os.Stderr, "search: %v\n", err)
			return exitFail
		}
		if hits == nil {
			hits = []retrieval.Hit{}
		}
		resp.Hits = hits
		resp.Count = len(hits)
	}

	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		fmt.Fprintf(os.Stderr, "search: %v\n", err)
		return exitFail
	}
	return exitOK
}
