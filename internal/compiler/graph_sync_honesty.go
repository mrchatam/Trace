package compiler

import (
	"context"

	"github.com/mrchatam/Trace/internal/gitcli"
	"github.com/mrchatam/Trace/internal/store"
)

const graphSyncHonestyNotice = "graph indexed at an older commit — run trace index"

// buildGraphSyncHonesty reports when git HEAD differs from last_indexed_commit.
// Omit when not a git repo, HEAD unavailable, or watermark matches HEAD.
func buildGraphSyncHonesty(st *store.Store) *GraphSyncHonesty {
	if st == nil {
		return nil
	}
	repo, err := gitcli.OpenWithStore(st.ProjectRoot(), st)
	if err != nil {
		return nil
	}
	defer repo.Close()

	head, err := repo.Head(context.Background())
	if err != nil || head == "" {
		return nil
	}

	state, err := st.GetGraphSyncState()
	if err != nil {
		return nil
	}
	if state.LastIndexedCommit == head {
		return nil
	}

	return &GraphSyncHonesty{
		StaleCommit:       true,
		Head:              head,
		LastIndexedCommit: state.LastIndexedCommit,
		Notice:            graphSyncHonestyNotice,
	}
}
