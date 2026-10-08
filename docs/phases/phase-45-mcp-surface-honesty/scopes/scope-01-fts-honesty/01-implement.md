# Phase 45 / S01 / Implement — FTS honesty

## Row `P45-S01-01` — **done** 2026-10-04. Issues [#149](https://github.com/mrchatam/Trace/issues/149), [#153](https://github.com/mrchatam/Trace/issues/153).

## 1. Contract

```go
type MatchMode int
const (
    MatchAnd MatchMode = iota  // legacy: operators stripped, terms AND-joined
    MatchOps                   // opt-in: real FTS5 boolean expression
)
func ParseMatchMode(s string) (MatchMode, error)   // "" → MatchAnd; unknown → error

type FTSQuery struct {
    Match       string   // what was actually sent to FTS5
    Operands    []string // user-supplied terms, in order
    SawOperator bool     // true when AND/OR/NOT/NEAR/phrase/prefix appeared
    Fallback    string   // non-empty when MatchMode could not be honored
}

type SearchFTSOpts struct { Query string; Limit int; Mode MatchMode }
func (s *Store) SearchFTSOpts(o SearchFTSOpts) ([]FTSHit, FTSQuery, error)
func CompileFTSQuery(q string, mode MatchMode) FTSQuery
func (s *Store) SearchFTS(query string, limit int) ([]FTSHit, error)  // unchanged wrapper
```

`SearchFTS` keeps its exact signature and delegates with `Mode: MatchAnd`. It is a **wrapper**, not a replacement — all six pre-existing call sites (`internal/retrieval/{concept,search,exact}.go`, `internal/store/{store,fts}_test.go`, `cmd/trace/cli_test.go`) are untouched.

## 2. New file — `internal/store/fts_ops.go`

| Symbol | Role |
|--------|------|
| `opsTokKind`, `opsTok` | Token stream: term / op / phrase / prefix / `NEAR(n)` |
| `nearOpsRe` | `^(?:NEAR\|Near\|near)\((\d{1,3})\)$` |
| `tokenizeOpsQuery(q) []opsTok` | Uppercase-only operator recognition |
| `buildOpsMatch(toks) (string, bool)` | Emits a **fully parenthesized** expression; `false` ⇒ caller must fall back |
| `opsOperands(toks) []string` | Extracts the bare terms for the response echo |

**Rejections** (all ⇒ `Fallback` set, expression is the AND form, response carries `query_note`):

- dangling operator (`a OR`, `a AND`, trailing `NOT`)
- doubled operator (`a OR OR b`, `AND a`)
- `NEAR` distance above `[1,64]` (own bound; SQLite itself accepts larger distances, this repo caps at 64 to bound expression size) or with no right operand
- doubled `NOT` (`NOT NOT a`)
- unbalanced quote

Full parenthesization matters: `a OR b AND c` means `(a OR b) AND c` to FTS5 but the opposite to a reader. Emitting `(a OR b) AND c` explicitly removes the ambiguity rather than relying on precedence the caller cannot see.

## 3. Modified

| File | Change |
|------|--------|
| `internal/store/fts.go` | `MatchMode`, `ParseMatchMode`, `String()`, `OpsFallbackNotice`, `FTSQuery`, `SearchFTSOpts`, `(*Store).SearchFTSOpts`, `CompileFTSQuery`, `ftsTokens`, `quoteFTS`. `SearchFTS` → wrapper. **64 limit cap preserved in both modes.** |
| `internal/retrieval/types.go` | `SearchOptions.Matcher store.MatchMode` |
| `internal/retrieval/search.go` | New `SearchWithQuery(...) (*Result, FTSQuery, error)`. `hitsFromFTS`/`appendUniqueHits`/`sortHits` refactored to take the compiled form — no second compile, no second `openStore`. |
| `internal/compiler/graph_sync_honesty.go` | `buildGraphSyncHonesty` → **exported** `GraphSyncHonestyFor(st)` so `internal/mcp` can reach it without duplicating the policy |
| `internal/compiler/explore.go` | `ExploreOpts.Matcher`; `ExploreResult.GraphSyncHonesty` |
| `internal/mcp/tools_search.go` | Rewritten. **Single `openStore`** (was two, `tools_search.go:28`). Echoes `limit` + `matcher`; adds `query_note` and `graph_sync_honesty`. |
| `internal/mcp/tools_impact.go` | Walk envelope (`~:291`) now emits `graph_sync_honesty` |
| `internal/mcp/tools_explore.go` | Rewritten: `Matcher` input + `resolveExploreBounds` range validation (#158 preview; completed in S04) |

## 4. Acceptance — all PASS

`internal/store/fts_ops_test.go`

| Test | Asserts |
|------|---------|
| `TestCompileFTSQueryModes` | `and` mode strips operators; `ops` mode preserves them; unknown mode errors |
| `TestBuildOpsMatchQuotedLiterals` | `"phrase"` stays one term; user input cannot inject operators |
| `TestSearchFTSOpsBooleanReturnsUnion` | `a OR b` returns the **union** (the #149 false-negative) |
| `TestParseMatchModeRejectsUnknown` | Unknown matcher is an error, not a silent `and` |
| `TestSearchFTSOptsCapsLimit` | The 64 cap holds in **both** modes |

`internal/mcp/tools_search_test.go`

| Test | Asserts |
|------|---------|
| `TestSearchAndModeNotesOperatorRewrite` | `matcher=and` + `OR` in the query ⇒ response says so (**never silent**) |
| `TestSearchOpsModeReturnsUnion` | MCP end-to-end union |
| `TestSearchRejectsUnknownMatcher` | Rejected at the tool boundary |
| `TestSearchResponseEchoesEffectiveLimit` | Cap is visible to the caller |
| `TestSearchResponseCarriesGraphSyncHonesty` | #153 on `trace_search` |
| `TestExploreRejectsOutOfRange` | Out-of-range bounds refused |

**DF-87 keepers re-run unchanged and green** (no regression in the existing search/explore behaviour).

## 5. Boundaries honoured

- No cap raised (Laws 6–7): the 64 limit is unchanged, not raised to accommodate `ops`.
- No silent fallback anywhere: `Fallback` non-empty ⇒ the response carries the notice.
- `SearchFTS` signature unchanged ⇒ all six existing call sites compile and pass unmodified.
- No `web/`; no `cmd/trace` changes.

## Next

`P45-S01-02`