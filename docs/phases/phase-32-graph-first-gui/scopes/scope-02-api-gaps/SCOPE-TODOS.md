# Scope 02 — board map

**S02 is not API-only.** Deliver:

1. Library/OpenAPI story → expect **`NO-GAPS.md`** (no new core `/v1` for inspector map)
2. **Client glue** → `getImpact` in `web/src/api/ops.ts`
3. **Always P32-PORT** → RESEARCH/OPEN-PORT-MULTI **#1** (friendly `EADDRINUSE` + `--addr` examples); #2 auto-port deferred

Serial: **S02-00 → S02-01 → S02-02**. Do not skip this scope when library gaps are empty.

| Board ID | Row | Prompt | Role |
|----------|-----|--------|------|
| 555 | P32-S02-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner (locks handoff) |
| 556 | P32-S02-01 | [01-implement.md](01-implement.md) | Implementer: `NO-GAPS.md` + `getImpact` + **P32-PORT #1** |
| 557 | P32-S02-02 | [02-review.md](02-review.md) | Reviewer (P32-PORT hard checklist) |

## Planner locks (P32-S02-00)

| Lock | Value |
|------|-------|
| Library API | `NO-GAPS.md` |
| Client | `getImpact(taskId, opt?)` → `/v1/impact?task_id=` |
| P32-PORT | #1 min; #2 defer; no `NO-PORT-CHANGE.md` without reason |
| Bind | Keep `127.0.0.1:7432`; fail-on-conflict; loopback policy unchanged |
| Forbidden | Invent `/v1/path`; require changes/regressions; graph UI (S03) |
