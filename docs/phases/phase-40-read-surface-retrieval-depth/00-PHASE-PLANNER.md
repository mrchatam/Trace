# Phase 40+ — Read surface & retrieval depth

**Phase planner.** Row `P40-00`. Human promotion required.

## Metadata
- id: P40-00
- todo_ids: [P40-00]
- role: planner
- skills: [planning-and-task-breakdown, incremental-implementation, context-engineering]
- verification: automated

## Mission

Lock Phase 40+ scopes against [Phase 38 REMEDIATION-PLAN](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) §2 G5/G2 + live repo + G1 dependency. Thicken scope planners; confirm board row order. Read [`INTAKE.md`](INTAKE.md).

## Gate

**Do not run P40-00 until human promotes.** Phase 39 must be CLOSED (2026-08-22).

## Scope sequence

| Scope | Theme | Rows | Deliverable |
|-------|-------|------|-------------|
| S00 | G5 GUI orient | P40-S00-00 → 02 | Graph-first onboarding UX adapter |
| S01 | G2 unified explore | P40-S01-00 → 02 | Task-aware capped `trace_explore` |
| S02 | VERIFY | P40-S02-00 → 02 | VERIFY-NOTES + DR-HANDOFF CLOSED |

## Locked defaults

| Item | Value |
|------|-------|
| Entry themes | **G5 + G2** (REMEDIATION-PLAN §3 Phase 40+) |
| M-001 moat | Unified explore merges into task loop — never query-only replacement |
| G5 | **Law 19 adapter only** — GUI over canonical library/API |
| G2 | **After G1 + law spike** — task-aware capped read; no mega-tool facade |
| G-004a vector | **Forbidden** — permanent defer |
| Law 6–7 | Progressive caps; no full-graph dump default |
| Law 19 | Library first; HTTP/MCP/GUI adapters thin |

## References

- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md) — Laws 6–7, 19
- [Phase 39 VERIFY-NOTES](../phase-39-context-orient-harness/scopes/scope-03-verify/VERIFY-NOTES.md)
- [REMEDIATION-PLAN §2 G5/G2](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-007/G-008](../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [Phase 39 DR-HANDOFF](../phase-39-context-orient-harness/DR-HANDOFF.md)

## Planner gate (P40-00)

- [x] Re-read live repo: `web/`, `internal/httpapi/`, `internal/mcp/`, G1 query merge paths
- [x] Thicken S00–S02 `00-PLANNER` + `01-*` + `02-*` prompts with file targets
- [x] Board `docs/TODO/phase-40.md` — scope rows filled after P40-00
- [x] Resolve INTAKE open questions (G6/G7 secondary scope cut, G2 law spike gate)
- [x] DR-HANDOFF scope checklist accurate
- [x] No product code in this row (planner only)

## Exit criteria

- Scope stubs runnable; board points to **P40-S00-00**
- G5/G2 locked with accept/reject per REMEDIATION-PLAN
- Secondary queue (G6, G7) documented in DR-HANDOFF forward note — not implement rows yet

## Next

`P40-S00-00`
