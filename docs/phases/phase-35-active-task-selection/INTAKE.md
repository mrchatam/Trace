# INTAKE — wrong active task (feet-seller dogfood)

**Human 2026-08-21.** Project: `/home/ali/Desktop/feet seller telegram app`

## Symptom

GUI / gate shows **task 1** (“Step 1: Market and feature research”) while the human is at **task 123** (Loop 112 / last work). Gate strip: `plan_missing` → PLAN. Same class of confusion seen when agents bind the wrong task.

## Light facts (2026-08-21)

| Fact | Value |
|------|--------|
| Task count | **123**, all **`DONE`** |
| List index 0 | `33247e2d-…` Step 1 (DONE) |
| List index 122 | `99d8fb92-…` Loop 112 Entitlements polish (DONE) |
| Edit gate first | `plan_missing` / PLAN / allowed=false |
| Edit gate last | `plan_missing` / PLAN / allowed=false |

## Likely product causes (must re-verify in S00)

1. **`Overview.pickActiveTask`** (`web/src/screens/Overview.tsx`): first `IN_PROGRESS`, else first non-terminal, else **`tasks[0]`**. With all DONE → always **Step 1**.
2. **`Loop.tsx`**: if no `?task_id=`, auto-selects **`res.items[0]`** — again Step 1.
3. **`listTasks({ limit: 100 })`** on Overview/Loop: client requests a cap; projects with **>100 tasks** risk never seeing task 123 if/when HTTP honors `limit`. **P35-00 note:** `handleListTasks` currently appears to ignore `limit` (returns full list) — S00 must re-verify end-to-end; selection bug still reproduces via `tasks[0]` even with a full list.
4. **No durable “current work”**: no last-touched / plan-current-scope / localStorage preference drives the gate seed.
5. Gate on **DONE** still reporting `plan_missing` may be correct library policy or a secondary honesty bug — investigate separately.

## Not the primary issue

- Phase 34 packaging/ports
- Agents “should create a plan” alone does not explain showing **task 1** vs **task 123**

## Desired outcome

Surfaces that show “the” task + gate must bind to **current work** (or an explicit user selection), not silently the oldest list row. Dogfood: feet-seller with 123 DONE tasks must not default gate to Step 1.
