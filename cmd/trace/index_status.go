package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

type indexStatusJSON struct {
	Head               string   `json:"head"`
	LastIndexedCommit  string   `json:"last_indexed_commit"`
	Stale              bool     `json:"stale"`
	HookInstalled      bool     `json:"hook_installed"`
	SupportedLanguages []string `json:"supported_languages"`
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
	if repo, rerr := tryOpenGit(abs, st); rerr == nil {
		defer repo.Close()
		head, herr := repo.Head(context.Background())
		if herr == nil {
			out.Head = head
			out.LastIndexedCommit = state.LastIndexedCommit
			out.Stale = head != state.LastIndexedCommit
		}
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
