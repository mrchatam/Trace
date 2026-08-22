# E02 Protocol — P25 validation

## Purpose

Validate **Phase 25 harness** (gap pass install + orchestrator rules), not a full Trace thought-process victory. G1 uses **two scored sessions**: Session A (build-only) and optional Session B (directed-gap).

## 0. Workspace

| Arm | Open Folder |
|-----|-------------|
| B0 | `experiments/ab-p25-gap-pass-validation/runs/B0` |
| G1 | `experiments/ab-p25-gap-pass-validation/runs/G1` |

Turn 1: `pwd` must match.

## 1. Prepare

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd /home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation
./prepare.sh          # first time both arms
# after B0 done:
./prepare.sh G1       # G1 only — keeps B0
```

G1 prepare runs `trace install cursor --write` + hook — must include **Phase 25** GapPassPrompt (rebuild trace binary from repo if install output lacks “mandatory gap pass”).

**Export:** `prepare.sh` does **not** export. After your agent session, export may happen during the session (per prompt) or automatically: **`score.sh` preflight exports `trace/graph.json` from the DB when the file is missing** before G2/P25 checks. **T02** then runs `seed export --strict --enforce` (both arms); thin graphs fail G2. Re-export after product commits if FM-07 warns about SHA drift.

**Git-sparsity (FM-07):** `score.sh` T03 compares `exported_at_commit` in `trace/graph.json` to `git HEAD`. Match → pass; missing git/SHA → skip; mismatch → **WARN** (“re-export recommended”) only. **warn ≠ fail** — FM-07 drift never fails G2 alone (unlike T02 thin/honesty enforce). Post-hoc SPEC commits without re-export are expected to WARN until you re-export. Decision SoT (remain warn-only; no plan-before-edit fail-closed without human approval): [`docs/phases/phase-28-residuals-validation/scopes/scope-06-r6-fm-residuals/FM07-DECISION.md`](../../docs/phases/phase-28-residuals-validation/scopes/scope-06-r6-fm-residuals/FM07-DECISION.md).

**Write-before-export (FM-02):** Agents must record ≥1 discovery OR ≥1 decision (linked to `TRACE_TASK_ID`) **before** `seed export --strict --enforce`. Do not export first and backfill. Plain `seed export` (no `--strict`) emits an early thin-graph stderr warn but still writes; `--strict --enforce` stays fail-closed.

**Parent / worker Trace (FM-04):** Multitask parent must set `TRACE_TASK_ID` before edits and before delegating; put task UUID + workspace in every worker prompt (env inheritance is not guaranteed). Parent owns gap pass / discoveries / decisions — do not offload graph-only work to workers while the parent edits without task. Option A hook denies empty task under `enforce=strict` per process; Trace cannot product-enforce Multitask inheritance (Option B deferred).

**Verify install before G1 session:**

```bash
grep -q "mandatory gap pass" runs/G1/.cursor/rules/trace-enforcement.mdc && echo OK
grep -q "Parent orchestrator" runs/G1/.cursor/rules/trace-enforcement.mdc && echo OK
```

## 2. B0 session

1. Open `runs/B0`
2. Paste `prompts/PROMPT-B0.md`
3. Single agent or Multitask — **build only**
4. `./score.sh B0 --test`

## 3. G1 Session A — build-only (critical)

1. Export env:

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
export PATH="/home/ali/Desktop/Trace/bin:$PATH"
export TRACE_TASK_ID=e0200000-0000-4000-8000-000000000010
export TRACE_PROJECT_ROOT=/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/G1
```

2. Open `runs/G1`
3. Paste `prompts/PROMPT-G1-BUILD.md` — **build only; no gap-analysis instruction**
4. If Multitask: orchestrator must read installed rules (gap pass + Parent orchestrator FM-04); set `TRACE_TASK_ID` before edits/delegation; put task UUID in every worker prompt (env inheritance not guaranteed)
5. **Do not** send a directed gap message before scoring Session A
6. Score build arm: `./score.sh G1 --p25 --arm build` (or `./score.sh G1 --p25` — same default)

**P25-3a:** Thin graph (`discoveries=0 decisions=0`) is an **expected FAIL** on build-only — documents install wiring without agent gap-pass behavior. See [RUBRIC.md](RUBRIC.md).

## 4. G1 Session B — directed-gap (optional confirmation)

Run **after** Session A is scored (or when operator wants P25-C behavior validation).

1. Same workspace (`runs/G1`) — **do not** re-run `./prepare.sh G1` unless invalidating Session A
2. Paste `prompts/PROMPT-G1-DIRECTED-GAP.md` — single directed message: run the mandatory gap pass from installed Trace rules
3. Agent may re-export during session; score uses the graph at `trace/graph.json`
4. Score directed arm: `./score.sh G1 --p25 --arm directed`

**P25-3b:** Same threshold as P25-3a (`discoveries≥1 OR decisions≥1`) but **PASS required** for gap-pass validation.

**P25-4 attestation:** Set arm-matched env **before** `./score.sh` so the harness can pass the check:

| Arm | Env | Meaning |
|-----|-----|---------|
| `--arm build` | `P25_ATTEST_BUILD=Y` | No human gap prompt before this score |
| `--arm directed` | `P25_ATTEST_DIRECTED=Y` | Gap prompt sent before this score |

Unset for the current arm → `score.sh` **skips** P25-4 (backward compatible). Wrong-arm env is ignored (e.g. `P25_ATTEST_DIRECTED=Y` while scoring `--arm build` does not pass build attestation).

Still record the human narrative in `experiments/RESULTS.md` (build: no gap prompt; directed: gap prompt sent). Example:

```bash
P25_ATTEST_BUILD=Y ./score.sh G1 --p25 --arm build
P25_ATTEST_DIRECTED=Y ./score.sh G1 --p25 --arm directed
```

### Invalidation rule

If the operator sends the **directed gap prompt** (`PROMPT-G1-DIRECTED-GAP.md`) **before** scoring Session A, the build arm is **invalidated**. Re-run with fresh `./prepare.sh G1` and repeat Session A only before any directed message.

## 5. Phase 26 gate

Record outcome in `experiments/RESULTS.md` using [RUBRIC.md](RUBRIC.md) decision table. Human promotes Phase 26 only if E02 null hypothesis is not rejected.

## Stopping conditions

- Agent cannot pass tests in 3 retries → partial record
- Operator accidentally sends gap prompt before build score → invalidate G1 build-only arm; re-run with fresh `./prepare.sh G1`
