# DESIGN-LOCKS — Phase 36

**Human-promoted 2026-08-22.** Reframed from GUI patch → **planning model + agent workflow alignment**.

| Lock | Value |
|------|-------|
| Dogfood fixture | `/home/ali/Desktop/feet seller telegram app` (123 DONE, 1 goal, 11 plan-changes, 0 progressive plan) |
| Theme | **Why `plan_missing` everywhere** — Trace product vs agent workflow; fix fundamentally |
| Must investigate | Two planning systems; MCP gap; enforce/install posture; mega-goal pattern; transition path without gate |
| Must fix (product) | Agents using Trace correctly **can** satisfy `PlanExists` without undocumented CLI-only steps |
| Must fix (honesty) | Gate/GUI must reflect **goal-level** plan gap clearly; terminal DONE tasks must not imply actionable "done blocked" for finished work |
| Must preserve | PLAN enforcement for **active** non-terminal work; fail-closed deliberation order |
| Law 19 | Policy in library; MCP/HTTP/GUI are adapters |
| Code anchors | `PlanExists` → `internal/loop/policy.go:45–49`; `PlanCritiqued` incl. plan-changes → `:60–62`; `evaluateDone` no terminal short-circuit → `internal/loop/gate.go:227–265`; MCP gap → `cmd/trace/plan.go:67` |
| Out of scope | Hosted SaaS; deleting feet-seller history; auto-LLM backlog generation |

## Candidate fundamental fixes (S01 chooses — not all required in one phase)

| Option | Description |
|--------|-------------|
| **MCP plan tools** | `trace_plan` MCP: create-coarse, set-current, deep, show |
| **PlanExists bridge** | Recognize plan-change density / bootstrap heuristic (careful — must not fake plan) |
| **Bootstrap command** | `trace plan bootstrap --goal` from plan-changes or seed import |
| **Install contract** | First goal task → mandatory create-coarse in AGENTS.md + cursor rules |
| **Enforce nudge** | Document or default `warn` when `.trace/` exists without config |
| **Goal structure** | Warn when >N tasks under goal with no coarse plan |
| **Terminal gate** | Library: DONE/SKIPPED tasks skip premature done-blocked for plan_missing |

## Success sketch

Greenfield agent via MCP: creates goal → **bootstrap coarse plan** → edit gate passes → work proceeds. Feet-seller: recovery path documented or executed; GUI no longer shows identical red alarm on all 123 historical tasks without explaining the **goal** has no progressive plan.
