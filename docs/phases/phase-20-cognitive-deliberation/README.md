# Phase 20 — Cognitive deliberation controller

**Status:** P20-S06-00 locked (2026-08-18) — human-promoted successor after Phase 19 close. Source vision: [`docs/TRACE_THOUGHTPROCESS.md`](../../TRACE_THOUGHTPROCESS.md). Coverage SoT: [`COVERAGE.md`](COVERAGE.md). Next runnable: **`P20-S06-01`**.

Phase 19 historical DR-HANDOFF remains **`no successor`** at time-of-close. This phase is a **forward** queue.

## Why this phase exists

Phase 19 shipped a harness-agnostic loop envelope (`trace loop next` / `apply` / `status`) that can saturate on task/plan churn. That is **not** yet project cognition.

[`TRACE_THOUGHTPROCESS.md`](../../TRACE_THOUGHTPROCESS.md) asks Trace to become an **externalized project cognition and engineering feedback system**: understand → gap → investigate → explore → decide → plan → critique → implement → test → verify → measure → score → detect regressions → reflect → update knowledge → replan.

Phase 20 turns the P19 loop into a **state-driven deliberation controller** and adds the missing durable artifacts (uncertainty, change/effects, verify≠test, gates, regression, reflection) **without** storing raw chain-of-thought and **without** a second database.

## Challenge (doc §30)

Do **not** implement every noun in §29B.

Risks if we copy the prompt literally:

- token/cost explosion from huge packets
- infinite ORIENT↔INVESTIGATE loops
- LLM-claimed “verified” without machine evidence (violates Law 2)
- false causal “this file caused latency”
- scoring theater (meaningless 0.87s)
- graph complexity that `why` cannot explain

Simplifications locked:

1. Deterministic policy table for phase selection — no ML.
2. Merge overlapping nouns (see COVERAGE.md).
3. Test/Verify/Evaluate are **result kinds**, not three product test runners in MVP.
4. Git remains canonical for diffs (Law 1). Change objects store SHA + path + expected/actual **semantic** effects.
5. P19 loop is the execution mechanism (doc §28). Do not replace `internal/loop`.
6. Agent output is untrusted. Promotion requires Claim→Evidence→Review or equivalent machine evidence (Laws 2, 14, 15).
7. Budgets: max phase hops per `loop apply`, max packet bytes, existing plan-churn N=5.

## Architecture (doc §29A)

```text
Agent (any harness)
    ↓  stdout JSON
trace loop next   →  deliberation packet (phase + context slice)
    ↓
Agent returns structured artifacts
    ↓
trace loop apply  →  writes artifacts, advances deliberation_state
    ↓
controller        →  next phase from deterministic policies
    ↓
trace loop status →  saturated | blocked | needs_<phase>
```

Components:

| Component | Role | Lives in |
|-----------|------|----------|
| `internal/loop` | P19 packet + apply/status | **reuse / extend** |
| `internal/deliberation` | phase enum, state, policy SelectNext | **new library** |
| domain/store | new thin tables + links | **extend SQLite** |
| CLI `trace loop …` | G19 adapter | **extend cmd/trace/loop.go** |
| MCP | inherit via library or thin later | **no new hosted tools required this phase** |
| Events | `deliberation.transition` audit trail | **reuse** existing `events` table via `AppendEvent` — do not fork a second event store |

## Phase locks (P20-00)

| Item | Lock |
|------|------|
| P19 loop | **Extend** `internal/loop` + `cmd/trace/loop.go`; do not replace; P19 Loop tests are keepers |
| New library | `internal/deliberation` — phase enum, `deliberation_state`, deterministic `SelectNext` |
| S01 owns | Controller + state persistence + transition events + hop budget (exact N locked in S01-00) |
| S06 owns | Phase-aware `loop next/apply/status` packet fields + fail-closed unknown apply keys + dual-version decision if needed |
| Assumption invalidate | Reuse provenance `STALE`/`SUPERSEDED` on existing Assumption rows — no second Assumption type |
| Task gates | Do **not** explode `work_state`; verification debt / gates via linked result records + BLOCKED semantics |
| MCP | Default **no new MCP tools** (G19 library inherit) |
| §16 / §18 | Experiments + risk-adaptive testing remain **Future** — not S01–S07 implement |

## Scope order (locked)

| Scope | Focus | Doc §§ |
|-------|--------|--------|
| S00 | Architecture lock vs live repo | 26–32, 27, 28 |
| S01 | Deliberation controller: phases, entry/exit, policies, events, budgets | 1, 3B, 4, 5, 6, 21, 25, 28 |
| S02 | Cognitive artifacts: uncertainty/question, assumption invalidate, hypothesis, decision alternatives + reconsider | 2, 3A, 7, 8, 9 |
| S03 | Change objects, expected vs actual effects, Git SHA refs | 10, 14, 29I, 29L |
| S04 | Test vs Verify vs Evaluate, gates, thin baseline, verification debt | 11, 12, 13, 20, 29F–H |
| S05 | Regression (correlated vs caused), reflection, observed relationships | 15, 17, 19, 29J |
| S06 | Phase-aware loop next/apply + context selection | 22, 23, 24, 29E |
| S07 | VERIFY: controller tests, protocol, mini-eval vs §31 story, DR-HANDOFF | 25, 29O–Q, 31 |

## MVP locks

- **Stdout-first**, harness-agnostic (P19).
- **Reuse** Goal, Task, Decision, DecisionAlternative, Assumption, Discovery, PlanChange, Claim, Evidence, Review, planner, `internal/loop`.
- Deliberation phases: `ORIENT INVESTIGATE EXPLORE PLAN CRITIQUE EXECUTE TEST VERIFY EVALUATE REFLECT REPLAN`.
- Controller is **inspectable**: every SelectNext writes an event `deliberation.transition` with inputs/scores/chosen phase.
- Stop: P19 saturation **or** phase=`CONTINUE` with no blocking uncertainty and verification policy met **or** max hops.
- No daemon, hosted MCP, embeddings, graph DB (Law 13).
- No raw CoT storage (doc §2).

## Out of scope unless promoted (still listed in COVERAGE.md)

- Experiment objects (§16)
- Risk-adaptive test matrix (§18)
- Model bake-off / learned policies
- Performance collectors as product
- Autonomous implementer that skips the human/agent

## Completion bar

A fresh agent, given only `trace loop next --task <id>`:

1. Sees a recommended **phase** and why (policy inputs).
2. Can record questions/uncertainties and cannot EXECUTE while blocking uncertainty is open.
3. Can record a Change with expected effects, then actual effects after work.
4. Cannot mark work complete on file edits alone — verify/eval/debt are distinct.
5. A contradicted expected effect or failed verification can open a Regression and force INVESTIGATE/REPLAN.
6. `trace why` / loop packet can answer “why did Trace choose INVESTIGATE?”
