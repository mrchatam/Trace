# Phase 44 / S03 / Planner — Inference

## Metadata
- id: P44-S03-00
- todo_ids: [P44-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, domain-modeling]
- mcps: [user-trace]
- verification: automated

## Objective

Plan opt-in inference rules (path, title, plan_scope) with `INFERRED` provenance per D4 (CLI only). Thicken `01-implement.md` with rule table + fail-closed conflicts. **No product code in this row.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) §4, D4
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)
- `internal/store/links.go` (`source_type`, `confidence`)
- `internal/store/plan_hierarchy.go` (`plan_scopes`)

## Session start

Follow agent-loop-protocol. Requires S02 APPROVE.

## Locked defaults

| Item | Value |
|------|-------|
| Trigger | CLI only (`trace graph infer` or locked equivalent) |
| Provenance | `source_type=INFERRED`; lower confidence than explicit |
| Conflicts | Explicit wins; never delete USER_ASSERTED / IMPORTED |
| Out | Index hook, daemon, LLM-only fill, gate-rel inference |
| DR-NOSSEM | No embeddings |

## Role work

1. Author rule table in `01-implement.md` (path glob → layer scope; title tokens; plan_scope co-membership).
2. Define idempotency + dry-run if cheap.
3. Thicken review rubric for honesty + caps.

## Exit criteria

- [x] `01-implement.md` with rule table + fail-closed conflicts
- [x] Next **P44-S03-01**

## Minimal todos

- [x] Thicken 01 + 02
- [x] Board done

## Next

`P44-S03-01`
