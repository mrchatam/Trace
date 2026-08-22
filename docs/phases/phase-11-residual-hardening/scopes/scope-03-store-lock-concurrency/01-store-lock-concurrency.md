# P11-S03-01 — Store lock / concurrency (STUB)

## Metadata
- id: P11-S03-01
- todo_ids: [P11-S03-01]
- role: implementer
- verification: automated

## Objective
Implement DF-47 per **P11-S03-00** FINAL locks.

## Depends-on
- P11-S03-00 done

## Exit criteria (outline)
- [ ] Named test(s) for parallel open / retry / UX per locks
- [ ] `TestConcurrentStoreOpenFailClosed` (or successor) still honest
- [ ] Compat checklist + CGO0 store/mcp + CGO1 product `./...` PASS
- [ ] No multi-writer SQLite; no daemon
