package store

import (
	"strings"
	"testing"
)

func TestOutcomeStoreRoundtrip(t *testing.T) {
	s, _ := openTempStore(t)

	b, err := s.UpsertBaseline(Baseline{
		GitCommit:  "abcdef0",
		ScoresJSON: `{"correctness":0.9}`,
		Label:      "B100",
	})
	if err != nil {
		t.Fatalf("UpsertBaseline: %v", err)
	}
	if b.ID == "" || b.CreatedAt == "" || b.UpdatedAt == "" {
		t.Fatalf("expected ids/timestamps: %+v", b)
	}
	gotB, err := s.GetBaseline(b.ID)
	if err != nil || gotB.GitCommit != "abcdef0" || gotB.Label != "B100" {
		t.Fatalf("GetBaseline: %+v err=%v", gotB, err)
	}

	testRow, err := s.UpsertOutcomeResult(OutcomeResult{
		TaskID:     "task-1",
		Kind:       OutcomeKindTest,
		TestName:   "TestFoo",
		TestStatus: TestStatusPass,
		Summary:    "ok",
		Actor:      "agent",
	})
	if err != nil {
		t.Fatalf("UpsertOutcomeResult test: %v", err)
	}
	if testRow.Kind != OutcomeKindTest || testRow.ScoresJSON != "{}" || testRow.ComparisonJSON != "{}" {
		t.Fatalf("test kind hygiene: %+v", testRow)
	}

	got, err := s.GetOutcomeResult(testRow.ID)
	if err != nil || got.TestName != "TestFoo" {
		t.Fatalf("GetOutcomeResult: %+v err=%v", got, err)
	}
	listed, err := s.ListOutcomeResultsByTaskID("task-1")
	if err != nil || len(listed) != 1 || listed[0].ID != testRow.ID {
		t.Fatalf("ListOutcomeResultsByTaskID: %+v err=%v", listed, err)
	}
	byKind, err := s.ListOutcomeResultsByTaskKind("task-1", OutcomeKindTest)
	if err != nil || len(byKind) != 1 {
		t.Fatalf("ListOutcomeResultsByTaskKind: %+v err=%v", byKind, err)
	}

	if _, err := s.UpsertOutcomeResult(OutcomeResult{Kind: OutcomeKindTest, TestName: "x", TestStatus: TestStatusPass}); err == nil {
		t.Fatal("expected task_id required")
	}
	if _, err := s.GetBaseline(""); err == nil {
		t.Fatal("expected baseline id required")
	}
	if _, err := s.GetBaseline("missing"); err == nil {
		t.Fatal("expected missing baseline")
	}
	if _, err := s.GetOutcomeResult(""); err == nil {
		t.Fatal("expected outcome id required")
	}

	signal, err := s.HasImplementationSignal("task-1")
	if err != nil || signal {
		t.Fatalf("no changes → no signal: %v %v", signal, err)
	}
}

func TestOutcomeResultsSchemaNoBlobColumns(t *testing.T) {
	s, _ := openTempStore(t)

	for _, table := range []string{"outcome_results", "baselines"} {
		var name string
		if err := s.db.QueryRow(`SELECT name FROM sqlite_master WHERE name=?`, table).Scan(&name); err != nil {
			t.Fatalf("missing table %s: %v", table, err)
		}
		cols, err := tableColumns(s.conn, table)
		if err != nil {
			t.Fatal(err)
		}
		if len(cols) == 0 {
			t.Fatalf("%s has no columns", table)
		}
		for _, c := range cols {
			lower := strings.ToLower(c)
			for _, bad := range []string{
				"blob", "patch", "diff", "content", "file_body", "source_text",
				"log", "stdout", "stderr", "raw_output", "full_log",
			} {
				if strings.Contains(lower, bad) && lower != "confidence" {
					t.Errorf("forbidden column %s.%s", table, c)
				}
			}
		}
	}

	baseCols, err := tableColumns(s.conn, "baselines")
	if err != nil {
		t.Fatal(err)
	}
	wantBase := map[string]bool{
		"id": true, "git_commit": true, "scores_json": true,
		"label": true, "source_type": true, "created_at": true, "updated_at": true,
		"status": true, "supersedes_id": true,
	}
	if len(baseCols) != len(wantBase) {
		t.Fatalf("baselines columns=%v want exactly oid+scores+label+meta+promotion", baseCols)
	}
	for _, c := range baseCols {
		if !wantBase[c] {
			t.Errorf("unexpected baselines column %q", c)
		}
	}

	bad, where, err := s.HasBlobLikeColumns()
	if err != nil {
		t.Fatalf("HasBlobLikeColumns: %v", err)
	}
	if bad {
		t.Fatalf("HasBlobLikeColumns true at %s", where)
	}
}

func TestUnknownOutcomeKindFailClosed(t *testing.T) {
	s, _ := openTempStore(t)

	_, err := s.UpsertOutcomeResult(OutcomeResult{
		TaskID: "task-1",
		Kind:   "experiment",
	})
	if err == nil {
		t.Fatal("expected unknown kind to fail closed")
	}

	_, err = s.UpsertOutcomeResult(OutcomeResult{
		TaskID:     "task-1",
		Kind:       OutcomeKindTest,
		TestName:   "TestFoo",
		TestStatus: TestStatusPass,
		GoalID:     "goal-1",
	})
	if err == nil {
		t.Fatal("expected kind=test with goal_id to fail closed")
	}

	_, err = s.db.Exec(`
		INSERT INTO outcome_results(
			id, task_id, kind, test_name, test_status, goal_id, verification_status,
			baseline_id, scores_json, comparison_json, summary, actor, source_type,
			confidence, created_at, updated_at
		) VALUES ('o1','t1','bogus','','','','','','{}','{}','','','',0,'now','now')
	`)
	if err == nil {
		t.Fatal("expected CHECK(kind) to reject unknown kind")
	}
}

func TestBaselineStoreRoundtrip(t *testing.T) {
	s, _ := openTempStore(t)
	if _, err := s.UpsertBaseline(Baseline{ScoresJSON: `{"a":1}`}); err == nil {
		t.Fatal("expected git_commit required")
	}
	b, err := s.UpsertBaseline(Baseline{
		GitCommit:  "abc1234",
		ScoresJSON: `{"correctness":1}`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(b.ScoresJSON), "blob") {
		t.Fatalf("scores look like a blob: %s", b.ScoresJSON)
	}
	if b.Status != BaselineStatusActive || b.SupersedesID != "" {
		t.Fatalf("defaults: status=%q supersedes_id=%q", b.Status, b.SupersedesID)
	}
	b2, err := s.UpsertBaseline(Baseline{
		ID:           "b-superseded",
		GitCommit:    "abc1234",
		ScoresJSON:   `{"correctness":0.9}`,
		Label:        "chain",
		Status:       BaselineStatusSuperseded,
		SupersedesID: "prior-id",
	})
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.GetBaseline(b2.ID)
	if err != nil || got.Status != BaselineStatusSuperseded || got.SupersedesID != "prior-id" {
		t.Fatalf("promotion columns roundtrip: %+v err=%v", got, err)
	}
	_, err = s.GetActiveBaselineByCommitLabel("abc1234", "chain")
	if err == nil {
		t.Fatal("superseded must not be returned as active")
	}
	if err := s.SetBaselinePromotion(b.ID, BaselineStatusSuperseded, ""); err != nil {
		t.Fatal(err)
	}
}
