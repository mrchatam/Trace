package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// cmdReview is a thin adapter: create|set|get|show|list|residual.
// Business logic lives in internal/domain only (G19).
func cmdReview(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace review create|set|get|show|list|residual …\n")
		return exitUsage
	}
	sub := args[0]
	switch sub {
	case "create":
		return cmdReviewCreate(root, args[1:])
	case "set":
		return cmdReviewSet(root, args[1:])
	case "get", "show":
		return cmdReviewGet(root, args[1:])
	case "list":
		return cmdReviewList(root, args[1:])
	case "residual":
		return cmdReviewResidual(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "review: unknown subcommand %q (want create|set|get|show|list|residual)\n", sub)
		return exitUsage
	}
}

// reviewGetRow is DF-45 get/show snake_case shape (includes body).
type reviewGetRow struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Result string `json:"result"`
	Status string `json:"status"`
}

// reviewListRow is DF-45 list snake_case shape (body omitted).
type reviewListRow struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Result string `json:"result"`
	Status string `json:"status"`
}

func cmdReviewGet(root string, args []string) int {
	fs := flag.NewFlagSet("review get", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	id := fs.String("id", "", "review UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *id == "" {
		fmt.Fprintf(os.Stderr, "usage: trace review get|show --id <id>\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "review", "review"); code != exitOK {
		return code
	}
	r, err := svc.GetReview(context.Background(), *id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(reviewGetRow{
		ID: r.ID, Title: r.Title, Body: r.Body, Result: r.Result, Status: r.Status,
	})
	return exitOK
}

func cmdReviewList(root string, args []string) int {
	fs := flag.NewFlagSet("review list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "optional task UUID filter (review_judges_task)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "review", "review"); code != exitOK {
		return code
	}
	ctx := context.Background()

	var list []store.Review
	if *taskID != "" {
		list, err = svc.ListReviewsByTaskID(ctx, *taskID)
	} else {
		list, err = svc.ListReviews(ctx)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	out := make([]reviewListRow, 0, len(list))
	for _, r := range list {
		out = append(out, reviewListRow{
			ID: r.ID, Title: r.Title, Result: r.Result, Status: r.Status,
		})
	}
	_ = json.NewEncoder(os.Stdout).Encode(out)
	return exitOK
}

func cmdReviewCreate(root string, args []string) int {
	fs := flag.NewFlagSet("review create", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	title := fs.String("title", "", "review title (required)")
	body := fs.String("body", "", "review body")
	id := fs.String("id", "", "optional UUID")
	taskID := fs.String("task", "", "optional task UUID to link (review_judges_task)")
	scopeID := fs.String("scope", "", "optional plan_scope UUID to link (review_judges_scope)")
	sourceType := fs.String("source-type", "", "provenance source_type")
	confidence := fs.Float64("confidence", 0, "provenance confidence")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if strings.TrimSpace(*title) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace review create --title <title> [--task <id>] [--scope <plan_scope_id>] [--id <uuid>]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "review", "review"); code != exitOK {
		return code
	}
	ctx := context.Background()

	r, err := svc.CreateReview(ctx, domain.ReviewInput{
		ID: *id, Title: *title, Body: *body,
		SourceType: *sourceType, Confidence: *confidence,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	if *taskID != "" {
		if err := svc.LinkReviewTask(ctx, r.ID, *taskID, domain.LinkMeta{}); err != nil {
			fmt.Fprintf(os.Stderr, "review: link task: %v\n", err)
			return exitFail
		}
	}
	if *scopeID != "" {
		if err := svc.LinkReviewScope(ctx, r.ID, *scopeID, domain.LinkMeta{}); err != nil {
			fmt.Fprintf(os.Stderr, "review: link scope: %v\n", err)
			return exitFail
		}
	}
	out := map[string]any{"ok": true, "id": r.ID, "type": domain.EntityReview, "result": r.Result}
	if *taskID != "" {
		out["task"] = *taskID
		out["rel"] = domain.RelReviewJudgesTask
	}
	if *scopeID != "" {
		out["scope"] = *scopeID
		out["scope_rel"] = domain.RelReviewJudgesScope
	}
	_ = json.NewEncoder(os.Stdout).Encode(out)
	return exitOK
}

func cmdReviewSet(root string, args []string) int {
	fs := flag.NewFlagSet("review set", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	id := fs.String("id", "", "review UUID")
	result := fs.String("result", "", "PASS|FAIL|UNCERTAIN")
	actor := fs.String("actor", "", "actor (required non-empty)")
	reason := fs.String("reason", "", "reason (required non-empty)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *id == "" || *result == "" || strings.TrimSpace(*reason) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace review set --id <id> --result PASS|FAIL|UNCERTAIN --reason <text> [--actor <a>]\n")
		return exitUsage
	}
	if *actor == "" {
		*actor = "cli"
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "review", "review"); code != exitOK {
		return code
	}
	err = svc.SetReviewResult(context.Background(), *id, strings.ToUpper(*result), domain.ReviewResultOptions{
		Actor: *actor, Reason: *reason,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "id": *id, "result": strings.ToUpper(*result),
	})
	return exitOK
}

func cmdReviewResidual(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace review residual add|list …\n")
		return exitUsage
	}
	switch args[0] {
	case "add":
		return cmdReviewResidualAdd(root, args[1:])
	case "list":
		return cmdReviewResidualList(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "review residual: unknown subcommand %q (want add|list)\n", args[0])
		return exitUsage
	}
}

func cmdReviewResidualAdd(root string, args []string) int {
	fs := flag.NewFlagSet("review residual add", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	reviewID := fs.String("review", "", "review UUID")
	code := fs.String("code", "", "residual code (required)")
	body := fs.String("body", "", "residual body")
	severity := fs.String("severity", "", "INFO|WARN|BLOCKING (default INFO)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *reviewID == "" || strings.TrimSpace(*code) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace review residual add --review <id> --code <CODE> [--body <text>] [--severity INFO|WARN|BLOCKING]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "review", "review"); code != exitOK {
		return code
	}
	res, err := svc.AddResidual(context.Background(), *reviewID, domain.ResidualInput{
		Code: *code, Body: *body, Severity: *severity,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "id": res.ID, "review": res.ReviewID, "code": res.Code,
		"severity": res.Severity, "status": res.Status,
	})
	return exitOK
}

func cmdReviewResidualList(root string, args []string) int {
	fs := flag.NewFlagSet("review residual list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	reviewID := fs.String("review", "", "review UUID")
	scopeID := fs.String("scope", "", "plan_scope UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if (*reviewID == "") == (*scopeID == "") {
		fmt.Fprintf(os.Stderr, "usage: trace review residual list --review <id> | --scope <plan_scope_id>\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "review", "review"); code != exitOK {
		return code
	}
	ctx := context.Background()

	var list []store.ReviewResidual
	if *reviewID != "" {
		list, err = svc.ListResidualsByReview(ctx, *reviewID)
	} else {
		list, err = svc.ListResidualsByScope(ctx, *scopeID)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "review: %v\n", err)
		return exitFail
	}
	if list == nil {
		list = []store.ReviewResidual{}
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "residuals": list, "count": len(list),
	})
	return exitOK
}
