# Phase 45 / S04 / Review — Envelope + schema + descriptions

## Row `P45-S04-02` — **done** 2026-10-04.

> **Self-review, not an independent pass.** Implemented in one session by the coding agent. Evidence, not a second opinion.

## Honesty invariants

| # | Invariant | Holds? | Evidence |
|---|-----------|--------|----------|
| H1 | Every generic history list returns the same envelope shape (#157) | **yes** | `TestListEnvelopeIsUniform` sweeps changes/regressions/agents list+recommend/tasks for {ok, items, count, truncated}, items-is-array, count==len(items) |
| H2 | `truncated` is never a hardcoded guess | **yes** | Audit A9: page-fill signal in changes/regressions; `trace_tasks` uses the store's own `res.Truncated`; capability list detects overflow by fetching limit+1 |
| H3 | Every enum-constrained input is visible to the caller **in both** the schema and the description (#159) | **yes** | `TestToolDescriptionsCoverSchemaEnums` walks `toolListSchemas()` × `ToolDescriptions()`; drift fails the build |
| H4 | Every tool description is on the server, not just in test fixtures (#159) | **yes** | `TestServerExposesEveryToolDescription` |
| H5 | Out-of-range inputs are rejected, not silently clamped, where the code rejects (#158) | **yes** | SDK validates against the declared ranges (conformance `TestProtocolRejectsOutOfEnumValue` + `resolveExploreBounds` negatives; `TestExploreRejectsOutOfRange`) |
| H6 | Where the code clamps, the caller is told the effective value | **yes** | `limit` echoed in changes/regressions/tasks envelopes; search echoes effective limit (`TestSearchResponseEchoesEffectiveLimit`) |
| H7 | No PascalCase keys anywhere in a tool response (#161, class-wide) | **yes** | `TestNoToolResponseKeyContainsUppercase` sweeps 21 real tool calls with `walkJSON`; `TestChangeShowRowIsSnakeCase` pins the show DTO |
| H8 | Internal rowids are not on the wire | **yes** | `changeShowRow` has no `RowID` field |
| H9 | Empty lists marshal as `[]`, never `null` | **yes** | `newListEnvelope` contract + `TestListEnvelopeIsUniform` items-is-array assertion |
| H10 | Schemas resolve at construction, not first call | **yes** | `schemaFor` panics in `NewServer` on missing property / unresolvable schema; every test constructs the server |

## Commands run

```
CGO_ENABLED=1 go test ./internal/mcp/ -count=1
```

All green (including the three reshaped `mcp_test.go` tests and the audit follow-ups A7–A12).

## Residuals (accepted, stated)

- **Entity-specific collections keep domain-named keys** (`capabilities`, `missing`, `findings`) — CLI parity beats uniformity here; a cross-surface rename is its own decision.
- **`trace_changes compare`** returns the raw diff document (not a list; no envelope).
- **Ranges are enforced twice** (SDK schema validation + handler checks). The handler check is still load-bearing for direct library callers and for the `float64`→`int` conversion edge; removing it would trade redundancy for a hole.

## Verdict

**PASS.** #157, #158, #159, #161 closed by behaviour + permanent tests.