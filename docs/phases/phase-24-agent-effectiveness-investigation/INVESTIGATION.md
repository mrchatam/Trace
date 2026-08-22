# Investigation — agent effectiveness (Phase 24 SoT)

**Status:** locked for Phase 24 scaffold (2026-08-20). Investigators follow scope prompts; this doc is the product SoT for *what* we are diagnosing and *what counts as done*.

## Problem statement

Trace’s **thought process** (P19–P22: loop next/apply, gap detection, progressive plan, deliberation phases, blocking uncertainties, verification debt) is implemented in library/CLI and partially enforceable (P23). In **live Cursor dogfood**, behavior splits by **session mode** (see below).

### Session modes (E01 evidence — must not conflate)

| Mode | Trigger | Trace graph | Product | Loop/gate |
|------|---------|-------------|---------|-----------|
| **A — Build** | PROMPT-G1 / Multitask implement | Thin; seed tasks only; 0 decisions | Often ≡ B0 | STOP / saturated early |
| **B — Directed gap** | Human: “gap analysis + plan + fix using Trace” | 7 discoveries, 2 decisions, 4 evidence | Diverges from B0; real fixes | Still blocked (`hop_budget_exceeded`) |

**Phase 24 goal:** explain why Mode A ≠ Mode B and what product/harness changes collapse them.

In **Mode A**, agents still:

1. **Treat seed tasks as the full backlog** — rarely `trace add` or promote discoveries → new tasks
2. **Skip rich graph writes** — few decisions, unresolved uncertainties, post-hoc export
3. **Implement through gates** — code ships while loop reports STOP / saturation / hop budget exceeded
4. **Bypass orchestrator Trace** — Multitask parent delegates; only workers touch `.trace/` (when they touch it at all)
5. **Converge with B0** — identical product trees; Trace does not change implementation path

In **Mode B**, agents **can** record discoveries and decisions and drive fixes — but **still** do not create new tasks, resolve uncertainties, or clear loop/gate state for verify work.

**Working as intended** (operational definition for Phase 24):

| Behavior | Intended | E01 Session A | E01 Session B |
|----------|----------|---------------|---------------|
| Gap → new work | Discovery → new or updated tasks | Seed only | Discoveries yes; **tasks still seed-only** |
| Loop before edit | EXECUTE-ready before product edits | Saturated early | Gate still STOP/hop_budget |
| Graph honesty | Decisions + resolved uncertainties | 0 decisions | 2 decisions; 0 uncertainties |
| Orchestrator | Parent owns graph before delegate | Workers-only | Single-agent directed session |
| A/B signal | G1 differs in process + causality | Same code as B0 | Code diverges; graph richer |

## Investigation questions (must answer in FINDINGS.md)

### A — Agent / harness

- Why do agents not call `trace add` / MCP `trace_add` when work splits — **even after recording discoveries in Mode B**?
- Does **Mode B** only work because the human prompt names “gap analysis” — i.e. prompt-dependent, not product-default?
- Do prompts, seed shape, or MCP tool discovery hide “create task”?
- Does Multitask split break `TRACE_TASK_ID` and hook enforcement on the orchestrator?
- Are cursor rules too long / advisory vs hook-only on subagents?
- After gap fixes, why does **hop budget** still block verify task edits — is deliberation state never reset?

### B — Product / policy

- Does `SelectNext` recommend STOP/saturated before agents perceive “permission” to implement?
- Is `p19_saturated` / hop budget miscalibrated for single-session full-app builds **and** for post-build gap passes?
- Does **`discovery_mentions_task` without new tasks** give a false sense of replanning?
- Does export omit decisions agents recorded in SQLite (export bug vs agent skip)? — **Session B had 2 decisions in export; Session A had 0**
- Does loop **apply** fail to promote plan scopes → executable tasks?
- Should gap closure **transition** or **reset hop_count** on verify task?

### C — Experiment protocol

- Do seed tasks with fixed UUIDs **anchor** agents to a closed set?
- Does identical starter `project/` cause B0/G1 convergence regardless of Trace?
- Is arm isolation insufficient (sibling run readable)?

### D — External comparables

- How do other agent-memory / planning systems force replanning (Understand-Anything, SWE-agent, Aider conventions, etc.)?
- What harness patterns **require** plan expansion before file edit?

## Evidence sources (priority)

1. **E01 Session A + B** — `runs/G1/trace/graph.json`, git log (`704e2ff`…`a37e7c0`), VERIFY.md, diff vs B0
2. **E01 summary** — `experiments/RESULTS.md`
2. **Prior dogfood** — D44, D45 notes in RESULTS / conversation transcripts
3. **Repo** — `internal/loop/`, `internal/deliberation/`, `cmd/trace/loop.go`, `cmd/trace-mcp/`, `internal/install/`
4. **External** — web research + `similar projects/` (e.g. Understand-Anything)
5. **Phase docs** — P19–P23 VERIFY notes, ENFORCEMENT.md limits

## Failure mode taxonomy (seed list — S01 validates/extends)

| ID | Name | Symptom |
|----|------|---------|
| FM-01 | Seed anchoring | Only seed UUIDs used; no new tasks — **confirmed in Session B** |
| FM-02 | Graph thin export | decisions=0, no uncertainties — **partial: Mode B has decisions/evidence, no uncertainties** |
| FM-03 | Loop saturation | STOP / p19_saturated / hop_budget while coding continues — **persists after gap fixes** |
| FM-04 | Orchestrator bypass | Parent does not Trace; workers fragment graph |
| FM-05 | Enforcement optional | Hook/env not set; gates skipped |
| FM-06 | Cross-arm leakage | G1 code identical to B0 — **Session A yes; Session B no (G1 diverged)** |
| FM-07 | Post-hoc planning | SPEC after code; graph filled at export |
| FM-08 | Tool surface gap | MCP/CLI does not make “add task on gap” the obvious next step |
| FM-09 | Mode-dependent effectiveness | Trace works for gap **recording** only when human directs; not default build |
| FM-10 | Discovery without task promotion | `discovery` + `discovery_mentions_task` links but fixes land without `trace add` |

## Intervention categories (S04 must classify each)

| Category | Examples |
|----------|----------|
| **Product** | Auto-spawn tasks from discovery; reset hop budget on gap pass; relax saturation for greenfield |
| **Harness** | Default “gap pass” prompt in install; orchestrator hook; MCP `trace_add` nudge after discovery |
| **Protocol** | E01 scorer fix; FM-* checks; arm isolation; **two-session rubric** (build vs directed) |
| **UX/docs** | When to add task vs discovery-only; loop recovery after STOP |
| **Experiment** | Scoring that fails closed on FM-01..FM-10 |

## Phase 24 outputs

| Artifact | Owner scope |
|----------|-------------|
| `FINDINGS.md` | S01 opens; S04 consolidates |
| `scopes/scope-01-*/POSTMORTEM.md` | S01 |
| `scopes/scope-02-*/CODEBASE-AUDIT.md` | S02 |
| `scopes/scope-03-*/EXTERNAL-RESEARCH.md` | S03 |
| `scopes/scope-04-*/INTERVENTION-MATRIX.md` | S04 |
| VERIFY evidence dir | S05 |

## Non-goals

- Shipping Phase 25 code under P24 investigate rows
- Proving Trace “wins” A/B — only diagnosing why it did not
