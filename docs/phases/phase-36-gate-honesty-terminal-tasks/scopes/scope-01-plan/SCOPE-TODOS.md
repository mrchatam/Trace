# Scope 01 — board map

**S01 plan** — fundamental fix set from S00 verdict. Serial: **P36-S01-00 → P36-S01-01**. Primary artifact: `PLAN.md`. S00-02 PASS (high confidence) — unblocked.

| Order | Board ID | Prompt | Role | Artifact |
|------:|----------|--------|------|----------|
| 626 | P36-S01-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock fix options; thicken 01-plan |
| 627 | P36-S01-01 | [01-plan.md](01-plan.md) | Implementer | Author `PLAN.md` |

## First-class fix options (PLAN.md must decide each)

From DESIGN-LOCKS + INVESTIGATION.md ranked recommendations. **MCP / bootstrap / install are primary axis** — not footnotes.

| Option | Description | S00 default | PLAN.md section |
|--------|-------------|-------------|-----------------|
| **MCP plan tools** | `trace_plan` MCP: create-coarse, set-current, deep, show — Law 19 adapter over `internal/planner` | **accept** | §2.1 |
| **Bootstrap command** | `trace plan bootstrap --goal` from plan-changes or seed import | **accept** | §2.2 |
| **Install contract** | First goal → mandatory create-coarse in AGENTS.md + cursor rules + install docs | **accept** | §2.3 |
| **Terminal gate** | Library: DONE/SKIPPED skip misleading done-blocked for goal-level `plan_missing` | **accept** | §2.5 |
| **Goal structure** | Warn when >N tasks under goal with no coarse plan | accept/defer | §2.7 |
| **Enforce nudge** | Document or default `warn` when `.trace/` exists without config | accept/defer | §2.6 |
| **PlanExists bridge** | Plan-change density heuristic (must not fake progressive plan) | **defer** | §2.4 |
| **Feet-seller recovery** | S02 execute vs S03 verify vs defer | decide | §2.8 |

## Locked touch-list order (S02 follows PLAN.md)

```text
library (loop/policy, gate) → planner library → MCP → CLI → install → config nudge → HTTP → GUI
```

Law 19: policy in library; MCP/HTTP/GUI are adapters only.

## Locked acceptance tests (PLAN.md §6 → S02 implements)

| Test | Purpose |
|------|---------|
| `TestGreenfield_MCPPlanBootstrap_EditGatePasses` | Agent-complete path: MCP bootstrap → edit gate passes |
| `TestLegacy_FeetSellerExport_GateHonestyUntilBootstrap` | Terminal honesty + post-bootstrap recovery |
| `TestActiveWork_PlanMissingStillBlocksEdit` | Preserve non-terminal PLAN enforcement |

Reference patterns: `cmd/trace/loop_test.go:1137–1167`, `internal/loop/gate_test.go:198–205`.

## Planner inputs (from S00 — locked)

- INVESTIGATION.md verdict: Trace **primary** (dual planning + MCP gap); agent+harness secondary; GUI tertiary
- Export counts: **11 plan-changes / 0 progressive planner** on feet-seller
- MCP: 15 tools, no `trace_plan`; `trace_loop` = next|apply|status only
- Install: no `create-coarse` in enforcement docs; fixture missing AGENTS/config/hooks
- `PlanExists`: `policy.go:45–49` — requires `current_scope_id` + `current_deep_plan`
- `evaluateDone`: no terminal `work_state` short-circuit (`gate.go:227–265`)

## Explicit non-goals (PLAN.md §7)

- GUI-only patch without agent-complete path
- Global `plan_missing` weaken for active work
- Deleting feet-seller history
- Hosted SaaS; auto-LLM backlog generation

## Out of scope

- Product code (S02)
- VERIFY live dogfood (S03)
- Authoring `PLAN.md` in planner row S01-00 (that's S01-01)
