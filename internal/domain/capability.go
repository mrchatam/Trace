package domain

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// Capability kind vocabulary (canonical; matches store.CapabilityKind*).
const (
	CapabilityKindSkill = store.CapabilityKindSkill
	CapabilityKindRule  = store.CapabilityKindRule
	CapabilityKindMCP   = store.CapabilityKindMCP
	CapabilityKindTool  = store.CapabilityKindTool
	CapabilityKindHook  = store.CapabilityKindHook
	CapabilityKindAgent = store.CapabilityKindAgent
)

// Capability status vocabulary.
const (
	CapabilityStatusAvailable   = store.CapabilityStatusAvailable
	CapabilityStatusUnavailable = store.CapabilityStatusUnavailable
	CapabilityStatusUnknown     = store.CapabilityStatusUnknown
)

// CapabilityInput creates or updates a catalog capability.
type CapabilityInput struct {
	ID     string
	Kind   string
	Slug   string
	Title  string
	Status string // empty → UNKNOWN
	Body   string
}

// CapabilitySpec describes a builtin capability without writing to the DB.
type CapabilitySpec struct {
	Kind   string
	Slug   string
	Title  string
	Status string
}

// ListCapabilitiesFilter optionally filters ListCapabilities.
type ListCapabilitiesFilter struct {
	Kind   string
	Status string
}

// NormalizeCapabilityKind returns a valid kind. Empty and unknown fail closed.
func NormalizeCapabilityKind(kind string) (string, error) {
	k := strings.ToUpper(strings.TrimSpace(kind))
	if k == "" {
		return "", &ErrValidation{Msg: "capability kind is required (SKILL|RULE|MCP|TOOL|HOOK|AGENT)"}
	}
	switch k {
	case CapabilityKindSkill, CapabilityKindRule, CapabilityKindMCP, CapabilityKindTool, CapabilityKindHook, CapabilityKindAgent:
		return k, nil
	default:
		return "", &ErrValidation{Msg: "capability kind must be SKILL, RULE, MCP, TOOL, HOOK, or AGENT"}
	}
}

// NormalizeCapabilityStatus returns a valid status. Empty → UNKNOWN; unknown values fail closed.
func NormalizeCapabilityStatus(status string) (string, error) {
	s := strings.ToUpper(strings.TrimSpace(status))
	if s == "" {
		return CapabilityStatusUnknown, nil
	}
	switch s {
	case CapabilityStatusAvailable, CapabilityStatusUnavailable, CapabilityStatusUnknown:
		return s, nil
	default:
		return "", &ErrValidation{Msg: "capability status must be AVAILABLE, UNAVAILABLE, or UNKNOWN"}
	}
}

// BuiltinMCPCapabilitySpecs returns the live MCP tool names as CapabilitySpec
// (kind=MCP, status=AVAILABLE, slug=mcp:<name>). Does NOT write to the DB.
func BuiltinMCPCapabilitySpecs() []CapabilitySpec {
	names := []string{
		"trace_why", "trace_context", "trace_add",
		"trace_link", "trace_transition", "trace_review",
		"trace_tasks", "trace_capability", "trace_impact", "trace_version",
		"trace_search", "trace_changes", "trace_regressions", "trace_loop",
		"trace_agents", "trace_plan", "trace_explore",
	}
	out := make([]CapabilitySpec, 0, len(names))
	for _, n := range names {
		out = append(out, CapabilitySpec{
			Kind:   CapabilityKindMCP,
			Slug:   "mcp:" + n,
			Title:  n,
			Status: CapabilityStatusAvailable,
		})
	}
	return out
}

// BuiltinCLICapabilitySpecs returns gated CLI commands as CapabilitySpec
// (kind=TOOL, status=AVAILABLE, slug=cli:<command>). Does NOT write to the DB.
// Not merged into BuiltinMCPCapabilitySpecs. reindex is an alias of index, not a separate spec.
func BuiltinCLICapabilitySpecs() []CapabilitySpec {
	titles := []string{
		"add", "link", "transition", "review", "why", "context",
		"tasks", "seed", "impact", "plan", "index", "loop", "agents", "changes", "patterns", "knowledge", "search", "explore", "test", "tests", "verify", "eval", "outcomes", "regressions",
	}
	out := make([]CapabilitySpec, 0, len(titles))
	for _, t := range titles {
		out = append(out, CapabilitySpec{
			Kind:   CapabilityKindTool,
			Slug:   "cli:" + t,
			Title:  t,
			Status: CapabilityStatusAvailable,
		})
	}
	return out
}

// UpsertCapability creates or updates a catalog capability.
// Empty ID resolves by slug (DF-41): re-declare without id updates the same row.
// Explicit different-id slug clash still fails at the store UNIQUE(slug) constraint.
func (s *Service) UpsertCapability(ctx context.Context, in CapabilityInput) (store.Capability, error) {
	_ = ctx
	kind, err := NormalizeCapabilityKind(in.Kind)
	if err != nil {
		return store.Capability{}, err
	}
	slug := strings.TrimSpace(in.Slug)
	if slug == "" {
		return store.Capability{}, &ErrValidation{Msg: "capability slug is required"}
	}
	status, err := NormalizeCapabilityStatus(in.Status)
	if err != nil {
		return store.Capability{}, err
	}
	id := strings.TrimSpace(in.ID)
	if id == "" {
		if existing, err := s.store.GetCapabilityBySlug(slug); err == nil {
			id = existing.ID
		} else if !errors.Is(err, sql.ErrNoRows) {
			return store.Capability{}, err
		}
	}
	return s.store.UpsertCapability(store.Capability{
		ID:     id,
		Kind:   kind,
		Slug:   slug,
		Title:  in.Title,
		Status: status,
		Body:   in.Body,
	})
}

// GetCapability loads a capability by id.
func (s *Service) GetCapability(ctx context.Context, id string) (store.Capability, error) {
	_ = ctx
	if id == "" {
		return store.Capability{}, &ErrValidation{Msg: "capability id is required"}
	}
	return s.store.GetCapability(id)
}

// GetCapabilityBySlug loads a capability by slug.
func (s *Service) GetCapabilityBySlug(ctx context.Context, slug string) (store.Capability, error) {
	_ = ctx
	if strings.TrimSpace(slug) == "" {
		return store.Capability{}, &ErrValidation{Msg: "capability slug is required"}
	}
	return s.store.GetCapabilityBySlug(strings.TrimSpace(slug))
}

// ListCapabilities returns catalog entries matching the optional filter (stable by slug).
func (s *Service) ListCapabilities(ctx context.Context, f ListCapabilitiesFilter) ([]store.Capability, error) {
	_ = ctx
	var kind, status string
	if f.Kind != "" {
		k, err := NormalizeCapabilityKind(f.Kind)
		if err != nil {
			return nil, err
		}
		kind = k
	}
	if f.Status != "" {
		st, err := NormalizeCapabilityStatus(f.Status)
		if err != nil {
			return nil, err
		}
		status = st
	}
	return s.store.ListCapabilities(store.CapabilityListFilter{Kind: kind, Status: status})
}

// RequireCapability attaches a capability requirement to a task.
// Task and capability must exist. UNIQUE(task,cap) is idempotent (returns existing).
func (s *Service) RequireCapability(ctx context.Context, taskID, capabilityID string) (store.TaskCapabilityRequirement, error) {
	_ = ctx
	if taskID == "" || capabilityID == "" {
		return store.TaskCapabilityRequirement{}, &ErrValidation{Msg: "taskID and capabilityID are required"}
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return store.TaskCapabilityRequirement{}, err
	}
	if _, err := s.store.GetCapability(capabilityID); err != nil {
		return store.TaskCapabilityRequirement{}, err
	}
	return s.store.InsertTaskCapabilityRequirement(store.TaskCapabilityRequirement{
		TaskID:       taskID,
		CapabilityID: capabilityID,
	})
}

// UnrequireCapability removes a task↔capability requirement.
// Missing row is a no-op (idempotent).
func (s *Service) UnrequireCapability(ctx context.Context, taskID, capabilityID string) error {
	_ = ctx
	if taskID == "" || capabilityID == "" {
		return &ErrValidation{Msg: "taskID and capabilityID are required"}
	}
	return s.store.DeleteTaskCapabilityRequirement(taskID, capabilityID)
}

// ListRequiredCapabilities returns capabilities required by a task (stable by slug).
func (s *Service) ListRequiredCapabilities(ctx context.Context, taskID string) ([]store.Capability, error) {
	_ = ctx
	if taskID == "" {
		return nil, &ErrValidation{Msg: "taskID is required"}
	}
	return s.store.ListCapabilitiesRequiredByTaskID(taskID)
}

// MissingCapabilities returns required capabilities that are absent or status != AVAILABLE.
// Prefer returning Capability rows when present; never silently drop UNKNOWN/UNAVAILABLE.
func (s *Service) MissingCapabilities(ctx context.Context, taskID string) ([]store.Capability, error) {
	_ = ctx
	if taskID == "" {
		return nil, &ErrValidation{Msg: "taskID is required"}
	}
	reqs, err := s.store.ListTaskCapabilityRequirementsByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	var missing []store.Capability
	for _, r := range reqs {
		cap, err := s.store.GetCapability(r.CapabilityID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// Row deleted mid-flight — surface a placeholder as missing.
				missing = append(missing, store.Capability{
					ID:     r.CapabilityID,
					Status: CapabilityStatusUnknown,
					Slug:   r.CapabilityID,
				})
				continue
			}
			return nil, err
		}
		if cap.Status != CapabilityStatusAvailable {
			missing = append(missing, cap)
		}
	}
	return missing, nil
}

// ResolveCapabilityIDOrSlug returns a capability by id, falling back to slug.
func (s *Service) ResolveCapabilityIDOrSlug(ctx context.Context, idOrSlug string) (store.Capability, error) {
	_ = ctx
	idOrSlug = strings.TrimSpace(idOrSlug)
	if idOrSlug == "" {
		return store.Capability{}, &ErrValidation{Msg: "capability id or slug is required"}
	}
	c, err := s.store.GetCapability(idOrSlug)
	if err == nil {
		return c, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		// Try slug anyway for non-UUID lookups that still 404 as ErrNoRows wrapped.
		if c2, err2 := s.store.GetCapabilityBySlug(idOrSlug); err2 == nil {
			return c2, nil
		}
		return store.Capability{}, err
	}
	return s.store.GetCapabilityBySlug(idOrSlug)
}
