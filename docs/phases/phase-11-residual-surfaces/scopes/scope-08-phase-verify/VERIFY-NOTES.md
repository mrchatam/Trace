# P11-S08-01 — Phase VERIFY notes (residual-surfaces closeout)

**Date:** 2026-08-16  
**Verifier:** independent re-run (does **not** trust S01–S07 Notes alone)  
**Verdict:** **Phase 11 VERIFY PASS / residual surfaces green**  
**Confidence:** high  
**Spawns:** none  

**Explicit claims:** All **18** canonical P11 DFs re-proved via live named tests — S01 DF-40 (+ DF-20 retain); S02 DF-43/44 + honesty Path C; S03 DF-47; S04 DF-41/51; S05 DF-49/35/48/42; S06 DF-50/22/37 (nine tools / `trace_version`; tip/docs only); S07 DF-33/30/46/45/28. Honesty Paths A/B/C + Gate G + Gate E + Gate F + capability ablation + Gate H + compat checklist + p0x + x0 + domain/store/planner/compiler/mcp/retrieval + product packages PASS. Gate C artifacts remain **Go** (`dry_run:false`, N=3; mean G1 **0.800** > B0 **0.000**).  

**Explicit non-claims:** Phase 01 dry-run is **not** Gate C, **not** Gate F, **not** Gate G, **not** ablation, **not** Gate H, and **not** the compat checklist. VerifiedFact still out. No plan/impact/index MCP dump. No product Go on this row. Phase 11 not marked complete here — **P11-S08-02** owns handoff close + phase complete.

**DR-HANDOFF = `no successor`.** Parallel dogfood may continue under `experiments/` (ab-* ladders / Cursor MCP reload manual, etc.). Do **not** scaffold Phase 12 unless Notes explicitly promote.

## Environment

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| `go version` | go1.24.2 linux/amd64 |
| S01 cmd / S02 CLI / S03 CLI / S04 CLI / S05 CLI / S06 CLI / S07 / Gate H / compat / p0x / x0 / full suite | `CGO_ENABLED=1` |
| S01 store / S02 domain+mcp / S03 store / S04 domain+mcp / S05 libs / S06 mcp / Honesty / Gate E / Gate F / ablation / support pkgs | `CGO_ENABLED=0` where locked |
| Gate C metrics | `docs/verification/gate-c-x0/` (`dry_run:false`, N=3/condition; git_sha `15fe50a1…`) |
| Gate C means (inspect only) | B0 mean **0.000**; G1 mean **0.800** — **not** re-scored |

## Evidence table (independent)

| Bucket / command | Result |
|------------------|--------|
| S01 `CGO_ENABLED=0 go test ./internal/store/... -count=1 -run 'TestListFilePathsByContentHash'` | **PASS** |
| S01 `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestIndexPartialArgvGCAfterRename\|TestIndexGCAfterPathRename\|TestIndexArgvMissingPathDeletesOnlyThatPath\|TestIndexIncrementalIsolation'` | **PASS** (DF-40 + DF-20 retain) |
| S02 `CGO_ENABLED=0 go test ./internal/domain/... ./internal/mcp/... -count=1 -run 'TestSiblingFailBlocksDone\|TestSiblingPassAloneAllowsDone\|TestSiblingPassPlusUncertainAllowsDone\|TestHatchBypassesSiblingFail\|TestOperatorDoneRequiresFlag\|TestAsOperatorSchemaIdentityDocs'` | **PASS** (DF-43/44) |
| S02 `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestAsOperatorFlagIdentityDocs'` | **PASS** |
| S03 `CGO_ENABLED=0 go test ./internal/store/... -count=1 -run 'TestOpenRetrySucceedsWhenLockReleasedSoon\|TestConcurrentStoreOpenFailClosed\|TestErrLockedSerializeGuidance'` | **PASS** (DF-47) |
| S03 `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestHelpSerializeLockGuidance\|TestInitFailClosedWhenStoreLocked'` | **PASS** |
| S04 `CGO_ENABLED=0 go test ./internal/domain/... ./internal/mcp/... -count=1 -run 'TestUpsertCapabilityBySlugUpdatesExisting\|TestUpsertCapabilityGetAndReject\|TestAllowDoneDoesNotBypassMissingCaps\|TestTransitionAllowDoneEmitsWarning'` | **PASS** (DF-41/51) |
| S04 `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestAllowDoneWarnsOnStderr'` | **PASS** |
| S05 `CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./internal/store/... ./internal/domain/... -count=1 -run 'TestWhySymbolExact\|TestGetSymbolByID\|TestExpandDepth2NoSiblingTaskBody\|TestExpandContextDepth2NoSiblingTaskBody\|TestDecisionMarkdownTrustLabels\|TestLinkDiscoveryMentionsTask\|TestWhyTaskDPCMultiGoalNoForeignPollution'` | **PASS** (DF-49/35/48/42) |
| S05 `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestLinkDiscoveryMentionsTaskCLI'` | **PASS** |
| S06 `CGO_ENABLED=0 go test ./internal/mcp/... -count=1 -run 'TestToolNamesRegistered\|TestTraceVersion'` | **PASS** (nine tools + version) |
| S06 `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallCursorPrintReloadTip\|TestInstallCursorWriteMergeBackup\|TestInstallCursorWriteCreateMissing'` | **PASS** (DF-50/22/37 tip parity) |
| S07 `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedImportFromIDAliases\|TestSeedImportMissingEndpointsMessage\|TestPlanShowSnakeCaseAndEmptyPhases\|TestPlanShowWithPhasesSnakeCase\|TestReviewGetShowList\|TestHelpHandoffSoT'` | **PASS** (DF-33/30/46/45/28) |
| `CGO_ENABLED=0 go test ./evals/honesty/... -count=1` | **PASS** |
| `CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim\|TestHonestyEscapeRateGateGPrelim'` | **PASS** (A/B/C + Gate G) |
| `CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan` | **PASS** (Gate E) |
| `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` | **PASS** (Gate F) |
| `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` | **PASS** |
| `CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH` | **PASS** (~7.2s named; ~11s in `./...`) |
| `CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist` | **PASS** |
| `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1` | **PASS** — p0x + x0 |
| `CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... ./internal/mcp/... ./internal/retrieval/... -count=1` | **PASS** — all six |
| `CGO_ENABLED=1 go test ./... -count=1` | **PASS** product pkgs; known FAIL only `similar projects/graphify` space-in-path (non-product residual) |
| MCP nine tools / no plan/impact/index | **PASS** — `TestToolNamesRegistered` wants exactly nine; boundary rejects `trace_plan`/`trace_impact`/`trace_index` |
| Actor≠auth spot-check | **PASS** — `TestOperatorDoneRequiresFlag` rejects `Actor:"operator"` without `AllowOperatorDone` |
| Gate C artifacts inspect | **PASS** — `dry_run:false`; N=3; mean G1 0.800 > B0 0.000; **not** re-scored |
| No `011_*` from Phase 11 | **PASS** — no `011_*` migration files |
| No committed `.trace/` under `fixtures/` / `evals/` | **PASS** |
| G19 library packages do not import `cmd/trace` | **PASS** (`go list` clean) |

## Law checks

| Check | Hold? |
|-------|-------|
| No daemon / always-on HTTP as primary surface | Yes |
| No committed `.trace/` under `fixtures/` or `evals/` | Yes |
| G19 — library packages do not import `cmd/trace` or `cmd/trace-mcp` | Yes |
| S01–S07 evidence is **named tests** — not Notes-only | Yes |
| MCP is **nine** tools; no plan/impact/index MCP | Yes |
| Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat green | Yes |
| Gate C evidence remains `dry_run:false` — not Phase 01 dry-run | Yes |
| Embeddings / VerifiedFact / `plan simulate` still out | Yes |
| No full-rebuild-on-any-change indexer architecture (DF-40 is hash-orphan delete, not full rebuild) | Yes |
| No new migration `011_*` from Phase 11 | Yes |
| **No Phase 12 scaffold** | Yes |
| Forward-only: do **not** rewrite Phase 00–10 `done` history; Phase 10 historical `no successor` left intact | Yes |

## Residuals / deferrals

- **Known `./...` nit:** `similar projects/graphify` space-in-path setup FAIL — pre-existing non-product; product pkgs PASS.
- **CGO0 analyzers FAIL OK** residual if present on zero-CGO analyzer path — product bar uses CGO1 for analyzers (PASS in `./...`).
- **S01 residual:** rename+edit (hash diverges) ghost until full-tree — forward note only.
- **S06 residual:** live Cursor MCP reload remains **manual** (tip/docs only; no PID kill) — parallel `experiments/`.
- Parallel dogfood: ab-* ladders — **not** board-blocking; stay on `experiments/`.
- VerifiedFact / embeddings / `plan simulate` / daemon-HTTP primary still out.

## Findings closeout (forward-only)

Flipped DF-22, 28, 30, 33, 35, 37, 40–51 to **closed** in [`experiments/DOGFOOD-FINDINGS.md`](../../../../../experiments/DOGFOOD-FINDINGS.md) citing this VERIFY re-prove. Live Cursor reload remains an ops residual note under DF-22/37 closed rows (tip/docs mitigation; not a Phase 12).

## Handoff

| Item | Value |
|------|-------|
| **DR-HANDOFF** | **`no successor`** (**started** this row; **closed** on P11-S08-02 — see [REVIEW-NOTES.md](REVIEW-NOTES.md)) |
| Phase 12 | **Do not scaffold** — intentional absence (no promotion) |
| Parallel dogfood | May continue under `experiments/` off-board |
| Completion owner | **P11-S08-02** — ✅ closed |
| Next board row | **none** (roadmap closed again) |
