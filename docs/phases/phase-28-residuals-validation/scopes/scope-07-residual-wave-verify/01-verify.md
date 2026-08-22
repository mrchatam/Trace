# P28-S07-01 — Residual-wave VERIFY

## Metadata
- id: P28-S07-01
- todo_ids: [P28-S07-01]
- role: verify
- skills: [incremental-implementation, documentation-and-adrs, test-driven-development]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated
- hooks: []

## Objective

Aggregate evidence that FR-P28-01…07 residual-wave work is green after S06 reviews APPROVE. Capture command outputs, spot-check FM Notes/decision artifacts, write **`VERIFY-NOTES-RESIDUAL-WAVE.md`** in this folder. Keep S05 `VERIFY-NOTES.md` **immutable**. **Does not** close DR-HANDOFF (P28-S07-02 owns). **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — S07-00 locks (FINAL)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [../scope-06-r6-fm-residuals/SCOPE-TODOS.md](../scope-06-r6-fm-residuals/SCOPE-TODOS.md)
- FR evidence under `../scope-06-r6-fm-residuals/`: `FM01-NOTES.md` … `FM10-NOTES.md`, `FM07-DECISION.md`
- [../scope-05-verify/VERIFY-NOTES.md](../scope-05-verify/VERIFY-NOTES.md) — S05 baseline (**do not edit**)
- [../scope-02-session-b-dogfood/SESSION-A-GRAPH-SNAPSHOT.json](../scope-02-session-b-dogfood/SESSION-A-GRAPH-SNAPSHOT.json) — thin baseline docs
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — Residual wave stays **OPEN** until S07-02
- [experiments/ab-p25-gap-pass-validation/](../../../../../../experiments/ab-p25-gap-pass-validation/) — optional directed score only (**no** `prepare.sh`)
- FM09 archive (if present): `experiments/runs/2026-08-20-p28-s06-11-fm09/evidence/`

## Session start

Follow agent-loop-protocol. Unattended: do not stop after planning. This row runs residual-wave verification and records evidence; it does **not** close DR-HANDOFF or decide successor.

## Locked defaults (FINAL — S07-00)

| Item | Value |
|------|-------|
| Precondition | P28-S06-02 … P28-S06-14 all `done` with APPROVE (or explicit skip with reason) |
| Binary | Rebuild `bin/trace` from repo HEAD before optional harness |
| **G1 wipe** | **FORBIDDEN** — never `./prepare.sh` / `./prepare.sh G1` |
| Dual-lane | Do **not** conflate thin (SESSION-A snapshot disc=0/dec=0) with rich (live G1 / FM09 evidence) |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p28-s07-01-verify/evidence/` |
| Notes artifact | `scopes/scope-07-residual-wave-verify/VERIFY-NOTES-RESIDUAL-WAVE.md` (**required**) |
| S05 VERIFY-NOTES | **Immutable** — never rewrite |
| Product Go | **Forbidden** |
| DR-HANDOFF | Residual wave stays **OPEN** — S07-02 closes |
| Successor | **Out of scope** — S07-02 only |
| Deferred (non-blocking) | FR-P28-D1…D4, X1 remain deferred (list in notes; do not fail VERIFY for them) |

## Locked verify command floor

Run from repo root unless noted. Tee outputs into evidence dir.

### Block 0 — Evidence dir

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p28-s07-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P28-S07-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "wave=residual-wave-S06-FR-P28-01-07"
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

**Pass:** both exit 0. Confirms Option A still green (strict + empty `TRACE_TASK_ID` → deny; non-strict → allow; INT-11 drift).

### Block 4 — FR evidence spot-check (required)

Confirm each artifact exists and board Notes for the matching review row say **APPROVE**. Record paths in VERIFY-NOTES.

| FR | FM | Artifact (under `scopes/scope-06-r6-fm-residuals/`) | Review row |
|----|-----|------------------------------------------------------|------------|
| FR-P28-01 | FM-01 | `FM01-NOTES.md` | P28-S06-02 |
| FR-P28-02 | FM-02 | `FM02-NOTES.md` | P28-S06-04 |
| FR-P28-03 | FM-04 | `FM04-NOTES.md` | P28-S06-06 |
| FR-P28-04 | FM-07 | `FM07-DECISION.md` (remain warn-only) | P28-S06-08 |
| FR-P28-05 | FM-08 | `FM08-NOTES.md` | P28-S06-10 |
| FR-P28-06 | FM-09 | `FM09-NOTES.md` + dual-lane archive (no prepare) | P28-S06-12 |
| FR-P28-07 | FM-10 | `FM10-NOTES.md` | P28-S06-14 |

```bash
SCOPE6=docs/phases/phase-28-residuals-validation/scopes/scope-06-r6-fm-residuals
for f in FM01-NOTES.md FM02-NOTES.md FM04-NOTES.md FM07-DECISION.md FM08-NOTES.md FM09-NOTES.md FM10-NOTES.md; do
  test -f "$SCOPE6/$f" || { echo "MISSING $f"; exit 1; }
done
# Spot-read FM09 for: thin via SESSION-A snapshot; directed PASS; rich build labeled post-directed; prepare.sh NOT RUN
```

**Pass:** all seven artifacts present; FM09 notes do not claim prepare wipe; board reviews APPROVE.

### Block 5 — Optional directed score smoke (no wipe)

Only if FM09 evidence needs refresh or live re-check. Prefer citing existing FM09 archive when still valid.

```bash
# FORBIDDEN: ./prepare.sh
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd /home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation
P25_ATTEST_DIRECTED=Y ./score.sh G1 --p25 --arm directed --test 2>&1 | tee "/home/ali/Desktop/Trace/$EVID/score-directed-optional.txt"
```

**If run — Pass:** `p25 arm: directed`; P25-3b PASS preferred.  
**If skipped — Pass:** VERIFY-NOTES cite FM09 evidence path + P28-S06-12 APPROVE.

**Never** treat rich live build as thin Session-A FAIL. Thin = SESSION-A snapshot / FM09 thin lane only.

## VERIFY-NOTES-RESIDUAL-WAVE.md template (required)

Write to `scopes/scope-07-residual-wave-verify/VERIFY-NOTES-RESIDUAL-WAVE.md`:

```markdown
# Phase 28 residual-wave VERIFY notes

**Run:** YYYY-MM-DD  
**Row:** P28-S07-01  
**Evidence dir:** experiments/runs/YYYY-MM-DD-p28-s07-01-verify/evidence/  
**Git SHA:** `<rev-parse HEAD or unknown>`  
**S05 VERIFY-NOTES:** immutable (not rewritten)

## Verdict

PASS | FAIL — confidence high | medium | low

## Per-block

| Block | Result | Notes |
|-------|--------|-------|
| 1 Build | PASS/FAIL | … |
| 2 Unit + cmd | … | … |
| 3 Install + hook smoke | … | Option A deny intact |
| 4 FR FM01–FM10 artifacts | … | all present; reviews APPROVE |
| 5 Directed score (optional) | PASS/SKIP | cite FM09 archive if skipped |

## FR-P28-01…07 disposition

| FR | FM | Disposition | Evidence |
|----|-----|-------------|---------|
| FR-P28-01 | FM-01 | closed | FM01-NOTES + S06-02 APPROVE |
| FR-P28-02 | FM-02 | closed | FM02-NOTES + S06-04 APPROVE |
| FR-P28-03 | FM-04 | closed | FM04-NOTES + S06-06 APPROVE |
| FR-P28-04 | FM-07 | closed (warn-only by design) | FM07-DECISION + S06-08 APPROVE |
| FR-P28-05 | FM-08 | closed | FM08-NOTES + S06-10 APPROVE |
| FR-P28-06 | FM-09 | closed | FM09 dual-lane + S06-12 APPROVE |
| FR-P28-07 | FM-10 | closed | FM10-NOTES + S06-14 APPROVE |

## Dual-lane (do not conflate)

| Lane | Source | Expected |
|------|--------|----------|
| Thin | SESSION-A snapshot / FM09 thin | disc=0/dec=0 |
| Directed rich | FM09 / optional re-score | P25-3b PASS |
| Build rich | FM09 labeled **post-directed** | P25-3a PASS — not Session-A thin FAIL |
| prepare.sh | **NOT RUN** | wipe forbidden |

## Deferred (non-blocking — remain open)

| ID | Topic |
|----|-------|
| FR-P28-D1 | Autonomous discovery→task spawn |
| FR-P28-D2 | Full Graphiti / temporal invalidation |
| FR-P28-D3 | RESULTS.md parser for P25-4 |
| FR-P28-D4 | Hook Option B |
| FR-P28-X1 | Daemon / HTTP / hosted MCP |

## Gaps / spawn

(none) | list for S07-02

## DR-HANDOFF status

**Residual wave OPEN** — S07-02 closes Residual wave section only (S05 CLOSED history intact).
```

## Pass / fail policy

| Condition | Row status |
|-----------|------------|
| Blocks 1–4 PASS; FR table closed | `done` |
| Unit / cmd / install / hook FAIL | `failed` — note for S07-02 repair spawn |
| Missing FR Notes/decision artifact | `failed` |
| FM09 claims prepare wipe or conflates thin/rich | `failed` — evidence tainted |
| G1 missing for optional score | SKIP Block 5; do **not** run prepare.sh; still `done` if 1–4 PASS |
| D1–D4/X1 still deferred | **PASS** — do not fail VERIFY |

## Todo updates

Status + notes on **P28-S07-01** only.

## Exit criteria

- [ ] Evidence dir populated with build/unit/cmd/install/hook (+ optional score)
- [ ] `VERIFY-NOTES-RESIDUAL-WAVE.md` complete (blocks + FR table + deferred)
- [ ] No `./prepare.sh` run; dual-lane not conflated
- [ ] S05 `VERIFY-NOTES.md` untouched
- [ ] Residual wave DR-HANDOFF still **OPEN**
- [ ] Board Notes cite verdict + evidence path

## Minimal todos

- [ ] Block 0–3: evidence dir, build, unit+cmd, install+hook smoke
- [ ] Block 4: spot-check FM01–FM10 / FM07-DECISION vs board APPROVE
- [ ] Block 5: optional directed score or cite FM09 archive
- [ ] Write VERIFY-NOTES-RESIDUAL-WAVE.md
- [ ] Mark P28-S07-01 `done` / `failed` / `blocked`

## Next

**P28-S07-02**
