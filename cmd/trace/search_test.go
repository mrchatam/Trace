package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func TestCLISearchUsesFTS(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc := domain.New(st)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "Task zephyrsearchtoken"})
	if err != nil {
		st.Close()
		t.Fatalf("CreateTask: %v", err)
	}
	chg, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "regression zephyrchangef token",
		Paths:  []domain.ChangePathInput{{Path: "internal/foo.go", Status: "M"}},
	})
	if err != nil {
		st.Close()
		t.Fatalf("CreateChange: %v", err)
	}
	st.Close()

	out := captureStdout(t, func() int {
		return run([]string{"-C", dir, "search", "zephyrchangef"})
	})
	var resp struct {
		OK    bool            `json:"ok"`
		Hits  []retrieval.Hit `json:"hits"`
		Count int             `json:"count"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &resp); err != nil {
		t.Fatalf("json: %v body=%q", err, out)
	}
	if !resp.OK || resp.Count < 1 {
		t.Fatalf("expected hits: %+v", resp)
	}
	var foundChange bool
	for _, h := range resp.Hits {
		if h.ReasonCode != retrieval.ReasonFTSMatch {
			t.Fatalf("expected fts_match, got %q in %+v", h.ReasonCode, h)
		}
		if h.EntityType == "change" && h.EntityID == chg.ID {
			foundChange = true
		}
	}
	if !foundChange {
		t.Fatalf("expected change entity in hits: %+v", resp.Hits)
	}

	emptyOut := captureStdout(t, func() int {
		return run([]string{"-C", dir, "search", "   "})
	})
	var emptyResp struct {
		OK    bool `json:"ok"`
		Count int  `json:"count"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(emptyOut)), &emptyResp); err != nil {
		t.Fatal(err)
	}
	if !emptyResp.OK || emptyResp.Count != 0 {
		t.Fatalf("empty query: %+v", emptyResp)
	}
}

func TestCLISearchHelp(t *testing.T) {
	out := captureStdout(t, func() int { return run([]string{"help"}) })
	if !strings.Contains(out, "search <query>") {
		t.Fatalf("help missing search: %q", out)
	}
}
