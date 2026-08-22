package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

// ImpactInput mirrors `trace impact finding|alternative|report|walk`.
type ImpactInput struct {
	Project     string   `json:"project,omitempty" jsonschema:"optional project root override"`
	Action      string   `json:"action" jsonschema:"finding|alternative|report|walk"`
	Op          string   `json:"op,omitempty" jsonschema:"add|list|recommend (finding/alternative)"`
	Decision    string   `json:"decision,omitempty" jsonschema:"UUID; alias decision_id"`
	DecisionID  string   `json:"decision_id,omitempty" jsonschema:"alias for decision"`
	Class       string   `json:"class,omitempty" jsonschema:"finding add SAFE|CAUTION|HIGH|DESTRUCTIVE|REVERSAL; alias impact_class"`
	ImpactClass string   `json:"impact_class,omitempty" jsonschema:"alias for class"`
	Kind        string   `json:"kind,omitempty" jsonschema:"finding add AFFECTED_WORK|INVALIDATED_ASSUMPTION|WORK_AT_RISK|NEW_WORK|UNRESOLVED"`
	Uncertainty string   `json:"uncertainty,omitempty" jsonschema:"empty → UNKNOWN"`
	Body        string   `json:"body,omitempty" jsonschema:"finding add / alternative add"`
	RelatedType string   `json:"related_type,omitempty" jsonschema:"finding add"`
	RelatedID   string   `json:"related_id,omitempty" jsonschema:"finding add"`
	Title       string   `json:"title,omitempty" jsonschema:"alternative add"`
	Recommended bool     `json:"recommended,omitempty" jsonschema:"alternative add bool"`
	ID          string   `json:"id,omitempty" jsonschema:"alternative recommend; optional on add"`
	Seeds       []string `json:"seeds,omitempty" jsonschema:"walk string array file:<uuid>|symbol:<uuid>"`
	Depth       float64  `json:"depth,omitempty" jsonschema:"walk 1|2; default library DefaultImpactDepth()"`
}

func (in ImpactInput) resolvedDecisionID() string {
	if strings.TrimSpace(in.Decision) != "" {
		return in.Decision
	}
	return strings.TrimSpace(in.DecisionID)
}

func (in ImpactInput) resolvedClass() string {
	if strings.TrimSpace(in.Class) != "" {
		return in.Class
	}
	return strings.TrimSpace(in.ImpactClass)
}

func (s *Server) toolImpact(ctx context.Context, _ *sdkmcp.CallToolRequest, in ImpactInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	if err := assertMCPToolAllowed(ctx, st, "trace_impact"); err != nil {
		_ = st.Close()
		return nil, nil, err
	}
	if err := st.Close(); err != nil {
		return nil, nil, err
	}

	switch strings.ToLower(strings.TrimSpace(in.Action)) {
	case "finding":
		return s.impactFinding(ctx, in)
	case "alternative":
		return s.impactAlternative(ctx, in)
	case "report":
		return s.impactReport(ctx, in)
	case "walk":
		return s.impactWalk(ctx, in)
	default:
		return nil, nil, fmt.Errorf("trace_impact: action must be finding|alternative|report|walk")
	}
}

func (s *Server) impactFinding(ctx context.Context, in ImpactInput) (*sdkmcp.CallToolResult, any, error) {
	switch strings.ToLower(strings.TrimSpace(in.Op)) {
	case "add":
		return s.impactFindingAdd(ctx, in)
	case "list":
		return s.impactFindingList(ctx, in)
	default:
		return nil, nil, fmt.Errorf("trace_impact finding: op must be add|list")
	}
}

func (s *Server) impactFindingAdd(ctx context.Context, in ImpactInput) (*sdkmcp.CallToolResult, any, error) {
	decisionID := in.resolvedDecisionID()
	class := in.resolvedClass()
	if decisionID == "" || class == "" || strings.TrimSpace(in.Kind) == "" {
		return nil, nil, fmt.Errorf("trace_impact finding add: decision, class, and kind are required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	svc := domain.New(st)
	f, err := svc.AddImpactFinding(ctx, decisionID, domain.ImpactFindingInput{
		ID: in.ID, ImpactClass: class, Uncertainty: in.Uncertainty, Kind: in.Kind,
		Body: in.Body, RelatedType: in.RelatedType, RelatedID: in.RelatedID,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_impact: %w", err)
	}
	b, err := json.Marshal(map[string]any{
		"ok": true, "id": f.ID, "decision": f.DecisionID,
		"impact_class": f.ImpactClass, "uncertainty": f.Uncertainty, "kind": f.Kind,
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) impactFindingList(ctx context.Context, in ImpactInput) (*sdkmcp.CallToolResult, any, error) {
	decisionID := in.resolvedDecisionID()
	if decisionID == "" {
		return nil, nil, fmt.Errorf("trace_impact finding list: decision is required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	list, err := domain.New(st).ListImpactFindings(ctx, decisionID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_impact: %w", err)
	}
	if list == nil {
		list = []store.DecisionImpactFinding{}
	}
	b, err := json.Marshal(map[string]any{
		"ok": true, "findings": list, "count": len(list),
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) impactAlternative(ctx context.Context, in ImpactInput) (*sdkmcp.CallToolResult, any, error) {
	switch strings.ToLower(strings.TrimSpace(in.Op)) {
	case "add":
		return s.impactAlternativeAdd(ctx, in)
	case "list":
		return s.impactAlternativeList(ctx, in)
	case "recommend":
		return s.impactAlternativeRecommend(ctx, in)
	default:
		return nil, nil, fmt.Errorf("trace_impact alternative: op must be add|list|recommend")
	}
}

func (s *Server) impactAlternativeAdd(ctx context.Context, in ImpactInput) (*sdkmcp.CallToolResult, any, error) {
	decisionID := in.resolvedDecisionID()
	if decisionID == "" || strings.TrimSpace(in.Title) == "" {
		return nil, nil, fmt.Errorf("trace_impact alternative add: decision and title are required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	a, err := domain.New(st).AddDecisionAlternative(ctx, decisionID, domain.AlternativeInput{
		ID: in.ID, Title: in.Title, Body: in.Body, Recommended: in.Recommended,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("trace_impact: %w", err)
	}
	b, err := json.Marshal(map[string]any{
		"ok": true, "id": a.ID, "decision": a.DecisionID,
		"title": a.Title, "is_recommended": a.IsRecommended,
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) impactAlternativeList(ctx context.Context, in ImpactInput) (*sdkmcp.CallToolResult, any, error) {
	decisionID := in.resolvedDecisionID()
	if decisionID == "" {
		return nil, nil, fmt.Errorf("trace_impact alternative list: decision is required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	list, err := domain.New(st).ListDecisionAlternatives(ctx, decisionID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_impact: %w", err)
	}
	if list == nil {
		list = []store.DecisionAlternative{}
	}
	b, err := json.Marshal(map[string]any{
		"ok": true, "alternatives": list, "count": len(list),
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) impactAlternativeRecommend(ctx context.Context, in ImpactInput) (*sdkmcp.CallToolResult, any, error) {
	decisionID := in.resolvedDecisionID()
	if decisionID == "" || strings.TrimSpace(in.ID) == "" {
		return nil, nil, fmt.Errorf("trace_impact alternative recommend: decision and id are required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	if err := domain.New(st).SetRecommendedAlternative(ctx, decisionID, in.ID); err != nil {
		return nil, nil, fmt.Errorf("trace_impact: %w", err)
	}
	b, err := json.Marshal(map[string]any{
		"ok": true, "decision": decisionID, "id": in.ID, "is_recommended": true,
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) impactReport(ctx context.Context, in ImpactInput) (*sdkmcp.CallToolResult, any, error) {
	decisionID := in.resolvedDecisionID()
	if decisionID == "" {
		return nil, nil, fmt.Errorf("trace_impact report: decision is required")
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	rep, err := domain.New(st).ImpactReport(ctx, decisionID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_impact: %w", err)
	}
	b, err := json.Marshal(map[string]any{
		"ok":                  true,
		"decision_id":         rep.Decision.ID,
		"affected_task_ids":   rep.AffectedTaskIDs,
		"findings":            rep.Findings,
		"alternatives":        rep.Alternatives,
		"overall_class":       rep.OverallClass,
		"overall_uncertainty": rep.OverallUncertainty,
		"has_unknown":         rep.HasUnknown,
		"incomplete":          rep.Incomplete,
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) impactWalk(ctx context.Context, in ImpactInput) (*sdkmcp.CallToolResult, any, error) {
	if len(in.Seeds) == 0 {
		return nil, nil, fmt.Errorf("trace_impact walk: seeds required (file:<uuid>|symbol:<uuid>)")
	}
	walkSeeds := make([]retrieval.ImpactSeed, 0, len(in.Seeds))
	for _, raw := range in.Seeds {
		typ, id, ok := strings.Cut(raw, ":")
		if !ok || typ == "" || id == "" {
			return nil, nil, fmt.Errorf("trace_impact walk: bad seed %q (want file:<uuid> or symbol:<uuid>)", raw)
		}
		typ = strings.ToLower(strings.TrimSpace(typ))
		if typ != "file" && typ != "symbol" {
			return nil, nil, fmt.Errorf("trace_impact walk: seed type must be file|symbol, got %q", typ)
		}
		walkSeeds = append(walkSeeds, retrieval.ImpactSeed{EntityType: typ, EntityID: strings.TrimSpace(id)})
	}
	depth := int(in.Depth)
	if depth == 0 {
		depth = retrieval.DefaultImpactDepth()
	}
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	res, err := retrieval.New(st).ImpactWalk(ctx, walkSeeds, depth)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_impact: %w", err)
	}
	b, err := json.Marshal(map[string]any{
		"ok":             true,
		"seeds":          res.Seeds,
		"blast":          res.Blast,
		"affected_tests": res.AffectedTests,
		"blast_total":    res.BlastTotal,
		"blast_kept":     res.BlastKept,
		"truncated":      res.Truncated,
		"depth":          res.Depth,
	})
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}
