package retrieval

import (
	"strings"
	"time"
)

// Locked reason codes (exact strings).
const (
	ReasonExactID                = "exact_id"
	ReasonExactPath              = "exact_path"
	ReasonExactSymbol            = "exact_symbol"
	ReasonFTSMatch               = "fts_match"
	ReasonGraphLabelMatch        = "graph_label_match"
	ReasonDirectTaskScope        = "direct_task_scope"
	ReasonGoalHasTask            = "goal_has_task"
	ReasonDecisionAffectsTask    = "decision_affects_task"
	ReasonDiscoveryCausesPlanChg = "discovery_causes_plan_change"
	ReasonClaimHasEvidence       = "claim_has_evidence"
	ReasonReviewJudgesTask       = "review_judges_task"
	ReasonReviewJudgesScope      = "review_judges_scope"
	ReasonGraphNeighbor          = "graph_neighbor"
	ReasonRecentEvent            = "recent_event"
	ReasonDeliberationTransition = "deliberation_transition"
	ReasonHistoricalVCS          = "historical_vcs"
	ReasonBudgetDropped          = "budget_dropped" // debug/drops only — never on included items
)

// Hit is one retrieval candidate with provenance.
type Hit struct {
	EntityType     string  `json:"entity_type"`
	EntityID       string  `json:"entity_id"`
	Title          string  `json:"title,omitempty"`
	Excerpt        string  `json:"excerpt,omitempty"`
	Path           string  `json:"path,omitempty"`
	ReasonCode     string  `json:"reason_code"`
	Score          float64 `json:"score,omitempty"`
	Distance       int     `json:"distance"`
	EdgeProvenance string  `json:"edge_provenance,omitempty"` // EXTRACTED|INFERRED|AMBIGUOUS on structural import hops
}

// WhyStep is one link/event in a causal explanation chain.
type WhyStep struct {
	EntityType     string `json:"entity_type"`
	EntityID       string `json:"entity_id"`
	Title          string `json:"title,omitempty"`
	ReasonCode     string `json:"reason_code"`
	Detail         string `json:"detail,omitempty"`
	Distance       int    `json:"distance"`
	EdgeProvenance string `json:"edge_provenance,omitempty"`
}

// ExactQuery selects exact lookup dimensions. Misses return empty slices, not fabricated hits.
type ExactQuery struct {
	EntityType string // goal|task|decision|review|… when EntityID set
	EntityID   string
	Path       string // file path
	SymbolName string // with Path, or unambiguous single symbol name
}

// SearchOptions controls FTS Search. Limit default ≤ 32; hard-capped at 64.
// When Intent is non-nil, Search builds the FTS query from ExtractIntent (G9);
// nil preserves legacy raw-q behavior for MCP direct search.
type SearchOptions struct {
	Limit  int
	Intent *IntentInput
}

// WhyResult is the ordered causal explanation for an entity.
type WhyResult struct {
	SeedType  string    `json:"seed_type"`
	SeedID    string    `json:"seed_id"`
	Steps     []WhyStep `json:"steps"`
	Generated time.Time `json:"generated_at"`
}

// NormalizeEntityType maps CLI/MCP aliases to canonical store/JSON types (DF-23).
// Emitted hits and Why steps always use the canonical form (e.g. plan_change).
func NormalizeEntityType(entityType string) string {
	switch strings.TrimSpace(entityType) {
	case "plan-change":
		return "plan_change"
	case "outcome":
		return "outcome_result"
	default:
		return strings.TrimSpace(entityType)
	}
}

// hitKey dedupes hits by type+id.
func hitKey(entityType, entityID string) string {
	return entityType + "\x00" + entityID
}
