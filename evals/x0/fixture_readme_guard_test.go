package x0

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Stable fixture UUIDs from seed/gt.json — must not appear in agent-facing README (GC-02).
var fixtureOracleUUIDs = []string{
	"11111111-1111-1111-1111-111111111111",
	"22222222-2222-2222-2222-222222222222",
	"33333333-3333-3333-3333-333333333333",
	"44444444-4444-4444-4444-444444444444",
	"55555555-5555-5555-5555-555555555555",
}

// TestFixtureReadmeHasNoGTUUIDOracle asserts fixtures/x0/README.md does not
// publish the ground-truth UUID table (Gate C GC-02 fairness residual).
func TestFixtureReadmeHasNoGTUUIDOracle(t *testing.T) {
	root := findModuleRoot()
	if root == "" {
		t.Fatal("module root not found")
	}
	path := filepath.Join(root, "fixtures", "x0", "README.md")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	body := string(raw)
	for _, id := range fixtureOracleUUIDs {
		if strings.Contains(body, id) {
			t.Errorf("fixtures/x0/README.md must not contain GT UUID %s (agent oracle)", id)
		}
	}
	// Evaluator map must exist under evals/x0/.
	gtMap := filepath.Join(root, "evals", "x0", "GT-MAP.md")
	if _, err := os.Stat(gtMap); err != nil {
		t.Fatalf("evaluator GT map missing: %v", err)
	}
}
