package agents_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mrchatam/Trace/internal/agents"
)

func defaultCatalogPath(t *testing.T) string {
	t.Helper()
	root := moduleRoot(t)
	return filepath.Join(root, "trace", "agents", "default.json")
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func TestDefaultCatalogValidJSON(t *testing.T) {
	doc, err := agents.LoadDefaultCatalog(defaultCatalogPath(t))
	if err != nil {
		t.Fatalf("LoadDefaultCatalog: %v", err)
	}
	if doc.SchemaVersion != 1 {
		t.Fatalf("schema_version: want 1 got %d", doc.SchemaVersion)
	}
	if doc.RegistryVersion == "" {
		t.Fatal("registry_version required")
	}
	if len(doc.Agents) != 6 {
		t.Fatalf("want 6 bundled agents, got %d", len(doc.Agents))
	}
	wantSlugs := map[string]bool{
		"agent:code-reviewer":        true,
		"agent:performance-reviewer": true,
		"agent:security-reviewer":    true,
		"agent:nested-reviewer":      true,
		"agent:explore":              true,
		"agent:generalPurpose":       true,
	}
	for _, a := range doc.Agents {
		if !wantSlugs[a.Slug] {
			t.Fatalf("unexpected bundled slug %q", a.Slug)
		}
		delete(wantSlugs, a.Slug)
	}
	if len(wantSlugs) != 0 {
		t.Fatalf("missing slugs: %v", wantSlugs)
	}
}

func TestEmbeddedDefaultCatalogMatchesFile(t *testing.T) {
	embedded, err := agents.LoadEmbeddedDefaultCatalog()
	if err != nil {
		t.Fatalf("LoadEmbeddedDefaultCatalog: %v", err)
	}
	fromFile, err := agents.LoadDefaultCatalog(defaultCatalogPath(t))
	if err != nil {
		t.Fatalf("LoadDefaultCatalog: %v", err)
	}
	if embedded.SchemaVersion != fromFile.SchemaVersion || embedded.RegistryVersion != fromFile.RegistryVersion {
		t.Fatalf("embedded vs file header mismatch: %+v vs %+v", embedded, fromFile)
	}
	if len(embedded.Agents) != len(fromFile.Agents) {
		t.Fatalf("agent count: embedded %d file %d", len(embedded.Agents), len(fromFile.Agents))
	}
}

func TestDefaultCatalogRejectsBadSchema(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{"schema_version":2,"registry_version":"x","agents":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := agents.LoadDefaultCatalog(path); err == nil {
		t.Fatal("schema_version 2 must fail")
	}
}
