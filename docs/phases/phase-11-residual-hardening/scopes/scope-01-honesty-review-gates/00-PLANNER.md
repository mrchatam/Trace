# P11-S01-00 — Honesty / review-gate tightening (STUB)

## Metadata
- id: P11-S01-00
- todo_ids: [P11-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S01 implement/review prompts for **DF-43, DF-51, DF-44**. **No product Go in this row.**

## References
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A1–A3
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- Phase 10 S04: [operator-capability-gates](../../../phase-10-integrity-surfaces/scopes/scope-04-operator-capability-gates/)

## Owns (phase lock — thicken here)
| DF | Intent |
|----|--------|
| DF-43 | Sibling FAIL must block DONE even if another PASS exists |
| DF-51 | Hatch vs missing-caps: clear WARNING; flags stay independent by default |
| DF-44 | Flag≠identity clarity (help/MCP); no real auth |

## Locked defaults (STUB — S01-00 FINALIZES)
- Package hint: `internal/domain` (+ thin CLI/MCP warn copy); prefer **no** mig
- Keep Gate G hatch + honesty Path C / operator-flag supersession green
- Forbidden: OAuth; weakening DF-17/18/24; daemon/HTTP

## Exit
- [ ] Thicken `01-honesty-review-gates.md` + `02-scope-review.md` + SCOPE-TODOS
- [ ] Board Notes; next **P11-S01-01**
- [ ] Product Go — **not** this row
