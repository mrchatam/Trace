package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/indexwalk"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

type indexStatusJSON struct {
	Head               string          `json:"head"`
	LastIndexedCommit  string          `json:"last_indexed_commit"`
	Stale              bool            `json:"stale"`
	HookInstalled      bool            `json:"hook_installed"`
	SupportedLanguages []string        `json:"supported_languages"`
	Drift              *indexDriftJSON `json:"drift"`
}

// indexDriftJSON is the wire shape of indexwalk.Drift (same fields, JSON
// tags). The walk-vs-store comparison itself lives in internal/indexwalk so
// the loop gate evaluates the identical definition.
type indexDriftJSON = indexwalk.Drift

func computeIndexDrift(abs string, useGitIgnore bool, st *store.Store) (*indexDriftJSON, error) {
	// useGitIgnore is honored for CLI parity with the index walk; the loop
	// gate passes true unconditionally (GitIgnored self-disables without git).
	_ = useGitIgnore
	return indexwalk.ComputeDrift(abs, st)
}

func cmdIndexStatus(root string, args []string) int {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace index status\n")
		return exitUsage
	}
	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "index status: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "index status: %v\n", err)
		return exitFail
	}
	defer st.Close()

	state, err := st.GetGraphSyncState()
	if err != nil {
		fmt.Fprintf(os.Stderr, "index status: %v\n", err)
		return exitFail
	}

	out := indexStatusJSON{
		HookInstalled:      state.HookInstalled,
		SupportedLanguages: analyzers.SupportedLanguages(),
	}
	var repo vcs.Repository
	if r, rerr := tryOpenGit(abs, st); rerr == nil {
		repo = r
		defer repo.Close()
		head, herr := repo.Head(context.Background())
		if herr == nil {
			out.Head = head
			out.LastIndexedCommit = state.LastIndexedCommit
			out.Stale = head != state.LastIndexedCommit
		}
	}

	drift, derr := computeIndexDrift(abs, repo != nil, st)
	if derr != nil {
		fmt.Fprintf(os.Stderr, "index status: drift: %v\n", derr)
		return exitFail
	}
	out.Drift = drift
	if drift.NotIndexed > 0 {
		fmt.Fprintf(os.Stderr, "index status: drift: %d of %d walkable files are not in the store (e.g. %s) — run `trace index` to add them\n",
			drift.NotIndexed, drift.Walkable, shortPathSample(drift.NotIndexedPaths, drift.NotIndexed))
	}
	if drift.NotOnDisk > 0 {
		fmt.Fprintf(os.Stderr, "index status: %d indexed files are gone from disk (or no longer indexable); a full `trace index` removes them (e.g. %s)\n",
			drift.NotOnDisk, shortPathSample(drift.NotOnDiskPaths, drift.NotOnDisk))
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "index status: %v\n", err)
		return exitFail
	}
	return exitOK
}

func updateGraphSyncWatermark(ctx context.Context, st *store.Store, repo vcs.Repository) error {
	head, err := repo.Head(ctx)
	if err != nil {
		return err
	}
	state, err := st.GetGraphSyncState()
	if err != nil {
		return err
	}
	state.LastIndexedCommit = head
	state.LastIndexedAt = time.Now().UTC().Format(time.RFC3339)
	return st.UpsertGraphSyncState(state)
}

func setHookInstalledFlag(st *store.Store, installed bool) error {
	state, err := st.GetGraphSyncState()
	if err != nil {
		return err
	}
	state.HookInstalled = installed
	return st.UpsertGraphSyncState(state)
}

// shortPathSample quotes up to indexwalk.MaxPathsSample paths for a stderr
// line, with an "…and N more" tail when the underlying list was truncated.
func shortPathSample(paths []string, total int) string {
	named := paths
	if len(named) > indexwalk.MaxPathsSample {
		named = named[:indexwalk.MaxPathsSample]
	}
	quoted := make([]string, 0, len(named)+1)
	for _, p := range named {
		quoted = append(quoted, strconv.Quote(p))
	}
	if more := total - len(named); more > 0 {
		quoted = append(quoted, fmt.Sprintf("…and %d more", more))
	}
	return strings.Join(quoted, ", ")
}
