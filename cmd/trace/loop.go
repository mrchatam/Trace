package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mrchatam/Trace/internal/compiler"
	"github.com/mrchatam/Trace/internal/config"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

const (
	gateSchemaVersion = "trace.loop.gate.v1"
	exitGateBlocked   = 1
)

type gateEnvelope struct {
	SchemaVersion    string           `json:"schema_version"`
	TaskID           string           `json:"task_id"`
	For              string           `json:"for"`
	Allowed          bool             `json:"allowed"`
	RecommendedPhase string           `json:"recommended_phase,omitempty"`
	ReasonCode       string           `json:"reason_code,omitempty"`
	Violations       []loop.Violation `json:"violations"`
}

func cmdLoop(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace loop next --task <id>\n")
		return exitUsage
	}
	switch args[0] {
	case "next":
		return cmdLoopNext(root, args[1:])
	case "apply":
		return cmdLoopApply(root, args[1:])
	case "status":
		return cmdLoopStatus(root, args[1:])
	case "gate":
		return cmdLoopGate(root, args[1:])
	case "reset":
		return cmdLoopReset(root, args[1:])
	case "help", "-h", "--help":
		printLoopHelp(os.Stdout)
		return exitOK
	default:
		fmt.Fprintf(os.Stderr, "loop: unknown subcommand %q\n", args[0])
		printLoopHelp(os.Stderr)
		return exitUsage
	}
}

func printLoopHelp(w *os.File) {
	fmt.Fprint(w, `trace loop — stdout-first loop helpers

Subcommands:
  next --task <id>
        Emit one bounded JSON packet for a seed task. Derives goal context from
        the task's goal_id; no stdin or interactive fallback.
  apply [--in <path>]
        Apply a trace.loop.apply.v1 write envelope from --in JSON file or stdin.
        For BLOCKING-discovery promotion, set writes.spawned_tasks[].discovery_id.
  status --task <id> [--goal <id>]
        Report trace.loop.status.v1 from persisted loop-step evidence.
        Includes violations[] (edit gate parity). Optional .trace/config.json
        enforce mode (off|warn|strict, default off): warn/strict print stderr
        hints when violations present; exit stays 0.
  gate --task <id> [--for orient|edit|execute|done|export]
        Check deliberation gate for a task. Emits trace.loop.gate.v1 JSON on stdout.
        Exit 0 when allowed, 1 when blocked, 2 on usage or internal error.
        Default --for is edit (pre-edit harness choke point).
  reset --task <id>
        Clear sticky STOP, hop_count, and consecutive empty-apply counter.
        Sets phase to EXECUTE; preserves plan_critiqued.
`)
}

func parseLoopGateFor(s string) (loop.GateFor, error) {
	switch s {
	case string(loop.GateForOrient):
		return loop.GateForOrient, nil
	case string(loop.GateForEdit):
		return loop.GateForEdit, nil
	case string(loop.GateForExecute):
		return loop.GateForExecute, nil
	case string(loop.GateForDone):
		return loop.GateForDone, nil
	case string(loop.GateForExport):
		return loop.GateForExport, nil
	default:
		return "", fmt.Errorf("loop gate: invalid --for %q (want orient|edit|execute|done|export)", s)
	}
}

func buildGateEnvelope(taskID, gateFor string, allowed bool, violations []loop.Violation) gateEnvelope {
	env := gateEnvelope{
		SchemaVersion: gateSchemaVersion,
		TaskID:        taskID,
		For:           gateFor,
		Allowed:       allowed,
		Violations:    []loop.Violation{},
	}
	if len(violations) > 0 {
		env.Violations = violations
	}
	if !allowed && len(violations) == 1 {
		env.RecommendedPhase = violations[0].RecommendedPhase
		env.ReasonCode = violations[0].ReasonCode
	}
	return env
}

func cmdLoopGate(root string, args []string) int {
	fs := flag.NewFlagSet("loop gate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "seed task UUID")
	gateFor := fs.String("for", "edit", "gate context: orient|edit|execute|done|export")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true, "for": true})); err != nil {
		return exitFail
	}
	if *taskID == "" || fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace loop gate --task <id> [--for orient|edit|execute|done|export]\n")
		return exitFail
	}
	gf, err := parseLoopGateFor(*gateFor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return exitFail
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop gate: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop gate: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "loop", "loop"); code != exitOK {
		return code
	}

	allowed, violations, err := loop.EvaluateGate(
		context.Background(),
		domain.New(st),
		planner.New(st),
		st,
		*taskID,
		gf,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return exitFail
	}

	env := buildGateEnvelope(*taskID, *gateFor, allowed, violations)
	if err := json.NewEncoder(os.Stdout).Encode(env); err != nil {
		fmt.Fprintf(os.Stderr, "loop gate: %v\n", err)
		return exitFail
	}
	if allowed {
		return exitOK
	}
	if len(violations) > 0 {
		fmt.Fprintln(os.Stderr, violations[0].Message)
	}
	return exitGateBlocked
}

func cmdLoopNext(root string, args []string) int {
	fs := flag.NewFlagSet("loop next", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "seed task UUID")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true})); err != nil {
		return exitUsage
	}
	if *taskID == "" || fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace loop next --task <id>\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop next: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop next: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "loop", "loop"); code != exitOK {
		return code
	}

	eng := retrieval.New(st)
	if repo, rerr := tryOpenGit(abs, st); rerr == nil {
		defer repo.Close()
		eng = eng.WithVCS(repo)
	}
	comp := compiler.New(st).WithRetrieval(eng)

	packet, err := loop.BuildNextPacket(context.Background(), loop.BuildNextInput{
		TaskID:    *taskID,
		Store:     st,
		Planner:   planner.New(st),
		Retrieval: eng,
		Compiler:  comp,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(packet); err != nil {
		fmt.Fprintf(os.Stderr, "loop next: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdLoopApply(root string, args []string) int {
	fs := flag.NewFlagSet("loop apply", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	inPath := fs.String("in", "", "path to trace.loop.apply.v1 JSON (reads stdin when omitted)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"in": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace loop apply [--in <path>]\n")
		return exitUsage
	}

	raw, err := readLoopApplyInput(*inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop apply: %v\n", err)
		return exitFail
	}
	env, err := loop.ParseApplyEnvelope(raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return exitFail
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop apply: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop apply: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "loop", "loop"); code != exitOK {
		return code
	}

	res, err := loop.Apply(context.Background(), st, planner.New(st), env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop apply: %v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(res); err != nil {
		fmt.Fprintf(os.Stderr, "loop apply: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdLoopStatus(root string, args []string) int {
	fs := flag.NewFlagSet("loop status", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "seed task UUID")
	goalID := fs.String("goal", "", "seed goal UUID (optional; derived from task)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true, "goal": true})); err != nil {
		return exitUsage
	}
	if *taskID == "" || fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace loop status --task <id> [--goal <id>]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop status: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop status: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "loop", "loop"); code != exitOK {
		return code
	}

	res, err := loop.Status(context.Background(), st, planner.New(st), loop.ApplySeed{TaskID: *taskID, GoalID: *goalID})
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop status: %v\n", err)
		return exitFail
	}
	mode := config.LoadEnforceMode(abs)
	if (mode == config.EnforceWarn || mode == config.EnforceStrict) && len(res.Violations) > 0 {
		for _, v := range res.Violations {
			fmt.Fprintf(os.Stderr, "loop status: %s\n", v.Message)
		}
	}
	if err := json.NewEncoder(os.Stdout).Encode(res); err != nil {
		fmt.Fprintf(os.Stderr, "loop status: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdLoopReset(root string, args []string) int {
	fs := flag.NewFlagSet("loop reset", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "seed task UUID")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true})); err != nil {
		return exitUsage
	}
	if *taskID == "" || fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace loop reset --task <id>\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop reset: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop reset: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "loop", "loop"); code != exitOK {
		return code
	}

	next, err := domain.New(st).ResetDeliberationState(context.Background(), *taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "loop reset: %v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(next); err != nil {
		fmt.Fprintf(os.Stderr, "loop reset: %v\n", err)
		return exitFail
	}
	return exitOK
}

func readLoopApplyInput(path string) ([]byte, error) {
	if path != "" {
		return os.ReadFile(path)
	}
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		return nil, fmt.Errorf("stdin is empty")
	}
	return raw, nil
}
