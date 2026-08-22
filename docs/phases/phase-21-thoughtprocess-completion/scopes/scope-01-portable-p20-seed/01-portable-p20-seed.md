# P21-S01-01 — Implement: portable P20 seed

## Metadata
- id: P21-S01-01
- todo_ids: [P21-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective
Extend seed export/import so a cloned repo can round-trip all **11 P20 cognition tables** per S01-00 FINAL locks. **Supersedes** P20 verify omit policy (D-05).

## Session start
Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Board edits: **status + notes only**.

## References
- [00-PLANNER.md](00-PLANNER.md) — FINAL locks (authoritative)
- [DECISION-LOG.md](../../DECISION-LOG.md) D-05
- Live: `internal/domain/seed_export.go`, `seed_import.go`, `cmd/trace/cli_test.go`
- Schema: `internal/store/schema/015`–`019` (max **019**, no **020+**)

## Locked defaults (from S01-00 — do not re-debate)

| Item | Value |
|------|-------|
| Seed version | **1** — additive JSON fields only |
| SQL migration | **None** — tables exist 015–019 |
| Compat ceiling | **19** (unchanged) |
| Denied surfaces | **Still omit:** transitions, task `work_state`, reviews, caps, tokens, index blobs |
| Nested paths | `change_paths` embedded as `paths[]` on each `changes[]` entry |
| Idempotency | Re-import same IDs → no duplicate rows; local `work_state` preserved on tasks |
| Git evidence | `exported_at_commit` unchanged (DF-85) |
| CLI | `trace seed export` / `trace seed import` — no new flags |
| Empty project | Omit keys or `[]` (match existing seed style) |

### New `SeedDocument` top-level keys (FINAL)

| JSON key | Source table(s) | Notes |
|----------|-----------------|-------|
| `deliberation_states` | `deliberation_state` | All columns; snake_case JSON |
| `uncertainties` | `uncertainties` | Full row |
| `hypotheses` | `hypotheses` | Full row |
| `decision_reconsiderations` | `decision_reconsiderations` | Full row |
| `changes` | `changes` + `paths[]` | Paths from `change_paths`; git SHA + path refs only |
| `effects` | `effects` | Top-level; keyed by `change_id` |
| `outcome_results` | `outcome_results` | All kinds: test / verification / evaluation |
| `baselines` | `baselines` | `git_commit` + `scores_json` (+ label metadata) |
| `regressions` | `regressions` | Full row |
| `reflections` | `reflections` | JSON array columns as string fields |

JSON field names mirror store column snake_case (e.g. `task_id`, `git_commit`, `scores_json`, `last_verified_at` optional/nullable).

## Requirements

1. Add `Seed*` structs + extend `SeedDocument` with the 10 keys above (11 tables; paths nested).
2. Add `ListAll*` store helpers per S01-00 table (no schema SQL).
3. Export: list all rows; for each change call `ListChangePaths`; embed as `paths[]`.
4. Import in dependency-safe order (after existing P17 import):
   - baselines → outcome_results → changes (+ paths) → effects → uncertainties/hypotheses/decision_reconsiderations → regressions/reflections → deliberation_states
5. Idempotent upsert by primary key (reuse existing upsert patterns in store).
6. Implement **all named tests** from S01-00.
7. Run P17 keepers + new tests before marking `done`.
8. Export `trace/graph.json` once (Notes: path + git SHA).

## Touch files

- `internal/domain/seed_export.go` — structs + `BuildSeedDocument`
- `internal/domain/seed_import.go` — import helpers + `ImportSeedDocument` ordering
- `internal/domain/seed_export_test.go` / `seed_import_test.go` (new or extend)
- `internal/store/deliberation.go` — `ListAllDeliberationStates`
- `internal/store/cognitive.go` — `ListAllUncertainties`, `ListAllHypotheses`, `ListAllDecisionReconsiderations`
- `internal/store/changes.go` — `ListAllChanges`, `ListAllEffects`
- `internal/store/outcomes.go` — `ListAllOutcomeResults`, `ListAllBaselines`
- `internal/store/regressions.go` — `ListAllRegressions`, `ListAllReflections`
- `cmd/trace/cli_test.go` — only if round-trip fixture needs P20 rows (keep P17 keeper names)

## Named tests (minimum)

| Test | Proves |
|------|--------|
| `TestSeedExportIncludesP20Cognition` | Export after loop-apply fixture includes all 10 keys with ≥1 row where seeded |
| `TestSeedImportP20RoundTrip` | import → export → fresh import preserves P20 IDs + `deliberation_state.current_phase` |
| `TestSeedExportRoundTrip` | P17 keeper — plan tree + causal entities |
| `TestSeedExportOmitsDeniedSurfaces` | P17 keeper — no transitions/work_state/reviews |
| `TestSeedExportWritesExportedAtCommit` | P17 keeper — git HEAD evidence |

```bash
go test ./internal/domain/... -count=1 -run 'TestSeedExportIncludesP20Cognition|TestSeedImportP20RoundTrip'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit'
```

Fixture hint: seed via domain/store upserts or `trace loop apply` (see `cmd/trace/loop_test.go` helpers) then export; assert JSON keys present.

## Exit criteria

- [ ] All 5 named tests PASS
- [ ] `trace seed export -o trace/graph.json` includes P20 keys when DB populated
- [ ] P20 loop keeper still PASS: `go test ./cmd/trace -count=1 -run 'TestLoopApplyDeliberationTransitionEvent'`
- [ ] Compat ceiling **19** unchanged (`evals/compat`)
- [ ] No mig 020+
- [ ] Board row Notes: test output + graph.json SHA

## Minimal todos

- [ ] Extend `SeedDocument` + export + `ListAll*` helpers
- [ ] Extend import + ordering
- [ ] Named tests green
- [ ] Export graph.json + board Notes
