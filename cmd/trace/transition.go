package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

// stringList is a repeatable flag.Value for --evidence.
type stringList []string

func (s *stringList) String() string { return strings.Join(*s, ",") }
func (s *stringList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func cmdTransition(root string, args []string) int {
	fs := flag.NewFlagSet("transition", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "task UUID")
	to := fs.String("to", "", "target work_state")
	actor := fs.String("actor", "", "actor")
	reason := fs.String("reason", "", "reason (required non-empty)")
	allowDone := fs.Bool("allow-done", false, "AllowDoneWithoutReview escape hatch (loud WARNING on success)")
	asOperator := fs.Bool("as-operator", false, "AllowOperatorDone (conscious claim; flag≠identity / not verified operator identity; Actor string ≠ auth)")
	allowMissingCaps := fs.Bool("allow-missing-caps", false, "AllowMissingCapabilities override")
	enforce := fs.Bool("enforce", false, "Run deliberation gate before DONE (--to DONE only)")
	var evidence stringList
	fs.Var(&evidence, "evidence", "evidence id (repeatable)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *taskID == "" || *to == "" || strings.TrimSpace(*reason) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace transition --task <id> --to <state> --reason <text> [--actor <a>] [--as-operator] [--allow-done] [--allow-missing-caps] [--enforce] [--evidence id]\n")
		fmt.Fprintf(os.Stderr, "note: DONE requires linked Review PASS + --as-operator (no linked FAIL), or --allow-done hatch; --as-operator is a conscious claim (flag≠identity / not verified identity); --evidence alone does not authorize DONE; Actor is not authorization; --allow-done does not bypass missing capabilities (use --allow-missing-caps); optional --enforce runs deliberation gate (GateForDone) before DONE\n")
		return exitUsage
	}
	if *actor == "" {
		*actor = "cli"
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "transition: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "transition: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "transition", "transition"); code != exitOK {
		return code
	}
	ctx := context.Background()
	if *enforce && strings.EqualFold(strings.TrimSpace(*to), store.WorkStateDone) {
		plan := planner.New(st)
		allowed, violations, gateErr := loop.EvaluateGate(ctx, svc, plan, st, *taskID, loop.GateForDone)
		if gateErr != nil {
			fmt.Fprintf(os.Stderr, "transition: gate: %v\n", gateErr)
			return exitFail
		}
		if !allowed {
			if len(violations) > 0 {
				fmt.Fprintf(os.Stderr, "transition: %s\n", violations[0].Message)
			} else {
				fmt.Fprintln(os.Stderr, "transition: gate blocked")
			}
			return exitGateBlocked
		}
	}
	err = svc.TransitionTask(ctx, *taskID, *to, domain.TransitionOptions{
		Actor:                    *actor,
		Reason:                   *reason,
		EvidenceIDs:              []string(evidence),
		AllowDoneWithoutReview:   *allowDone,
		AllowOperatorDone:        *asOperator,
		AllowMissingCapabilities: *allowMissingCaps,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "transition: %v\n", err)
		return exitFail
	}
	if *allowDone {
		fmt.Fprintf(os.Stderr, "WARNING: --allow-done escape hatch used; Review PASS and --as-operator were bypassed; missing capabilities still need --allow-missing-caps\n")
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "task": *taskID, "to": *to,
	})
	return exitOK
}
