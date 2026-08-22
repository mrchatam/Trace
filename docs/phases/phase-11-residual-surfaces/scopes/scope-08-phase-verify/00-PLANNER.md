# P11-S08-00 — Phase VERIFY / residual closeout (FINAL)

## Metadata
- id: P11-S08-00
- todo_ids: [P11-S08-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock Phase 11 VERIFY evidence table: **S01–S07 named DF regressions** (all **18** canonical P11 DFs) + **carry-forward gates** + product `./...`. Decide **DR-HANDOFF** = **`no successor`** unless Notes promote Phase 12. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — all P11 DF regressions + carry-forward
- Pattern: Phase 10 S05 [`../../../phase-10-integrity-surfaces/scopes/scope-05-phase-verify/`](../../../phase-10-integrity-surfaces/scopes/scope-05-phase-verify/)
- Sibling REVIEW-NOTES: S01–S07 under `../scope-0{1..7}-*/REVIEW-NOTES.md`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — grill only if A1–A7 conflict.

## Depends-on (S01–S07 — landed)

| Scope | Board | DF locks imported |
|-------|-------|-------------------|
| S01 | P11-S01-02 APPROVE high | DF-40 `TestIndexPartialArgvGCAfterRename` (+ P10 DF-20 retain) |
| S02 | P11-S02-02 APPROVE high | DF-43 sibling FAIL+PASS; DF-44 conscious flag≠identity |
| S03 | P11-S03-02 APPROVE high | DF-47 bounded retry + exclusivity + serialize UX |
| S04 | P11-S04-02 APPROVE high | DF-41 slug upsert; DF-51 hatch≠missing-caps |
| S05 | P11-S05-02 APPROVE high | DF-49/35/48/42 symbol / depth / trust / discovery→task |
| S06 | P11-S06-02 APPROVE high | DF-50/22/37 install tip parity + nine tools/`trace_version` |
| S07 | P11-S07-02 APPROVE high | DF-33/30/46/45/28 seed aliases; plan show; review CLI; thin handoff SoT |

## Live residuals → DR-HANDOFF decision (2026-08-16)

| Bucket | Items | Phase implication |
|--------|-------|-------------------|
| Product gaps scheduled in Phase 11 | 18 DFs → S01–S07 | Closed by S01–S07 APPROVE high — VERIFY must **re-prove named tests** |
| Explicit residual OK into VERIFY | rename+edit ghost until full-tree; empty-result review path; soft OR in older budget test; CGO0 analyzers path | Forward notes only — **not** a successor phase |
| Parallel dogfood (not board-blocking) | ab-* ladders / Cursor MCP reload manual | Stay in `experiments/` — **not** Phase 12 unless promoted |
| Known `./...` nit | `similar projects/graphify` space-in-path FAIL; optional CGO0 analyzers FAIL | Pre-existing non-product — VERIFY records **product pkgs PASS** |

**DR-HANDOFF = `no successor`.** Reopen only with explicit promotion + scaffold (same posture as Phase 10 historical close / Phase 11 forward reopen).

## Planner work
1. Lock VERIFY command set (S01–S07 named tests + carry-forward + product `./...`).
2. Thicken `01-verify.md` evidence table + spawn 01a/b/c + handoff **start**.
3. Thicken `02-scope-review.md` owns DR-HANDOFF **completion** (`no successor`).
4. SCOPE-TODOS + board sync; light stamp `DR-HANDOFF.md` ownership.

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Phase gate | Phase 11 residual-surfaces closeout — **not** a new planted eval gate |
| DF home | all **18** P11 DF regressions: DF-22, 28, 30, 33, 35, 37, 40–51 |
| Migration | **None** — no `011_*` from Phase 11 VERIFY |
| S01 DF-40 | `TestIndexPartialArgvGCAfterRename` (+ store `TestListFilePathsByContentHash` if present); retain `TestIndexGCAfterPathRename` / `TestIndexArgvMissingPathDeletesOnlyThatPath` / `TestIndexIncrementalIsolation` |
| S02 DF-43/44 | `TestSiblingFailBlocksDone` / `TestSiblingPassAloneAllowsDone` / `TestSiblingPassPlusUncertainAllowsDone` / `TestHatchBypassesSiblingFail` + `TestAsOperatorFlagIdentityDocs` / `TestAsOperatorSchemaIdentityDocs` / `TestOperatorDoneRequiresFlag`; honesty Path C supersession via A/B/C suite |
| S03 DF-47 | `TestOpenRetrySucceedsWhenLockReleasedSoon` / `TestConcurrentStoreOpenFailClosed` / `TestErrLockedSerializeGuidance` / `TestHelpSerializeLockGuidance` / `TestInitFailClosedWhenStoreLocked` + compat `trace_lock_ok` |
| S04 DF-41/51 | `TestUpsertCapabilityBySlugUpdatesExisting` / `TestUpsertCapabilityGetAndReject` / `TestAllowDoneDoesNotBypassMissingCaps` / `TestAllowDoneWarnsOnStderr` / `TestTransitionAllowDoneEmitsWarning` |
| S05 DF-49/35/48/42 | `TestWhySymbolExact` / `TestGetSymbolByID` / `TestExpandDepth2NoSiblingTaskBody` / `TestExpandContextDepth2NoSiblingTaskBody` / `TestDecisionMarkdownTrustLabels` / `TestLinkDiscoveryMentionsTask` / `TestLinkDiscoveryMentionsTaskCLI` (+ multi-goal DPC still green) |
| S06 DF-50/22/37 | `TestInstallCursorPrintReloadTip` / `TestInstallCursorWriteMergeBackup` / `TestInstallCursorWriteCreateMissing` + `TestToolNamesRegistered` / `TestTraceVersion` (nine tools retained; no kill/daemon) |
| S07 DF-33/30/46/45/28 | `TestSeedImportFromIDAliases` / `TestSeedImportMissingEndpointsMessage` / `TestPlanShowSnakeCaseAndEmptyPhases` / `TestPlanShowWithPhasesSnakeCase` / `TestReviewGetShowList` / `TestHelpHandoffSoT` |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x 7/7; x0; Gate C `dry_run:false` N=3; prior P10 DF-17…32 surfaces stay green |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist |
| Full bar | `CGO_ENABLED=1 go test ./... -count=1` — **product pkgs PASS**; known FAIL only `similar projects/graphify` space (non-product); CGO0 analyzers FAIL OK residual if present |
| Allowed Go on VERIFY | **None** for features — re-run + evidence docs only; spawn remediation if fail |
| Spawn | On fail: `01a` implement / `01b` review (+`01c` re-VERIFY if needed) immediately below |
| Findings closeout | On PASS: S08-01/02 may flip the 18 DF rows in `DOGFOOD-FINDINGS.md` to closed (forward-only; no ID reuse) |
| DR-HANDOFF | **`no successor`** — **S08-01 starts** Notes; **S08-02 owns completion**. Do **not** scaffold Phase 12 unless Notes explicitly promote. |
| Forbidden | Daemon/HTTP/embeddings; full-rebuild indexer; rewriting Phase 00–10 `done` history; claiming Phase 10 historical handoff was wrong |

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable (thickened)
- [x] VERIFY commands + DR-HANDOFF locked
- [x] SCOPE-TODOS + board Notes; next `P11-S08-01`
- [x] Product Go — **not** this row

## Out of scope
- Running VERIFY (S08-01)
- Product Go / new MCP tools / daemon / mig
- Scaffolding Phase 12 without explicit promotion
- Closing parallel dogfood experiments
- Claiming Phase 10 historical handoff was wrong
