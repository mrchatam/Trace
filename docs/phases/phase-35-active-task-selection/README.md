# Phase 35 — Active task selection (test + fix)

Human-promoted **2026-08-21**. GUI/loop gate surfaces bind **task 1** (oldest list row) while dogfood work is at **task ~123**.

## Design SoT

| Doc | Role |
|-----|------|
| [`INTAKE.md`](INTAKE.md) | Symptom + light facts (feet-seller) |
| [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md) | Must-fix / must-test / out of scope |
| [`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md) | Phase planner row `P35-00` |
| [`DR-HANDOFF.md`](DR-HANDOFF.md) | Open until `P35-S03-02` |

Board: [`docs/TODO/phase-35.md`](../../TODO/phase-35.md).

## Dogfood fixture

`/home/ali/Desktop/feet seller telegram app` — **do not delete** `.trace/` data.

Live baseline (P35-00 spot-check 2026-08-21):

| Fact | Value |
|------|--------|
| Task count | **123**, all **`DONE`** |
| List index 0 | `33247e2d-…` Step 1: Market and feature research |
| List index 122 | `99d8fb92-…` Loop 112: Entitlements polish + RESUME STOP |
| Store order | `ORDER BY created_at ASC, id ASC` → oldest first |

## Live code pointers (planner baseline — S00 must re-verify)

| Surface | Path | Behavior today |
|---------|------|----------------|
| Overview pick | `web/src/screens/Overview.tsx` `pickActiveTask` | First `IN_PROGRESS`, else first non-terminal, else **`tasks[0]`** → all-DONE ⇒ Step 1 |
| Loop default | `web/src/screens/Loop.tsx` | No `?task_id=` ⇒ **`res.items[0]`** |
| Client list | Overview/Loop `listTasks({ limit: 100 })` | Client requests limit |
| HTTP list | `internal/httpapi/handlers_tasks.go` | **Currently ignores `limit`/`cursor`** (returns full list) — S00 must confirm; truncation is still a honesty risk if/when pagination lands |
| Store | `internal/store/helpers.go` `ListTasks` | Oldest-first |
| Agent bind | `cmd/trace/AGENTS.md` | Expects `TRACE_TASK_ID` — no durable GUI “current work” |

Explore (`/`) is graph-first; gate confusion is primarily **Overview** + **Loop** (+ any shared pick helpers). `plan_missing` on DONE is **secondary** — do not weaken PLAN-phase gates (DESIGN-LOCKS).

## Scope sequence (board SoT)

```
S00 investigate (repro + root cause)
 → S01 plan selection policy + API/GUI
 → S02 implement + tests + review
 → S03 VERIFY (feet-seller live) + DR-HANDOFF
```

| Scope | Board rows | Primary artifact |
|-------|------------|------------------|
| S00 | P35-S00-00 → S00-01 → S00-02 | `INVESTIGATION.md` |
| S01 | P35-S01-00 → S01-01 | `PLAN.md` (+ acceptance tests sketch) |
| S02 | P35-S02-00 → S02-01 → S02-02 | Product fix + tests |
| S03 | P35-S03-00 → S03-01 → S03-02 | `VERIFY-NOTES.md` + DR-HANDOFF CLOSED |

Serial across scopes. No product code in planner rows.

## Hard constraints

- **Law 19** — HTTP + `web/` are adapters; prefer library-backed “current work” if one exists; else explicit GUI policy calling API (no business-logic fork)
- Do not delete feet-seller data
- Out of scope: weakening `plan_missing` for true PLAN work; hosted SaaS

## Success sketch

Opening `trace gui -C "<feet-seller>"` does **not** present Step 1 as the implied current task/gate when work has progressed to task ~123.
