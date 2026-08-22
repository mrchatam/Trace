# P22-S02-04 — Review: graph sync honesty

## Metadata
- id: P22-S02-04
- todo_ids: [P22-S02-04]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C04** (graph synchronized **or** lag is honest and closable via `trace index` / hook).

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Review checklist

| # | Check |
|---|-------|
| 1 | Schema max **23** exactly (23 sql files); **no 024+** |
| 2 | `graph_sync_state` separate from `vcs_meta` / `MetaVCSWatermark` |
| 3 | `trace index status` reports head vs watermark |
| 4 | Commit lag surfaces on context packet and/or loop freshness **dirty** |
| 5 | Existing disk `IndexHonesty` tests still PASS (`TestIndexStaleBanner`, etc.) |
| 6 | Index still file-local — no full-rebuild glob added |
| 7 | Watermark updates after successful index, not before |

```bash
go test ./internal/store/... ./internal/compiler/... ./cmd/trace/... -count=1 -run 'TestGraphSyncStaleWhenHeadDiffers|TestHookIndexUpdatesLastIndexedCommit|TestOpenCreatesDBAndMigratesIdempotent|TestIndexStaleBanner'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l   # expect 23
```

## Spawn policy

If C04 unmet: spawn **`P22-S02-04a` + `P22-S02-04b`**. Do not close with residuals.

## Exit criteria

- [ ] C04 closed or spawned
- [ ] Confidence **high**
- [ ] Board Notes → **Next `P22-S02-05`**
