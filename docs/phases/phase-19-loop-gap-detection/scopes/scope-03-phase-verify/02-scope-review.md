# P19-S03-02 — Review Phase 19 verify + close handoff

## Metadata
- id: P19-S03-02
- todo_ids: [P19-S03-02]
- role: reviewer
- skills: [documentation-and-adrs, writing-for-agents]
- mcps: []
- verification: automated

## Objective

Independently review Phase 19 verification evidence and either:
- approve and close DR-HANDOFF, or
- spawn forward fixes with full prompt stubs

## Review focus

The reviewer should confirm:

- all three loop commands were evidenced from ordinary CLI usage
- the mini-eval actually demonstrates loop usefulness, not just command existence
- saturation is derived from explicit persisted state
- residual risks are named if verification stayed fixture-scale
- `DR-HANDOFF.md` ends in an explicit successor state (`OPEN`, closed with successor, or closed with `no successor`)

Required acceptance checks:

- verify evidence references current schemas exactly:
  - `trace.loop.next.v1`
  - `trace.loop.apply.v1`
  - `trace.loop.status.v1`
- verify status reasons are constrained to `tasks_and_plan_unchanged`, `max_iterations_reached`, `insufficient_history`
- verify malformed apply rejection and replay no-op were both evidenced (not only claimed)
- verify mini-eval produced at least one forward write (`spawned_task` or `plan_change`) from packet-driven reasoning
- verify S03-01 did not bypass ordinary CLI with hidden harness-only behavior

## DR-HANDOFF close protocol (locked)

At row completion, reviewer must update `docs/phases/phase-19-loop-gap-detection/DR-HANDOFF.md` as follows:

1. If verify bar passes with acceptable residual risks, set status to closed and record explicit successor decision.
2. If verify bar fails, keep handoff OPEN and spawn forward fix rows directly below this review row.
3. Never leave `Successor decision` as `TBD` once row is marked done.

## Exit criteria

- no blocker/high gaps remain unaccounted for
- DR-HANDOFF state is explicit
- successor decision is explicit (named successor queue or `no successor`)
