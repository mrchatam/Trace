# Phase 39 — Context orient & harness

**Phase planner.** Row `P39-00`. Human promotion required.

## Metadata
- id: P39-00
- todo_ids: [P39-00]
- role: planner
- skills: [planning-and-task-breakdown, incremental-implementation, context-engineering]
- verification: automated

## Mission

Lock Phase 39 scopes against [Phase 38 REMEDIATION-PLAN](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) §2 G1/G3/G4 + live repo. Thicken scope planners; confirm board row order. Read [`INTAKE.md`](INTAKE.md).

## Gate

**Do not run P39-00 until human promotes.** Phase 38 must be CLOSED (2026-08-22).

## Scope sequence

| Scope | Theme | Rows | Deliverable |
|-------|-------|------|-------------|
| S00 | G1 context merge | P39-S00-00 → 02 | Query+task orient in compiler/MCP |
| S01 | G3 harness | P39-S01-00 → 02 | MCP playbook + bootstrap orient |
| S02 | G4 docs | P39-S02-00 → 02 | Dual-stack CONTRIBUTING/AGENTS |
| S03 | VERIFY | P39-S03-00 → 02 | VERIFY-NOTES + DR-HANDOFF CLOSED |

## Locked defaults

| Item | Value |
|------|-------|
| Entry co-wave | **G1 + G3 + G4** (REMEDIATION-PLAN §3) |
| M-001 moat | All changes merge into task loop — never replace moat |
| G4 | **Doc-only** — no product dual-index default |
| G2 explore | **Forbidden** in P39 — Phase 40+ after G1 + law spike |
| Law 6–7 | Progressive caps; no full-graph dump default |
| Law 19 | Library first; HTTP/MCP/GUI adapters thin |

## References

- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md) — Laws 6–7, 19
- [REMEDIATION-PLAN §2 G1/G3/G4](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY.md](../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [Phase 38 DR-HANDOFF](../phase-38-retrieval-context-peer-gaps/DR-HANDOFF.md)

## Planner gate (P39-00)

- [x] Re-read live repo: `internal/compiler/`, `internal/mcp/`, `cmd/trace/context`, harness install paths
- [x] Thicken S00–S03 `00-PLANNER` + `01-*` + `02-*` prompts with file targets
- [x] Board `docs/TODO/phase-39.md` — scope rows filled after P39-00
- [x] Resolve INTAKE open questions (scope cut, 9/16 fix docs vs code)
- [x] DR-HANDOFF scope checklist accurate
- [x] No product code in this row (planner only)

## Exit criteria

- Scope stubs runnable; board points to **P39-S00-00**
- G1/G3/G4 locked with accept/reject per REMEDIATION-PLAN
- Secondary queue (G5, G2) documented in DR-HANDOFF forward note — not implement rows yet

## Next

`P39-S00-00`
