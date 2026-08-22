package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

const (
	// VCSCaptureTaskID is the sentinel task_id for VCS-promoted changes (no FK).
	VCSCaptureTaskID = "trace:vcs-capture"
	// SourceTypeVCS marks changes promoted from vcs_commits.
	SourceTypeVCS = "VCS"
)

// PromoteVCSCommitOptions controls path filtering when promoting indexed commits.
type PromoteVCSCommitOptions struct {
	// AllPaths skips the meaningful-path filter (CLI --all).
	AllPaths bool
}

const (
	maxChangeTextBytes      = 8192
	maxEffectDimensionBytes = 64
)

var gitCommitOID = regexp.MustCompile(`^[0-9a-fA-F]{7,64}$`)

// ChangePathInput is one repo-relative path on a change (no file contents).
type ChangePathInput struct {
	Path     string
	Status   string
	SymbolID string
}

// ExpectedEffectInput declares an expected effect. Comparison must stay empty.
type ExpectedEffectInput struct {
	Dimension  string
	Expected   string
	Confidence float64
	SourceType string
}

// ChangeInput creates a Change.
type ChangeInput struct {
	TaskID         string
	GitCommit      string
	ParentChangeID string
	Actor          string
	Reason         string
	Paths          []ChangePathInput
	Expected       []ExpectedEffectInput
	DecisionID     string
	SourceType     string
	Confidence     float64
}

// RecordActualEffectInput records actual + caller-supplied comparison.
type RecordActualEffectInput struct {
	Dimension        string
	Actual           string
	Comparison       string
	Confidence       float64
	EvidenceIDs      []string
	EmitDiscovery    bool
	DiscoveryTitle   string
	DiscoveryBody    string
	HypothesisID     string
	CreateHypothesis bool
	HypothesisTitle  string
	HypothesisBody   string
}

func normalizeGitCommit(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", nil
	}
	if !gitCommitOID.MatchString(s) {
		return "", &ErrValidation{Msg: "git_commit must be a hex OID (7-64 chars)"}
	}
	return strings.ToLower(s), nil
}

func failClosedMaxBytes(label, s string, max int) error {
	if len(s) > max {
		return &ErrValidation{Msg: label + " exceeds " + strconv.Itoa(max) + " bytes"}
	}
	return nil
}

func normalizeEffectComparison(raw string) (string, error) {
	c := strings.TrimSpace(raw)
	switch c {
	case store.EffectComparisonSupported, store.EffectComparisonPartiallySupported, store.EffectComparisonContradicted:
		return c, nil
	default:
		return "", &ErrValidation{Msg: "comparison must be supported, partially_supported, or contradicted"}
	}
}

func normalizeDimension(raw string) (string, error) {
	d := strings.TrimSpace(raw)
	if d == "" {
		return "", &ErrValidation{Msg: "dimension is required"}
	}
	if err := failClosedMaxBytes("dimension", d, maxEffectDimensionBytes); err != nil {
		return "", err
	}
	return d, nil
}

func (s *Service) maybeMarkChangeCompared(c store.Change) (store.Change, error) {
	if c.Status != store.ChangeStatusRecorded {
		return c, nil
	}
	effects, err := s.store.ListEffectsByChangeID(c.ID)
	if err != nil {
		return store.Change{}, err
	}
	if len(effects) == 0 {
		return c, nil
	}
	for _, e := range effects {
		if e.Comparison == store.EffectComparisonNone {
			return c, nil
		}
	}
	c.Status = store.ChangeStatusCompared
	return s.store.UpsertChange(c)
}

func (s *Service) maybeDemoteCompared(c store.Change) (store.Change, error) {
	if c.Status != store.ChangeStatusCompared {
		return c, nil
	}
	if strings.TrimSpace(c.GitCommit) == "" {
		c.Status = store.ChangeStatusOpen
	} else {
		c.Status = store.ChangeStatusRecorded
	}
	return s.store.UpsertChange(c)
}

// CreateChange persists a change with ≥1 path. Empty SHA → OPEN; SHA present → RECORDED.
func (s *Service) CreateChange(ctx context.Context, in ChangeInput) (store.Change, error) {
	_ = ctx
	taskID := strings.TrimSpace(in.TaskID)
	if taskID == "" {
		return store.Change{}, &ErrValidation{Msg: "task_id is required"}
	}
	if err := failClosedMaxBytes("reason", in.Reason, maxChangeTextBytes); err != nil {
		return store.Change{}, err
	}
	gitCommit, err := normalizeGitCommit(in.GitCommit)
	if err != nil {
		return store.Change{}, err
	}
	if len(in.Paths) < 1 {
		return store.Change{}, &ErrValidation{Msg: "at least one path is required"}
	}
	seenPaths := make(map[string]struct{}, len(in.Paths))
	var paths []store.ChangePath
	for _, p := range in.Paths {
		path := store.NormalizePath(strings.TrimSpace(p.Path))
		if path == "" {
			return store.Change{}, &ErrValidation{Msg: "path is required"}
		}
		if _, dup := seenPaths[path]; dup {
			return store.Change{}, &ErrValidation{Msg: "duplicate path"}
		}
		seenPaths[path] = struct{}{}
		paths = append(paths, store.ChangePath{
			Path:     path,
			Status:   strings.TrimSpace(p.Status),
			SymbolID: strings.TrimSpace(p.SymbolID),
		})
	}
	parentID := strings.TrimSpace(in.ParentChangeID)
	if parentID != "" {
		if _, err := s.store.GetChange(parentID); err != nil {
			return store.Change{}, err
		}
	}
	decisionID := strings.TrimSpace(in.DecisionID)
	if decisionID != "" {
		if _, err := s.store.GetDecision(decisionID); err != nil {
			return store.Change{}, err
		}
	}
	for _, exp := range in.Expected {
		if _, err := normalizeDimension(exp.Dimension); err != nil {
			return store.Change{}, err
		}
		if strings.TrimSpace(exp.Expected) == "" {
			return store.Change{}, &ErrValidation{Msg: "expected is required"}
		}
		if err := failClosedMaxBytes("expected", exp.Expected, maxChangeTextBytes); err != nil {
			return store.Change{}, err
		}
	}

	src := in.SourceType
	if src == "" {
		src = DefaultSourceType
	}
	status := store.ChangeStatusOpen
	if gitCommit != "" {
		status = store.ChangeStatusRecorded
	}
	id := uuid.NewString()
	c, err := s.store.UpsertChange(store.Change{
		ID:             id,
		TaskID:         taskID,
		GitCommit:      gitCommit,
		ParentChangeID: parentID,
		Actor:          strings.TrimSpace(in.Actor),
		Reason:         in.Reason,
		Status:         status,
		SourceType:     src,
		Confidence:     in.Confidence,
	})
	if err != nil {
		return store.Change{}, err
	}
	for _, p := range paths {
		p.ChangeID = c.ID
		if _, err := s.store.InsertChangePath(p); err != nil {
			return store.Change{}, err
		}
	}
	for _, exp := range in.Expected {
		if _, err := s.RecordExpectedEffect(ctx, c.ID, exp); err != nil {
			return store.Change{}, err
		}
	}
	if decisionID != "" {
		if err := s.insertTypedLink(EntityChange, c.ID, RelChangeImplementsDecision, EntityDecision, decisionID); err != nil {
			return store.Change{}, err
		}
	}
	if err := s.appendCreated(EntityChange, c.ID, c.Reason); err != nil {
		return store.Change{}, err
	}
	return s.store.GetChange(c.ID)
}

// RecordChangeCommit sets git_commit on OPEN (or empty SHA). Fail closed if replacing a different SHA.
func (s *Service) RecordChangeCommit(ctx context.Context, changeID, gitCommit string) (store.Change, error) {
	_ = ctx
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return store.Change{}, &ErrValidation{Msg: "change id is required"}
	}
	sha, err := normalizeGitCommit(gitCommit)
	if err != nil {
		return store.Change{}, err
	}
	if sha == "" {
		return store.Change{}, &ErrValidation{Msg: "git_commit is required"}
	}
	c, err := s.store.GetChange(changeID)
	if err != nil {
		return store.Change{}, err
	}
	if c.Status == store.ChangeStatusSuperseded {
		return store.Change{}, &ErrValidation{Msg: "SUPERSEDED change is terminal"}
	}
	firstSet := c.GitCommit == ""
	if c.GitCommit != "" && c.GitCommit != sha {
		return store.Change{}, &ErrValidation{Msg: "git_commit already set; create a child change"}
	}
	c.GitCommit = sha
	if c.Status == store.ChangeStatusOpen {
		c.Status = store.ChangeStatusRecorded
	}
	out, err := s.store.UpsertChange(c)
	if err != nil {
		return store.Change{}, err
	}
	if firstSet {
		if err := s.appendNamed(EventChangeRecorded, EntityChange, out.ID, map[string]string{
			"git_commit": out.GitCommit,
		}); err != nil {
			return store.Change{}, err
		}
	}
	return s.maybeMarkChangeCompared(out)
}

// SupersedeChange marks a change SUPERSEDED. Reason is required. Terminal afterwards.
func (s *Service) SupersedeChange(ctx context.Context, changeID, reason string) (store.Change, error) {
	_ = ctx
	changeID = strings.TrimSpace(changeID)
	reason = strings.TrimSpace(reason)
	if changeID == "" {
		return store.Change{}, &ErrValidation{Msg: "change id is required"}
	}
	if reason == "" {
		return store.Change{}, &ErrValidation{Msg: "reason is required"}
	}
	if err := failClosedMaxBytes("reason", reason, maxChangeTextBytes); err != nil {
		return store.Change{}, err
	}
	c, err := s.store.GetChange(changeID)
	if err != nil {
		return store.Change{}, err
	}
	switch c.Status {
	case store.ChangeStatusOpen, store.ChangeStatusRecorded, store.ChangeStatusCompared:
	case store.ChangeStatusSuperseded:
		return store.Change{}, &ErrValidation{Msg: "SUPERSEDED change is terminal"}
	default:
		return store.Change{}, &ErrValidation{Msg: "change cannot be superseded"}
	}
	c.Status = store.ChangeStatusSuperseded
	c.Reason = reason
	return s.store.UpsertChange(c)
}

// RecordExpectedEffect inserts an expected row. Comparison stays empty.
func (s *Service) RecordExpectedEffect(ctx context.Context, changeID string, in ExpectedEffectInput) (store.Effect, error) {
	_ = ctx
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return store.Effect{}, &ErrValidation{Msg: "change id is required"}
	}
	dim, err := normalizeDimension(in.Dimension)
	if err != nil {
		return store.Effect{}, err
	}
	expected := in.Expected
	if strings.TrimSpace(expected) == "" {
		return store.Effect{}, &ErrValidation{Msg: "expected is required"}
	}
	if err := failClosedMaxBytes("expected", expected, maxChangeTextBytes); err != nil {
		return store.Effect{}, err
	}
	c, err := s.store.GetChange(changeID)
	if err != nil {
		return store.Effect{}, err
	}
	if c.Status == store.ChangeStatusSuperseded {
		return store.Effect{}, &ErrValidation{Msg: "SUPERSEDED change is terminal"}
	}
	if _, err := s.store.GetEffectByChangeDimension(changeID, dim); err == nil {
		return store.Effect{}, &ErrValidation{Msg: "expected effect already exists for dimension"}
	} else if !errors.Is(err, sql.ErrNoRows) {
		return store.Effect{}, err
	}
	src := in.SourceType
	if src == "" {
		src = DefaultSourceType
	}
	e, err := s.store.UpsertEffect(store.Effect{
		ChangeID:   changeID,
		Dimension:  dim,
		Expected:   expected,
		Comparison: store.EffectComparisonNone,
		Confidence: in.Confidence,
		SourceType: src,
	})
	if err != nil {
		return store.Effect{}, err
	}
	if _, err := s.maybeDemoteCompared(c); err != nil {
		return store.Effect{}, err
	}
	return e, nil
}

// RecordActualEffect writes actual + comparison onto an existing expected dimension.
func (s *Service) RecordActualEffect(ctx context.Context, changeID string, in RecordActualEffectInput) (store.Effect, *store.Discovery, error) {
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return store.Effect{}, nil, &ErrValidation{Msg: "change id is required"}
	}
	dim, err := normalizeDimension(in.Dimension)
	if err != nil {
		return store.Effect{}, nil, err
	}
	if strings.TrimSpace(in.Actual) == "" {
		return store.Effect{}, nil, &ErrValidation{Msg: "actual is required"}
	}
	if err := failClosedMaxBytes("actual", in.Actual, maxChangeTextBytes); err != nil {
		return store.Effect{}, nil, err
	}
	cmp, err := normalizeEffectComparison(in.Comparison)
	if err != nil {
		return store.Effect{}, nil, err
	}
	hypID := strings.TrimSpace(in.HypothesisID)
	if hypID != "" && in.CreateHypothesis {
		return store.Effect{}, nil, &ErrValidation{Msg: "HypothesisID and CreateHypothesis are mutually exclusive"}
	}
	if cmp != store.EffectComparisonContradicted {
		if in.EmitDiscovery || in.CreateHypothesis || hypID != "" {
			return store.Effect{}, nil, &ErrValidation{Msg: "Discovery/Hypothesis flags are only valid when comparison is contradicted"}
		}
	}
	if in.EmitDiscovery && strings.TrimSpace(in.DiscoveryTitle) == "" {
		return store.Effect{}, nil, &ErrValidation{Msg: "DiscoveryTitle is required when EmitDiscovery"}
	}
	if in.CreateHypothesis && strings.TrimSpace(in.HypothesisTitle) == "" {
		return store.Effect{}, nil, &ErrValidation{Msg: "HypothesisTitle is required when CreateHypothesis"}
	}
	c, err := s.store.GetChange(changeID)
	if err != nil {
		return store.Effect{}, nil, err
	}
	if c.Status == store.ChangeStatusSuperseded {
		return store.Effect{}, nil, &ErrValidation{Msg: "SUPERSEDED change is terminal"}
	}
	e, err := s.store.GetEffectByChangeDimension(changeID, dim)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return store.Effect{}, nil, &ErrValidation{Msg: "expected effect for dimension is required"}
		}
		return store.Effect{}, nil, err
	}
	for _, eid := range in.EvidenceIDs {
		eid = strings.TrimSpace(eid)
		if eid == "" {
			return store.Effect{}, nil, &ErrValidation{Msg: "evidence id is required"}
		}
		if _, err := s.store.GetEvidence(eid); err != nil {
			return store.Effect{}, nil, err
		}
	}
	if hypID != "" {
		if _, err := s.store.GetHypothesis(hypID); err != nil {
			return store.Effect{}, nil, err
		}
	}

	e.Actual = in.Actual
	e.Comparison = cmp
	e.Confidence = in.Confidence
	out, err := s.store.UpsertEffect(e)
	if err != nil {
		return store.Effect{}, nil, err
	}
	if err := s.appendNamed(EventEffectCompared, EntityEffect, out.ID, map[string]string{
		"dimension":  out.Dimension,
		"comparison": out.Comparison,
		"change_id":  out.ChangeID,
	}); err != nil {
		return store.Effect{}, nil, err
	}

	if _, err := s.maybeMarkChangeCompared(c); err != nil {
		return store.Effect{}, nil, err
	}

	for _, eid := range in.EvidenceIDs {
		eid = strings.TrimSpace(eid)
		if err := s.insertTypedLink(EntityEffect, out.ID, RelEffectSupportedBy, EntityEvidence, eid); err != nil {
			return store.Effect{}, nil, err
		}
	}

	if cmp != store.EffectComparisonContradicted {
		return out, nil, nil
	}

	if err := s.appendNamed(EventEffectContradicted, EntityChange, c.ID, map[string]string{
		"effect_id":  out.ID,
		"dimension":  out.Dimension,
		"comparison": out.Comparison,
	}); err != nil {
		return store.Effect{}, nil, err
	}

	links, err := s.store.ListLinksFrom(EntityChange, c.ID)
	if err != nil {
		return store.Effect{}, nil, err
	}
	for _, l := range links {
		if l.Rel != RelChangeImplementsDecision || l.ToType != EntityDecision {
			continue
		}
		if _, err := s.RecordDecisionReconsideration(ctx, l.ToID, ReconsiderationInput{
			Trigger:     store.ReconsiderTriggerContradictedEffect,
			Status:      store.ReconsiderStatusFired,
			Reason:      "effect contradicted: " + out.Dimension,
			RelatedType: EntityEffect,
			RelatedID:   out.ID,
		}); err != nil {
			return store.Effect{}, nil, err
		}
	}

	if in.CreateHypothesis {
		h, err := s.CreateHypothesis(ctx, HypothesisInput{
			Title: strings.TrimSpace(in.HypothesisTitle),
			Body:  in.HypothesisBody,
		})
		if err != nil {
			return store.Effect{}, nil, err
		}
		hypID = h.ID
	}
	if hypID != "" {
		if err := s.insertTypedLink(EntityHypothesis, hypID, RelHypothesisExplainsEffect, EntityEffect, out.ID); err != nil {
			return store.Effect{}, nil, err
		}
	}

	var disc *store.Discovery
	if in.EmitDiscovery {
		d, err := s.CreateDiscovery(ctx, DiscoveryInput{
			Title:    strings.TrimSpace(in.DiscoveryTitle),
			Body:     in.DiscoveryBody,
			Severity: SeverityPlanAffecting,
		})
		if err != nil {
			return store.Effect{}, nil, err
		}
		if err := s.insertTypedLink(EntityDiscovery, d.ID, RelDiscoveryFromContradictedEffect, EntityEffect, out.ID); err != nil {
			return store.Effect{}, nil, err
		}
		disc = &d
	}
	return out, disc, nil
}

// PromoteVCSCommitToChange creates a RECORDED change from vcs_commits + paths.
// Idempotent on OID. Zero qualifying paths → empty Change, nil error.
func (s *Service) PromoteVCSCommitToChange(ctx context.Context, oid string) (store.Change, error) {
	return s.PromoteVCSCommitToChangeOpts(ctx, oid, PromoteVCSCommitOptions{})
}

// PromoteVCSCommitToChangeOpts promotes with optional --all semantics.
func (s *Service) PromoteVCSCommitToChangeOpts(ctx context.Context, oid string, opts PromoteVCSCommitOptions) (store.Change, error) {
	_ = ctx
	sha, err := normalizeGitCommit(oid)
	if err != nil {
		return store.Change{}, err
	}
	if existing, err := s.store.GetChangeByGitCommit(sha); err == nil {
		return existing, nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return store.Change{}, err
	}

	commit, err := s.store.GetIndexedCommit(sha)
	if err != nil {
		return store.Change{}, fmt.Errorf("domain: promote vcs commit: %w", err)
	}
	rawPaths, err := s.store.ListIndexedCommitPaths(sha)
	if err != nil {
		return store.Change{}, err
	}

	var paths []ChangePathInput
	for _, p := range rawPaths {
		if !opts.AllPaths && !s.isMeaningfulChangePath(p.Path) {
			continue
		}
		paths = append(paths, ChangePathInput{
			Path:   p.Path,
			Status: p.Status,
		})
	}
	if len(paths) == 0 {
		return store.Change{}, nil
	}

	return s.CreateChange(ctx, ChangeInput{
		TaskID:     VCSCaptureTaskID,
		GitCommit:  sha,
		Reason:     commit.Subject,
		SourceType: SourceTypeVCS,
		Paths:      paths,
	})
}

// PromoteRecentVCSCommits promotes indexed commits after sinceOID (exclusive).
func (s *Service) PromoteRecentVCSCommits(ctx context.Context, sinceOID string, opts PromoteVCSCommitOptions) ([]store.Change, error) {
	commits, err := s.store.ListIndexedCommitsSince(sinceOID)
	if err != nil {
		return nil, err
	}
	out := make([]store.Change, 0, len(commits))
	for _, c := range commits {
		ch, err := s.PromoteVCSCommitToChangeOpts(ctx, c.OID, opts)
		if err != nil {
			return out, err
		}
		if ch.ID != "" {
			out = append(out, ch)
		}
	}
	return out, nil
}

func (s *Service) isMeaningfulChangePath(path string) bool {
	path = store.NormalizePath(strings.TrimSpace(path))
	if path == "" {
		return false
	}
	if _, ok := analyzers.DetectLanguage(path); ok {
		return true
	}
	if _, err := s.store.GetFileByPath(path); err == nil {
		return true
	}
	return false
}

// GetChange loads a change by id.
func (s *Service) GetChange(ctx context.Context, id string) (store.Change, error) {
	_ = ctx
	return s.store.GetChange(strings.TrimSpace(id))
}

// ListChangesByTaskID lists changes for a task.
func (s *Service) ListChangesByTaskID(ctx context.Context, taskID string) ([]store.Change, error) {
	_ = ctx
	taskID = strings.TrimSpace(taskID)
	if taskID == "" {
		return nil, &ErrValidation{Msg: "task_id is required"}
	}
	return s.store.ListChangesByTaskID(taskID)
}

// ListChanges returns recent changes newest-first. limit defaults to 32 (cap 64).
// Empty taskID lists across all tasks; unknown taskID returns an empty slice.
func (s *Service) ListChanges(ctx context.Context, limit int, taskID string) ([]store.Change, error) {
	_ = ctx
	return s.store.ListChangesRecent(limit, strings.TrimSpace(taskID))
}

// ListChangePaths lists path refs for a change.
func (s *Service) ListChangePaths(ctx context.Context, changeID string) ([]store.ChangePath, error) {
	_ = ctx
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return nil, &ErrValidation{Msg: "change id is required"}
	}
	return s.store.ListChangePaths(changeID)
}

// ListEffectsByChangeID lists effects for a change.
func (s *Service) ListEffectsByChangeID(ctx context.Context, changeID string) ([]store.Effect, error) {
	_ = ctx
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return nil, &ErrValidation{Msg: "change id is required"}
	}
	return s.store.ListEffectsByChangeID(changeID)
}

// ResolveChangePath returns Git bytes via vcs.Repository.ShowFile. Never persists them.
func (s *Service) ResolveChangePath(ctx context.Context, changeID, path string, repo vcs.Repository) ([]byte, error) {
	changeID = strings.TrimSpace(changeID)
	if changeID == "" {
		return nil, &ErrValidation{Msg: "change id is required"}
	}
	if repo == nil {
		return nil, &ErrValidation{Msg: "repository is required"}
	}
	c, err := s.store.GetChange(changeID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(c.GitCommit) == "" {
		return nil, &ErrValidation{Msg: "git_commit is required to resolve a path"}
	}
	norm := store.NormalizePath(strings.TrimSpace(path))
	if norm == "" {
		return nil, &ErrValidation{Msg: "path is required"}
	}
	if _, err := s.store.GetChangePath(changeID, norm); err != nil {
		return nil, &ErrValidation{Msg: "path is not on this change"}
	}
	return repo.ShowFile(ctx, c.GitCommit, norm)
}
