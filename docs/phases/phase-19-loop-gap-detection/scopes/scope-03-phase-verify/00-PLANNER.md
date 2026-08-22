# P19-S03-00 — Phase 19 verify planner

## Metadata
- id: P19-S03-00
- todo_ids: [P19-S03-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Grep, Glob, Write]
- verification: automated

## Objective

Lock the verification bar for the loop MVP and own DR-HANDOFF closeout.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- D42 evidence in `experiments/RESULTS.md`

## Verify focus

- plain CLI proof for `loop next/apply/status`
- mini-eval on an existing repo/task (likely taskboard continuation or fixture-scale proof)
- harness-agnostic output
- no daemon, no special REPL assumptions

## Session start

Follow agent-loop-protocol session start (Agent -> clarify if needed -> Plan -> execute). This row is approved for unattended planner execution: perform internal clarify/plan and execute without stopping at planning.

## Locked defaults (S03)

| Item | Value |
|------|-------|
| Verify command surface | `trace loop next --task <id>`, `trace loop apply [--in <path>]`, `trace loop status --task <id> [--goal <id>]` |
| Expected schema versions | `trace.loop.next.v1`, `trace.loop.apply.v1`, `trace.loop.status.v1` |
| Status reasons allowed | `tasks_and_plan_unchanged`, `max_iterations_reached`, `insufficient_history` |
| Apply replay expectation | same `apply_id` + seed returns replay no-op and does not add duplicate tasks/plan changes/links |
| Evidence transport | ordinary CLI commands from repo root (or explicit `-C`) with stdout JSON evidence |
| Out-of-scope | no product Go edits in planner row; no daemon/MCP dependency for verify |

## Live repo anchors (must match proof bar)

- CLI entrypoints/help in `cmd/trace/loop.go`:
  - `loop next --task <id>` emits bounded JSON packet (stdout)
  - `loop apply [--in <path>]` reads file or stdin, validates fail-closed
  - `loop status --task <id> [--goal <id>]` derives saturation from persisted loop-step events
- Packet schema anchor in `internal/loop/next.go`: `trace.loop.next.v1` with explicit section freshness and related-neighborhood availability.
- Apply/status schema anchors in `internal/loop/apply.go`:
  - apply envelope `trace.loop.apply.v1` requires `schema_version`, `apply_id`, `seed`, `writes`
  - status output `trace.loop.status.v1` and reason set above
  - append-only loop-step event type `loop.step.applied` as status evidence source
- Existing CLI tests in `cmd/trace/loop_test.go` already cover packet shape, fail-closed apply validation, replay/idempotency, and status reasons.

## Required verify packet/examples for S03-01

S03-01 must capture one evidence packet each for:

1. `loop next` JSON packet example proving machine-readable bounded packet with explicit schema and freshness/availability sections.
2. `loop apply` JSON writeback example proving structured apply result (`new_*` counters, `replay`).
3. `loop status` JSON saturation example proving persisted-state reasoning (`insufficient_history` and at least one saturated reason path).

The evidence must be ordinary CLI invocation and must not rely on internal in-memory state or test-only helper code.

## Mini-eval bar (S03-01)

Mini-eval must prove a fresh agent can consume `loop next` output and produce at least one forward write (`spawned_task` or `plan_change`) through `loop apply`, then observe status progression via `loop status`.

Acceptable evidence paths:

- Preferred: continuation on a real local Trace goal/task that already has planning context.
- Allowed fallback: fixture-scale repo/task synthesized for controlled proof.

If fallback is used, S03-01 Notes must include an explicit residual-risk statement that live continuation variance remains to be validated.

## DR-HANDOFF ownership contract

- S03-00 (this row) locks close expectations and review criteria only.
- S03-01 collects verify evidence; it must not close handoff.
- S03-02 must make the explicit handoff successor decision in `DR-HANDOFF.md`:
  - close with successor queue identified, or
  - close with explicit `no successor`.

No silent `TBD` handoff state is allowed after S03-02 completion.

## Locked verification intent

The verify planner must leave a deterministic bar for S03-01:

- one ordinary CLI packet example from `loop next`
- one ordinary CLI writeback example from `loop apply`
- one ordinary CLI saturation example from `loop status`
- a mini-eval showing a fresh agent can identify at least one gap and create at least one forward task or plan change from the packet
- explicit residual-risk statement if the mini-eval uses a fixture instead of a live continuation of D42

S03 also owns the successor decision in `DR-HANDOFF.md`; it must not silently leave the phase open-ended.

## Exit criteria

- [ ] verify commands locked
- [ ] handoff policy stated
- [ ] mini-eval bar and acceptable evidence path locked
- [ ] no product Go in this planner row

## Minimal todos

- [ ] Lock S03-01 CLI proof packet expectations to live loop schemas/flags
- [ ] Lock S03-01 mini-eval acceptable evidence paths + residual-risk wording
- [ ] Lock S03-02 handoff-close acceptance criteria and successor decision requirement
- [ ] Update `DR-HANDOFF.md` close-policy language for S03 reviewer execution
- [ ] Mark board row status/notes with next runnable row
