package store

import (
	"strings"
	"testing"
)

func TestRegressionStoreRoundtrip(t *testing.T) {
	s, _ := openTempStore(t)

	row, err := s.UpsertRegression(Regression{
		TaskID:     "task-1",
		SourceKind: RegressionSourceEvaluation,
		SourceID:   "outcome-1",
		Dimension:  "overall",
		Summary:    "latency up",
		Actor:      "agent",
	})
	if err != nil {
		t.Fatalf("UpsertRegression: %v", err)
	}
	if row.ID == "" || row.CreatedAt == "" || row.UpdatedAt == "" {
		t.Fatalf("expected ids/timestamps: %+v", row)
	}
	if row.Attribution != RegressionAttributionCorrelated || row.Status != RegressionStatusOpen {
		t.Fatalf("defaults: %+v", row)
	}

	got, err := s.GetRegression(row.ID)
	if err != nil || got.SourceID != "outcome-1" || got.Dimension != "overall" {
		t.Fatalf("GetRegression: %+v err=%v", got, err)
	}
	bySrc, err := s.GetRegressionBySource(RegressionSourceEvaluation, "outcome-1", "overall")
	if err != nil || bySrc.ID != row.ID {
		t.Fatalf("GetRegressionBySource: %+v err=%v", bySrc, err)
	}

	open, err := s.HasOpenRegression("task-1")
	if err != nil || !open {
		t.Fatalf("HasOpenRegression: %v %v", open, err)
	}
	n, err := s.CountOpenRegressionsByTaskID("task-1")
	if err != nil || n != 1 {
		t.Fatalf("CountOpenRegressionsByTaskID: %d %v", n, err)
	}
	listed, err := s.ListOpenRegressions("task-1")
	if err != nil || len(listed) != 1 || listed[0].ID != row.ID {
		t.Fatalf("ListOpenRegressions: %+v err=%v", listed, err)
	}

	got.Status = RegressionStatusResolved
	updated, err := s.UpsertRegression(got)
	if err != nil || updated.Status != RegressionStatusResolved {
		t.Fatalf("resolve upsert: %+v err=%v", updated, err)
	}
	open, err = s.HasOpenRegression("task-1")
	if err != nil || open {
		t.Fatalf("resolved must clear open: %v %v", open, err)
	}

	if _, err := s.UpsertRegression(Regression{SourceKind: RegressionSourceEvaluation, SourceID: "x"}); err == nil {
		t.Fatal("expected task_id required")
	}
	if _, err := s.GetRegression(""); err == nil {
		t.Fatal("expected id required")
	}
	if _, err := s.GetRegression("missing"); err == nil {
		t.Fatal("expected missing regression")
	}
	if _, err := s.HasOpenRegression(""); err == nil {
		t.Fatal("expected task_id required")
	}
}

func TestRegressionCreateRequiresCorrelated(t *testing.T) {
	s, _ := openTempStore(t)

	_, err := s.UpsertRegression(Regression{
		TaskID:      "task-1",
		SourceKind:  RegressionSourceEvaluation,
		SourceID:    "outcome-1",
		Dimension:   "overall",
		Attribution: RegressionAttributionCaused,
	})
	if err == nil {
		t.Fatal("create with caused must fail closed")
	}
	_, err = s.UpsertRegression(Regression{
		TaskID:      "task-1",
		SourceKind:  RegressionSourceContradictedEffect,
		SourceID:    "effect-1",
		Dimension:   "latency",
		Attribution: RegressionAttributionHypothesized,
	})
	if err == nil {
		t.Fatal("create with hypothesized must fail closed")
	}

	row, err := s.UpsertRegression(Regression{
		TaskID:     "task-1",
		SourceKind: RegressionSourceEvaluation,
		SourceID:   "outcome-1",
		Dimension:  "overall",
	})
	if err != nil {
		t.Fatal(err)
	}
	row.Attribution = RegressionAttributionHypothesized
	updated, err := s.UpsertRegression(row)
	if err != nil || updated.Attribution != RegressionAttributionHypothesized {
		t.Fatalf("update to hypothesized: %+v err=%v", updated, err)
	}
}

func TestRegressionUnknownAttributionFailClosed(t *testing.T) {
	s, _ := openTempStore(t)

	_, err := s.UpsertRegression(Regression{
		TaskID:      "task-1",
		SourceKind:  RegressionSourceEvaluation,
		SourceID:    "outcome-1",
		Attribution: "inferred",
	})
	if err == nil {
		t.Fatal("unknown attribution must fail closed")
	}

	_, err = s.db.Exec(`
		INSERT INTO regressions(
			id, task_id, source_kind, source_id, dimension, attribution, status,
			summary, actor, source_type, confidence, created_at, updated_at
		) VALUES ('r1','t1','evaluation','o1','overall','guessed','OPEN','','','',0,'now','now')
	`)
	if err == nil {
		t.Fatal("expected CHECK(attribution) to reject unknown value")
	}
}

func TestReflectionStoreRoundtrip(t *testing.T) {
	s, _ := openTempStore(t)

	row, err := s.UpsertReflection(Reflection{
		TaskID:                     "task-1",
		Summary:                    "learned",
		InvalidatedAssumptionsJSON: `["a-1"]`,
		NewDependenciesJSON:        `[{"kind":"path","ref":"internal/foo.go"}]`,
		UsefulTestsJSON:            `["TestFoo"]`,
		BroadenTestsNote:           "add soak",
	})
	if err != nil {
		t.Fatalf("UpsertReflection: %v", err)
	}
	got, err := s.GetReflection(row.ID)
	if err != nil || got.Summary != "learned" || !strings.Contains(got.UsefulTestsJSON, "TestFoo") {
		t.Fatalf("GetReflection: %+v err=%v", got, err)
	}
	listed, err := s.ListReflectionsByTaskID("task-1")
	if err != nil || len(listed) != 1 || listed[0].ID != row.ID {
		t.Fatalf("ListReflectionsByTaskID: %+v err=%v", listed, err)
	}

	if _, err := s.UpsertReflection(Reflection{Summary: "x"}); err == nil {
		t.Fatal("expected task_id required")
	}
	if _, err := s.GetReflection(""); err == nil {
		t.Fatal("expected id required")
	}
	if _, err := s.ListReflectionsByTaskID(""); err == nil {
		t.Fatal("expected task_id required")
	}
}

func TestReflectionSchemaNoBodyColumn(t *testing.T) {
	s, _ := openTempStore(t)

	var name string
	if err := s.db.QueryRow(`SELECT name FROM sqlite_master WHERE name=?`, "reflections").Scan(&name); err != nil {
		t.Fatalf("missing table reflections: %v", err)
	}
	cols, err := tableColumns(s.conn, "reflections")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cols {
		if strings.EqualFold(c, "body") {
			t.Fatalf("reflections must not have body column; got %v", cols)
		}
	}

	_, err = s.UpsertReflection(Reflection{
		TaskID:                     "task-1",
		InvalidatedAssumptionsJSON: `"this is an essay"`,
	})
	if err == nil {
		t.Fatal("string essay in array column must fail closed")
	}
	_, err = s.UpsertReflection(Reflection{
		TaskID:              "task-1",
		NewDependenciesJSON: `{"kind":"path","ref":"x"}`,
	})
	if err == nil {
		t.Fatal("object root in array column must fail closed")
	}

	regCols, err := tableColumns(s.conn, "regressions")
	if err != nil {
		t.Fatal(err)
	}
	if len(regCols) == 0 {
		t.Fatal("regressions has no columns")
	}
	for _, c := range regCols {
		lower := strings.ToLower(c)
		for _, bad := range []string{"blob", "patch", "diff", "content", "log", "cot", "essay"} {
			if strings.Contains(lower, bad) {
				t.Errorf("forbidden column regressions.%s", c)
			}
		}
	}
}
