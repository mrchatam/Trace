package domain

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/mrchatam/Trace/internal/store"
)

// ClaimInput creates a Claim (honesty path; create+link).
type ClaimInput struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	LastVerifiedAt *string
}

// EvidenceInput creates Evidence (honesty path).
type EvidenceInput struct {
	ID             string
	Title          string
	Body           string
	SourceType     string
	Confidence     float64
	Status         string
	LastVerifiedAt *string
}

// CreateClaim persists a claim and appends entity.created.
func (s *Service) CreateClaim(ctx context.Context, in ClaimInput) (store.Claim, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, in.SourceType, in.Status)
	if err != nil {
		return store.Claim{}, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	c, err := s.store.UpsertClaim(store.Claim{
		ID:             id,
		Title:          strings.TrimSpace(in.Title),
		Body:           in.Body,
		SourceType:     src,
		Confidence:     in.Confidence,
		Status:         status,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.Claim{}, err
	}
	if err := s.appendCreated(EntityClaim, c.ID, c.Title); err != nil {
		return store.Claim{}, err
	}
	return c, nil
}

// CreateEvidence persists an evidence and appends entity.created.
func (s *Service) CreateEvidence(ctx context.Context, in EvidenceInput) (store.Evidence, error) {
	_ = ctx
	src, status, err := applyProvenance(in.Title, in.SourceType, in.Status)
	if err != nil {
		return store.Evidence{}, err
	}
	id := in.ID
	if id == "" {
		id = uuid.NewString()
	}
	e, err := s.store.UpsertEvidence(store.Evidence{
		ID:             id,
		Title:          strings.TrimSpace(in.Title),
		Body:           in.Body,
		SourceType:     src,
		Confidence:     in.Confidence,
		Status:         status,
		LastVerifiedAt: in.LastVerifiedAt,
	})
	if err != nil {
		return store.Evidence{}, err
	}
	if err := s.appendCreated(EntityEvidence, e.ID, e.Title); err != nil {
		return store.Evidence{}, err
	}
	return e, nil
}

// LinkClaimEvidence inserts entity_links rel=claim_has_evidence.
func (s *Service) LinkClaimEvidence(ctx context.Context, claimID, evidenceID string, meta LinkMeta) error {
	_ = ctx
	if claimID == "" || evidenceID == "" {
		return &ErrValidation{Msg: "claimID and evidenceID are required"}
	}
	if _, err := s.store.GetClaim(claimID); err != nil {
		return err
	}
	if _, err := s.store.GetEvidence(evidenceID); err != nil {
		return err
	}
	meta = meta.withDefaults()
	if _, err := s.store.InsertLink(store.EntityLink{
		FromType:   EntityClaim,
		FromID:     claimID,
		Rel:        RelClaimHasEvidence,
		ToType:     EntityEvidence,
		ToID:       evidenceID,
		SourceType: meta.SourceType,
		Confidence: meta.Confidence,
	}); err != nil {
		return err
	}
	return s.appendLinked(EntityClaim, claimID, RelClaimHasEvidence, EntityEvidence, evidenceID, meta)
}

// GetClaim / GetEvidence thin wrappers.

func (s *Service) GetClaim(ctx context.Context, id string) (store.Claim, error) {
	_ = ctx
	return s.store.GetClaim(id)
}

func (s *Service) GetEvidence(ctx context.Context, id string) (store.Evidence, error) {
	_ = ctx
	return s.store.GetEvidence(id)
}
