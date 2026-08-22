# Subagent delegation template (G1 Multitask)

Packet format for **any** Trace task — seed IDs or tasks the orchestrator added later.

This is **not** a closed list of four agents. Fill `{TASK_UUID}` from `trace tasks` / `loop next`, including new tasks.

---

## Subagent task — E01 G1

**Workspace (mandatory):** `/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/G1`

**Do not use any other directory.** Run `pwd` on turn 1; stop if it does not match exactly.

**Seed / Trace task UUID:** `{TASK_UUID}`  
**Task:** `{TASK_TITLE}`  
**Your slice:** `{SLICE}`

### Environment

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
export PATH="/home/ali/Desktop/Trace/bin:$PATH"
export TRACE_TASK_ID={TASK_UUID}
export TRACE_PROJECT_ROOT=/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/G1
```

### Before coding

1. Read **SPEC.md** and **PLANNING-MATRIX.md** if present; also `trace context` for `{TASK_UUID}` — Trace is source of order if they diverge.
2. `$TRACE_BIN context` (or Trace MCP `trace_context`) for `{TASK_UUID}`.
3. `$TRACE_BIN loop gate --task "$TRACE_TASK_ID" --for edit` — do not edit product code if blocked.

### Use Trace fully for your slice

You may record decisions, uncertainties, discoveries, evidence, reviews, and spawn or request follow-up tasks. You are not limited to “implement then DONE.”

### Deliver

- Implement your slice (or review/test if that is the task)
- Tests if you change product code
- Evidence in Trace; `--enforce` when marking **this** task DONE

### Seed IDs (examples only — not a roster)

| UUID | Seed title |
|------|------------|
| `e0100000-0000-4000-8000-000000000010` | Design architecture |
| `e0100000-0000-4000-8000-000000000020` | REST API and auth |
| `e0100000-0000-4000-8000-000000000030` | Responder/admin dashboard |
| `e0100000-0000-4000-8000-000000000040` | Public status page |
| `e0100000-0000-4000-8000-000000000050` | End-to-end verification |

Plus any UUID the orchestrator created in `.trace/`.
