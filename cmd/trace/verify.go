package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/testrun"
)

func cmdVerify(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace verify run|invariants\n")
		return exitUsage
	}
	switch args[0] {
	case "run":
		return cmdVerifyRun(root, args[1:])
	case "invariants":
		return cmdVerifyInvariants(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown verify subcommand: %s\n", args[0])
		return exitUsage
	}
}

func cmdVerifyInvariants(root string, args []string) int {
	fs := flag.NewFlagSet("verify invariants", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "task UUID")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 || strings.TrimSpace(*taskID) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace verify invariants --task <id>\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	if code := failCLIDenied(svc, "verify", "verify invariants"); code != exitOK {
		return code
	}

	result, err := svc.CheckArchitecturalInvariants(context.Background(), *taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify invariants: %v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(result); err != nil {
		fmt.Fprintf(os.Stderr, "verify invariants: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdVerifyRun(root string, args []string) int {
	fs := flag.NewFlagSet("verify run", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "task UUID")
	forceEval := fs.Bool("force-eval", false, "run evaluation even when latest test failed")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true, "force-eval": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 || strings.TrimSpace(*taskID) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace verify run --task <id> [--force-eval]\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	if code := failCLIDenied(svc, "verify", "verify run"); code != exitOK {
		return code
	}

	result, err := svc.CoordinateVerification(context.Background(), *taskID, domain.CoordinateOptions{
		ForceEval:  *forceEval,
		RunTests:   testrun.CoordinateTestRun(st, svc, testrun.Options{Actor: "cli", SourceType: "CLI"}),
		Actor:      "cli",
		SourceType: "CLI",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "verify run: %v\n", err)
		return exitFail
	}

	payload := map[string]any{
		"ok":                  true,
		"task":                result.TaskID,
		"verification_found":  result.VerificationFound,
		"evaluation_recorded": result.EvaluationRecorded,
		"stopped_early":       result.StoppedEarly,
		"stop_reason":         result.StopReason,
		"regression_detected": result.Regression.Detected,
		"test_name":           result.Regression.TestName,
		"steps":               result.Steps,
	}
	if len(result.TestOutcomes) > 0 {
		rows := make([]map[string]any, 0, len(result.TestOutcomes))
		for _, o := range result.TestOutcomes {
			rows = append(rows, map[string]any{
				"id": o.ID, "test_name": o.TestName, "test_status": o.TestStatus,
			})
		}
		payload["test_outcomes"] = rows
	}
	if err := json.NewEncoder(os.Stdout).Encode(payload); err != nil {
		fmt.Fprintf(os.Stderr, "verify run: %v\n", err)
		return exitFail
	}
	return exitOK
}
