# P38-S06-00 — Scope planner (remediation plan only)

## Metadata
- id: P38-S06-00
- todo_ids: [P38-S06-00]
- role: planner
- skills: [planning-and-task-breakdown, analyst]
- verification: automated

## Objective

Lock S06: **`REMEDIATION-PLAN.md` only** — ranked future phases/themes. **No product code. No implement rows in P38.**

## Prerequisite

S05-02 APPROVE (saturated) — `ready_for_REMEDIATION_PLAN: true`.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — § REMEDIATION-PLAN shape
- [GAP-REGISTRY.md](../scope-04-gap-registry/GAP-REGISTRY.md)
- [SATURATION-NOTES.md](../scope-05-saturation-gate/SATURATION-NOTES.md) — §4 H7 defer owner

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Inputs (required)

- GAP-REGISTRY.md, SATURATION-NOTES.md, DESIGN-LOCKS.md
- TRACE-AUDIT.md, PEER-CG.md, PEER-UA-GF.md (incl. MP §3)
- INVESTIGATION-INDEX.md
- [`h7-compose-desk-check.md`](../../../../../experiments/runs/2026-08-22-p38-s05-663/evidence/h7-compose-desk-check.md)
- [`h11-stack-docs.md`](../../../../../experiments/runs/2026-08-22-p38-s04-660/evidence/h11-stack-docs.md)

## Locked defaults

| Item | Value |
|------|-------|
| Artifact | `scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md` (S06-01 only) |
| Theme IDs | **G1–G9** consolidating G-001…G-011 |
| Ranking rubric | **(impact × law_fit) ÷ effort** — axes 1–5 |
| Rank order | G1 → G3 → G4 → G5 → G2 → G6 → G7 → G8 → G9 |
| H7 owner | **Compose-first UX before unified `trace_explore`** |
| H11 stack | **Doc-only** — not product integration |
| G-004a vector | **Reject/defer** — not a theme |
| Product edits | **Forbidden** in P38 |

## Must answer for 01 handoff

### 1. G1–G9 themes — ranked

See [SCOPE-TODOS.md](SCOPE-TODOS.md) theme registry + [01-plan.md](01-plan.md) pre-computed scores.

### 2. Each theme (summary)

| Theme | GAP ids | Peer pattern | Phase sketch | Not P38 |
|-------|---------|--------------|--------------|---------|
| G1 | G-001, G-002 | UA context-builder; MP wake_up | Phase 39 orient merge | No compiler/MCP code |
| G2 | G-007 | CG explore; compose split | Phase 39 compose / 40+ explore | No trace_explore ship |
| G3 | G-006, G-010 | CG SERVER_INSTRUCTIONS | Phase 39 harness | No server.go |
| G4 | G-011 | PEER-CG §5 complement | Phase 39 docs only | Doc-only H11 |
| G5 | G-008 | GF graph.html; UA onboard | Phase 39–40 GUI | No web/ code |
| G6 | G-004b | GF EXTRACTED/INFERRED | Phase 40+ law gate | No vector |
| G7 | G-005 | CG watcher | Phase 40+ index | No daemon |
| G8 | G-003 | MP layers (contrast) | Phase 41+ | No layer ship |
| G9 | G-009 | MP fact_checker contrast | Phase 41+ or doc-revise | No retrieval code |

### 3. Reject list (≥12)

SCOPE-TODOS reject seeds — PEER-CG §4 + SATURATION §5 + plan-specific (dual-stack product, compose-equivalence claim, P38 implement).

### 4. H11 stack recommendation

**Doc-only** — exhaustive grep found zero user dual-stack workflow ([h11-stack-docs.md](../../../../../experiments/runs/2026-08-22-p38-s04-660/evidence/h11-stack-docs.md)). Product integration rejected.

### 5. H7 owner

**Compose-first UX ranks above unified `trace_explore`** — desk-check 7/7 not equivalent; unified explore Phase 40+ after G1 + law spike.

## Planner gate

- [x] `01-plan.md` thickened (T0–T7, rubric, G1–G9 registry, H7/H11 locks)
- [x] `02-review.md` requires Checklists A–F
- [x] SCOPE-TODOS IDs 665–667 + theme registry
- [x] H7: compose-first before explore; H11: doc-only

## Exit criteria

- [x] S06-01/02 prompts runnable alone
- [x] Board `P38-S06-00` → `done`

## Next

`P38-S06-01`
