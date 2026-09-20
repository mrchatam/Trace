package planner

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDeepPlanResultJSONSnakeCase(t *testing.T) {
	raw, err := json.Marshal(DeepPlanResult{
		RevisionID:       "rev-1",
		Document:         DeepPlanDocument{ScopeID: "scope-1"},
		SupersededCount:  2,
		LookaheadScopeID: "look-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	for _, want := range []string{
		`"revision_id"`,
		`"document"`,
		`"superseded_count"`,
		`"lookahead_scope_id"`,
	} {
		if !strings.Contains(s, want) {
			t.Fatalf("missing %s in %s", want, s)
		}
	}
	for _, bad := range []string{`"RevisionID"`, `"SupersededCount"`, `"LookaheadScopeID"`, `"Document"`} {
		if strings.Contains(s, bad) {
			t.Fatalf("unexpected PascalCase key %s in %s", bad, s)
		}
	}
}
