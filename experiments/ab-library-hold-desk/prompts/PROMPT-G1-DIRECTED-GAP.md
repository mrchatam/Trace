# PROMPT-G1-DIRECTED-GAP — Mandatory gap pass (Session B, optional)

> **Arm:** G1 Session B (directed-gap)  
> **Experiment:** E03 — only if Session A P25-3a FAIL  
> **Trace binary:** `/home/ali/Desktop/Trace/bin/trace`

## STOP — verify workspace first

Your Cursor workspace **must** be:

`/home/ali/Desktop/Trace/experiments/ab-library-hold-desk/runs/G1`

**Turn 1:** Run `pwd`. If wrong, **stop**.

## Before you start (operator)

Session A must already be scored. **Do not** `./prepare.sh G1` (would wipe evidence).

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
export PATH="/home/ali/Desktop/Trace/bin:$PATH"
export TRACE_TASK_ID=e0300000-0000-4000-8000-000000000010
export TRACE_PROJECT_ROOT=/home/ali/Desktop/Trace/experiments/ab-library-hold-desk/runs/G1
```

## Requirement

Run the **mandatory gap pass** from your installed Trace rules (`.cursor/rules/trace-enforcement.mdc`).

Record discoveries/decisions linked to `$TRACE_TASK_ID`, promote BLOCKING gaps if needed, then:

```bash
$TRACE_BIN seed export -o trace/graph.json --strict --enforce
```

## Constraints

- Do not wipe the workspace
- Do not re-run the build-only product from scratch unless fixing gaps
- Prefer fixing gaps and enriching the graph
