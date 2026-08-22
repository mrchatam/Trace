package retrieval

import (
	"github.com/mrchatam/Trace/internal/store"
)

// discoveryPlanChangeHitsForTask returns discovery↔plan_change endpoints scoped to
// a task Expand seed (DF-19). Replaces global all-project attach.
//
// Rules (planner locks):
//  1. Pair-completion: if either DPC endpoint is already in alreadyHits, include both.
//  2. Goal-scope: for task with goal_id=G, include a pair iff it is not foreign to G
//     (foreign = either endpoint linked to a task whose goal_id is set and ≠ G).
//  3. Unattributed (neither endpoint linked to any task): include only when the store
//     has exactly one goal and that goal is G.
//  4. Task with nil goal_id: pair-completion only.
func (e *Engine) discoveryPlanChangeHitsForTask(taskID string, alreadyHits map[string]Hit) ([]Hit, error) {
	task, err := e.store.GetTask(taskID)
	if err != nil {
		return nil, err
	}
	var goalID *string
	if task.GoalID != nil && *task.GoalID != "" {
		goalID = task.GoalID
	}

	links, err := e.store.ListLinksByRel(ReasonDiscoveryCausesPlanChg)
	if err != nil {
		return nil, err
	}

	var goalCount int
	var soleGoalID string
	if goalID != nil {
		goals, err := e.store.ListGoals()
		if err != nil {
			return nil, err
		}
		goalCount = len(goals)
		if goalCount == 1 {
			soleGoalID = goals[0].ID
		}
	}

	out := make([]Hit, 0)
	seen := map[string]struct{}{}
	addSide := func(typ, id string) error {
		k := hitKey(typ, id)
		if _, ok := seen[k]; ok {
			return nil
		}
		seen[k] = struct{}{}
		h, err := e.lookupEntity(typ, id, ReasonDiscoveryCausesPlanChg, 0, 0.85)
		if err != nil {
			if isNotFound(err) {
				return nil
			}
			return err
		}
		out = append(out, h)
		return nil
	}

	for _, l := range links {
		fromK := hitKey(l.FromType, l.FromID)
		toK := hitKey(l.ToType, l.ToID)
		_, fromPresent := alreadyHits[fromK]
		_, toPresent := alreadyHits[toK]
		if fromPresent || toPresent {
			if err := addSide(l.FromType, l.FromID); err != nil {
				return nil, err
			}
			if err := addSide(l.ToType, l.ToID); err != nil {
				return nil, err
			}
			continue
		}

		if goalID == nil {
			continue
		}
		G := *goalID

		attributed, foreign, err := e.dpcAttribution(l, G)
		if err != nil {
			return nil, err
		}
		if attributed {
			if foreign {
				continue
			}
			if err := addSide(l.FromType, l.FromID); err != nil {
				return nil, err
			}
			if err := addSide(l.ToType, l.ToID); err != nil {
				return nil, err
			}
			continue
		}

		// Unattributed: only single-goal store matching G.
		if goalCount == 1 && soleGoalID == G {
			if err := addSide(l.FromType, l.FromID); err != nil {
				return nil, err
			}
			if err := addSide(l.ToType, l.ToID); err != nil {
				return nil, err
			}
		}
	}
	return out, nil
}

// dpcAttribution reports whether either DPC endpoint is linked to any task, and
// whether either endpoint is linked to a task with goal_id set and ≠ G.
func (e *Engine) dpcAttribution(l store.EntityLink, G string) (attributed, foreign bool, err error) {
	for _, side := range []struct{ typ, id string }{
		{l.FromType, l.FromID},
		{l.ToType, l.ToID},
	} {
		attr, forgn, err := e.endpointTaskGoals(side.typ, side.id, G)
		if err != nil {
			return false, false, err
		}
		if attr {
			attributed = true
		}
		if forgn {
			foreign = true
		}
	}
	return attributed, foreign, nil
}

func (e *Engine) endpointTaskGoals(entityType, entityID, G string) (linkedToTask, foreign bool, err error) {
	fromLinks, err := e.store.ListLinksFrom(entityType, entityID)
	if err != nil {
		return false, false, err
	}
	toLinks, err := e.store.ListLinksTo(entityType, entityID)
	if err != nil {
		return false, false, err
	}
	check := func(otherType, otherID string) error {
		if otherType != "task" {
			return nil
		}
		linkedToTask = true
		t, err := e.store.GetTask(otherID)
		if err != nil {
			if isNotFound(err) {
				return nil
			}
			return err
		}
		if t.GoalID != nil && *t.GoalID != "" && *t.GoalID != G {
			foreign = true
		}
		return nil
	}
	for _, l := range fromLinks {
		if err := check(l.ToType, l.ToID); err != nil {
			return false, false, err
		}
	}
	for _, l := range toLinks {
		if err := check(l.FromType, l.FromID); err != nil {
			return false, false, err
		}
	}
	return linkedToTask, foreign, nil
}
