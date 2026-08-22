# P02 / S01 / 00-PLANNER — X0 Gate C evaluation

## Metadata
- id: P02-S01-00
- todo_ids: [P02-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: mixed

## Objective
Finalize sibling `01-gate-c-eval.md` for **running Experiment X0** (beyond dry-run) and drafting the Gate C Go/No-Go evidence. Lock scoring, N-runs, corpus, and kill criteria against live `evals/x0` + `I_BENCHMARK_PLAN`. No product code.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 2
- Live: `evals/x0`, `evals/x0/schema.json`, `fixtures/x0`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (Gate C gaps vs Phase 01 dry-run)
| Item | Today (post–Phase 01) | Gate C need |
|------|----------------------|-------------|
| Dry-run | `TestX0DryRunMetricsB0AndG1` schema-valid B0/G1 with `dry_run:true` | Keep green as regression |
| Agent scoring | Absent — B0 stub file-read; G1 live why/context only | Understanding accuracy/misses vs human GT |
| Efficiency | Schema has latency/tokens; not agent-run aggregated | Capture on real runs |
| N≥3 agent runs | Not yet | ≥3 per condition |
| Go/No-Go report | Not yet | Evidence table + pass/fail/iterate |
| Kill criteria | Documented in `I_BENCHMARK_PLAN` | Apply honestly |

**Phase 01 dry-run ≠ Gate C pass.**

## Phase defaults already locked (respect; refine paths only)
| Item | Value |
|------|-------|
| Conditions | B0 = repo tools; G1 = `trace why`/`context` (+ repo OK) |
| Corpus | `fixtures/x0` + abs seed; expand only with explicit lock |
| Instrument | Prefer extend `evals/x0`; schema bump if breaking |
| Honesty / p0x | Keep separate + green |
| MCP | Optional; CLI path sufficient (DR-AGENT) |
| Kill | G1 ≤ B0 within error **and** non-trivial seeding cost → thesis endangered |

## Planner work
- Lock Gate C run protocol + metrics fields + evidence artifact path in `01-*`.
- Thicken exit criteria enough to run alone (do not invent unfair scoring).
- Light-update **upcoming** S02 stubs with expected issue-list shape.
- Sync SCOPE-TODOS.md + board Notes.

## Exit criteria
- [x] `01-gate-c-eval.md` runnable alone
- [x] Locked: B0/G1, corpus, scoring, N≥3, kill criteria, CLI-not-MCP, evidence path
- [x] No product Go; no false Gate C pass from dry-run

## Minimal todos
- [x] Diff dry-run harness vs Gate C needs (live)
- [x] Thicken 01 + light S02 notes
- [x] Board + SCOPE-TODOS
