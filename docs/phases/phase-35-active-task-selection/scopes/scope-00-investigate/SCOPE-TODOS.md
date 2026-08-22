# Scope 00 — board map

**S00 investigate** — repro feet-seller + root-cause cites only. Serial: **P35-S00-00 → P35-S00-01 → P35-S00-02**. Primary artifact: `INVESTIGATION.md` (written in **S00-01**, reviewed in **S00-02**). Do **not** start S01 until S00-02 PASS. Do **not** write product code. Planner (**S00-00**) does **not** author `INVESTIGATION.md`.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 611 | P35-S00-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 + this file |
| 612 | P35-S00-01 | [01-investigate.md](01-investigate.md) | Implementer | Author `INVESTIGATION.md` + red-capable repro notes |
| 613 | P35-S00-02 | [02-review.md](02-review.md) | Reviewer | Checklist vs DESIGN-LOCKS + INTAKE + live cites |

## Planner-locked live facts (for 01) — verified 2026-08-21

- Fixture: `/home/ali/Desktop/feet seller telegram app` — **read-only** against `.trace/` (no delete/reset).
- CLI: `trace -C "<feet-seller>" tasks` → **123** tasks, **all DONE**.
  - First: `33247e2d-aa10-4b25-b194-4b7afb5a6359` — Step 1: Market and feature research
  - Last: `99d8fb92-65ac-462c-82c4-21bcf198c09e` — Loop 112: Entitlements polish + RESUME STOP
- Store order: `internal/store/helpers.go` `ListTasks` — `ORDER BY created_at ASC, id ASC` (oldest first).
- Overview: `web/src/screens/Overview.tsx` L17–20 `pickActiveTask` → IN_PROGRESS else non-terminal else **`tasks[0]`**; L38 `listTasks({ limit: 100 })` then gate/status for pick.
- Loop: `web/src/screens/Loop.tsx` L51–54 — no `?task_id=` ⇒ `setParams` to **`res.items[0].id`**.
- Client: `web/src/api/ops.ts` L41–45 forwards `limit`/`cursor`.
- HTTP: `internal/httpapi/handlers_tasks.go` L18–48 `handleListTasks` — OpenAPI documents `limit`/`cursor`, but handler **does not read them** (returns full filtered list). **S00-01 must prove with live HTTP** `items.length` for `?limit=100` (expect 123 today if still unpaginated).
- Related (not primary gate seed): `web/src/lib/overviewCompose.ts` `prioritizeTaskSeeds` — DONE used when no active; confirm Explore vs Overview surfaces.
- No durable “current work” entity found at plan time; agents use `TRACE_TASK_ID` (`cmd/trace/AGENTS.md`). AppChrome localStorage = theme/token only.
- Secondary: gate `plan_missing` on DONE — honesty note only; **do not** propose weakening PLAN-phase gates.

## Investigation rejects (01 must document)

1. “Agents should just set TRACE_TASK_ID” as the **only** fix — GUI still seeds wrong without URL/env.
2. Blaming Phase 34 packaging/ports for Step-1 default.
3. Deleting or rewriting feet-seller task history to “fix” the demo.
4. Treating `plan_missing` as the primary bug that explains Step 1 vs 123.

## Out of this scope

- Writing `PLAN.md` (S01), product fixes (S02), VERIFY (S03).
- Authoring `INVESTIGATION.md` in the planner row.
