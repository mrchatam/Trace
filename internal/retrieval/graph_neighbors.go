package retrieval

import (
	"github.com/mrchatam/Trace/internal/store"
)

// graphWalkNeighbor pairs a neighbor hit with the edge that reaches it (links + goal_id).
type graphWalkNeighbor struct {
	neighbor Hit
	edge     GraphEdge
}

// ProvenanceFromSourceType maps entity_links.source_type → wire provenance.
// Empty when unset / unknown (S02 may leave provenance empty).
func ProvenanceFromSourceType(sourceType string) string {
	switch sourceType {
	case "USER_ASSERTED", "IMPORTED", "AGENT_PROPOSED":
		return "explicit"
	case "INFERRED":
		return "inferred"
	default:
		return ""
	}
}

// graphWalkNeighbors returns causal neighbors for graph walking: entity_links plus
// goal↔task via tasks.goal_id (aligned with Expand; goal_id is not stored in entity_links).
func (e *Engine) graphWalkNeighbors(h Hit) ([]graphWalkNeighbor, error) {
	var out []graphWalkNeighbor

	fromLinks, err := e.store.ListLinksFrom(h.EntityType, h.EntityID)
	if err != nil {
		return nil, err
	}
	for _, l := range fromLinks {
		nh, err := e.hitFromLinkNeighbor(l.ToType, l.ToID, l.Rel, 0.5)
		if err != nil {
			if isNotFound(err) {
				continue
			}
			return nil, err
		}
		out = append(out, graphWalkNeighbor{
			neighbor: nh,
			edge: GraphEdge{
				Rel: l.Rel, From: l.FromID, To: l.ToID,
				Provenance: ProvenanceFromSourceType(l.SourceType),
			},
		})
	}

	toLinks, err := e.store.ListLinksTo(h.EntityType, h.EntityID)
	if err != nil {
		return nil, err
	}
	for _, l := range toLinks {
		nh, err := e.hitFromLinkNeighbor(l.FromType, l.FromID, l.Rel, 0.5)
		if err != nil {
			if isNotFound(err) {
				continue
			}
			return nil, err
		}
		out = append(out, graphWalkNeighbor{
			neighbor: nh,
			edge: GraphEdge{
				Rel: l.Rel, From: l.FromID, To: l.ToID,
				Provenance: ProvenanceFromSourceType(l.SourceType),
			},
		})
	}

	switch h.EntityType {
	case "task":
		t, err := e.store.GetTask(h.EntityID)
		if err != nil {
			return nil, err
		}
		if t.GoalID != nil && *t.GoalID != "" {
			gh, err := e.lookupEntity("goal", *t.GoalID, ReasonGoalHasTask, 0, 0.9)
			if err != nil {
				if !isNotFound(err) {
					return nil, err
				}
			} else {
				out = append(out, graphWalkNeighbor{
					neighbor: gh,
					edge:     GraphEdge{Rel: ReasonGoalHasTask, From: *t.GoalID, To: h.EntityID, Provenance: "explicit"},
				})
			}
		}
	case "goal":
		taskPage, err := e.store.ListTasksFiltered(store.TaskListFilter{
			GoalID: h.EntityID, Limit: store.DefaultTaskListLimit,
		})
		if err != nil {
			return nil, err
		}
		for _, t := range taskPage.Tasks {
			out = append(out, graphWalkNeighbor{
				neighbor: Hit{
					EntityType: "task",
					EntityID:   t.ID,
					Title:      t.Title,
					ReasonCode: ReasonGoalHasTask,
					Score:      0.9,
				},
				edge: GraphEdge{Rel: ReasonGoalHasTask, From: h.EntityID, To: t.ID, Provenance: "explicit"},
			})
		}
	}

	return out, nil
}
