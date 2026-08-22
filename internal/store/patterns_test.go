package store

import (
	"testing"
)

func TestReplaceChangePatternsAndList(t *testing.T) {
	s, _ := openTempStore(t)
	rows := []ChangePattern{
		{ChangeKind: "seg:internal", OutcomeKind: "improvement", CountPositive: 2, LastSeen: "2026-01-02T00:00:00Z"},
		{ChangeKind: "seg:internal", OutcomeKind: "regression", CountNegative: 1, LastSeen: "2026-01-01T00:00:00Z"},
	}
	if err := s.ReplaceChangePatterns(rows); err != nil {
		t.Fatal(err)
	}
	all, err := s.ListChangePatterns(64)
	if err != nil || len(all) != 2 {
		t.Fatalf("ListChangePatterns: %+v err=%v", all, err)
	}
	byKind, err := s.ListChangePatternsByKind("seg:internal", 64)
	if err != nil || len(byKind) != 2 {
		t.Fatalf("ListChangePatternsByKind: %+v err=%v", byKind, err)
	}
	if err := s.ReplaceChangePatterns(nil); err != nil {
		t.Fatal(err)
	}
	empty, err := s.ListChangePatterns(64)
	if err != nil || len(empty) != 0 {
		t.Fatalf("after clear: %+v err=%v", empty, err)
	}
}

func TestListChangesByPathPrefixExcludesSuperseded(t *testing.T) {
	s, _ := openTempStore(t)
	taskID := "task-1"
	open, err := s.UpsertChange(Change{
		TaskID: taskID, Reason: "open", Status: ChangeStatusOpen, SourceType: "USER_ASSERTED",
		CreatedAt: "2026-01-01T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.InsertChangePath(ChangePath{ChangeID: open.ID, Path: "internal/open.go"}); err != nil {
		t.Fatal(err)
	}
	sup, err := s.UpsertChange(Change{
		TaskID: taskID, Reason: "old", Status: ChangeStatusSuperseded, SourceType: "USER_ASSERTED",
		CreatedAt: "2026-01-02T00:00:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.InsertChangePath(ChangePath{ChangeID: sup.ID, Path: "internal/old.go"}); err != nil {
		t.Fatal(err)
	}

	got, err := s.ListChangesByPathPrefix("internal/", 64)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != open.ID {
		t.Fatalf("superseded excluded: %+v", got)
	}
}
