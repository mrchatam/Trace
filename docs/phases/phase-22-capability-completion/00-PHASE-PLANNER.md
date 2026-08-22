# P22-00 — Phase 22 scaffold: capability completion

## Metadata
- id: P22-00
- todo_ids: [P22-00]
- role: planner
- skills: [planning-and-task-breakdown, spec-driven-development, writing-for-agents, incremental-implementation, brainstorming]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective

Lock Phase 22 against live repo + P21 close artifacts + the 43 unchecked checklist bullets. Confirm [`WORK-MAP.md`](WORK-MAP.md), [`DECISION-LOG.md`](DECISION-LOG.md), and README coverage matrix. Thicken S01–S08 scope prompts. **No product Go this row.**

This file **is** the planner output. A later agent re-running this row only re-confirms live inventory and thickens **upcoming** prompts; they do not rewrite `done` history.

## References

- [`DECISION-LOG.md`](DECISION-LOG.md)
- [`WORK-MAP.md`](WORK-MAP.md)
- [`README.md`](README.md) — coverage matrix (authoritative ownership)
- [`docs/CAPABILITIES_CHECKLIST.md`](../../CAPABILITIES_CHECKLIST.md) — only `- [ ]` items
- P21 [`DR-HANDOFF.md`](../phase-21-thoughtprocess-completion/DR-HANDOFF.md) (historical `no successor`)
- Protocol: [`docs/rules/agent-loop-protocol.md`](../../rules/agent-loop-protocol.md)
- Live: schema max **021**, `internal/loop/policy.go`, `internal/analyzers/`, `internal/mcp/server.go` (10 tools), `internal/install/` (cursor+claude only), `internal/store/fts.go` (`SearchFTS` unwired to CLI)

## Live inventory (2026-08-18 — planner)

| Surface | Location | P22 action |
|---------|----------|------------|
| Schema max | `internal/store/schema/` **021** files | S01 adds **022**; later scopes 023–026 |
| Compat ceiling | `evals/compat/compat_test.go` **21** (forbid 022+) | Bump per owning implement row |
| `BuildPolicyInputs` | `internal/loop/policy.go` | Does **not** set ExecutePending/TestPending/EvaluationPending/ReflectPending |
| SelectNext | `internal/deliberation/select.go` 14-row table | Keep; S03-01 feeds flags |
| MCP | `internal/mcp/server.go` **10** tools; `TestToolNamesRegistered` | S05/S08 add stdio tools; G19 library-first |
| Install registry | `internal/install/` cursor+claude | S02-01 add `git-hook` CONDITIONAL |
| FTS | `store.SearchFTS` / `retrieval.Engine.Search` | S05 CLI+MCP |
| Impact walk | `retrieval/impact_walk.go` file\|symbol only | S01-07 add tests via `validates` |
| Index | `analyzers.IndexFile` symbols+imports | S01 replace/add `code_edges` |
| Changes | `domain.CreateChange` via loop apply | S02-05 VCS-promoted path |
| Outcomes | `domain` test/verification/evaluation; comment “no test runner” | S03-03 explicit invoke (D-22-03) |
| Seed | P20 cognition keys; omits index | Keep omit code graph (D-22-06); add knowledge/eval keys in S06/S07 |

## Planner work

1. [x] Read skills, protocol, laws, P21/P20/P17 conventions.
2. [x] Inventory live packages vs 43 bullets + 15 audit findings.
3. [x] Lock 8 scopes, coverage matrix, mig sequence, DF-86 promote, D-16 supersede.
4. [x] Write thickened `00-PLANNER` / `01-*` / `02-*` (and extra pairs) for S01–S08.
5. [x] Board `docs/TODO/phase-22.md`; index + AGENTS Current focus.
6. [ ] Product Go — none this row.

## Exit criteria

- [x] S01–S08 prompts runnable alone (objective, files, named tests, exit criteria, spawn policy)
- [x] Coverage matrix lists all 43 bullets with owner row ids
- [x] DECISION-LOG records locks + supersessions
- [x] No product Go
- [x] DR-HANDOFF **OPEN** (S08-08 closes)

## Next

**P22-S01-00** after this row is `done`.
