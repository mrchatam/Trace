package compiler

import (
	"sort"
	"unicode/utf8"

	"github.com/mrchatam/Trace/internal/retrieval"
)

// estimateTokens is a cheap heuristic (~4 runes per token). No tokenizer dependency.
func estimateTokens(s string) int {
	n := utf8.RuneCountInString(s)
	if n == 0 {
		return 0
	}
	t := (n + 3) / 4
	if t < 1 {
		return 1
	}
	return t
}

func itemTokens(it Item) int {
	return estimateTokens(it.EntityType) + estimateTokens(it.EntityID) +
		estimateTokens(it.Title) + estimateTokens(it.Excerpt) +
		estimateTokens(it.ReasonCode) + 8
}

// trimToBudget enforces MaxItems and TokenLimit. Prefers lower layer (L0→L3),
// then lower distance, then higher score. Sets truncated when anything is dropped.
func trimToBudget(items []Item, tokenLimit, maxItems int) (kept []Item, tokensEst int, truncated bool) {
	if maxItems <= 0 {
		maxItems = DefaultMaxItems
	}
	if maxItems > DefaultMaxItems {
		maxItems = DefaultMaxItems
	}
	if tokenLimit <= 0 {
		tokenLimit = DefaultTokenBudget
	}

	// Stable priority: already expected ordered by caller; we still cap.
	for _, it := range items {
		cost := itemTokens(it)
		if len(kept) >= maxItems {
			truncated = true
			break
		}
		if tokensEst+cost > tokenLimit && len(kept) > 0 {
			// Always keep at least Layer-0 core if possible; if first items exceed alone, still add until max one if empty.
			truncated = true
			break
		}
		kept = append(kept, it)
		tokensEst += cost
	}
	if len(kept) < len(items) {
		truncated = true
	}
	return kept, tokensEst, truncated
}

// sortItemsForBudget groups items by layer (L0→L3) for trim; order within a layer is preserved.
func sortItemsForBudget(items []Item) {
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Layer < items[j].Layer
	})
}

// rankHits sorts hits for compile: exact/direct > fts > graph distance; score desc.
func hitPriority(h retrieval.Hit) int {
	switch h.ReasonCode {
	case retrieval.ReasonExactID, retrieval.ReasonExactPath, retrieval.ReasonExactSymbol,
		retrieval.ReasonDirectTaskScope:
		return 0
	case retrieval.ReasonGoalHasTask, retrieval.ReasonDecisionAffectsTask,
		retrieval.ReasonDiscoveryCausesPlanChg, retrieval.ReasonClaimHasEvidence:
		return 1
	case retrieval.ReasonFTSMatch:
		return 2
	default:
		return 3 + h.Distance
	}
}

func sortHits(hits []retrieval.Hit) {
	// insertion sort for small N — deterministic
	for i := 1; i < len(hits); i++ {
		j := i
		for j > 0 {
			a, b := hits[j-1], hits[j]
			pa, pb := hitPriority(a), hitPriority(b)
			swap := false
			if pa != pb {
				swap = pa > pb
			} else if a.Distance != b.Distance {
				swap = a.Distance > b.Distance
			} else if a.Score != b.Score {
				swap = a.Score < b.Score
			} else if a.EntityType != b.EntityType {
				swap = a.EntityType > b.EntityType
			} else {
				swap = a.EntityID > b.EntityID
			}
			if !swap {
				break
			}
			hits[j-1], hits[j] = hits[j], hits[j-1]
			j--
		}
	}
}
