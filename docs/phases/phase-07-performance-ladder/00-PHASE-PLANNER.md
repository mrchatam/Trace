# P07 / 00-PHASE-PLANNER — Performance ladder phase scaffold

## Metadata
- id: P07-00
- todo_ids: [P07-00]
- role: planner
- skills: [planning-and-task-breakdown, grilling]
- mcps: [Read, Grep, Glob, Write]
- agents: []
- verification: automated

## Objective

Light replan of **Phase 07 — Performance ladder & language plugins** after Phase 06 capability-selection ablation VERIFYs green and S03-02 closes DR-HANDOFF. Confirm scope order, refresh stub prompts (00/01/02), lock phase defaults from live repo + `A_PROJECT_PLAN` Phase 7, and sync board rows. **Do not** implement product Go. **Do not** deep-finalize every implement prompt — each scope’s `00-PLANNER` does that.

## References
- [agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [project-rules.md](../../rules/project-rules.md)
- [phase README](./README.md)
- [docs/init/A_PROJECT_PLAN.md](../../init/A_PROJECT_PLAN.md) Phase 7
- [docs/STORAGE_AND_PERFORMANCE.md](../../STORAGE_AND_PERFORMANCE.md)
- [docs/init/D_DECISION_REGISTER.md](../../init/D_DECISION_REGISTER.md) DR-HANDOFF
- [docs/TODO.md](../../TODO.md)
- Phase 06 VERIFY: [VERIFY-NOTES.md](../phase-06-environment-capability/scopes/scope-03-phase-verify/VERIFY-NOTES.md)

## Session start
Agent → clarify if needed → Plan → execute (planner only).

## Prior locks to respect (from Phase 06)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Capability-selection ablation | Green via `evals/capability` `TestPlantedCapabilitySelectionAblation` — do not weaken |
| Gate F / Gate G / Gate E / Gate C | Green / Go — do not silently reopen |
| Honesty / p0x / x0 / replan / impact / capability | Keep green |
| GC-03/04 | Deferred unless promoted with measurement |
| Daemon/HTTP/embeddings | Still forbidden as primary |
| VerifiedFact / `plan simulate` | Out unless explicitly promoted with Notes |
| Full-rebuild-on-any-change | Still forbidden |

## Planner work
1. Inspect live indexer/analyzer/CLI walk vs performance + language-plugin needs.
2. Confirm or revise S01→S03 order (incremental/ignore → language plugins → VERIFY/Gate H/handoff).
3. Ensure each scope folder has `00-PLANNER` / `01-*` / `02-scope-review` / `SCOPE-TODOS.md`.
4. Patch **upcoming** stubs only; never edit Phase 06 `done` prompt bodies.
5. Unblock/sync `docs/TODO.md` Phase 07 section; set this row `done` with Notes.
6. Update phase README if order/locks change.
7. Lock Gate H harness path (prefer planted/synthetic ladders under `evals/` — do not invent commercial multi-model perf theater; set thresholds after first measurements).

## Live inventory (confirmed 2026-08-16)
| Surface | Finding |
|---------|---------|
| CLI index | `cmd/trace/index.go` — `walkIndexable` skips `.git`/`.trace`; language filter; best-effort `git check-ignore`; file-local `indexOne` → `IndexFile`/`IndexFileAtRev` |
| Analyzers | `internal/analyzers` — JS/TS/TSX/Python only; SHA-256; file-local ReplaceSymbols/Imports; no plugin registry |
| Incremental proof | `cmd/trace` `TestIndexIncrementalIsolation` exists; no size-ladder / Gate H harness |
| Evals | `evals/{capability,honesty,impact,p0x,replan,x0}` — **no** `evals/perf` |
| Store | mig through `010_capability_surface.sql` |
| Gaps | No T0 ignore tiers beyond gitignore; no extra languages; no Gate H measurements |

## Locked phase defaults (do not weaken — set 2026-08-16)
| Item | Value |
|------|-------|
| Goal | Performance ladder & language plugins (`A_PROJECT_PLAN` Phase 7 / Gate H) |
| Scope order | S01 incremental indexing / ignore tiers → S02 language plugins → S03 VERIFY + Gate H + Phase 08 handoff |
| Validation gate | Gate H — prefer **`evals/perf`** planted ladders; preferred names `TestPlantedPerfLadderGateH` + `schema-gate-h.json` v1 + temp `metrics-gate-h.json` (**S03-00 finalizes**; thresholds **after first measurements**) |
| Package hint | Prefer `internal/analyzers` + CLI walk + optional `evals/perf` — avoid premature infra rewrite |
| Ignore-tier hint | STORAGE T0–T3; S01 owns T0 always-skip + measurable file-local incremental (**S01-00 finalizes**) |
| Language hint | Prefer **Go** first additional language (**S02-00 finalizes** grammar/module) |
| Migration hint | Prefer no mig; additive `011_*` only if needed |
| Phase 08 folder | **`phase-08-ecosystem-hardening`** |
| Review policy | Every scope gets `02-review` before next scope implement |
| Carry-forward bars | Honesty A/B/C; Gate G; Gate E; Gate F; capability ablation; p0x 7/7; Gate C artifacts intact |
| Perf policy | Measure first — no premature optimization theater |

## Cross-scope blast radius
- S01 indexing/ignore quality thickens S02 plugin load path and S03 Gate H ladders.
- S02 must not weaken Gate C / ablation / honesty / Gate F/G/E evidence integrity.
- S03 scaffolds Phase 08 `phase-08-ecosystem-hardening` (or records explicit `no successor`).

## Exit criteria
- [x] Phase README accurate vs live Phase 06 outcomes
- [x] Scopes S01–S03 each have 00/01/02 stubs (+ SCOPE-TODOS.md)
- [x] `docs/TODO.md` Phase 07 lists planner + scope rows; `P07-00` done after Phase 06 S03-02
- [x] No product Go in this row
- [x] Board Notes record locks + next runnable (`P07-S01-00`)

## Minimal todos
- [x] Inventory live indexer/analyzers/CLI vs perf + language gaps
- [x] Write/refresh S01–S03 stub prompts
- [x] Sync TODO.md Phase 07 section
- [x] Update README + mark P07-00 done

## Out of scope
- Product Go (indexers, new languages, Gate H harness implementation)
- Declaring Gate H pass without measurements
- Starting Phase 08 implement wave
