package eval

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/eval/rulectx"
	"github.com/mrchatam/Trace/internal/store"
)

const rulesFileVersion = 1

const evalRulesRelPath = "trace/eval-rules.json"

// RulesPath returns the committed eval-rules file path under project root.
func RulesPath(root string) string {
	return filepath.Join(root, evalRulesRelPath)
}

// InvariantRule toggles one architectural invariant for the eval mechanism.
type InvariantRule struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

// RulesFile is the v1 eval-rules.json schema.
type RulesFile struct {
	Version    int             `json:"version"`
	Mechanisms []string        `json:"mechanisms"`
	Invariants []InvariantRule `json:"invariants"`
}

// RulesLoadResult is the outcome of LoadRules for CLI and tests.
type RulesLoadResult struct {
	Path       string          `json:"path"`
	Loaded     bool            `json:"loaded"`
	Rules      RulesFile       `json:"-"`
	Mechanisms []string        `json:"mechanisms"`
	Invariants []InvariantRule `json:"invariants"`
	CachedAt   string          `json:"cached_at,omitempty"`
}

// DefaultRules returns built-in mechanism ids and the default invariant enabled.
func DefaultRules() RulesFile {
	return RulesFile{
		Version: rulesFileVersion,
		Mechanisms: []string{
			MechanismStoredTest,
			MechanismStoredVerification,
			MechanismStoredEvaluation,
			MechanismArchitecturalInvariant,
		},
		Invariants: []InvariantRule{
			{ID: domain.RuleInternalMustNotImportCmd, Enabled: true},
		},
	}
}

// IsInvariantEnabled reports whether an invariant id runs (default true when absent).
func (r *RulesFile) IsInvariantEnabled(id string) bool {
	for _, inv := range r.Invariants {
		if inv.ID == id {
			return inv.Enabled
		}
	}
	return true
}

// FilterRegisteredMechanisms preserves file order and skips unknown ids.
func (r *RulesFile) FilterRegisteredMechanisms(registered map[string]struct{}) []string {
	out := make([]string, 0, len(r.Mechanisms))
	for _, id := range r.Mechanisms {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := registered[id]; ok {
			out = append(out, id)
		} else {
			log.Printf("eval: skip unknown mechanism id %q in eval-rules.json", id)
		}
	}
	return out
}

// LoadRules reads trace/eval-rules.json when present; missing file returns defaults.
// Successful parse upserts eval_rule_sets id=default. Invalid JSON or version fails closed.
func LoadRules(ctx context.Context, root string, st *store.Store) (RulesLoadResult, error) {
	_ = ctx
	path := RulesPath(root)
	res := RulesLoadResult{Path: evalRulesRelPath}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			rules := DefaultRules()
			res.Rules = rules
			res.Mechanisms = append([]string(nil), rules.Mechanisms...)
			res.Invariants = append([]InvariantRule(nil), rules.Invariants...)
			return res, nil
		}
		return RulesLoadResult{}, fmt.Errorf("eval: read rules: %w", err)
	}

	var parsed RulesFile
	if err := json.Unmarshal(data, &parsed); err != nil {
		return RulesLoadResult{}, &domain.ErrValidation{Msg: "eval-rules.json: invalid JSON"}
	}
	if parsed.Version != rulesFileVersion {
		return RulesLoadResult{}, &domain.ErrValidation{Msg: fmt.Sprintf("eval-rules.json: unsupported version %d (want %d)", parsed.Version, rulesFileVersion)}
	}

	bodyJSON, err := json.Marshal(parsed)
	if err != nil {
		return RulesLoadResult{}, err
	}
	row, err := st.UpsertEvalRuleSet(store.EvalRuleSet{
		ID:         store.EvalRuleSetDefaultID,
		SourcePath: evalRulesRelPath,
		BodyJSON:   string(bodyJSON),
	})
	if err != nil {
		return RulesLoadResult{}, err
	}

	res.Loaded = true
	res.Rules = parsed
	res.Mechanisms = append([]string(nil), parsed.Mechanisms...)
	res.Invariants = append([]InvariantRule(nil), parsed.Invariants...)
	res.CachedAt = row.UpdatedAt
	return res, nil
}

// WithRules attaches eval rules to ctx for built-in mechanisms.
func WithRules(ctx context.Context, rules *RulesFile) context.Context {
	return rulectx.With(ctx, rules)
}

// RulesFromContext returns eval rules previously attached with WithRules.
func RulesFromContext(ctx context.Context) *RulesFile {
	v := rulectx.From(ctx)
	if v == nil {
		return nil
	}
	r, ok := v.(*RulesFile)
	if !ok {
		return nil
	}
	return r
}
