package store_test

import (
	"testing"

	"github.com/mrchatam/Trace/internal/store"
)

func TestListPromotionCandidateDiscoveriesSQLFilter(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })

	info, err := st.UpsertDiscovery(store.Discovery{
		ID: "d-info", Title: "info", Severity: store.SeverityINFO, Status: store.StatusActive,
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = info
	block1, err := st.UpsertDiscovery(store.Discovery{
		ID: "d-block-1", Title: "b1", Severity: store.SeverityBlocking, Status: store.StatusActive,
	})
	if err != nil {
		t.Fatal(err)
	}
	block2, err := st.UpsertDiscovery(store.Discovery{
		ID: "d-block-2", Title: "b2", Severity: store.SeverityBlocking, Status: store.StatusActive,
	})
	if err != nil {
		t.Fatal(err)
	}
	task, err := st.UpsertTask(store.Task{
		ID: "t1", Title: "linked", WorkState: store.WorkStatePending, Status: store.StatusActive,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.InsertLink(store.EntityLink{
		FromType: "discovery", FromID: block1.ID, Rel: "discovery_mentions_task",
		ToType: "task", ToID: task.ID,
	}); err != nil {
		t.Fatal(err)
	}

	got, err := st.ListPromotionCandidateDiscoveries(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != block2.ID {
		t.Fatalf("want only unlinked BLOCKING d-block-2, got %+v", got)
	}

	// limit+1 style: ask for 1 when 1 candidate → no extra
	got, err = st.ListPromotionCandidateDiscoveries(1)
	if err != nil || len(got) != 1 {
		t.Fatalf("limit 1: %+v err=%v", got, err)
	}
}

func TestListPromotionCandidateDiscoveriesDeadTaskLinkStillCandidate(t *testing.T) {
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.Close() })

	d, err := st.UpsertDiscovery(store.Discovery{
		ID: "d-dead", Title: "orphan link", Severity: store.SeverityBlocking, Status: store.StatusActive,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.InsertLink(store.EntityLink{
		FromType: "discovery", FromID: d.ID, Rel: "discovery_mentions_task",
		ToType: "task", ToID: "missing-task-id",
	}); err != nil {
		t.Fatal(err)
	}
	got, err := st.ListPromotionCandidateDiscoveries(5)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != d.ID {
		t.Fatalf("dead task link must still be a candidate: %+v", got)
	}
}
