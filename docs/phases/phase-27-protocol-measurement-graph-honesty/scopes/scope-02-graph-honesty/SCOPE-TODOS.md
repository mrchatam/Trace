# Scope 02 todos — graph honesty (INT-07)

Board: [`docs/TODO/phase-27.md`](../../../../TODO/phase-27.md)

| Order | ID | Role | Prompt | Status | Notes |
|------:|----|------|--------|--------|-------|
| 459 | P27-S02-00 | planner | [00-PLANNER.md](00-PLANNER.md) | done | Locked defaults in [01-implement.md](01-implement.md); review checklist in [02-review.md](02-review.md) |
| 460 | P27-S02-01 | implementer | [01-implement.md](01-implement.md) | pending | S02-T01..T06; product Go only |
| 461 | P27-S02-02 | reviewer | [02-review.md](02-review.md) | pending | Verify locks; APPROVE → P27-S03-00 |

## Locked defaults (P27-S02-00)

| Decision | Lock |
|----------|------|
| Thin-graph rule | Min count `discoveries≥1 OR decisions≥1` + orphan link validation when non-empty |
| Export vs done gate | Honesty in `seed.go`; `GateForExport` ≡ `evaluateDone` unchanged |
| BLOCKING rule | Store BLOCKING discovery must have `discovery_mentions_task` link |
| Done-task gate skip | Keep skip for done/skipped/stale |
| Strict vs enforce | Violations always on `--strict`; fail closed only with `--enforce` |
| Harness | No `score.sh` change on S02-01; S03 VERIFY enables `--enforce` |
| Fixtures | P26 `p26-export-snippet.json` thin-graph baseline |

## Task seeds (AUDIT → implement prompt)

| Task | Files | Lock |
|------|-------|------|
| S02-T01 | `cmd/trace/seed.go`, `internal/domain/seed_export_honesty.go` | `collectExportGraphHonestyViolations` |
| S02-T02 | — (no gate.go split) | Document honesty only |
| S02-T03 | `cmd/trace/seed.go` | Keep gate skip L121–124 |
| S02-T04 | `internal/domain/*_test.go` | Unit tests; no eval-rules body |
| S02-T05 | `cmd/trace/enforce_test.go`, `testdata/p26-export-snippet.json` | Thin-graph enforce block |
| S02-T06 | `internal/loop/gate_test.go` | Export≡done test unchanged |
