# P19-S01-00 — Loop next packet planner

## Metadata
- id: P19-S01-00
- todo_ids: [P19-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective

Replan against live CLI/domain/planner/retrieval surfaces and lock the MVP for `trace loop next`. Thicken sibling prompts. **No product Go in this row.**

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [../../00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- Live seams: `cmd/trace/tasks.go`, `cmd/trace/why.go`, `cmd/trace/context.go`, `cmd/trace/plan.go`, `cmd/trace/impact.go`, `cmd/trace/root.go`
- Research inputs from `similar projects/Understand-Anything`

## Locked defaults

`trace loop next` must emit a single machine-readable packet containing:
- seed task / goal
- relevant tasks summary
- why chain
- plan snapshot
- context snapshot
- freshness state
- related files / symbols / impacted neighborhood
- loop hints (max_iterations, current_iteration if available)

No stdin wait. No agent-side reasoning inside Trace.

## Live inventory to reuse

`P19-00` confirmed these reusable surfaces already exist:

- `trace tasks [--goal]` for bounded task summaries
- `trace why <type> <id>` for causal chain + attached impact summaries
- `trace context <task-id> [--depth 1|2] [--include-why]` for bounded task context packets
- `trace plan show --goal <id>` for current scope, deep plan, lookahead, and goal tasks
- `trace impact walk --seed file:<id>|symbol:<id>` for related neighborhood without a full graph dump

S01 should compose these into a loop packet and only add new product code where composition glue is missing.

## Packet locks

The planner for this scope must lock an output envelope that stays small, deterministic, and wrapper-friendly. At minimum the implementer prompt must name:

- top-level identifiers: goal, seed task, iteration metadata
- bounded task summary list, not full task bodies for the whole project
- a plan snapshot derived from `plan show`
- a context snapshot derived from `context`
- a why snapshot derived from `why`
- related file/symbol neighborhood derived from `impact walk` or equivalent graph access
- explicit freshness classification for every reused packet section

Freshness vocabulary for MVP remains: `fresh`, `dirty`, `stale`, `unknown`.

## Planner work

1. [ ] Inventory existing APIs that can be reused.
2. [ ] Lock output schema shape and freshness semantics.
3. [ ] Decide seed/goal flag set and fail-closed behavior for missing task or missing plan context.
4. [ ] Thicken `01-loop-next-packet.md` and `02-scope-review.md`.

## Exit criteria

- [ ] FINAL prompt for S01 implement/review
- [ ] Packet schema shape locked
- [ ] Reused live seams named so S01-01 does not re-debate foundations
- [ ] No product Go
