# S01 — FTS honesty — SCOPE-TODOS

Issues: [#149](https://github.com/mrchatam/Trace/issues/149), [#153](https://github.com/mrchatam/Trace/issues/153).

| ID | Item | Status |
|----|------|--------|
| S01-0 | P45-S01-00 thicken 01/02 | **done** 2026-10-04 |
| S01-1 | `MatchAnd`/`MatchOps` + `matcher` input + `graph_sync_honesty` on 3 tools | **done** 2026-10-04 (P45-S01-01) |
| S01-2 | Verification record (self-review) | **done** 2026-10-04 **PASS** (P45-S01-02) |

## Acceptance

- [x] `a OR b` under `matcher=ops` returns the union — `TestSearchFTSOpsBooleanReturnsUnion`
- [x] Any operator rewrite in `and` mode is reported — `TestSearchAndModeNotesOperatorRewrite`
- [x] Unknown matcher rejected, not defaulted — `TestSearchRejectsUnknownMatcher`
- [x] 64 cap holds in both modes — `TestSearchFTSOptsCapsLimit`
- [x] `SearchFTS` unchanged for all six pre-existing call sites
- [x] `graph_sync_honesty` on `trace_search`, `trace_explore`, `trace_impact` walk — #153
- [x] Single `openStore` per `trace_search` call