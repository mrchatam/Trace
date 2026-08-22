# P28-S07-00 — Scope planner (residual-wave VERIFY)

## Metadata
- id: P28-S07-00
- todo_ids: [P28-S07-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: automated
- hooks: []

## Objective

Plan residual-wave VERIFY after S06 FR-P28-01…07 reviews APPROVE. Lock verify floor + DR-HANDOFF residual-wave close policy. Do **not** rewrite S05 close history — extend `DR-HANDOFF.md` Residual wave section only.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [VERIFY-NOTES.md](../scope-05-verify/VERIFY-NOTES.md) — S05 baseline (immutable history)
- [SCOPE-TODOS.md](../scope-06-r6-fm-residuals/SCOPE-TODOS.md)
- S05 pattern: `../scope-05-verify/00-PLANNER.md`

## Session start

Follow agent-loop-protocol Session start.

## Locked defaults

| Item | Value |
|------|-------|
| Prerequisite | P28-S06-02…14 all `done` with APPROVE (or explicit skip with reason) |
| Verify floor | `go test ./internal/... ./cmd/trace/...`; install/hook smoke; spot-check FR evidence artifacts |
| Dual-lane | No `prepare.sh`; do not conflate thin/rich arms |
| DR-HANDOFF | Close **Residual wave** section only; keep S05 CLOSED history intact |
| Successor | Explicit: `no successor` **or** human-promoted Phase 29 — never TBD |

## Planner gate

- [x] `01-verify.md` + `02-dr-handoff.md` thickened for residual wave
- [x] Board next = **P28-S07-01**

## Exit criteria

- [x] Sibling prompts runnable
- [x] Notes cite verify floor locks

## Todo updates

Status + notes on **P28-S07-00** only.

## Next

`P28-S07-01`
