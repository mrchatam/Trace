package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
)

func TestImportSeedScopeBadKind(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	badKinds := []string{"Feature", "", "domain", "layerz", "businessx", "FEATURE"}
	for _, k := range badKinds {
		_, _, err := svc.ImportSeedScope(ctx, domain.SeedScope{
			ID:    "11111111-1111-1111-1111-111111111111",
			Slug:  "bad-kind-" + k,
			Title: "Bad Kind",
			Kind:  k,
		})
		if err == nil {
			t.Fatalf("expected error for kind %q", k)
		}
		var ve *domain.ErrValidation
		if !errors.As(err, &ve) {
			t.Fatalf("kind %q: expected *ErrValidation, got %T: %v", k, err, err)
		}
		if ve.Msg == "" || !contains(ve.Msg, "feature, layer, or business") {
			t.Fatalf("kind %q: bad message: %q", k, ve.Msg)
		}
	}

	scopes, err := st.ListScopes()
	if err != nil {
		t.Fatal(err)
	}
	if len(scopes) != 0 {
		t.Fatalf("no scope row should be persisted for bad kind, got %d", len(scopes))
	}
}

func TestImportSeedScopeSlugCollisionDifferentID(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	existing, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug:  "collision",
		Title: "Existing",
		Kind:  "feature",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, _, err = svc.ImportSeedScope(ctx, domain.SeedScope{
		ID:    "22222222-2222-2222-2222-222222222222",
		Slug:  "collision",
		Title: "Different",
		Kind:  "layer",
	})
	if err == nil {
		t.Fatal("expected error for slug collision with different id")
	}
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("expected *ErrValidation, got %T: %v", err, err)
	}
	if ve.Msg == "" || !contains(ve.Msg, "already exists with id") || !contains(ve.Msg, existing.ID) || !contains(ve.Msg, "22222222-2222-2222-2222-222222222222") {
		t.Fatalf("bad message: %q", ve.Msg)
	}

	got, err := st.GetScopeBySlug("collision")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != existing.ID || got.Title != "Existing" || got.Kind != "feature" {
		t.Fatalf("original scope mutated: %+v", got)
	}
}

func TestImportSeedScopeReimportSameIDSlug(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	seed := domain.SeedScope{
		ID:    "33333333-3333-3333-3333-333333333333",
		Slug:  "reimport",
		Title: "Reimport Test",
		Kind:  "business",
	}

	sc1, inserted1, err := svc.ImportSeedScope(ctx, seed)
	if err != nil {
		t.Fatalf("first import: %v", err)
	}
	if !inserted1 {
		t.Fatal("first import should be inserted")
	}
	if sc1.ID != seed.ID || sc1.Slug != seed.Slug || sc1.Title != seed.Title || sc1.Kind != seed.Kind {
		t.Fatalf("first import mismatch: %+v", sc1)
	}

	sc2, inserted2, err := svc.ImportSeedScope(ctx, seed)
	if err != nil {
		t.Fatalf("re-import: %v", err)
	}
	if inserted2 {
		t.Fatal("re-import should not be inserted")
	}
	if sc2.ID != sc1.ID || sc2.Slug != sc1.Slug || sc2.Title != sc1.Title || sc2.Kind != sc1.Kind {
		t.Fatalf("re-import mismatch: %+v", sc2)
	}

	evs, err := st.ListEventsByEntity(domain.EntityScope, sc1.ID)
	if err != nil {
		t.Fatal(err)
	}
	createdCount := 0
	for _, e := range evs {
		if e.Type == domain.EventEntityCreated {
			createdCount++
		}
	}
	if createdCount != 1 {
		t.Fatalf("expected exactly one entity.created event, got %d", createdCount)
	}
}

func TestImportSeedScopeMissingSlugTitle(t *testing.T) {
	svc, _ := openDomain(t)
	ctx := context.Background()

	for _, tc := range []struct {
		name string
		sc   domain.SeedScope
	}{
		{"empty slug", domain.SeedScope{ID: "1", Slug: "", Title: "T", Kind: "feature"}},
		{"empty title", domain.SeedScope{ID: "1", Slug: "s", Title: "", Kind: "feature"}},
		{"both empty", domain.SeedScope{ID: "1", Slug: "", Title: "", Kind: "feature"}},
	} {
		_, _, err := svc.ImportSeedScope(ctx, tc.sc)
		if err == nil {
			t.Fatalf("%s: expected error", tc.name)
		}
		var ve *domain.ErrValidation
		if !errors.As(err, &ve) {
			t.Fatalf("%s: expected *ErrValidation, got %T: %v", tc.name, err, err)
		}
		if ve.Msg == "" || !contains(ve.Msg, "slug and title required") {
			t.Fatalf("%s: bad message: %q", tc.name, ve.Msg)
		}
	}
}

func TestImportSeedScopeEmptyIDExistingSlug(t *testing.T) {
	svc, st := openDomain(t)
	ctx := context.Background()

	existing, err := svc.CreateScope(ctx, domain.ScopeInput{
		Slug:  "existing-slug",
		Title: "Existing Scope",
		Kind:  "feature",
	})
	if err != nil {
		t.Fatal(err)
	}

	sc, inserted, err := svc.ImportSeedScope(ctx, domain.SeedScope{
		ID:    "",
		Slug:  "existing-slug",
		Title: "Different Title",
		Kind:  "layer",
	})
	if err != nil {
		t.Fatalf("import with empty id and existing slug: %v", err)
	}
	if inserted {
		t.Fatal("import with empty id and existing slug should not be inserted")
	}
	if sc.ID != existing.ID {
		t.Fatalf("expected reused id %s, got %s", existing.ID, sc.ID)
	}
	if sc.Slug != "existing-slug" || sc.Title != "Different Title" || sc.Kind != "layer" {
		t.Fatalf("scope fields not updated: %+v", sc)
	}

	scopes, err := st.ListScopes()
	if err != nil {
		t.Fatal(err)
	}
	if len(scopes) != 1 {
		t.Fatalf("expected exactly 1 scope, got %d", len(scopes))
	}

	evs, err := st.ListEventsByEntity(domain.EntityScope, existing.ID)
	if err != nil {
		t.Fatal(err)
	}
	createdCount := 0
	for _, e := range evs {
		if e.Type == domain.EventEntityCreated {
			createdCount++
		}
	}
	if createdCount != 1 {
		t.Fatalf("expected exactly one entity.created event, got %d", createdCount)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
