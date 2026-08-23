# Phase 44 — Scoped graph linking & scope semantics

Human-promoted successor to Phase 43 idle. **Plan-only scaffold** (2026-08-23).

## Goal

Make the project graph show **meaningful scoped interconnection** — not a hairball radiating from the first goal. Introduce scope as a first-class graph concept (tags, clusters, or typed link relations), explicit + inferred edges, and GUI layout that reads clusters instead of star topology.

## Problem (user intake)

- Current `mode=project` graph: nodes mostly connect via `goal_has_task` / causal rels to the **first goal** — weak logical scope grouping.
- Desired: frontend auth ↔ frontend auth; backend auth ↔ backend auth; front↔back via API contract edges; broader business scopes (design, marketing, operator, goals/tasks).
- **Not a UI-only fix** — requires domain link taxonomy, store/API, optional inference, then GUI.

## Design SoT

| Doc | Role |
|-----|------|
| [`INTAKE.md`](INTAKE.md) | User problem, examples, in/out |
| [`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md) | Row `P44-00` |
| [`01-DESIGN-LOCKS.md`](01-DESIGN-LOCKS.md) | Scope model, link taxonomy, inference policy, API cap, non-goals |
| [`DR-HANDOFF.md`](DR-HANDOFF.md) | OPEN until VERIFY close |

Board: [`docs/TODO/phase-44.md`](../../TODO/phase-44.md).

## Scope sequence

```
S00 Intake + research (current link model, gap vs vision)
 → S01 Design locks (finalize 01-DESIGN-LOCKS)
 → S02 Backend links (entity_links rels, domain, MCP, seed export)
 → S03 Inference (derived edges from paths, titles, plan scopes — explicit provenance)
 → S04 GUI (scope-aware layout, edge priority, orient copy)
 → S05 VERIFY + DR-HANDOFF
```

| Scope | Theme | Artifact focus |
|-------|-------|----------------|
| S00 | Research | Link creation paths, hairball root cause, gap matrix |
| S01 | Design locks | Taxonomy + cap policy locked for implement wave |
| S02 | Backend | New rels, `trace_link` parity, graph walk, OpenAPI |
| S03 | Inference | Opt-in derived links; never silent overwrite of explicit |
| S04 | GUI | Force layout uses scope edges; not cosmetic-only |
| S05 | VERIFY | Laws 6–7, M-001, portable graph export if schema changes |

## In scope (phase 1 cut — draft)

- Scope concept: **`scope` tag/cluster** on entities and/or **`scope_member`** + typed cross-scope rels (`api_contract`, `implements`, `blocks`, `same_feature_front_back`, …).
- Explicit agent links via MCP/CLI; inference as **secondary** with provenance label.
- Project graph API: scope-aware edges in `BoundedGraph`; cap policy (retain 500 GUI default vs tiles/pagination — locked in S01).
- GUI: layout groups by scope / rel priority (extends Phase 40 G5 orient route).

## Out of scope (phase 1)

- Vector/semantic clustering (DR-NOSSEM)
- Full Graphify port or Codegraph cross-index
- Always-on daemon or unbounded graph dump (Laws 6–7)
- Replacing task moat or plan tree as source of order

## Moat charter (M-001)

Scoped linking **enriches** the causal graph; it does not replace Tasks → Loop → Gate → Review. New rels must merge into existing `entity_links` discipline (Law 13 — no parallel relationship store).
