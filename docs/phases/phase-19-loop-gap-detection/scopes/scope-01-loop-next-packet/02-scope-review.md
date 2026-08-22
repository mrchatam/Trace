# P19-S01-02 — Review `trace loop next`

## Metadata
- id: P19-S01-02
- todo_ids: [P19-S01-02]
- role: reviewer
- skills: [documentation-and-adrs, writing-for-agents]
- mcps: []
- verification: automated

## Objective

Independently review `trace loop next` for:
- packet usefulness to a fresh agent
- freshness correctness
- graph/context grounding quality
- CLI compatibility

This review owns forward fixes if the implementer had to leave structural gaps in Notes.

## Review focus

The reviewer should confirm the implementation actually reused the intended live seams, or justify any divergence:

- `tasks` for task summaries
- `why` for causal trace
- `context` for bounded packet context
- `plan show` for active plan state
- `impact walk` or equivalent bounded graph traversal for related files/symbols

The reviewer should also verify the planner locks from `P19-S01-00` were honored:

- CLI entry stays narrow: `trace loop next --task <id>`
- goal context is derived from the seed task's `goal_id`, not from a second required planning seed
- missing task, task with no `goal_id`, and missing goal-plan context fail closed with non-zero exit
- packet sections have explicit presence/freshness semantics rather than silent omission

Also verify the command did not quietly reintroduce forbidden behavior:

- no full-graph dump by default
- no stdin-interactive primary protocol
- no daemon / hosted dependency
- no silent success when freshness or packet sections are unavailable

## Evidence to gather

- inspect help/usage text and root wiring
- exercise success case for a task that has a goal and plan context
- exercise failure cases for missing task / missing `goal_id` / missing plan context where feasible in automated tests
- inspect packet shape for deterministic section naming and explicit freshness markers
- verify any related-neighborhood section is bounded and provenance-preserving

## Spawn guidance

If the review finds a structural packet-shape or freshness gap that is more than a trivial inline fix, spawn the next forward implement/review pair immediately below this row instead of loosening the lock or rewriting earlier history.

## Exit criteria

- blocker/high issues fixed or spawned forward
- review confidence medium/high with explicit residuals
