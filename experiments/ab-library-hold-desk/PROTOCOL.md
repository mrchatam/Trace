# E03 Protocol — full-stack regression

## Purpose

Validate **Phases 25–28** on a new domain (library hold desk). Primary arm: **build-only**. Session-B only if build P25-3a FAIL.

## 0. Workspace

| Arm | Open Folder |
|-----|-------------|
| B0 | `experiments/ab-library-hold-desk/runs/B0` |
| G1 | `experiments/ab-library-hold-desk/runs/G1` |

Turn 1: `pwd` must match.

## 1. Prepare

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace

cd experiments/ab-library-hold-desk
./prepare.sh          # first time
# after B0:
./prepare.sh G1       # keeps B0
```

Verify install:

```bash
grep -q "mandatory gap pass" runs/G1/.cursor/rules/trace-enforcement.mdc && echo OK
grep -q "Parent orchestrator" runs/G1/.cursor/rules/trace-enforcement.mdc && echo OK
```

## 2. B0

1. Open `runs/B0` → paste `prompts/PROMPT-B0.md`
2. `./score.sh B0 --test`

## 3. G1 Session A — build-only (required)

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
export PATH="/home/ali/Desktop/Trace/bin:$PATH"
export TRACE_TASK_ID=e0300000-0000-4000-8000-000000000010
export TRACE_PROJECT_ROOT=/home/ali/Desktop/Trace/experiments/ab-library-hold-desk/runs/G1
```

1. Open `runs/G1`
2. Paste `prompts/PROMPT-G1-BUILD.md` — **no directed gap message**
3. Score:

```bash
P25_ATTEST_BUILD=Y ./score.sh G1 --test --p25 --arm build
```

**Do not** wipe G1 after this score.

## 4. G1 Session B — only if Session A thin

If P25-3a FAIL (`discoveries=0 decisions=0`):

1. Same `runs/G1` — **no** `./prepare.sh G1`
2. Paste `prompts/PROMPT-G1-DIRECTED-GAP.md`
3. Score:

```bash
P25_ATTEST_DIRECTED=Y ./score.sh G1 --test --p25 --arm directed
```

## 5. Phase 29 gate

Record RESULTS.md. Promote Phase 29 only per [README.md](README.md) decision table.

## Stopping conditions

- Wrong workspace → invalidate arm
- Accidental gap prompt before build score → invalidate Session A; re-prepare G1
- Operator wipes G1 after Session A → lose dual-lane evidence
