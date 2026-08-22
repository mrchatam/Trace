# E01 Protocol

Sequential: **prepare → B0 session → G1 session**.

## 0. Critical: workspace root

**Do not** paste prompts with Cursor opened at the Trace monorepo root.

| Arm | File → Open Folder |
|-----|-------------------|
| B0 | `/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/B0` |
| G1 | `/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/G1` |

Agent turn 1: run `pwd` — must match the table.

## 1. Prepare (operator)

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd /home/ali/Desktop/Trace/experiments/ab-incident-tracker
./prepare.sh          # first time only (both arms)
```

Creates fresh `runs/B0` and `runs/G1`. **After B0 is done, never run bare `./prepare.sh`** — that used to wipe B0. Prepare G1 with:

```bash
./prepare.sh G1       # does not touch runs/B0
```

Bare `./prepare.sh` now **refuses** if B0 already has product files (`PREPARE_FORCE=1` to override).

## 2. B0 session

1. Open Folder → `runs/B0`
2. Paste [prompts/PROMPT-B0.md](prompts/PROMPT-B0.md)
3. Agent builds incident tracker — **no Trace**
4. `./score.sh B0 --test`

## 3. G1 session

1. Operator: `./prepare.sh G1` if G1 is still empty, then env:

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
export PATH="/home/ali/Desktop/Trace/bin:$PATH"
export TRACE_TASK_ID=e0100000-0000-4000-8000-000000000010
export TRACE_PROJECT_ROOT=/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/G1
```

2. Open Folder → `runs/G1`
3. Paste [prompts/PROMPT-G1-ENFORCE.md](prompts/PROMPT-G1-ENFORCE.md)
4. **Multitask:** orchestrator completes deliberation on task `…0010` and writes SPEC/PLANNING-MATRIX **before** spawning implementers — see [MULTITASK.md](MULTITASK.md)
5. `./score.sh G1 --test --gate`

## Stopping conditions

- Agent cannot resolve crash loop in 3 retries → record partial
- G1 codes before deliberation completes → protocol violation

## Post-run

Update [../RESULTS.md](../RESULTS.md) with verdict and `./score.sh` output summary.
