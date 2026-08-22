package planner

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// Event type for human ack of auto-replan churn budget.
const EventReplanAcked = "plan.replan_acked"

// ApplyDiscoveryReplan applies a discovery against a scope's deep plan with severity
// and churn gates. INFO may link a PlanChange but never supersedes or increments.
// PLAN_AFFECTING and BLOCKING auto-open replan via SupersedeDeepPlan + LinkDiscoveryPlanChange.
func (s *Service) ApplyDiscoveryReplan(ctx context.Context, in ApplyDiscoveryReplanInput) (ApplyDiscoveryReplanResult, error) {
	discID := strings.TrimSpace(in.DiscoveryID)
	scopeID := strings.TrimSpace(in.ScopeID)
	if discID == "" {
		return ApplyDiscoveryReplanResult{}, &ErrValidation{Msg: "discovery_id is required"}
	}
	if scopeID == "" {
		return ApplyDiscoveryReplanResult{}, &ErrValidation{Msg: "scope_id is required"}
	}

	disc, err := s.store.GetDiscovery(discID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ApplyDiscoveryReplanResult{}, fmt.Errorf("planner: discovery %q: %w", discID, ErrNotFound)
		}
		return ApplyDiscoveryReplanResult{}, err
	}

	sev := strings.TrimSpace(disc.Severity)
	if sev == "" {
		sev = SeverityINFO
	}
	switch sev {
	case SeverityINFO:
		return s.applyInfoDiscovery(ctx, in, disc, scopeID)
	case SeverityPlanAffecting, SeverityBlocking:
		return s.applyAutoReplanDiscovery(ctx, in, disc, scopeID, sev)
	default:
		return ApplyDiscoveryReplanResult{}, &ErrValidation{
			Msg: "discovery severity must be INFO, PLAN_AFFECTING, or BLOCKING",
		}
	}
}

func (s *Service) applyInfoDiscovery(ctx context.Context, in ApplyDiscoveryReplanInput, disc store.Discovery, scopeID string) (ApplyDiscoveryReplanResult, error) {
	out := ApplyDiscoveryReplanResult{
		DiscoveryID:       disc.ID,
		ScopeID:           scopeID,
		AutoReplanApplied: false,
		Reason:            "severity_info",
	}
	sc, err := s.store.GetPlanScope(scopeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ApplyDiscoveryReplanResult{}, fmt.Errorf("planner: scope %q: %w", scopeID, ErrNotFound)
		}
		return ApplyDiscoveryReplanResult{}, err
	}
	out.AutoReplanCount = sc.AutoReplanCount

	pcID, linked, err := s.ensurePlanChangeLink(ctx, in, disc, false)
	if err != nil {
		return ApplyDiscoveryReplanResult{}, err
	}
	if linked {
		out.PlanChangeID = pcID
	}
	return out, nil
}

func (s *Service) applyAutoReplanDiscovery(ctx context.Context, in ApplyDiscoveryReplanInput, disc store.Discovery, scopeID, sev string) (ApplyDiscoveryReplanResult, error) {
	maxN := in.MaxAutoReplans
	if maxN <= 0 {
		maxN = DefaultMaxAutoReplans
	}

	sc, err := s.store.GetPlanScope(scopeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ApplyDiscoveryReplanResult{}, fmt.Errorf("planner: scope %q: %w", scopeID, ErrNotFound)
		}
		return ApplyDiscoveryReplanResult{}, err
	}
	if sc.AutoReplanCount >= maxN {
		return ApplyDiscoveryReplanResult{}, fmt.Errorf(
			"%w: count=%d max=%d scope=%s", ErrReplanBudgetExceeded, sc.AutoReplanCount, maxN, scopeID,
		)
	}

	pcID, _, err := s.ensurePlanChangeLink(ctx, in, disc, true)
	if err != nil {
		return ApplyDiscoveryReplanResult{}, err
	}

	sup, err := s.SupersedeDeepPlan(ctx, SupersedeInput{
		ScopeID:          scopeID,
		ExitCriteria:     in.ExitCriteria,
		Constraints:      in.Constraints,
		WorkItems:        in.WorkItems,
		LookaheadScopeID: in.LookaheadScopeID,
		LookaheadSummary: in.LookaheadSummary,
		Actor:            in.Actor,
	})
	if err != nil {
		return ApplyDiscoveryReplanResult{}, err
	}

	newCount, err := s.store.IncrementAutoReplanCount(scopeID)
	if err != nil {
		return ApplyDiscoveryReplanResult{}, err
	}

	return ApplyDiscoveryReplanResult{
		DiscoveryID:       disc.ID,
		ScopeID:           scopeID,
		PlanChangeID:      pcID,
		AutoReplanApplied: true,
		Reason:            "severity_" + strings.ToLower(sev),
		RevisionID:        sup.RevisionID,
		AutoReplanCount:   newCount,
		SupersededCount:   sup.SupersededCount,
	}, nil
}

// ensurePlanChangeLink creates or reuses a PlanChange and links discovery→plan_change.
// When requireCreate is false (INFO), linking is skipped unless PlanChangeID or title/body is provided.
func (s *Service) ensurePlanChangeLink(ctx context.Context, in ApplyDiscoveryReplanInput, disc store.Discovery, requireCreate bool) (planChangeID string, linked bool, err error) {
	dom := domain.New(s.store)
	pcID := strings.TrimSpace(in.PlanChangeID)
	title := strings.TrimSpace(in.PlanChangeTitle)
	body := in.PlanChangeBody

	if pcID == "" {
		if !requireCreate && title == "" && strings.TrimSpace(body) == "" {
			return "", false, nil
		}
		if title == "" {
			title = disc.Title
			if title == "" {
				title = "Plan change from discovery"
			}
		}
		if body == "" {
			body = disc.Body
		}
		pc, err := dom.CreatePlanChange(ctx, domain.PlanChangeInput{Title: title, Body: body})
		if err != nil {
			return "", false, err
		}
		pcID = pc.ID
	} else {
		if _, err := s.store.GetPlanChange(pcID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", false, fmt.Errorf("planner: plan_change %q: %w", pcID, ErrNotFound)
			}
			return "", false, err
		}
	}

	if err := dom.LinkDiscoveryPlanChange(ctx, disc.ID, pcID, domain.LinkMeta{}); err != nil {
		return "", false, err
	}
	return pcID, true, nil
}

// AckReplan resets auto_replan_count to 0 for a scope (human ack) and appends a thin event.
func (s *Service) AckReplan(ctx context.Context, scopeID string) error {
	_ = ctx
	scopeID = strings.TrimSpace(scopeID)
	if scopeID == "" {
		return &ErrValidation{Msg: "scope_id is required"}
	}
	if _, err := s.store.GetPlanScope(scopeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("planner: scope %q: %w", scopeID, ErrNotFound)
		}
		return err
	}
	if err := s.store.AckAutoReplan(scopeID); err != nil {
		return err
	}
	payload, _ := json.Marshal(map[string]string{"scope_id": scopeID, "actor": "ack"})
	_, _ = s.store.AppendEvent(store.Event{
		Type: EventReplanAcked, EntityType: entityScope, EntityID: scopeID,
		PayloadJSON: string(payload),
	})
	return nil
}
