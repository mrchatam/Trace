package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

// AgentListItem is the summary shape for trace agents list / trace_agents list.
type AgentListItem struct {
	Slug               string   `json:"slug"`
	Title              string   `json:"title"`
	SubagentType       string   `json:"subagent_type"`
	DeliberationPhases []string `json:"deliberation_phases"`
	Requirements       []string `json:"requirements"`
}

// AgentDescribe is the full profile for trace agents describe / trace_agents describe.
type AgentDescribe struct {
	ID                 string   `json:"id"`
	Slug               string   `json:"slug"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	SubagentType       string   `json:"subagent_type"`
	DeliberationPhases []string `json:"deliberation_phases"`
	TaskKeywords       []string `json:"task_keywords"`
	RecommendSubagent  bool     `json:"recommend_subagent"`
	RegistrySource     string   `json:"registry_source"`
	RegistryVersion    string   `json:"registry_version"`
	ExternalURL        string   `json:"external_url,omitempty"`
	Requirements       []string `json:"requirements"`
	CreatedAt          string   `json:"created_at"`
	UpdatedAt          string   `json:"updated_at"`
}

// ListAgentSummaries returns catalog rows ordered by slug (empty catalog → []).
func ListAgentSummaries(ctx context.Context, st *store.Store) ([]AgentListItem, error) {
	_ = ctx
	rows, err := st.ListHarnessAgents()
	if err != nil {
		return nil, err
	}
	out := make([]AgentListItem, 0, len(rows))
	for _, a := range rows {
		phases, err := parseStringJSONArray(a.DeliberationPhases)
		if err != nil {
			return nil, fmt.Errorf("agents list: %s deliberation_phases: %w", a.Slug, err)
		}
		reqSlugs, err := requirementSlugs(st, a.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, AgentListItem{
			Slug:               a.Slug,
			Title:              a.Title,
			SubagentType:       a.SubagentType,
			DeliberationPhases: phases,
			Requirements:       reqSlugs,
		})
	}
	if out == nil {
		out = []AgentListItem{}
	}
	return out, nil
}

// DescribeAgent returns one full profile by slug.
func DescribeAgent(ctx context.Context, st *store.Store, slug string) (AgentDescribe, error) {
	_ = ctx
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return AgentDescribe{}, fmt.Errorf("agents describe: slug is required")
	}
	a, err := st.GetHarnessAgentBySlug(slug)
	if err != nil {
		return AgentDescribe{}, fmt.Errorf("agents describe: unknown slug %q", slug)
	}
	phases, err := parseStringJSONArray(a.DeliberationPhases)
	if err != nil {
		return AgentDescribe{}, fmt.Errorf("agents describe: %s deliberation_phases: %w", a.Slug, err)
	}
	keywords, err := parseStringJSONArray(a.TaskKeywords)
	if err != nil {
		return AgentDescribe{}, fmt.Errorf("agents describe: %s task_keywords: %w", a.Slug, err)
	}
	reqSlugs, err := requirementSlugs(st, a.ID)
	if err != nil {
		return AgentDescribe{}, err
	}
	return AgentDescribe{
		ID:                 a.ID,
		Slug:               a.Slug,
		Title:              a.Title,
		Description:        a.Description,
		SubagentType:       a.SubagentType,
		DeliberationPhases: phases,
		TaskKeywords:       keywords,
		RecommendSubagent:  a.RecommendSubagent,
		RegistrySource:     a.RegistrySource,
		RegistryVersion:    a.RegistryVersion,
		ExternalURL:        a.ExternalURL,
		Requirements:       reqSlugs,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}, nil
}

func requirementSlugs(st *store.Store, agentID string) ([]string, error) {
	reqs, err := st.ListHarnessAgentRequirements(agentID)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(reqs))
	for _, r := range reqs {
		out = append(out, r.RequiredCapabilitySlug)
	}
	if out == nil {
		out = []string{}
	}
	return out, nil
}

func parseStringJSONArray(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "[]" {
		return []string{}, nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = []string{}
	}
	return out, nil
}
