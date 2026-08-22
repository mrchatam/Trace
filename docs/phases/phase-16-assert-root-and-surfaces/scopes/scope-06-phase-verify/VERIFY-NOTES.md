# P16-S06-01 — Phase VERIFY notes (assert-root-and-surfaces closeout)

**Date:** 2026-08-17  
**Verifier:** independent re-run (does **not** trust S01–S05 Notes alone)  
**Verdict:** **Phase 16 VERIFY PASS / assert-root-and-surfaces green**  
**Confidence:** high  
**Spawns:** none  

**Explicit claims:** S01–S05 named DF regressions green on live packages. Catalog **10** including `trace_version` (`trace_impact` before version; slug `mcp:trace_impact`). DF-72 thin `trace_impact` is a **fail bar** — `TestMCPTraceImpactReport` + `TestMCPImpactDeniedBlocksCallTool` PASS. Carry-forward honesty Paths A/B/C + Gate G + Gate E + Gate F + capability ablation + Gate H + compat (mig ceiling **14**) + p0x 7/7 + x0 + product `./cmd|internal|evals` (CGO1) PASS. Gate C artifacts remain **Go** (`dry_run:false`, N=3; mean G1 **0.800** > B0 **0.000**).

**Explicit non-claims:** DF-67 / P14 R2 / P15 R3/R4 **not** fixed. S05-02 `attachTaskImpact` swallow **not** fixed. 014 nine-Name `IN` list **not** edited. Phase 01 dry-run is **not** Gate C, **not** Gate F, **not** Gate G, **not** ablation, **not** Gate H, and **not** the compat checklist. VerifiedFact still out. No install/decide/plan/index MCP. No product Go on this row. No research S05 / `plan simulate` / D21+ scaffold. Phase 16 not marked complete here — **P16-S06-02** owns handoff close + phase complete. Phase 15 historical `no successor` left intact as history. **Phase 17 is independently queued and is not this successor.**

**DR-HANDOFF = `no successor`.** Phase 17 remains independently queued after the P16 board block — **not** this VERIFY successor; P17 rows were **not** rewritten. Research S05 / `plan simulate` / D21+ stay off-board unless Notes explicitly promote (they do not). Parallel dogfood may continue under `experiments/` off-board.

## Environment

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` (`go.mod`) |
| `go version` | go1.24.2 linux/amd64 |
| Named S01–S05 library / honesty / E / F / ablation | `CGO_ENABLED=0` + `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off` on mcp/library lines |
| CLI named / Gate H / compat / p0x / x0 / product bar | `CGO_ENABLED=1` |
| Product bar env | `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off` |
| Gate C metrics | `docs/verification/gate-c-x0/` (`dry_run:false`, N=3/condition) |
| Gate C means (inspect only) | B0 mean **0.000**; G1 mean **0.800** — **not** re-scored |
| Optional dogfood | `experiments/` **not** run this row (non-blocking) |

## Evidence table (independent)

| Bucket / command | Result |
|------------------|--------|
| S01 DF-76 + P15 keepers `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered\|TestMCPAssert\|TestMCPVirgin\|TestMCPInitialized\|TestOpenExisting\|TestOpenCreates\|TestBuiltinMCPCapabilitySpecs\|TestCapabilityDecision'` | **PASS** — mcp ok; domain ok; store ok |
| S02 DF-75/78 `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestCapabilityToolDecisionCheckRejectsYOLO\|TestCapabilityToolDecisionMigrateHealsYOLOToPending\|TestResolveYOLOBuiltinDoesNotAutoAllow\|TestDecideUnprefixedMCPNameCanonicalizes\|TestCanonicalizeCustomAndCLISlugsUnchanged\|TestMigrateUnprefixedDeniedFoldsOverAutoAllowed\|TestMCPUnprefixedDecideGatesCallTool\|TestMCPAssert\|TestToolNamesRegistered\|TestOpenCreates\|TestMCPVirgin'` | **PASS** — mcp ok; domain ok; store ok |
| S03 DF-77 library `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestCLIAddDeniedDoesNotBlockMCPAdd\|TestUnprefixedAddDecideDoesNotGateCLI\|TestCapabilityDecisionAutoAllowBuiltinCLI\|TestCanonicalizeCLIReindexFoldsToIndex\|TestCanonicalizeCustomAndCLISlugsUnchanged\|TestCapabilityDecision\|TestMCPAssert\|TestMCPUnprefixed\|TestToolNamesRegistered\|TestMCPVirgin\|TestOpenCreates'` | **PASS** — mcp ok; domain ok; store ok |
| S03 DF-77 CLI `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestCLIAddSucceedsWhenMCPAddDenied\|TestCLIAddDeniedFailClosed\|TestCLIWhySucceedsWhenMCPWhyDenied\|TestCLIWhyDeniedFailClosed\|TestUngatedCapabilityDecideWhenCLIAddDenied\|TestCLIIndexAliasDenied\|TestUnprefixedAddDecideDoesNotGateCLI'` | **PASS** — cmd/trace ok |
| S04 DF-68 library `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/install/... ./internal/mcp/... ./internal/store/... -count=1 -run 'TestInstallDetectListsCursorStable\|TestInstallConditional\|TestToolNamesRegistered\|TestMCPVirgin\|TestOpenCreates'` | **PASS** — install ok; mcp ok; store ok |
| S04 DF-68 CLI `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallClaudeDashC\|TestInstallDetectDashC\|TestInstallUninstallClaudeDashC\|TestInstallCursorPrintReloadTip\|TestInstallCursorWriteMergeBackup\|TestCLIAddDeniedFailClosed'` | **PASS** — cmd/trace ok |
| S05 DF-70…74 + catalog 10 library `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/compiler/... ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestContextIncludesImpactOverallClass\|TestMCPTraceImpactReport\|TestMCPImpactDeniedBlocksCallTool\|TestToolNamesRegistered\|TestBuiltinMCPCapabilitySpecs\|TestImportBoundaryMCPNoPlanImpactIndexTools\|TestMCPVirgin\|TestOpenCreates'` | **PASS** — compiler ok; mcp ok; domain ok; store ok (DF-72 named tests are a fail bar) |
| S05 DF-70…74 CLI `CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedImportDiscoveryMentionsTask\|TestSeedImportImpactFindings\|TestImpactReportJSONSnakeCase\|TestWhyIncludesImpactOverallClass\|TestCLIAddDeniedFailClosed\|TestLinkDiscoveryMentionsTaskCLI\|TestInstallClaudeDashCRefuseCitesProjectRoot'` | **PASS** — cmd/trace ok |
| Honesty full `CGO_ENABLED=0 go test ./evals/honesty/... -count=1` | **PASS** |
| Honesty A/B/C + G `CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim\|TestHonestyEscapeRateGateGPrelim'` | **PASS** |
| Gate E `CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan` | **PASS** |
| Gate F `CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim` | **PASS** |
| Ablation `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` | **PASS** |
| Gate H `CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH` | **PASS** (~5.2s) |
| Compat `CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist` | **PASS** (mig ceiling **14**; no 015+) |
| P0-X + X0 `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1` | **PASS** — p0x + x0 |
| Product `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1` | **PASS** — all product pkgs ok (incl. analyzers under CGO1) |
| Catalog 10 grep `RegisteredToolNames` / `registerTools` | **PASS** — exactly 10: why, context, add, link, transition, review, tasks, capability, **impact**, **version**; `trace_impact` before `trace_version`; no `trace_plan` / `trace_index` / `trace_install` / `trace_decide` MCP |
| Gate C artifacts inspect | **PASS** — `dry_run:false`; N=3 (b0-gatec-1..3, g1-gatec-1..3); mean G1 0.800 > B0 0.000; **not** re-scored |
| Mig 014 / no VERIFY mig | **PASS** — no new migration from this row; ceiling **14** |
| No committed `.trace/` under `fixtures/` / `evals/` | **PASS** |
| G19 library packages do not import `cmd/trace` / `cmd/trace-mcp` | **PASS** — only test-string mentions in `internal/mcp/mcp_test.go` and `evals/compat/compat_test.go` |
| Phase 17 rows intact | **PASS** — P17 board rows 232–244 not rewritten; not claimed as this successor |

## Law checks

| Check | Hold? |
|-------|-------|
| No daemon / always-on HTTP as primary surface | Yes |
| No committed `.trace/` under `fixtures/` or `evals/` | Yes |
| G19 — library packages do not import `cmd/trace` or `cmd/trace-mcp` | Yes |
| S01–S05 evidence is **named tests** — not Notes-only | Yes |
| MCP remains **ten** tools including `trace_version`; `trace_impact` **in** catalog; no install/decide/plan/index MCP | Yes |
| Honesty A/B/C + Gate G + Gate E + Gate F + ablation + Gate H + compat **14** still green | Yes |
| Gate C evidence remains `dry_run:false` — not Phase 01 dry-run | Yes |
| Embeddings / VerifiedFact / Neo4j SoT still out | Yes |
| No full-rebuild-on-any-change indexer architecture | Yes |
| No new migration from VERIFY; mig 014 already in tree; compat ceiling **14** | Yes |
| No YOLO / AllowAll defaults | Yes |
| **DF-67 / R2 / R3 / R4 / attachTaskImpact swallow / 014 nine-Name do not fail VERIFY** | Yes |
| **No research S05 / plan simulate / D21+ scaffold** unless Notes explicitly promote | Yes — not promoted |
| **Phase 17 board rows left intact** — not claimed as this successor | Yes |
| Forward-only: do **not** rewrite Phase 00–15 `done` history; Phase 15 historical `no successor` left intact as history | Yes |

## Residuals / deferrals (non-blocking)

| Residual | Disposition | VERIFY note |
|----------|-------------|-------------|
| **DF-67** symbol-entity honesty | **defer** | Still out of `index_honesty` bar (compiler honesty has no symbol-entity coverage) — **not** fail criteria; **not** claimed fixed |
| **P14 R2** `allowContainsOut` late-upgrade | **defer** | Spot-checked still present in `internal/retrieval/impact_walk.go` — **not** fail criteria; **not** claimed fixed |
| **P15 R3** graphify space-in-path on full `./...` | **wontfix** | Product bar is `./cmd\|internal\|evals` (PASS) — do **not** fail VERIFY |
| **P15 R4** CGO0 analyzers / CGO0 `cmd/trace` FAIL | **wontfix** | Product bar is CGO1 (analyzers PASS under product suite); CGO0 `cmd/trace` was **not** used as a fail bar |
| **S05-02** `attachTaskImpact` swallow | residual | `internal/compiler/compiler.go` still returns on helper error (capability-style) — **not** fail criteria |
| **014 nine-Name `IN` list** | historical | `014_capability_tool_decision_enum.sql` still lists nine unprefixed Names (no `trace_impact`) — **not** edited; **not** fail criteria |
| **DF-22/37** Cursor reload | ops | tip keepers only (`TestInstallCursorPrintReloadTip` / `TestInstallCursorWriteMergeBackup`) — **not** fail criteria |
| Goals #2–#4 | deferred | Research S05 / `plan simulate` / D21+ stay **off-board** |
| Parallel dogfood | off-board | `experiments/` may continue; **not** board-blocking; **not** run this row |
| VerifiedFact / embeddings / daemon-HTTP primary | out | unchanged |

## Dry-run ≠ gates

**Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H / ≠ checklist.** Gate C artifacts remain Mode-B `dry_run:false` (inspect only; not re-scored). Phase 01 dry-run is regression-only.

## Catalog 10 + DF-72 (fail bar)

Live `RegisteredToolNames()` / `TestToolNamesRegistered`: **10** names — `trace_why`, `trace_context`, `trace_add`, `trace_link`, `trace_transition`, `trace_review`, `trace_tasks`, `trace_capability`, `trace_impact`, `trace_version`. `trace_impact` is registered **before** `trace_version`. Boundary `TestImportBoundaryMCPNoPlanImpactIndexTools` allows `trace_impact` only (forbids plan/index/install/decide MCP). DF-72 named tests `TestMCPTraceImpactReport` and `TestMCPImpactDeniedBlocksCallTool` **PASS** — missing these would **fail** VERIFY (historical P16-00 defer is not the live lock).

## Handoff

| Item | Value |
|------|-------|
| **DR-HANDOFF** | **`no successor`** (**started** P16-S06-01; close owned by **P16-S06-02**) |
| Phase 17 | Independently queued (board rows 232–244) — **not** this successor; **not** rewritten |
| Research S05 / plan simulate / D21+ | **Do not scaffold** — no promotion |
| Parallel dogfood / research FUTURE | May continue off-board under `experiments/` / research docs |
| Completion owner | **P16-S06-02** — refuse `done` until VERIFY-NOTES + fresh evidence agree **and** handoff is explicitly `no successor` |
| Spawns | **none** |
| Next board | **P16-S06-02** |
