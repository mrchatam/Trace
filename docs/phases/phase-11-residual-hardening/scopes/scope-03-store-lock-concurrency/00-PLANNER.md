# P11-S03-00 — Store lock / concurrency (STUB)

## Metadata
- id: P11-S03-00
- todo_ids: [P11-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: automated

## Objective
Finalize S03 for **DF-47** (exclusive `.trace` lock under parallel MCP/CLI). **No product Go.**

## References
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A5
- Phase 08 S02: [worktrees](../../../phase-08-ecosystem-hardening/scopes/scope-02-worktrees/) — path-local + exclusive lock

## Owns
| DF | Intent |
|----|--------|
| DF-47 | Parallel CLI↔CLI / CLI↔MCP on same root: clearer fail / retry / documented single-writer — without inventing multi-writer SQLite |

## Locked defaults (STUB — S03-00 FINALIZES)
- Keep path-local `.trace` bind
- Prefer short retry + clearer `ErrLocked` UX over dropping exclusivity
- Compat checklist lock tests must stay meaningful
- No daemon; no HTTP primary

## Exit
- [ ] Thicken 01+02 + SCOPE-TODOS; next **P11-S03-01**
- [ ] Product Go — **not** this row
