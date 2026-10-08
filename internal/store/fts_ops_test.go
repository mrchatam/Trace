package store

import (
	"strings"
	"testing"
)

// TestCompileFTSQueryModes locks the AND default (DF-87) and the opt-in
// boolean subset. AND mode must stay byte-identical to the pre-existing
// sanitizer so existing callers and cursors are unaffected.
func TestCompileFTSQueryModes(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		mode        MatchMode
		wantMatch   string
		wantOps     []string
		wantSawOp   bool
		wantFBEmpty bool
	}{
		{
			name:      "and default slash is a separator",
			in:        "GET /notes",
			mode:      MatchAnd,
			wantMatch: `"GET" AND "notes"`,
			wantOps:   []string{"GET", "notes"},
		},
		{
			name:      "and default treats OR as a term",
			in:        "zeppelin OR cache",
			mode:      MatchAnd,
			wantMatch: `"zeppelin" AND "OR" AND "cache"`,
			wantOps:   []string{"zeppelin", "OR", "cache"},
			wantSawOp: true,
		},
		{
			name:      "ops honors OR",
			in:        "zeppelin OR cache",
			mode:      MatchOps,
			wantMatch: `("zeppelin" OR "cache")`,
			wantOps:   []string{"zeppelin", "cache"},
			wantSawOp: true,
		},
		{
			name:      "ops folds mixed operators left to right",
			in:        `a OR b AND c`,
			mode:      MatchOps,
			wantMatch: `(("a" OR "b") AND "c")`,
			wantOps:   []string{"a", "b", "c"},
			wantSawOp: true,
		},
		{
			name:      "ops supports NOT",
			in:        "NOT alpha",
			mode:      MatchOps,
			wantMatch: `NOT ("alpha")`,
			wantOps:   []string{"alpha"},
			wantSawOp: true,
		},
		{
			name:      "ops supports NEAR",
			in:        "alpha NEAR(4) beta",
			mode:      MatchOps,
			wantMatch: `("alpha" NEAR(4) "beta")`,
			wantOps:   []string{"alpha", "beta"},
			wantSawOp: true,
		},
		{
			name:      "ops supports phrases and prefixes",
			in:        `"exact phrase" OR cach*`,
			mode:      MatchOps,
			wantMatch: `("exact phrase" OR "cach"*)`,
			wantOps:   []string{"exact phrase", "cach"},
			wantSawOp: true,
		},
		{
			name:      "ops leaves plain prose alone",
			in:        "cache eviction",
			mode:      MatchOps,
			wantMatch: `"cache" AND "eviction"`,
			wantOps:   []string{"cache", "eviction"},
		},
		{
			name:        "ops falls back on a dangling operator",
			in:          "alpha AND",
			mode:        MatchOps,
			wantMatch:   `"alpha" AND "AND"`,
			wantOps:     []string{"alpha", "AND"},
			wantSawOp:   true,
			wantFBEmpty: false,
		},
		{
			name:        "ops falls back on a doubled operator",
			in:          "alpha AND OR beta",
			mode:        MatchOps,
			wantMatch:   `"alpha" AND "AND" AND "OR" AND "beta"`,
			wantOps:     []string{"alpha", "AND", "OR", "beta"},
			wantSawOp:   true,
			wantFBEmpty: false,
		},
		{
			name:      "empty query compiles to nothing",
			in:        "   ",
			mode:      MatchOps,
			wantMatch: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CompileFTSQuery(tc.in, tc.mode)
			if got.Match != tc.wantMatch {
				t.Errorf("Match = %q, want %q", got.Match, tc.wantMatch)
			}
			if got.SawOperator != tc.wantSawOp {
				t.Errorf("SawOperator = %v, want %v", got.SawOperator, tc.wantSawOp)
			}
			if tc.wantFBEmpty && got.Fallback == "" {
				t.Error("want a Fallback notice, got none")
			}
			if !tc.wantFBEmpty && got.Fallback != "" && got.Match != tc.wantMatch {
				t.Errorf("unexpected Fallback %q", got.Fallback)
			}
			if len(got.Operands) != len(tc.wantOps) {
				t.Fatalf("Operands = %v, want %v", got.Operands, tc.wantOps)
			}
			for i := range tc.wantOps {
				if got.Operands[i] != tc.wantOps[i] {
					t.Errorf("Operands[%d] = %q, want %q", i, got.Operands[i], tc.wantOps[i])
				}
			}
		})
	}
}

// TestBuildOpsMatchQuotedLiterals proves no caller text reaches MATCH as
// syntax: punctuation inside a term is stripped by the letter/number
// sanitizer, and only the closed operator vocabulary is emitted.
func TestBuildOpsMatchQuotedLiterals(t *testing.T) {
	expr, ok := buildOpsMatch(tokenizeOpsQuery(`cache OR "x'); DROP TABLE--"`))
	if !ok {
		t.Fatal("want a usable expression")
	}
	if want := `("cache" OR "x DROP TABLE")`; expr != want {
		t.Fatalf("expr = %q, want %q", expr, want)
	}
	// The injected punctuation must not survive inside the quoted term. The
	// trailing ")" belongs to the builder's own grouping, so check the operand.
	for _, forbidden := range []string{"');", "--", `x');`} {
		if strings.Contains(expr, forbidden) {
			t.Fatalf("expression leaked %q: %s", forbidden, expr)
		}
	}
}

// TestSearchFTSOpsBooleanReturnsUnion is the behavior proof for #149: the same
// query that returns nothing under the default AND mode returns the union
// under matcher=ops, and the AND mode still reports the rewrite.
func TestSearchFTSOpsBooleanReturnsUnion(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	if _, err := s.UpsertTask(Task{Title: "zeppelin cache layer"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.UpsertTask(Task{Title: "unrelated widget"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.UpsertTask(Task{Title: "cache eviction policy"}); err != nil {
		t.Fatal(err)
	}

	andHits, andCompiled, err := s.SearchFTSOpts(SearchFTSOpts{Query: "zeppelin OR cache", Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(andHits) != 0 {
		t.Fatalf("AND mode should find nothing for this query, got %d hits", len(andHits))
	}
	if !andCompiled.SawOperator {
		t.Error("AND mode must report that it saw an operator")
	}

	opsHits, opsCompiled, err := s.SearchFTSOpts(SearchFTSOpts{Query: "zeppelin OR cache", Limit: 10, Mode: MatchOps})
	if err != nil {
		t.Fatal(err)
	}
	if len(opsHits) != 2 {
		t.Fatalf("ops mode want 2 union hits, got %d", len(opsHits))
	}
	if opsCompiled.Fallback != "" {
		t.Fatalf("unexpected fallback: %q", opsCompiled.Fallback)
	}

	onlyZeppelin, _, err := s.SearchFTSOpts(SearchFTSOpts{Query: "zeppelin NOT cache", Limit: 10, Mode: MatchOps})
	if err != nil {
		t.Fatal(err)
	}
	if len(onlyZeppelin) != 0 {
		t.Fatalf("NOT should exclude the matching row, got %d", len(onlyZeppelin))
	}
}

// TestParseMatchModeRejectsUnknown keeps the wire vocabulary closed.
func TestParseMatchModeRejectsUnknown(t *testing.T) {
	for _, in := range []string{"", "and", "AND", " ops "} {
		if _, err := ParseMatchMode(in); err != nil {
			t.Errorf("ParseMatchMode(%q) unexpected error: %v", in, err)
		}
	}
	if _, err := ParseMatchMode("regex"); err == nil {
		t.Fatal("want an error for an unknown matcher")
	}
}

// TestSearchFTSOptsCapsLimit keeps the 64 cap on both paths.
func TestSearchFTSOptsCapsLimit(t *testing.T) {
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if _, err := s.UpsertTask(Task{Title: "cap token"}); err != nil {
		t.Fatal(err)
	}
	hits, compiled, err := s.SearchFTSOpts(SearchFTSOpts{Query: "cap", Limit: 5000, Mode: MatchOps})
	if err != nil {
		t.Fatal(err)
	}
	if len(hits) != 1 || compiled.Match == "" {
		t.Fatalf("hits=%d compiled=%+v", len(hits), compiled)
	}
}
