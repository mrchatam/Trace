# Phase 45 / S01 / Review — FTS honesty

## Row `P45-S01-02` — **done** 2026-10-04.

> **Self-review, not an independent pass.** This phase was implemented in one session by the coding agent; no separate reviewer sub-agent examined it. What follows is the verification actually run and the invariants actually checked — evidence, not a second opinion.

## Honesty invariants

| # | Invariant | Holds? | Evidence |
|---|-----------|--------|----------|
| H1 | No query narrowing is silent — every path that changes the meaning of a query emits a notice | **yes** | `CompileFTSQuery` sets `Fallback` in all five rejection branches (`fts.go:148–178`); `TestSearchAndModeNotesOperatorRewrite` asserts the notice reaches the MCP response body |
| H2 | `matcher=ops` genuinely changes results vs `and` (i.e. it is not a no-op decorator) | **yes** | `TestSearchFTSOpsBooleanReturnsUnion` + `TestSearchOpsModeReturnsUnion` — `a OR b` returns the union, not the intersection |
| H3 | Unknown matcher fails closed rather than defaulting | **yes** | `ParseMatchMode` returns an error (`fts.go:60`); `TestSearchRejectsUnknownMatcher` at the tool boundary |
| H4 | `SearchFTS` behaviour is unchanged for all pre-existing callers | **yes** | Wrapper delegates with `Mode: MatchAnd`; all six call sites compile unmodified; store + CLI tests green |
| H5 | The 64-result cap holds in both modes | **yes** | `TestSearchFTSOptsCapsLimit` |
| H6 | All three walk tools report index staleness (#153) | **yes** | `GraphSyncHonestyFor` exported from `internal/compiler`; search, explore and impact all emit it; `TestSearchResponseCarriesGraphSyncHonesty` |
| H7 | Operators cannot be injected through a quoted literal | **yes** | `TestBuildOpsMatchQuotedLiterals` — `"a OR b"` is one term |
| H8 | Expression precedence is explicit, not implied | **yes** | `buildOpsMatch` fully parenthesizes; `TestCompileFTSQueryModes` covers mixed `OR`/`AND` |
| H9 | One call ⇒ one `openStore` | **yes** | `tools_search.go` rewritten around a single store handle; the pre-existing double-open at `tools_search.go:28` is gone |
| H10 | No cap raised | **yes** | `grep` — no numeric limit constant added or widened in this scope |

## Commands run

```
CGO_ENABLED=1 go test ./internal/store/ -count=1 -run 'TestCompileFTSQueryModes|TestBuildOpsMatchQuotedLiterals|TestSearchFTSOpsBooleanReturnsUnion|TestParseMatchModeRejectsUnknown|TestSearchFTSOptsCapsLimit'
CGO_ENABLED=1 go test ./internal/mcp/ -count=1 -run 'TestSearch|TestExplore'
CGO_ENABLED=1 go test ./internal/store/ ./internal/retrieval/ ./internal/compiler/ ./internal/mcp/ -count=1
```

All green. DF-87 keepers re-run unchanged and green.

## Residuals (accepted, not defects)

- **`ops` is not the default.** An agent that passes `foo OR bar` without `matcher=ops` still gets the intersection — but it now gets told so, in the response, every time. This is the deliberate D1 trade: correctness-by-default for existing callers, full power opt-in. Revisiting it is a breaking change and belongs in a phase that can migrate callers.
- **`NEAR` upper bound of 64** is this repo's own cap (bounds expression size); SQLite itself places no 10-distance limit — `NEAR`'s *default* is 10 when no N is supplied.
- Lowercase `near(n)` is accepted by `nearOpsRe` but is **not** treated as an operator; it falls through as a literal term, consistent with the uppercase-only policy.

## Verdict

**PASS.** #149 and #153 closed by behaviour + permanent tests.