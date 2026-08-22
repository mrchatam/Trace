package domain

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// Tool decision statuses (canonical).
const (
	ToolDecisionAutoAllowed = store.ToolDecisionAutoAllowed
	ToolDecisionPending     = store.ToolDecisionPending
	ToolDecisionAllowed     = store.ToolDecisionAllowed
	ToolDecisionDenied      = store.ToolDecisionDenied
)

// ToolDecision is the resolved allowlist outcome for a tool slug.
type ToolDecision struct {
	Slug     string
	Decision string
	Reason   string
	Actor    string
	Durable  bool // true when backed by a store row
}

// NormalizeToolDecision returns a valid human/auto status. Empty/unknown fail closed.
func NormalizeToolDecision(decision string) (string, error) {
	d := strings.ToUpper(strings.TrimSpace(decision))
	if d == "" {
		return "", &ErrValidation{Msg: "tool decision is required (ALLOWED|DENIED)"}
	}
	switch d {
	case ToolDecisionAutoAllowed, ToolDecisionPending, ToolDecisionAllowed, ToolDecisionDenied:
		return d, nil
	default:
		return "", &ErrValidation{Msg: "tool decision must be AUTO_ALLOWED, PENDING, ALLOWED, or DENIED"}
	}
}

// isBuiltinMCPSlug reports exact match against BuiltinMCPCapabilitySpecs() slugs only.
func isBuiltinMCPSlug(slug string) bool {
	slug = strings.TrimSpace(slug)
	for _, spec := range BuiltinMCPCapabilitySpecs() {
		if spec.Slug == slug {
			return true
		}
	}
	return false
}

// isBuiltinCLISlug reports exact match against BuiltinCLICapabilitySpecs() slugs only.
func isBuiltinCLISlug(slug string) bool {
	slug = strings.TrimSpace(slug)
	for _, spec := range BuiltinCLICapabilitySpecs() {
		if spec.Slug == slug {
			return true
		}
	}
	return false
}

// canonicalizeToolSlug maps an exact registered MCP Name or mcp:Name to mcp:Name.
// Exact match only (no globs, no case-fold). Also folds cli:reindex → cli:index.
// Unprefixed add/why/index/reindex and other cli: slugs are unchanged.
func canonicalizeToolSlug(slug string) string {
	slug = strings.TrimSpace(slug)
	for _, spec := range BuiltinMCPCapabilitySpecs() {
		if slug == spec.Title || slug == spec.Slug {
			return spec.Slug
		}
	}
	if slug == "cli:reindex" {
		return "cli:index"
	}
	return slug
}

// ResolveToolDecision returns the durable or graduated decision for a slug.
// Persisted ALLOWED/DENIED win; else exact builtin MCP slug → persist AUTO_ALLOWED;
// else PENDING (not an error by itself; no durable PENDING required).
// Unknown persisted statuses fail closed as PENDING and never upsert AUTO_ALLOWED.
func (s *Service) ResolveToolDecision(ctx context.Context, slug string) (ToolDecision, error) {
	_ = ctx
	slug = canonicalizeToolSlug(slug)
	if slug == "" {
		return ToolDecision{}, &ErrValidation{Msg: "tool slug is required"}
	}

	existing, err := s.store.GetCapabilityToolDecisionBySlug(slug)
	if err == nil {
		switch existing.Decision {
		case ToolDecisionAllowed, ToolDecisionDenied, ToolDecisionAutoAllowed, ToolDecisionPending:
			return ToolDecision{
				Slug:     existing.Slug,
				Decision: existing.Decision,
				Reason:   existing.Reason,
				Actor:    existing.Actor,
				Durable:  true,
			}, nil
		default:
			return ToolDecision{
				Slug:     existing.Slug,
				Decision: ToolDecisionPending,
				Reason:   existing.Reason,
				Actor:    existing.Actor,
				Durable:  true,
			}, nil
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return ToolDecision{}, err
	}

	if isBuiltinMCPSlug(slug) || isBuiltinCLISlug(slug) {
		reason := "builtin MCP capability"
		if isBuiltinCLISlug(slug) {
			reason = "builtin CLI command"
		}
		row, err := s.store.UpsertCapabilityToolDecision(store.CapabilityToolDecision{
			Slug:     slug,
			Decision: ToolDecisionAutoAllowed,
			Reason:   reason,
			Actor:    "system",
		})
		if err != nil {
			return ToolDecision{}, err
		}
		return ToolDecision{
			Slug:     row.Slug,
			Decision: row.Decision,
			Reason:   row.Reason,
			Actor:    row.Actor,
			Durable:  true,
		}, nil
	}

	return ToolDecision{
		Slug:     slug,
		Decision: ToolDecisionPending,
		Reason:   "unknown tool; awaiting human decision",
		Durable:  false,
	}, nil
}

// AssertToolAllowed fails closed when the resolved decision is PENDING or DENIED.
func (s *Service) AssertToolAllowed(ctx context.Context, slug string) error {
	d, err := s.ResolveToolDecision(ctx, slug)
	if err != nil {
		return err
	}
	switch d.Decision {
	case ToolDecisionAutoAllowed, ToolDecisionAllowed:
		return nil
	case ToolDecisionDenied:
		return &ErrValidation{Msg: "tool " + slug + " is DENIED"}
	case ToolDecisionPending:
		return &ErrValidation{Msg: "tool " + slug + " is PENDING (fail-closed)"}
	default:
		return &ErrValidation{Msg: "tool " + slug + " is not allowed"}
	}
}

// DecideToolInput is a human ALLOWED|DENIED decision.
type DecideToolInput struct {
	Slug     string
	Decision string // ALLOWED|DENIED only
	Reason   string
	Actor    string // default cli
}

// DecideTool persists a human ALLOWED or DENIED decision (actor default "cli").
func (s *Service) DecideTool(ctx context.Context, in DecideToolInput) (store.CapabilityToolDecision, error) {
	_ = ctx
	slug := canonicalizeToolSlug(in.Slug)
	if slug == "" {
		return store.CapabilityToolDecision{}, &ErrValidation{Msg: "tool slug is required"}
	}
	decision, err := NormalizeToolDecision(in.Decision)
	if err != nil {
		return store.CapabilityToolDecision{}, err
	}
	if decision != ToolDecisionAllowed && decision != ToolDecisionDenied {
		return store.CapabilityToolDecision{}, &ErrValidation{Msg: "human decide accepts ALLOWED or DENIED only"}
	}
	actor := strings.TrimSpace(in.Actor)
	if actor == "" {
		actor = "cli"
	}
	return s.store.UpsertCapabilityToolDecision(store.CapabilityToolDecision{
		Slug:     slug,
		Decision: decision,
		Reason:   in.Reason,
		Actor:    actor,
	})
}

// ListToolDecisions returns durable audit rows.
func (s *Service) ListToolDecisions(ctx context.Context) ([]store.CapabilityToolDecision, error) {
	_ = ctx
	return s.store.ListCapabilityToolDecisions()
}
