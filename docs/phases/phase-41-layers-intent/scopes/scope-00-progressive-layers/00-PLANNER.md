# P41-S00-00 — Scope planner (G8 progressive layers)

## Metadata
- id: P41-S00-00
- todo_ids: [P41-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, context-engineering, api-and-interface-design]
- mcps: [user-trace, user-codegraph]
- verification: automated

## Objective

Lock S00 **G8** against live repo: progressive layers L2–L3 (G-003). Thicken `01-implement.md` + `02-review.md` with file targets, acceptance map, and rejects. **No product code in this row.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md) — P41-00 Q3 resolution (**ship**)
- [REMEDIATION-PLAN §2 G8](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-003](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [RETRIEVAL_AND_CONTEXT.md §4](../../../../RETRIEVAL_AND_CONTEXT.md) — L0–L3 definitions
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors (verified 2026-08-22 P41-00):
  - `internal/compiler/doc.go:7` — "Layers 2–3 are not auto-loaded in P0-X"
  - `internal/compiler/packet.go:44,97` — `Item.Layer` and `Packet.Layer` cap at 0–1 today
  - `internal/compiler/compiler.go:66–287` — `TaskContext` / `ExpandContext` depth 1..2; all items layer 0–1
  - `internal/compiler/compiler.go:376–391` — `layer1AdmitKey` only
  - `internal/compiler/budget.go:28–59` — trim prefers Layer 0 over Layer 1
  - `cmd/trace/context.go:18–77` — `--depth 1|2` (graph expand, **not** L2/L3 promotion)
  - `internal/mcp/tools_context.go:18–70` — MCP depth 1..2 mirror
  - Evidence: [h3-layers-designed-vs-shipped.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h3-layers-designed-vs-shipped.md)

## Session start

Follow agent-loop-protocol Session start. Unattended: INTAKE + P41-00 locks are authority.

## Locked defaults (FINAL — P41-00)

| Item | Value |
|------|-------|
| GAP ids | G-003 |
| Verdict | **Accept** per REMEDIATION-PLAN G8 |
| P41-00 decision | **Ship L2–L3 in compiler** (bounded); spec-revise only if S00-01 blocked |
| M-001 | Layer expansion merges into task packet — never full-graph dump or query-only |
| Law 6–7 | Progressive caps; L2–L3 **opt-in** — default packet stays L0–L1 |
| Layer policy | Per RETRIEVAL_AND_CONTEXT §4: deeper layers requested/justified, not auto-loaded |
| Depth vs layer | `--depth` = graph expand hops (1..2, unchanged); **`max_layer`** = progressive layer ceiling (default 1; opt-in 2 or 3) |
| L2 content | Direct dependents, recent discoveries, related future tasks, architectural neighbors |
| L3 content | Deeper historical decisions, cross-module impact, older evidence |
| Library first | Admission + compile in `internal/compiler/`; retrieval helpers in `internal/retrieval/` if needed |
| Adapters | Thin CLI/MCP `--max-layer` on existing `trace_context`; no new MCP tool required |
| Out | Auto-load L2–L3 on default context; semantic/vector; dump API; G-004a |

## Live repo gap (re-verified P41-00)

| Check | Shipped | Gap |
|-------|---------|-----|
| `packet.layer` max | 1 | Must reach 2 or 3 when opted in |
| `item.layer` values | 0, 1 only | Need 2, 3 with honest reason_codes |
| `--depth 2` effect | More L0–L1 hits via expand | Does **not** promote to L2/L3 |
| `doc.go` defer line | Explicit L2–L3 defer | Update when shipped |

## Accept / reject (G8)

| Decision | Item |
|----------|------|
| **Accept** | `ContextOptions.MaxLayer` (default 1; allowed 2, 3) on compiler path |
| **Accept** | `layer2AdmitKey` / `layer3AdmitKey` admission rules aligned to §4 |
| **Accept** | L2/L3-specific reason_codes (document in `retrieval/doc.go` if new) |
| **Accept** | Budget trim honors L0 > L1 > L2 > L3 priority |
| **Accept** | CLI `--max-layer 2|3` + MCP `max_layer` on `trace_context` |
| **Accept** | Tests G8-L1–L7 (see thickened `01-implement.md`) |
| **Accept (fallback)** | Spec-revise: update RETRIEVAL_AND_CONTEXT §4 + ADR if ship blocked — documented alternative only |
| **Reject** | Default `max_layer=3` on TaskContext |
| **Reject** | Reusing `--depth` as layer selector without separate `max_layer` |
| **Reject** | Full-graph dump; cap default inflation beyond scoped L2/L3 budget |

## Must lock for S00-01 (delivered in thickened 01-implement)

1. Touch-list: compiler admission → packet layer field → CLI/MCP opt-in → tests.
2. Seven acceptance tests G8-L1–L7.
3. Update `doc.go` defer honesty when L2–L3 ship.

## Exit criteria

- [ ] `01-implement.md` + `02-review.md` runnable with file targets + G8 accept map
- [ ] SCOPE-TODOS updated
- [ ] Board row → `done` with Notes

## Next

`P41-S00-01`
