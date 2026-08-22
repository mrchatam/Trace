# PROMPT-G1-BUILD — Equipment checkout desk (Trace + Phase 25 install)

> **Arm:** G1  
> **Experiment:** E02 — validates **Phase 25** gap-pass install (INT-03/04)  
> **Trace binary:** `/home/ali/Desktop/Trace/bin/trace`

## STOP — verify workspace first

Your Cursor workspace **must** be:

`/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/G1`

**Turn 1:** Run `pwd`. If the path is not exactly the line above, **stop**.

## Before you start (operator)

After B0 completes:

```bash
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace

cd experiments/ab-p25-gap-pass-validation
./prepare.sh G1

export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
export PATH="/home/ali/Desktop/Trace/bin:$PATH"
export TRACE_TASK_ID=e0200000-0000-4000-8000-000000000010
export TRACE_PROJECT_ROOT=/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/G1
```

First-time: `./prepare.sh` (both arms).

## Requirement

Build a full **equipment checkout desk** (same feature list as B0 — see PROMPT-B0).

Underspecified — you decide stack and design.

## Trace (mandatory for G1)

1. Use `$TRACE_BIN`; keep `TRACE_TASK_ID` on the active task.
2. **Before product code edits:** `$TRACE_BIN loop gate --task "$TRACE_TASK_ID" --for edit`
3. Follow installed rules in `.cursor/rules/trace-enforcement.mdc` and hooks (including the mandatory gap pass).
4. **Write-before-export:** if the gap pass records discoveries/decisions, write them (linked to `$TRACE_TASK_ID`) **before** export — do not export a thin graph first.
5. Export only after those writes: `$TRACE_BIN seed export -o trace/graph.json --strict --enforce`

**Important for E02:** This prompt is **build-only**. Do **not** wait for the operator to ask for gap analysis — follow your **installed Trace rules** (including any post-build gap pass) as part of finishing the session.

## Workspace

**Project root:** `/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/G1`

- Edit only under this path
- Do not read or modify `runs/B0`

## Multitask orchestrator (INT-04 / FM-04)

If you delegate:

1. Set `TRACE_TASK_ID` to the active seed task **before** any product-code edit and before each subagent starts
2. Include workspace path + task UUID in every subagent prompt; tell workers to `export TRACE_TASK_ID=…` and `TRACE_PROJECT_ROOT=…` — do not assume Multitask inherits parent env
3. Parent owns Trace loop/gate **and** graph writes (gap pass, discoveries, decisions) — do not offload graph-only work to workers while the parent edits without task
4. Option A hook denies empty `TRACE_TASK_ID` under `enforce=strict` per process; Trace cannot product-detect parent orchestrators (Option B out)

## Constraints

- `go test ./...` passing
- README with run instructions
