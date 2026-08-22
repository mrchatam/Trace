# P28-S06-02 — FM-01 / FR-P28-01 reviewer

## Metadata
- id: P28-S06-02
- todo_ids: [P28-S06-02]
- role: reviewer
- skills: [code-review-and-quality]
- verification: mixed
- hooks: []

## Objective

Independent **fresh-session** review of P28-S06-01 (FR-P28-01 / FM-01). Re-run spot-checks; do not trust Notes alone. Approve, request forward spawn, or fail with evidence.

## References

- [01-fm01-implement.md](01-fm01-implement.md)
- [00-PLANNER.md](00-PLANNER.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md)
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)

## Session start

Fresh context — must not be S06-01 implementer.

## Checklist

- [ ] Acceptance: post-import BLOCKING path → task or explicit decline without invented UUIDs
- [ ] Human gate preserved (no silent auto-spawn)
- [ ] Tests/dogfood evidence present; product diffs justified
- [ ] No daemon/HTTP; no scope creep into other FMs
- [ ] Forward-only (no rewrite of S00–S05 done history)

## Exit criteria

- [ ] Verdict in Notes: APPROVE / spawn / fail
- [ ] Next runnable **P28-S06-03** on APPROVE

## Todo updates

Status + notes on **P28-S06-02** only; may spawn forward rows if blocked.

## Next

`P28-S06-03`
