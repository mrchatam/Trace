package domain_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func TestOperatorDoneRequiresFlag(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "DF-17"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "agent", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}
	rev, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "ok"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, rev.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, rev.ID, store.ReviewResultPass, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "pass",
	}); err != nil {
		t.Fatal(err)
	}

	// PASS without operator flag → reject (Actor spoof insufficient)
	err = svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "operator", Reason: "spoof",
	})
	if err == nil {
		t.Fatal("DONE after PASS without AllowOperatorDone must fail")
	}
	var inv *domain.ErrInvalidTransition
	if !errors.As(err, &inv) {
		t.Fatalf("want ErrInvalidTransition, got %T: %v", err, err)
	}
	if !strings.Contains(inv.Reason, "AllowOperatorDone") && !strings.Contains(inv.Reason, "as-operator") {
		t.Fatalf("reason should mention AllowOperatorDone/--as-operator: %q", inv.Reason)
	}

	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "agent", Reason: "promote", AllowOperatorDone: true,
	}); err != nil {
		t.Fatalf("PASS + AllowOperatorDone: %v", err)
	}
	got, _ := svc.GetTask(ctx, task.ID)
	if got.WorkState != store.WorkStateDone {
		t.Fatalf("want DONE, got %q", got.WorkState)
	}
}

func TestOperatorDoneHatchBypassesOperator(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "hatch"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "go",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "a", Reason: "escape", AllowDoneWithoutReview: true,
	}); err != nil {
		t.Fatalf("hatch without AllowOperatorDone: %v", err)
	}
}

func TestReopenInvalidatesPassReviews(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "DF-18"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}
	rev, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "first pass"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, rev.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, rev.ID, store.ReviewResultPass, domain.ReviewResultOptions{
		Actor: "r", Reason: "ok",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "op", Reason: "done", AllowOperatorDone: true,
	}); err != nil {
		t.Fatal(err)
	}

	// Leave DONE → PENDING: PASS becomes UNCERTAIN
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStatePending, domain.TransitionOptions{
		Actor: "op", Reason: "reopen",
	}); err != nil {
		t.Fatalf("reopen: %v", err)
	}
	gotRev, err := svc.GetReview(ctx, rev.ID)
	if err != nil {
		t.Fatal(err)
	}
	if gotRev.Result != store.ReviewResultUncertain {
		t.Fatalf("want UNCERTAIN after reopen, got %q", gotRev.Result)
	}

	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "restart",
	}); err != nil {
		t.Fatal(err)
	}
	// Sticky PASS gone: DONE rejected even with operator flag
	err = svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "op", Reason: "retry sticky", AllowOperatorDone: true,
	})
	if err == nil {
		t.Fatal("DONE with invalidated PASS must fail")
	}

	// New PASS + operator → OK
	rev2, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "second pass"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, rev2.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, rev2.ID, store.ReviewResultPass, domain.ReviewResultOptions{
		Actor: "r", Reason: "ok again",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "op", Reason: "done again", AllowOperatorDone: true,
	}); err != nil {
		t.Fatalf("new PASS + operator: %v", err)
	}
}

func TestMissingCapabilitiesBlockTransition(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "DF-24"})
	if err != nil {
		t.Fatal(err)
	}
	cap, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: "MCP", Slug: "mcp:needed", Title: "Needed", Status: store.CapabilityStatusUnavailable,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RequireCapability(ctx, task.ID, cap.ID); err != nil {
		t.Fatal(err)
	}

	err = svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "start",
	})
	if err == nil {
		t.Fatal("transition with missing caps must fail")
	}
	var inv *domain.ErrInvalidTransition
	if !errors.As(err, &inv) {
		t.Fatalf("want ErrInvalidTransition, got %T: %v", err, err)
	}
	if !strings.Contains(inv.Reason, "missing") {
		t.Fatalf("reason should mention missing caps: %q", inv.Reason)
	}

	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "override", AllowMissingCapabilities: true,
	}); err != nil {
		t.Fatalf("AllowMissingCapabilities: %v", err)
	}
}

// TestAllowDoneDoesNotBypassMissingCaps — DF-51: hatch ≠ AllowMissingCapabilities.
func TestAllowDoneDoesNotBypassMissingCaps(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "DF-51"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}
	cap, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind: "MCP", Slug: "mcp:hatch-caps", Title: "Needed", Status: store.CapabilityStatusUnavailable,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RequireCapability(ctx, task.ID, cap.ID); err != nil {
		t.Fatal(err)
	}

	err = svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "a", Reason: "hatch alone", AllowDoneWithoutReview: true,
	})
	if err == nil {
		t.Fatal("AllowDoneWithoutReview alone must not bypass missing caps")
	}
	var inv *domain.ErrInvalidTransition
	if !errors.As(err, &inv) {
		t.Fatalf("want ErrInvalidTransition, got %T: %v", err, err)
	}
	if !strings.Contains(inv.Reason, "missing") {
		t.Fatalf("reason should mention missing caps: %q", inv.Reason)
	}

	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "a", Reason: "hatch+override", AllowDoneWithoutReview: true, AllowMissingCapabilities: true,
	}); err != nil {
		t.Fatalf("hatch + AllowMissingCapabilities: %v", err)
	}
}
