# P11-S02-00 — Index partial-path GC (STUB)

## Metadata
- id: P11-S02-00
- todo_ids: [P11-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: automated

## Objective
Finalize S02 for **DF-40** (path-scoped `index <new-path>` rename ghosts). **No product Go.**

## References
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A4
- Phase 10 S03: [index-gc](../../../phase-10-integrity-surfaces/scopes/scope-03-index-gc/) — DF-20 full-tree GC must stay green

## Owns
| DF | Intent |
|----|--------|
| DF-40 | After rename, indexing only the new path must not leave ghost old path/symbols (or document + implement equivalent GC) |

## Locked defaults (STUB — S02-00 FINALIZES)
- Prefer incremental delete-on-missing / rename awareness; **no** full-rebuild-on-any-change
- Full-tree `trace index` GC (DF-20) must remain correct
- Prefer **no** mig; no MCP index tool

## Exit
- [ ] Thicken 01+02 + SCOPE-TODOS; next **P11-S02-01**
- [ ] Product Go — **not** this row
