# P31-S02-00 — Scope planner (VERIFY)

## Metadata
- id: P31-S02-00
- todo_ids: [P31-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock VERIFY floor and DR-HANDOFF close policy for Phase 31. **No product code.** Thicken `01-verify.md` + `02-dr-handoff.md` if still thin. DR-HANDOFF stays **OPEN** until S02-02.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [Phase 31 README](../../README.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [GAPS.md](../scope-00-inventory/GAPS.md)
- S01 review Notes / tests

## Session start

Follow agent-loop-protocol. Gate: S01-02 PASS (or remediation closed).

## Locked defaults

| Item | Value |
|------|-------|
| VERIFY owner | P31-S02-01 |
| Close owner | P31-S02-02 only |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p31-s02-01-verify/evidence/` |
| Floor | Focused stray tests + `go test ./internal/...` + repro script if shipped |
| Fail vs residual | Regression / path change / silent delete / missing must-add → FAIL; deferred GAPS with reason → residual OK |
| Successor | **Phase 32** / first runnable **P32-00** (never TBD) |
| Product code | None in S02 |

## VERIFY floor (lock into 01)

1. Preflight + evidence dir
2. Focused store stray tests (existing + any new names from S01)
3. `go test ./internal/...`
4. Repro script if present; else temp one-shot matching P30 block-3 shape
5. Docs/gitignore/join spot-check (no path redesign)
6. Residuals listed; overall PASS/FAIL

## Planner gate

- [x] `01-verify.md` has commands + fail criteria + VERIFY-NOTES template
- [x] `02-dr-handoff.md` has successor table Phase 32 + spot-check floor
- [x] `SCOPE-TODOS.md` updated
- [x] DR-HANDOFF remains OPEN

## Exit criteria

- [x] Verifier can run unattended
- [x] Board Notes; next **P31-S02-01**

## Todo updates

Status + notes on **P31-S02-00** only.

## Next

`P31-S02-01`
