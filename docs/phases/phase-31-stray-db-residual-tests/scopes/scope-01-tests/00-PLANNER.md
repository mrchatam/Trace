# P31-S01-00 — Scope planner (tests)

## Metadata
- id: P31-S01-00
- todo_ids: [P31-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, test-driven-development]
- mcps: []
- verification: automated
- hooks: []

## Objective

Validate [`../scope-00-inventory/GAPS.md`](../scope-00-inventory/GAPS.md) and lock implement defaults so `01-implement.md` can ship must-add items alone. **No product code in this row.** Thicken `01-implement.md` + `02-review.md` if still thin.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [Phase 31 README](../../README.md)
- [GAPS.md](../scope-00-inventory/GAPS.md) — required input
- Live: `internal/store/stray_trace_db_test.go`, `internal/store/open.go`

## Session start

Follow agent-loop-protocol. If `GAPS.md` missing or empty must-add with no defer reasons, **blocked** — send back to S00 (do not invent gaps).

## Locked defaults

| Item | Value |
|------|-------|
| Input | `GAPS.md` must-add only — **frozen: G1, G5, G6** |
| G1 | Extend `internal/store/stray_trace_db_test.go` — dir-named root stub quiet |
| G5 | **Required** `scripts/repro-stray-trace-db.sh` (create `scripts/` if needed) |
| G6 | Docs-only: `CONTRIBUTING.md` **and** `AGENTS.md` — once-per-`openStore`; no suppress flag |
| G2 | Nice-to-have — **not** required for S01-01 exit |
| G3 / G4 | Deferred — leave alone |
| Canonical path | Unchanged `.trace/trace.db` |
| Silent delete / GUI / path redesign / suppress flag | **Forbidden** |
| Serve note (S00-00) | HTTP `openStore` → `store.Open` per request (`internal/httpapi/server.go`); do not invent a process-startup open |

## Planner gate

- [x] `GAPS.md` present; must-add list frozen into `01-implement.md` task cards (G1, G5, G6)
- [x] `01-implement.md` has exit criteria + hard checks + `go test ./internal/...` green
- [x] `02-review.md` has evidence checklist + spawn `02a`/`02b` on blocker/high
- [x] `SCOPE-TODOS.md` updated
- [x] No Go written in this row

## Exit criteria

- [x] Implementer can run without re-debating stack/paths
- [x] Board Notes cite must-add IDs locked
- [x] Next: **P31-S01-01**

## Todo updates

Status + notes on **P31-S01-00** only.

## Next

`P31-S01-01`
