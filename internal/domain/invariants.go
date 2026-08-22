package domain

import (
	"context"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

const (
	// RuleInternalMustNotImportCmd is the single default architectural invariant (S07 extends rules).
	RuleInternalMustNotImportCmd = "internal_must_not_import_cmd"
)

// InvariantViolation is one failed architectural check on a change-path import.
type InvariantViolation struct {
	Rule         string `json:"rule"`
	ImporterPath string `json:"importer_path"`
	ImportedPath string `json:"imported_path"`
	FromLayer    string `json:"from_layer"`
	ToLayer      string `json:"to_layer"`
}

// InvariantCheckResult is the advisory JSON from CheckArchitecturalInvariants.
type InvariantCheckResult struct {
	Passed     bool                 `json:"passed"`
	Violations []InvariantViolation `json:"violations"`
}

// CheckArchitecturalInvariants applies the default forbidden rule (internal/ must
// not import cmd/) to indexed imports on paths from the task's latest
// RECORDED/COMPARED change. It does not glob the repo (Law 12).
func (s *Service) CheckArchitecturalInvariants(ctx context.Context, taskID string) (InvariantCheckResult, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return InvariantCheckResult{}, &ErrValidation{Msg: "task_id is required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return InvariantCheckResult{}, err
	}
	out := InvariantCheckResult{Passed: true, Violations: []InvariantViolation{}}
	changes, err := s.store.ListChangesByTaskID(taskID)
	if err != nil {
		return InvariantCheckResult{}, err
	}
	latest, ok := latestRecordedOrComparedChange(changes)
	if !ok {
		return out, nil
	}
	paths, err := s.store.ListChangePaths(latest.ID)
	if err != nil {
		return InvariantCheckResult{}, err
	}
	touched := make(map[string]struct{}, len(paths))
	for _, p := range paths {
		touched[store.NormalizePath(p.Path)] = struct{}{}
	}
	if len(touched) == 0 {
		return out, nil
	}

	cross, err := s.store.ListCrossLayerImports()
	if err != nil {
		return InvariantCheckResult{}, err
	}
	for _, im := range cross {
		importer := store.NormalizePath(im.ImporterPath)
		if _, ok := touched[importer]; !ok {
			continue
		}
		if im.FromLayer != "internal" || im.ToLayer != "cmd" {
			continue
		}
		out.Violations = append(out.Violations, InvariantViolation{
			Rule:         RuleInternalMustNotImportCmd,
			ImporterPath: importer,
			ImportedPath: im.ImportedPath,
			FromLayer:    im.FromLayer,
			ToLayer:      im.ToLayer,
		})
	}
	out.Passed = len(out.Violations) == 0
	return out, nil
}
