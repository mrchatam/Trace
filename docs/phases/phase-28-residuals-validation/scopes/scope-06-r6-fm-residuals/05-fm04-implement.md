# P28-S06-05 — FM-04 / FR-P28-03 implementer

## Metadata
- id: P28-S06-05
- todo_ids: [P28-S06-05]
- role: implementer
- skills: [incremental-implementation, security-and-hardening, writing-for-agents]
- mcps: [user-codegraph]
- verification: mixed
- hooks: []

## Objective

**FR-P28-03 / FM-04:** Stop worker-only Trace when parent delegates graph work to subagents without `TRACE_TASK_ID` / loop gate. Extend harness parent Multitask / orchestrator rules, worker inheritance docs/tests, or parent-must-set-task guidance. Option A hook already denies empty-task under strict — do not reopen Option B (FR-P28-D4) unless human re-opens.

## References

- [00-PLANNER.md](00-PLANNER.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — FM-04; INT-04
- S03 Option A: `internal/install/enforcement.go` `CursorLoopGateHookScript`, `ParentOrchestratorRule`
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)

## Session start

Follow agent-loop-protocol Session start.

## Acceptance hint

Live or scripted: parent cannot complete edit path by offloading graph to workers while parent edits without task; document Cursor Multitask limits if product-unfixable.

## Preflight

```bash
cd /home/ali/Desktop/Trace
grep -n 'ParentOrchestrator\|TRACE_TASK_ID\|FailClosed' internal/install/enforcement.go | head
GOPROXY=direct go test ./internal/install/... -count=1 -run 'CursorLoopGate|HookDrift'
```

## Suggested work

1. Survey parent/worker task inheritance gaps vs Option A.
2. Harness/docs/tests: parent-must-set-task; optional worker env inheritance guidance.
3. Evidence in Notes + `FM04-NOTES.md` if Multitask limits are documentation-only.
4. Do **not** implement Option B parent-orchestrator detection.

## Out of scope

- FR-P28-D4 Option B; daemon/HTTP; rewriting S03 done prompts

## Exit criteria

- [ ] Acceptance hint met or Multitask limits documented with evidence
- [ ] Next runnable **P28-S06-06**

## Todo updates

Status + notes on **P28-S06-05** only.

## Next

`P28-S06-06`
