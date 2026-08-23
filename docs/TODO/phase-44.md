Index: [`docs/TODO.md`](../TODO.md)

## Phase 44 — Scoped graph linking & scope semantics

**Active — promoted 2026-08-23.** Plan scaffold from user intake (post Phase 43 idle).

Design SoT: [`phases/phase-44-scoped-graph-linking/00-PHASE-PLANNER.md`](../phases/phase-44-scoped-graph-linking/00-PHASE-PLANNER.md) · [`INTAKE.md`](../phases/phase-44-scoped-graph-linking/INTAKE.md) · [`01-DESIGN-LOCKS.md`](../phases/phase-44-scoped-graph-linking/01-DESIGN-LOCKS.md)

Scope sequence:
- **S00** — Intake + research (link model, hairball root cause, gap matrix)
- **S01** — Design locks (scope model, rel taxonomy, cap policy)
- **S02** — Backend links (entity_links, domain, MCP, OpenAPI, graph walk)
- **S03** — Inference (opt-in derived edges + provenance)
- **S04** — GUI (scope-aware layout; Law 19 adapter)
- **S05** — VERIFY + DR-HANDOFF

| Order | ID | Status | Prompt | Notes |
|------:|----|--------|--------|-------|
| 721 | P44-00 | pending | [phases/phase-44-scoped-graph-linking/00-PHASE-PLANNER.md](../phases/phase-44-scoped-graph-linking/00-PHASE-PLANNER.md) | **2026-08-23 scaffold:** Phase 44 promoted from user intake — project graph lacks scoped interconnection (auth FE/BE, business scopes); not UI-only. Draft locks in 01-DESIGN-LOCKS; 500 cap = GUI default (Law 6–7), API max 5000. Successor TBD at VERIFY. Next **P44-S00-00** after P44-00. |
| 722 | P44-S00-00 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-00-intake-research/00-PLANNER.md](../phases/phase-44-scoped-graph-linking/scopes/scope-00-intake-research/00-PLANNER.md) | Research current link model: project_graph, graph_neighbors, entity_links, trace_link, goal_id hairball. Deliver gap matrix vs INTAKE. |
| 723 | P44-S00-01 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-00-intake-research/01-research.md](../phases/phase-44-scoped-graph-linking/scopes/scope-00-intake-research/01-research.md) | Implement research doc + board notes; no product code unless trivial doc fix. |
| 724 | P44-S00-02 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-00-intake-research/02-review.md](../phases/phase-44-scoped-graph-linking/scopes/scope-00-intake-research/02-review.md) | Review research completeness vs INTAKE examples. |
| 725 | P44-S01-00 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-01-design-locks/00-PLANNER.md](../phases/phase-44-scoped-graph-linking/scopes/scope-01-design-locks/00-PLANNER.md) | Lock D1–D5 in 01-DESIGN-LOCKS; thicken S02–S04 implement prompts. |
| 726 | P44-S01-01 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-01-design-locks/01-lock.md](../phases/phase-44-scoped-graph-linking/scopes/scope-01-design-locks/01-lock.md) | Finalize design locks doc; close open decisions. |
| 727 | P44-S01-02 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-01-design-locks/02-review.md](../phases/phase-44-scoped-graph-linking/scopes/scope-01-design-locks/02-review.md) | Independent lock review before implement wave. |
| 728 | P44-S02-00 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-02-backend-links/00-PLANNER.md](../phases/phase-44-scoped-graph-linking/scopes/scope-02-backend-links/00-PLANNER.md) | Backend: migration, domain rels, store, retrieval walk, MCP trace_link, OpenAPI. |
| 729 | P44-S02-01 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-02-backend-links/01-implement.md](../phases/phase-44-scoped-graph-linking/scopes/scope-02-backend-links/01-implement.md) | Implement scope link rels + graph API fields. |
| 730 | P44-S02-02 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-02-backend-links/02-review.md](../phases/phase-44-scoped-graph-linking/scopes/scope-02-backend-links/02-review.md) | Review backend + export round-trip. |
| 731 | P44-S03-00 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-03-inference/00-PLANNER.md](../phases/phase-44-scoped-graph-linking/scopes/scope-03-inference/00-PLANNER.md) | Inference rules, provenance, opt-in CLI/hook. |
| 732 | P44-S03-01 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-03-inference/01-implement.md](../phases/phase-44-scoped-graph-linking/scopes/scope-03-inference/01-implement.md) | Implement inference pass (explicit wins). |
| 733 | P44-S03-02 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-03-inference/02-review.md](../phases/phase-44-scoped-graph-linking/scopes/scope-03-inference/02-review.md) | Review inference honesty + caps. |
| 734 | P44-S04-00 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-04-gui/00-PLANNER.md](../phases/phase-44-scoped-graph-linking/scopes/scope-04-gui/00-PLANNER.md) | GUI scope layout, edge priority, orient copy (Law 19). |
| 735 | P44-S04-01 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-04-gui/01-implement.md](../phases/phase-44-scoped-graph-linking/scopes/scope-04-gui/01-implement.md) | Implement scope-aware graph presentation. |
| 736 | P44-S04-02 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-04-gui/02-review.md](../phases/phase-44-scoped-graph-linking/scopes/scope-04-gui/02-review.md) | Review GUI + a11y + perf at 500 cap. |
| 737 | P44-S05-00 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-05-verify/00-PLANNER.md](../phases/phase-44-scoped-graph-linking/scopes/scope-05-verify/00-PLANNER.md) | VERIFY planner — blocks, evidence paths. |
| 738 | P44-S05-01 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-05-verify/01-verify.md](../phases/phase-44-scoped-graph-linking/scopes/scope-05-verify/01-verify.md) | Phase VERIFY gate. |
| 739 | P44-S05-02 | pending | [phases/phase-44-scoped-graph-linking/scopes/scope-05-verify/02-dr-handoff.md](../phases/phase-44-scoped-graph-linking/scopes/scope-05-verify/02-dr-handoff.md) | DR-HANDOFF close + successor decision. |
