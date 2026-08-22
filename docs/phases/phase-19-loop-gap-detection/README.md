# Phase 19 — Loop gap detection + harness-agnostic packets

**Status:** planner locked (2026-08-18) — human-promoted successor after Phase 18 close, then finalized by `P19-00`. Next runnable: **`P19-S01-00`**.

## Why this phase exists

Taskboard D42 showed a real Trace advantage in structured planning continuity, but also exposed the next product gap:

1. Trace can store goals/tasks/decisions/why/context, but it cannot yet drive a **repeatable loop** for “understand -> find gaps -> add tasks -> re-evaluate”.
2. `trace why` is useful, but today it mostly returns governance edges and recent events unless the graph already contains richer planning facts.
3. A fresh agent needs a **single context packet** that combines plan/why/context/impact/freshness in a harness-agnostic form.
4. Long-run portability matters: the loop surface must work across shell, CLI, MCP wrappers, IDEs, and offline harnesses without requiring an interactive stdin session.

## Goal

Ship a thin, local-first Trace MVP for iterative planning loops:

- `trace loop next` emits a structured packet on stdout
- `trace loop apply` records discoveries / plan changes / spawned tasks from agent output
- `trace loop status` reports whether the loop is saturated

This phase is about **planning-state orchestration**, not hosted services, not daemonization, and not turning Trace into a coding runtime.

## Scope order (locked)

| Scope | Focus |
|-------|-------|
| S00 / phase planner | Lock MVP and scaffold rows |
| S01 | `trace loop next` packet: task/why/plan/context/freshness/related-files |
| S02 | `trace loop apply` + `trace loop status`: structured writes + saturation stop |
| S03 | VERIFY: harness-agnostic CLI proof + taskboard continuation mini-eval + DR-HANDOFF |

## MVP locks

- Core interface is **stdout-first**, not stdin-first
- No daemon / hosted service / product MCP dependency
- Reuse existing Trace primitives where possible: `tasks`, `why`, `context`, `plan`, `impact`
- Reuse existing CLI/domain seams where possible:
  - `trace tasks` for task list / goal filter
  - `trace why` for causal chain + decision impact summaries
  - `trace context` for bounded task context packets
  - `trace plan show` and `trace plan apply-discovery` for plan state and controlled replan
  - `trace impact walk` for related file/symbol neighborhood
- Loop stop rule for MVP: **stop when no new tasks and no new plan changes are generated** (or max iterations)
- File-relation / blast-radius info should be included in `loop next` using existing graph/index data; no full-graph dumps by default
- Freshness must be explicit (`fresh` / `dirty` / `stale` / `unknown`) so loops do not silently act on stale context
- Apply path must stay **structured and fail-closed**: malformed loop output cannot create partial state silently

## Live inventory baseline

`P19-00` confirmed the current repo already has enough planner-facing primitives to support a thin loop surface without inventing a second planning system:

- `cmd/trace/tasks.go` emits JSON task rows (`id`, `title`, `work_state`, `goal_id`)
- `cmd/trace/why.go` emits `WhyResult` plus decision impact summaries
- `cmd/trace/context.go` emits bounded JSON/markdown task context with optional why inclusion
- `cmd/trace/plan.go` already exposes `show`, `apply-discovery`, `deep`, and current-scope state
- `cmd/trace/impact.go` already exposes bounded impacted-neighborhood walks for file/symbol seeds

So Phase 19 should layer a **loop-oriented envelope** over existing retrieval/planner logic, not fork product semantics.

## Out of scope unless promoted

- Interactive stdin-driven loop protocol as the primary core surface
- Hosted coordination, multi-user server state, daemon workers
- Full autonomous coding loop that self-runs tests and edits code
- Embeddings / vector DB / semantic web retrieval
- Rewriting Phase 18 history or claiming Phase 18 handoff was wrong

## Borrowed patterns (research input)

Research candidate: [`similar projects/Understand-Anything/`](../../../similar%20projects/Understand-Anything/)

Stealable ideas for this phase:
- search -> 1-hop context expansion
- diff ripple / impacted neighborhood
- freshness classification before using graph outputs
- deterministic packet building for repeatable agent loops

## Completion bar

The phase passes when a fresh agent can:

1. call `trace loop next --task <id>`
2. receive a machine-readable packet with enough context to reason about gaps
3. return structured changes through `trace loop apply`
4. call `trace loop status` and observe saturation when no further tasks/plan changes are produced

And this works through ordinary CLI invocation without requiring a special harness.
