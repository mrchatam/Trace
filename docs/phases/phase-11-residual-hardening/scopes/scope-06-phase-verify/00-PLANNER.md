# P11-S06-00 — Phase 11 VERIFY (STUB)

## Metadata
- id: P11-S06-00
- todo_ids: [P11-S06-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: automated

## Objective
Lock Phase 11 VERIFY evidence table: **S01–S05 named DF regressions** + **carry-forward gates** + product `./...`. Decide **DR-HANDOFF** = **`no successor`** unless Notes promote. No product Go.

## References
- [phase README](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: [phase-10 VERIFY](../../../phase-10-integrity-surfaces/scopes/scope-05-phase-verify/)

## Depends-on (S01–S05 — when landed)
| Scope | DFs |
|-------|-----|
| S01 | DF-43, DF-51, DF-44 |
| S02 | DF-40 |
| S03 | DF-47 |
| S04 | DF-41, DF-49, DF-50 |
| S05 | DF-35, DF-48, DF-42 |

## Locked defaults (STUB — S06-00 FINALIZES)
- Carry-forward: honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false`; Phase 10 DF-17…32 spot checks
- Dry-run ≠ Gate C/F/G/ablation/H/checklist
- Spawn 01a/b/c on fail
- DR-HANDOFF default **`no successor`** (S06-01 starts, S06-02 owns)

## Exit
- [ ] Thicken 01-verify + 02-scope-review + SCOPE-TODOS; next **P11-S06-01**
- [ ] Product Go — **not** this row
