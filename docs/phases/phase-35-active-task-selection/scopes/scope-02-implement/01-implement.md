# P35-S02-01 — Implement active task pick

## Metadata
- id: P35-S02-01
- todo_ids: [P35-S02-01]
- role: implementer
- skills: [test-driven-development, incremental-implementation, frontend-ui-engineering]
- verification: automated
- mcps: []

## Objective

Ship **placement B** from [`PLAN.md`](../scope-01-plan/PLAN.md): one shared pure helper `pickActiveTask` (P1→P2→P3a), wire **Overview** + **Loop**, remove inline `tasks[0]` / `items[0]` auto-pick fallbacks, and green automated unit tests (incl. all-DONE ≠ Step1 + limit honesty). **No** HTTP pagination / OpenAPI / Go library current-work in this row.

## References

- [`PLAN.md`](../scope-01-plan/PLAN.md) — **normative**
- [`00-PLANNER.md`](00-PLANNER.md) — locked files/cmds (this scope)
- [`INVESTIGATION.md`](../scope-00-investigate/INVESTIGATION.md)
- [`DESIGN-LOCKS.md`](../../DESIGN-LOCKS.md)
- `docs/rules/agent-loop-protocol.md` (implementer loop)
- Law 19: adapters thin; no Overview/Loop fork

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Do **not** re-debate placement A vs B or P3a vs P3b/c. If a lock conflicts with live paths, note on the board and continue with the closest lawful fix.

## Locked defaults (do not re-debate)

| Item | Locked value |
|------|----------------|
| Placement | **B** — `web/src/lib/pickActiveTask.ts` (pure). Placement **A** (Go) deferred |
| Semantics | **P1** first `IN_PROGRESS` → **P2** first non-`DONE`/`SKIPPED` → **P3a** `tasks[tasks.length - 1]` → empty → `null` |
| Signature | `export function pickActiveTask(tasks: TaskRow[]): TaskRow \| null` (`TaskRow` from `../api/ops`) |
| Must wire | `web/src/screens/Overview.tsx` — delete local `pickActiveTask` (L17–20); import shared helper |
| Must wire | `web/src/screens/Loop.tsx` — when `!taskId`, `setParams({ task_id: pickActiveTask(items)?.id })` — **not** `items[0]` |
| Overrides | Explicit `?task_id=` **wins** (Loop must not overwrite when present). Agent `TRACE_TASK_ID` stays process bind; do not make env the sole fix |
| Fetch-for-pick | Auto-pick callers **must not** silently P3a on a truncated first page. Prefer thin **`listTasksForPick`** in `web/src/api/ops.ts`: page `listTasks({ limit: 100, cursor })` while `next_cursor` is set; concatenate. Today handler ignores limit/cursor → one full page (OK). Display lists may keep `limit: 100`; **pick path** uses `listTasksForPick` |
| Fail-closed (only if paging impossible) | Do **not** bind index 0 as “current” when the list is known truncated; prefer (a) full fetch over (b) |
| Tests | **New** `web/src/lib/pickActiveTask.test.ts` — same pattern as `overviewCompose.test.ts` (**node:test** + strip-types). **No** vitest/`npm test` (web has neither) |
| Test cmd | `cd web && node --experimental-strip-types --test src/lib/pickActiveTask.test.ts` |
| Narrative UUIDs (fixtures only) | Step1=`33247e2d-aa10-4b25-b194-4b7afb5a6359`; Loop112=`99d8fb92-65ac-462c-82c4-21bcf198c09e` — synthetic lists in unit tests; **never** mutate feet-seller |
| Optional | One short `cmd/trace/AGENTS.md` note TRACE_TASK_ID ↔ GUI default — **not** blocking |
| Do **not** touch | `handlers_tasks.go`, `openapi.yaml` pagination, Go pick helper, `overviewCompose.ts` / Graph, feet-seller data, `plan_missing` gate behavior |

## Preflight / Plan

1. Confirm Overview still has local pick ~L17–20 and Loop still uses `items[0]` ~L51–54.
2. TDD: write failing `pickActiveTask.test.ts` first (red), then helper, then wire screens + fetch helper.
3. Run locked test cmd; cite exit 0 in board Notes.

## Role work

Implementer loop on Minimal todos below. Product code allowed. Status + Notes on own board row only.

## Todo updates

Board rights: **status + notes only** on `P35-S02-01`. Do not rewrite this prompt or spawn rows (reviewer owns that).

## Exit criteria

- [ ] `pickActiveTask.ts` implements P1→P2→P3a; no duplicate pick policy in Overview/Loop
- [ ] Overview + Loop use shared helper; no `tasks[0]` / `items[0]` auto-pick fallbacks remain on those paths
- [ ] Explicit `?task_id=` preserved on Loop
- [ ] Auto-pick path uses fetch-for-pick completeness (`listTasksForPick` or equivalent page-until-exhausted / omit-limit that does not silently truncate for P3a)
- [ ] Unit tests green via locked cmd; cover acceptance #1–6 from PLAN (synthetic)
- [ ] All-DONE oldest-first fixture: pick **≠** Step1; prefer last id (Loop112 narrative in synthetic list)
- [ ] No feet-seller mutation; no HTTP/OpenAPI pagination work; no Go library pick
- [ ] Board Notes: test cmd + pass evidence; residual risk if HTTP later honors limit without client page-through (should be mitigated by `listTasksForPick`)
- [ ] Next **P35-S02-02**

## Minimal todos

- [ ] **RED** — Add `web/src/lib/pickActiveTask.test.ts` with cases:
  - `allDone_picksLast` — oldest-first all DONE → last id; **≠** Step1 (`33247e2d-…`); encode Loop112 as last when using narrative ids
  - `inProgress_wins`
  - `nonTerminal_beforeDone`
  - `empty_returnsNull`
  - `truncatedList_honesty` — length >100: either pick on full list after fetch strategy, or prove helper+caller contract does not silently treat truncated page’s `[0]` as current (prefer testing full synthetic >100 → last)
  - Red→green seed: INVESTIGATION assert shape (`defaultPick ≠ Step1` on all-DONE oldest-first)
- [ ] **GREEN helper** — Add `web/src/lib/pickActiveTask.ts` exporting `pickActiveTask`
- [ ] **Fetch-for-pick** — Add `listTasksForPick` (or equivalent) in `web/src/api/ops.ts`; Overview + Loop auto-pick paths call it instead of bare `listTasks({ limit: 100 })` for the list fed to `pickActiveTask`
- [ ] **Wire Overview** — Remove local function; import helper; set active from pick on fetched-for-pick list
- [ ] **Wire Loop** — If `!taskId`, set `task_id` from `pickActiveTask(items)?.id` only; never overwrite existing `?task_id=`
- [ ] **GREEN** — Run `cd web && node --experimental-strip-types --test src/lib/pickActiveTask.test.ts` (exit 0)
- [ ] **Sanity** — `cd web && npm run lint` (or oxlint on touched files) if cheap; fix regressions you introduce
- [ ] **Optional** — AGENTS.md one-liner TRACE_TASK_ID ↔ GUI; skip if timeboxed
- [ ] Board Notes + mark `P35-S02-01` **done** (or failed/blocked with reason)

## Next

`P35-S02-02`
