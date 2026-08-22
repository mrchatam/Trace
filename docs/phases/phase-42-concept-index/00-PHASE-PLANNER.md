# Phase 42+ — Concept & index

**Phase planner.** Row `P42-00`. Human promotion required.

## Metadata
- id: P42-00
- todo_ids: [P42-00]
- role: planner
- skills: [planning-and-task-breakdown, incremental-implementation, context-engineering]
- verification: automated

## Mission

Lock Phase 42+ scopes against [Phase 38 REMEDIATION-PLAN](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) §2 G6/G7 + live repo + Phase 41 prerequisites. Thicken scope planners; confirm board row order. Read [`INTAKE.md`](INTAKE.md).

## Gate

**Do not run P42-00 until human promotes.** Phase 41 must be CLOSED (2026-08-22).

## Scope sequence

| Scope | Theme | Rows | Deliverable |
|-------|-------|------|-------------|
| S00 | G6 non-semantic concept | P42-S00-00 → 02 | Graph-label retrieval channel |
| S01 | G7 index freshness & langs | P42-S01-00 → 02 | Lang policy + index honesty |
| S02 | VERIFY | P42-S02-00 → 02 | VERIFY-NOTES + DR-HANDOFF CLOSED |

## Locked defaults

| Item | Value |
|------|-------|
| Entry themes | **G6 + G7** (REMEDIATION-PLAN rank 6–7) |
| M-001 moat | Concept/index changes merge into task loop — never replace moat |
| G6 | Graph-label channel under **DR-NOSSEM** — no vector (G-004a forbidden) |
| G7 | Lang expansion policy; optional local watch/hook (no always-on daemon) |
| G-004a vector | **Forbidden** — permanent defer |
| Law 6–7 | Progressive caps; no full-graph dump default |
| Law 19 | Library first; HTTP/MCP/GUI adapters thin |

## References

- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md) — Laws 6–7, 19
- [Phase 41 VERIFY-NOTES](../phase-41-layers-intent/scopes/scope-02-verify/VERIFY-NOTES.md)
- [REMEDIATION-PLAN §2 G6/G7](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-004b/G-005](../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [Phase 41 DR-HANDOFF](../phase-41-layers-intent/DR-HANDOFF.md)

## Planner gate (P42-00)

- [x] Re-read live repo: `internal/retrieval/`, indexer/analyzer paths, lang coverage
- [x] Thicken S00–S02 `00-PLANNER` + `01-*` + `02-*` prompts with file targets
- [x] Board `docs/TODO/phase-42.md` — scope rows filled after P42-00
- [x] Resolve INTAKE open questions (G6 law review gate; G7 lang policy cut)
- [x] DR-HANDOFF scope checklist accurate
- [x] No product code in this row (planner only)

## P42-00 locks (2026-08-22)

| Item | Locked value |
|------|--------------|
| G6 law review | Desk-check at S00-00 → `LAW-REVIEW-NOTES.md`; implement S00-01 |
| G6 channel | `graph_label_match` FTS over concept entity types; merge compile/explore |
| G7 langs | Tier-1: go/js/ts/tsx/py (frozen); Tier-2 defer; Tier-3 path-only |
| G7 freshness | Git-hook primary; optional foreground `trace index watch`; no daemon |
| Successor | **`no successor` default** at S02-02 (G1–G9 complete) |

## Exit criteria

- Scope stubs runnable; board points to **P42-S00-00**
- G6/G7 locked with accept/reject per REMEDIATION-PLAN
- Successor sketch for Phase 43+ documented in DR-HANDOFF forward note

## Next

`P42-S00-00`
