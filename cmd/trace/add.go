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

func cmdAdd(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace add <kind> --title <title> [flags]\n")
		return exitUsage
	}
	kind := args[0]
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	title := fs.String("title", "", "entity title (required)")
	body := fs.String("body", "", "entity body")
	id := fs.String("id", "", "optional UUID")
	sourceType := fs.String("source-type", "", "provenance source_type")
	confidence := fs.Float64("confidence", 0, "provenance confidence")
	status := fs.String("status", "", "provenance status")
	goalID := fs.String("goal-id", "", "goal id (task only)")
	fromDiscovery := fs.String("from-discovery", "", "promote BLOCKING discovery id into a linked task (task only)")
	severity := fs.String("severity", "", "discovery severity: INFO|PLAN_AFFECTING|BLOCKING (default INFO)")
	if err := fs.Parse(args[1:]); err != nil {
		return exitUsage
	}
	if strings.TrimSpace(*title) == "" && !(kind == "task" && strings.TrimSpace(*fromDiscovery) != "") {
		fmt.Fprintf(os.Stderr, "add: --title is required\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "add: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "add: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	if code := failCLIDenied(svc, "add", "add"); code != exitOK {
		return code
	}

	var outID, outType string
	switch kind {
	case "goal":
		g, err := svc.CreateGoal(ctx, domain.GoalInput{
			ID: *id, Title: *title, Body: *body,
			SourceType: *sourceType, Confidence: *confidence, Status: *status,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "add: %v\n", err)
			return exitFail
		}
		outID, outType = g.ID, domain.EntityGoal
	case "task":
		if strings.TrimSpace(*fromDiscovery) != "" {
			taskID, _, err := svc.PromoteBlockingDiscovery(ctx, *fromDiscovery, *goalID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "add: %v\n", err)
				return exitFail
			}
			outID, outType = taskID, domain.EntityTask
			break
		}
		var gid *string
		if *goalID != "" {
			g := *goalID
			gid = &g
		}
		t, err := svc.CreateTask(ctx, domain.TaskInput{
			ID: *id, GoalID: gid, Title: *title, Body: *body,
			SourceType: *sourceType, Confidence: *confidence, Status: *status,
			// WorkState omitted → PENDING; use transition for DONE
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "add: %v\n", err)
			return exitFail
		}
		outID, outType = t.ID, domain.EntityTask
	case "decision":
		d, err := svc.CreateDecision(ctx, domain.DecisionInput{
			ID: *id, Title: *title, Body: *body,
			SourceType: *sourceType, Confidence: *confidence, Status: *status,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "add: %v\n", err)
			return exitFail
		}
		outID, outType = d.ID, domain.EntityDecision
	case "assumption":
		a, err := svc.CreateAssumption(ctx, domain.AssumptionInput{
			ID: *id, Title: *title, Body: *body,
			SourceType: *sourceType, Confidence: *confidence, Status: *status,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "add: %v\n", err)
			return exitFail
		}
		outID, outType = a.ID, domain.EntityAssumption
	case "discovery":
		d, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{
			ID: *id, Title: *title, Body: *body,
			SourceType: *sourceType, Confidence: *confidence, Status: *status,
			Severity: *severity,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "add: %v\n", err)
			return exitFail
		}
		outID, outType = d.ID, domain.EntityDiscovery
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"ok": true, "type": outType, "id": outID, "severity": d.Severity,
		})
		return exitOK
	case "plan-change":
		p, err := svc.CreatePlanChange(ctx, domain.PlanChangeInput{
			ID: *id, Title: *title, Body: *body,
			SourceType: *sourceType, Confidence: *confidence, Status: *status,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "add: %v\n", err)
			return exitFail
		}
		outID, outType = p.ID, domain.EntityPlanChange
	case "claim":
		c, err := svc.CreateClaim(ctx, domain.ClaimInput{
			ID: *id, Title: *title, Body: *body,
			SourceType: *sourceType, Confidence: *confidence, Status: *status,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "add: %v\n", err)
			return exitFail
		}
		outID, outType = c.ID, domain.EntityClaim
	case "evidence":
		e, err := svc.CreateEvidence(ctx, domain.EvidenceInput{
			ID: *id, Title: *title, Body: *body,
			SourceType: *sourceType, Confidence: *confidence, Status: *status,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "add: %v\n", err)
			return exitFail
		}
		outID, outType = e.ID, domain.EntityEvidence
	default:
		fmt.Fprintf(os.Stderr, "add: unknown kind %q\n", kind)
		return exitUsage
	}

	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(map[string]any{"ok": true, "type": outType, "id": outID})
	return exitOK
}
