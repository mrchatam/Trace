# P19-S01-01 — Implement `trace loop next`

## Metadata
- id: P19-S01-01
- todo_ids: [P19-S01-01]
- role: implementer
- skills: [writing-for-agents]
- mcps: []
- verification: automated

## Objective

Implement `trace loop next` as a stdout-first packet builder using existing Trace primitives where possible.

## Session start

Follow agent-loop-protocol. Unattended: do not stop after planning. No daemon, no hosted service, no stdin-driven primary interaction. This row may add product Go and tests, but board edits remain status + notes only.

## Locked defaults

| Item | Value |
|------|-------|
| Core shape | single JSON packet on stdout |
| Primary seed | task-centric MVP with `--task <id>` required |
| Goal handling | derive goal context from the seed task's `goal_id`; do not add a parallel required `--goal` flag in S01 |
| Reuse first | build on `tasks`, `why`, `context`, `plan show`, and `impact walk` semantics before inventing new retrieval paths |
| Boundedness | no full-project graph dump; packet stays targeted to the seed task / goal |
| Freshness enum | `fresh` \| `dirty` \| `stale` \| `unknown` |
| Failure mode | malformed args, missing seed task, seed task without `goal_id`, missing goal plan context, or impossible packet construction must fail closed with non-zero exit |

## Required CLI contract

The implementer should keep the first CLI surface intentionally narrow:

- `trace loop next --task <id>`
- no stdin input
- no hidden interactive fallback
- no alternate seed kind in this scope unless already required by the existing library seams

The command should reuse the seed task as the routing anchor:

1. load the task row
2. derive `goal_id` from that task
3. use the derived goal to obtain bounded task summaries and plan state
4. derive any related file/symbol seeds from already available task/why/context data before inventing new discovery flows

## Required packet sections

At minimum the packet should expose:

- packet metadata / schema version if needed for deterministic wrappers
- seed identity: task id, title, work_state, and non-optional goal id once the packet is emitted
- task summary slice for the relevant goal / nearby work, reusing the existing `tasks` row shape where practical
- plan snapshot from current planner state, reusing `trace plan show --goal <id>` semantics where practical
- why snapshot for the seed task (including impact summaries when available)
- context snapshot for the seed task, reusing `trace context <task-id> --depth 1 --include-why` or equivalent library behavior where practical
- related neighborhood focused on files/symbols/impact, not raw global graph output; if no usable seeds are discoverable, return an explicit empty/unknown section instead of fabricating one
- freshness metadata for every major section, or an equivalent explicit section-level status map
- loop hints such as iteration count / max iteration if such data is already available without inventing a separate runtime; otherwise emit explicit absence rather than synthetic counters

If a section is unavailable, the packet should say so explicitly rather than silently omitting provenance-critical context.

## Freshness expectations

Lock the vocabulary at section granularity:

- `fresh`: derived from currently available repo / planner state in this invocation
- `dirty`: underlying state appears to have changed enough that the section should be treated as suspect but still renderable
- `stale`: previously derived state is known to be outdated
- `unknown`: freshness cannot yet be proven

S01 does not need to solve perfect invalidation. It does need to make every section's freshness explicit and deterministic.

## Implementation guidance

- CLI surface under `trace loop next`
- machine-readable output (JSON)
- accepts exactly the seed task id plus only minimal support flags justified by the current command architecture
- includes enough data for a fresh agent to reason about gaps without rereading the whole repo
- includes explicit freshness classification
- includes related file/symbol context without dumping the entire graph
- no hosted services, no daemon, no stdin interaction

Prefer a thin adapter shape similar to other `cmd/trace/*` commands:

- wire command in `cmd/trace/root.go` and help text
- compose existing retrieval/planner services first
- keep loop-specific policy in a narrow library/helper instead of scattering it across unrelated commands
- prefer existing store/planner/retrieval/compiler/domain seams over shelling out to sibling CLI commands

Files likely touched:

- `cmd/trace/root.go`
- new `cmd/trace/loop*.go` adapter file(s)
- a narrow helper under `internal/` if composition glue is needed
- focused test files near the new command/helper

## Verification expectations

- focused tests for packet shape, missing-task behavior, missing-`goal_id` fail-closed behavior, missing-plan-context fail-closed behavior, and freshness / related-neighborhood presence
- verify existing commands remain backward compatible
- avoid brittle golden output that over-specifies incidental ordering unless ordering is part of the lock

## Exit criteria

- command works from normal CLI invocation
- tests cover packet shape and fail-closed behavior
- existing commands remain backward compatible
