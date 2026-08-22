# P07 / S01 / 00-PLANNER — Incremental indexing / ignore tiers

## Metadata
- id: P07-S01-00
- todo_ids: [P07-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-incremental-indexing.md` for **incremental indexing quality** and **ignore tiers** against live Trace CLI/analyzers/store. Lock package paths, measurement hooks, and exit criteria. No product code in this planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 7
- [docs/STORAGE_AND_PERFORMANCE.md](../../../../STORAGE_AND_PERFORMANCE.md) §5 incremental + §9 tiers
- Live: `cmd/trace/index.go` `walkIndexable` + `gitIgnored`; `internal/analyzers` file-local incremental; store mig `001`…`010`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (gaps — from P07-00)
| Item | Today | S01 need |
|------|-------|----------|
| Walk skip | `.git` / `.trace` dirs only | **T0 always-skip** paths (vendor/`node_modules`/generated/binary) beyond gitignore |
| Ignore | Best-effort `git check-ignore` | Durable ignore-tier model aligned STORAGE T0–T3 (S01 owns T0; T1–T3 may be stubs/notes) |
| Incremental | File-local `IndexFile` + hash upsert; `TestIndexIncrementalIsolation` | Measurable incremental quality (unchanged sibling isolation + skip correctness + optional latency hooks for Gate H) |
| Perf harness | Absent | Prefer seed hooks / fixtures usable by later `evals/perf` (S03 owns Gate H thresholds) |
| Migration | Through `010` | Prefer **no** mig (path filter); additive `011_*` only if store needs tier metadata |

## Phase defaults already locked (respect — P07-00)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Honesty / Gates / ablation / p0x / x0 / Gate C | Keep green / intact |
| Full-rebuild-on-any-change | Forbidden |
| Gate H path | Prefer `evals/perf` — thresholds after measurements; S03 finalizes names |
| Daemon/HTTP/embeddings | Forbidden as primary |
| VerifiedFact / `plan simulate` | Out unless explicitly promoted |

## Planner work
1. Inventory live index/reindex/ignore behavior vs STORAGE_AND_PERFORMANCE incremental pipeline.
2. Lock ignore-tier model (T0 always-skip vs cold/warm notes) without inventing megastore.
3. Thickened `01-incremental-indexing.md` exit criteria + measurement surface (prefer evals or timed unit harness).
4. Light notes on S02/S03 Depends; sync SCOPE-TODOS.

## Locked defaults (finalized 2026-08-16 — S01-00)
| Item | Value |
|------|-------|
| Package | `cmd/trace` walk/index helpers; keep `internal/analyzers` file-local + NUL binary `SkipError` — **no** `internal/indexer` |
| Architecture | File-local incremental only |
| Migration | **None** — no `011_*` (path filter only) |
| T0 dirs | `.git`/`.trace`/`node_modules`/`vendor`/`__pycache__`/`.venv`/`venv`/`dist`/`.next`/`target`/`coverage` (no `build`/`bin` this row) |
| T0 files | basename `.min.js`/`.min.mjs`/`.min.cjs` + any path component matching T0 dir |
| Walk order | T0 SkipDir → DetectLanguage → T0 file/path → gitignore → collect |
| Measurement | Required: `TestIndexIncrementalIsolation` + T0 walk + explicit T0 argv tests; optional `evals/perf` fixture/timing seed — **no** Gate H pass/thresholds |
| T1–T3 | Notes only (not implemented) |
| Gate H | Do **not** declare Gate H pass |
| Carry-forward | Honesty A/B/C; Gate G/E/F; capability ablation; p0x; x0; Gate C `dry_run:false` |

## Exit criteria
- [x] `01-incremental-indexing.md` runnable alone with locked paths + VERIFY commands
- [x] No product Go; board Notes + next `P07-S01-01`

## Out of scope
- New language grammars (S02)
- Gate H threshold finalization (S03)
- Commercial multi-model perf theater
