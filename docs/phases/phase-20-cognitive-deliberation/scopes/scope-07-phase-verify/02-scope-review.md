# P20-S07-02 — Review Phase 20 verify + DR-HANDOFF close

## Metadata
- id: P20-S07-02
- todo_ids: [P20-S07-02]
- role: reviewer
- skills: [code-review-and-quality, documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective

Independent **fresh-session** review of S07-01 evidence; re-run locked verify floor (do **not** trust Notes alone); **close DR-HANDOFF** with explicit successor decision.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff (mandatory)
- [00-PLANNER.md](00-PLANNER.md) — FINAL evidence bar
- [01-verify.md](01-verify.md) — locked floor + Must checklist
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [COVERAGE.md](../../COVERAGE.md)
- Pattern: [P18 S04-02 REVIEW-NOTES](../../../phase-18-fts-clone-honesty/scopes/scope-04-phase-verify/02-scope-review.md)

## Session start

Follow agent-loop-protocol. **Fresh session** — reviewer must not be the S07-01 implementer. Unattended: execute review loop until blocker/high clear or spawned forward.

## Locked DR-HANDOFF close policy (FINAL — S07-00)

| Field | Locked value |
|-------|--------------|
| Who closes | **S07-02 only** (S07-01 gathers evidence; S07-00 locked policy) |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Successor decision | **`no successor`** |
| Rationale | Phase 20 delivers TRACE_THOUGHTPROCESS MVP on P19 loop. §16 Experiments + §18 Risk-adaptive verification remain **Future** (COVERAGE.md). Hosted MCP is **Later developments** (separate repo; not a board phase). Human may promote a forward phase later (as Phase 20 followed Phase 19) — that is **not** this close's job. |
| Must not | Leave `Successor decision: TBD`; rewrite Phase 19 historical `no successor`; scaffold Phase 21 from this row without human promotion |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF CLOSED |
| Next runnable after close | **none** (board queue empty until human promotes) |

### DR-HANDOFF.md update template (on APPROVE)

```markdown
**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Closed | YYYY-MM-DD |
| Successor decision | **no successor** |
| Phase 20 outcome | State-driven deliberation controller + S02–S05 artifacts; P19 loop extended additively v1; schema through 019 |
| Residuals (non-blocking) | Seed export omits P20 cognition tables; FTS sync for apply-created entities; non-tx upsert+event paths |
| Forward (human queue) | Portable graph for P20 entities; §16/§18 if promoted; hosted MCP separate repo |
```

If verify bar **fails**: keep DR-HANDOFF **OPEN**, spawn `P20-S07-02a` implement + `02b` review immediately below this row; successor stays **`no successor` intent** until pass.

## Review focus

Confirm independently:

1. **Verify floor** — all command blocks from [01-verify.md](01-verify.md) re-run PASS (session ≠ S07-01).
2. **Compat ceiling 19** — `TestCompatibilitySecurityChecklist` sees 015–019, forbids 020+.
3. **P19 keepers** — six named loop tests still PASS (no regression from S06).
4. **SelectNext table** — `TestSelectNext` covers blocking uncertainty, open regression, verification_incomplete, hop budget, p19_saturated; never EXECUTE on blocking uncertainty.
5. **Gates** — `TestTestPassAloneCannotSatisfyVerificationGate` cited; CLI shows verification debt.
6. **Contradiction path** — S03 contradiction + S05 correlated regression; never auto-`caused`.
7. **CLI evidence** — archived under `experiments/runs/…-p20-s07-01-verify/`; ordinary CLI not test-only helpers.
8. **Must checklist** — S07-01 Notes map COVERAGE Must rows to evidence; reviewer spot-checks ≥10 anchors.
9. **§31 mini-eval** — fixture-scale acceptable with explicit residual if used.
10. **Seed export** — P17 keepers PASS; P20 tables omitted by design (documented residual, **not** fail).
11. **§29O** — hallucinated/incomplete evidence blocked by gates (reviewer cites 2+ tests).
12. **No scope creep** — no mig 020, no new MCP, no raw CoT storage.

## Locked re-verify commands

Re-run the **same floor** as S07-01 (minimum):

```bash
go test ./internal/deliberation/... -count=1 -run 'TestSelectNext|TestSelectNextNeverExecuteOnBlockingUncertainty|TestApplyTransitionHopBudgetDoesNotIncrementPastN'
go test ./internal/domain/ -count=1 -run 'TestTestPassAloneCannotSatisfyVerificationGate|TestCorrelationAndContradictionNeverAutoSetCaused|TestCountBlockingUncertaintiesFeedsApplyDeliberationTransition|TestContradictedEffectDoesNotCreateRegressionOrAutoReplan'
go test ./internal/loop/... -count=1 -run 'TestLoopApplyUncertaintyWriteAffectsNextSelectNext|TestLoopNextVerifySurfacesVerificationDebt|TestLoopApplyUnknownWriteKeyFailsClosed'
go test ./cmd/trace -count=1 -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestHelpIncludesLoopNext'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit'
```

Broader re-run encouraged; **minimum** `-run` above is mandatory for APPROVE.

## Evidence artifacts

- Read S07-01 `experiments/runs/…/evidence/*` + optional `VERIFY-NOTES.md`
- Write **`REVIEW-NOTES.md`** in this scope folder (recommended)
- Update board Notes: confidence, successor, residuals

## Spawn policy

| Severity | Action |
|----------|--------|
| blocker / high | Small inline fix **or** spawn `P20-S07-02a` + `02b` immediately below |
| medium | Prefer spawn unless trivial one-liner |
| low / nit | List in REVIEW-NOTES; do not block close |

## Exit criteria

- [ ] Independent re-verify floor PASS
- [ ] S07-01 evidence reviewed (CLI files exist, schemas correct)
- [ ] No open blocker/high without pending spawn
- [ ] **`DR-HANDOFF.md` CLOSED** with successor **`no successor`**
- [ ] Phase 20 marked complete in board Notes
- [ ] Confidence **high** (or **medium** with explicit residuals — never silent)
- [ ] Update [`docs/TODO.md`](../../../../TODO.md) index if phase complete (status line only)

## Forbidden

- Leaving successor **TBD** when row is `done`
- Closing DR-HANDOFF on S07-01 evidence without independent re-run
- Rewriting Phase 19 `done` history or DR-HANDOFF
- Implementing seed export extension or Phase 21 scaffold as part of review
