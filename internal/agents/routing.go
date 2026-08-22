package agents

import (
	"context"
	"sort"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

const (
	defaultMaxResults  = 4
	subagentCapSlug    = "harness:subagent"
	subagentPromptStub = "Fresh subagent for independent review — not the implementer session."
)

var perfKeywords = []string{"perf", "performance", "latency", "benchmark", "slow", "memory"}
var securityKeywords = []string{"auth", "injection", "owasp", "secret", "xss", "csrf"}

// RecommendInput carries routing signals for deterministic agent recommendations.
type RecommendInput struct {
	Phase        string // deliberation phase e.g. CRITIQUE
	TaskTitle    string
	TaskTags     []string
	GoalKeywords []string
	HarnessCaps  map[string]string // slug → status (AVAILABLE|UNAVAILABLE|UNKNOWN)
}

// Recommendation is one ranked harness agent suggestion (recommend-only; no spawn).
type Recommendation struct {
	AgentSlug           string   `json:"agent_slug"`
	SubagentType        string   `json:"subagent_type"`
	Reason              string   `json:"reason"`
	Confidence          string   `json:"confidence"` // high|medium|low
	UseSubagent         bool     `json:"use_subagent"`
	PromptStub          string   `json:"prompt_stub,omitempty"`
	MissingCapabilities []string `json:"missing_capabilities,omitempty"`
}

type candidate struct {
	slug       string
	score      int
	reason     string
	confidence string
}

// RecommendAgents returns ranked agent suggestions from the catalog (capped at 4).
func RecommendAgents(ctx context.Context, st *store.Store, in RecommendInput) ([]Recommendation, error) {
	_ = ctx
	maxResults := defaultMaxResults

	agents, err := st.ListHarnessAgents()
	if err != nil {
		return nil, err
	}
	if len(agents) == 0 {
		return []Recommendation{}, nil
	}

	bySlug := make(map[string]store.HarnessAgent, len(agents))
	for _, a := range agents {
		bySlug[a.Slug] = a
	}

	cands := collectCandidates(in)
	sort.SliceStable(cands, func(i, j int) bool {
		if cands[i].score != cands[j].score {
			return cands[i].score > cands[j].score
		}
		return cands[i].slug < cands[j].slug
	})

	var out []Recommendation
	for _, c := range cands {
		if len(out) >= maxResults {
			break
		}
		agent, ok := bySlug[c.slug]
		if !ok {
			continue
		}
		rec := Recommendation{
			AgentSlug:    agent.Slug,
			SubagentType: agent.SubagentType,
			Reason:       c.reason,
			Confidence:   c.confidence,
		}
		reqs, err := st.ListHarnessAgentRequirements(agent.ID)
		if err != nil {
			return nil, err
		}
		for _, req := range reqs {
			status := capStatus(in.HarnessCaps, req.RequiredCapabilitySlug)
			if status != store.CapabilityStatusAvailable {
				rec.MissingCapabilities = append(rec.MissingCapabilities, req.RequiredCapabilitySlug)
			}
		}
		if agent.RecommendSubagent && capStatus(in.HarnessCaps, subagentCapSlug) == store.CapabilityStatusAvailable {
			rec.UseSubagent = true
			rec.PromptStub = subagentPromptStub
		}
		out = append(out, rec)
	}
	if out == nil {
		out = []Recommendation{}
	}
	return out, nil
}

func capStatus(caps map[string]string, slug string) string {
	if caps == nil {
		return store.CapabilityStatusUnknown
	}
	s, ok := caps[slug]
	if !ok {
		return store.CapabilityStatusUnknown
	}
	return strings.ToUpper(strings.TrimSpace(s))
}

func collectCandidates(in RecommendInput) []candidate {
	phase := strings.ToUpper(strings.TrimSpace(in.Phase))
	text := buildSearchText(in)

	var cands []candidate
	seen := map[string]bool{}
	add := func(slug string, score int, reason, confidence string) {
		if seen[slug] {
			return
		}
		seen[slug] = true
		cands = append(cands, candidate{slug: slug, score: score, reason: reason, confidence: confidence})
	}

	switch phase {
	case "CRITIQUE":
		add("agent:code-reviewer", 100, "CRITIQUE phase prefers code-reviewer", "high")
		add("agent:nested-reviewer", 90, "CRITIQUE phase secondary nested-reviewer", "medium")
	case "VERIFY":
		if containsKeyword(text, perfKeywords) {
			add("agent:performance-reviewer", 85, "VERIFY with performance keywords", "high")
		}
	}

	if phase == "CRITIQUE" || phase == "VERIFY" {
		if containsKeyword(text, securityKeywords) {
			add("agent:security-reviewer", 80, "security keywords in task context", "medium")
		}
	}

	if phase == "INVESTIGATE" || phase == "ORIENT" {
		add("agent:explore", 70, "investigation or orientation phase", "high")
	}

	if len(cands) == 0 {
		add("agent:generalPurpose", 10, "no phase or keyword match; general fallback", "low")
	}
	return cands
}

func buildSearchText(in RecommendInput) string {
	parts := []string{strings.ToLower(in.TaskTitle)}
	for _, t := range in.TaskTags {
		parts = append(parts, strings.ToLower(t))
	}
	for _, k := range in.GoalKeywords {
		parts = append(parts, strings.ToLower(k))
	}
	return strings.Join(parts, " ")
}

func containsKeyword(text string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}
