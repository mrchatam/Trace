# Scope 00 — board map

Inventory-only scope. Serial: **S00-00 → S00-01**. No S00 review row — **S01-00** validates `GAPS.md` before implement. Do not start S01 until S00-01 writes `GAPS.md`.

| Board ID | Row | Prompt | Role | Status |
|----------|-----|--------|------|--------|
| 540 | P31-S00-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner (lock prompts; no product code) | done (this planner) |
| 541 | P31-S00-01 | [01-inventory.md](01-inventory.md) | Implementer → write `GAPS.md` | done → [GAPS.md](GAPS.md) (must-add×3: G1/G5/G6) |

Primary artifact: `GAPS.md` (written by **P31-S00-01** only — not this planner).

## Planner lock (P31-S00-00, 2026-08-21)

- Live anchors re-verified: `warnIfStrayRootTraceDB` in `openStore` (`open.go` L85); Stat-only `IsRegular` quiet (`open.go` L144–149); join `.trace`+`trace.db`; four unit cases in `stray_trace_db_test.go`; `/trace.db` only in `.gitignore` + `fixtures/x0/.gitignore`.
- `01-inventory.md` thickened: preflight, must-answer×5, live snapshot, candidate hint table, GAPS template + ignore-scaffold audit, exit criteria.
- Phase locks unchanged: canonical `.trace/trace.db`; no silent delete; no GUI; successor Phase 32.
- S01 lightly noted: serve warn is request-scoped `store.Open`, not process startup — do not invent a startup open.
