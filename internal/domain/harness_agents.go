package domain

import (
	"context"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

const harnessAgentSlugPrefix = "agent:"

// HarnessAgentInput upserts one harness agent and optional requirement slugs.
type HarnessAgentInput struct {
	ID                 string
	Slug               string
	Title              string
	Description        string
	SubagentType       string
	DeliberationPhases string // JSON array
	TaskKeywords       string // JSON array
	RecommendSubagent  bool
	RegistrySource     string
	RegistryVersion    string
	ExternalURL        string
	Requirements       []string // required capability slugs
}

// ValidateHarnessAgentSlug ensures slug uses agent: prefix.
func ValidateHarnessAgentSlug(slug string) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return &ErrValidation{Msg: "harness agent slug is required"}
	}
	if !strings.HasPrefix(slug, harnessAgentSlugPrefix) {
		return &ErrValidation{Msg: "harness agent slug must start with agent:"}
	}
	return nil
}

// UpsertHarnessAgent persists an agent and replaces its requirements.
func (s *Service) UpsertHarnessAgent(ctx context.Context, in HarnessAgentInput) (store.HarnessAgent, error) {
	_ = ctx
	if err := ValidateHarnessAgentSlug(in.Slug); err != nil {
		return store.HarnessAgent{}, err
	}
	row, err := s.store.UpsertHarnessAgent(store.HarnessAgent{
		ID:                 in.ID,
		Slug:               strings.TrimSpace(in.Slug),
		Title:              in.Title,
		Description:        in.Description,
		SubagentType:       in.SubagentType,
		DeliberationPhases: in.DeliberationPhases,
		TaskKeywords:       in.TaskKeywords,
		RecommendSubagent:  in.RecommendSubagent,
		RegistrySource:     in.RegistrySource,
		RegistryVersion:    in.RegistryVersion,
		ExternalURL:        in.ExternalURL,
	})
	if err != nil {
		return store.HarnessAgent{}, err
	}
	if err := s.store.DeleteHarnessAgentRequirementsForAgent(row.ID); err != nil {
		return store.HarnessAgent{}, err
	}
	for _, slug := range in.Requirements {
		slug = strings.TrimSpace(slug)
		if slug == "" {
			continue
		}
		if _, err := s.store.InsertHarnessAgentRequirement(store.HarnessAgentRequirement{
			AgentID:                row.ID,
			RequiredCapabilitySlug: slug,
		}); err != nil {
			return store.HarnessAgent{}, err
		}
	}
	return row, nil
}

// GetHarnessAgent loads one harness agent by id.
func (s *Service) GetHarnessAgent(ctx context.Context, id string) (store.HarnessAgent, error) {
	_ = ctx
	if strings.TrimSpace(id) == "" {
		return store.HarnessAgent{}, &ErrValidation{Msg: "harness agent id is required"}
	}
	return s.store.GetHarnessAgent(id)
}

// GetHarnessAgentBySlug loads one harness agent by slug.
func (s *Service) GetHarnessAgentBySlug(ctx context.Context, slug string) (store.HarnessAgent, error) {
	_ = ctx
	if err := ValidateHarnessAgentSlug(slug); err != nil {
		return store.HarnessAgent{}, err
	}
	return s.store.GetHarnessAgentBySlug(strings.TrimSpace(slug))
}

// ListHarnessAgents returns all harness agents.
func (s *Service) ListHarnessAgents(ctx context.Context) ([]store.HarnessAgent, error) {
	_ = ctx
	return s.store.ListHarnessAgents()
}

// ListHarnessAgentRequirements returns requirements for one agent.
func (s *Service) ListHarnessAgentRequirements(ctx context.Context, agentID string) ([]store.HarnessAgentRequirement, error) {
	_ = ctx
	if strings.TrimSpace(agentID) == "" {
		return nil, &ErrValidation{Msg: "harness agent id is required"}
	}
	return s.store.ListHarnessAgentRequirements(agentID)
}
