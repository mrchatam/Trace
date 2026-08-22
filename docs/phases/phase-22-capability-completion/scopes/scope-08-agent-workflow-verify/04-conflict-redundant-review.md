# P22-S08-04 — Review: conflict / redundant work

## Metadata
- id: P22-S08-04
- todo_ids: [P22-S08-04]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C28** — advisory overlap visible to agents; no locking or blocking behavior.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md), [03-conflict-redundant.md](03-conflict-redundant.md)
- D-22-13: promote advisory only

## Review checklist

### C28 — reduce conflicting/redundant work

- [x] `TestDetectOverlappingOpenTasks`, `TestRedundantSimilarTitleSameGoal`, `TestNoConflictWhenTasksDisjoint` PASS
- [x] `TestTasksConflictsCLI` + `TestLoopNextIncludesWorkConflicts` PASS
- [x] `trace tasks conflicts` documented in help
- [x] Loop `work_conflicts` capped at **8**; omitted or empty when none
- [x] **Advisory only** — grep: no lock table, no transition block, no auto-DENIED on overlap
- [x] Active task filter matches planner (excludes DONE/SKIPPED/STALE/FAILED)

### Landmines

- [x] Detection logic lives in `internal/domain` (G19); loop/CLI thin
- [x] Schema **26**; no 027+ migrations this row
- [x] S03/S08-01 keepers still green
- [x] Optional MCP `include_conflicts` — if present, read-only parity with CLI (not implemented; optional per spec)

## Spawn policy

If unmet: spawn **`P22-S08-04a`** + **`P22-S08-04b`**. Do not close with residuals.

## Re-run commands

```bash
go test ./internal/domain/... ./internal/loop/... -count=1 -run 'TestDetectOverlapping|TestRedundantSimilar|TestNoConflictWhenTasksDisjoint|TestLoopNextIncludesWorkConflicts'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run TestTasksConflicts
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [x] C28 closed (checklist `[x]`) or spawned
- [x] Confidence **high**
- [x] Board Notes: evidence + checklist boxed
