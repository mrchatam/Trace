# Scope 01 — board map

Tests + review. Serial: **S01-00 → S01-01 → S01-02**. Input: [`../scope-00-inventory/GAPS.md`](../scope-00-inventory/GAPS.md). Do not start S02 until S01-02 PASS (or spawned 02a/02b close).

| Board ID | Row | Prompt | Role |
|----------|-----|--------|------|
| 542 | P31-S01-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner (lock must-add; no product code) |
| 543 | P31-S01-01 | [01-implement.md](01-implement.md) | Implementer → G1 + G5 + G6 |
| 544 | P31-S01-02 | [02-review.md](02-review.md) | Reviewer → PASS/FAIL; spawn 02a/02b on blocker/high |

## Locked must-add (P31-S01-00, 2026-08-21)

| ID | What | Home |
|----|------|------|
| G1 | Dir-named root `trace.db` quiet unit | `internal/store/stray_trace_db_test.go` |
| G5 | Durable dogfood repro | `scripts/repro-stray-trace-db.sh` |
| G6 | Multi-open warn intentional | `CONTRIBUTING.md` + `AGENTS.md` |

- Nice-to-have (not required for S01-01 exit): **G2** CLI stderr harness
- Deferred: **G3** serve startup, **G4** extra ignore scaffolds

## Hard bans

- Store-path redesign (canonical remains `<root>/.trace/trace.db`)
- Silent / flagged delete of root stub
- GUI / Phase 32 explorer work
- Warn suppress flag
- Inventing a `trace serve` process-startup open (HTTP uses per-request `store.Open`)

## Serve note (from S00)

`trace serve` / HTTP opens via `store.Open(s.root)` **per handler request** (`internal/httpapi/server.go`) — not a dedicated process-startup open.
