package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

func cmdLink(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace link <rel> --from <id> --to <id>\n")
		return exitUsage
	}
	rel := args[0]
	fs := flag.NewFlagSet("link", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	from := fs.String("from", "", "from entity UUID")
	to := fs.String("to", "", "to entity UUID")
	sourceType := fs.String("source-type", "", "optional source_type")
	if err := fs.Parse(args[1:]); err != nil {
		return exitUsage
	}
	if *from == "" || *to == "" {
		fmt.Fprintf(os.Stderr, "link: --from and --to are required\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "link: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "link: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	if code := failCLIDenied(svc, "link", "link"); code != exitOK {
		return code
	}
	meta := domain.LinkMeta{SourceType: *sourceType}

	switch rel {
	case "goal-task":
		err = svc.LinkGoalTask(ctx, *from, *to, meta)
	case "decision-task":
		err = svc.LinkDecisionTask(ctx, *from, *to, meta)
	case "discovery-plan-change":
		err = svc.LinkDiscoveryPlanChange(ctx, *from, *to, meta)
	case "discovery-mentions-task":
		err = svc.LinkDiscoveryMentionsTask(ctx, *from, *to, meta)
	case "claim-evidence":
		err = svc.LinkClaimEvidence(ctx, *from, *to, meta)
	default:
		fmt.Fprintf(os.Stderr, "link: unknown rel %q\n", rel)
		return exitUsage
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "link: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "rel": rel, "from": *from, "to": *to,
	})
	return exitOK
}
