# Scope 00 — board map

Investigation-only scope. Serial: **S00-00 → S00-01 → S00-02**. Do not start S01 until S00-02 PASS (or spawned 02a/02b close).

| Board ID | Row | Prompt | Role |
|----------|-----|--------|------|
| 528 | P30-S00-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner (lock prompts; no product code) |
| 529 | P30-S00-01 | [01-investigate.md](01-investigate.md) | Implementer → write `INVESTIGATION.md` |
| 530 | P30-S00-02 | [02-review.md](02-review.md) | Reviewer → PASS/FAIL; spawn 02a/02b only on blocker/high |

Primary artifact: `INVESTIGATION.md` (written by **P30-S00-01**, reviewed by **P30-S00-02**).

Planner lock (P30-S00-00, 2026-08-21): live anchors verified (`open.go` `.trace`+`trace.db`; MCP `OpenExisting`; HTTP `store.Open(s.root)`; help documents `.trace/trace.db`). Canonical path assumption until overturned: `<root>/.trace/trace.db`.

**S00 closed** (`P30-S00-02` PASS, 2026-08-21): verdict **agent hygiene**; path lock **confirmed** (not overturned). Next board row: **P30-S01-00**.
