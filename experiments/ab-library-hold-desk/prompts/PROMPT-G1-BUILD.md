# PROMPT-G1-BUILD — Library hold desk (Trace + post-P28 stack)

> **Arm:** G1 Session A (build-only)  
> **Experiment:** E03 — full-stack regression after Phases 25–28  
> **Trace binary:** `/home/ali/Desktop/Trace/bin/trace`

## STOP — verify workspace first

Your Cursor workspace **must** be:

`/home/ali/Desktop/Trace/experiments/ab-library-hold-desk/runs/G1`

**Turn 1:** Run `pwd`. If the path is not exactly the line above, **stop**.

## Before you start (operator)

After B0 completes:

```bash
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace

cd experiments/ab-library-hold-desk
./prepare.sh G1

export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
export PATH="/home/ali/Desktop/Trace/bin:$PATH"
export TRACE_TASK_ID=e0300000-0000-4000-8000-000000000010
export TRACE_PROJECT_ROOT=/home/ali/Desktop/Trace/experiments/ab-library-hold-desk/runs/G1
```

First-time: `./prepare.sh` (both arms).

## Requirement

Build a full **library hold / waitlist desk** (same feature list as B0 — see PROMPT-B0).

Underspecified — you decide stack and design.

## Trace (mandatory for G1)

1. Use `$TRACE_BIN`; keep `TRACE_TASK_ID` on the active task.
2. **Before product code edits:** `$TRACE_BIN loop gate --task "$TRACE_TASK_ID" --for edit`
3. Follow installed rules in `.cursor/rules/trace-enforcement.mdc` and hooks (mandatory gap pass, Parent orchestrator).
4. **Write-before-export:** record ≥1 discovery OR ≥1 decision linked to `$TRACE_TASK_ID` **before** export — do not export a thin graph first.
5. After discoveries: prefer `trace add task --from-discovery` or `loop apply` with `spawned_tasks[].discovery_id` when a BLOCKING gap needs a backlog row.
6. Export: `$TRACE_BIN seed export -o trace/graph.json --strict --enforce`

**Important for E03:** This prompt is **build-only**. Do **not** wait for the operator to ask for gap analysis — follow **installed Trace rules** (including post-build gap pass) as part of finishing.

## Workspace

**Project root:** `/home/ali/Desktop/Trace/experiments/ab-library-hold-desk/runs/G1`

- Edit only under this path
- Do not read or modify `runs/B0`

## Multitask orchestrator (INT-04 / FM-04)

If you delegate:

1. Set `TRACE_TASK_ID` **before** any product-code edit and before each subagent starts
2. Include workspace path + task UUID in every subagent prompt; tell workers to export `TRACE_TASK_ID` / `TRACE_PROJECT_ROOT` — do not assume Multitask inherits env
3. Parent owns loop/gate **and** graph writes — do not offload graph-only work while parent edits without task
4. Option A hook denies empty `TRACE_TASK_ID` under `enforce=strict` per process

## Constraints

- `go test ./...` passing
- README with run instructions
