# P19-S03-01 — Verify Phase 19

## Metadata
- id: P19-S03-01
- todo_ids: [P19-S03-01]
- role: verify
- skills: [writing-for-agents]
- mcps: []
- verification: automated

## Objective

Verify that the loop MVP works from ordinary CLI use and produces usable packets and saturation results.

## Session start

Follow agent-loop-protocol. Unattended: do not stop after planning. This row verifies and records evidence; it does not open a new product direction beyond the Phase 19 locks.

## Required evidence

Capture and archive evidence for all four artifacts below (stdout JSON + concise Notes summary):

- `trace loop next` packet example with `schema_version="trace.loop.next.v1"` and explicit freshness/availability sections (`seed`, `tasks`, `plan`, `why`, `context`, `related`, `loop_hints`)
- `trace loop apply` structured writeback example from `trace.loop.apply.v1` envelope (file input and/or stdin input)
- `trace loop status` progression examples including:
  - unsaturated/no-history path: `reason="insufficient_history"`
  - saturated path: `reason` in `tasks_and_plan_unchanged` or `max_iterations_reached`
- mini-eval showing a fresh agent can identify at least one gap and create at least one forward task or plan change from the packet

## Verification bar

- run via normal CLI invocation from the repo root or an explicit `-C` target
- show that packet output is machine-readable and bounded
- show that apply writes are reflected in repo state inspectable by existing Trace commands
- show that status can distinguish saturated from not saturated using stored state, not a live in-memory loop session
- prove fail-closed behavior remains intact (at least one malformed apply envelope rejection, no partial writes)
- prove replay-safe behavior (same `apply_id` + seed yields replay no-op with stable counts)
- keep evidence honest: if the mini-eval is fixture-scale rather than a live continuation, say so explicitly and record residual risk

## Locked command/test floor

S03-01 must run and report the result of:

- `go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'`

It may run broader tests too, but this floor must be present in Notes.

## Mini-eval expectations

- Use one seed task with a valid goal + current deep plan context (required by `loop next`).
- Feed the `loop next` packet to a fresh-agent reasoning step (manual or scripted) to produce one structured apply envelope.
- Apply the envelope via ordinary CLI, then verify the write surfaced in existing Trace views (`trace tasks`, `trace why`, `trace plan show`, or equivalent).
- Run `loop status` after apply and record whether it saturated and why.
- If fixture-scale is used, include a one-line residual risk describing what live-repo variance remains unproven.

## Failure guidance

- if any required command cannot produce plain CLI evidence, mark the row `failed` or `blocked` with exact reason
- do not close `DR-HANDOFF.md`; that belongs to S03-02 review
