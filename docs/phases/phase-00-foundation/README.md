# Phase 00 — Foundation (P0-X)

## Goal

Build the smallest **correct** foundation that passes **P0-X 7/7** and closes roadmap P0 — without MCP/daemon, without embeddings, without treating CRUD-only as success.

Authoritative bar: [`docs/init/C_FIRST_SCOPE.md`](../../init/C_FIRST_SCOPE.md), [`docs/init/I_BENCHMARK_PLAN.md`](../../init/I_BENCHMARK_PLAN.md), decisions in [`docs/init/D_DECISION_REGISTER.md`](../../init/D_DECISION_REGISTER.md).

## Locked phase defaults (do not weaken)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Surface | Library + `cmd/trace` CLI only |
| Store | `.trace/` in bound repo, gitignored; one SQLite per project |
| Analyzers | tree-sitter; TS/JS + Python |
| Git | `git` CLI behind VCS adapter interface; no blob duplication |
| P0-X bar | **7/7** including incremental localized file update |
| Incremental law | Full-rebuild-on-any-change is **forbidden** (G12 / DR-P0X #7; historical alias `DR-INCREMENTAL`) |
| MCP / daemon / HTTP | **Forbidden** this phase |
| Review policy | Every scope: `00-PLANNER` → `01` → `02-review` before next scope’s implement completes board order |

## P0-X must all pass

1. Goal/Task/Decision/Discovery round-trip + provenance  
2. Files + minimal symbols/imports (tree-sitter)  
3. `trace why` causal chain + reason codes  
4. `trace context` bounded task context  
5. Human seed matches fixture GT  
6. Several deterministic understanding queries (no LLM; default ≥5)  
7. Incremental update of one changed file **without** full fixture rebuild  

## Scope run order (P0-X critical path)

Confirmed against `C_FIRST_SCOPE` / historical `T001→T012a`:

| Scope | Theme | Board IDs | Folder |
|-------|--------|-----------|--------|
| S01 | Go module + `cmd/trace` scaffold | P00-S01-00/01/02 | `scopes/scope-01-go-scaffold/` |
| S02 | `.trace/` SQLite store + events | P00-S02-00/01/02 | `scopes/scope-02-store/` |
| S03 | VCS interface + git CLI + history index | P00-S03-00/01/02 | `scopes/scope-03-vcs/` |
| S04 | tree-sitter TS/JS + Python (incremental) | P00-S04-00/01/02 | `scopes/scope-04-analyzers/` |
| S05 | Work/causal API | P00-S05-00/01/02 | `scopes/scope-05-causal/` |
| S06 | Retrieval + context compiler | P00-S06-00/01/02 | `scopes/scope-06-retrieval-context/` |
| S07 | CLI commands for P0-X | P00-S07-00/01/02 | `scopes/scope-07-cli/` |
| S08 | Synthetic fixture + human seed + P0-X harness | P00-S08-00/01/02 | `scopes/scope-08-fixture-p0x/` |
| S09 | Phase verify / P0 close gate | P00-S09-00/01/02 | `scopes/scope-09-phase-verify/` |

Each scope folder has `00-PLANNER.md`, `01-*.md`, `02-scope-review.md`, and `SCOPE-TODOS.md`. Scope planners thicken `01-*`; this phase planner only light-locks order and cross-scope assumptions.

## Expected package sketch (locked by S01 planner)

```text
cmd/trace                 — thin CLI (no business logic)
internal/store            — .trace/ SQLite (not internal/sqlite)
internal/vcs              — VCS adapter interface
internal/gitcli           — git CLI implementation
internal/analyzers        — tree-sitter TS/JS + Python
internal/domain           — work/causal API
internal/retrieval        — exact / FTS / graph
internal/compiler         — context compiler (not context / contextx)
fixtures/x0               — synthetic project + ground truth (S08)
evals/p0x                 — deterministic foundation harness (S08)
```

Library must not import CLI. Future MCP (post–Phase 00) imports library only.

## Phase rules

- Phase planner first (`00-PHASE-PLANNER.md`), then scopes in order.
- Each scope: `00-PLANNER` → `01-implement` → `02-review` → spawned `a/b/…` until review confidence high.
- No Phase 01 until this phase VERIFY is `done`.
- Forward-only: do not rewrite `done` prompts; spawn remediation ahead.

## Out of scope (this phase)

MCP, daemon/HTTP, embeddings, env graph, impact engine, agent Gate C (X0), honesty/claim full path (Phase 01).
