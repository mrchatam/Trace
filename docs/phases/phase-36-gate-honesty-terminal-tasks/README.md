# Phase 36 — Planning model alignment + plan_missing root cause

Human-promoted **2026-08-22**. Symptom: every feet-seller task shows **`plan_missing`**. User wants **root cause** (Trace vs agent vs harness), not a GUI patch.

## Design SoT

| Doc | Role |
|-----|------|
| [`INTAKE.md`](INTAKE.md) | Symptom + investigation facts + verdict sketch |
| [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md) | Fundamental fix options + must-preserve |
| [`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md) | Phase planner `P36-00` |
| [`DR-HANDOFF.md`](DR-HANDOFF.md) | Open until `P36-S03-02` |

Board: [`docs/TODO/phase-36.md`](../../TODO/phase-36.md).

## Confirmed root cause (P36-00 spot-check — S00 must formalize)

```
Agents wrote plan-change entities (11)  ──×──►  PlanExists gate
Agents never ran progressive planner  ──────►  plan_missing forever
MCP has no plan tools                 ──────►  agents can't fix via primary surface
feet-seller: no trace install         ──────►  enforce off; 123 DONE without gate block
```

**Code anchors (live repo):**

| Signal | Location | Behavior |
|--------|----------|----------|
| `PlanExists` | `internal/loop/policy.go:45–49` | Requires `current_scope_id` **and** `current_deep_plan` from progressive planner |
| `PlanCritiqued` | `internal/loop/policy.go:51–62` | Deliberation state **or** plan-change writes on apply path |
| `evaluateDone` | `internal/loop/gate.go:227–265` | No terminal `work_state` short-circuit; PLAN phase → `plan_missing` |
| MCP surface | `cmd/trace/plan.go:67`, `internal/mcp/server.go` | CLI `create-coarse` exists; MCP exposes `trace_add kind=plan-change` only |
| Install | `cmd/trace/install.go`, `internal/install/` | Optional; feet-seller never ran; no config → enforce off |

## Dogfood fixture

`/home/ali/Desktop/feet seller telegram app` — read-only unless S02 recovery explicitly scoped.

| Metric | Value |
|--------|-------|
| plan_phases / scopes / deep | 0 / 0 / 0 |
| plan_changes | 11 |
| tasks / goal | 123 / 1 |
| config.json | absent |
| trace install | not present |
| Step 1 / Loop 112 | `33247e2d-…` / `99d8fb92-…` |
| CLI gate (either, `--for done`) | `allowed=false`, `reason_code=plan_missing` (identical) |

## Scope sequence

| Scope | Board rows | Artifact |
|-------|------------|----------|
| S00 | P36-S00-00 → S00-02 | `INVESTIGATION.md` — Trace/agent/harness verdict |
| S01 | P36-S01-00 → S01-01 | `PLAN.md` — fundamental fix set |
| S02 | P36-S02-00 → S02-02 | Product + MCP/install + tests |
| S03 | P36-S03-00 → S03-02 | `VERIFY-NOTES.md` + DR-HANDOFF CLOSED |

```
S00 investigate — prove A/B/C/D; INVESTIGATION.md with verdict table
 → S01 plan — pick fundamental fix set (MCP plan, bootstrap, install, PlanExists bridge, terminal gate)
 → S02 implement + tests + review
 → S03 VERIFY — feet-seller + greenfield agent path + DR-HANDOFF
```

## In scope

- Two planning systems alignment (plan-change vs progressive planner)
- MCP plan tools and/or bootstrap command
- Install contract (AGENTS.md, cursor rules, bootstrap on new goal)
- Gate honesty for terminal DONE tasks (library first, GUI adapter)
- Feet-seller recovery path (backfill or documented legacy)

## Out of scope

- Hosted SaaS / multi-tenant
- Deleting feet-seller history
- Weakening PLAN gate for **active** non-terminal work
- GUI-only hide of `plan_missing` without agent-complete path

## Hard constraints

- Law 19 — library policy first; MCP/HTTP/GUI are adapters
- No "hide the error" without fixing agent-complete path
- Preserve PLAN gate for active work
