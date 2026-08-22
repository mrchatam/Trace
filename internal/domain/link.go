package domain

import (
	"context"
	"encoding/json"

	"github.com/mrchatam/Trace/internal/store"
)

func (s *Service) appendLinked(fromType, fromID, rel, toType, toID string, meta LinkMeta) error {
	payload, _ := json.Marshal(map[string]any{
		"rel":       rel,
		"to_type":   toType,
		"to_id":     toID,
		"from_type": fromType,
		"from_id":   fromID,
	})
	_, err := s.store.AppendEvent(store.Event{
		Type:        EventEntityLinked,
		EntityType:  fromType,
		EntityID:    fromID,
		PayloadJSON: string(payload),
	})
	return err
}

// LinkGoalTask sets tasks.goal_id and appends entity.linked (rel goal_has_task in payload only).
func (s *Service) LinkGoalTask(ctx context.Context, goalID, taskID string, meta LinkMeta) error {
	_ = ctx
	if goalID == "" || taskID == "" {
		return &ErrValidation{Msg: "goalID and taskID are required"}
	}
	if _, err := s.store.GetGoal(goalID); err != nil {
		return err
	}
	task, err := s.store.GetTask(taskID)
	if err != nil {
		return err
	}
	gid := goalID
	task.GoalID = &gid
	if _, err := s.store.UpsertTask(task); err != nil {
		return err
	}
	meta = meta.withDefaults()
	_ = meta // provenance for Goal→Task is on the task row; event records the link
	return s.appendLinked(EntityGoal, goalID, RelGoalHasTaskEvent, EntityTask, taskID, meta)
}

// LinkDecisionTask inserts entity_links rel=decision_affects_task and appends entity.linked.
func (s *Service) LinkDecisionTask(ctx context.Context, decisionID, taskID string, meta LinkMeta) error {
	_ = ctx
	if decisionID == "" || taskID == "" {
		return &ErrValidation{Msg: "decisionID and taskID are required"}
	}
	if _, err := s.store.GetDecision(decisionID); err != nil {
		return err
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return err
	}
	meta = meta.withDefaults()
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   EntityDecision,
		FromID:     decisionID,
		Rel:        RelDecisionAffectsTask,
		ToType:     EntityTask,
		ToID:       taskID,
		SourceType: meta.SourceType,
		Confidence: meta.Confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(EntityDecision, decisionID, RelDecisionAffectsTask, EntityTask, taskID, meta)
}

// LinkDiscoveryMentionsTask inserts entity_links rel=discovery_mentions_task (DF-42).
// Used for multi-goal DPC attribution (discovery↔task endpoint).
func (s *Service) LinkDiscoveryMentionsTask(ctx context.Context, discoveryID, taskID string, meta LinkMeta) error {
	_ = ctx
	if discoveryID == "" || taskID == "" {
		return &ErrValidation{Msg: "discoveryID and taskID are required"}
	}
	if _, err := s.store.GetDiscovery(discoveryID); err != nil {
		return err
	}
	if _, err := s.store.GetTask(taskID); err != nil {
		return err
	}
	meta = meta.withDefaults()
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   EntityDiscovery,
		FromID:     discoveryID,
		Rel:        RelDiscoveryMentionsTask,
		ToType:     EntityTask,
		ToID:       taskID,
		SourceType: meta.SourceType,
		Confidence: meta.Confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(EntityDiscovery, discoveryID, RelDiscoveryMentionsTask, EntityTask, taskID, meta)
}

// LinkDiscoveryPlanChange inserts entity_links rel=discovery_causes_plan_change.
func (s *Service) LinkDiscoveryPlanChange(ctx context.Context, discoveryID, planChangeID string, meta LinkMeta) error {
	_ = ctx
	if discoveryID == "" || planChangeID == "" {
		return &ErrValidation{Msg: "discoveryID and planChangeID are required"}
	}
	if _, err := s.store.GetDiscovery(discoveryID); err != nil {
		return err
	}
	if _, err := s.store.GetPlanChange(planChangeID); err != nil {
		return err
	}
	meta = meta.withDefaults()
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   EntityDiscovery,
		FromID:     discoveryID,
		Rel:        RelDiscoveryCausesPlanChange,
		ToType:     EntityPlanChange,
		ToID:       planChangeID,
		SourceType: meta.SourceType,
		Confidence: meta.Confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(EntityDiscovery, discoveryID, RelDiscoveryCausesPlanChange, EntityPlanChange, planChangeID, meta)
}

// ListLinksFrom returns entity_links from an entity.
func (s *Service) ListLinksFrom(ctx context.Context, entityType, entityID string) ([]store.EntityLink, error) {
	_ = ctx
	return s.store.ListLinksFrom(entityType, entityID)
}
