# Phase 24 findings (consolidated)

Consolidated in **S04** (2026-08-20). S01 opens with E01 evidence (two sessions).

## Executive summary

Trace’s P19–P22 thought process is **implemented** in library/CLI and partially enforceable (Phase 23), but **live Cursor dogfood splits by session mode**. In E01 Session A (Multitask build), agents shipped a full app with a thin seed-only graph, early STOP/saturation, and product ≡ B0. In Session B (human-directed gap analysis), agents recorded 7 discoveries, 2 decisions, and real product fixes — yet still created **zero new tasks**, left all deliberation states STOP, and could not clear the verify gate. Root causes span **harness defaults** (no default gap pass, orchestrator bypass without `TRACE_TASK_ID`), **product policy** (P19 saturation on first empty apply, no deliberation reset, discovery without task promotion), and **protocol gaps** (two-session rubric, arm isolation, export/DB drift).

No single fix closes all ten failure modes. The ranked [INTERVENTION-MATRIX.md](scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) recommends Phase 25 start with **P25-C** (default gap pass + orchestrator Trace-first harness), then **P25-A** (discovery→task promotion) and **P25-B** (loop recalibration + deliberation reset). Measurement themes **P25-D/E** strengthen experiment scoring and export honesty but do not alone change agent-visible build behavior.

## Status

| Section | Owner | Status |
|---------|-------|--------|
| Two-mode model (build vs directed gap) | S01 | **done** — [POSTMORTEM.md](scopes/scope-01-dogfood-postmortem/POSTMORTEM.md) §1–§2 |
| Failure taxonomy | S01 | **done** — [POSTMORTEM.md](scopes/scope-01-dogfood-postmortem/POSTMORTEM.md) §3 |
| Codebase audit | S02 | **done** — [CODEBASE-AUDIT.md](scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md) |
| External research | S03 | **done** — [EXTERNAL-RESEARCH.md](scopes/scope-03-external-research/EXTERNAL-RESEARCH.md) |
| Intervention matrix | S04 | **done** — [INTERVENTION-MATRIX.md](scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) |

## Two-mode synthesis

Session A and Session B must remain **separate evidence rows** (see tables below). Synthesis: Trace is **mode-dependent** — directed gap mode records rich graph entities when the human names gap analysis; default build mode converges to seed anchoring, thin export, loop saturation, and B0 product parity **without** that prompt. Collapsing the modes requires harness defaults (gap pass, parent graph ownership) **and** product loop/task changes (promotion, reset, recalibrated saturation), not documentation alone.

## E01 — two agent sessions (same G1 workspace)

Phase 24 must treat these as **separate behaviors**, not one blended verdict.

### Session A — Initial build (Multitask / PROMPT-G1)

| Signal | Observation |
|--------|-------------|
| Tasks | Seed `…0010`–`…0050` only |
| Graph | 0 decisions; no discoveries; plan tree present |
| Loop | STOP / `p19_saturated` on seed tasks (export) |
| Product | Shipped; `internal/` **≡ B0** (scored — [RESULTS.md](../../../experiments/RESULTS.md)) |
| Harness | Multitask orchestrator; hook + strict config; verify gate blocked |
| Git | Anchor `f70aaea` — starter-only in git (4 files); product not at SHA |

### Session B — Directed gap analysis (human prompt: gap + plan + fix with Trace)

| Signal | Observation |
|--------|-------------|
| Discoveries | **7** gap-titled (`trace add` discovery entities) |
| Decisions | **2** (unassign semantics, monolith web) |
| Evidence | **4** (test pass, loop reviews) |
| Links | 2× `discovery_mentions_task` → verify task `…0050` |
| Tasks | **Still 5 seed only** — gaps fixed in place, no new tasks |
| Uncertainties | **0** |
| Loop / gate | Verify task STOP / `p19_saturated` in export; live gate `task_not_found` (DB sync) / `hop_budget_exceeded` (with `-C`) |
| Product | **Diverged from B0** — unassign, assignee filter, `started_at`, `store_test.go`, expanded web tests |
| Git | Commits `704e2ff` … `a37e7c0`; VERIFY manual table filled |
| Tests | `go test ./...` PASS (domain, store, web) |

**S01 explains:** Session B succeeds at gap **recording** but not task **expansion** or loop **recovery** — see [POSTMORTEM.md](scopes/scope-01-dogfood-postmortem/POSTMORTEM.md) Must answer §1–§3.

## Preliminary conclusion (full investigation)

Trace is **mode-dependent**:

- **Directed gap mode** — discoveries, decisions, evidence, real fixes **when human asks**
- **Default build mode** — seed anchoring, thin graph, saturation, B0 convergence **without explicit gap instruction**

Phase 25 should target **making build mode behave like directed gap mode by default** (harness gap pass + orchestrator Trace-first), then **sustain** that richness with task promotion and loop recovery (product). See [INTERVENTION-MATRIX.md §1](scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md#1-executive-summary--ranking-rationale) for ranked interventions and [DR-HANDOFF.md](DR-HANDOFF.md) for recommended Phase 25 themes.

## Failure taxonomy (FM-01..FM-10)

Validated in S01 — full matrix: [POSTMORTEM.md](scopes/scope-01-dogfood-postmortem/POSTMORTEM.md) §3. Per-session status:

| FM | Name | Session A | Session B |
|----|------|-----------|-----------|
| FM-01 | Seed anchoring | confirmed | confirmed |
| FM-02 | Graph thin export | confirmed | partial |
| FM-03 | Loop saturation | confirmed | confirmed |
| FM-04 | Orchestrator bypass | confirmed | rejected |
| FM-05 | Enforcement optional | partial | partial |
| FM-06 | Cross-arm leakage (G1≡B0) | confirmed | rejected |
| FM-07 | Post-hoc planning | confirmed | partial |
| FM-08 | Tool surface gap | confirmed | confirmed |
| FM-09 | Mode-dependent effectiveness | confirmed | partial |
| FM-10 | Discovery without task promotion | rejected | confirmed |

**Both sessions:** FM-01, FM-03, FM-08. **Session A only:** FM-04, FM-06. **Session B only:** FM-10.

See [INVESTIGATION.md](INVESTIGATION.md) for seed definitions and investigation questions A–D.

## Codebase audit (S02)

Full mechanism table: [CODEBASE-AUDIT.md §2](scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md#2-fm-mechanism-table-required).

- **All 10 FMs** mapped to file:line mechanisms; change levers tagged product \| harness \| protocol.
- **FM-03 root cause:** P19 saturation on first empty `loop apply`; sticky `Stopped` remaps live gate to `hop_budget_exceeded` while export retains `p19_saturated`.
- **FM-01 + FM-10:** Task roster fixed by seed import; only `loop apply` `spawned_tasks[]` expands backlog; mentions-task links do not create tasks.
- **FM-04/05:** Hook bypass without `TRACE_TASK_ID`; install rules gate-only — no gap-pass or task-promotion nudge.
- **No deliberation reset API** after gap pass; export/DB drift requires `-C` + import discipline.

## External comparables (S03)

Full table + Q-D answers: [EXTERNAL-RESEARCH.md](scopes/scope-03-external-research/EXTERNAL-RESEARCH.md).

- **Plan-before-edit (hard gates):** OpenHands planning agent; Aider architect→editor; Cursor `preToolUse` deny — Trace hook still allows parent without `TRACE_TASK_ID`.
- **Progressive graph context:** CodeGraph, UA, Graphify — scoped `trace context`; avoid full-graph dumps.
- **Task promotion:** AgentRQ `createTask`/`publishEvent` clearest peer; Trace FM-10 sits in middle (discoveries without spawn).
- **Deliberation reset vocabulary:** Graphiti episodes + literature memory evolution — Trace sticky STOP lacks peer-equivalent reset.
- **Anti-patterns flagged:** Graphiti/AgentRQ HTTP MCP as core path, full rebuilds, advisory-only rules.

## Intervention matrix (S04)

Ranked options: [INTERVENTION-MATRIX.md](scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) (11 rows INT-01..INT-11).

**Top 3:**

| Rank | ID | Intervention |
|------|-----|--------------|
| 1 | INT-03 | Default gap-pass install bundle (collapse Mode A→B without custom human prompt) |
| 2 | INT-04 | Orchestrator Trace-first — parent `TRACE_TASK_ID` + failClosed hook |
| 3 | INT-01 | Discovery→task promotion via `loop apply` spawn / guided `trace add task` |

**Recommended Phase 25 themes:** P25-C → P25-A → P25-B ([DR-HANDOFF.md](DR-HANDOFF.md)).
