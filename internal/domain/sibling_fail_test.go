package domain_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// TestSiblingFailBlocksDone locks DF-43: linked FAIL blocks →DONE even when a
// sibling PASS exists and AllowOperatorDone is set. UNCERTAIN/empty do not block;
// hatch bypasses FAIL; PASS alone still authorizes with the operator flag.
func TestSiblingFailBlocksDone(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "DF-43 sibling FAIL+PASS"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}

	revFail, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "fail review"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, revFail.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, revFail.ID, store.ReviewResultFail, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "not ready",
	}); err != nil {
		t.Fatal(err)
	}

	revPass, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "pass review"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, revPass.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, revPass.ID, store.ReviewResultPass, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "looks good",
	}); err != nil {
		t.Fatal(err)
	}

	err = svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "op", Reason: "promote despite FAIL", AllowOperatorDone: true,
	})
	if err == nil {
		t.Fatal("DONE with linked FAIL+PASS + AllowOperatorDone must reject")
	}
	var inv *domain.ErrInvalidTransition
	if !errors.As(err, &inv) {
		t.Fatalf("want ErrInvalidTransition, got %T: %v", err, err)
	}
	if !strings.Contains(inv.Reason, "FAIL") {
		t.Fatalf("reason must mention FAIL: %q", inv.Reason)
	}
	got, _ := svc.GetTask(ctx, task.ID)
	if got.WorkState != store.WorkStateInProgress {
		t.Fatalf("want IN_PROGRESS after reject, got %q", got.WorkState)
	}

	// Explicit supersession: clear FAIL → UNCERTAIN, then PASS+operator DONE OK.
	if err := svc.SetReviewResult(ctx, revFail.ID, store.ReviewResultUncertain, domain.ReviewResultOptions{
		Actor: "reviewer", Reason: "superseded by later review",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "op", Reason: "promote after supersede", AllowOperatorDone: true,
	}); err != nil {
		t.Fatalf("PASS+UNCERTAIN + AllowOperatorDone after supersede: %v", err)
	}
	got, _ = svc.GetTask(ctx, task.ID)
	if got.WorkState != store.WorkStateDone {
		t.Fatalf("want DONE after supersede, got %q", got.WorkState)
	}
}

func TestSiblingPassAloneAllowsDone(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "DF-43 PASS alone"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}
	rev, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "only pass"})
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
		t.Fatalf("PASS alone + flag: %v", err)
	}
}

func TestSiblingPassPlusUncertainAllowsDone(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "DF-43 PASS+UNCERTAIN"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}
	revU, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "uncertain"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, revU.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, revU.ID, store.ReviewResultUncertain, domain.ReviewResultOptions{
		Actor: "r", Reason: "maybe",
	}); err != nil {
		t.Fatal(err)
	}
	revP, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "pass"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, revP.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, revP.ID, store.ReviewResultPass, domain.ReviewResultOptions{
		Actor: "r", Reason: "ok",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "op", Reason: "done", AllowOperatorDone: true,
	}); err != nil {
		t.Fatalf("PASS+UNCERTAIN + flag: %v", err)
	}
}

func TestHatchBypassesSiblingFail(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	task, err := svc.CreateTask(ctx, domain.TaskInput{Title: "DF-43 hatch+FAIL"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateInProgress, domain.TransitionOptions{
		Actor: "a", Reason: "start",
	}); err != nil {
		t.Fatal(err)
	}
	revFail, err := svc.CreateReview(ctx, domain.ReviewInput{Title: "fail"})
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.LinkReviewTask(ctx, revFail.ID, task.ID, domain.LinkMeta{}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SetReviewResult(ctx, revFail.ID, store.ReviewResultFail, domain.ReviewResultOptions{
		Actor: "r", Reason: "no",
	}); err != nil {
		t.Fatal(err)
	}
	if err := svc.TransitionTask(ctx, task.ID, store.WorkStateDone, domain.TransitionOptions{
		Actor: "a", Reason: "escape", AllowDoneWithoutReview: true,
	}); err != nil {
		t.Fatalf("hatch with linked FAIL must allow DONE: %v", err)
	}
}
