# PROMPT-G1-DIRECTED-GAP — Mandatory gap pass (Session B)

> **Arm:** G1 Session B (directed-gap)  
> **Experiment:** E02 — validates **P25-C** gap-pass behavior after build-only Session A  
> **Trace binary:** `/home/ali/Desktop/Trace/bin/trace`

## STOP — verify workspace first

Your Cursor workspace **must** be:

`/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/G1`

**Turn 1:** Run `pwd`. If the path is not exactly the line above, **stop**.

## Before you start (operator)

Session A (build-only) should already be scored. Do **not** run this prompt before scoring Session A — that invalidates the build arm.

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
export PATH="/home/ali/Desktop/Trace/bin:$PATH"
export TRACE_TASK_ID=e0200000-0000-4000-8000-000000000010
export TRACE_PROJECT_ROOT=/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/G1
```

## Requirement

Run the **mandatory gap pass** from your installed Trace rules.

Follow `.cursor/rules/trace-enforcement.mdc` and hooks — record discoveries, decisions, and plan changes as your rules require.

**Order (FM-02 write-before-export):**

1. Write ≥1 discovery OR ≥1 decision linked to `$TRACE_TASK_ID` (`trace add discovery` / `trace add decision` with task links).
2. Only then re-export:

```bash
$TRACE_BIN seed export -o trace/graph.json --strict --enforce
```

Do **not** call `--strict --enforce` on a thin graph (discoveries=0 decisions=0) and backfill afterward.

## Trace (mandatory)

1. Use `$TRACE_BIN`; keep `TRACE_TASK_ID` on the active task.
2. **Before product code edits:** `$TRACE_BIN loop gate --task "$TRACE_TASK_ID" --for edit`
3. Execute the gap pass per installed rules — this is the directed confirmation that P25-C wiring works when prompted.
4. Discoveries/decisions **before** export (step order above).

## Workspace

**Project root:** `/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/G1`

- Edit only under this path
- Do not read or modify `runs/B0`

## After session

Operator scores directed arm:

```bash
./score.sh G1 --p25 --arm directed
```

Record P25-4 directed attestation in `experiments/RESULTS.md`.
