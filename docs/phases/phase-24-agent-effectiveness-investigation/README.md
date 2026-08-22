# Phase 24 — Agent effectiveness investigation (thought process not working in dogfood)

**Status:** scaffolded (2026-08-20, P24-00) — human-promoted after E01 + Phase 23 close. Design SoT: [`INVESTIGATION.md`](INVESTIGATION.md). Next runnable: **`P24-S01-00`**.

Phase 23 delivered **machine-checkable enforcement** (gate CLI, `--enforce`, harness install). Phase 24 asks why **live agents still do not behave as Trace intends** in **default build mode** — and what differs when humans **explicitly** ask for gap analysis with Trace.

## Why this phase exists

E01 (`ab-incident-tracker`) produced **two observable modes** in the same G1 workspace:

**Session A — initial build:** seed-task-only graph, 0 decisions, loop STOP/saturated, G1 code identical to B0.

**Session B — directed gap analysis** (human asked agent to gap-analyze + plan + fix with Trace): **7 discoveries**, **2 decisions**, **4 evidence**, product fixes (unassign, assignee filter, public `started_at`, store tests), G1 **diverges from B0** — but **still no new tasks**, **0 uncertainties**, verify task **still gate-blocked** (`hop_budget_exceeded`).

Enforcement **blocks** bad transitions when invoked; it does not **teach** agents to loop, expand tasks, or recover hop budget after gap work. This phase explains that gap and ranks fixes so **build mode** gains directed-mode behavior without a second human prompt.

## Architecture (investigation flow)

```text
E01 Session A (build) + Session B (directed gap) + prior dogfood
    ↓
S01 Two-mode post-mortem + failure taxonomy
    ↓
S02 Codebase audit (loop, SelectNext, task add, saturation, MCP/CLI UX)
    ↓
S03 External research (similar projects, harness patterns, papers/docs)
    ↓
S04 Intervention matrix (product / harness / experiment / docs — ranked)
    ↓
S05 VERIFY + DR-HANDOFF → Phase 25+ implementation queue (human promotes)
```

## Phase locks (P24-00)

| Item | Lock |
|------|------|
| Goal | Explain **why** Trace thought process fails in live agents; propose **evidence-backed** fixes |
| Method | Dogfood post-mortem + codebase trace + external research + synthesis |
| Deliverables | `FINDINGS.md`, per-scope notes, intervention matrix, VERIFY evidence |
| Product Go | **Deferred** — S01–S04 are docs/research rows unless a row explicitly says "spike ≤N lines" |
| Forbidden | daemon; hosted MCP; rewriting Phase 23 history; claiming fixes without evidence |
| Successor | DR-HANDOFF lists **candidate Phase 25 themes** — human picks one |

## Scope order (locked)

| Scope | Focus |
|-------|--------|
| S01 | E01 + prior dogfood **failure taxonomy** |
| S02 | **Codebase audit** — loop, policy, task creation, saturation, install surfaces |
| S03 | **External research** — similar projects + harness/orchestration patterns |
| S04 | **Intervention matrix** — ranked product vs harness vs protocol changes |
| S05 | VERIFY + DR-HANDOFF |

## Out of scope unless promoted

- Implementing Phase 25 features inside P24 investigate rows
- New dogfood experiment runs (note recommendations only)
- LLM-as-judge for agent quality

## Completion bar

1. `FINDINGS.md` documents **Session A vs Session B** with cited artifacts.
2. `FINDINGS.md` lists ≥5 distinct failure modes with E01 evidence.
2. Codebase audit maps each failure mode to **specific files/APIs** (loop, deliberation, add, MCP).
3. External research cites ≥3 comparable systems or papers with **actionable** deltas.
4. Intervention matrix ranks ≥8 interventions with effort/impact/risk and owner (product/harness/docs).
5. DR-HANDOFF recommends 1–3 **Phase 25** implementation themes (not a flat backlog).
