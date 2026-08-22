# P28-S05-01 — VERIFY (full regression)

## Metadata
- id: P28-S05-01
- todo_ids: [P28-S05-01]
- role: verify
- skills: [incremental-implementation, documentation-and-adrs]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed
- hooks: []

## Objective

Run Phase 28 full regression VERIFY after S00–S04. Capture evidence, write **`VERIFY-NOTES.md`** with per-block PASS/FAIL and final R1–R8 disposition. **Does not** close DR-HANDOFF (P28-S05-02 owns). **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — S05-00 locks (FINAL)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md)
- [TEST-MATRIX.md](../scope-01-integration-tests/TEST-MATRIX.md) — M-01..M-16
- [SESSION-A-GRAPH-SNAPSHOT.json](../scope-02-session-b-dogfood/SESSION-A-GRAPH-SNAPSHOT.json) — thin baseline
- [SESSION-B-SCORE.txt](../scope-02-session-b-dogfood/SESSION-B-SCORE.txt) — directed P25-3b
- [Phase 27 VERIFY-NOTES](../../../phase-27-protocol-measurement-graph-honesty/scopes/scope-03-verify/VERIFY-NOTES.md) — baseline
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S05-02
- [experiments/ab-p25-gap-pass-validation/](../../../../../../experiments/ab-p25-gap-pass-validation/) — `score.sh` only (**no** `prepare.sh`)
- [experiments/RESULTS.md](../../../../../../experiments/RESULTS.md)

## Session start

Follow agent-loop-protocol. Unattended: do not stop after planning. This row runs verification and records evidence; it does **not** close DR-HANDOFF or decide successor.

## Locked defaults (FINAL — S05-00)

| Item | Value |
|------|-------|
| Precondition | P28-S00-02 … P28-S04-02 all `done` |
| Binary | Rebuild `bin/trace` from repo HEAD before harness |
| **G1 wipe** | **FORBIDDEN** — never `./prepare.sh` / `./prepare.sh G1` |
| Score cwd | `experiments/ab-p25-gap-pass-validation/` with `TRACE_BIN=/home/ali/Desktop/Trace/bin/trace` |
| Live G1 state | **Rich** after Session-B (disc≥1 / dec≥1) — expected |
| Thin Session-A | Docs-only via `SESSION-A-GRAPH-SNAPSHOT.json` (disc=0/dec=0) |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p28-s05-01-verify/evidence/` |
| Notes artifact | `scopes/scope-05-verify/VERIFY-NOTES.md` (**required**) |
| Product Go | **Forbidden** |
| DR-HANDOFF | Stays **OPEN** — S05-02 closes |
| Successor | **Out of scope** — S05-02 only |

### Dual-lane G1 score strategy (FINAL — never wipe)

| Lane | Command / evidence | How to interpret |
|------|--------------------|------------------|
| **Directed (primary)** | `P25_ATTEST_DIRECTED=Y ./score.sh G1 --p25 --arm directed` on **live rich** G1 | Expect **P25-3b PASS**; P25-4 PASS via env. Regression if FAIL. |
| **Build live (required)** | `P25_ATTEST_BUILD=Y ./score.sh G1 --p25 --arm build` on **same rich** on-disk G1 | Expect **P25-3a PASS** (rich). Label in notes as **post-Session-B rich build-arm**. Not Session-A thin. |
| **Build thin baseline (docs)** | Cite `SESSION-A-GRAPH-SNAPSHOT.json` + prior E02/P27 Session-A build score | Historical thin P25-3a FAIL remains valid. Rich build PASS does **not** invalidate thin baseline. |

**Fail VERIFY for:** directed regression, unit/cmd/install/matrix FAIL, hook deny smoke FAIL.  
**Do not fail VERIFY because** rich `--arm build` no longer shows thin P25-3a FAIL.

## Locked verify command floor

Run from repo root unless noted. Tee outputs into evidence dir.

### Block 0 — Evidence dir

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p28-s05-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P28-S05-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
} > "$EVID/99-run-metadata.txt"
```

### Block 1 — Build (required)

```bash
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
```

**Pass:** exit 0.

### Block 2 — Unit + cmd (required)

```bash
GOPROXY=direct go test ./internal/... -count=1 2>&1 | tee "$EVID/unit.txt"
GOPROXY=direct go test ./cmd/trace/... -count=1 2>&1 | tee "$EVID/cmd.txt"
```

**Pass:** both exit 0.

### Block 3 — Install + hook smoke (required)

```bash
GOPROXY=direct go test ./internal/install/... -count=1 2>&1 | tee "$EVID/install.txt"
GOPROXY=direct go test ./internal/install/... -run 'CursorLoopGateFailClosed|HookDrift|CursorLoopGateAllowNonStrict' -count=1 2>&1 | tee "$EVID/hook-smoke.txt"
```

**Pass:** both exit 0. Confirms Option A: strict + empty `TRACE_TASK_ID` → deny; non-strict → allow; INT-11 drift tests.

### Block 4 — TEST-MATRIX / evals (required)

```bash
bash evals/p28-regression/score_arm_labels_test.sh 2>&1 | tee "$EVID/matrix-m16.txt"
# Document in VERIFY-NOTES: TEST-MATRIX.md M-01..M-16 still cited; M-16 = this smoke (not P25-3b dogfood)
```

**Pass:** bash smoke exit 0; notes confirm matrix path exists and M-01..M-16 listed.

### Block 5 — Score directed (primary harness)

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd /home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation
# FORBIDDEN: ./prepare.sh
P25_ATTEST_DIRECTED=Y ./score.sh G1 --p25 --arm directed 2>&1 | tee "/home/ali/Desktop/Trace/$EVID/score-directed.txt"
```

**Pass:** `p25 arm: directed`; **P25-3b PASS**; P25-1/2 PASS; G2 `--strict --enforce` PASS; P25-4 PASS (env); VERDICT PASS preferred. Compare to S02 `SESSION-B-SCORE.txt`.

### Block 6 — Score build (rich live + thin docs)

```bash
# Still in harness cwd; still NO prepare.sh
P25_ATTEST_BUILD=Y ./score.sh G1 --p25 --arm build 2>&1 | tee "/home/ali/Desktop/Trace/$EVID/score-build-rich.txt"
```

**Pass (live rich):** exit 0; `p25 arm: build`; **P25-3a PASS** expected on rich graph; P25-4 PASS via env. Record as **post-Session-B rich build-arm**.

**Thin baseline (required in VERIFY-NOTES, no command wipe):**

| Source | Expected |
|--------|----------|
| `../scope-02-session-b-dogfood/SESSION-A-GRAPH-SNAPSHOT.json` | discoveries=0 decisions=0 |
| Prior E02 / P27 Session-A build score | P25-3a thin FAIL (historical) |

### Block 7 — Archive + RESULTS

Copy score outputs already teed. Optionally add a Phase 28 VERIFY row to `experiments/RESULTS.md` if S02 did not already cover the verify re-score (do not overwrite `E02-SB`).

## VERIFY-NOTES.md template (required)

Write to `scopes/scope-05-verify/VERIFY-NOTES.md`:

```markdown
# Phase 28 VERIFY notes

**Run:** YYYY-MM-DD  
**Row:** P28-S05-01  
**Evidence dir:** experiments/runs/YYYY-MM-DD-p28-s05-01-verify/evidence/  
**Git SHA:** `<rev-parse HEAD or unknown>`  
**Harness:** score.sh G1 --p25 (dual-lane; no prepare.sh)

## Verdict

PASS | FAIL — confidence high | medium | low

## Per-block

| Block | Result | Notes |
|-------|--------|-------|
| 1 Build | PASS/FAIL | … |
| 2 Unit + cmd | … | … |
| 3 Install + hook smoke | … | Option A deny |
| 4 Matrix M-16 | … | TEST-MATRIX M-01..M-16 |
| 5 Score directed | … | P25-3b; P25_ATTEST_DIRECTED |
| 6 Score build rich | … | post-Session-B; P25-3a PASS expected |
| 6 Thin baseline docs | … | SESSION-A snapshot disc=0/dec=0 |

## Dual-lane G1

| Lane | Result | Interpretation |
|------|--------|----------------|
| Directed live | P25-3b … | Primary dogfood regression |
| Build live rich | P25-3a … | post-Session-B — not thin Session-A |
| Build thin (snapshot) | disc=0/dec=0 | Historical thin baseline; do not wipe |

## Residual register R1–R8 (final)

| ID | Final status | Evidence |
|----|--------------|----------|
| R1 | closed | Session-B / directed re-score |
| R2 | closed | Option A hook |
| R3 | closed | FM-05 script deny |
| R4 | closed | honesty single source |
| R5 | closed | P25_ATTEST_* |
| R6 | partial/deferred | FM gaps — list remaining |
| R7 | closed | TEST-MATRIX + this VERIFY |
| R8 | closed | hook_drift_test |

## vs Phase 27 VERIFY baseline

| Topic | P27 | P28 |
|-------|-----|-----|
| P25-3b | FAIL OK (no Session-B) | PASS (Session-B + re-score) |
| Hook failClosed | residual | Option A closed |
| Honesty dup | residual | R4 closed |
| P25-4 | skip/manual | env attestation |
| Build thin P25-3a | FAIL expected | snapshot docs + rich live PASS labeled |

## Gaps / spawn

(none) | list for S05-02

## DR-HANDOFF status

**OPEN** — S05-02 closes with successor decision (never TBD).
```

## Pass / fail policy

| Condition | Row status |
|-----------|------------|
| All required blocks PASS under dual-lane rules | `done` |
| Directed / unit / install / hook FAIL | `failed` — note for S05-02 repair spawn |
| G1 missing / no `.trace/` | `blocked` — write reason; **still no prepare.sh** |
| Rich build P25-3a PASS | **PASS** for Block 6 live — never treat as thin-baseline failure |

## Todo updates

Status + notes on **P28-S05-01** only.

## Exit criteria

- [ ] Evidence dir populated with build/unit/cmd/install/hook/matrix/score outputs
- [ ] `VERIFY-NOTES.md` complete (blocks + dual-lane + R1–R8)
- [ ] No `./prepare.sh` run
- [ ] DR-HANDOFF still **OPEN**
- [ ] Board Notes cite verdict + evidence path

## Minimal todos

- [ ] Block 0–4: build, tests, install, hook smoke, matrix
- [ ] Block 5: directed score with `P25_ATTEST_DIRECTED=Y`
- [ ] Block 6: rich build score with `P25_ATTEST_BUILD=Y` + snapshot thin docs
- [ ] Write VERIFY-NOTES.md; optional RESULTS row
- [ ] Mark P28-S05-01 `done` / `failed` / `blocked`

## Next

**P28-S05-02**
