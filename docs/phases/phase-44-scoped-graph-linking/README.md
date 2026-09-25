# Phase 44 — Scoped graph linking & scope semantics

Human-promoted successor to Phase 43 idle. **P44-00 phase planner complete** 2026-08-23 (docs scaffold; no product code yet).

## Goal

Make the project graph show **meaningful scoped interconnection** — not a hairball radiating from the first goal. Introduce scope as a first-class graph concept (thin `scope` records + typed link relations), explicit + inferred edges, and GUI layout that reads clusters instead of star topology.

## Problem (user intake)

- Current `mode=project` graph: nodes mostly connect via `goal_has_task` / causal rels to the **first goal** — weak logical scope grouping.
- Desired: frontend auth ↔ frontend auth; backend auth ↔ backend auth; front↔back via API contract edges; broader business scopes (design, marketing, operator, goals/tasks).
- **Not a UI-only fix** — requires domain link taxonomy, store/API, optional inference, then GUI.

## Design SoT

| Doc | Role |
|-----|------|
| [`INTAKE.md`](INTAKE.md) | User problem, examples, in/out |
| [`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md) | Row `P44-00` (done) |
| [`01-DESIGN-LOCKS.md`](01-DESIGN-LOCKS.md) | Light-locked D1–D5 pending S01 APPROVE |
| [`DR-HANDOFF.md`](DR-HANDOFF.md) | OPEN until VERIFY close |

Board: [`docs/TODO/phase-44.md`](../../TODO/phase-44.md).

## Scope sequence (run order)

```
S00 Intake + research (current link model, gap vs vision)
 → S01 Design locks (D1–D5 → LOCKED; thicken S02–S04)
 → S02 Backend links (entity_links rels, domain, MCP, seed export)
 → S03 Inference (CLI-only derived edges; explicit provenance)
 → S04 GUI (scope-aware layout, edge priority, orient copy)
 → S05 VERIFY + DR-HANDOFF
```

| Scope | Theme | Artifact focus | First row |
|-------|-------|----------------|-----------|
| S00 | Research | Link creation paths, hairball root cause, gap matrix | P44-S00-00 |
| S01 | Design locks | Taxonomy + cap policy locked for implement wave | P44-S01-00 |
| S02 | Backend | New rels, `trace_link` parity, graph walk, OpenAPI | P44-S02-00 |
| S03 | Inference | Opt-in CLI derived links; never silent overwrite of explicit | P44-S03-00 |
| S04 | GUI | Force layout uses scope edges; not cosmetic-only | P44-S04-00 |
| S05 | VERIFY | Laws 6–7, M-001, portable graph export if schema changes | P44-S05-00 |

## In scope (phase 1 — light-locked)

- Scope concept: thin **`scope` records** + **`scope_member`** + MVP typed rels (`api_contract`, `implements`, `blocks`).
- Explicit agent links via MCP/CLI; inference as **secondary** with `INFERRED` provenance (CLI only).
- Project graph API: `scope` filter + optional `scope_id` / edge provenance on `BoundedGraph`; **500 GUI default**, API max 5000.
- GUI: layout groups by scope / rel priority (extends Phase 40 G5 orient route); Law 19 adapter only.

## Out of scope (phase 1)

- Vector/semantic clustering (DR-NOSSEM)
- Full Graphify port or Codegraph cross-index
- Always-on daemon, index-hook inference, or unbounded graph dump (Laws 6–7)
- Raising GUI default above 500 without VERIFY evidence
- Pagination/tiles for project graph
- Replacing task moat or plan tree as source of order

## Moat charter (M-001)

Scoped linking **enriches** the causal graph; it does not replace Tasks → Loop → Gate → Review. New rels must merge into existing `entity_links` discipline (Law 13 — no parallel relationship store).

## Next runnable

**P44-S00-00** — scope planner for intake research.
