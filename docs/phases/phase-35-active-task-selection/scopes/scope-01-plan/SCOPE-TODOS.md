# Scope 01 — board map

**S01 plan** — selection policy + acceptance tests. Serial: **P35-S01-00 → P35-S01-01**. Artifact: `PLAN.md`. Requires **S00-02 PASS** + `INVESTIGATION.md`. No product code. No separate S01 review row — S02 review covers implementation fidelity to PLAN.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 614 | P35-S01-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock placement/semantics defaults; thicken 01 + this file + S02 seed note |
| 615 | P35-S01-01 | [01-plan.md](01-plan.md) | Implementer | Author `PLAN.md` only (no product code) |

## Inputs

- `scopes/scope-00-investigate/INVESTIGATION.md` (required, S00-02 PASS)
- DESIGN-LOCKS + INTAKE + Phase README

## Locked for S01-01 (summary)

- Root cause frozen: Overview L17–20 + Loop L51–54 → Step1; no library current-work; HTTP ignores limit
- Placement closed set: **A** library+adapters vs **B** shared `web/src/lib` helper (lean B unless A fits S02)
- Semantics lean: IN_PROGRESS → non-terminal → **last** DONE (not `tasks[0]`)
- Must-test: all-DONE ≠ Step1; >100/`limit` honesty; explicit `task_id` wins
- Out: `plan_missing` weaken; TRACE_TASK_ID-only; delete dogfood; Overview/Loop fork

## Out of this scope

- Shipping code (S02), live VERIFY (S03), writing `PLAN.md` in the planner row (S01-00).
