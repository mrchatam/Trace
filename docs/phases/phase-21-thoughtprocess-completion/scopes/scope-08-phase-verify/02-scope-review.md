# P21-S08-02 — Review: phase verify + close DR-HANDOFF

## Metadata
- id: P21-S08-02
- todo_ids: [P21-S08-02]
- role: reviewer
- skills: [code-review-and-quality, documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective

Independent **fresh-session** review of S08-01 evidence; re-run locked verify floor (do **not** trust Notes alone); **close DR-HANDOFF** with explicit successor decision (**never TBD**).

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff (mandatory)
- [00-PLANNER.md](00-PLANNER.md) — FINAL evidence bar + DR-HANDOFF policy
- [01-verify.md](01-verify.md) — locked floor Blocks A–C + §31 live mini-eval
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [WORK-MAP.md](../../WORK-MAP.md) W-14
- [DECISION-LOG.md](../../DECISION-LOG.md) — D-01…D-15 closures
- [COVERAGE.md](../../../phase-20-cognitive-deliberation/COVERAGE.md)
- Pattern: [P20 S07-02 review](../../../phase-20-cognitive-deliberation/scopes/scope-07-phase-verify/02-scope-review.md)

## Session start

Follow agent-loop-protocol. **Fresh session** — reviewer must not be the S08-01 implementer. Unattended: execute review loop until blocker/high clear or spawned forward.

## Locked DR-HANDOFF close policy (FINAL — S08-00)

| Field | Locked value |
|-------|--------------|
| Who closes | **S08-02 only** (S08-01 gathers evidence; S08-00 locked policy) |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Successor decision | **`no successor`** — default unless **human operator** has already named a Phase 22 before this row runs |
| Rationale | Phase 21 completes P20 residuals + TRACE_THOUGHTPROCESS promoted gaps (WORK-MAP W-01…W-15). Hosted MCP, daemon, HTTP remain **Later developments** (D-19). Human may promote Phase 22 later (as Phase 21 followed Phase 20) — that is **not** this close's job unless pre-named. |
| Must not | Leave `Successor decision: TBD`, `later`, or empty; rewrite Phase 20 historical `no successor`; auto-scaffold Phase 22 without human promotion |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF CLOSED + all Phase 21 board rows done |
| Next runnable after close | **none** (board queue empty until human promotes) |

### Residuals to list on close (non-blocking)

| Topic | Disposition |
|-------|-------------|
| Experiments not in seed export | Operational records; S08 verify uses domain/CLI |
| Trace does not run tests autonomously | D-16 **out** — harness-agnostic by design |
| Hosted MCP / daemon / HTTP | D-19 **out** — separate repo / Later developments |
| ML / learned policies | D-21 **out** |
| Graph DB / embeddings | D-22 **out** |
| Requirement as own table | D-17 **keep** merged into Goal |
| Full §16 bake-off engine | D-01 thin promote only — not multi-agent runner |
| Full §18 ML risk matrix | D-02 minimal hints only |

### DR-HANDOFF.md update template (on APPROVE)

```markdown
**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Closed | YYYY-MM-DD |
| Successor decision | **no successor** |
| Phase 21 outcome | P20 residuals closed (seed, retrieval, cycle, promotion, why, tx apply); thin §16/§18 promoted; schema through **021** |
| Residuals (non-blocking) | experiments not in seed; autonomous test execution out-of-scope; hosted MCP separate track |
| Forward (human queue) | Phase 22 only if human promotes; MCP/daemon if human approves Later developments |
```

If human has **named Phase 22** before S08-02: replace successor line with `Phase 22 — <human-given name>` and cite promotion evidence in board Notes.

If verify bar **fails**: keep DR-HANDOFF **OPEN**, spawn `P21-S08-02a` implement + `02b` review immediately below this row; successor stays **`no successor` intent** until pass.

## Review focus

Confirm independently:

1. **P20 floor (Block A)** — all command blocks from [01-verify.md](01-verify.md) re-run PASS (session ≠ S08-01).
2. **P21 deltas (Block B)** — S01–S07 named keepers PASS; no regression on P20 keepers.
3. **Compat ceiling 21** — `TestCompatibilitySecurityChecklist` sees 015–021, forbids 022+; **21** embed files.
4. **Schema** — max **021** (`020_baselines_promotion`, `021_experiments`); no 022+.
5. **Seed export** — P20 keys round-trip (D-05 closed); P17 keepers PASS; denied surfaces still omitted.
6. **SelectNext 14-row cycle** — `TestSelectNextFullCycleOrdering`; EXECUTE/TEST/VERIFY/EVALUATE/REFLECT/REPLAN reachable (D-03).
7. **Baseline promotion** — `TestPromoteBaselineSupersedesPrior`; eval regression → `promotion_blocked` advisory (D-09, D-10).
8. **Retrieval + FTS** — no `retrieval: unknown entity type` on INVESTIGATE stderr (D-06, D-07).
9. **Observability** — `trace why` P20 types; `historical_relationships` in loop next (D-11, D-12).
10. **Transactional apply** — `TestLoopApplyTransactionalRollbackOnFailure`; goal_id guards (D-08, D-13, D-15).
11. **Experiments + risk hints** — thin table only; no runner; `risk_hints` cap 4 (D-01, D-02).
12. **CLI evidence** — archived under `experiments/runs/…-p21-s08-01-verify/`; ordinary CLI not test-only helpers.
13. **§31 live mini-eval** — **≥5** distinct phases in CLI JSON; **live repo** primary (D-14 promoted; fixture-only **not** sufficient alone).
14. **Must checklist** — S08-01 Notes map COVERAGE Must rows; reviewer spot-checks ≥10 anchors including §31.
15. **§29O** — hallucinated/incomplete evidence blocked (cite 2+ tests).
16. **No scope creep** — no mig 022+, no new MCP tools (still **10**), no raw CoT storage.
17. **P20 DR-HANDOFF** — historical `no successor` **unchanged** (D-24).

## Review checklist

| Check | Evidence |
|-------|----------|
| Block A P20 floor | All P20 S07-01 named tests PASS (reviewer re-run) |
| Block B S01 seed | `TestSeedExportIncludesP20Cognition`, `TestSeedImportP20RoundTrip` |
| Block B S02 retrieval | `TestExactLookupUncertainty`, `TestLoopNextInvestigateNoRetrievalStderr`, FTS sync tests |
| Block B S03 cycle | `TestSelectNextFullCycleOrdering`, EXPLORE/EXECUTE row tests |
| Block B S04 promotion | `TestPromoteBaselineSupersedesPrior`, `TestEvalRegressionBlocksPromotionGate` |
| Block B S05 why/historical | `TestCLIWhyUncertainty`, `TestLoopNextHistoricalRelationshipsSection` |
| Block B S06 tx apply | `TestLoopApplyTransactionalRollbackOnFailure`, goal_id tests; internal/loop ≥8 apply tests |
| Block B S07 experiments/hints | 6 S07 named tests; no `os/exec` runner |
| Compat ceiling **21** | `CGO_ENABLED=1` compat + `TestMigrationStatusReportsEmbedMax` |
| Loop schemas | `trace.loop.next.v1`, `trace.loop.apply.v1`, `trace.loop.status.v1` unchanged strings |
| MCP catalog | **10** tools — no new registrations |
| CLI evidence files | #1–#7 + `99-run-metadata.txt` exist under evidence dir |
| §31 live eval | CLI #2 shows **≥5** phases; live `trace` binary used |
| Seed pre-PR | `trace/graph.json` export attempted; P20 keys when populated |
| D-05…D-15 | Each promoted item has test or CLI anchor cited in Notes |

## Locked re-verify commands

Re-run the **same floor** as S08-01 (minimum — reviewer session):

```bash
# P20 core (abbreviated minimum)
go test ./internal/deliberation/... -count=1 -run 'TestSelectNext|TestSelectNextNeverExecuteOnBlockingUncertainty|TestApplyTransitionHopBudgetDoesNotIncrementPastN'
go test ./internal/domain/ -count=1 -run 'TestTestPassAloneCannotSatisfyVerificationGate|TestCorrelationAndContradictionNeverAutoSetCaused|TestCountBlockingUncertaintiesFeedsApplyDeliberationTransition|TestContradictedEffectDoesNotCreateRegressionOrAutoReplan'
go test ./internal/loop/... -count=1 -run 'TestLoopApplyUncertaintyWriteAffectsNextSelectNext|TestLoopNextVerifySurfacesVerificationDebt|TestLoopApplyTransactionalRollbackOnFailure|TestLoopApplyDeliberationTransitionEvent'
go test ./cmd/trace -count=1 -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestHelpIncludesLoopNext'

# P21 deltas (mandatory)
go test ./internal/domain/... -count=1 -run 'TestSeedExportIncludesP20Cognition|TestSeedImportP20RoundTrip|TestPromoteBaselineSupersedesPrior|TestEvalRegressionBlocksPromotionGate|TestCreateExperimentLinksOutcome'
go test ./internal/deliberation/... -count=1 -run 'TestSelectNextFullCycleOrdering|TestSelectNextExecuteWhenPending'
go test ./internal/retrieval/... -count=1 -run 'TestExactLookupUncertainty|TestWhyUncertaintyIncludesGraphSteps'
CGO_ENABLED=1 go test ./cmd/trace -count=1 -run 'TestLoopNextInvestigateNoRetrievalStderr|TestCLIWhyUncertainty|TestLoopNextHistoricalRelationshipsSection'
go test ./internal/loop/... -count=1 -run 'TestRiskHintsManyPaths|TestLoopNextRiskHintsBounded|TestLoopApplyGoalIDMismatchFailsClosed'

# Compat + seed keepers
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestSeedExportWritesExportedAtCommit'
go test ./internal/store/... -count=1 -run TestMigrationStatusReportsEmbedMax
```

Broader re-run of Blocks A–C encouraged; **minimum** `-run` above is mandatory for APPROVE.

## Evidence artifacts

- Read S08-01 `experiments/runs/…-p21-s08-01-verify/evidence/*` + optional `VERIFY-NOTES.md`
- Write **`REVIEW-NOTES.md`** in this scope folder (recommended)
- Update board Notes: confidence, successor, residuals, schema count (**21**)

## Spawn policy

| Severity | Action |
|----------|--------|
| blocker / high | Small inline fix **or** spawn `P21-S08-02a` + `02b` immediately below |
| medium | Prefer spawn unless trivial one-liner |
| low / nit | List in REVIEW-NOTES; do not block close |

## Exit criteria

- [ ] Independent re-verify floor PASS (Blocks A minimum + P21 delta minimum above)
- [ ] S08-01 evidence reviewed (CLI files exist, schemas correct, ≥5 phases in #2)
- [ ] No open blocker/high without pending spawn
- [ ] **`DR-HANDOFF.md` CLOSED** with successor **`no successor`** (or human-named Phase 22)
- [ ] Phase 21 all rows `done` in board Notes
- [ ] Confidence **high** (or **medium** with explicit residuals — never silent)
- [ ] Update [`docs/TODO.md`](../../../../TODO.md) index if phase complete (status line only)

## Forbidden

- Leaving successor **TBD** when row is `done`
- Closing DR-HANDOFF on S08-01 evidence without independent re-run
- Rewriting Phase 20 `done` history or Phase 20 DR-HANDOFF
- Accepting fixture-only §31 mini-eval as sole evidence (D-14 requires live repo)
- Implementing mig 022+, seed experiments export, or Phase 22 scaffold as part of review

## Next

**none** after close (human promotes forward work)
