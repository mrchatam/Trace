# P22-S05-05 — Implement: context / planning evidence

## Metadata
- id: P22-S05-05
- todo_ids: [P22-S05-05]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Surface evaluations, reflections, and accumulated evidence in **context** and **loop next** so agents can query evidence when planning (**C35, C42-surface**). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md). Agent → clarify → Plan → execute.

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- `internal/compiler/packet.go`, `compiler.go` — TaskContext/ExpandContext
- `internal/loop/next.go` — BuildNextPacket sections
- S03/S04: `outcome_results`, `reflections`, `improvements`, `regressions` store APIs

## Live baseline (do not re-ship)

| Present | Absent |
|---------|--------|
| Context `items[]` + honesty + why_trace | Top-level `evaluations`, `reflections`, `planning_evidence` |
| Loop `recent_changes`, `historical_relationships` | `planning_evidence` section |
| `ListOutcomeResultsByTaskKind`, `ListReflectionsByTaskID` | Compiler aggregation for packet |
| MCP `trace_context` mirrors compiler | Automatic inheritance once compiler updated (G19) |
| Compat **24** | **025+** |

## Locked defaults

| Item | Value |
|------|-------|
| SQL | **None**; compat stays **24** |
| Scope | **Task-scoped** — use compile `taskID` / loop seed task |
| Cap | **8** rows per section (evaluations, reflections, planning_evidence) |
| Order | Newest first: `created_at DESC, rowid DESC` (match S03 regression tie-break) |
| evaluations | `kind=evaluation` outcome_results for task — fields: id, task_id, summary, scores_json (truncate 512 if needed), created_at |
| reflections | `ListReflectionsByTaskID` — fields: id, summary, created_at (no full assumption blobs in packet) |
| planning_evidence | Mixed slice for planning: **open regressions** for task + **recent failed test outcomes** + **improvements** for task; cap 8 total; each item: `{entity_type, entity_id, title, summary, created_at}` |
| Packet JSON | Add optional arrays on `compiler.Packet` — **additive**; `schema_version` stays `"0.2"` |
| Loop next | Add `PlanningEvidenceSection` with same three arrays + `freshness`; **`trace.loop.next.v1`** unchanged |
| Markdown | Render short labeled sections when `IncludeMarkdown` (mirror existing packet MD style) |
| MCP | No new tools — `trace_context` picks up compiler fields automatically |
| C42-surface | Evaluations visible to future agents via context/MCP — full eval contract remains S07 |

## Requirements

1. **`internal/compiler/evidence_sections.go`** (new) — build helpers called from `compileAtDepth`.
2. Extend `Packet` struct + JSON/Markdown render in `packet.go`.
3. Extend `BuildNextPacket` — populate `planning_evidence` section (may reuse same builder with task id).
4. Named tests in `compiler_test.go`, `loop/next_test.go` or `cmd/trace` integration tests.
5. Checklist C35, C42-surface **unboxed** until S05-06.

## Touch files

- `internal/compiler/evidence_sections.go`, `evidence_sections_test.go` (new)
- `internal/compiler/packet.go`, `compiler.go`
- `internal/loop/next.go`, `next_test.go` (extend)
- `cmd/trace/context_test.go` or `loop_test.go` (extend if needed)

## Named tests

| Test | Proves |
|------|--------|
| `TestContextIncludesEvaluationsAndReflections` | TaskContext JSON has non-empty evaluations + reflections when seeded |
| `TestLoopNextPlanningEvidenceSection` | loop next includes planning_evidence with regressions/failures/improvements when seeded |
| `TestTaskContextAndBudgets` | keeper — existing budget behavior unchanged |

```bash
go test ./internal/compiler/... ./internal/loop/... -count=1 -run 'TestContextIncludesEvaluations|TestLoopNextPlanningEvidence|TestTaskContextAndBudgets'
CGO_ENABLED=1 go test ./cmd/trace -count=1 -run 'TestContext|TestLoopNext'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # expect 24
```

## Exit criteria

- [ ] C35 and C42-surface true
- [ ] Named tests PASS
- [ ] `trace context` + MCP `trace_context` expose new fields (grep/spot-check)
- [ ] Board Notes

## Minimal todos

- [ ] Compiler evidence sections (cap 8)
- [ ] Loop next planning_evidence section
- [ ] Tests + MD render
- [ ] Board notes

## Residual risks (carry to S05-06)

- Packet token budget — evidence sections are **outside** items[] budget (additive); document in test if totals grow
- Empty sections omit vs `[]` — lock: emit empty arrays for stable JSON shape
- Cross-task leakage — evaluations/reflections must filter by task_id strictly
