package vcs_test

import (
	"context"
	"testing"

	"github.com/mrchatam/Trace/internal/vcs"
)

func TestFakeRepositoryShowFileAndHistory(t *testing.T) {
	ctx := context.Background()
	fake := &vcs.Fake{
		IsGit:   true,
		HeadOID: "aaa",
		Files: map[string][]byte{
			"aaa:readme.md": []byte("hello"),
		},
		PathHistory: map[string][]vcs.CommitMeta{
			"readme.md": {
				{OID: "aaa", Subject: "add readme", CommittedAt: "2026-01-02T00:00:00Z"},
				{OID: "bbb", Subject: "init", CommittedAt: "2026-01-01T00:00:00Z"},
			},
		},
		PathsByCommit: map[string][]vcs.PathChange{
			"aaa": {{Path: "readme.md", Status: "M"}},
		},
		Commits: []vcs.CommitMeta{
			{OID: "bbb", Subject: "init"},
			{OID: "aaa", Subject: "add readme"},
		},
	}

	var repo vcs.Repository = fake

	ok, err := repo.IsRepo(ctx)
	if err != nil || !ok {
		t.Fatalf("IsRepo: ok=%v err=%v", ok, err)
	}

	head, err := repo.Head(ctx)
	if err != nil || head != "aaa" {
		t.Fatalf("Head: got %q err=%v", head, err)
	}

	body, err := repo.ShowFile(ctx, "aaa", "readme.md")
	if err != nil {
		t.Fatalf("ShowFile: %v", err)
	}
	if string(body) != "hello" {
		t.Fatalf("ShowFile body: %q", body)
	}

	hist, err := repo.History(ctx, "readme.md", 1)
	if err != nil {
		t.Fatalf("History: %v", err)
	}
	if len(hist) != 1 || hist[0].OID != "aaa" {
		t.Fatalf("History: %+v", hist)
	}

	last, err := repo.LastChanged(ctx, "readme.md")
	if err != nil || last.OID != "aaa" {
		t.Fatalf("LastChanged: %+v err=%v", last, err)
	}

	changes, err := repo.Changes(ctx, "aaa")
	if err != nil || len(changes) != 1 || changes[0].Status != "M" {
		t.Fatalf("Changes: %+v err=%v", changes, err)
	}

	between, err := repo.CommitsBetween(ctx, "bbb", "aaa")
	if err != nil || len(between) != 1 || between[0].OID != "aaa" {
		t.Fatalf("CommitsBetween: %+v err=%v", between, err)
	}

	res, err := repo.Refresh(ctx)
	if err != nil || res.NewCommits != 0 {
		t.Fatalf("Refresh: %+v err=%v", res, err)
	}

	if err := repo.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if !fake.Closed {
		t.Fatal("expected Fake.Closed")
	}
}
