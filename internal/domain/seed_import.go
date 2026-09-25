package domain

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

// SeedImportSummary is the import result for CLI JSON output.
type SeedImportSummary struct {
	OK                  bool                 `json:"ok"`
	Created             map[string][]string  `json:"created"`
	Links               int                  `json:"links"`
	Findings            int                  `json:"findings"`
	Alternatives        int                  `json:"alternatives"`
	Transitions         int                  `json:"transitions"`
	PromotionCandidates          []PromotionCandidate `json:"promotion_candidates"`
	PromotionCandidatesTruncated bool                 `json:"promotion_candidates_truncated,omitempty"`
	PromotionHint               string               `json:"promotion_hint,omitempty"`
}

// ImportSeedDocument idempotently imports seed JSON v1 (DF-81/83/84).
func (s *Service) ImportSeedDocument(ctx context.Context, doc SeedDocument) (SeedImportSummary, error) {
	summary := SeedImportSummary{
		OK:                  true,
		Created:             map[string][]string{},
		PromotionCandidates: []PromotionCandidate{},
	}
	err := s.store.WithTx(func(stx *store.Store) error {
		return s.withStore(stx).importSeedDocumentBody(ctx, doc, &summary)
	})
	if err != nil {
		summary.OK = false
		summary.Created = map[string][]string{}
		summary.Links = 0
		summary.Findings = 0
		summary.Alternatives = 0
		summary.Transitions = 0
		summary.PromotionCandidates = []PromotionCandidate{}
		summary.PromotionCandidatesTruncated = false
		summary.PromotionHint = ""
		return summary, err
	}
	candidates, truncated, err := s.ListPromotionCandidatesLimited(DefaultPromotionCandidateLimit)
	if err != nil {
		summary.OK = false
		return summary, err
	}
	summary.PromotionCandidates = candidates
	summary.PromotionCandidatesTruncated = truncated
	if len(candidates) > 0 {
		summary.PromotionHint = SeedImportPromotionHint
	}
	return summary, nil
}

func (s *Service) importSeedDocumentBody(ctx context.Context, doc SeedDocument, summary *SeedImportSummary) error {
	addCreated := func(kind, id string, inserted bool) {
		if inserted {
			summary.Created[kind] = append(summary.Created[kind], id)
		}
	}

	for _, g := range doc.Goals {
		ent, inserted, err := s.ImportSeedGoal(ctx, g)
		if err != nil {
			return err
		}
		addCreated(EntityGoal, ent.ID, inserted)
	}
	for _, t := range doc.Tasks {
		var gid *string
		if t.GoalID != "" {
			g := t.GoalID
			gid = &g
		}
		ent, inserted, err := s.ImportSeedTask(ctx, SeedTask{
			ID: t.ID, GoalID: t.GoalID, Title: t.Title, Body: t.Body,
		}, gid)
		if err != nil {
			return err
		}
		addCreated(EntityTask, ent.ID, inserted)
	}
	for _, d := range doc.Decisions {
		ent, inserted, err := s.ImportSeedDecision(ctx, d)
		if err != nil {
			return err
		}
		addCreated(EntityDecision, ent.ID, inserted)
	}
	for _, a := range doc.Assumptions {
		ent, inserted, err := s.ImportSeedAssumption(ctx, a)
		if err != nil {
			return err
		}
		addCreated(EntityAssumption, ent.ID, inserted)
	}
	for _, d := range doc.Discoveries {
		ent, inserted, err := s.ImportSeedDiscovery(ctx, d)
		if err != nil {
			return err
		}
		addCreated(EntityDiscovery, ent.ID, inserted)
	}
	for _, p := range doc.PlanChanges {
		ent, inserted, err := s.ImportSeedPlanChange(ctx, p)
		if err != nil {
			return err
		}
		addCreated(EntityPlanChange, ent.ID, inserted)
	}
	for _, c := range doc.Claims {
		ent, inserted, err := s.ImportSeedClaim(ctx, c)
		if err != nil {
			return err
		}
		addCreated(EntityClaim, ent.ID, inserted)
	}
	for _, e := range doc.Evidence {
		ent, inserted, err := s.ImportSeedEvidence(ctx, e)
		if err != nil {
			return err
		}
		addCreated(EntityEvidence, ent.ID, inserted)
	}

	for _, sc := range doc.Scopes {
		ent, inserted, err := s.ImportSeedScope(ctx, sc)
		if err != nil {
			return err
		}
		addCreated(EntityScope, ent.ID, inserted)
	}

	for _, l := range doc.Links {
		if err := s.ImportSeedLink(ctx, l); err != nil {
			return err
		}
		summary.Links++
	}

	for _, f := range doc.Findings {
		if strings.TrimSpace(f.DecisionID) == "" {
			return &ErrValidation{Msg: "finding decision_id required"}
		}
		if _, err := s.ImportSeedFinding(ctx, f.DecisionID, f); err != nil {
			return err
		}
		summary.Findings++
	}
	for _, a := range doc.Alternatives {
		if strings.TrimSpace(a.DecisionID) == "" {
			return &ErrValidation{Msg: "alternative decision_id required"}
		}
		if _, err := s.ImportSeedAlternative(ctx, a.DecisionID, a); err != nil {
			return err
		}
		summary.Alternatives++
	}

	for _, p := range doc.PlanPhases {
		if _, err := s.ImportSeedPlanPhase(ctx, p); err != nil {
			return err
		}
	}
	for _, sc := range doc.PlanScopes {
		if _, err := s.ImportSeedPlanScope(ctx, sc); err != nil {
			return err
		}
	}
	for _, d := range doc.ScopeDeepPlans {
		if _, err := s.ImportSeedScopeDeepPlan(ctx, d); err != nil {
			return err
		}
	}
	for _, g := range doc.GoalPlanState {
		if _, err := s.ImportSeedGoalPlanState(ctx, g); err != nil {
			return err
		}
	}

	for _, b := range doc.Baselines {
		if _, err := s.ImportSeedBaseline(ctx, b); err != nil {
			return err
		}
	}
	for _, o := range doc.OutcomeResults {
		if _, err := s.ImportSeedOutcomeResult(ctx, o); err != nil {
			return err
		}
	}
	for _, c := range doc.Changes {
		if _, err := s.ImportSeedChange(ctx, c); err != nil {
			return err
		}
	}
	for _, e := range doc.Effects {
		if _, err := s.ImportSeedEffect(ctx, e); err != nil {
			return err
		}
	}
	for _, u := range doc.Uncertainties {
		if _, err := s.ImportSeedUncertainty(ctx, u); err != nil {
			return err
		}
	}
	for _, h := range doc.Hypotheses {
		if _, err := s.ImportSeedHypothesis(ctx, h); err != nil {
			return err
		}
	}
	for _, r := range doc.DecisionReconsiderations {
		if _, err := s.ImportSeedDecisionReconsideration(ctx, r); err != nil {
			return err
		}
	}
	for _, r := range doc.Regressions {
		if _, err := s.ImportSeedRegression(ctx, r); err != nil {
			return err
		}
	}
	for _, im := range doc.Improvements {
		if _, err := s.ImportSeedImprovement(ctx, im); err != nil {
			return err
		}
	}
	for _, r := range doc.Reflections {
		if _, err := s.ImportSeedReflection(ctx, r); err != nil {
			return err
		}
	}
	for _, p := range doc.ChangePatterns {
		if _, err := s.ImportSeedChangePattern(ctx, p); err != nil {
			return err
		}
	}
	for _, k := range doc.EngineeringKnowledge {
		if _, err := s.ImportSeedEngineeringKnowledge(ctx, k); err != nil {
			return err
		}
	}
	for _, a := range doc.HarnessAgents {
		if _, err := s.ImportSeedHarnessAgent(ctx, a); err != nil {
			return err
		}
	}
	if doc.EvalRulesPath != "" {
		if _, err := s.store.UpsertEvalRuleSet(store.EvalRuleSet{
			ID:         store.EvalRuleSetDefaultID,
			SourcePath: doc.EvalRulesPath,
			BodyJSON:   "{}",
		}); err != nil {
			return err
		}
	}
	for _, ds := range doc.DeliberationStates {
		if _, err := s.ImportSeedDeliberationState(ctx, ds); err != nil {
			return err
		}
	}

	for _, tr := range doc.Transitions {
		if strings.TrimSpace(tr.Reason) == "" {
			return &ErrValidation{Msg: "transition reason required"}
		}
		if err := s.ImportSeedTransition(ctx, tr); err != nil {
			return err
		}
		summary.Transitions++
	}

	return nil
}


func seedEntityExists(s *store.Store, entityType, id string) (bool, error) {
	if id == "" {
		return false, nil
	}
	var err error
	switch entityType {
	case EntityGoal:
		_, err = s.GetGoal(id)
	case EntityTask:
		_, err = s.GetTask(id)
	case EntityDecision:
		_, err = s.GetDecision(id)
	case EntityAssumption:
		_, err = s.GetAssumption(id)
	case EntityDiscovery:
		_, err = s.GetDiscovery(id)
	case EntityPlanChange:
		_, err = s.GetPlanChange(id)
	case EntityClaim:
		_, err = s.GetClaim(id)
	case EntityEvidence:
		_, err = s.GetEvidence(id)
	default:
		return false, &ErrValidation{Msg: "unknown entity type " + entityType}
	}
	if err == nil {
		return true, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return false, err
}

func (s *Service) upsertEntityCreated(entityType, entityID, title string, existed bool) error {
	if existed {
		return nil
	}
	return s.appendCreated(entityType, entityID, title)
}

// ImportSeedGoal upserts a goal; entity.created only on first insert.
func (s *Service) ImportSeedGoal(ctx context.Context, in SeedEntity) (store.Goal, bool, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, "", "")
	if err != nil {
		return store.Goal{}, false, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	existed, err := seedEntityExists(s.store, EntityGoal, id)
	if err != nil {
		return store.Goal{}, false, err
	}
	g, err := s.store.UpsertGoal(store.Goal{
		ID: id, Title: strings.TrimSpace(in.Title), Body: in.Body,
		SourceType: src, Status: status,
	})
	if err != nil {
		return store.Goal{}, false, err
	}
	if err := s.upsertEntityCreated(EntityGoal, g.ID, g.Title, existed); err != nil {
		return store.Goal{}, false, err
	}
	return g, !existed, nil
}

// ImportSeedTask upserts a task preserving local work_state on conflict.
func (s *Service) ImportSeedTask(ctx context.Context, in SeedTask, goalID *string) (store.Task, bool, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, "", "")
	if err != nil {
		return store.Task{}, false, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	existed, err := seedEntityExists(s.store, EntityTask, id)
	if err != nil {
		return store.Task{}, false, err
	}
	t, err := s.store.UpsertTaskFromSeed(store.Task{
		ID: id, GoalID: goalID, Title: strings.TrimSpace(in.Title), Body: in.Body,
		SourceType: src, Status: status, WorkState: store.WorkStatePending,
	})
	if err != nil {
		return store.Task{}, false, err
	}
	if err := s.upsertEntityCreated(EntityTask, t.ID, t.Title, existed); err != nil {
		return store.Task{}, false, err
	}
	return t, !existed, nil
}

// ImportSeedDecision upserts a decision; entity.created only on first insert.
func (s *Service) ImportSeedDecision(ctx context.Context, in SeedEntity) (store.Decision, bool, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, "", "")
	if err != nil {
		return store.Decision{}, false, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	existed, err := seedEntityExists(s.store, EntityDecision, id)
	if err != nil {
		return store.Decision{}, false, err
	}
	d, err := s.store.UpsertDecision(store.Decision{
		ID: id, Title: strings.TrimSpace(in.Title), Body: in.Body,
		SourceType: src, Status: status,
	})
	if err != nil {
		return store.Decision{}, false, err
	}
	if err := s.upsertEntityCreated(EntityDecision, d.ID, d.Title, existed); err != nil {
		return store.Decision{}, false, err
	}
	return d, !existed, nil
}

// ImportSeedAssumption upserts an assumption; entity.created only on first insert.
func (s *Service) ImportSeedAssumption(ctx context.Context, in SeedEntity) (store.Assumption, bool, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, "", "")
	if err != nil {
		return store.Assumption{}, false, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	existed, err := seedEntityExists(s.store, EntityAssumption, id)
	if err != nil {
		return store.Assumption{}, false, err
	}
	a, err := s.store.UpsertAssumption(store.Assumption{
		ID: id, Title: strings.TrimSpace(in.Title), Body: in.Body,
		SourceType: src, Status: status,
	})
	if err != nil {
		return store.Assumption{}, false, err
	}
	if err := s.upsertEntityCreated(EntityAssumption, a.ID, a.Title, existed); err != nil {
		return store.Assumption{}, false, err
	}
	return a, !existed, nil
}

// ImportSeedDiscovery upserts a discovery; entity.created only on first insert.
func (s *Service) ImportSeedDiscovery(ctx context.Context, in SeedEntity) (store.Discovery, bool, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, "", "")
	if err != nil {
		return store.Discovery{}, false, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	existed, err := seedEntityExists(s.store, EntityDiscovery, id)
	if err != nil {
		return store.Discovery{}, false, err
	}
	d, err := s.store.UpsertDiscovery(store.Discovery{
		ID: id, Title: strings.TrimSpace(in.Title), Body: in.Body,
		SourceType: src, Status: status, Severity: store.SeverityINFO,
	})
	if err != nil {
		return store.Discovery{}, false, err
	}
	if err := s.upsertEntityCreated(EntityDiscovery, d.ID, d.Title, existed); err != nil {
		return store.Discovery{}, false, err
	}
	return d, !existed, nil
}

// ImportSeedPlanChange upserts a plan_change; entity.created only on first insert.
func (s *Service) ImportSeedPlanChange(ctx context.Context, in SeedEntity) (store.PlanChange, bool, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, "", "")
	if err != nil {
		return store.PlanChange{}, false, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	existed, err := seedEntityExists(s.store, EntityPlanChange, id)
	if err != nil {
		return store.PlanChange{}, false, err
	}
	p, err := s.store.UpsertPlanChange(store.PlanChange{
		ID: id, Title: strings.TrimSpace(in.Title), Body: in.Body,
		SourceType: src, Status: status,
	})
	if err != nil {
		return store.PlanChange{}, false, err
	}
	if err := s.upsertEntityCreated(EntityPlanChange, p.ID, p.Title, existed); err != nil {
		return store.PlanChange{}, false, err
	}
	return p, !existed, nil
}

// ImportSeedClaim upserts a claim; entity.created only on first insert.
func (s *Service) ImportSeedClaim(ctx context.Context, in SeedEntity) (store.Claim, bool, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, "", "")
	if err != nil {
		return store.Claim{}, false, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	existed, err := seedEntityExists(s.store, EntityClaim, id)
	if err != nil {
		return store.Claim{}, false, err
	}
	c, err := s.store.UpsertClaim(store.Claim{
		ID: id, Title: strings.TrimSpace(in.Title), Body: in.Body,
		SourceType: src, Status: status,
	})
	if err != nil {
		return store.Claim{}, false, err
	}
	if err := s.upsertEntityCreated(EntityClaim, c.ID, c.Title, existed); err != nil {
		return store.Claim{}, false, err
	}
	return c, !existed, nil
}

// ImportSeedEvidence upserts evidence; entity.created only on first insert.
func (s *Service) ImportSeedEvidence(ctx context.Context, in SeedEntity) (store.Evidence, bool, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, "", "")
	if err != nil {
		return store.Evidence{}, false, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	existed, err := seedEntityExists(s.store, EntityEvidence, id)
	if err != nil {
		return store.Evidence{}, false, err
	}
	e, err := s.store.UpsertEvidence(store.Evidence{
		ID: id, Title: strings.TrimSpace(in.Title), Body: in.Body,
		SourceType: src, Status: status,
	})
	if err != nil {
		return store.Evidence{}, false, err
	}
	if err := s.upsertEntityCreated(EntityEvidence, e.ID, e.Title, existed); err != nil {
		return store.Evidence{}, false, err
	}
	return e, !existed, nil
}

// ImportSeedLink idempotently imports a seed link (duplicate endpoints no-op).
func (s *Service) ImportSeedLink(ctx context.Context, l SeedLink) error {
	_ = ctx
	from := strings.TrimSpace(l.From)
	if from == "" {
		from = strings.TrimSpace(l.FromID)
	}
	to := strings.TrimSpace(l.To)
	if to == "" {
		to = strings.TrimSpace(l.ToID)
	}
	if from == "" || to == "" {
		return &ErrValidation{Msg: "link endpoints required (from/to or from_id/to_id)"}
	}
	meta := LinkMeta{}
	if st := strings.TrimSpace(l.SourceType); st != "" {
		meta.SourceType = st
	}
	if l.Confidence != 0 {
		meta.Confidence = l.Confidence
	}
	switch l.Rel {
	case RelGoalHasTaskEvent, "goal-task":
		return s.importSeedGoalTaskLink(ctx, from, to, meta)
	case RelDecisionAffectsTask, "decision-task":
		return s.importSeedEntityLink(ctx, EntityDecision, from, RelDecisionAffectsTask, EntityTask, to, meta)
	case RelDiscoveryCausesPlanChange, "discovery-plan-change":
		return s.importSeedEntityLink(ctx, EntityDiscovery, from, RelDiscoveryCausesPlanChange, EntityPlanChange, to, meta)
	case RelClaimHasEvidence, "claim-evidence":
		return s.importSeedEntityLink(ctx, EntityClaim, from, RelClaimHasEvidence, EntityEvidence, to, meta)
	case RelDiscoveryMentionsTask, "discovery-mentions-task":
		return s.importSeedEntityLink(ctx, EntityDiscovery, from, RelDiscoveryMentionsTask, EntityTask, to, meta)
	case RelRegressionAssociatedChange, "regression-associated-change":
		return s.importSeedRegressionChangeLink(ctx, from, to, meta)
	case RelScopeMember, "scope-member":
		return s.importSeedScopeMemberLink(ctx, from, to, meta)
	case RelAPIContract, "api-contract":
		return s.importSeedEntityLink(ctx, EntityTask, from, RelAPIContract, EntityTask, to, meta)
	case RelImplements:
		return s.importSeedImplementsLink(ctx, from, to, meta)
	case RelBlocks:
		return s.importSeedEntityLink(ctx, EntityTask, from, RelBlocks, EntityTask, to, meta)
	default:
		return &ErrValidation{Msg: "unknown link rel " + l.Rel}
	}
}

func (s *Service) importSeedGoalTaskLink(ctx context.Context, goalID, taskID string, meta LinkMeta) error {
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
	if task.GoalID != nil && *task.GoalID == goalID {
		return nil
	}
	gid := goalID
	task.GoalID = &gid
	if _, err := s.store.UpsertTaskFromSeed(task); err != nil {
		return err
	}
	meta = meta.withDefaults()
	return s.appendLinked(EntityGoal, goalID, RelGoalHasTaskEvent, EntityTask, taskID, meta)
}

func (s *Service) importSeedRegressionChangeLink(ctx context.Context, regressionID, changeID string, meta LinkMeta) error {
	_ = ctx
	if regressionID == "" || changeID == "" {
		return &ErrValidation{Msg: "regression and change ids are required"}
	}
	if _, err := s.store.GetRegression(regressionID); err != nil {
		return err
	}
	if _, err := s.store.GetChange(changeID); err != nil {
		return err
	}
	return s.importSeedEntityLink(ctx, EntityRegression, regressionID, RelRegressionAssociatedChange, EntityChange, changeID, meta)
}

func (s *Service) importSeedEntityLink(ctx context.Context, fromType, fromID, rel, toType, toID string, meta LinkMeta) error {
	_ = ctx
	if fromID == "" || toID == "" {
		return &ErrValidation{Msg: "link endpoints required"}
	}
	if err := s.validateLinkEndpoints(fromType, fromID, toType, toID); err != nil {
		return err
	}
	meta = meta.withDefaults()
	inserted, _, err := s.store.InsertLinkOrIgnore(store.EntityLink{
		FromType: fromType, FromID: fromID, Rel: rel,
		ToType: toType, ToID: toID,
		SourceType: meta.SourceType, Confidence: meta.Confidence,
	})
	if err != nil {
		return err
	}
	if !inserted {
		return nil
	}
	return s.appendLinked(fromType, fromID, rel, toType, toID, meta)
}

func (s *Service) validateLinkEndpoints(fromType, fromID, toType, toID string) error {
	if err := s.validateEntityExists(fromType, fromID); err != nil {
		return err
	}
	return s.validateEntityExists(toType, toID)
}

func (s *Service) validateEntityExists(entityType, id string) error {
	switch entityType {
	case EntityDecision:
		_, err := s.store.GetDecision(id)
		return err
	case EntityDiscovery:
		_, err := s.store.GetDiscovery(id)
		return err
	case EntityClaim:
		_, err := s.store.GetClaim(id)
		return err
	case EntityRegression:
		_, err := s.store.GetRegression(id)
		return err
	case EntityTask:
		_, err := s.store.GetTask(id)
		return err
	case EntityPlanChange:
		_, err := s.store.GetPlanChange(id)
		return err
	case EntityEvidence:
		_, err := s.store.GetEvidence(id)
		return err
	case EntityChange:
		_, err := s.store.GetChange(id)
		return err
	case EntityGoal:
		_, err := s.store.GetGoal(id)
		return err
	case EntityAssumption:
		_, err := s.store.GetAssumption(id)
		return err
	case EntityReview:
		_, err := s.store.GetReview(id)
		return err
	case EntityScope:
		_, err := s.store.GetScope(id)
		return err
	default:
		return &ErrValidation{Msg: "unsupported link entity type " + entityType}
	}
}

// ImportSeedScope upserts a thin graph scope; entity.created only on first insert.
func (s *Service) ImportSeedScope(ctx context.Context, in SeedScope) (store.Scope, bool, error) {
	_ = ctx
	slug := strings.TrimSpace(in.Slug)
	title := strings.TrimSpace(in.Title)
	kind := strings.TrimSpace(in.Kind)
	if slug == "" || title == "" {
		return store.Scope{}, false, &ErrValidation{Msg: "scope slug and title required"}
	}
	id := strings.TrimSpace(in.ID)
	if id == "" {
		id = uuid.NewString()
	}
	existed, err := scopeExists(s.store, id)
	if err != nil {
		return store.Scope{}, false, err
	}
	sc, err := s.store.UpsertScope(store.Scope{
		ID: id, Slug: slug, Title: title, Kind: kind,
	})
	if err != nil {
		return store.Scope{}, false, err
	}
	if err := s.upsertEntityCreated(EntityScope, sc.ID, sc.Title, existed); err != nil {
		return store.Scope{}, false, err
	}
	return sc, !existed, nil
}

func (s *Service) importSeedScopeMemberLink(ctx context.Context, fromID, scopeID string, meta LinkMeta) error {
	_ = ctx
	fromType, err := s.resolveLinkableEntityType(fromID)
	if err != nil {
		return err
	}
	return s.importSeedEntityLink(ctx, fromType, fromID, RelScopeMember, EntityScope, scopeID, meta)
}

func (s *Service) importSeedImplementsLink(ctx context.Context, fromID, toTaskID string, meta LinkMeta) error {
	fromType, err := s.resolveTaskOrDecision(fromID)
	if err != nil {
		return err
	}
	return s.importSeedEntityLink(ctx, fromType, fromID, RelImplements, EntityTask, toTaskID, meta)
}

// ImportSeedFinding upserts an impact finding by id.
func (s *Service) ImportSeedFinding(ctx context.Context, decisionID string, in SeedFinding) (store.DecisionImpactFinding, error) {
	_ = ctx
	class, err := NormalizeImpactClass(in.ImpactClass)
	if err != nil {
		return store.DecisionImpactFinding{}, err
	}
	kind, err := NormalizeFindingKind(in.Kind)
	if err != nil {
		return store.DecisionImpactFinding{}, err
	}
	unc, err := NormalizeUncertainty(in.Uncertainty)
	if err != nil {
		return store.DecisionImpactFinding{}, err
	}
	if _, err := s.store.GetDecision(decisionID); err != nil {
		return store.DecisionImpactFinding{}, err
	}
	return s.store.UpsertDecisionImpactFinding(store.DecisionImpactFinding{
		ID: in.ID, DecisionID: decisionID,
		ImpactClass: class, Uncertainty: unc, Kind: kind, Body: in.Body,
		RelatedType: strings.TrimSpace(in.RelatedType), RelatedID: strings.TrimSpace(in.RelatedID),
	})
}

// ImportSeedAlternative upserts a decision alternative by id.
func (s *Service) ImportSeedAlternative(ctx context.Context, decisionID string, in SeedAlternative) (store.DecisionAlternative, error) {
	_ = ctx
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return store.DecisionAlternative{}, &ErrValidation{Msg: "alternative title is required"}
	}
	if _, err := s.store.GetDecision(decisionID); err != nil {
		return store.DecisionAlternative{}, err
	}
	return s.store.UpsertDecisionAlternative(store.DecisionAlternative{
		ID: in.ID, DecisionID: decisionID,
		Title: title, Body: in.Body, IsRecommended: in.Recommended,
	})
}

// ImportSeedPlanPhase upserts a plan phase by id.
func (s *Service) ImportSeedPlanPhase(ctx context.Context, in SeedPlanPhase) (store.PlanPhase, error) {
	_ = ctx
	return s.store.UpsertPlanPhase(store.PlanPhase{
		ID: in.ID, GoalID: in.GoalID, Title: in.Title, Body: in.Body,
		Ord: in.Ord, Status: in.Status,
	})
}

// ImportSeedPlanScope upserts a plan scope by id.
func (s *Service) ImportSeedPlanScope(ctx context.Context, in SeedPlanScope) (store.PlanScope, error) {
	_ = ctx
	return s.store.UpsertPlanScope(store.PlanScope{
		ID: in.ID, PhaseID: in.PhaseID, Title: in.Title, Body: in.Body,
		Ord: in.Ord, Status: in.Status, AutoReplanCount: in.AutoReplanCount,
	})
}

// ImportSeedScopeDeepPlan upserts a scope deep plan by id.
func (s *Service) ImportSeedScopeDeepPlan(ctx context.Context, in SeedScopeDeepPlan) (store.ScopeDeepPlan, error) {
	_ = ctx
	return s.store.UpsertScopeDeepPlan(store.ScopeDeepPlan{
		ID: in.ID, ScopeID: in.ScopeID, ContentJSON: in.ContentJSON, Status: in.Status,
	})
}

// ImportSeedGoalPlanState upserts goal plan state (last-wins on goal_id).
func (s *Service) ImportSeedGoalPlanState(ctx context.Context, in SeedGoalPlanState) (store.GoalPlanState, error) {
	_ = ctx
	return s.store.UpsertGoalPlanState(store.GoalPlanState{
		GoalID: in.GoalID, CurrentScopeID: in.CurrentScopeID,
	})
}

// ImportSeedTransition applies a transition, skipping when task already at target work_state.
func (s *Service) ImportSeedTransition(ctx context.Context, tr SeedTransition) error {
	task, err := s.store.GetTask(tr.TaskID)
	if err != nil {
		return err
	}
	if task.WorkState == tr.To {
		return nil
	}
	actor := tr.Actor
	if actor == "" {
		actor = "seed"
	}
	return s.TransitionTask(ctx, tr.TaskID, tr.To, TransitionOptions{
		Actor:                    actor,
		Reason:                   tr.Reason,
		AllowDoneWithoutReview:   tr.AllowDone,
		AllowOperatorDone:        tr.AsOperator,
		AllowMissingCapabilities: tr.AllowMissingCaps,
	})
}

// ImportSeedBaseline upserts a baseline by id.
func (s *Service) ImportSeedBaseline(ctx context.Context, in SeedBaseline) (store.Baseline, error) {
	_ = ctx
	status := in.Status
	if status == "" {
		status = store.BaselineStatusActive
	}
	return s.store.UpsertBaseline(store.Baseline{
		ID: in.ID, GitCommit: in.GitCommit, ScoresJSON: in.ScoresJSON, Label: in.Label,
		SourceType: in.SourceType, Status: status, SupersedesID: in.SupersedesID,
		CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	})
}

// ImportSeedOutcomeResult upserts an outcome result by id.
func (s *Service) ImportSeedOutcomeResult(ctx context.Context, in SeedOutcomeResult) (store.OutcomeResult, error) {
	_ = ctx
	return s.store.UpsertOutcomeResult(store.OutcomeResult{
		ID: in.ID, TaskID: in.TaskID, Kind: in.Kind, TestName: in.TestName, TestStatus: in.TestStatus,
		GoalID: in.GoalID, VerificationStatus: in.VerificationStatus, BaselineID: in.BaselineID,
		ScoresJSON: in.ScoresJSON, ComparisonJSON: in.ComparisonJSON, Summary: in.Summary,
		Actor: in.Actor, SourceType: in.SourceType, Confidence: in.Confidence,
		CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	})
}

// ImportSeedChange upserts a change and nested paths by id.
func (s *Service) ImportSeedChange(ctx context.Context, in SeedChange) (store.Change, error) {
	_ = ctx
	c, err := s.store.UpsertChange(store.Change{
		ID: in.ID, TaskID: in.TaskID, GitCommit: in.GitCommit, ParentChangeID: in.ParentChangeID,
		Actor: in.Actor, Reason: in.Reason, Status: in.Status, SourceType: in.SourceType,
		Confidence: in.Confidence, CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.Change{}, err
	}
	for _, p := range in.Paths {
		if _, err := s.store.UpsertChangePath(store.ChangePath{
			ChangeID: c.ID, Path: p.Path, Status: p.Status, SymbolID: p.SymbolID,
		}); err != nil {
			return store.Change{}, err
		}
	}
	return c, nil
}

// ImportSeedEffect upserts an effect by id.
func (s *Service) ImportSeedEffect(ctx context.Context, in SeedEffect) (store.Effect, error) {
	_ = ctx
	return s.store.UpsertEffect(store.Effect{
		ID: in.ID, ChangeID: in.ChangeID, Dimension: in.Dimension, Expected: in.Expected,
		Actual: in.Actual, Comparison: in.Comparison, Confidence: in.Confidence,
		SourceType: in.SourceType, CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	})
}

// ImportSeedUncertainty upserts an uncertainty by id.
func (s *Service) ImportSeedUncertainty(ctx context.Context, in SeedUncertainty) (store.Uncertainty, error) {
	_ = ctx
	return s.store.UpsertUncertainty(store.Uncertainty{
		ID: in.ID, Title: in.Title, Body: in.Body, Severity: in.Severity, Status: in.Status,
		Kind: in.Kind, Confidence: in.Confidence, SourceType: in.SourceType, Resolution: in.Resolution,
		CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt, LastVerifiedAt: in.LastVerifiedAt,
	})
}

// ImportSeedHypothesis upserts a hypothesis by id.
func (s *Service) ImportSeedHypothesis(ctx context.Context, in SeedHypothesis) (store.Hypothesis, error) {
	_ = ctx
	return s.store.UpsertHypothesis(store.Hypothesis{
		ID: in.ID, Title: in.Title, Body: in.Body, Status: in.Status, Confidence: in.Confidence,
		SourceType: in.SourceType, CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
		LastVerifiedAt: in.LastVerifiedAt,
	})
}

// ImportSeedDecisionReconsideration upserts a decision reconsideration by id.
func (s *Service) ImportSeedDecisionReconsideration(ctx context.Context, in SeedDecisionReconsideration) (store.DecisionReconsideration, error) {
	_ = ctx
	return s.store.UpsertDecisionReconsideration(store.DecisionReconsideration{
		ID: in.ID, DecisionID: in.DecisionID, Trigger: in.Trigger, Status: in.Status,
		Reason: in.Reason, RelatedType: in.RelatedType, RelatedID: in.RelatedID,
		ReconsiderAt: in.ReconsiderAt, CreatedAt: in.CreatedAt,
	})
}

// ImportSeedRegression upserts a regression by id.
func (s *Service) ImportSeedRegression(ctx context.Context, in SeedRegression) (store.Regression, error) {
	_ = ctx
	return s.store.UpsertRegression(store.Regression{
		ID: in.ID, TaskID: in.TaskID, SourceKind: in.SourceKind, SourceID: in.SourceID,
		Dimension: in.Dimension, Attribution: in.Attribution, Status: in.Status, Summary: in.Summary,
		Actor: in.Actor, SourceType: in.SourceType, Confidence: in.Confidence,
		CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	})
}

// ImportSeedImprovement upserts an improvement by id.
func (s *Service) ImportSeedImprovement(ctx context.Context, in SeedImprovement) (store.Improvement, error) {
	_ = ctx
	evJSON := in.EvidenceIDsJSON
	if evJSON == "" {
		evJSON = "[]"
	}
	return s.store.UpsertImprovement(store.Improvement{
		ID: in.ID, ChangeID: in.ChangeID, TaskID: in.TaskID, Dimension: in.Dimension,
		Summary: in.Summary, EvidenceIDsJSON: evJSON, SourceType: in.SourceType,
		Confidence: in.Confidence, CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	})
}

// ImportSeedReflection upserts a reflection by id.
func (s *Service) ImportSeedReflection(ctx context.Context, in SeedReflection) (store.Reflection, error) {
	_ = ctx
	return s.store.UpsertReflection(store.Reflection{
		ID: in.ID, TaskID: in.TaskID, Summary: in.Summary,
		InvalidatedAssumptionsJSON: in.InvalidatedAssumptionsJSON,
		NewDependenciesJSON:        in.NewDependenciesJSON, UsefulTestsJSON: in.UsefulTestsJSON,
		BroadenTestsNote: in.BroadenTestsNote, Actor: in.Actor, SourceType: in.SourceType,
		Confidence: in.Confidence, CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	})
}

// ImportSeedChangePattern upserts one change_patterns row.
func (s *Service) ImportSeedChangePattern(ctx context.Context, in SeedChangePattern) (store.ChangePattern, error) {
	_ = ctx
	if err := s.store.UpsertChangePattern(store.ChangePattern{
		ChangeKind: in.ChangeKind, OutcomeKind: in.OutcomeKind,
		CountPositive: in.CountPositive, CountNegative: in.CountNegative, LastSeen: in.LastSeen,
	}); err != nil {
		return store.ChangePattern{}, err
	}
	rows, err := s.store.ListChangePatternsByKind(in.ChangeKind, 64)
	if err != nil {
		return store.ChangePattern{}, err
	}
	for _, r := range rows {
		if r.OutcomeKind == in.OutcomeKind {
			return r, nil
		}
	}
	return store.ChangePattern{
		ChangeKind: in.ChangeKind, OutcomeKind: in.OutcomeKind,
		CountPositive: in.CountPositive, CountNegative: in.CountNegative, LastSeen: in.LastSeen,
	}, nil
}

// ImportSeedHarnessAgent upserts a harness agent and its requirements by id.
func (s *Service) ImportSeedHarnessAgent(ctx context.Context, in SeedHarnessAgent) (store.HarnessAgent, error) {
	_ = ctx
	phases := in.DeliberationPhases
	if phases == "" {
		phases = "[]"
	}
	keywords := in.TaskKeywords
	if keywords == "" {
		keywords = "[]"
	}
	reqs := make([]string, 0, len(in.Requirements))
	for _, r := range in.Requirements {
		if slug := strings.TrimSpace(r.RequiredCapabilitySlug); slug != "" {
			reqs = append(reqs, slug)
		}
	}
	return s.UpsertHarnessAgent(ctx, HarnessAgentInput{
		ID: in.ID, Slug: in.Slug, Title: in.Title, Description: in.Description,
		SubagentType: in.SubagentType, DeliberationPhases: phases, TaskKeywords: keywords,
		RecommendSubagent: in.RecommendSubagent, RegistrySource: in.RegistrySource,
		RegistryVersion: in.RegistryVersion, ExternalURL: in.ExternalURL,
		Requirements: reqs,
	})
}

// ImportSeedEngineeringKnowledge upserts an engineering_knowledge row by id.
func (s *Service) ImportSeedEngineeringKnowledge(ctx context.Context, in SeedEngineeringKnowledge) (store.EngineeringKnowledge, error) {
	_ = ctx
	evJSON := in.EvidenceIDsJSON
	if evJSON == "" {
		evJSON = "[]"
	}
	bodyJSON := in.BodyJSON
	if bodyJSON == "" {
		bodyJSON = "{}"
	}
	status := in.Status
	if status == "" {
		status = store.KnowledgeStatusActive
	}
	return s.store.UpsertEngineeringKnowledge(store.EngineeringKnowledge{
		ID: in.ID, Title: in.Title, BodyJSON: bodyJSON, Topic: in.Topic,
		EvidenceIDsJSON: evJSON, Confidence: in.Confidence, Status: status,
		SourceType: in.SourceType, CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	})
}

// ImportSeedDeliberationState upserts deliberation state keyed by task_id.
func (s *Service) ImportSeedDeliberationState(ctx context.Context, in SeedDeliberationState) (store.DeliberationState, error) {
	_ = ctx
	return s.store.UpsertDeliberationState(store.DeliberationState{
		TaskID: in.TaskID, GoalID: in.GoalID, CurrentPhase: in.CurrentPhase,
		HopCount: in.HopCount, LastPhase: in.LastPhase, PlanCritiqued: in.PlanCritiqued,
		Stopped: in.Stopped, StopReason: in.StopReason,
		ConsecutiveEmptyApplies: in.ConsecutiveEmptyApplies, UpdatedAt: in.UpdatedAt,
	})
}
