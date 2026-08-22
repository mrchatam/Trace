# P26-S05-01 — VERIFY (E03 re-score)

## Metadata
- id: P26-S05-01
- todo_ids: [P26-S05-01]
- role: verify
- skills: [incremental-implementation, documentation-and-adrs]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed
- hooks: []

## Objective

Run post-implementation verification for Phase 26 deliverables **D1–D6** (PLAN.md VERIFY mapping). **Closure signal:** E02 **P25-2 PASS** (was FAIL). Archive evidence; write **`VERIFY-NOTES.md`**. **Does not** close DR-HANDOFF (P26-S05-02 owns).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — S05-00 locks (FINAL)
- [PLAN.md](../scope-01-planning/PLAN.md) — VERIFY mapping D1–D6
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S05-02
- [experiments/ab-p25-gap-pass-validation/](../../../../../../experiments/ab-p25-gap-pass-validation/) — `prepare.sh`, `score.sh`, `RUBRIC.md`
- [experiments/RESULTS.md](../../../../../../experiments/RESULTS.md)

## Session start

Follow agent-loop-protocol. Unattended: do not stop after planning. This row runs verification and records evidence; it does not close DR-HANDOFF or decide successor.

## Locked defaults (FINAL — S05-00)

| Item | Value |
|------|-------|
| Precondition | P26-S02-02, P26-S03-02, P26-S04-02 all `done` |
| Binary | Rebuild `bin/trace` from repo HEAD before any harness step |
| **Option A (preferred)** | `prepare.sh G1` + `score.sh G1 --p25` with `TRACE_BIN=…/bin/trace` |
| **Option B (partial fallback)** | Temp-dir `trace install cursor --write` greps only — use when G1 workspace missing/stale and `prepare.sh G1` blocked |
| Closure signal | **P25-2 PASS** (E02 recorded FAIL — `ParentOrchestratorRule` unwired) |
| Unit tests | `go test ./internal/...` — **must PASS** |
| Schema embed | Expect migration **028** (`TestMigrateBackupAuthCLI` / compat) |
| G1 verify task | `e0200000-0000-4000-8000-000000000050` (seed task 5) |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p26-s05-01-verify/evidence/` |
| Notes artifact | `scopes/scope-05-verify/VERIFY-NOTES.md` (**required**) |
| Results row | Update `experiments/RESULTS.md` (E02 re-score or new E03 row with date + verdict) |
| Product Go | **Forbidden** except blocking doc typo ≤5 lines |
| DR-HANDOFF | Stays **OPEN** — S05-02 closes |
| Successor | **Out of scope** — S05-02 only |

### Option selection (locked)

| Condition | Use |
|-----------|-----|
| `bin/trace` builds; `prepare.sh G1` succeeds | **Option A** — full harness |
| G1 missing `.trace/` / no `trace/graph.json`; `prepare.sh` blocked | **Option B** for P25-1/2 only; **document P25-3 SKIP** in VERIFY-NOTES + board Notes |
| Option B only | Row may be `done` only if unit tests PASS **and** P25-1/2 greps PASS **and** residuals explicit; prefer `failed` if Option A was achievable |

## Deliverable map (D1–D6)

| D | Description | Primary evidence | Locked command |
|---|-------------|------------------|----------------|
| **D1** | `ParentOrchestratorRule` in MDC output | S04 install | `go test ./internal/install/... -count=1`; `rg ParentOrchestratorRule internal/install/` |
| **D2** | BLOCKING discovery → task path | S02 promotion | `go test ./internal/loop/... ./internal/mcp/... ./internal/domain/... -count=1 -run 'Promote|spawned_task|promotion_candidate'` |
| **D3** | Greenfield no sticky STOP at hop 1 | S03 saturation | G1 task `…0050`: one empty apply → `stopped=false` (or not `p19_saturated`) |
| **D4** | Reset API clears STOP | S03 reset | `trace loop reset --task …0050` after saturation → `stopped=false`, `hop_count=0` |
| **D5** | Unified STOP reason | S03 INT-09 | After saturation: gate `reason_code` == export `stop_reason` (`p19_saturated`) |
| **D6** | Full PASS | all | `go test ./internal/...`; Option A: `score.sh G1 --p25` VERDICT PASS |

**Non-blocking residual (document, do not fail D6):** INT-04 hook may still be permissive when `TRACE_TASK_ID` unset — VERIFY scores **text presence** only (PLAN cross-scope risk).

## Locked verify command floor

Run from repo root unless noted. Capture stdout/stderr snippets in evidence dir (`spot-checks.txt`).

### Block 1 — Rebuild + unit tests (required)

```bash
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
go test ./internal/... -count=1
```

**Pass:** exit 0 on both.

### Block 2 — D1 install regression (required)

```bash
go test ./internal/install/... -count=1
rg -n 'ParentOrchestratorRule' internal/install/enforcement.go
```

**Pass:** tests PASS; `rg` shows definition + usage inside `cursorRulesMDCContent` and `claudeFallbackRulesContent`.

### Block 3 — Option A harness (preferred)

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
cd /home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation
./prepare.sh G1
./score.sh G1 --p25 2>&1 | tee /tmp/p26-s05-score.txt
```

**Pass thresholds:**

| Check | Target |
|-------|--------|
| P25-1 GapPassPrompt | PASS |
| P25-2 Parent orchestrator | **PASS** (closure signal) |
| P25-3 graph richness | PASS (or document residual if graph from prior E02 session) |
| G2/G3/E1–E3 harness | PASS |
| VERDICT | **PASS** |

Copy `/tmp/p26-s05-score.txt` into evidence dir.

### Block 4 — Option B partial fallback (only if Block 3 blocked)

```bash
cd /home/ali/Desktop/Trace
tmpdir=$(mktemp -d)
./bin/trace -C "$tmpdir" init
./bin/trace -C "$tmpdir" install cursor --write
grep -qi 'mandatory gap pass' "$tmpdir/.cursor/rules/trace-enforcement.mdc" && echo P25-1_PASS
grep -qi 'Parent orchestrator' "$tmpdir/.cursor/rules/trace-enforcement.mdc" && echo P25-2_PASS
rm -rf "$tmpdir"
```

**Pass:** both greps succeed. **Must** note `P25-3 SKIP` and reason in VERIFY-NOTES.

### Block 5 — D3–D5 loop spot-checks (after Option A prepare, or fresh G1)

Use G1 workspace `experiments/ab-p25-gap-pass-validation/runs/G1` and task `e0200000-0000-4000-8000-000000000050`.

```bash
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
WS=/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/G1
TASK=e0200000-0000-4000-8000-000000000050

# D3 — first empty apply must not sticky-STOP (may still block edit for plan_critiqued)
"$TRACE_BIN" -C "$WS" loop apply --task "$TASK" --writes '{}' 2>/dev/null || true
"$TRACE_BIN" -C "$WS" loop status --task "$TASK" | tee /tmp/p26-d3-status.json

# D4 — saturate then reset (two consecutive empties if needed)
"$TRACE_BIN" -C "$WS" loop apply --task "$TASK" --writes '{}' 2>/dev/null || true
"$TRACE_BIN" -C "$WS" loop reset --task "$TASK"
"$TRACE_BIN" -C "$WS" loop status --task "$TASK" | tee /tmp/p26-d4-reset.json

# D5 — gate vs export stop_reason alignment after saturation STOP
"$TRACE_BIN" -C "$WS" loop gate --task "$TASK" --for edit 2>/tmp/p26-gate.json || true
"$TRACE_BIN" -C "$WS" seed export -o /tmp/p26-export.json
# Compare gate reason_code vs export deliberation_states[].stop_reason for $TASK
```

**Pass (D3):** after **one** pure-empty apply, not stopped with `p19_saturated` (second empty may STOP — expected).

**Pass (D4):** after `loop reset`, `stopped=false`, `hop_count=0` (edit may still block on `plan_critiqued` — OK).

**Pass (D5):** when stopped from saturation, gate `reason_code` matches export `stop_reason` (`p19_saturated`).

Skip Block 5 only if G1 not prepared; document SKIP in VERIFY-NOTES.

### Block 6 — Evidence archive

```bash
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p26-s05-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P26-S05-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "trace_bin=$TRACE_BIN"
  echo "option=A|B"
} > "$EVID/99-run-metadata.txt"
# Copy score output, spot-check JSON, spot-checks.txt into $EVID
```

## VERIFY-NOTES.md template (required)

Write to `scopes/scope-05-verify/VERIFY-NOTES.md`:

```markdown
# Phase 26 VERIFY notes

**Run:** YYYY-MM-DD  
**Row:** P26-S05-01  
**Evidence dir:** experiments/runs/YYYY-MM-DD-p26-s05-01-verify/evidence/  
**Git SHA:** `<rev-parse HEAD>`  
**Option:** A | B

## Verdict

PASS | FAIL — confidence high | medium | low

## Closure signal

| Check | E02 (Phase 25) | This run |
|-------|----------------|----------|
| P25-2 Parent orchestrator | FAIL | PASS/FAIL |

## D1–D6 map

| D | Result | Evidence |
|---|--------|----------|
| D1 | PASS/FAIL/SKIP | … |
| D2 | … | … |
| D3 | … | … |
| D4 | … | … |
| D5 | … | … |
| D6 | … | score.sh / unit tests |

## P25 harness (--p25)

| ID | Result | Notes |
|----|--------|-------|
| P25-1 | … | … |
| P25-2 | … | **closure** |
| P25-3 | … | … |

## Residuals (non-blocking)

| Topic | Disposition |
|-------|-------------|
| INT-04 hook permissive | Text-only verify; hook hardening deferred |
| P25-4 operator attestation | Manual; not blocking |
| Option B partial | List what was skipped |

## Gaps / spawn

(none | P26-S05-01a implement + 01b review)

## DR-HANDOFF status

**OPEN** — S05-02 closes with successor decision.
```

## Update experiments/RESULTS.md

Add or update a row:

```markdown
| YYYY-MM-DD | E02-rescore / E03 | ab-p25-gap-pass-validation | **PASS/FAIL** | Phase 26 S05 verify; P25-2 was FAIL on E02; Option A/B; see VERIFY-NOTES |
```

## Do not

- Close [DR-HANDOFF.md](../../DR-HANDOFF.md) — S05-02 owns
- Decide Phase 27 successor — S05-02 only
- Mark P25-3 PASS from Option B alone without graph evidence
- Implement product features (spawn forward on regression)
- Rewrite S02–S04 `done` history

## Exit criteria

- [ ] `bin/trace` rebuilt; `go test ./internal/...` PASS
- [ ] D1 install regression PASS
- [ ] Option A run **or** Option B with documented SKIP residuals
- [ ] **P25-2 PASS** (closure signal)
- [ ] D3–D5 spot-checks PASS or documented SKIP with reason
- [ ] Evidence dir + `VERIFY-NOTES.md` complete
- [ ] `experiments/RESULTS.md` updated
- [ ] Board Notes: option used, commands, verdict, evidence path
- [ ] DR-HANDOFF remains **OPEN**

## Minimal todos

- [ ] Preflight: confirm S02/S03/S04 reviews `done`
- [ ] Block 1–2: build + unit + install tests
- [ ] Block 3 (preferred) or Block 4 (fallback)
- [ ] Block 5 loop spot-checks when G1 available
- [ ] Block 6 archive + VERIFY-NOTES + RESULTS.md
- [ ] Set row `done` or `failed` with evidence

## Next

**P26-S05-02**
