package domain

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

// PromotionCandidate is a BLOCKING discovery with no discovery_mentions_task link.
// Surfaced on seed import and loop next so agents can promote without inventing UUIDs.
type PromotionCandidate struct {
	DiscoveryID string `json:"discovery_id"`
	Title       string `json:"title"`
	Severity    string `json:"severity"`
}

// SeedImportPromotionHint guides agents when promotion_candidates is non-empty.
// Import never auto-spawns (human gate / FR-P28-D1).
const SeedImportPromotionHint = "BLOCKING discoveries listed in promotion_candidates need an explicit promote: trace add task --from-discovery <discovery_id> or loop apply spawned_tasks[].discovery_id (or decline). Do not invent task UUIDs."

const (
	DefaultPromotionCandidateLimit = 32
	MaxPromotionCandidateLimit     = 64
)

// ListPromotionCandidates returns BLOCKING discoveries that lack a live
// discovery_mentions_task target, capped at DefaultPromotionCandidateLimit.
// Prefer ListPromotionCandidatesLimited when callers need truncation honesty.
func (s *Service) ListPromotionCandidates() ([]PromotionCandidate, error) {
	items, _, err := s.ListPromotionCandidatesLimited(DefaultPromotionCandidateLimit)
	return items, err
}

// ListPromotionCandidatesLimited returns up to limit BLOCKING discoveries that
// lack a live discovery_mentions_task target. truncated is true when more exist.
// limit <= 0 uses DefaultPromotionCandidateLimit; values above Max are clamped.
func (s *Service) ListPromotionCandidatesLimited(limit int) ([]PromotionCandidate, bool, error) {
	if limit <= 0 {
		limit = DefaultPromotionCandidateLimit
	}
	if limit > MaxPromotionCandidateLimit {
		limit = MaxPromotionCandidateLimit
	}
	discoveries, err := s.store.ListDiscoveries()
	if err != nil {
		return nil, false, err
	}
	candidates := make([]PromotionCandidate, 0)
	truncated := false
	for _, d := range discoveries {
		if d.Severity != SeverityBlocking {
			continue
		}
		links, err := s.store.ListLinksFrom(EntityDiscovery, d.ID)
		if err != nil {
			return nil, false, err
		}
		linked := false
		for _, l := range links {
			if l.Rel == RelDiscoveryMentionsTask && l.ToType == EntityTask {
				if _, err := s.store.GetTask(l.ToID); err == nil {
					linked = true
					break
				} else if !errors.Is(err, sql.ErrNoRows) {
					return nil, false, err
				}
			}
		}
		if linked {
			continue
		}
		if len(candidates) >= limit {
			truncated = true
			break
		}
		candidates = append(candidates, PromotionCandidate{
			DiscoveryID: d.ID,
			Title:       d.Title,
			Severity:    d.Severity,
		})
	}
	return candidates, truncated, nil
}

// PromoteBlockingDiscovery promotes a BLOCKING discovery into a task and links
// discovery_mentions_task. It is idempotent: if the link already exists, it
// returns the existing task and inserted=false.
func (s *Service) PromoteBlockingDiscovery(ctx context.Context, discoveryID, goalID string) (taskID string, inserted bool, err error) {
	_ = ctx
	discoveryID = strings.TrimSpace(discoveryID)
	if discoveryID == "" {
		return "", false, &ErrValidation{Msg: "discovery_id is required"}
	}
	disc, err := s.store.GetDiscovery(discoveryID)
	if err != nil {
		return "", false, err
	}
	if disc.Severity != SeverityBlocking {
		return "", false, &ErrValidation{Msg: "discovery must have BLOCKING severity"}
	}

	links, err := s.store.ListLinksFrom(EntityDiscovery, discoveryID)
	if err != nil {
		return "", false, err
	}
	for _, link := range links {
		if link.Rel != RelDiscoveryMentionsTask || link.ToType != EntityTask {
			continue
		}
		if _, err := s.store.GetTask(link.ToID); err == nil {
			return link.ToID, false, nil
		} else if !errors.Is(err, sql.ErrNoRows) {
			return "", false, err
		}
	}

	goalID = strings.TrimSpace(goalID)
	if goalID == "" {
		return "", false, &ErrValidation{Msg: "goal_id is required when discovery is not yet linked"}
	}

	task, inserted, err := s.ImportSeedTask(ctx, SeedTask{
		ID:     disc.ID,
		GoalID: goalID,
		Title:  strings.TrimSpace(disc.Title),
	}, &goalID)
	if err != nil {
		return "", false, err
	}
	if err := s.LinkDiscoveryMentionsTask(ctx, disc.ID, task.ID, LinkMeta{}); err != nil {
		return "", false, err
	}
	return task.ID, inserted, nil
}
