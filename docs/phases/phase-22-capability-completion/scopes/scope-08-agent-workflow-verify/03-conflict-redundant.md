# P22-S08-03 — Implement: conflict / redundant work

## Metadata
- id: P22-S08-03
- todo_ids: [P22-S08-03]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

**Reduce conflicting or redundant work** (**C28**) via advisory overlap detection visible to agents. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — D-22-13, conflict semantics
- Live: `internal/store/changes.go` (`ListChangesByTaskID`, `ListChangePaths`), `internal/domain/impact_compare.go` (`seedsFromChangePaths`)

## Live baseline

| Present | Absent |
|---------|--------|
| Tasks + changes + paths in store | `internal/domain/conflicts.go` |
| Loop next rich packet | `work_conflicts[]` section |
| `trace tasks` list only | `trace tasks conflicts` |

## Locked defaults

| Item | Value |
|------|-------|
| Advisory | Surface overlaps; **never** block transitions, apply, or take locks |
| Active tasks | `work_state` ∈ {PENDING, IN_PROGRESS, AWAITING_REVIEW, BLOCKED} |
| Path overlap | Tasks A,B conflict if any repo-relative path prefix overlaps (union paths from each task’s changes via `ListChangesByTaskID` → `ListChangePaths`) |
| Impact overlap | OR intersecting `ImpactSeedRef` (file/symbol entity_id) from `seedsFromChangePaths` |
| Title redundancy | Same non-nil `goal_id`; normalized titles equal **or** mutual substring (min **8** runes after lower+trim+collapse space) |
| CLI | `trace tasks conflicts [--task <id>]` → JSON `{ok,conflicts:[{task_a,task_b,reason,paths[]?,seed_overlap?}]}` |
| Loop next | Add **`work_conflicts`** section to `NextPacket` (schema **`trace.loop.next.v1` unchanged**); cap **8** entries for seed task |
| MCP | **Optional**: `trace_tasks` `include_conflicts=true` — implement if trivial; **CLI required** |
| Schema | **No SQL** — computed only; compat stays **26** |

## Requirements

1. **`domain.DetectWorkConflicts(ctx, svc, opts)`** — returns bounded conflict list; stable sort by task id pair.
2. **`normalizeTitle(s string) string`** — shared helper for title redundancy tests.
3. Wire **`buildWorkConflictsSection`** in `internal/loop/next.go` for seed task.
4. Extend **`cmd/trace/tasks.go`** with `conflicts` subcommand; gate `cli:tasks`.
5. Unit tests with fixtures: two IN_PROGRESS tasks, overlapping paths, disjoint control, same-goal similar titles.

## Touch files

- `internal/domain/conflicts.go` (new) + `conflicts_test.go`
- `internal/loop/next.go` — `WorkConflictsSection`, builder
- `cmd/trace/tasks.go` — conflicts subcommand + test
- `cmd/trace/help.go` — document `tasks conflicts`
- Optional: `internal/mcp/tools_parity.go` — `include_conflicts` on `trace_tasks`

## Named tests

| Test | Proves |
|------|--------|
| `TestDetectOverlappingOpenTasks` | path prefix overlap between active tasks |
| `TestRedundantSimilarTitleSameGoal` | same goal + similar title |
| `TestNoConflictWhenTasksDisjoint` | negative — different paths, different goals |
| `TestTasksConflictsCLI` | CLI JSON shape + `--task` filter |
| `TestLoopNextIncludesWorkConflicts` | loop next packet includes capped conflicts when seeded |

```bash
go test ./internal/domain/... -count=1 -run 'TestDetectOverlapping|TestRedundantSimilar|TestNoConflictWhenTasksDisjoint'
go test ./internal/loop/... -count=1 -run TestLoopNextIncludesWorkConflicts
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run TestTasksConflicts
```

## Exit criteria

- [ ] C28 true (advisory, visible in CLI + loop next)
- [ ] Named tests PASS
- [ ] Checklist C28 unboxed on evidence (review boxes)
- [ ] Board Notes

## Minimal todos

- [ ] Detector + CLI + loop section
- [ ] Tests
- [ ] Help string
- [ ] Board notes

## Residual risks (carry to S08-04)

- Path overlap false positives on short prefixes (e.g. `a` vs `ab`) — prefer path segment boundary or require `/` prefix rules
- Performance on many open tasks — cap pairwise scan or limit to same-goal first
- S09 E02 routing may also reference conflicts — must not duplicate divergent logic
