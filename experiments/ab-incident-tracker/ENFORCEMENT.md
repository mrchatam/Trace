# E01 enforcement — Phase 23 product + harness

Same semantics as [docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md](../../docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md).

## G1 install (via prepare.sh)

- `.trace/config.json` → `{ "enforce": "strict" }`
- `trace install cursor --write` → `.cursor/rules/trace-enforcement.mdc`
- `trace install cursor-hook --write` → pre-edit gate when `TRACE_TASK_ID` set

## Seed task IDs (E01)

| UUID | Title |
|------|-------|
| `e0100000-0000-4000-8000-000000000010` | Design architecture (planning — gate blocks edit until plan exists) |
| `e0100000-0000-4000-8000-000000000020` | REST API and auth |
| `e0100000-0000-4000-8000-000000000030` | Responder/admin dashboard |
| `e0100000-0000-4000-8000-000000000040` | Public status page |
| `e0100000-0000-4000-8000-000000000050` | End-to-end verification |

## Agent obligations (G1)

1. `$TRACE_BIN loop gate --task "$TRACE_TASK_ID" --for edit` before product code edits
2. `transition … --to DONE --enforce` when closing tasks
3. `seed export -o trace/graph.json --strict --enforce` at end

## Verify mechanics (no app build)

```bash
./run-enforcement-demo.sh
```

Fresh seed must block edit with `plan_missing` on task `…0010`.
