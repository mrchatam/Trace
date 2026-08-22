# INVESTIGATION — Phase 35 S00

## Verdict (1 paragraph)

Feet-seller has **123** tasks, all **`DONE`**, ordered oldest-first. Overview and Loop both bind **Step 1** (`33247e2d-…`) because `pickActiveTask` falls through to `tasks[0]` and Loop auto-sets `?task_id=` to `res.items[0].id` when absent. HTTP `GET /v1/tasks?limit=100` returns **123** items today (`handleListTasks` ignores `limit`/`cursor`); selection still fails with a full list. There is **no** library/HTTP “current work” / last-touched / GUI-persisted task preference — only agent `TRACE_TASK_ID`. A red assert (`defaultPick ≠ Step1`) fails with exit 2 while this bug exists.

## Repro steps + evidence

**Fixture (read-only):** `/home/ali/Desktop/feet seller telegram app`

### CLI baseline (2026-08-21)

```bash
trace -C "/home/ali/Desktop/feet seller telegram app" tasks
```

| Fact | Value |
|------|--------|
| Count | **123** |
| States | all **`DONE`** (0 non-DONE) |
| First | `33247e2d-aa10-4b25-b194-4b7afb5a6359` — Step 1: Market and feature research |
| Last | `99d8fb92-65ac-462c-82c4-21bcf198c09e` — Loop 112: Entitlements polish + RESUME STOP |

### HTTP (live)

```bash
trace serve --root "/home/ali/Desktop/feet seller telegram app" --addr "127.0.0.1:$PORT"
curl -s "http://127.0.0.1:$PORT/v1/tasks?limit=100"
curl -s "http://127.0.0.1:$PORT/v1/tasks"
```

| Request | `items.length` | First id | Last id |
|---------|----------------|----------|---------|
| `?limit=100` | **123** | Step 1 | Loop 112 |
| no limit | **123** | Step 1 | Loop 112 |

`limit_honored=False`; lengths identical.

### Simulated Overview / Loop bind (no GUI required)

Same list as client would receive (`listTasks({ limit: 100 })` → full 123):

- **Overview `pickActiveTask`:** → `33247e2d-…` Step 1
- **Loop `items[0]`:** → `33247e2d-…` Step 1

Gate strip on that bind (secondary honesty): both Step 1 and Loop 112 return `allowed=false`, `recommended_phase=PLAN`, `reason_code=plan_missing` for `for=edit`.

Optional GUI: open `/overview` and `/loop` without `?task_id=` — expect the same Step 1 bind (logic-proven above; not required for this row).

## Root causes (with cites)

### Confirmed primary causes

1. **`Overview.pickActiveTask`** — `web/src/screens/Overview.tsx` L17–20: first `IN_PROGRESS`, else first non-`DONE`/`SKIPPED`, else **`tasks[0]`**. With all DONE → always oldest row (Step 1). Used at L44–50 to load status/gate for that pick. Client fetch: L38 `listTasks({ limit: 100 })`.

2. **`Loop.tsx` default** — `web/src/screens/Loop.tsx` L51–54: `listTasks({ limit: 100 })`; if no `?task_id=` and `res.items?.[0]`, **`setParams({ task_id: res.items[0].id })`**. Oldest-first list → Step 1.

3. **Store / API order** — `internal/store/helpers.go` `ListTasks` L39–43: `ORDER BY created_at ASC, id ASC` (oldest first). HTTP list returns that order unchanged.

### INTAKE “likely causes” 1–5

| # | Claim | Verdict |
|---|--------|---------|
| 1 | Overview `pickActiveTask` → `tasks[0]` when all DONE | **Confirmed** (L17–20 + live sim) |
| 2 | Loop auto-selects `items[0]` | **Confirmed** (L51–54 + live sim) |
| 3 | `listTasks({ limit: 100 })` may hide task 123 if HTTP honors limit | **Partially:** client requests limit (**ops.ts** L41–45); handler **ignores** limit today so 123 is present — selection bug still via `[0]`. **If** pagination were honored, Loop 112 is index **122** → **would be hidden** from a 100-item page |
| 4 | No durable “current work” preference | **Confirmed** — see library search; AppChrome localStorage = theme/token only |
| 5 | Gate `plan_missing` on DONE | **Confirmed secondary** on both Step 1 and Loop 112; does **not** explain Step 1 vs 123 |

### Related (not Overview/Loop gate seed)

- `web/src/lib/overviewCompose.ts` `prioritizeTaskSeeds` (L50–73): used by **Graph** seed composition (`Graph.tsx`), not by Overview `pickActiveTask`. With only terminal tasks it still starts from list order → Step 1 among seeds. Separate surface; same “oldest DONE first” smell.

## limit / pagination honesty

- **Live:** `GET /v1/tasks?limit=100` → **123** items (not 100).
- **Handler:** `internal/httpapi/handlers_tasks.go` L18–48 `handleListTasks` reads only `goal_id` / `work_state`; **does not read `limit`/`cursor`**; returns full filtered `items`.
- **Contract:** `api/openapi.yaml` `/v1/tasks` documents optional `limit`/`cursor` (default 100) — **contract vs handler gap**.
- **Client:** `web/src/api/ops.ts` L41–45 forwards `limit`/`cursor` in the query string.
- **Hypo:** if limit were honored at 100, Loop 112 (index 122) would **not** appear in the first page — a future pagination honesty risk for “current work” at the end of a long DONE list. Today’s dogfood symptom does **not** require truncation; `[0]` is enough.

## Library / API “current work” search

**Result: none** for durable GUI/agent “current work” task selection.

Searched `internal/` + HTTP routes + web chrome:

| Candidate | Finding |
|-----------|---------|
| HTTP `/v1/tasks` | List/get only; no current-work field |
| HTTP plans | `GET /v1/plans` list — not a focused-task API |
| `planner.GetCurrentScope` | Plan **scope** current ref for a goal (`plan.current_set`); **not** exposed as GUI “active task”, not used by Overview/Loop |
| `TRACE_TASK_ID` | Agent bind documented in `cmd/trace/AGENTS.md` / install prompts — process env, not durable store |
| `AppChrome.tsx` | localStorage: theme + API token only — **no task preference key** |

Law 19 implication for S01/S02: prefer a library/API “current work” if introduced; else **one** shared GUI pick helper (no Overview/Loop fork).

## Red-capable feedback loop (for S02)

**Already run (RED, exit 2)** — pure HTTP + pick logic (no product edit):

```bash
FIXTURE="/home/ali/Desktop/feet seller telegram app"
STEP1="33247e2d-aa10-4b25-b194-4b7afb5a6359"
LOOP112="99d8fb92-65ac-462c-82c4-21bcf198c09e"
PORT=…  # free loopback
trace serve --root "$FIXTURE" --addr "127.0.0.1:$PORT" &
# then:
curl -s "http://127.0.0.1:$PORT/v1/tasks?limit=100" | python3 -c '
import json,sys
STEP1="33247e2d-aa10-4b25-b194-4b7afb5a6359"
LOOP112="99d8fb92-65ac-462c-82c4-21bcf198c09e"
items=json.load(sys.stdin)["items"]
def pick(tasks):
  for t in tasks:
    if t["work_state"]=="IN_PROGRESS": return t
  for t in tasks:
    if t["work_state"] not in ("DONE","SKIPPED"): return t
  return tasks[0] if tasks else None
p=pick(items)
assert items[0]["id"]==STEP1 and items[-1]["id"]==LOOP112
# DESIGN-LOCKS must-test: default bound ≠ Step 1 when later DONE exist
sys.exit(0 if p["id"] != STEP1 else 2)
'
```

**Observed:** `assert_exit=2` — `bound=33247e2d-…` (Step 1); substrate first=Step1, last=Loop112.

S02 should turn this into a unit/integration test of the **shared** pick policy (green when default ≠ Step 1 under all-DONE + later tasks). Optional: assert `limit` honesty separately when/if handler gains pagination.

## Secondary notes (plan_missing, etc.)

- Edit gate on **DONE** tasks returns `plan_missing` / PLAN / `allowed=false` for both first and last tasks — confusing UX but **orthogonal** to wrong-task selection.
- Do not weaken PLAN-phase gates for true PLAN work (DESIGN-LOCKS out of scope).

## Rejected alternatives

From SCOPE-TODOS (investigation rejects):

1. **“Agents should just set TRACE_TASK_ID” as the only fix** — GUI Overview/Loop still seed wrong without URL/env; rejected as sole remedy.
2. **Blaming Phase 34 packaging/ports** — selection is client pick + oldest-first list; not embed/port related.
3. **Deleting or rewriting feet-seller task history** — forbidden; hides the bug.
4. **Treating `plan_missing` as the primary bug** — explains gate strip tone, not Step 1 vs Loop 112.

## Handoff to S01

- **Artifact for review:** this file (`P35-S00-02`).
- **Policy problem to solve in S01 `PLAN.md`:** define “current work” when all tasks are DONE (e.g. last DONE / last meaningful / plan current scope / explicit pick) and place it in **library/API** if possible (Law 19), else one shared GUI helper used by Overview + Loop.
- **Must not:** keep divergent `tasks[0]` / `items[0]` fallbacks; ignore future `limit` honesty for >100 tasks.
- **S02 implement:** change pick + tests so red assert above goes green; pagination honesty if in scope of implement locks.
- **Do not start S01 until S00-02 PASS.**
