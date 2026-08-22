# P07 / S02 / 00-PLANNER — Language plugins / adapters

## Metadata
- id: P07-S02-00
- todo_ids: [P07-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-language-plugins.md` for **additional language adapters** via the analyzers boundary. Lock which language(s), adapter interface, and exit criteria. No product code in this planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 7
- Live: `internal/analyzers` tree-sitter JS/TS/TSX/Python (`DetectLanguage` / `extract_*`)
- S01 Depends: incremental/ignore surface from S01
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (gaps — from P07-00)
| Item | Today | S02 need |
|------|-------|----------|
| Languages | `.js/.jsx/.mjs/.cjs`, `.ts`, `.tsx`, `.py` only | ≥1 additional adapter |
| Boundary | Switch in `DetectLanguage` + `extract(lang)` | Keep adapter-shaped extension — not universal language theater |
| CGO | Analyzers-only (store/domain CGO-free) | Preserve |
| Phase hint | Prefer **Go** (self-host Trace) | **S02-00 finalizes** language + official grammar module + go 1.24 floor |

## Phase defaults already locked (respect)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Carry-forward bars | Keep green |
| Plugin policy | Adapter boundary — not universal language theater |
| Language hint | Prefer Go first (**finalize here**) |
| Daemon/HTTP/embeddings | Forbidden as primary |
| S01 ignore/walk | Consume S01 T0 surface (dirs/`*.min.js`/path-segment + gitignore order); do not regress T0 skip when adding languages — **live-confirm (S01-02 APPROVE):** `cmd/trace` `isT0SkipDir`/`isT0SkipPath`/`walkIndexable` order T0→lang→T0 file/path→gitignore + explicit argv `skipped++` |

## Planner work
1. [x] Choose first additional language(s) + grammar deps with go 1.24 floor → **Go** + `tree-sitter-go` v0.25.0.
2. [x] Lock DetectLanguage / IndexFile extension points → DetectLanguage + `extract` switch + `extract_go.go`.
3. [x] Thickened `01-language-plugins.md` + light S03 Depends; SCOPE-TODOS sync.

## Locked defaults (FINAL — 2026-08-16)
| Item | Value |
|------|-------|
| Language | **Go** (exactly one additional adapter) |
| Grammar | `github.com/tree-sitter/tree-sitter-go` **v0.25.0** (`bindings/go`); go 1.22 floor ≤ Trace **1.24.0** |
| Package | `internal/analyzers` — extend `DetectLanguage` + `extract` switch + `extract_go.go` |
| Adapter policy | Switch/extension only — **not** plugin registry / universal language theater |
| Extension | `.go` → lang id `"go"` |
| CGO | Analyzers-only |
| Migration | No `011_*` |
| S01 | Consume T0 walk; do not regress (`vendor` stays T0) |
| Golden tests | Required for Go smoke/symbols/imports |
| Gate H | May extend ladder fixtures with `.go` later; **no** threshold invent in S02 |
| Carry-forward | Honesty; Gates E/F/G; capability ablation; p0x; x0; Gate C intact |

## Exit criteria
- [x] `01-language-plugins.md` runnable alone
- [x] No product Go; board Notes + next `P07-S02-01`

## Out of scope
- Gate H threshold finalization (S03)
- Rewriting JS/TS/Python extractors without need
