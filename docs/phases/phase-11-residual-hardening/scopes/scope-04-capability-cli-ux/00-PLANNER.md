# P11-S04-00 — Capability upsert + small CLI UX (STUB)

## Metadata
- id: P11-S04-00
- todo_ids: [P11-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: automated

## Objective
Finalize S04 for **DF-41, DF-49, DF-50** (+ thin DF-22/37 install tip). **No product Go.**

## References
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A6
- Phase 10 S01 Exact/Why; S02 install tip

## Owns
| DF | Intent |
|----|--------|
| DF-41 | `capability declare` upsert by slug without requiring `--id` |
| DF-49 | `trace why symbol <id>` works (Exact/Why) |
| DF-50 | Print-only `install cursor` emits reload tip (pairs DF-22/37 ops) |

## Locked defaults (STUB — S04-00 FINALIZES)
- Prefer no mig if Upsert-by-slug fits existing UNIQUE
- DF-22/37: **docs/tip only** — do not invent Cursor process control
- No new MCP tools beyond existing nine unless S04-00 explicitly adds thin parity for symbol why (prefer inherit Exact)

## Exit
- [ ] Thicken 01+02 + SCOPE-TODOS; next **P11-S04-01**
- [ ] Product Go — **not** this row
