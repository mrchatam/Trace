package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

type createEntityRequest struct {
	Kind       string  `json:"kind"`
	Title      string  `json:"title"`
	Body       string  `json:"body"`
	ID         string  `json:"id"`
	GoalID     string  `json:"goal_id"`
	SourceType string  `json:"source_type"`
	Confidence float64 `json:"confidence"`
	Status     string  `json:"status"`
}

type entitySummary struct {
	ID        string  `json:"id"`
	Kind      string  `json:"kind"`
	Title     string  `json:"title"`
	Body      string  `json:"body,omitempty"`
	WorkState *string `json:"work_state,omitempty"`
}

func (s *Server) handleCreateEntity(w http.ResponseWriter, r *http.Request) {
	var in createEntityRequest
	if err := decodeJSON(r, &in); err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON body", nil)
		return
	}
	if strings.TrimSpace(in.Kind) == "" || strings.TrimSpace(in.Title) == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "kind and title are required", nil)
		return
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := r.Context()

	var out entitySummary
	switch in.Kind {
	case "goal":
		g, err := svc.CreateGoal(ctx, domain.GoalInput{
			ID: in.ID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		out = entitySummary{ID: g.ID, Kind: domain.EntityGoal, Title: g.Title, Body: g.Body}
	case "task":
		var gid *string
		if in.GoalID != "" {
			g := in.GoalID
			gid = &g
		}
		t, err := svc.CreateTask(ctx, domain.TaskInput{
			ID: in.ID, GoalID: gid, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		ws := t.WorkState
		out = entitySummary{ID: t.ID, Kind: domain.EntityTask, Title: t.Title, Body: t.Body, WorkState: &ws}
	case "decision":
		d, err := svc.CreateDecision(ctx, domain.DecisionInput{
			ID: in.ID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		out = entitySummary{ID: d.ID, Kind: domain.EntityDecision, Title: d.Title, Body: d.Body}
	case "assumption":
		a, err := svc.CreateAssumption(ctx, domain.AssumptionInput{
			ID: in.ID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		out = entitySummary{ID: a.ID, Kind: domain.EntityAssumption, Title: a.Title, Body: a.Body}
	case "discovery":
		d, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{
			ID: in.ID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		out = entitySummary{ID: d.ID, Kind: domain.EntityDiscovery, Title: d.Title, Body: d.Body}
	case "plan-change":
		p, err := svc.CreatePlanChange(ctx, domain.PlanChangeInput{
			ID: in.ID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		out = entitySummary{ID: p.ID, Kind: domain.EntityPlanChange, Title: p.Title, Body: p.Body}
	case "claim":
		c, err := svc.CreateClaim(ctx, domain.ClaimInput{
			ID: in.ID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		out = entitySummary{ID: c.ID, Kind: domain.EntityClaim, Title: c.Title, Body: c.Body}
	case "evidence":
		e, err := svc.CreateEvidence(ctx, domain.EvidenceInput{
			ID: in.ID, Title: in.Title, Body: in.Body,
			SourceType: in.SourceType, Confidence: in.Confidence, Status: in.Status,
		})
		if err != nil {
			mapDomainErr(w, err)
			return
		}
		out = entitySummary{ID: e.ID, Kind: domain.EntityEvidence, Title: e.Title, Body: e.Body}
	default:
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "unknown kind", nil)
		return
	}
	writeJSON(w, http.StatusCreated, out)
}

func (s *Server) handleGetEntity(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("entity_id")
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	out, err := lookupEntitySummary(st, id)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func lookupEntitySummary(st *store.Store, id string) (entitySummary, error) {
	try := []func() (entitySummary, error){
		func() (entitySummary, error) {
			t, err := st.GetTask(id)
			if err != nil {
				return entitySummary{}, err
			}
			ws := t.WorkState
			return entitySummary{ID: t.ID, Kind: domain.EntityTask, Title: t.Title, Body: t.Body, WorkState: &ws}, nil
		},
		func() (entitySummary, error) {
			g, err := st.GetGoal(id)
			if err != nil {
				return entitySummary{}, err
			}
			return entitySummary{ID: g.ID, Kind: domain.EntityGoal, Title: g.Title, Body: g.Body}, nil
		},
		func() (entitySummary, error) {
			d, err := st.GetDecision(id)
			if err != nil {
				return entitySummary{}, err
			}
			return entitySummary{ID: d.ID, Kind: domain.EntityDecision, Title: d.Title, Body: d.Body}, nil
		},
		func() (entitySummary, error) {
			a, err := st.GetAssumption(id)
			if err != nil {
				return entitySummary{}, err
			}
			return entitySummary{ID: a.ID, Kind: domain.EntityAssumption, Title: a.Title, Body: a.Body}, nil
		},
		func() (entitySummary, error) {
			d, err := st.GetDiscovery(id)
			if err != nil {
				return entitySummary{}, err
			}
			return entitySummary{ID: d.ID, Kind: domain.EntityDiscovery, Title: d.Title, Body: d.Body}, nil
		},
		func() (entitySummary, error) {
			p, err := st.GetPlanChange(id)
			if err != nil {
				return entitySummary{}, err
			}
			return entitySummary{ID: p.ID, Kind: domain.EntityPlanChange, Title: p.Title, Body: p.Body}, nil
		},
		func() (entitySummary, error) {
			c, err := st.GetClaim(id)
			if err != nil {
				return entitySummary{}, err
			}
			return entitySummary{ID: c.ID, Kind: domain.EntityClaim, Title: c.Title, Body: c.Body}, nil
		},
		func() (entitySummary, error) {
			e, err := st.GetEvidence(id)
			if err != nil {
				return entitySummary{}, err
			}
			return entitySummary{ID: e.ID, Kind: domain.EntityEvidence, Title: e.Title, Body: e.Body}, nil
		},
		func() (entitySummary, error) {
			rv, err := st.GetReview(id)
			if err != nil {
				return entitySummary{}, err
			}
			return entitySummary{ID: rv.ID, Kind: domain.EntityReview, Title: rv.Title, Body: rv.Body}, nil
		},
	}
	var last error
	for _, fn := range try {
		out, err := fn()
		if err == nil {
			return out, nil
		}
		if isNotFoundErr(err) {
			last = err
			continue
		}
		return entitySummary{}, err
	}
	if last == nil {
		last = errors.New("entity not found")
	}
	return entitySummary{}, last
}

type createLinkRequest struct {
	Rel        string `json:"rel"`
	From       string `json:"from"`
	To         string `json:"to"`
	SourceType string `json:"source_type"`
}

func (s *Server) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	var in createLinkRequest
	if err := decodeJSON(r, &in); err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON body", nil)
		return
	}
	if in.Rel == "" || in.From == "" || in.To == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "rel, from, and to are required", nil)
		return
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	svc := domain.New(st)
	meta := domain.LinkMeta{SourceType: in.SourceType}
	ctx := r.Context()
	switch in.Rel {
	case "goal-task":
		err = svc.LinkGoalTask(ctx, in.From, in.To, meta)
	case "decision-task":
		err = svc.LinkDecisionTask(ctx, in.From, in.To, meta)
	case "discovery-plan-change":
		err = svc.LinkDiscoveryPlanChange(ctx, in.From, in.To, meta)
	case "discovery-mentions-task":
		err = svc.LinkDiscoveryMentionsTask(ctx, in.From, in.To, meta)
	case "claim-evidence":
		err = svc.LinkClaimEvidence(ctx, in.From, in.To, meta)
	default:
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "unknown rel", nil)
		return
	}
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"rel": in.Rel, "from": in.From, "to": in.To})
}

type createTransitionRequest struct {
	TaskID           string   `json:"task_id"`
	ToState          string   `json:"to_state"`
	Reason           string   `json:"reason"`
	Actor            string   `json:"actor"`
	AllowDone        bool     `json:"allow_done"`
	AsOperator       bool     `json:"as_operator"`
	AllowMissingCaps bool     `json:"allow_missing_caps"`
	EvidenceIDs      []string `json:"evidence_ids"`
}

func (s *Server) handleCreateTransition(w http.ResponseWriter, r *http.Request) {
	var in createTransitionRequest
	if err := decodeJSON(r, &in); err != nil {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "invalid JSON body", nil)
		return
	}
	if in.TaskID == "" || in.ToState == "" || strings.TrimSpace(in.Reason) == "" {
		writeEnvelope(w, http.StatusBadRequest, "VALIDATION_ERROR", "task_id, to_state, and reason are required", nil)
		return
	}
	actor := in.Actor
	if actor == "" {
		actor = "http"
	}
	st, err := s.openStore()
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	defer st.Close()
	svc := domain.New(st)
	err = svc.TransitionTask(r.Context(), in.TaskID, in.ToState, domain.TransitionOptions{
		Actor: actor, Reason: in.Reason, EvidenceIDs: in.EvidenceIDs,
		AllowDoneWithoutReview: in.AllowDone, AllowOperatorDone: in.AsOperator,
		AllowMissingCapabilities: in.AllowMissingCaps,
	})
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	t, err := st.GetTask(in.TaskID)
	if err != nil {
		mapDomainErr(w, err)
		return
	}
	out := map[string]any{"task_id": t.ID, "work_state": t.WorkState}
	if in.AllowDone {
		out["warning"] = "allow_done escape hatch used; Review PASS and as_operator were bypassed; missing capabilities still need allow_missing_caps"
	}
	writeJSON(w, http.StatusOK, out)
}
