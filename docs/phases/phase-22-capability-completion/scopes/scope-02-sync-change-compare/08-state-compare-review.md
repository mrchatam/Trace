# P22-S02-08 — Review: state compare + S02 close

## Metadata
- id: P22-S02-08
- todo_ids: [P22-S02-08]
- role: reviewer
- skills: [code-review-and-quality, writing-for-agents]
- mcps: [Shell, Read, Grep]
- verification: automated

## Objective

Confirm **C06** and that **C04, C05, C25** still hold. S02 complete only if all four checklist capabilities are met.

## Session start

**Fresh subagent.** Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## S02 close matrix

| Capability | Checklist line | Proof |
|------------|----------------|-------|
| **C04** | Keep graph state synchronized… | hook/index path + honesty when HEAD ≠ watermark |
| **C05** | Record every meaningful change… | VCS promote + loop apply |
| **C06** | Allow comparison between project states | `CompareStates` + CLI |
| **C25** | Continuously update Trace… | `git-hook` install + docs |

Update **`docs/CAPABILITIES_CHECKLIST.md`** — box all four lines when closing.

## Review checklist

| # | Check |
|---|-------|
| 1 | No daemon / no HTTP / no wrap `git commit` |
| 2 | Schema **23**; no **024+** |
| 3 | S01 capabilities **C01–C03, C07** keepers still PASS (spot-check impact/validates) |
| 4 | MCP catalog still **10** tools |
| 5 | Compare + promote store no blobs |

```bash
go test ./internal/domain/... ./internal/install/... ./internal/store/... ./internal/compiler/... -count=1 -run 'TestCompareStates|TestPromoteVCSCommit|TestInstallGitHook|TestGraphSyncStaleWhenHeadDiffers'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/mcp/... -count=1 -run TestToolNamesRegistered
```

## Spawn policy

If any of C04/C05/C06/C25 unmet: spawn **`P22-S02-08a` + `P22-S02-08b`**. Do not close with residuals.

## Exit criteria

- [ ] C04, C05, C06, C25 closed or spawned
- [ ] Checklist boxes updated when closed
- [ ] Confidence **high**
- [ ] Board Notes → **S02 complete. Next `P22-S03-00`**
