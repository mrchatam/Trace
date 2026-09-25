package store

import "testing"

func TestScopesCRUD(t *testing.T) {
	s, _ := openTempStore(t)

	sc, err := s.UpsertScope(Scope{
		Slug:  "auth-fe",
		Title: "Auth frontend",
		Kind:  ScopeKindFeature,
	})
	if err != nil {
		t.Fatalf("UpsertScope: %v", err)
	}
	if sc.ID == "" || sc.Slug != "auth-fe" || sc.Kind != ScopeKindFeature {
		t.Fatalf("unexpected scope: %+v", sc)
	}

	got, err := s.GetScope(sc.ID)
	if err != nil {
		t.Fatalf("GetScope: %v", err)
	}
	if got.Title != "Auth frontend" {
		t.Fatalf("GetScope title: %q", got.Title)
	}

	bySlug, err := s.GetScopeBySlug("auth-fe")
	if err != nil {
		t.Fatalf("GetScopeBySlug: %v", err)
	}
	if bySlug.ID != sc.ID {
		t.Fatalf("slug mismatch: %s vs %s", bySlug.ID, sc.ID)
	}

	_, err = s.UpsertScope(Scope{Slug: "x", Title: "y", Kind: "bogus"})
	if err == nil {
		t.Fatal("expected invalid kind rejection")
	}

	list, err := s.ListScopes()
	if err != nil || len(list) != 1 {
		t.Fatalf("ListScopes: %v %#v", err, list)
	}
}

func TestScopesMigrationApplied(t *testing.T) {
	s, _ := openTempStore(t)
	var name string
	err := s.db.QueryRow(`SELECT name FROM sqlite_master WHERE name=?`, "scopes").Scan(&name)
	if err != nil {
		t.Fatalf("scopes table missing: %v", err)
	}
	st, err := s.MigrationStatus()
	if err != nil {
		t.Fatal(err)
	}
	if st.MaxApplied < 29 {
		t.Fatalf("expected migration 29 applied, max=%d", st.MaxApplied)
	}
}
