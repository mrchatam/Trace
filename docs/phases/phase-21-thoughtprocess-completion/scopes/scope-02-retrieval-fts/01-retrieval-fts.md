# P21-S02-01 — Implement: retrieval + FTS for P20 types

## Metadata
- id: P21-S02-01
- todo_ids: [P21-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective
Register all **8 P20 entity types** in Exact/Why/Expand (`lookupEntity`, `NormalizeEntityType`); extend FTS (`loadEntityText`, `RebuildFTS`, `ensureFTSPopulated`, cognitive upsert sync); eliminate INVESTIGATE `TaskContext` degradation caused by unknown-type Expand failures. **No seed changes** (S01 owns seed). **No mig 020.**

## Session start
Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: **status + notes only**.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- [DECISION-LOG.md](../../DECISION-LOG.md) D-06, D-07
- [WORK-MAP.md](../../WORK-MAP.md) W-02, W-03
- P20 residual: [VERIFY-NOTES.md](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/VERIFY-NOTES.md) item 4
- Live: `internal/retrieval/{exact.go,why.go,types.go,expand.go}`, `internal/store/{fts.go,cognitive.go}`, `internal/loop/{next.go,deliberation_packet.go}`, `cmd/trace/{why.go,loop_test.go}`

## Locked defaults (from S02-00 — do not re-debate)

| Item | Value |
|------|-------|
| Canonical types (Exact/Why/Expand/FTS entity_type) | `uncertainty`, `hypothesis`, `change`, `effect`, `regression`, `reflection`, `baseline`, `outcome_result` |
| Pre-P20 types | Unchanged — `goal`…`review`, `capability`, `file`, `symbol` |
| CLI `trace why` | Same types as Exact (via `cmd/trace/why.go` + `NormalizeEntityType`) |
| Aliases (`NormalizeEntityType`) | `plan-change`→`plan_change`; `outcome`→`outcome_result`; pass-through for canonical P20 strings |
| Fail-closed | Unknown type → error at Exact/Why/`SyncEntityFTS` boundary (no silent empty + stderr elsewhere) |
| INVESTIGATE packet | `loop next` with blocking uncertainty: exit 0, valid JSON, **no** `retrieval: unknown entity type` on stderr |
| Seed JSON | **No change** |
| Schema / compat | Max mig **019**; compat ceiling **19** |
| FTS Must index | `uncertainty`, `hypothesis`, `change`, `regression`, `reflection` |
| FTS Should (same PR if cheap) | `baseline` (label + truncated `scores_json`), `outcome_result` (summary + kind) |
| FTS skip | `effect` — Exact/Why only (no FTS row) |
| JSON/blob Law 1 | FTS body fields: text columns only; truncate JSON excerpts ≤512 runes; **never** index file blobs or path contents |

### Entity text mapping

**Exact/Why `lookupEntity` hit fields**

| type | Title | Excerpt |
|------|-------|---------|
| uncertainty | `title` | `body` (+ optional severity/kind prefix in excerpt) |
| hypothesis | `title` | `body` (+ status when non-empty) |
| change | `git_commit` (fallback: short id) | `reason` + `status` |
| effect | `dimension` | `comparison` + truncated `expected`/`actual` |
| regression | `dimension` | `summary` + `attribution` |
| reflection | `summary` | truncated `invalidated_assumptions_json` |
| baseline | `label` (fallback: `git_commit`) | truncated `scores_json` |
| outcome_result | `test_name` (fallback: `kind`) | `summary` |

**FTS `loadEntityText` / `RebuildFTS`**

| type | table | title field | body field |
|------|-------|-------------|------------|
| uncertainty | uncertainties | title | body + severity + kind (space-joined) |
| hypothesis | hypotheses | title | body + status |
| change | changes | git_commit | reason + status |
| regression | regressions | dimension | summary + attribution |
| reflection | reflections | summary | invalidated_assumptions_json (truncated ≤512) |
| baseline | baselines | label or git_commit | scores_json (truncated ≤512) |
| outcome_result | outcome_results | test_name or kind | summary |

Use a small shared `truncateText(s string, max int)` helper in `fts.go` if needed (do not over-abstract).

### INVESTIGATE stderr root cause (fix target)

1. Blocking uncertainty creates `uncertainty_blocks_task` link (domain).
2. `loop.BuildNextPacket` → `compiler.TaskContext` → `retrieval.Expand` → `neighbors` → `hitFromLinkNeighbor("uncertainty", …)` → `lookupEntity` **unknown type** today.
3. `next.go` degrades to `minimalTaskContextPacket` when blocking uncertainty (lines ~193–196) — packet valid but context thin.
4. **Fix:** register P20 types in `lookupEntity` so Expand succeeds; TaskContext returns full packet without degradation.

Optional doc-only: update comment in `deliberation_packet.go` line 43 (`ContextIncludeWhy: false` — why built separately in `next.go`; no longer "compiler lacks uncertainty").

## Requirements

1. **`lookupEntity`** — add `switch` cases for all 8 P20 types using existing store getters (`GetUncertainty`, `GetHypothesis`, `GetChange`, `GetEffect`, `GetRegression`, `GetReflection`, `GetBaseline`, `GetOutcomeResult`).
2. **`NormalizeEntityType`** — add aliases per locked table; emitted hits always use canonical form.
3. **`loadEntityText`** — extend map + custom SQL for non-`title`/`body` tables per mapping above.
4. **`RebuildFTS`** — append P20 Must (+ Should) tables after pre-P20 loop; reuse `rebuildTextTable` only where schema has `title`+`body`; custom rebuild helpers for `changes`, `regressions`, `reflections`, optional `baselines`/`outcome_results`.
5. **`ensureFTSPopulated`** — include P20 table counts in empty-index backfill guard.
6. **`SyncEntityFTS`** — call after `UpsertUncertainty` and `UpsertHypothesis` (remove fail-closed comments); return sync errors to caller.
7. **Tests** — implement all **8 named tests** below; keep P17/P19/P20 loop keepers green.
8. **No seed / schema / compat changes.**

## Touch files

- `internal/retrieval/exact.go` — `lookupEntity` cases
- `internal/retrieval/types.go` — `NormalizeEntityType`
- `internal/retrieval/why.go` — inherits lookupEntity (verify only)
- `internal/retrieval/expand.go` — inherits via hitFromLinkNeighbor (verify only)
- `internal/store/fts.go` — `loadEntityText`, `RebuildFTS`, `ensureFTSPopulated`, truncate helper
- `internal/store/cognitive.go` — `SyncEntityFTS` after uncertainty/hypothesis upserts
- `internal/loop/deliberation_packet.go` — comment only (optional)
- `internal/retrieval/exact_test.go` or extend `retrieval_test.go` — Exact/Why/Normalize tests
- `internal/store/fts_test.go` (new) — Sync + Rebuild tests
- `cmd/trace/loop_test.go` — `TestLoopNextInvestigateNoRetrievalStderr`

## Named tests (minimum)

| Test | Package | Proves |
|------|---------|--------|
| `TestExactLookupUncertainty` | `internal/retrieval` | Exact hit for seeded uncertainty by id |
| `TestExactLookupHypothesis` | `internal/retrieval` | Exact hit for hypothesis |
| `TestWhyUncertaintyIncludesGraphSteps` | `internal/retrieval` | Why on uncertainty seed returns steps (seed + ≥0 expand) |
| `TestNormalizeEntityTypeP20Aliases` | `internal/retrieval` | `outcome`→`outcome_result`, `plan-change`→`plan_change` |
| `TestSyncEntityFTSUncertainty` | `internal/store` | Upsert/sync → `SearchFTS` finds title token |
| `TestRebuildFTSIncludesP20Types` | `internal/store` | `RebuildFTS` indexes seeded uncertainty |
| `TestLoopNextInvestigateNoRetrievalStderr` | `cmd/trace` | Blocking uncertainty INVESTIGATE: exit 0, stderr lacks `retrieval: unknown entity type` |
| `TestCausalWhyContextRoundTrip` | `cmd/trace` | P17 why keeper still green |

Fixture hint: seed uncertainty via `domain.CreateUncertainty` with `TaskID` + `BLOCKING` severity (see `loop_test.go` helpers).

## Keeper floor

```bash
go test ./internal/retrieval/... -count=1 -run 'TestExactLookupUncertainty|TestExactLookupHypothesis|TestWhyUncertaintyIncludesGraphSteps|TestNormalizeEntityTypeP20Aliases'
go test ./internal/store/... -count=1 -run 'TestSyncEntityFTSUncertainty|TestRebuildFTSIncludesP20Types'
go test ./cmd/trace -count=1 -run 'TestLoopNextInvestigateNoRetrievalStderr|TestCausalWhyContextRoundTrip|TestLoopNextInvestigateEmphasizesUncertainties'
go test ./internal/deliberation/... -count=1
go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [ ] All 8 named tests PASS
- [ ] `./internal/deliberation/...` + P19/P20 loop keepers still PASS
- [ ] Compat ceiling **19** unchanged
- [ ] No mig 020+; no seed export/import edits
- [ ] Board row Notes: test command output summary

## Minimal todos

- [ ] Extend `lookupEntity` + `NormalizeEntityType`
- [ ] Extend FTS load/rebuild/sync + cognitive upsert wiring
- [ ] Add 8 named tests
- [ ] Run keeper floor + board Notes
