package loop

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

var allowedApplyWriteKeys = map[string]struct{}{
	"discoveries":   {},
	"plan_changes":  {},
	"spawned_tasks": {},
	"stop":          {},
	"uncertainties": {},
	"hypotheses":    {},
	"changes":       {},
	"effects":       {},
	"test_results":  {},
	"verifications": {},
	"evaluations":   {},
	"regressions":   {},
	"reflections":   {},
}

func validateWritesKeys(raw json.RawMessage) error {
	if len(raw) == 0 {
		return nil
	}
	var writes map[string]json.RawMessage
	if err := json.Unmarshal(raw, &writes); err != nil {
		return fmt.Errorf("loop apply: parse writes: %w", err)
	}
	for key := range writes {
		if _, ok := allowedApplyWriteKeys[key]; !ok {
			return fmt.Errorf("loop apply: unknown writes key %q", key)
		}
	}
	return nil
}

func applyCognitiveWrites(ctx context.Context, dom *domain.Service, st *store.Store, seed ApplySeed, w ApplyWrites) error {
	for _, u := range w.Uncertainties {
		if err := applyUncertainty(ctx, dom, seed, u); err != nil {
			return err
		}
	}
	for _, h := range w.Hypotheses {
		if err := applyHypothesis(ctx, dom, h); err != nil {
			return err
		}
	}
	for _, c := range w.Changes {
		if err := applyChange(ctx, dom, st, seed, c); err != nil {
			return err
		}
	}
	for _, e := range w.Effects {
		if err := applyEffect(ctx, dom, e); err != nil {
			return err
		}
	}
	for _, tr := range w.TestResults {
		if err := applyTestResult(ctx, dom, seed, tr); err != nil {
			return err
		}
	}
	for _, v := range w.Verifications {
		if err := applyVerification(ctx, dom, seed, v); err != nil {
			return err
		}
	}
	for _, ev := range w.Evaluations {
		if err := applyEvaluation(ctx, dom, seed, ev); err != nil {
			return err
		}
	}
	for _, r := range w.Regressions {
		if err := applyRegression(ctx, dom, seed, r); err != nil {
			return err
		}
	}
	for _, rf := range w.Reflections {
		if err := applyReflection(ctx, dom, seed, rf); err != nil {
			return err
		}
	}
	return nil
}

func applyUncertainty(ctx context.Context, dom *domain.Service, seed ApplySeed, u ApplyUncertainty) error {
	status := strings.TrimSpace(u.Status)
	if status == store.UncertaintyStatusResolved {
		resolution := strings.TrimSpace(u.Body)
		if resolution == "" {
			resolution = strings.TrimSpace(u.Title)
		}
		if resolution == "" {
			return fmt.Errorf("loop apply: uncertainty %q resolve requires body or title as resolution", u.ID)
		}
		_, err := dom.ResolveUncertainty(ctx, u.ID, resolution)
		return err
	}
	taskID := strings.TrimSpace(u.TaskID)
	if taskID == "" && strings.TrimSpace(u.Severity) == store.UncertaintySeverityBlocking {
		taskID = seed.TaskID
	}
	_, err := dom.CreateUncertainty(ctx, domain.UncertaintyInput{
		ID:       u.ID,
		Title:    u.Title,
		Body:     u.Body,
		Severity: u.Severity,
		Status:   u.Status,
		Kind:     u.Kind,
		TaskID:   taskID,
		GoalID:   seed.GoalID,
	})
	return err
}

func applyHypothesis(ctx context.Context, dom *domain.Service, h ApplyHypothesis) error {
	_, err := dom.CreateHypothesis(ctx, domain.HypothesisInput{
		ID:            h.ID,
		Title:         h.Title,
		Body:          h.Body,
		EvidenceIDs:   h.EvidenceIDs,
		UncertaintyID: h.UncertaintyID,
	})
	return err
}

func applyChange(ctx context.Context, dom *domain.Service, st *store.Store, seed ApplySeed, c ApplyChange) error {
	if _, err := st.GetChange(c.ID); err == nil {
		if c.GitCommit != "" {
			_, err := dom.RecordChangeCommit(ctx, c.ID, c.GitCommit)
			return err
		}
		return nil
	}
	if len(c.Paths) < 1 {
		return fmt.Errorf("loop apply: change %q requires at least one path", c.ID)
	}
	status := store.ChangeStatusOpen
	gitCommit := strings.TrimSpace(c.GitCommit)
	if gitCommit != "" {
		status = store.ChangeStatusRecorded
	}
	ch, err := st.UpsertChange(store.Change{
		ID:             c.ID,
		TaskID:         seed.TaskID,
		GitCommit:      strings.ToLower(gitCommit),
		ParentChangeID: strings.TrimSpace(c.ParentChangeID),
		Actor:          strings.TrimSpace(c.Actor),
		Reason:         c.Reason,
		Status:         status,
	})
	if err != nil {
		return err
	}
	for _, p := range c.Paths {
		path := store.NormalizePath(strings.TrimSpace(p))
		if path == "" {
			return fmt.Errorf("loop apply: change %q path is required", c.ID)
		}
		if _, err := st.InsertChangePath(store.ChangePath{ChangeID: ch.ID, Path: path}); err != nil {
			return err
		}
	}
	for _, e := range c.Expected {
		if _, err := dom.RecordExpectedEffect(ctx, ch.ID, domain.ExpectedEffectInput{
			Dimension: e.Dimension,
			Expected:  e.Expected,
		}); err != nil {
			return err
		}
	}
	return nil
}

func applyEffect(ctx context.Context, dom *domain.Service, e ApplyEffect) error {
	if strings.TrimSpace(e.Expected) != "" && strings.TrimSpace(e.Actual) == "" {
		_, err := dom.RecordExpectedEffect(ctx, e.ChangeID, domain.ExpectedEffectInput{
			Dimension: e.Dimension,
			Expected:  e.Expected,
		})
		return err
	}
	if strings.TrimSpace(e.Actual) != "" {
		_, _, err := dom.RecordActualEffect(ctx, e.ChangeID, domain.RecordActualEffectInput{
			Dimension:    e.Dimension,
			Actual:       e.Actual,
			Comparison:   e.Comparison,
			EvidenceIDs:  e.EvidenceIDs,
			HypothesisID: e.HypothesisID,
		})
		return err
	}
	return fmt.Errorf("loop apply: effect on change %q requires expected or actual", e.ChangeID)
}

func applyTestResult(ctx context.Context, dom *domain.Service, seed ApplySeed, tr ApplyTestResult) error {
	_, err := dom.RecordTestOutcome(ctx, domain.TestOutcomeInput{
		TaskID:      seed.TaskID,
		TestName:    tr.TestName,
		TestStatus:  tr.TestStatus,
		Summary:     tr.Summary,
		EvidenceIDs: tr.EvidenceIDs,
	})
	return err
}

func applyVerification(ctx context.Context, dom *domain.Service, seed ApplySeed, v ApplyVerification) error {
	goalID := strings.TrimSpace(v.GoalID)
	if goalID == "" {
		goalID = seed.GoalID
	}
	_, err := dom.RecordVerificationOutcome(ctx, domain.VerificationOutcomeInput{
		TaskID:             seed.TaskID,
		GoalID:             goalID,
		VerificationStatus: v.VerificationStatus,
		EvidenceIDs:        v.EvidenceIDs,
		Summary:            v.Summary,
	})
	return err
}

func applyEvaluation(ctx context.Context, dom *domain.Service, seed ApplySeed, ev ApplyEvaluation) error {
	_, err := dom.RecordEvaluationOutcome(ctx, domain.EvaluationOutcomeInput{
		TaskID:     seed.TaskID,
		BaselineID: ev.BaselineID,
		ScoresJSON: ev.ScoresJSON,
	})
	return err
}

func applyRegression(ctx context.Context, dom *domain.Service, seed ApplySeed, r ApplyRegression) error {
	switch strings.TrimSpace(r.SourceKind) {
	case store.RegressionSourceEvaluation:
		_, err := dom.RecordRegressionFromEvaluation(ctx, domain.EvaluationRegressionInput{
			OutcomeID: r.SourceID,
			TaskID:    seed.TaskID,
			Summary:   r.Summary,
		})
		return err
	case store.RegressionSourceContradictedEffect:
		_, err := dom.RecordRegressionFromContradictedEffect(ctx, domain.EffectRegressionInput{
			EffectID: r.SourceID,
			TaskID:   seed.TaskID,
			Summary:  r.Summary,
		})
		return err
	default:
		return fmt.Errorf("loop apply: regression source_kind must be evaluation or contradicted_effect")
	}
}

func applyReflection(ctx context.Context, dom *domain.Service, seed ApplySeed, rf ApplyReflection) error {
	deps := make([]domain.DependencyRef, 0, len(rf.NewDependencies))
	for _, d := range rf.NewDependencies {
		deps = append(deps, domain.DependencyRef{Kind: d.Kind, Ref: d.Ref})
	}
	_, err := dom.CreateReflection(ctx, domain.ReflectionInput{
		TaskID:                   seed.TaskID,
		Summary:                  rf.Summary,
		InvalidatedAssumptionIDs: rf.InvalidatedAssumptionIDs,
		NewDependencies:          deps,
		UsefulTests:              rf.UsefulTests,
		BroadenTestsNote:         rf.BroadenTestsNote,
	})
	return err
}
