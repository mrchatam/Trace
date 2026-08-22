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

type testRunRow struct {
	ID         string `json:"id"`
	TestName   string `json:"test_name"`
	TestStatus string `json:"test_status"`
	Summary    string `json:"summary,omitempty"`
}

func cmdTest(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace test run\n")
		return exitUsage
	}
	switch args[0] {
	case "run":
		return cmdTestRun(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown test subcommand: %s\n", args[0])
		return exitUsage
	}
}

func cmdTestRun(root string, args []string) int {
	fs := flag.NewFlagSet("test run", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	taskID := fs.String("task", "", "task UUID")
	paths := fs.String("paths", "", "comma-separated seed paths (optional)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"task": true, "paths": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 || strings.TrimSpace(*taskID) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace test run --task <id> [--paths path,...]\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	if code := failCLIDenied(svc, "test", "test run"); code != exitOK {
		return code
	}

	var seedPaths []string
	if strings.TrimSpace(*paths) != "" {
		for _, p := range strings.Split(*paths, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				seedPaths = append(seedPaths, p)
			}
		}
	}

	outcomes, err := testrun.RunRelevantTests(context.Background(), st, svc, *taskID, testrun.Options{
		Paths: seedPaths,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "test run: %v\n", err)
		return exitFail
	}

	rows := make([]testRunRow, 0, len(outcomes))
	for _, o := range outcomes {
		rows = append(rows, testRunRow{
			ID: o.ID, TestName: o.TestName, TestStatus: o.TestStatus, Summary: o.Summary,
		})
	}
	if err := json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "task": *taskID, "outcomes": rows, "count": len(rows),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "test run: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdTestVerifying(root string, args []string) int {
	fs := flag.NewFlagSet("test verifying", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	symbolID := fs.String("symbol", "", "target symbol UUID")
	filePath := fs.String("file", "", "target file path")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"symbol": true, "file": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace tests verifying --symbol <uuid> | --file <path>\n")
		return exitUsage
	}
	hasSymbol := strings.TrimSpace(*symbolID) != ""
	hasFile := strings.TrimSpace(*filePath) != ""
	if hasSymbol == hasFile {
		fmt.Fprintf(os.Stderr, "usage: trace tests verifying --symbol <uuid> | --file <path>\n")
		return exitUsage
	}

	svc, st, code := openDomain(root)
	if code != exitOK {
		return code
	}
	defer st.Close()
	if code := failCLIDenied(svc, "test", "tests verifying"); code != exitOK {
		return code
	}

	ctx := context.Background()
	var rows []domain.ValidatingTest
	var err error
	if hasSymbol {
		rows, err = svc.ListTestsValidatingSymbol(ctx, *symbolID)
	} else {
		rows, err = svc.ListTestsValidatingFile(ctx, *filePath)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "tests verifying: %v\n", err)
		return exitFail
	}
	if rows == nil {
		rows = []domain.ValidatingTest{}
	}
	if err := json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "tests": rows, "count": len(rows),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "tests verifying: %v\n", err)
		return exitFail
	}
	return exitOK
}
