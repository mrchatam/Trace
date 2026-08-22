# Phase 36 — Planning model alignment + plan_missing root cause

**Phase planner.** Row `P36-00`.

## Metadata
- id: P36-00
- todo_ids: [P36-00]
- role: planner
- skills: [diagnosing-bugs, domain-modeling, planning-and-task-breakdown]
- verification: automated

## Mission

Find **why** every feet-seller task shows `plan_missing`; determine Trace vs agent vs harness blame; fix **fundamentally** (not GUI-only).

Read [`INTAKE.md`](INTAKE.md) + [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md).

## Gate

Phase 35 **done**. Proceed.

## Scope sequence

| Scope | Rows | Artifact |
|-------|------|----------|
| S00 | P36-S00-00 → 02 | `INVESTIGATION.md` — verdict table Trace/agent/harness |
| S01 | P36-S01-00 → 01 | `PLAN.md` — pick fix set from DESIGN-LOCKS candidates |
| S02 | P36-S02-00 → 02 | Product + MCP/install + tests |
| S03 | P36-S03-00 → 02 | VERIFY + DR-HANDOFF |

```
S00 investigate — prove A/B/C/D; INVESTIGATION.md with verdict table
 → S01 plan — pick fundamental fix set (MCP plan, bootstrap, install, PlanExists bridge, terminal gate)
 → S02 implement + tests + review
 → S03 VERIFY — feet-seller + greenfield agent path + DR-HANDOFF
```

## Hard constraints

- No GUI-only patch as sole deliverable
- MCP or documented bootstrap must let agents satisfy PlanExists
- Preserve active-work PLAN enforcement

## Planner gate (P36-00)

- [x] Phase folder README + locks reflect fundamental scope
- [x] S00 prompts ask Trace vs agent vs harness questions
- [x] S01 lists MCP plan / bootstrap / install as first-class options
- [x] DR-HANDOFF OPEN

## P36-00 outcome (2026-08-22)

Gate **PASS**. Thickened README (code anchors, scope table, in/out). DESIGN-LOCKS + live cites (`policy.go`, `gate.go`, `plan.go`). Protocol-thickened S00–S03 prompts + `SCOPE-TODOS.md`. S01-00/01 rewritten: MCP plan / bootstrap / install as first-class fix options (not GUI-only). DR-HANDOFF theme aligned to planning-model scope. Spot-checked feet-seller: 123 tasks all DONE; Step1 + Loop112 gate JSON identical `plan_missing`; progressive planner empty; MCP catalog has no plan tools. No product code. Next: **P36-S00-00**.

## Next

`P36-S00-00`
