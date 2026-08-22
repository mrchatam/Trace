package retrieval

import (
	"context"

	"github.com/mrchatam/Trace/internal/domain"
)

// DomainImpactWalker adapts Engine.ImpactWalk to domain.ImpactWalker (breaks import cycle).
type DomainImpactWalker struct {
	Engine *Engine
}

// ImpactWalk implements domain.ImpactWalker.
func (w DomainImpactWalker) ImpactWalk(ctx context.Context, seeds []domain.ImpactSeedRef, depth int) (*domain.ImpactWalkSnapshot, error) {
	if w.Engine == nil {
		return nil, &domain.ErrValidation{Msg: "impact engine is nil"}
	}
	in := make([]ImpactSeed, 0, len(seeds))
	for _, s := range seeds {
		in = append(in, ImpactSeed{EntityType: s.EntityType, EntityID: s.EntityID})
	}
	res, err := w.Engine.ImpactWalk(ctx, in, depth)
	if err != nil {
		return nil, err
	}
	out := &domain.ImpactWalkSnapshot{
		BlastTotal: res.BlastTotal,
		BlastKept:  res.BlastKept,
		Truncated:  res.Truncated,
		Depth:      res.Depth,
	}
	for _, s := range res.Seeds {
		out.Seeds = append(out.Seeds, domain.ImpactSeedRef{EntityType: s.EntityType, EntityID: s.EntityID})
	}
	for _, h := range res.Blast {
		out.Blast = append(out.Blast, domain.ImpactEntityRef{EntityType: h.EntityType, EntityID: h.EntityID})
	}
	for _, h := range res.AffectedTests {
		out.AffectedTests = append(out.AffectedTests, domain.ImpactEntityRef{EntityType: h.EntityType, EntityID: h.EntityID})
	}
	return out, nil
}

// WireDomainImpactWalker configures svc for predict/compare walks.
func WireDomainImpactWalker(svc *domain.Service, eng *Engine) {
	if svc == nil || eng == nil {
		return
	}
	svc.SetImpactWalker(DomainImpactWalker{Engine: eng})
}
