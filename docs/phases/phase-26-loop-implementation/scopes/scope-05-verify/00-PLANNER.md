# P26-S05-00 — Scope planner (VERIFY)

## Metadata
- id: P26-S05-00
- todo_ids: [P26-S05-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Define VERIFY / DR-HANDOFF criteria after S02–S04 reviews are `done`. No product features on this row.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [Phase 26 README](../../README.md)
- `experiments/ab-p25-gap-pass-validation/`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Precondition | P26-S02-02, P26-S03-02, P26-S04-02 all `done` |
| Preferred verify | Rebuild `bin/trace`; smoke install + `score.sh G1 --p25` |
| Closure signal | P25-2 PASS (was FAIL on E02) |
| Successor | Decided only on P26-S05-02 |

## Verify strategy

- **A (preferred):** rebuild binary; `prepare.sh G1` + `score.sh G1 --p25`
- **B (partial):** temp-dir `trace install cursor --write` greps if E02 workspace stale

## Success criteria

| Check | Target |
|-------|--------|
| P25-1 GapPassPrompt | PASS |
| P25-2 Parent orchestrator | **PASS** |
| P25-3 graph richness | PASS (or document residual) |
| `go test ./internal/...` | PASS |

## Planner gate

- [ ] `01-verify.md` + `02-dr-handoff.md` runnable
- [ ] `SCOPE-TODOS.md` current

## Exit criteria

- [ ] VERIFY implementer knows Option A vs B; own Notes updated
