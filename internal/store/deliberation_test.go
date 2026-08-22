package store

import (
	"database/sql"
	"errors"
	"testing"
)

func TestDeliberationStateRoundtrip(t *testing.T) {
	s, _ := openTempStore(t)

	got, err := s.GetDeliberationState("missing")
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("missing: got %+v err=%v", got, err)
	}

	in := DeliberationState{
		TaskID:                  "task-seed",
		GoalID:                  "goal-1",
		CurrentPhase:            "INVESTIGATE",
		HopCount:                3,
		LastPhase:               "ORIENT",
		PlanCritiqued:           true,
		Stopped:                 false,
		StopReason:              "",
		ConsecutiveEmptyApplies: 1,
	}
	out, err := s.UpsertDeliberationState(in)
	if err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	if out.UpdatedAt == "" {
		t.Fatal("expected updated_at")
	}
	if out.TaskID != in.TaskID || out.GoalID != in.GoalID || out.CurrentPhase != in.CurrentPhase {
		t.Fatalf("out identity: %+v", out)
	}
	if out.HopCount != 3 || !out.PlanCritiqued || out.Stopped || out.LastPhase != "ORIENT" {
		t.Fatalf("out fields: %+v", out)
	}
	if out.ConsecutiveEmptyApplies != 1 {
		t.Fatalf("consecutive_empty_applies: %d want 1", out.ConsecutiveEmptyApplies)
	}

	loaded, err := s.GetDeliberationState("task-seed")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if loaded.HopCount != 3 || !loaded.PlanCritiqued {
		t.Fatalf("loaded: %+v", loaded)
	}

	loaded.Stopped = true
	loaded.StopReason = "hop_budget_exceeded"
	loaded.HopCount = 12
	loaded.CurrentPhase = "STOP"
	updated, err := s.UpsertDeliberationState(loaded)
	if err != nil {
		t.Fatalf("Upsert stop: %v", err)
	}
	if !updated.Stopped || updated.StopReason != "hop_budget_exceeded" || updated.HopCount != 12 {
		t.Fatalf("stopped row: %+v", updated)
	}

	if _, err := s.UpsertDeliberationState(DeliberationState{TaskID: "x"}); err == nil {
		t.Fatal("expected goal_id required")
	}
	if _, err := s.UpsertDeliberationState(DeliberationState{GoalID: "g"}); err == nil {
		t.Fatal("expected task_id required")
	}
}

func TestDeliberationStateTableMigrated(t *testing.T) {
	s, _ := openTempStore(t)
	var name string
	err := s.db.QueryRow(`SELECT name FROM sqlite_master WHERE name=?`, "deliberation_state").Scan(&name)
	if err != nil {
		t.Fatalf("deliberation_state missing: %v", err)
	}
	st, err := s.MigrationStatus()
	if err != nil {
		t.Fatalf("MigrationStatus: %v", err)
	}
	if st.EmbedExpected != 28 || st.MaxApplied != 28 {
		t.Fatalf("embed/applied: %+v want 28", st)
	}
}
