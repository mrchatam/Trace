# Phase 02 — Gate C evaluation & slice hardening

## Goal

Run Experiment **X0** for real (not dry-run): measure understanding accuracy, misses, tokens, latency, and task success for **B0** (repo tools) vs **G1** (`trace` CLI `why`/`context`). Produce a documented **Gate C** Go / No-Go / iterate decision. Harden the vertical slice only where measurement exposes concrete gaps.

**Explicit:** Phase 01 VERIFY (X0 dry-run readiness) **≠** Gate C pass. Do not claim product-thesis success from `dry_run:true` metrics alone.

## Prior phase outcomes (live — carry forward)

| Item | Live value |
|------|------------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| Layout | `internal/{store,vcs,gitcli,analyzers,domain,retrieval,compiler,mcp}` + thin `cmd/trace` + `cmd/trace-mcp` |
| Store | `.trace/trace.db`, modernc sqlite, migrations `001`–`005` |
| DONE / Review | Review PASS via `review_judges_task` **or** explicit `AllowDoneWithoutReview` |
| Honesty | `evals/honesty` Paths A/B/C fail-closed (keep green) |
| X0 dry-run | `evals/x0` `TestX0DryRunMetricsB0AndG1`; schema v1; temp metrics `dry_run:true` — **instrument only**, not Gate C |
| P0-X | `evals/p0x` **7/7** — regression-keep; do not replace with agent scoring |
| MCP | Thin stdio `trace-mcp` (go-sdk v1.4.0); six tools; G19; **optional** for Gate C (DR-AGENT — CLI path sufficient) |
| Daemon / HTTP / embeddings | Still forbidden as primary surface |

## Live inventory — dry-run vs Gate C gaps (P02-00)

| Need | Phase 01 dry-run | Gate C (this phase) |
|------|------------------|---------------------|
| Schema-valid B0/G1 metrics | Yes (`dry_run:true`) | Real condition runs; `dry_run:false` (or sibling Gate C artifacts) |
| Agent understanding scoring vs GT | Absent (B0 stub read; G1 live CLI only) | Accuracy / misses primary |
| Tokens / latency / task success | Schema fields exist; not agent-aggregated | Quality + efficiency secondary |
| N≥3 agent runs / condition | No | Required (`I_BENCHMARK_PLAN`) |
| Go / No-Go / iterate report | No | Evidence table required |
| Kill criteria applied | Documented only | Apply: G1 ≤ B0 within error **and** non-trivial seeding cost → thesis endangered |

Prefer **extending** `evals/x0` (keep dry-run regression green). Keep `evals/honesty` and `evals/p0x` separate.

**S01 scope locks (P02-S01-00):** evidence = `scopes/scope-01-x0-gate-c/GATE-C-NOTES.md` + `docs/verification/gate-c-x0/`; ≥5 shared query bank; mean `understanding_accuracy`; kill when G1 mean ≤ B0 mean (N≥3) with non-trivial `fixtures/x0` seed; issue list shape `GC-NN` for S02.

## Locked phase defaults (do not weaken)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Gate | **Gate C** — documented pass / fail / iterate (`A_PROJECT_PLAN` Phase 2) |
| Instrument | Extend `evals/x0` (or sibling) beyond dry-run; keep `evals/p0x` + `evals/honesty` green |
| Conditions | **B0** = agent + ordinary repo tools; **G1** = agent + `trace` CLI `why`/`context` (may still read repo) |
| Corpus | Start with `fixtures/x0` + abs `seed/gt.json`; expand only with explicit scope-planner lock |
| Scoring | Understanding primary; implementation + honesty secondary; quality + efficiency metrics — **no** silent “G1 beats B0” without evidence table |
| N-runs | ≥3 agent runs per condition (`I_BENCHMARK_PLAN`) |
| Kill criteria | G1 understanding ≤ B0 within error **and** seeding cost non-trivial → thesis endangered |
| MCP | Optional surface; X0/Gate C remain CLI-path (DR-AGENT) |
| Daemon / HTTP | Forbidden as primary surface |
| Embeddings | Forbidden |
| Review policy | Every scope: `00-PLANNER` → `01` → `02-review` before next scope implement |
| DR-HANDOFF | Before phase VERIFY closes: scaffold Phase 03 (progressive planner) on Gate C **Go** / iterate-with-continue; on **No-Go** record explicit stop / `no successor` (or user override) |

## Scope run order (locked — P02-00)

| Scope | Theme | Board IDs | Folder |
|-------|--------|-----------|--------|
| S01 | X0 Gate C run + metrics + Go/No-Go draft | P02-S01-00/01/02 | `scopes/scope-01-x0-gate-c/` |
| S02 | Slice hardening from Gate C issue list | P02-S02-00/01/02 | `scopes/scope-02-slice-hardening/` |
| S03 | Phase VERIFY + Phase 03 handoff | P02-S03-00/01/02 | `scopes/scope-03-phase-verify/` |

Order is **Gate C run → slice hardening → phase verify/handoff**. Scope planners thicken `01-*`; this phase planner light-locks only.

## Cross-scope blast radius

- **S01** defines scoring + run protocol → thickens S02 issue-backlog shape.
- **S02** must not reopen P0-X bar or weaken honesty Paths A/B/C; measurement-driven fixes only.
- **S03** owns DR-HANDOFF for Phase 03 (or explicit stop on Gate C No-Go). S02 may be `skipped` with reason after No-Go if no hardening is warranted.

## Phase rules

- Run `00-PHASE-PLANNER` (`P02-00`) first, then scopes in order.
- Forward-only: do not rewrite Phase 00/01 `done` prompts.
- A1 / H1 remain EXPERIMENT_REQUIRED until Gate C evidence is written.
- Honesty + p0x bars stay green throughout.

## Out of scope (this phase)

- Progressive planner product (Phase 03+)
- Daemon / always-on HTTP
- Embeddings / env graph / impact engine
- Declaring commercial packaging / scale ceiling settled
- Replacing `evals/p0x` with agent eval
- Claiming Gate C from Phase 01 dry-run alone

## References

- [`docs/init/A_PROJECT_PLAN.md`](../../init/A_PROJECT_PLAN.md) Phase 2
- [`docs/init/I_BENCHMARK_PLAN.md`](../../init/I_BENCHMARK_PLAN.md) Experiment X0
- [`docs/init/H_VERIFICATION_STRATEGY.md`](../../init/H_VERIFICATION_STRATEGY.md)
- [`docs/init/D_DECISION_REGISTER.md`](../../init/D_DECISION_REGISTER.md) DR-AGENT, DR-HANDOFF
- Phase 01 VERIFY: [`../phase-01-x0-readiness/scopes/scope-05-phase-verify/VERIFY-NOTES.md`](../phase-01-x0-readiness/scopes/scope-05-phase-verify/VERIFY-NOTES.md)
- Protocol: [`docs/rules/agent-loop-protocol.md`](../../rules/agent-loop-protocol.md)
