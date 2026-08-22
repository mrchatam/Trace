# S06 — Phase VERIFY — scope todos

**Depends-on:** P16-S05-02 APPROVE. Owns Phase 16 close + DR-HANDOFF.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | VERIFY planner | **done** — FINAL locks + DR-HANDOFF = `no successor` |
| 2 | 01-verify | verifier | **done** — Phase 16 VERIFY PASS; DR-HANDOFF started = `no successor`; next **P16-S06-02** |
| 3 | 02-scope-review | review / handoff close | **done** — APPROVE high; DR-HANDOFF closed = `no successor`; Phase 16 complete; Phase 17 also complete |

## Depends (from S01–S05 — imported at S06-00; live `func Test*` confirmed 2026-08-17)

VERIFY must include (named `-run` in `00-PLANNER.md` / `01-verify.md`):

- **S01 DF-76:** `TestMCPVirginProjectDoesNotMkdir`, `TestMCPInitializedOtherRootIsolated`, `TestOpenExistingMissingReturnsErrNotInitialized`, `TestOpenExistingEmptyTraceDir` + P15 keepers `TestMCPAssertDeniedBlocksCallTool` / `TestMCPAssertBuiltinAutoAllowedSucceeds` / `TestToolNamesRegistered` (**live length 10**)
- **S02 DF-75/78:** `TestCapabilityToolDecisionCheckRejectsYOLO`, `TestCapabilityToolDecisionMigrateHealsYOLOToPending`, `TestResolveYOLOBuiltinDoesNotAutoAllow`, `TestDecideUnprefixedMCPNameCanonicalizes`, `TestCanonicalizeCustomAndCLISlugsUnchanged`, `TestMigrateUnprefixedDeniedFoldsOverAutoAllowed`, `TestMCPUnprefixedDecideGatesCallTool` + keepers `TestMCPAssert*` / `TestOpenCreates*` (v**14**) / `TestMCPVirgin*`
- **S03 DF-77:** `TestCLIAddSucceedsWhenMCPAddDenied`, `TestCLIAddDeniedFailClosed`, `TestCLIWhySucceedsWhenMCPWhyDenied`, `TestCLIWhyDeniedFailClosed`, `TestCLIAddDeniedDoesNotBlockMCPAdd`, `TestUngatedCapabilityDecideWhenCLIAddDenied`, `TestUnprefixedAddDecideDoesNotGateCLI`, `TestCapabilityDecisionAutoAllowBuiltinCLI`, `TestCanonicalizeCLIReindexFoldsToIndex`, `TestCLIIndexAliasDenied` (CLI under **CGO1**)
- **S04 DF-68:** `TestInstallClaudeDashCRefuseCitesProjectRoot`, `TestInstallClaudeDashCIgnoresCwdMarker`, `TestInstallClaudeDashCWriteUsesProjectRoot`, `TestInstallDetectDashCClaudeReasonCitesRoot`, `TestInstallDetectDashCCursorHomeUnchanged`, `TestInstallUninstallClaudeDashCUsesProjectRoot`; DF-22/37 keepers `TestInstallCursorPrintReloadTip`, `TestInstallCursorWriteMergeBackup`; library `TestInstallDetectListsCursorStable` / `TestInstallConditional*`
- **S05 DF-70…74:** `TestSeedImportDiscoveryMentionsTask`, `TestContextIncludesImpactOverallClass`, `TestWhyIncludesImpactOverallClass`, `TestSeedImportImpactFindings`, `TestImpactReportJSONSnakeCase`, `TestMCPTraceImpactReport`, `TestMCPImpactDeniedBlocksCallTool`; catalog **10** `TestToolNamesRegistered` / `TestBuiltinMCPCapabilitySpecs`; boundary `TestImportBoundaryMCPNoPlanImpactIndexTools`; CLI keepers `TestCLIAddDeniedFailClosed` / `TestLinkDiscoveryMentionsTaskCLI` / `TestInstallClaudeDashCRefuseCitesProjectRoot`; Gate F keeper in carry-forward
- **Carry-forward:** honesty A/B/C+G; Gate E/F; ablation; Gate H; compat **14**; p0x; x0; Gate C `dry_run:false`; product `./cmd|internal|evals`

## Residuals (non-blocking — do **not** fail VERIFY)
- **DF-67 defer:** symbol-entity honesty out of `index_honesty` bar
- **P14 R2 defer:** `allowContainsOut` late-upgrade — Notes only
- **P15 R3 wontfix:** `similar projects/graphify` space-in-path on full `./...`
- **P15 R4 wontfix:** CGO0 analyzers FAIL (tree-sitter) — product bar is CGO1; CGO0 `cmd/trace` FAIL OK
- **S05-02:** `attachTaskImpact` swallows helper errors (capability-style); 014 SQL nine-Name `IN` list historical (do not edit); seed summary always emits `findings`/`alternatives` ints (0 if none)
- **DF-22/37:** manual Cursor reload — tip keepers only

## Reminders
- Default DR-HANDOFF = **`no successor`** — **S06-01 started** Notes; **S06-02 completed** close
- **Phase 17** is independently queued after P16 rows (board 232–244) — **not** rewritten; **not** this VERIFY successor; do **not** auto-board research S05 / plan simulate / D21+
- Dry-run ≠ Gate C / F / G / ablation / H / checklist
- Catalog **10** including `trace_version`; `trace_impact` **in** catalog; no install/decide/plan/index MCP
- Phase 15 historical handoff stays intact
- Spawn on fail: `P16-S06-01a` / `01b` (+`01c`) immediately below (unused — VERIFY green)
- Next after Phase 16 complete: **none** unless human promotes follow-on (Phase 17 also closed 2026-08-17; independent queue — **not** DR-HANDOFF promotion)
