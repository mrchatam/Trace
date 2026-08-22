# P28-S05-02 — DR-HANDOFF Phase 28 close

## Metadata
- id: P28-S05-02
- todo_ids: [P28-S05-02]
- role: reviewer
- skills: [code-review-and-quality, documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: mixed
- hooks: []

## Objective

Independent **fresh-session** review of S05-01 verify evidence. Re-run minimum spot-checks (do **not** trust Notes alone). **Close Phase 28 DR-HANDOFF** with explicit successor decision (**never TBD**). Scaffold Phase 29 **only** if human promotes. Phase 28 complete when this row is `done`.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- [00-PLANNER.md](00-PLANNER.md) — S05-00 locks
- [01-verify.md](01-verify.md) — locked verify floor + dual-lane score strategy
- [VERIFY-NOTES.md](VERIFY-NOTES.md) — produced by S05-01
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [docs/TODO.md](../../../../TODO.md)
- [AGENTS.md](../../../../../AGENTS.md)
- Pattern: [P27 S03-02](../../../phase-27-protocol-measurement-graph-honesty/scopes/scope-03-verify/02-dr-handoff.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S05-01 verifier. Unattended: execute review loop until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-05-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p28-s05-01-verify/evidence/` |
| Session-A thin | `scopes/scope-02-session-b-dogfood/SESSION-A-GRAPH-SNAPSHOT.json` |
| Session-B score | `scopes/scope-02-session-b-dogfood/SESSION-B-SCORE.txt` |
| TEST-MATRIX | `scopes/scope-01-integration-tests/TEST-MATRIX.md` |
| Phase handoff | `DR-HANDOFF.md` |

## Locked DR-HANDOFF close policy (FINAL — S05-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S05-01** — verify floor + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S05-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Closure prerequisite | S05-01 `done`; R1–R5 closed; dual-lane score + unit/install/hook green per VERIFY-NOTES |
| Default successor | **`no successor`** when R1–R5 closed + regression green |
| Phase 29 | **Human promote only** — do not invent a phase; scaffold only if human asks |
| Regression path | Spawn `P28-S05-02a` implement + `02b` review; **do not** close Phase 28 |
| Must not | Leave `Successor decision: TBD`; run `./prepare.sh`; rewrite S00–S04 `done` history; ship product in this row |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 28 board rows `done` |

### Successor decision table (locked — pick exactly one)

| Outcome (from S05-01 + independent spot-check) | Decision | Phase 29 scaffold? |
|------------------------------------------------|----------|-------------------|
| R1–R5 closed; unit/cmd/install/matrix/hook green; directed P25-3b PASS; rich build labeled correctly | **`no successor`** | **No** |
| Hook deny / honesty / attestation regression | **Do not close** — spawn repair | No |
| Directed P25-3b FAIL on re-score | **Do not close** — spawn repair | No |
| Unit/cmd/install FAIL on re-run | **Do not close** — spawn repair | No |
| VERIFY used `./prepare.sh` wipe | **Do not close** — spawn repair; evidence tainted | No |
| VERIFY PASS but human requests a new theme | **Phase 29** — theme named by human in Notes | **Yes** (minimal runnable) only if human promotes |
| Human named different successor before this row | Document in Notes with evidence | Per human |

**Never** leave successor as `TBD` when marking this row `done`. If blocked on repair, write successor as **`pending repair spawn`** (still not TBD).

### Independent spot-check floor (minimum)

```bash
cd /home/ali/Desktop/Trace
# Do NOT prepare.sh
GOPROXY=direct go test ./internal/install/... -run 'CursorLoopGateFailClosed|HookDrift' -count=1
test -f docs/phases/phase-28-residuals-validation/scopes/scope-05-verify/VERIFY-NOTES.md
test -f docs/phases/phase-28-residuals-validation/scopes/scope-02-session-b-dogfood/SESSION-A-GRAPH-SNAPSHOT.json
# Optional directed re-score (no wipe):
# export TRACE_BIN=$PWD/bin/trace
# cd experiments/ab-p25-gap-pass-validation && P25_ATTEST_DIRECTED=Y ./score.sh G1 --p25 --arm directed
```

Confirm VERIFY-NOTES dual-lane: rich build labeled post-Session-B; thin baseline cited from snapshot; no prepare wipe claimed.

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S00: Residual audit complete (`RESIDUAL-AUDIT.md`)
- [ ] S01: Integration test matrix implemented + reviewed
- [ ] S02: Session-B dogfood run + P25-3b scored
- [ ] S03: Hook failClosed hardening + reviewed
- [ ] S04: Product polish + reviewed
- [ ] S05: Full regression VERIFY + successor documented (**never TBD**)

### Residuals to list on close (non-blocking OK)

| Topic | Disposition |
|-------|-------------|
| R6 FM matrix gaps (partial) | Deferred measurement / future human theme |
| FM-07 warn-only | By design |
| Live G1 rich | Expected after Session-B; thin = SESSION-A snapshot only |
| Autonomous discovery→task spawn | Out of scope (project laws) |

### DR-HANDOFF.md update template (on APPROVE — default `no successor`)

```markdown
**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Closed | YYYY-MM-DD |
| Successor decision | **no successor** |
| Phase 28 outcome | R1–R5 closed; Session-B P25-3b PASS; Option A hook failClosed; honesty dedupe; P25_ATTEST_*; TEST-MATRIX + full regression VERIFY |
| Verify delta vs Phase 27 | P25-3b validated; hook/honesty/attestation residuals closed; dual-lane G1 score (no prepare wipe) |
| Residuals (non-blocking) | R6 FM partial; FM-07 warn-only |
| Forward | Human promotes Phase 29 only if needed |
```

If **Phase 29** (human only):

```markdown
| Successor decision | **Phase 29** — <explicit theme from human> |
```

If verify **failed**: keep DR-HANDOFF **OPEN**; spawn repair; successor = **`pending repair spawn`**.

## Role work

1. Fresh-session re-verify S05-01 evidence (spot-checks above).
2. Write `REVIEW-NOTES.md` in this folder (findings + confidence + successor pick).
3. On APPROVE: update `DR-HANDOFF.md` → CLOSED; tick scope checklist; set successor (never TBD).
4. Update `docs/TODO.md`: Phase 28 → done; orchestrator paste for idle / next human promote.
5. Update `AGENTS.md` current focus.
6. Ensure RESIDUAL-AUDIT R1–R8 final dispositions match VERIFY-NOTES (forward notes only if needed).

## Todo updates

Status + notes on **P28-S05-02** only; may spawn repair rows below this row if needed.

## Exit criteria

- [ ] Independent spot-check recorded in `REVIEW-NOTES.md`
- [ ] `DR-HANDOFF.md` CLOSED with successor **not** TBD
- [ ] `docs/TODO.md` + `AGENTS.md` updated
- [ ] All P28 board rows `done` (or repair spawn pending — then do not close)
- [ ] Confidence medium or high with evidence

## Minimal todos

- [ ] Spot-check VERIFY-NOTES + dual-lane + hook smoke
- [ ] Write REVIEW-NOTES.md
- [ ] Close DR-HANDOFF or spawn repair
- [ ] Update TODO.md + AGENTS.md
- [ ] Mark P28-S05-02 `done` / `failed` / `blocked`

## Next

Phase complete when this row `done` + DR-HANDOFF CLOSED (**no successor** default).
