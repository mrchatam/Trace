# P16-S06-00 — Phase VERIFY / assert-root-and-surfaces (FINAL)

## Metadata
- id: P16-S06-00
- todo_ids: [P16-S06-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock Phase 16 VERIFY evidence: **S01–S05 named DF regressions** + **carry-forward gates** + product pkgs. Decide **DR-HANDOFF** = **`no successor`** (default). **No product Go.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — historical; DF-72 live lock is [../../DF-72-FORWARD.md](../../DF-72-FORWARD.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: Phase 15 S02 VERIFY [`../../../phase-15-p14-residual-plan/scopes/scope-02-phase-verify/00-PLANNER.md`](../../../phase-15-p14-residual-plan/scopes/scope-02-phase-verify/00-PLANNER.md)
- S01–S05 REVIEW-NOTES (all **APPROVE high**): [S01](../scope-01-mcp-project-root/REVIEW-NOTES.md), [S02](../scope-02-tool-decision-enum/REVIEW-NOTES.md), [S03](../scope-03-cli-mcp-allowlist-parity/REVIEW-NOTES.md), [S04](../scope-04-install-project-root/REVIEW-NOTES.md), [S05](../scope-05-seed-impact-packet/REVIEW-NOTES.md)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Depends-on: S01–S05 APPROVE (landed). Default handoff remains **`no successor`**. Unattended: no Plan-mode switch; no product Go.

## Depends-on (S01–S05 — landed; live named tests confirmed 2026-08-17)

| Scope | Board | Named tests imported (live `func Test*` exists) |
|-------|-------|--------------------------------------------------|
| S01 DF-76 | **APPROVE high** (P16-S01-02) | `TestMCPVirginProjectDoesNotMkdir`; `TestMCPInitializedOtherRootIsolated`; `TestOpenExistingMissingReturnsErrNotInitialized`; `TestOpenExistingEmptyTraceDir`; P15 keepers `TestMCPAssertDeniedBlocksCallTool` / `TestMCPAssertBuiltinAutoAllowedSucceeds` / `TestToolNamesRegistered` (**live length 10** after S05) |
| S02 DF-75/78 | **APPROVE high** (P16-S02-02) | `TestCapabilityToolDecisionCheckRejectsYOLO`; `TestCapabilityToolDecisionMigrateHealsYOLOToPending`; `TestResolveYOLOBuiltinDoesNotAutoAllow`; `TestDecideUnprefixedMCPNameCanonicalizes`; `TestCanonicalizeCustomAndCLISlugsUnchanged`; `TestMigrateUnprefixedDeniedFoldsOverAutoAllowed`; `TestMCPUnprefixedDecideGatesCallTool`; keepers `TestMCPAssert*` / `TestToolNamesRegistered` / `TestOpenCreates*` (v**14**) / `TestMCPVirgin*` |
| S03 DF-77 | **APPROVE high** (P16-S03-02) | CLI CGO1: `TestCLIAddSucceedsWhenMCPAddDenied`; `TestCLIAddDeniedFailClosed`; `TestCLIWhySucceedsWhenMCPWhyDenied`; `TestCLIWhyDeniedFailClosed`; `TestUngatedCapabilityDecideWhenCLIAddDenied`; `TestCLIIndexAliasDenied`; `TestUnprefixedAddDecideDoesNotGateCLI`. Library: `TestCLIAddDeniedDoesNotBlockMCPAdd`; `TestUnprefixedAddDecideDoesNotGateCLI`; `TestCapabilityDecisionAutoAllowBuiltinCLI`; `TestCanonicalizeCLIReindexFoldsToIndex`; `TestCanonicalizeCustomAndCLISlugsUnchanged` |
| S04 DF-68 | **APPROVE high** (P16-S04-02) | CLI CGO1: `TestInstallClaudeDashCRefuseCitesProjectRoot`; `TestInstallClaudeDashCIgnoresCwdMarker`; `TestInstallClaudeDashCWriteUsesProjectRoot`; `TestInstallDetectDashCClaudeReasonCitesRoot`; `TestInstallDetectDashCCursorHomeUnchanged`; `TestInstallUninstallClaudeDashCUsesProjectRoot`; DF-22/37 keepers `TestInstallCursorPrintReloadTip` / `TestInstallCursorWriteMergeBackup`; `TestCLIAddDeniedFailClosed`. Library: `TestInstallDetectListsCursorStable` / `TestInstallConditional*` |
| S05 DF-70…74 | **APPROVE high** (P16-S05-02) | `TestSeedImportDiscoveryMentionsTask`; `TestContextIncludesImpactOverallClass`; `TestWhyIncludesImpactOverallClass`; `TestSeedImportImpactFindings`; `TestImpactReportJSONSnakeCase`; `TestMCPTraceImpactReport`; `TestMCPImpactDeniedBlocksCallTool`; catalog **10** `TestToolNamesRegistered` / `TestBuiltinMCPCapabilitySpecs`; boundary `TestImportBoundaryMCPNoPlanImpactIndexTools` (**allows** `trace_impact`; **forbids** plan/index/install/decide). CLI keepers: `TestCLIAddDeniedFailClosed` / `TestLinkDiscoveryMentionsTaskCLI` / `TestInstallClaudeDashCRefuseCitesProjectRoot`. Gate F keeper in carry-forward. |

## Live residuals → DR-HANDOFF decision (2026-08-17)

| Bucket | Items | Phase implication |
|--------|-------|-------------------|
| Product gaps scheduled in Phase 16 | DF-76/75/78/77/68/70/71/72/73/74 | Closed by S01–S05 APPROVE — VERIFY must **re-prove named tests** (DF-72 `trace_impact` **is** a fail bar) |
| Explicit residual OK into VERIFY | DF-67 defer; P14 R2 defer; P15 R3/R4 wontfix; S05-02 `attachTaskImpact` swallow; 014 nine-Name `IN` list historical; seed summary 0-counts; DF-22/37 ops (tip keepers only) | **Do not fail VERIFY** for these |
| Independently queued Phase 17 | Portable graph git — board rows **232–244 already exist** | **Do not rewrite P17.** **Do not** claim P17 as this VERIFY successor |
| Goals sequence #2–#4 | Research S05 / `plan simulate` / D21+ | Stay off-board — **not** auto-boarded |
| Product bar | `./cmd\|internal\|evals` with CGO1 | Prefer product pkgs over full-module `./...` when graphify space FAIL present |

**DR-HANDOFF = `no successor`.** Phase 17 is a separate human queue after P16 rows; S06 Notes must not promote it. Reopen research S05 / plan simulate / D21+ only with explicit human promotion.

## Planner work
1. [x] Confirm live named tests exist; import S01–S05 into 01-verify
2. [x] Lock verify command set + DR-HANDOFF = `no successor`
3. [x] Thicken 01-verify / 02-scope-review / SCOPE-TODOS
4. [x] Mark this prompt **FINAL**; next **P16-S06-01**

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Phase 16 assert-root-and-surfaces closeout — S01–S05 named DF regressions — **not** a new planted eval gate |
| Catalog | **Ten** MCP names **including** `trace_version`: why, context, add, link, transition, review, tasks, capability, **`trace_impact`**, version. Slug `mcp:trace_impact`. **No** install/decide/plan/index MCP |
| Migration | **None** from VERIFY — mig **014** already landed; compat ceiling **14** (no 015+) |
| Carry-forward | honesty A/B/C+G; E/F; ablation; H; compat **14**; p0x; x0; Gate C `dry_run:false`; product `./cmd\|internal\|evals` |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist |
| Full bar | `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1` — **product pkgs PASS**; R3 graphify space FAIL on full `./...` OK; R4 CGO0 analyzers / `cmd/trace` CGO0 FAIL OK |
| Residuals OK | DF-67 defer; P14 R2 defer; P15 R3/R4 wontfix; S05-02 `attachTaskImpact` swallow; 014 nine-Name list; DF-22/37 tip-only — **do not fail VERIFY** |
| Allowed Go on VERIFY | **None** for features — re-run + evidence docs only; spawn remediation if fail |
| Optional strong evidence | Grep `registerTools` / `RegisteredToolNames` (10, impact before version); Gate C artifact inspect; G19 — non-blocking |
| Spawn | On fail: `01a` implement / `01b` review (+`01c` re-VERIFY if needed) immediately below |
| DR-HANDOFF | **`no successor`** — **S06-01 starts** Notes; **S06-02 owns completion**. Phase 17 independently queued — **do not rewrite**. Do **not** auto-board research S05 / plan simulate / D21+ |
| Forbidden | Product features on VERIFY; claiming DF-67/R2/R3/R4 fixed; rewriting P15 history; Mode-B Gate C rewrite; daemon/HTTP/embeddings; full-rebuild indexer; YOLO/AllowAll; new MCP install/decide/plan/index tools; CGO0 `./cmd/trace/...` as a fail bar |

### Locked verify command set (FINAL)

Per-scope named `-run` lines match S01–S05 REVIEW-NOTES (independent DF re-proof). Overlap is intentional. **CGO0 `./cmd/trace/...` is R4 — do not use.** `GOMODCACHE`+`GOPROXY=off` on mcp/library lines (S03 sandbox 403 class).

```bash
# --- S01 DF-76 named + P15 keepers ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered|TestMCPAssert|TestMCPVirgin|TestMCPInitialized|TestOpenExisting|TestOpenCreates|TestBuiltinMCPCapabilitySpecs|TestCapabilityDecision'

# --- S02 DF-75/78 CHECK + mcp: canonicalize; compat 14 keepers ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestCapabilityToolDecisionCheckRejectsYOLO|TestCapabilityToolDecisionMigrateHealsYOLOToPending|TestResolveYOLOBuiltinDoesNotAutoAllow|TestDecideUnprefixedMCPNameCanonicalizes|TestCanonicalizeCustomAndCLISlugsUnchanged|TestMigrateUnprefixedDeniedFoldsOverAutoAllowed|TestMCPUnprefixedDecideGatesCallTool|TestMCPAssert|TestToolNamesRegistered|TestOpenCreates|TestMCPVirgin'

# --- S03 DF-77 dual-slug CLI-first ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestCLIAddDeniedDoesNotBlockMCPAdd|TestUnprefixedAddDecideDoesNotGateCLI|TestCapabilityDecisionAutoAllowBuiltinCLI|TestCanonicalizeCLIReindexFoldsToIndex|TestCanonicalizeCustomAndCLISlugsUnchanged|TestCapabilityDecision|TestMCPAssert|TestMCPUnprefixed|TestToolNamesRegistered|TestMCPVirgin|TestOpenCreates'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestCLIAddSucceedsWhenMCPAddDenied|TestCLIAddDeniedFailClosed|TestCLIWhySucceedsWhenMCPWhyDenied|TestCLIWhyDeniedFailClosed|TestUngatedCapabilityDecideWhenCLIAddDenied|TestCLIIndexAliasDenied|TestUnprefixedAddDecideDoesNotGateCLI'

# --- S04 DF-68 -C ProjectRoot + DF-22/37 tip keepers ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/install/... ./internal/mcp/... ./internal/store/... -count=1 -run 'TestInstallDetectListsCursorStable|TestInstallConditional|TestToolNamesRegistered|TestMCPVirgin|TestOpenCreates'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallClaudeDashC|TestInstallDetectDashC|TestInstallUninstallClaudeDashC|TestInstallCursorPrintReloadTip|TestInstallCursorWriteMergeBackup|TestCLIAddDeniedFailClosed'

# --- S05 DF-70/71/72/73/74 + catalog 10 + boundary ---
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/compiler/... ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestContextIncludesImpactOverallClass|TestMCPTraceImpactReport|TestMCPImpactDeniedBlocksCallTool|TestToolNamesRegistered|TestBuiltinMCPCapabilitySpecs|TestImportBoundaryMCPNoPlanImpactIndexTools|TestMCPVirgin|TestOpenCreates'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedImportDiscoveryMentionsTask|TestSeedImportImpactFindings|TestImpactReportJSONSnakeCase|TestWhyIncludesImpactOverallClass|TestCLIAddDeniedFailClosed|TestLinkDiscoveryMentionsTaskCLI|TestInstallClaudeDashCRefuseCitesProjectRoot'

# Honesty: Paths A/B/C + Gate G (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# Gate E / F / capability ablation
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# Gate H + compat (compat covers mig 014 ceiling)
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# P0-X + X0
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1

# Product regression bar (prefer over full-module ./... when R3 graphify present)
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Optional (strong evidence, **not** substitutes for package PASS / not Mode-B Gate C):

```bash
# Grep: RegisteredToolNames / registerTools — exactly 10; trace_impact before trace_version; no trace_plan/trace_index/trace_install/trace_decide
# Gate C artifact inspect (jq/grep OK): docs/verification/gate-c-x0/ metrics dry_run:false, N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# Spot-check only (non-blocking): DF-67 still out of index_honesty; R2 allowContainsOut still in impact_walk.go; 014 Name IN list still nine historical Names
# Phase 17 board rows 232–244 already exist — do not rewrite; do not claim as this successor
# Do NOT fail for R3 graphify space FAIL or R4 CGO0 analyzers / CGO0 cmd/trace FAIL
# Do NOT fail for attachTaskImpact swallow / seed 0-counts / DF-22/37 ops
```

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable (thickened)
- [x] VERIFY commands + DR-HANDOFF locked (`no successor`)
- [x] SCOPE-TODOS + board Notes; next `P16-S06-01`
- [x] Product Go — **not** this row

## Out of scope
- Running VERIFY (S06-01)
- Product Go / new MCP tools / daemon / mig 015
- Rewriting Phase 17 board rows or claiming P17 as this successor
- Auto-boarding research S05 / plan simulate / D21+
- Claiming DF-67 / P14 R2 / P15 R3/R4 fixed
- Closing parallel dogfood experiments
- Claiming Phase 15 historical handoff was wrong

## Next
**P16-S06-01** (independent VERIFY run).
