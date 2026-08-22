package retrieval

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mrchatam/Trace/internal/deliberation"
	"github.com/mrchatam/Trace/internal/domain"
)

// Why walks the causal neighborhood (links + goal_id) and recent events into
// an ordered explanation chain. Every step has a reason_code.
func (e *Engine) Why(ctx context.Context, entityType, entityID string) (WhyResult, error) {
	entityType = NormalizeEntityType(entityType)
	if entityType == "" || entityID == "" {
		return WhyResult{}, fmt.Errorf("retrieval: Why: entityType and entityID required")
	}

	seed, err := e.lookupEntity(entityType, entityID, ReasonExactID, 0, 1.0)
	if err != nil {
		return WhyResult{}, err
	}

	res := WhyResult{
		SeedType:  seed.EntityType,
		SeedID:    entityID,
		Generated: time.Now().UTC(),
		Steps: []WhyStep{{
			EntityType: seed.EntityType,
			EntityID:   seed.EntityID,
			Title:      seed.Title,
			ReasonCode: ReasonExactID,
			Detail:     "seed",
			Distance:   0,
		}},
	}

	expanded, err := e.Expand(ctx, []Hit{seed}, defaultExpandDepth)
	if err != nil {
		return WhyResult{}, err
	}
	for _, h := range expanded {
		if h.EntityType == seed.EntityType && h.EntityID == seed.EntityID {
			continue
		}
		res.Steps = append(res.Steps, WhyStep{
			EntityType:     h.EntityType,
			EntityID:       h.EntityID,
			Title:          h.Title,
			ReasonCode:     h.ReasonCode,
			Detail:         "graph expand",
			Distance:       h.Distance,
			EdgeProvenance: h.EdgeProvenance,
		})
	}

	// Events on seed
	events, err := e.store.ListEventsByEntity(seed.EntityType, entityID)
	if err != nil {
		return WhyResult{}, err
	}
	// Prefer most recent few
	start := 0
	if len(events) > 8 {
		start = len(events) - 8
	}
	for _, ev := range events[start:] {
		if ev.Type == domain.EventDeliberationTransition {
			var payload deliberation.TransitionPayload
			if err := json.Unmarshal([]byte(ev.PayloadJSON), &payload); err == nil {
				res.Steps = append(res.Steps, WhyStep{
					EntityType: ev.EntityType,
					EntityID:   ev.EntityID,
					Title:      string(payload.ToPhase),
					ReasonCode: ReasonDeliberationTransition,
					Detail:     string(payload.ReasonCode),
					Distance:   0,
				})
				continue
			}
		}
		res.Steps = append(res.Steps, WhyStep{
			EntityType: ev.EntityType,
			EntityID:   ev.EntityID,
			Title:      ev.Type,
			ReasonCode: ReasonRecentEvent,
			Detail:     excerpt(ev.PayloadJSON),
			Distance:   0,
		})
	}

	// Optional VCS temporal note for file seeds / path-bearing hits
	if e.vcs != nil {
		path := seed.Path
		if path == "" && seed.EntityType == "file" {
			path = seed.Title
		}
		if path != "" {
			meta, err := e.vcs.LastChanged(ctx, path)
			if err == nil && meta.OID != "" {
				res.Steps = append(res.Steps, WhyStep{
					EntityType: "commit",
					EntityID:   meta.OID,
					Title:      meta.Subject,
					ReasonCode: ReasonHistoricalVCS,
					Detail:     "LastChanged ref only",
					Distance:   0,
				})
			}
		}
	}

	return res, nil
}
