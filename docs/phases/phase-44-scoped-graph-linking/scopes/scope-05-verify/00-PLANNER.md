# Phase 44 / S05 / VERIFY planner

## Metadata
- id: P44-S05-00
- todo_ids: [P44-S05-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [user-trace]
- verification: automated

## Objective

Thicken VERIFY blocks: M-001, Laws 6–7, Law 19, portable graph, auth FE/BE fixture smoke. Prepare DR-HANDOFF close path (successor default `no successor`). **No product features.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) §10
- [`DR-HANDOFF.md`](../../DR-HANDOFF.md)
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)

## Session start

Follow agent-loop-protocol. Requires S00–S04 reviews APPROVE (or documented residuals).

## Locked defaults

| Item | Value |
|------|-------|
| Evidence | `VERIFY-NOTES.md` in this folder |
| Successor | Default `no successor` unless human names next phase |
| Portable | Require `trace/graph.json` export evidence if schema changed |

## Role work

1. Enumerate VERIFY blocks in `01-verify.md` with commands/paths.
2. Thicken `02-dr-handoff.md` checklist (AGENTS.md / TODO.md orchestrator update).

## Exit criteria

- [ ] `01-verify.md` runnable; next **P44-S05-01**

## Minimal todos

- [ ] Thicken 01 + 02
- [ ] Board done

## Next

`P44-S05-01`
