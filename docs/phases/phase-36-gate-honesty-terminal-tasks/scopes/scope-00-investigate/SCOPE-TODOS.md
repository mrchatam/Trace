# Scope 00 — board map

**S00 investigate** — repro feet-seller gate behavior + policy root cause. Serial: **P36-S00-00 → P36-S00-01 → P36-S00-02**. Primary artifact: `INVESTIGATION.md` (written in **S00-01**, reviewed in **S00-02**). Do **not** start S01 until S00-02 PASS. Do **not** write product code. Planner (**S00-00**) does **not** author `INVESTIGATION.md`.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 623 | P36-S00-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 + this file |
| 624 | P36-S00-01 | [01-investigate.md](01-investigate.md) | Implementer | Author `INVESTIGATION.md` + red-capable repro notes |
| 625 | P36-S00-02 | [02-review.md](02-review.md) | Reviewer | Checklist vs DESIGN-LOCKS + INTAKE + live cites |

## Planner-locked live facts (for 01) — verified 2026-08-22 (P36-S00-00 replan)

- Fixture: `/home/ali/Desktop/feet seller telegram app` — **read-only** against `.trace/` (gate/export need DB read; do not mutate).
- CLI: `trace -C "<fixture>" tasks` → **123** tasks, **all DONE**, **1 goal**.
- Goal: `353b12a4-57dd-4d68-8379-b2024e064733` — plan empty (`phases: []`, `current_scope_id: null`, `current_deep_plan: null`).
- Step 1: `33247e2d-aa10-4b25-b194-4b7afb5a6359`; Loop 112: `99d8fb92-65ac-462c-82c4-21bcf198c09e`.
- `trace loop gate --task <either> --for done` → `plan_missing`, `allowed=false`, `recommended_phase=PLAN` (**identical** JSON on both tasks).
- `trace seed export -o /tmp/feet-export.json` → `plan_changes: 11`; `plan_phases/plan_scopes/scope_deep_plans`: **0**.
- Fixture missing: `AGENTS.md`, `.trace/config.json`, `.cursor/rules/trace-enforcement.mdc` (enforce **off** by default).
- `PlanExists`: `internal/loop/policy.go:45–49` — needs `current_scope_id` + `current_deep_plan` (progressive planner only).
- `PlanCritiqued`: `internal/loop/policy.go:60–62` — deliberation state **or** plan-changes on apply path; **not** PlanExists.
- `deliberation/select.go:28–29` — `!PlanExists` → Phase PLAN / `plan_missing`.
- `evaluateDone`: `internal/loop/gate.go:227–265` — no terminal `work_state` short-circuit.
- MCP: `trace_add` incl. `plan-change`; **no** plan tools (`cmd/trace/plan.go:67`); `trace_loop` = next|apply|status only (no gate).
- Transition: `--enforce` opt-in on DONE (`cmd/trace/transition.go:36–69`); MCP `trace_transition` mirrors CLI.
- TaskDetail: `web/src/screens/TaskDetail.tsx:76–86`, `:198–210` — always fetches done gate + GateStrip.
- GateStrip: `web/src/components/GateStrip.tsx:41–56` — blocked styling + reason_code display.

## Investigation rejects (01 must document)

1. Blaming Phase 35 pick logic — selection is fixed; this is gate **meaning** on terminal tasks.
2. "Just create a plan" as the **only** fix — does not explain misleading DONE gate on finished work.
3. Weakening `plan_missing` globally for active PLAN work.
4. Deleting feet-seller history.

## Out of this scope

- Writing `PLAN.md` (S01), product fixes (S02), VERIFY (S03).
- Authoring `INVESTIGATION.md` in the planner row (S00-00).
