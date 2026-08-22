package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/loop"
	"github.com/mrchatam/Trace/internal/planner"
	"github.com/mrchatam/Trace/internal/store"
)

// resolveSeedEndpoint picks primary (from/to) or alias (from_id/to_id).
// If both are set and differ → usage error.
func resolveSeedEndpoint(primary, alias, key string) (string, error) {
	primary = strings.TrimSpace(primary)
	alias = strings.TrimSpace(alias)
	if primary != "" && alias != "" && primary != alias {
		return "", fmt.Errorf("seed: link %s and %s_id both set and differ", key, key)
	}
	if primary != "" {
		return primary, nil
	}
	return alias, nil
}

var seedImportAllowedKeys = map[string]bool{
	"version": true, "goals": true, "tasks": true, "decisions": true,
	"assumptions": true, "discoveries": true, "plan_changes": true,
	"claims": true, "evidence": true, "links": true, "transitions": true,
	"findings": true, "alternatives": true,
	"plan_phases": true, "plan_scopes": true, "scope_deep_plans": true,
	"goal_plan_state": true, "exported_at_commit": true,
	// P20 cognition (S01 portable seed)
	"deliberation_states": true, "uncertainties": true, "hypotheses": true,
	"decision_reconsiderations": true, "changes": true, "effects": true,
	"outcome_results": true, "baselines": true, "regressions": true, "reflections": true,
}

func cmdSeed(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace seed import <file> | trace seed export [-o <file>] [--strict] [--enforce] [--task <id>]\n")
		return exitUsage
	}
	switch args[0] {
	case "import":
		return cmdSeedImport(root, args[1:])
	case "export":
		return cmdSeedExport(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "usage: trace seed import <file> | trace seed export [-o <file>] [--strict] [--enforce] [--task <id>]\n")
		return exitUsage
	}
}

type exportViolation struct {
	TaskID  string
	Message string
}

func collectExportStructuralViolations(doc domain.SeedDocument) []exportViolation {
	var out []exportViolation
	if doc.Version != 1 {
		out = append(out, exportViolation{
			Message: fmt.Sprintf("version must be 1 (got %d)", doc.Version),
		})
	}
	if doc.Goals == nil {
		out = append(out, exportViolation{Message: "goals must be present"})
	}
	if doc.Tasks == nil {
		out = append(out, exportViolation{Message: "tasks must be present"})
	}
	return out
}

func collectExportViolations(
	ctx context.Context,
	dom *domain.Service,
	plan *planner.Service,
	st *store.Store,
	taskFilter string,
) ([]exportViolation, error) {
	var out []exportViolation
	evaluate := func(taskID string) error {
		allowed, violations, err := loop.EvaluateGate(ctx, dom, plan, st, taskID, loop.GateForExport)
		if err != nil {
			return err
		}
		if !allowed {
			msg := "export gate blocked"
			if len(violations) > 0 {
				msg = violations[0].Message
			}
			out = append(out, exportViolation{TaskID: taskID, Message: msg})
		}
		return nil
	}

	if taskFilter != "" {
		if _, err := st.GetTask(taskFilter); err != nil {
			return nil, fmt.Errorf("task %q not found: %w", taskFilter, err)
		}
		if err := evaluate(taskFilter); err != nil {
			return nil, err
		}
		return out, nil
	}

	tasks, err := st.ListTasks()
	if err != nil {
		return nil, err
	}
	for _, task := range tasks {
		ws := task.WorkState
		if ws == store.WorkStateDone || ws == store.WorkStateSkipped || ws == store.WorkStateStale {
			continue
		}
		if err := evaluate(task.ID); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// collectExportGraphHonestyViolations maps document-level honesty checks only.
// Sole source: domain.CollectSeedDocumentHonestyViolations (no store-backed BLOCKING re-check).
func collectExportGraphHonestyViolations(doc domain.SeedDocument) []exportViolation {
	var out []exportViolation
	for _, v := range domain.CollectSeedDocumentHonestyViolations(doc) {
		out = append(out, exportViolation{Message: v.Message})
	}
	return out
}

func cmdSeedExport(root string, args []string) int {
	fs := flag.NewFlagSet("seed export", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	outPath := fs.String("o", "", "write seed JSON to file (default stdout)")
	strict := fs.Bool("strict", false, "Validate export honesty before write")
	enforce := fs.Bool("enforce", false, "Fail closed on strict violations (requires --strict)")
	taskFilter := fs.String("task", "", "Evaluate export gate for one task only (optional)")
	args = flagsFirst(args, map[string]bool{
		"o": true, "task": true,
	})
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *enforce && !*strict {
		fmt.Fprintf(os.Stderr, "usage: trace seed export [-o <file>] [--strict] [--enforce] [--task <id>]\n")
		fmt.Fprintf(os.Stderr, "seed export: --enforce requires --strict\n")
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace seed export [-o <file>] [--strict] [--enforce] [--task <id>]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	if code := failCLIDenied(svc, "seed", "seed"); code != exitOK {
		return code
	}

	doc, err := domain.BuildSeedDocument(ctx, st, domain.ExportOpts{ProjectRoot: abs})
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed: export: %v\n", err)
		return exitFail
	}

	if !*strict {
		// FM-02 early nudge: warn on thin graph before agents hit --strict --enforce.
		// Warn-only — still writes and exits 0; does not weaken enforce fail-closed.
		for _, v := range collectExportGraphHonestyViolations(doc) {
			if strings.Contains(v.Message, "thin graph") {
				fmt.Fprintf(os.Stderr, "seed export: warn: %s; write discoveries/decisions before seed export --strict --enforce\n", v.Message)
			}
		}
	}

	if *strict {
		violations := collectExportStructuralViolations(doc)
		violations = append(violations, collectExportGraphHonestyViolations(doc)...)
		gateViolations, gateErr := collectExportViolations(ctx, svc, planner.New(st), st, *taskFilter)
		if gateErr != nil {
			fmt.Fprintf(os.Stderr, "seed export strict: %v\n", gateErr)
			return exitFail
		}
		violations = append(violations, gateViolations...)
		for _, v := range violations {
			if v.TaskID != "" {
				fmt.Fprintf(os.Stderr, "seed export strict: task %s: %s\n", v.TaskID, v.Message)
			} else {
				fmt.Fprintf(os.Stderr, "seed export strict: %s\n", v.Message)
			}
		}
		if *enforce && len(violations) > 0 {
			return exitGateBlocked
		}
	}

	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed: encode: %v\n", err)
		return exitFail
	}
	raw = append(raw, '\n')

	if *outPath == "" {
		if _, err := os.Stdout.Write(raw); err != nil {
			fmt.Fprintf(os.Stderr, "seed: write stdout: %v\n", err)
			return exitFail
		}
		return exitOK
	}

	dest, err := resolveSeedPath(abs, *outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed: %v\n", err)
		return exitFail
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "seed: mkdir: %v\n", err)
		return exitFail
	}
	if err := os.WriteFile(dest, raw, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "seed: write: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdSeedImport(root string, args []string) int {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: trace seed import <file>\n")
		return exitUsage
	}
	path, err := resolveSeedPath(root, args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed: %v\n", err)
		return exitFail
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed: %v\n", err)
		return exitFail
	}

	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		fmt.Fprintf(os.Stderr, "seed: parse: %v\n", err)
		return exitFail
	}
	for k := range top {
		if !seedImportAllowedKeys[k] {
			fmt.Fprintf(os.Stderr, "seed: unknown top-level key %q\n", k)
			return exitUsage
		}
	}

	var doc domain.SeedDocument
	if err := json.Unmarshal(raw, &doc); err != nil {
		fmt.Fprintf(os.Stderr, "seed: parse: %v\n", err)
		return exitFail
	}
	if doc.Version != 1 {
		fmt.Fprintf(os.Stderr, "seed: version must be 1 (got %d)\n", doc.Version)
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	if code := failCLIDenied(svc, "seed", "seed"); code != exitOK {
		return code
	}

	for i := range doc.Links {
		from, errFrom := resolveSeedEndpoint(doc.Links[i].From, doc.Links[i].FromID, "from")
		if errFrom != nil {
			fmt.Fprintf(os.Stderr, "%v\n", errFrom)
			return exitUsage
		}
		to, errTo := resolveSeedEndpoint(doc.Links[i].To, doc.Links[i].ToID, "to")
		if errTo != nil {
			fmt.Fprintf(os.Stderr, "%v\n", errTo)
			return exitUsage
		}
		if from == "" || to == "" {
			fmt.Fprintf(os.Stderr, "seed: link endpoints required (from/to or from_id/to_id)\n")
			return exitUsage
		}
		doc.Links[i].From = from
		doc.Links[i].To = to
		doc.Links[i].FromID = ""
		doc.Links[i].ToID = ""
	}

	summary, err := svc.ImportSeedDocument(ctx, doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "seed: import: %v\n", err)
		var val *domain.ErrValidation
		if errors.As(err, &val) {
			return exitUsage
		}
		return exitFail
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(summary)
	return exitOK
}

// resolveSeedPath joins relative import paths under the -C project root.
// Absolute paths are returned unchanged.
func resolveSeedPath(projectRoot, path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	absRoot, err := resolveRoot(projectRoot)
	if err != nil {
		return "", err
	}
	return filepath.Abs(filepath.Join(absRoot, path))
}
