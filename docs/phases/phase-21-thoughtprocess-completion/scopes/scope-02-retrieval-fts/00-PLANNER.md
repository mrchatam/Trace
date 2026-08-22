# P21-S02-00 — Planner: retrieval + FTS for P20 types

## Metadata
- id: P21-S02-00
- todo_ids: [P21-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock retrieval Exact/Why/compiler + FTS sync for P20 entity types; eliminate INVESTIGATE stderr (`unknown entity type "uncertainty"`). **No product Go this row.**

## References
- [DECISION-LOG.md](../../DECISION-LOG.md) D-06, D-07
- [WORK-MAP.md](../../WORK-MAP.md) W-02, W-03
- P20 verify residual: [VERIFY-NOTES.md](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/VERIFY-NOTES.md) item 4
- Live: `internal/retrieval/{exact.go,why.go,types.go}`, `internal/store/{fts.go,cognitive.go}`, `internal/compiler/`, `internal/loop/next.go`

## Live inventory (2026-08-18 — S02-00 reconfirmed post-S01)

| Surface | Today | S02 action |
|---------|-------|------------|
| `lookupEntity` | goal…review,capability,file,symbol only | **Add** uncertainty, hypothesis, change, effect, regression, reflection, baseline, outcome_result |
| `NormalizeEntityType` | plan-change alias only | **Add** aliases: `uncertainty`, `hypothesis`, `change`, `regression`, `reflection`, `baseline`, `outcome`, `outcome_result` |
| `loadEntityText` / `RebuildFTS` | pre-P20 tables | **Add** uncertainties, hypotheses, changes (+ summary fields), regressions, reflections |
| `UpsertUncertainty/Hypothesis` | `cognitive.go` comment "Does not SyncEntityFTS" | **Wire** `SyncEntityFTS` after upsert |
| Loop INVESTIGATE | `TestLoopNextInvestigateEmphasizesUncertainties` passes; Expand fails on `uncertainty_blocks_task` link → `TaskContext` degrades to `minimalTaskContextPacket` | **Fix** `lookupEntity` for P20 types so Expand succeeds; add stderr assertion test |
| Schema | max **019** | **No mig 020** — FTS content from existing columns |
| Compat ceiling | **19** | Unchanged |

## FINAL locked defaults (S02-01 must not re-debate)

| Item | Value |
|------|-------|
| Canonical entity types (Exact/Why/FTS) | `uncertainty`, `hypothesis`, `change`, `effect`, `regression`, `reflection`, `baseline`, `outcome_result` |
| CLI `trace why` | Same types as Exact (via existing `cmd/trace/why.go`) |
| FTS title/body | uncertainty/hypothesis: title+body; change: reason+status; regression/reflection: summary + JSON excerpts truncated |
| Fail-closed | Unknown type → error at Exact/Why boundary (no silent empty for malformed type string) |
| INVESTIGATE packet | `loop next` with blocking uncertainty produces **no** `retrieval: unknown entity type` on stderr |
| Seed JSON | **No change** (S01 owns seed); retrieval only |
| Compat ceiling | **19** |

### Entity text mapping (FTS)

| type | table | title field | body field |
|------|-------|-------------|------------|
| uncertainty | uncertainties | title | body + severity + kind |
| hypothesis | hypotheses | title | body + status |
| change | changes | git_commit | reason + status |
| regression | regressions | dimension | summary + attribution |
| reflection | reflections | summary | invalidated_assumptions_json (truncated) |

`effect`, `baseline`, `outcome_result`: Exact/Why only (optional FTS — **Should** index baseline label + outcome summary if cheap).

## Named tests (S02-01 implement; S02-02 re-run)

| Test | Package | Proves |
|------|---------|--------|
| `TestExactLookupUncertainty` | `internal/retrieval` | Exact hit for seeded uncertainty |
| `TestExactLookupHypothesis` | `internal/retrieval` | Exact hit for hypothesis |
| `TestWhyUncertaintyIncludesGraphSteps` | `internal/retrieval` | Why chain for uncertainty seed |
| `TestNormalizeEntityTypeP20Aliases` | `internal/retrieval` | Aliases map to canonical |
| `TestSyncEntityFTSUncertainty` | `internal/store` | FTS row after upsert |
| `TestRebuildFTSIncludesP20Types` | `internal/store` | Rebuild indexes uncertainties |
| `TestLoopNextInvestigateNoRetrievalStderr` | `cmd/trace` or `internal/loop` | stderr clean on INVESTIGATE next |
| `TestCausalWhyContextRoundTrip` | `cmd/trace` | P17 why keeper still green |

## Touch files

- `internal/retrieval/exact.go` — `lookupEntity` cases
- `internal/retrieval/types.go` — `NormalizeEntityType`
- `internal/retrieval/why.go` — inherits lookupEntity
- `internal/store/fts.go` — `loadEntityText`, `RebuildFTS`
- `internal/store/cognitive.go` — SyncEntityFTS calls
- `internal/compiler/` — packet paths referencing entity types
- `internal/retrieval/*_test.go`, `internal/store/fts_test.go`
- `cmd/trace/loop_test.go` — stderr assertion test

## Exit criteria

- [x] 01/02 thickened with named tests + touch files
- [x] P20 type list locked
- [x] No mig 020
- [x] No product Go

## Next

**P21-S02-01**
