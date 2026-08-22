# P22-S05-00 — Planner: query / history CLI+MCP

## Metadata
- id: P22-S05-00
- todo_ids: [P22-S05-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep]
- verification: automated

## Objective

Lock S05. Owned: **C17, C29, C30, C31, C32, C33, C34, C35, C37, C42-surface**. FTS exists unwired to CLI. **No product Go.**

## Live inventory (2026-08-18, post-S04)

| Surface | Live state |
|---------|------------|
| Schema max | **024** (`024_impact_compare.sql`; 24 embed sql files) |
| Compat ceiling | **24** (`evals/compat/compat_test.go` — no 025+) |
| FTS | `internal/store/fts.go` — `SearchFTS` (unicode61, sanitize AND-quoted tokens, default limit 32, cap 64); indexes goals/tasks/decisions/…/changes/regressions/reflections/outcome_results/files/symbols |
| Retrieval search | `internal/retrieval/search.go` — `Engine.Search` wraps FTS → `[]Hit` with `reason_code=fts_match`; used by **compiler** task-title FTS only |
| CLI search | **Absent** — no `trace search` |
| CLI changes | `cmd/trace/changes.go` — **`capture\|compare` only**; no `list\|show`; `cli:changes` gated |
| CLI tests | `cmd/trace/test.go` — **`run` only**; no `verifying` query |
| CLI outcomes | `cmd/trace/outcomes.go` — **`compare\|improvements` only**; no `failed\|worked` |
| CLI regressions | **Absent** — no `trace regressions` root |
| MCP catalog | **10** tools (`RegisteredToolNames` in `internal/mcp/server.go`); **no** `trace_search`, `trace_changes`, `trace_regressions` |
| Context packet | `internal/compiler/packet.go` — `items`, `why_trace`, capabilities, honesty; **no** top-level `evaluations` / `reflections` / `planning_evidence` |
| Loop next | `internal/loop/next.go` — has `recent_changes`, `historical_relationships`; **no** `planning_evidence` section |
| `trace why` | Live (`cmd/trace/why.go`, `trace_why` MCP) — causal neighborhood; **not** C29 (history is search + changes list/show) |
| S01 deps | `ListValidatesForSymbol`, `ListValidatesForFile` — C31 |
| S02 deps | `ListAllChanges`, `GetChange`, `ListChangePaths`, `CompareStates` — C30 |
| S03 deps | `outcome_results` kinds test/evaluation/verification — C32 query surface |
| S04 deps | `ListRegressionsByChangeID`, `ListAllRegressions`, `ListImprovementsByTaskID`/`ByChangeID` — C17/C33/C34 data |

S01–S04 closed — do not reopen schema, impact compare, or verification cycle in S05 prompts.

## References

- [DECISION-LOG.md](../../DECISION-LOG.md) D-22-05 (MCP catalog may grow for CLI mirrors), G19
- [WORK-MAP.md](../../WORK-MAP.md) W-17…W-22
- Coverage: [README.md](../../README.md) C17, C29–C35, C37, C42 rows

## FINAL locked defaults

| Item | Value |
|------|-------|
| SQL | **None** — reuse FTS + existing tables |
| Compat | Stays **24** (forbid **025+** entire S05) |
| G19 | Domain/retrieval query APIs; CLI + MCP thin encode only |
| Bounded | Default limit **32**, hard cap **64** (match `SearchFTS` / `Engine.Search`) |
| CLI search | `trace search <query> [--limit N]` → JSON `{ok,hits[],count}` using `retrieval.Engine.Search` |
| CLI changes | Extend `trace changes` → **`list\|show\|capture\|compare`**; `list [--task <id>] [--limit N]` newest-first; `show <change-id>` → change + `paths[]` (no blobs) |
| Evidence CLI | `trace tests verifying --symbol <uuid> \| --file <path>`; `trace outcomes failed [--task <id>] [--limit N]`; `trace outcomes worked [--task <id>] [--limit N]`; `trace regressions list [--task <id>] [--change <id>] [--limit N]` |
| MCP S05-01 | **`trace_search`** (query + limit); **`trace_changes`** (action=`list\|show\|compare`, same flags as CLI) |
| MCP S05-03 | **`trace_regressions`** (action=`list`, filters); tests/outcomes evidence **CLI-only** (avoid catalog explosion) |
| MCP catalog | S05-01 → **12** tools; S05-03 → **13** tools; update `RegisteredToolNames` + `TestToolNamesRegistered` + `BuiltinMCPCapabilitySpecs` each row |
| Capability slugs | Add `cli:search`, `cli:regressions`; extend `cli:tests` / `cli:outcomes` / `cli:changes` gating for new subcommands |
| Context | Add packet fields **`evaluations`**, **`reflections`**, **`planning_evidence`** — task-scoped, cap **8** each, compact JSON (id, kind, summary, created_at; no score blobs) |
| Loop next | Add **`planning_evidence`** section mirroring packet evidence slices (freshness from seed); schema string **`trace.loop.next.v1`** unchanged |
| C29 vs why | History questions → **search + changes**; do not extend `trace why` for list semantics |
| Checklist | Implementers **unbox** owned caps; reviewers **box** after review rows |

## Named tests

| Test | Row |
|------|-----|
| `TestCLISearchUsesFTS` | S05-01 |
| `TestCLIChangesList` | S05-01 |
| `TestCLIChangesShow` | S05-01 |
| `TestMCPSearchRegistered` | S05-01 |
| `TestMCPChangesList` | S05-01 |
| `TestToolNamesRegistered` | S05-01 (expect **12**) |
| `TestTestsVerifyingQuery` | S05-03 |
| `TestOutcomesFailedAndWorked` | S05-03 |
| `TestRegressionsListQueryable` | S05-03 |
| `TestMCPRegressionsRegistered` | S05-03 |
| `TestToolNamesRegistered` | S05-03 (expect **13**) |
| `TestContextIncludesEvaluationsAndReflections` | S05-05 |
| `TestLoopNextPlanningEvidenceSection` | S05-05 |

## Exit criteria

- [x] 01–06 thickened
- [x] No product Go

## Next

**P22-S05-01**
