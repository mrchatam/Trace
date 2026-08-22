package store

import (
	"strings"
	"testing"
)

func TestFTSFindsEntityTitleAndPathSymbol(t *testing.T) {
	s, _ := openTempStore(t)

	goal, err := s.UpsertGoal(Goal{
		Title: "Ship zephyrunique retrieval",
		Body:  "lexical index for agents",
	})
	if err != nil {
		t.Fatalf("UpsertGoal: %v", err)
	}

	hits, err := s.SearchFTS("zephyrunique", 10)
	if err != nil {
		t.Fatalf("SearchFTS: %v", err)
	}
	found := false
	for _, h := range hits {
		if h.EntityType == "goal" && h.EntityID == goal.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected goal in FTS hits, got %+v", hits)
	}

	_, err = s.UpsertFile("pkg/zephyrpath.go", "abc123", nil)
	if err != nil {
		t.Fatalf("UpsertFile: %v", err)
	}
	if err := s.ReplaceFileSymbols("pkg/zephyrpath.go", []Symbol{
		{Name: "ZephyrSymbol", Kind: "function", StartLine: 1, EndLine: 3},
	}); err != nil {
		t.Fatalf("ReplaceFileSymbols: %v", err)
	}

	pathHits, err := s.SearchFTS("zephyrpath", 10)
	if err != nil {
		t.Fatalf("SearchFTS path: %v", err)
	}
	foundPath := false
	for _, h := range pathHits {
		if h.EntityType == "file" && strings.Contains(h.Path, "zephyrpath") {
			foundPath = true
			break
		}
	}
	if !foundPath {
		t.Fatalf("expected file path in FTS hits, got %+v", pathHits)
	}

	symHits, err := s.SearchFTS("ZephyrSymbol", 10)
	if err != nil {
		t.Fatalf("SearchFTS symbol: %v", err)
	}
	foundSym := false
	for _, h := range symHits {
		if h.EntityType == "symbol" && h.SymbolName == "ZephyrSymbol" {
			foundSym = true
			break
		}
	}
	if !foundSym {
		t.Fatalf("expected symbol in FTS hits, got %+v", symHits)
	}

	// RebuildFTS must still find the same tokens (no source BLOBs).
	if err := s.RebuildFTS(); err != nil {
		t.Fatalf("RebuildFTS: %v", err)
	}
	after, err := s.SearchFTS("zephyrunique", 10)
	if err != nil || len(after) == 0 {
		t.Fatalf("after rebuild: %v %+v", err, after)
	}

	// Cap proof: limit never exceeds 64.
	many, err := s.SearchFTS("zephyrunique", 1000)
	if err != nil {
		t.Fatalf("SearchFTS cap: %v", err)
	}
	if len(many) > 64 {
		t.Fatalf("limit hard cap exceeded: %d", len(many))
	}
}

func TestFTSBackfillOnOpenWhenIndexEmpty(t *testing.T) {
	s, root := openTempStore(t)

	goal, err := s.UpsertGoal(Goal{Title: "preexisting backfilltokenxyz"})
	if err != nil {
		t.Fatalf("UpsertGoal: %v", err)
	}
	// Simulate post-migrate empty FTS (004 applied onto a DB that already had rows).
	if _, err := s.db.Exec(`DELETE FROM fts_docs`); err != nil {
		t.Fatalf("clear fts: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	s2, err := Open(root)
	if err != nil {
		t.Fatalf("re-Open: %v", err)
	}
	t.Cleanup(func() { _ = s2.Close() })

	hits, err := s2.SearchFTS("backfilltokenxyz", 10)
	if err != nil {
		t.Fatalf("SearchFTS after backfill: %v", err)
	}
	found := false
	for _, h := range hits {
		if h.EntityType == "goal" && h.EntityID == goal.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected Open backfill to reindex goal, got %+v", hits)
	}
}

func TestListTasksByGoalID(t *testing.T) {
	s, _ := openTempStore(t)
	g, err := s.UpsertGoal(Goal{Title: "G"})
	if err != nil {
		t.Fatalf("goal: %v", err)
	}
	gid := g.ID
	t1, err := s.UpsertTask(Task{Title: "T1", GoalID: &gid})
	if err != nil {
		t.Fatalf("task1: %v", err)
	}
	_, err = s.UpsertTask(Task{Title: "orphan"})
	if err != nil {
		t.Fatalf("orphan: %v", err)
	}
	listed, err := s.ListTasksByGoalID(g.ID)
	if err != nil {
		t.Fatalf("ListTasksByGoalID: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != t1.ID {
		t.Fatalf("listed: %+v", listed)
	}
}

func TestListTasks(t *testing.T) {
	s, _ := openTempStore(t)
	empty, err := s.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks empty: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("empty ListTasks: %+v", empty)
	}
	g, err := s.UpsertGoal(Goal{Title: "G"})
	if err != nil {
		t.Fatalf("goal: %v", err)
	}
	gid := g.ID
	t1, err := s.UpsertTask(Task{Title: "T1", GoalID: &gid})
	if err != nil {
		t.Fatalf("task1: %v", err)
	}
	t2, err := s.UpsertTask(Task{Title: "orphan"})
	if err != nil {
		t.Fatalf("orphan: %v", err)
	}
	listed, err := s.ListTasks()
	if err != nil {
		t.Fatalf("ListTasks: %v", err)
	}
	if len(listed) != 2 {
		t.Fatalf("ListTasks len: %+v", listed)
	}
	ids := map[string]bool{listed[0].ID: true, listed[1].ID: true}
	if !ids[t1.ID] || !ids[t2.ID] {
		t.Fatalf("ListTasks ids: %+v want %s and %s", listed, t1.ID, t2.ID)
	}
	// Same-second inserts: ORDER BY created_at ASC, id ASC
	if listed[0].CreatedAt == listed[1].CreatedAt && listed[0].ID > listed[1].ID {
		t.Fatalf("ListTasks id order: %+v", listed)
	}
}

func TestSanitizeFTSQueryPunctuationClass(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"GET /notes", `"GET" AND "notes"`},
		{"GET /notes/search", `"GET" AND "notes" AND "search"`},
		{"zephyrunique", `"zephyrunique"`},
		{"foo AND bar", `"foo" AND "AND" AND "bar"`},
		{"title:secret", `"title" AND "secret"`},
		{"NEAR foo", `"NEAR" AND "foo"`},
		{"work_state", `"work" AND "state"`},
		{"  GET   /notes  ", `"GET" AND "notes"`},
		{"///", ""},
		{"C++", `"C"`},
	}
	for _, tc := range cases {
		got := sanitizeFTSQuery(tc.in)
		if got != tc.want {
			t.Errorf("sanitizeFTSQuery(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestSyncEntityFTSChange(t *testing.T) {
	s, _ := openTempStore(t)

	task, err := s.UpsertTask(Task{Title: "FTS change task"})
	if err != nil {
		t.Fatalf("UpsertTask: %v", err)
	}
	chg, err := s.UpsertChange(Change{
		TaskID: task.ID,
		Reason: "zephyrftschange token on upsert",
		Status: ChangeStatusRecorded,
	})
	if err != nil {
		t.Fatalf("UpsertChange: %v", err)
	}

	hits, err := s.SearchFTS("zephyrftschange", 10)
	if err != nil {
		t.Fatalf("SearchFTS: %v", err)
	}
	found := false
	for _, h := range hits {
		if h.EntityType == "change" && h.EntityID == chg.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected change in FTS after upsert, got %+v", hits)
	}
}

func TestSyncEntityFTSUncertainty(t *testing.T) {
	s, _ := openTempStore(t)

	u, err := s.UpsertUncertainty(Uncertainty{
		Title: "zephyrftsuncertainty token",
		Body:  "indexed on upsert",
	})
	if err != nil {
		t.Fatalf("UpsertUncertainty: %v", err)
	}

	hits, err := s.SearchFTS("zephyrftsuncertainty", 10)
	if err != nil {
		t.Fatalf("SearchFTS: %v", err)
	}
	found := false
	for _, h := range hits {
		if h.EntityType == "uncertainty" && h.EntityID == u.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected uncertainty in FTS after upsert, got %+v", hits)
	}
}

func TestRebuildFTSIncludesP20Types(t *testing.T) {
	s, _ := openTempStore(t)

	u, err := s.UpsertUncertainty(Uncertainty{
		Title: "zephyrrebuilduncertainty token",
		Body:  "must survive rebuild",
	})
	if err != nil {
		t.Fatalf("UpsertUncertainty: %v", err)
	}

	if _, err := s.db.Exec(`DELETE FROM fts_docs`); err != nil {
		t.Fatalf("clear fts: %v", err)
	}
	if err := s.RebuildFTS(); err != nil {
		t.Fatalf("RebuildFTS: %v", err)
	}

	hits, err := s.SearchFTS("zephyrrebuilduncertainty", 10)
	if err != nil {
		t.Fatalf("SearchFTS: %v", err)
	}
	found := false
	for _, h := range hits {
		if h.EntityType == "uncertainty" && h.EntityID == u.ID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected uncertainty after RebuildFTS, got %+v", hits)
	}
}

func TestSearchFTSSlashInQuery(t *testing.T) {
	for _, title := range []string{"GET /notes", "GET /notes/search"} {
		t.Run(title, func(t *testing.T) {
			s, _ := openTempStore(t)
			task, err := s.UpsertTask(Task{Title: title})
			if err != nil {
				t.Fatalf("UpsertTask: %v", err)
			}
			hits, err := s.SearchFTS(title, 10)
			if err != nil {
				t.Fatalf("SearchFTS(%q): %v", title, err)
			}
			found := false
			for _, h := range hits {
				if h.EntityType == "task" && h.EntityID == task.ID {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("expected task %s in SearchFTS hits, got %+v", task.ID, hits)
			}
		})
	}
}
