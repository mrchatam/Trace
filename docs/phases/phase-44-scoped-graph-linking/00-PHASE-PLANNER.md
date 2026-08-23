# Phase 44 — Scoped graph linking & scope semantics

**Phase planner.** Row `P44-00`. Human promotion required.

## Metadata
- id: P44-00
- todo_ids: [P44-00]
- role: planner
- skills: [planning-and-task-breakdown, domain-modeling, incremental-implementation]
- mcps: [trace_context, trace_why]
- verification: automated

## Mission

Lock Phase 44 scopes against [`INTAKE.md`](INTAKE.md) and draft [`01-DESIGN-LOCKS.md`](01-DESIGN-LOCKS.md). Thicken scope planners; confirm board row order. **No product code.**

## Gate

**Do not run P44-00 until human promotes.** Phase 43 complete (2026-08-22).

## Scope sequence

| Scope | Theme | Rows | Deliverable |
|-------|-------|------|-------------|
| S00 | Intake + research | P44-S00-00 → 02 | Current link model + gap matrix |
| S01 | Design locks | P44-S01-00 → 02 | Finalize 01-DESIGN-LOCKS (D1–D5 closed) |
| S02 | Backend links | P44-S02-00 → 02 | entity_links rels, domain, MCP, OpenAPI, graph walk |
| S03 | Inference | P44-S03-00 → 02 | Opt-in derived links + provenance |
| S04 | GUI | P44-S04-00 → 02 | Scope-aware layout + orient (Law 19) |
| S05 | VERIFY | P44-S05-00 → 02 | VERIFY-NOTES + DR-HANDOFF |

## Locked defaults (draft)

| Item | Value |
|------|-------|
| Entry | User intake 2026-08-23 — scoped interconnection, not UI-only |
| M-001 moat | Scope links enrich graph; never replace task loop |
| Law 6–7 | Bounded graph API; 500 GUI default unless S01 raises with scope filter |
| Law 13 | Reuse `entity_links`; no parallel relationship store |
| Law 19 | GUI/HTTP adapters thin; walk logic in `internal/retrieval` |
| DR-NOSSEM | No vector/semantic scope clustering in P44 |
| Phase 43 | Do not reopen GitHub hygiene rows |

## Research anchors (S00)

| File | Question |
|------|----------|
| `internal/retrieval/project_graph.go` | Center + truncation |
| `internal/retrieval/graph_neighbors.go` | Walk edges |
| `internal/domain/doc.go` | Documented rels |
| `internal/mcp/tools_write.go` | MCP link surface |
| `internal/store/links.go` | Persistence |
| `web/src/lib/overviewCompose.ts` | 500 cap |
| `api/openapi.yaml` | Contract caps |

## References

- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md) — Laws 6–7, 13, 19
- [G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md)
- [Phase 40 G5](../phase-40-read-surface-retrieval-depth/README.md) — graph orient baseline

## Planner gate (P44-00)

- [ ] Re-read live repo files in Research anchors
- [ ] Thicken S00–S05 `00-PLANNER` + `01-*` + `02-*` prompts
- [ ] Board `docs/TODO/phase-44.md` — scope rows filled
- [ ] Close open decisions D1–D5 in 01-DESIGN-LOCKS or defer with reason
- [ ] No product code in this row

## Exit criteria

- Scope stubs runnable; board points to **P44-S00-00**
- 01-DESIGN-LOCKS ready for S01 lock pass
- Successor after P44: **TBD at VERIFY** (default `no successor`)

## Next

`P44-S00-00`
