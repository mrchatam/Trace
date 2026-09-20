package domain_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
)

func TestListPromotionCandidatesLimitedTruncates(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()
	const n = domain.DefaultPromotionCandidateLimit + 5
	for i := 0; i < n; i++ {
		if _, err := svc.CreateDiscovery(ctx, domain.DiscoveryInput{
			Title:    fmt.Sprintf("block-%d", i),
			Severity: domain.SeverityBlocking,
		}); err != nil {
			t.Fatal(err)
		}
	}
	items, truncated, err := svc.ListPromotionCandidatesLimited(domain.DefaultPromotionCandidateLimit)
	if err != nil {
		t.Fatal(err)
	}
	if !truncated {
		t.Fatal("expected truncated=true")
	}
	if len(items) != domain.DefaultPromotionCandidateLimit {
		t.Fatalf("len=%d want %d", len(items), domain.DefaultPromotionCandidateLimit)
	}
	// Unlimited wrapper still capped at default.
	all, err := svc.ListPromotionCandidates()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != domain.DefaultPromotionCandidateLimit {
		t.Fatalf("ListPromotionCandidates len=%d want default cap %d", len(all), domain.DefaultPromotionCandidateLimit)
	}
}
