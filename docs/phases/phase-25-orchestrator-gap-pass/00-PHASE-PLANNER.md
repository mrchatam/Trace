# P25-00 — Phase 25 planner (orchestrator + default gap pass)

## Metadata
- id: P25-00
- todo_ids: [P25-00]
- role: planner
- verification: automated

## Objective

Lock Phase 25 scopes against live repo state and Phase 24 handoff so implementers can execute `INT-03`, `INT-04`, and `INT-11` without re-debating scope.

## References

- [docs/rules/agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../rules/project-rules.md)
- [README.md](README.md)
- [GAP-PASS.md](GAP-PASS.md)
- [Phase 24 DR-HANDOFF](../phase-24-agent-effectiveness-investigation/DR-HANDOFF.md)
- [Phase 24 matrix](../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify if needed → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Theme | P25-C only |
| Top interventions | INT-03, INT-04, INT-11 |
| Core boundaries | Local-first; no daemon/HTTP core path |
| Product Go on this row | No |

## Exit criteria

- [ ] Scope list and run order verified against live repo
- [ ] Scope stubs are runnable and references are valid
- [ ] Board points to next runnable scope planner (`P25-S01-00`)

## Next

`P25-S01-00`
