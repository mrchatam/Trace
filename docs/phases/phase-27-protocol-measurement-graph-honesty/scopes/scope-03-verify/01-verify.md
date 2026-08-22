# P27-S03-01 — Phase 27 VERIFY

## Metadata
- id: P27-S03-01
- todo_ids: [P27-S03-01]
- role: verify
- skills: [incremental-implementation, documentation-and-adrs]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed
- hooks: []

## Objective

Run the locked verify floor for Phase 27 **INT-07 / INT-08 / INT-10** deliverables. **Own** the harness upgrade: `score.sh` T02 `--strict` warn-only → **`--strict --enforce` on both arms**. Archive evidence; write **`VERIFY-NOTES.md`**. **Does not** close DR-HANDOFF (P27-S03-02 owns).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — S03-00 locks (FINAL)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S03-02
- [experiments/ab-p25-gap-pass-validation/](../../../../../../experiments/ab-p25-gap-pass-validation/) — `prepare.sh`, `score.sh`, `RUBRIC.md`, `PROTOCOL.md`
- [cmd/trace/testdata/p26-export-snippet.json](../../../../../../cmd/trace/testdata/p26-export-snippet.json)
- Pattern: [P26 S05-01](../../../phase-26-loop-implementation/scopes/scope-05-verify/01-verify.md)

## Session start

Follow agent-loop-protocol. Unattended: do not stop after planning. This row runs verification + harness enforce upgrade and records evidence; it does **not** close DR-HANDOFF or decide successor.

## Locked defaults (FINAL — S03-00)

| Item | Value |
|------|-------|
| Precondition | P27-S01-02 + P27-S02-02 both `done` |
| Binary | Rebuild `bin/trace` from repo HEAD before any harness/product demo |
| Unit tests | `go test ./internal/... -count=1` — **must PASS** |
| Cmd honesty tests | `go test ./cmd/trace/... -run 'SeedExport\|Enforce\|Strict' -count=1` — **must PASS** |
| Harness `--enforce` | **This row owns** — upgrade `score.sh` T02 to `--strict --enforce` on **both** arms |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p27-s03-01-verify/evidence/` |
| Notes artifact | `scopes/scope-03-verify/VERIFY-NOTES.md` (**required**) |
| Product Go | **Forbidden** except ≤5-line doc/comment typo; honesty logic already shipped in S02 |
| Experiments allow | **`score.sh` (required enforce upgrade)**; optional RUBRIC/PROTOCOL one-line note if enforce semantics need clarity |
| DR-HANDOFF | Stays **OPEN** — S03-02 closes |
| Successor | **Out of scope** — S03-02 only |

## Deliverable map (INT-07 / 08 / 10)

| Theme | INT | Primary evidence | Locked check |
|-------|-----|------------------|--------------|
| Graph honesty | INT-07 | cmd tests + thin demo | Thin `--strict --enforce` blocked; clean fixture still passes tests |
| Protocol / score | INT-08 | `score.sh` preflight + enforce | Preflight export; T02 `--strict --enforce` fails G2 on thin |
| Two-session rubric | INT-10 | `--arm build\|directed` | P25-3a vs P25-3b labels; build thin FAIL expected |

## Pass / fail policy (VERIFY row)

### VERIFY **PASS** (phase-closeable after S03-02) when all hold

1. Rebuild + `go test ./internal/...` PASS
2. Cmd honesty/enforce tests PASS (incl. P26 thin block + warn-only write)
3. Product thin-graph demo PASS (Block 3)
4. **`score.sh` enforce upgrade landed** (T02 uses `--strict --enforce`; thin build produces G2 **FAIL**, not WARN-only)
5. `score.sh G1 --p25 --arm build`: **P25-3a** labeled; thin → **P25-3a FAIL expected**; P25-1/2 **PASS**
6. `score.sh G1 --p25 --arm directed`: **P25-3b** labeled (rich PASS only if Session-B graph exists; thin FAIL OK)
7. Evidence dir + `VERIFY-NOTES.md` complete; DR-HANDOFF still **OPEN**

### Expected residuals (**do not** FAIL VERIFY)

| Residual | Disposition |
|----------|-------------|
| P25-3a FAIL on build-only thin graph | **Expected** (RUBRIC baseline) |
| G2 FAIL after enforce on thin build | **Expected** WARN→FAIL after S03-01 upgrade |
| P25-3b FAIL without Session-B dogfood | **OK** — protocol machinery verified by label + arm isolation |
| P25-4 operator attestation | Manual; `skip` in score.sh — non-blocking |
| FM-07 SHA drift WARN | Warn-only — never fails alone |
| S02-02 BLOCKING duplicate orphan message | Low residual — document only |
| Overall `VERDICT FAIL` on build arm | **Acceptable** when caused only by expected thin/enforce residuals |

### VERIFY **FAIL** (blocking — spawn forward or mark `failed`)

| Condition | Action |
|-----------|--------|
| Unit or honesty tests FAIL | `failed` or spawn repair |
| Thin graph **passes** `--strict --enforce` | Product regression — spawn; do not close phase |
| `score.sh` still warn-only `--strict` (no `--enforce`) | Incomplete — finish upgrade before `done` |
| P25-1 or P25-2 FAIL | Install regression — spawn |
| Arms conflated (no P25-3a/3b / invalid `--arm`) | Protocol regression |
| DR-HANDOFF closed by this row | Protocol violation — revert status to OPEN |

## Locked verify command floor

Run from repo root unless noted. Capture stdout/stderr in evidence dir.

### Block 0 — Harness `--enforce` upgrade (required first)

Edit `experiments/ab-p25-gap-pass-validation/score.sh` T02 (today warn-only around the `--strict` preflight):

```bash
# Target behavior (both arms):
#   "$TRACE_BIN" -C "$WS" seed export -o trace/graph.json --strict --enforce
# On failure: fail "G2 …" (or named honesty check) — do NOT only WARN
# Keep FM-07 warn-only; keep P25-3a/3b count checks
```

**Pass:** `rg -n 'strict --enforce|--enforce' experiments/ab-p25-gap-pass-validation/score.sh` shows T02 path; comment updated (remove “until S02”).

### Block 1 — Rebuild + unit + honesty tests (required)

```bash
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
go test ./internal/... -count=1
go test ./cmd/trace/... -run 'SeedExport|Enforce|Strict' -count=1
```

**Pass:** exit 0 on all three.

### Block 2 — Product thin-graph enforce demo (required)

Reproduce S02 demo with P26-class thin graph (import snippet or use existing test fixture pattern):

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
# Prefer: go test already covers; plus one live CLI demo:
# 1) Import/seed thin workspace OR use cmd testdata path via temp import
# 2) $TRACE_BIN -C "$THIN_WS" seed export -o /tmp/p27-thin-out.json --strict --enforce
#    → non-zero exit; no output file; stderr mentions discoveries=0 decisions=0 (honesty)
# 3) $TRACE_BIN -C "$THIN_WS" seed export -o /tmp/p27-thin-warn.json --strict
#    → exit 0; file written; stderr honesty WARN lines
```

Document exact commands used in VERIFY-NOTES. **Pass:** enforce blocks; strict-only writes.

### Block 3 — Harness score dry runs (required)

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd /home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation

# Ensure G1 workspace exists (prepare if needed)
./prepare.sh G1 2>&1 | tee /tmp/p27-s03-prepare.txt

# Build arm (default) — thin expected
./score.sh G1 --p25 --arm build 2>&1 | tee /tmp/p27-s03-score-build.txt

# Directed arm — label check; rich PASS only if Session-B graph present
./score.sh G1 --p25 --arm directed 2>&1 | tee /tmp/p27-s03-score-directed.txt

# Invalid arm
./score.sh G1 --p25 --arm foo; echo "exit=$?"  # expect exit 2
```

**Pass thresholds:**

| Check | Target |
|-------|--------|
| P25-1 / P25-2 | **PASS** both arms |
| Build: `p25 arm: build` + **P25-3a** | FAIL on thin = **expected** |
| Build: G2 / enforce honesty | **FAIL** on thin after upgrade (not WARN-only) |
| Directed: `p25 arm: directed` + **P25-3b** | Label present; PASS iff rich graph |
| Invalid `--arm` | exit 2 |
| Overall VERDICT | May be **FAIL** on thin build — acceptable per policy above |

Copy score tees into evidence dir.

### Block 4 — Evidence archive (required)

```bash
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p27-s03-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P27-S03-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "trace_bin=$TRACE_BIN"
  echo "enforce_upgrade=yes"
} > "$EVID/99-run-metadata.txt"
# Copy: score-build, score-directed, prepare, thin demo snippets, spot-checks.txt, score.sh diff hunk
```

### Block 5 — Optional RESULTS.md note

Add a short row to `experiments/RESULTS.md` for Phase 27 verify (date + arms + expected thin FAIL + enforce note). Non-blocking if file busy; prefer update.

## VERIFY-NOTES.md template (required)

Write to `scopes/scope-03-verify/VERIFY-NOTES.md`:

```markdown
# Phase 27 VERIFY notes

**Run:** YYYY-MM-DD  
**Row:** P27-S03-01  
**Evidence dir:** experiments/runs/YYYY-MM-DD-p27-s03-01-verify/evidence/  
**Git SHA:** `<rev-parse HEAD>`  
**Harness:** score.sh T02 `--strict --enforce` (both arms)

## Verdict

PASS | FAIL — confidence high | medium | low

## Closure signals (Phase 27)

| Check | Target | This run |
|-------|--------|----------|
| INT-07 product honesty | Thin enforce blocked | PASS/FAIL |
| INT-08 score enforce | G2 FAIL on thin (not WARN-only) | PASS/FAIL |
| INT-10 arms | P25-3a/3b labels + arm isolation | PASS/FAIL |
| Unit + cmd tests | PASS | PASS/FAIL |
| P25-1 / P25-2 | PASS | PASS/FAIL |

## Harness score (--p25)

| ID / arm | Result | Notes |
|----------|--------|-------|
| P25-1 | … | … |
| P25-2 | … | … |
| P25-3a (build) | FAIL expected if thin | discoveries=… decisions=… |
| P25-3b (directed) | … | rich only if Session-B |
| G2 enforce | FAIL expected if thin | WARN→FAIL after upgrade |
| FM-07 | WARN/PASS/SKIP | warn-only |

## Residuals (non-blocking)

| Topic | Disposition |
|-------|-------------|
| Build thin P25-3a FAIL | Expected measurement baseline |
| Thin G2 enforce FAIL | Expected after S03-01 |
| Directed without Session-B | P25-3b FAIL OK |
| P25-4 attestation | Manual / skip |
| S02-02 BLOCKING duplicate msg | Note only |
| Overall VERDICT FAIL (thin) | Acceptable if closure signals PASS |

## Gaps / spawn

(none | P27-S03-01a implement + 01b review)

## DR-HANDOFF status

**OPEN** — S03-02 closes with successor decision (never TBD).
```

## Do not

- Close [DR-HANDOFF.md](../../DR-HANDOFF.md) — S03-02 owns
- Decide Phase 28 / `no successor` — S03-02 only
- Mark VERIFY PASS if thin graph still passes product `--strict --enforce`
- Leave `score.sh` on warn-only `--strict`
- Implement new product honesty features (spawn forward on regression)
- Rewrite S01–S02 `done` history
- Require Session-B dogfood or P25-3b PASS to mark this row `done`

## Exit criteria

- [ ] `score.sh` upgraded to `--strict --enforce` (both arms); thin → G2 FAIL not WARN-only
- [ ] `bin/trace` rebuilt; `go test ./internal/...` PASS
- [ ] Cmd honesty/enforce tests PASS
- [ ] Product thin-graph enforce demo PASS
- [ ] Build + directed score dry runs archived; P25-3a/3b labels confirmed
- [ ] P25-1/2 PASS; expected thin residuals documented
- [ ] Evidence dir + `VERIFY-NOTES.md` complete
- [ ] Board Notes: commands, verdict, evidence path, enforce upgrade summary
- [ ] DR-HANDOFF remains **OPEN**

## Minimal todos

- [ ] Preflight: confirm S01-02 + S02-02 `done`
- [ ] Block 0: upgrade `score.sh` T02 to `--strict --enforce`
- [ ] Block 1: build + unit + honesty tests
- [ ] Block 2: product thin demo
- [ ] Block 3: score build + directed + invalid arm
- [ ] Block 4–5: archive + VERIFY-NOTES (+ optional RESULTS.md)
- [ ] Set row `done` or `failed` with evidence

## Next

**P27-S03-02**
