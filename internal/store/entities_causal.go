package store

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Decision is an accepted (or proposed) choice with provenance.
type Decision struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

// Assumption is a working belief with provenance.
type Assumption struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

// Discovery severity vocabulary (mig 007). Fail closed on other values in domain.
const (
	SeverityINFO          = "INFO"
	SeverityPlanAffecting = "PLAN_AFFECTING"
	SeverityBlocking      = "BLOCKING"
)

// Discovery is a finding that may cause plan changes.
type Discovery struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	Severity       string // INFO | PLAN_AFFECTING | BLOCKING; default INFO
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

// PlanChange records an intentional plan delta.
type PlanChange struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

// Claim is an agent or human assertion (honesty path; not authority alone).
type Claim struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

// Evidence supports a claim (honesty path; not DONE authority alone).
type Evidence struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

// ReviewResult vocabulary for reviews.result (mig 005).
const (
	ReviewResultOpen      = ""
	ReviewResultPass      = "PASS"
	ReviewResultFail      = "FAIL"
	ReviewResultUncertain = "UNCERTAIN"
)

// Review judges a task (linked via review_judges_task). Result starts open ("").
type Review struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	Result         string // PASS | FAIL | UNCERTAIN | "" (open)
	CreatedAt      string
	UpdatedAt      string
	LastVerifiedAt *string
}

// UpsertDecision inserts or replaces a decision by id.
func (s *Store) UpsertDecision(d Decision) (Decision, error) {
	now := nowRFC3339()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	if d.Status == "" {
		d.Status = StatusActive
	}
	if d.CreatedAt == "" {
		d.CreatedAt = now
	}
	d.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO decisions(id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			body = excluded.body,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			status = excluded.status,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, d.ID, d.Title, d.Body, d.SourceType, d.Confidence, d.Status, d.CreatedAt, d.UpdatedAt, nullStr(d.LastVerifiedAt))
	if err != nil {
		return Decision{}, fmt.Errorf("store: upsert decision: %w", err)
	}
	out, err := s.GetDecision(d.ID)
	if err != nil {
		return Decision{}, err
	}
	if err := s.SyncEntityFTS("decision", out.ID); err != nil {
		return Decision{}, err
	}
	return out, nil
}

// GetDecision loads a decision by id.
func (s *Store) GetDecision(id string) (Decision, error) {
	var d Decision
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at
		FROM decisions WHERE id = ?
	`, id).Scan(
		&d.ID, &d.Title, &d.Body, &d.SourceType, &d.Confidence, &d.Status,
		&d.CreatedAt, &d.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Decision{}, fmt.Errorf("store: decision %q: %w", id, err)
	}
	if err != nil {
		return Decision{}, fmt.Errorf("store: get decision: %w", err)
	}
	d.LastVerifiedAt = nullStrPtr(lastVerified)
	return d, nil
}

// UpsertAssumption inserts or replaces an assumption by id.
func (s *Store) UpsertAssumption(a Assumption) (Assumption, error) {
	now := nowRFC3339()
	if a.ID == "" {
		a.ID = uuid.NewString()
	}
	if a.Status == "" {
		a.Status = StatusActive
	}
	if a.CreatedAt == "" {
		a.CreatedAt = now
	}
	a.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO assumptions(id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			body = excluded.body,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			status = excluded.status,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, a.ID, a.Title, a.Body, a.SourceType, a.Confidence, a.Status, a.CreatedAt, a.UpdatedAt, nullStr(a.LastVerifiedAt))
	if err != nil {
		return Assumption{}, fmt.Errorf("store: upsert assumption: %w", err)
	}
	out, err := s.GetAssumption(a.ID)
	if err != nil {
		return Assumption{}, err
	}
	if err := s.SyncEntityFTS("assumption", out.ID); err != nil {
		return Assumption{}, err
	}
	return out, nil
}

// GetAssumption loads an assumption by id.
func (s *Store) GetAssumption(id string) (Assumption, error) {
	var a Assumption
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at
		FROM assumptions WHERE id = ?
	`, id).Scan(
		&a.ID, &a.Title, &a.Body, &a.SourceType, &a.Confidence, &a.Status,
		&a.CreatedAt, &a.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Assumption{}, fmt.Errorf("store: assumption %q: %w", id, err)
	}
	if err != nil {
		return Assumption{}, fmt.Errorf("store: get assumption: %w", err)
	}
	a.LastVerifiedAt = nullStrPtr(lastVerified)
	return a, nil
}

// UpsertDiscovery inserts or replaces a discovery by id.
func (s *Store) UpsertDiscovery(d Discovery) (Discovery, error) {
	now := nowRFC3339()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	if d.Status == "" {
		d.Status = StatusActive
	}
	if d.Severity == "" {
		d.Severity = SeverityINFO
	}
	if d.CreatedAt == "" {
		d.CreatedAt = now
	}
	d.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO discoveries(id, title, body, source_type, confidence, status, severity, created_at, updated_at, last_verified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			body = excluded.body,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			status = excluded.status,
			severity = excluded.severity,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, d.ID, d.Title, d.Body, d.SourceType, d.Confidence, d.Status, d.Severity, d.CreatedAt, d.UpdatedAt, nullStr(d.LastVerifiedAt))
	if err != nil {
		return Discovery{}, fmt.Errorf("store: upsert discovery: %w", err)
	}
	out, err := s.GetDiscovery(d.ID)
	if err != nil {
		return Discovery{}, err
	}
	if err := s.SyncEntityFTS("discovery", out.ID); err != nil {
		return Discovery{}, err
	}
	return out, nil
}

// GetDiscovery loads a discovery by id.
func (s *Store) GetDiscovery(id string) (Discovery, error) {
	var d Discovery
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, title, body, source_type, confidence, status, severity, created_at, updated_at, last_verified_at
		FROM discoveries WHERE id = ?
	`, id).Scan(
		&d.ID, &d.Title, &d.Body, &d.SourceType, &d.Confidence, &d.Status, &d.Severity,
		&d.CreatedAt, &d.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Discovery{}, fmt.Errorf("store: discovery %q: %w", id, err)
	}
	if err != nil {
		return Discovery{}, fmt.Errorf("store: get discovery: %w", err)
	}
	d.LastVerifiedAt = nullStrPtr(lastVerified)
	return d, nil
}

// UpsertPlanChange inserts or replaces a plan_change by id.
func (s *Store) UpsertPlanChange(p PlanChange) (PlanChange, error) {
	now := nowRFC3339()
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	if p.Status == "" {
		p.Status = StatusActive
	}
	if p.CreatedAt == "" {
		p.CreatedAt = now
	}
	p.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO plan_changes(id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			body = excluded.body,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			status = excluded.status,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, p.ID, p.Title, p.Body, p.SourceType, p.Confidence, p.Status, p.CreatedAt, p.UpdatedAt, nullStr(p.LastVerifiedAt))
	if err != nil {
		return PlanChange{}, fmt.Errorf("store: upsert plan_change: %w", err)
	}
	out, err := s.GetPlanChange(p.ID)
	if err != nil {
		return PlanChange{}, err
	}
	if err := s.SyncEntityFTS("plan_change", out.ID); err != nil {
		return PlanChange{}, err
	}
	return out, nil
}

// GetPlanChange loads a plan_change by id.
func (s *Store) GetPlanChange(id string) (PlanChange, error) {
	var p PlanChange
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at
		FROM plan_changes WHERE id = ?
	`, id).Scan(
		&p.ID, &p.Title, &p.Body, &p.SourceType, &p.Confidence, &p.Status,
		&p.CreatedAt, &p.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return PlanChange{}, fmt.Errorf("store: plan_change %q: %w", id, err)
	}
	if err != nil {
		return PlanChange{}, fmt.Errorf("store: get plan_change: %w", err)
	}
	p.LastVerifiedAt = nullStrPtr(lastVerified)
	return p, nil
}

// UpsertClaim inserts or replaces a claim by id.
func (s *Store) UpsertClaim(c Claim) (Claim, error) {
	now := nowRFC3339()
	if c.ID == "" {
		c.ID = uuid.NewString()
	}
	if c.Status == "" {
		c.Status = StatusActive
	}
	if c.CreatedAt == "" {
		c.CreatedAt = now
	}
	c.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO claims(id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			body = excluded.body,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			status = excluded.status,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, c.ID, c.Title, c.Body, c.SourceType, c.Confidence, c.Status, c.CreatedAt, c.UpdatedAt, nullStr(c.LastVerifiedAt))
	if err != nil {
		return Claim{}, fmt.Errorf("store: upsert claim: %w", err)
	}
	out, err := s.GetClaim(c.ID)
	if err != nil {
		return Claim{}, err
	}
	if err := s.SyncEntityFTS("claim", out.ID); err != nil {
		return Claim{}, err
	}
	return out, nil
}

// GetClaim loads a claim by id.
func (s *Store) GetClaim(id string) (Claim, error) {
	var c Claim
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at
		FROM claims WHERE id = ?
	`, id).Scan(
		&c.ID, &c.Title, &c.Body, &c.SourceType, &c.Confidence, &c.Status,
		&c.CreatedAt, &c.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Claim{}, fmt.Errorf("store: claim %q: %w", id, err)
	}
	if err != nil {
		return Claim{}, fmt.Errorf("store: get claim: %w", err)
	}
	c.LastVerifiedAt = nullStrPtr(lastVerified)
	return c, nil
}

// UpsertEvidence inserts or replaces evidence by id.
func (s *Store) UpsertEvidence(e Evidence) (Evidence, error) {
	now := nowRFC3339()
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	if e.Status == "" {
		e.Status = StatusActive
	}
	if e.CreatedAt == "" {
		e.CreatedAt = now
	}
	e.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO evidence(id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			body = excluded.body,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			status = excluded.status,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, e.ID, e.Title, e.Body, e.SourceType, e.Confidence, e.Status, e.CreatedAt, e.UpdatedAt, nullStr(e.LastVerifiedAt))
	if err != nil {
		return Evidence{}, fmt.Errorf("store: upsert evidence: %w", err)
	}
	out, err := s.GetEvidence(e.ID)
	if err != nil {
		return Evidence{}, err
	}
	if err := s.SyncEntityFTS("evidence", out.ID); err != nil {
		return Evidence{}, err
	}
	return out, nil
}

// GetEvidence loads evidence by id.
func (s *Store) GetEvidence(id string) (Evidence, error) {
	var e Evidence
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, title, body, source_type, confidence, status, created_at, updated_at, last_verified_at
		FROM evidence WHERE id = ?
	`, id).Scan(
		&e.ID, &e.Title, &e.Body, &e.SourceType, &e.Confidence, &e.Status,
		&e.CreatedAt, &e.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Evidence{}, fmt.Errorf("store: evidence %q: %w", id, err)
	}
	if err != nil {
		return Evidence{}, fmt.Errorf("store: get evidence: %w", err)
	}
	e.LastVerifiedAt = nullStrPtr(lastVerified)
	return e, nil
}

// UpsertReview inserts or replaces a review by id. Empty Result stays open ("").
func (s *Store) UpsertReview(r Review) (Review, error) {
	now := nowRFC3339()
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	if r.Status == "" {
		r.Status = StatusActive
	}
	if r.CreatedAt == "" {
		r.CreatedAt = now
	}
	r.UpdatedAt = now

	_, err := s.db.Exec(`
		INSERT INTO reviews(id, title, body, source_type, confidence, status, result, created_at, updated_at, last_verified_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title,
			body = excluded.body,
			source_type = excluded.source_type,
			confidence = excluded.confidence,
			status = excluded.status,
			result = excluded.result,
			updated_at = excluded.updated_at,
			last_verified_at = excluded.last_verified_at
	`, r.ID, r.Title, r.Body, r.SourceType, r.Confidence, r.Status, r.Result, r.CreatedAt, r.UpdatedAt, nullStr(r.LastVerifiedAt))
	if err != nil {
		return Review{}, fmt.Errorf("store: upsert review: %w", err)
	}
	out, err := s.GetReview(r.ID)
	if err != nil {
		return Review{}, err
	}
	if err := s.SyncEntityFTS("review", out.ID); err != nil {
		return Review{}, err
	}
	return out, nil
}

// GetReview loads a review by id (includes result after mig 005).
func (s *Store) GetReview(id string) (Review, error) {
	var r Review
	var lastVerified sql.NullString
	err := s.db.QueryRow(`
		SELECT id, title, body, source_type, confidence, status, result, created_at, updated_at, last_verified_at
		FROM reviews WHERE id = ?
	`, id).Scan(
		&r.ID, &r.Title, &r.Body, &r.SourceType, &r.Confidence, &r.Status, &r.Result,
		&r.CreatedAt, &r.UpdatedAt, &lastVerified,
	)
	if err == sql.ErrNoRows {
		return Review{}, fmt.Errorf("store: review %q: %w", id, err)
	}
	if err != nil {
		return Review{}, fmt.Errorf("store: get review: %w", err)
	}
	r.LastVerifiedAt = nullStrPtr(lastVerified)
	return r, nil
}

// ListReviews returns all reviews ordered by created_at, then id.
func (s *Store) ListReviews() ([]Review, error) {
	rows, err := s.db.Query(`
		SELECT id, title, body, source_type, confidence, status, result, created_at, updated_at, last_verified_at
		FROM reviews
		ORDER BY created_at ASC, id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("store: list reviews: %w", err)
	}
	defer rows.Close()
	var out []Review
	for rows.Next() {
		var r Review
		var lastVerified sql.NullString
		if err := rows.Scan(
			&r.ID, &r.Title, &r.Body, &r.SourceType, &r.Confidence, &r.Status, &r.Result,
			&r.CreatedAt, &r.UpdatedAt, &lastVerified,
		); err != nil {
			return nil, fmt.Errorf("store: scan review: %w", err)
		}
		r.LastVerifiedAt = nullStrPtr(lastVerified)
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if out == nil {
		out = []Review{}
	}
	return out, nil
}
