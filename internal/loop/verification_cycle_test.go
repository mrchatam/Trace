package loop_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
)

func TestVerificationCycleBlocksSkipInStatus(t *testing.T) {
	st, psvc, dsvc := openLoopTestStore(t)
	ctx := context.Background()
	goalID, taskID, _ := seedGoalTaskPlan(t, psvc, dsvc)
	markPlanCritiqued(t, st, taskID, goalID)

	if _, err := dsvc.CreateChange(ctx, domain.ChangeInput{
		TaskID:    taskID,
		GitCommit: "abc1234",
		Paths:     []domain.ChangePathInput{{Path: "main.go"}},
	}); err != nil {
		t.Fatalf("CreateChange: %v", err)
	}

	status, err := loop.Status(ctx, st, psvc, loop.ApplySeed{TaskID: taskID, GoalID: goalID})
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if status.VerificationCycle == nil {
		t.Fatal("verification_cycle block missing")
	}
	if !status.VerificationCycle.TestPending {
		t.Fatalf("test_pending want true: %#v", status.VerificationCycle)
	}
	if status.VerificationCycle.IncompleteReason == "" {
		t.Fatalf("incomplete_reason required when test pending: %#v", status.VerificationCycle)
	}
	if !strings.Contains(status.VerificationCycle.IncompleteReason, "test_pending") {
		t.Fatalf("incomplete_reason=%q must include test_pending", status.VerificationCycle.IncompleteReason)
	}

	raw, err := json.Marshal(status)
	if err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	vc, ok := decoded["verification_cycle"].(map[string]any)
	if !ok {
		t.Fatalf("JSON verification_cycle missing: %#v", decoded)
	}
	if vc["incomplete_reason"] == nil || !strings.Contains(vc["incomplete_reason"].(string), "test_pending") {
		t.Fatalf("JSON incomplete_reason=%v", vc["incomplete_reason"])
	}
}
