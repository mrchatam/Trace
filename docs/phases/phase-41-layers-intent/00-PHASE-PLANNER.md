# Phase 41+ — Layers & intent

**Phase planner.** Row `P41-00`. Human promotion required.

## Metadata
- id: P41-00
- todo_ids: [P41-00]
- role: planner
- skills: [planning-and-task-breakdown, incremental-implementation, context-engineering]
- verification: automated

## Mission

Lock Phase 41+ scopes against [Phase 38 REMEDIATION-PLAN](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) §2 G8/G9 + live repo + Phase 40 prerequisites. Thicken scope planners; confirm board row order. Read [`INTAKE.md`](INTAKE.md).

## Gate

**Do not run P41-00 until human promotes.** Phase 40 must be CLOSED (2026-08-22).

## Scope sequence

| Scope | Theme | Rows | Deliverable |
|-------|-------|------|-------------|
| S00 | G8 progressive layers | P41-S00-00 → 02 | L2–L3 ship or spec revise |
| S01 | G9 intent pipeline | P41-S01-00 → 02 | Intent implement or doc-revise |
| S02 | VERIFY | P41-S02-00 → 02 | VERIFY-NOTES + DR-HANDOFF CLOSED |

## Locked defaults

| Item | Value |
|------|-------|
| Entry themes | **G8 + G9** (REMEDIATION-PLAN §3 Phase 41+) |
| M-001 moat | Layer/intent changes merge into task loop — never replace moat |
| G8 | Ship L2–L3 in compiler **or** documented spec alternative |
| G9 | Implement intent extraction **or** mark RETRIEVAL_AND_CONTEXT §3 aspirational + supersede |
| G6/G7 | **Secondary queue** — document in INTAKE; no S03/S04 rows at P41-00 unless human promotes |
| G-004a vector | **Forbidden** — permanent defer |
| Law 6–7 | Progressive caps; no full-graph dump default |
| Law 19 | Library first; HTTP/MCP/GUI adapters thin |

## References

- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md) — Laws 6–7, 19
- [Phase 40 VERIFY-NOTES](../phase-40-read-surface-retrieval-depth/scopes/scope-02-verify/VERIFY-NOTES.md)
- [REMEDIATION-PLAN §2 G8/G9](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-003/G-009](../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [Phase 40 DR-HANDOFF](../phase-40-read-surface-retrieval-depth/DR-HANDOFF.md)

## Planner gate (P41-00)

- [ ] Re-read live repo: `internal/compiler/`, layer docs, `docs/RETRIEVAL_AND_CONTEXT.md` §3
- [ ] Thicken S00–S02 `00-PLANNER` + `01-*` + `02-*` prompts with file targets
- [ ] Board `docs/TODO/phase-41.md` — scope rows filled after P41-00
- [ ] Resolve INTAKE open questions (G9 implement vs doc-revise; G6/G7 secondary cut)
- [ ] DR-HANDOFF scope checklist accurate
- [ ] No product code in this row (planner only)

## Exit criteria

- Scope stubs runnable; board points to **P41-S00-00**
- G8/G9 locked with accept/reject per REMEDIATION-PLAN
- Secondary queue (G6, G7) documented in DR-HANDOFF forward note — not implement rows yet

## Next

`P41-S00-00`
