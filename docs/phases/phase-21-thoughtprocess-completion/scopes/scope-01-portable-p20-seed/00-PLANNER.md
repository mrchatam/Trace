# P21-S01-00 — Planner: portable P20 seed

## Metadata
- id: P21-S01-00
- todo_ids: [P21-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock seed JSON v1 extension for **11 P20 cognition tables** (export + idempotent import). Retire P20 verify omit policy. **No product Go this row.**

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DECISION-LOG.md](../../DECISION-LOG.md) D-05
- [WORK-MAP.md](../../WORK-MAP.md) W-01
- P20 verify: [01-verify.md](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/01-verify.md) seed export bar (now **in scope**)
- Live: `internal/domain/seed_export.go`, `internal/domain/seed_import.go`, `cmd/trace/cli_test.go` (P17 keepers)

## Live inventory (2026-08-18)

| Surface | Location | S01 action |
|---------|----------|------------|
| `SeedDocument` | `seed_export.go` | **Extend** with P20 arrays (below) |
| `BuildSeedDocument` | `seed_export.go` | List + serialize 11 tables |
| `ImportSeedDocument` | `seed_import.go` | Idempotent upsert for each new array |
| P17 keepers | `cmd/trace/cli_test.go` | **Must stay green** (`TestSeedExportRoundTrip`, `TestSeedExportOmitsDeniedSurfaces`, `TestSeedExportWritesExportedAtCommit`) |
| Schema max | `internal/store/schema/` | **019** — no SQL migration required for export (tables exist 015–019) |
| Compat ceiling | `evals/compat/compat_test.go` | **19** unchanged unless S01 adds mig 020 (planner lock: **no mig 020**) |
| Omitted today | P20 verify list | `deliberation_state`, `uncertainties`, `hypotheses`, `decision_reconsiderations`, `changes`, `change_paths`, `effects`, `outcome_results`, `baselines`, `regressions`, `reflections` |

## FINAL locked defaults (S01-01 must not re-debate)

| Item | Value |
|------|-------|
| Seed version | **1** — additive JSON fields only (backward-compatible import of old seeds) |
| SQL migration | **None** — reuse 015–019 columns; forbid `020+` this scope |
| Compat ceiling | **19** (unchanged) |
| Denied surfaces | **Still omit:** transitions, task `work_state`, reviews, caps, tokens, index blobs (keep `TestSeedExportOmitsDeniedSurfaces`) |
| Child rows | `change_paths` nested under each `changes[]` entry; `effects[]` separate top-level keyed by `change_id` |
| Idempotency | Re-import same IDs → no duplicate rows; local `work_state` on tasks preserved (existing task import rule) |
| Git evidence | `exported_at_commit` unchanged (DF-85) |
| CLI | `trace seed export` / `trace seed import` — no new flags |

### New `SeedDocument` top-level keys (FINAL)

| JSON key | Source table | Notes |
|----------|--------------|-------|
| `deliberation_states` | `deliberation_state` | One row per task; include all columns |
| `uncertainties` | `uncertainties` | Full row |
| `hypotheses` | `hypotheses` | Full row |
| `decision_reconsiderations` | `decision_reconsiderations` | Full row |
| `changes` | `changes` + embedded `paths[]` | Paths from `change_paths` |
| `effects` | `effects` | All rows |
| `outcome_results` | `outcome_results` | All kinds test/verification/evaluation |
| `baselines` | `baselines` | Commit OID + scores_json only |
| `regressions` | `regressions` | Full row |
| `reflections` | `reflections` | Structured JSON array columns as strings |

Empty project → omit keys or `[]` (match existing seed style).

## Named tests (S01-01 must implement; S01-02 must re-run)

| Test | Proves |
|------|--------|
| `TestSeedExportIncludesP20Cognition` | Export after loop apply fixture includes all 11 keys with ≥1 row where seeded |
| `TestSeedImportP20RoundTrip` | import → export → fresh import preserves P20 IDs + deliberation phase |
| `TestSeedExportRoundTrip` | P17 keeper — plan tree + causal entities still round-trip |
| `TestSeedExportOmitsDeniedSurfaces` | P17 keeper — no transitions/work_state/reviews |
| `TestSeedExportWritesExportedAtCommit` | P17 keeper — git HEAD evidence |

## Store helpers (S01-01 — add `ListAll*` only; no schema change)

Export needs full-table scans (mirror `ListAllPlanPhases` pattern in `plan_hierarchy.go`):

| Helper | Table |
|--------|-------|
| `ListAllDeliberationStates` | `deliberation_state` |
| `ListAllUncertainties` | `uncertainties` |
| `ListAllHypotheses` | `hypotheses` |
| `ListAllDecisionReconsiderations` | `decision_reconsiderations` |
| `ListAllChanges` | `changes` |
| `ListAllEffects` | `effects` |
| `ListAllOutcomeResults` | `outcome_results` |
| `ListAllBaselines` | `baselines` |
| `ListAllRegressions` | `regressions` |
| `ListAllReflections` | `reflections` |

Reuse existing `ListChangePaths(changeID)` per change during export (no separate top-level `change_paths` key).

## Import order (FINAL — FK-safe)

After existing P17 entity import (goals → tasks → decisions → … → goal_plan_state):

1. `baselines`
2. `outcome_results` (may reference `baseline_id`, `task_id`, `goal_id`)
3. `changes` (+ nested `paths[]`)
4. `effects`
5. `uncertainties`, `hypotheses`, `decision_reconsiderations` (decisions already imported)
6. `regressions`, `reflections`
7. `deliberation_states` (last — one row per task)

## Policy supersession

**Retires** P20-S07-01 verify omit policy ([`01-verify.md`](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/01-verify.md) seed export bar). D-05 **promote** closes when S01-02 approves. Old seeds without P20 keys remain importable (additive v1).

## Touch files (FINAL)

- `internal/domain/seed_export.go` — structs + `BuildSeedDocument`
- `internal/domain/seed_import.go` — import helpers + `ImportSeedDocument` ordering
- `internal/domain/seed_*_test.go` (new or extend) — unit tests above
- `internal/store/{deliberation,cognitive,changes,outcomes,regressions}.go` — `ListAll*` helpers only
- `cmd/trace/cli_test.go` — extend round-trip fixture if needed (keep P17 keeper names)
- `trace/graph.json` — export after S01-01 (implementer Notes; not planner)

## Planner work

1. [x] Confirm 11 tables + no mig 020.
2. [x] Lock JSON key names + nested change_paths shape.
3. [x] Thicken `01-portable-p20-seed.md` + `02-scope-review.md`.
4. [x] Update `SCOPE-TODOS.md`.

## Exit criteria

- [x] 01/02 runnable alone with named tests + touch files
- [x] P17 keeper test names locked
- [x] Compat ceiling 19 locked
- [x] No product Go

## Next

**P21-S01-01** after this row is `done`.
