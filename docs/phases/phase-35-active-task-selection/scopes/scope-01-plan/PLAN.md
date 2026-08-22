# PLAN — Phase 35 active task selection

## Problem restatement (from INVESTIGATION)

Feet-seller has **123** tasks, all **`DONE`**, store-ordered oldest-first (`created_at ASC, id ASC`). Overview and Loop both bind **Step 1** (`33247e2d-aa10-4b25-b194-4b7afb5a6359`) while later work ends at **Loop 112** (`99d8fb92-65ac-462c-82c4-21bcf198c09e`):

- Overview `pickActiveTask` (`web/src/screens/Overview.tsx` ~L17–20): `IN_PROGRESS` → first non-terminal → **`tasks[0]`**.
- Loop (`web/src/screens/Loop.tsx` ~L51–54): no `?task_id=` → **`items[0].id`**.
- Both call `listTasks({ limit: 100 })`. Today `handleListTasks` **ignores** `limit`/`cursor`, so the client still receives all 123 — the dogfood symptom is the `[0]` fallback, not truncation.
- There is **no** library/HTTP durable “current work” preference (only agent `TRACE_TASK_ID`).
- INVESTIGATION red assert (`defaultPick ≠ Step1`) exits **2**.

Phase 35 must share one pick policy across Overview + Loop so all-DONE lists bind the **last** task (Loop 112 under P3a), keep explicit URL/agent overrides winning, and not leave “current work past page 1” silent if/when `limit` is honored.

## Options (ranked) + rejected

### Placement (closed set)

| Rank | Id | Placement | Verdict |
|-----:|----|-----------|---------|
| 1 | **B** | One shared TS helper under `web/src/lib/` (name: `pickActiveTask.ts`); Overview + Loop import it; no duplicate `[0]` fallbacks | **Chosen** — fits one S02 implement row |
| 2 | **A** | Go library pure pick (+ thin HTTP if needed) + GUI adapters call it | Deferred — see Future library promotion |
| — | Reject | Divergent Overview/Loop logic; TRACE_TASK_ID-only fix; reverse list order without shared policy; localStorage-only without shared pick | Violates DESIGN-LOCKS / Law 19 |

**Sizing for B over A:** The bug is already pure client pick over a list the GUI holds. Placement A requires Go helper + OpenAPI/HTTP surface + client wire before adapters can call it — that does not fit a single S02 vertical slice without scope creep. Law 19 still holds: one shared helper, no Overview/Loop fork; A remains the upgrade path when a library current-work API is warranted.

### Semantics (closed set)

URL / explicit `?task_id=` (and documented agent `TRACE_TASK_ID`) **always win** when present.

| Rank | Id | Rule | Verdict |
|-----:|----|------|---------|
| 1 | **P1** | First `IN_PROGRESS` | Keep |
| 2 | **P2** | Else first non-terminal (not `DONE`/`SKIPPED`) | Keep |
| 3 | **P3a** | Else **last** list item (`items[n-1]` on oldest-first list) | **Chosen** all-DONE fallback — fixes feet-seller |
| — | P3b | Else last “meaningful” (title/body/plan-link heuristic) | Rejected as default — no cheap predicate needed for gate |
| — | P3c | Else `planner.GetCurrentScope` / plan current | Rejected as default — does not apply when all DONE + no plan |
| — | P3d | Optional localStorage last-focused id | Enhancement only; **not** sole fix; out of S02 must-ship |

**Rejected:** keeping `tasks[0]` / `items[0]` as the all-DONE fallback.

## Chosen policy (normative)

**Auto-pick from a task list (oldest-first):**

1. **P1** — return first task with `work_state === 'IN_PROGRESS'`.
2. **P2** — else return first task whose `work_state` is not `DONE` and not `SKIPPED`.
3. **P3a** — else return `tasks[tasks.length - 1]` (last by store order = newest among terminals on feet-seller).
4. Empty list → `null` / no auto-set.

**Overrides (always win when present):**

- Loop / GUI: existing `?task_id=` must not be overwritten by auto-pick.
- Agents: `TRACE_TASK_ID` remains the process bind; GUI default must not contradict “current work” once P3a ships (document alignment; do not make env the sole fix).

**List completeness for P3a:** The helper picks over whatever list it is given. Callers that auto-pick **must** obtain a list that includes the P3a candidate (see limit honesty below). Do not assume newest-first store order.

## Law 19 placement (library vs adapter)

**Chosen: B** — shared adapter helper, not a second business-logic fork.

| Surface | Role |
|---------|------|
| `web/src/lib/pickActiveTask.ts` | Canonical pick policy (pure function) |
| `Overview.tsx` | Delete local `pickActiveTask`; import helper |
| `Loop.tsx` | When `!taskId`, set `task_id` from `pickActiveTask(items)?.id` — **not** `items[0]` |

HTTP handlers and OpenAPI stay adapters only; S02 does **not** add a Go current-work API.

### Future library promotion (A)

When a durable library/HTTP “current work” exists (or pagination forces server-side pick), promote the same P1→P2→P3a semantics into Go and have GUI call that API. Phase 35 ship does not require that promotion.

## API / GUI / docs touch list

**In scope for S02 (finalize ⊆ candidate set):**

| Layer | Path | Action |
|-------|------|--------|
| GUI must | `web/src/screens/Overview.tsx` | Use shared helper for active pick |
| GUI must | `web/src/screens/Loop.tsx` | Use shared helper for default `?task_id=` |
| GUI shared | `web/src/lib/pickActiveTask.ts` (**new**) | Export `pickActiveTask(tasks)` |
| GUI shared | `web/src/lib/pickActiveTask.test.ts` (**new**, or project test sibling) | Unit tests for acceptance cases 1–4 + truncated-list honesty |
| GUI fetch | `Overview.tsx` / `Loop.tsx` (and optionally thin helper in `web/src/api/ops.ts` only if needed for paging) | Auto-pick path: obtain full list for pick (omit/raise `limit`, or page with `cursor` until exhausted). Display list may still use a page; pick must not silently use a truncated first page as P3a input |

**Explicitly not in S02 touch list:**

| Layer | Path | Reason |
|-------|------|--------|
| GUI optional | `web/src/lib/overviewCompose.ts` / Graph | Related smell only; not Phase-35 gate seed |
| HTTP | `internal/httpapi/handlers_tasks.go`, `api/openapi.yaml` | Full `limit`/`cursor` compliance **deferred** (later phase) |
| Library | new `internal/…` pick helper | Placement A deferred |
| Docs optional | `cmd/trace/AGENTS.md` | Optional one-paragraph TRACE_TASK_ID ↔ GUI default alignment; not blocking if S02 notes defer to VERIFY docs |

## Acceptance tests

Use **synthetic** in-repo fixtures for CI. Narrative UUIDs (feet-seller, read-only — never mutate):

- Step1 = `33247e2d-aa10-4b25-b194-4b7afb5a6359`
- Loop112 = `99d8fb92-65ac-462c-82c4-21bcf198c09e`

| # | Case | Expect |
|---|------|--------|
| 1 | **All-DONE + later tasks** — oldest-first `[Step1, …, Loop112]` all `DONE` | `pickActiveTask` → Loop112; **≠** Step1 |
| 2 | **IN_PROGRESS wins** — mixed list | First `IN_PROGRESS` |
| 3 | **Non-terminal over DONE** — no `IN_PROGRESS`; PENDING (etc.) before DONE rows | First non-terminal; not P3a |
| 4 | **Explicit `task_id`** — Loop/URL already has `task_id` | Auto-pick does not overwrite |
| 5 | **>100 / limit honesty** — synthetic list length >100; caller would otherwise pass only first 100 (index 122 missing) | Either (a) fetch strategy supplies full/last page so P3a = true last task, **or** (b) documented fail-closed behavior that does **not** silently bind index 0 as “current” when the list is known truncated. Prefer (a). |
| 6 | **Red→green seed** — INVESTIGATION assert (`defaultPick ≠ Step1` on all-DONE oldest-first list) | Automated unit (or thin integration) test; green under shared helper |

**S02 does not implement** HTTP handler `limit`/`cursor` compliance. Honesty is satisfied by client fetch-for-pick + tests above. Optional later-phase note: when handler gains pagination, re-run case 5 against live HTTP.

Suggested unit cmds (S02 locks exact): `cd web && npm test -- pickActiveTask` (or repo-equivalent vitest/jest target).

## Agent TRACE_TASK_ID alignment

- Agents continue to bind via `TRACE_TASK_ID` (documented in `cmd/trace/AGENTS.md` / install prompts).
- GUI default after P3a should match “last progressed work” on all-DONE boards, reducing contradiction with agents who left off at the end of the list.
- **Not** a substitute for the GUI fix: env-only remediation remains rejected.
- Optional S02 docs touch: one short note that GUI auto-pick and `TRACE_TASK_ID` are complementary (URL/env override vs list policy).

## Out of scope / non-goals

- Weakening `plan_missing` / PLAN-phase gate behavior for true PLAN work.
- TRACE_TASK_ID as the **only** fix.
- Deleting or rewriting feet-seller / dogfood data.
- Hosted SaaS / multi-tenant current-work.
- Divergent Overview vs Loop fallbacks.
- Graph `prioritizeTaskSeeds` smell (follow-on; not gate).
- Full OpenAPI pagination implementation in S02.
- Placement A Go/HTTP current-work API in Phase 35 ship.
- P3b/P3c defaults; P3d as sole or must-ship fix.

## Handoff to S02

### Implement files (must)

1. **New** `web/src/lib/pickActiveTask.ts` — pure `pickActiveTask(tasks: TaskRow[]): TaskRow | null` implementing P1→P2→P3a.
2. **New** unit test sibling covering acceptance #1–4 and #5–6 (synthetic fixtures; no dogfood mutation).
3. `web/src/screens/Overview.tsx` — remove local pick; import helper; ensure fetch-for-pick completeness.
4. `web/src/screens/Loop.tsx` — replace `items[0]` default with helper; preserve `?task_id=` override; ensure fetch-for-pick completeness.

### Optional

- Thin list-all / page-until-exhausted helper near `web/src/api/ops.ts` if both screens need the same fetch loop.
- One paragraph in `cmd/trace/AGENTS.md` on TRACE_TASK_ID ↔ GUI default.

### Do not touch in S02 (unless spawn)

- `overviewCompose.ts` / Graph
- `handlers_tasks.go` / `openapi.yaml` for pagination
- Go library pick / new current-work route

### Test names / cmds sketch (S02 locks)

- `pickActiveTask` unit: `allDone_picksLast`, `inProgress_wins`, `nonTerminal_beforeDone`, `empty_returnsNull`, `truncatedList_honesty` (or equivalent).
- Red→green: encode INVESTIGATION seed (Step1 first, Loop112 last, all DONE → pick ≠ Step1).
- Cmd sketch: `cd web && npm test -- pickActiveTask` (adjust to repo runner).

### What VERIFY (S03) re-checks on live feet-seller

- Read-only fixture `/home/ali/Desktop/feet seller telegram app`.
- With no `?task_id=`, Overview active / Loop default bind **≠** Step1 (`33247e2d-…`); prefer Loop112 (`99d8fb92-…`) under all-DONE.
- Do not require mutating dogfood; do not weaken `plan_missing` gates.
- Confirm client pick path still correct if HTTP still returns full 123 ignoring limit (today) **and** document residual risk if pagination lands later without client page-through.

### Blockers for P35-S02-00

None from S01. S02-00 should thicken implement prompts from this PLAN (placement **B**, policy **P1→P2→P3a**, no HTTP pagination in-slice, fetch-for-pick honesty required).
