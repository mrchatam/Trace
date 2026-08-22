# PROMPT-G1-ENFORCE — Incident tracker (Trace + enforcement)

> **Arm:** G1 (enforcement)  
> **Experiment:** E01 — ab-incident-tracker  
> **Trace binary:** `/home/ali/Desktop/Trace/bin/trace`

## STOP — verify workspace first

Your Cursor workspace **must** be:

`/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/G1`

**Turn 1:** Run `pwd`. If the path is not exactly the line above, **stop** and tell the operator to use File → Open Folder on that directory.

## Before you start (operator)

**Do not run `./prepare.sh` if B0 already completed.** Bare `./prepare.sh` used to wipe both arms. After B0:

```bash
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace

# G1 only — leaves runs/B0 untouched
cd experiments/ab-incident-tracker
./prepare.sh G1

export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
export PATH="/home/ali/Desktop/Trace/bin:$PATH"
export TRACE_TASK_ID=e0100000-0000-4000-8000-000000000010
export TRACE_PROJECT_ROOT=/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/G1
```

If G1 is already prepared, skip `./prepare.sh G1` and only export the env vars.

First-time (no B0 product yet): `./prepare.sh` or `./prepare.sh both`.

## Requirement

Build a full **on-call incident tracker** (same feature list as B0 — see PROMPT-B0).

This is intentionally underspecified. Decide stack, schema, API, auth, layout, and testing during implementation.

## Enforcement (mandatory for this arm)

1. Use `$TRACE_BIN` (repo `bin/` first on PATH).
2. Keep `TRACE_TASK_ID` set to the **active** task UUID (starts as planning `…0010`; change it when you switch tasks).
3. **Before every product code edit:** `$TRACE_BIN loop gate --task "$TRACE_TASK_ID" --for edit` — do not edit if exit ≠ 0.
4. Before marking DONE: `$TRACE_BIN loop status --task "$TRACE_TASK_ID"` and `$TRACE_BIN transition … --to DONE --enforce`.
5. Export with `$TRACE_BIN seed export -o trace/graph.json --strict --enforce`.

Harness rules + hook and `{ "enforce": "strict" }` are installed by `prepare.sh G1`.

## Workspace

**Project root:** `/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/G1`

- Edit only under this path
- Do not modify `runs/B0`
- Seed at `runs/G1/seed/gt.json` — **starting graph only**, not a cap on work
- Trace state: `runs/G1/.trace/`
- Export graph: `runs/G1/trace/graph.json`

## Use the full Trace surface (not four workers)

The five seed tasks (`…0010`–`…0050`) are a **starter backlog**, not the only work allowed.

**Orchestrator and workers may use all Trace capabilities**, including (CLI and/or MCP):

- goals, tasks, add, link, transition
- loop next / apply / status / gate
- decisions, assumptions, uncertainties, discoveries, plan_changes
- context, why, review, capability, version
- new tasks and new agents when the plan needs them (reviewer, tester, schema, UI, etc.)

When you record a **discovery** for new work, prefer **`trace add` a task** (or MCP equivalent) and link it — do not only attach discoveries to existing seed rows.

Do **not** treat “delegate only `…0020`–`…0050`” as a rule. Spawn more tasks in Trace when you discover work.

[SUBAGENT-DELEGATION.md](SUBAGENT-DELEGATION.md) is a **packet format** (workspace + task UUID + gate). It is not a closed roster of four agents.

## Multitask orchestrator (you)

**You (the parent agent) must use Trace**, not only subagents:

1. Complete deliberation on the current planning task first (`loop next`, `apply`, decisions, uncertainties).
2. Write **SPEC.md** and **PLANNING-MATRIX.md** before product implementation (you may still add Trace tasks later).
3. Delegate **any** Trace task UUID with the delegation template — seed IDs or ones you created.
4. Serialize Trace graph writes (one writer at a time).

Subagents gate and implement their slice; **you** own the planning graph and may add work as you go.

## Constraints

- Use Trace CLI and/or Trace MCP throughout
- Record goals, tasks, decisions, uncertainties, and evidence as work progresses
- Deliver working app with tests passing (`go test ./...`)
