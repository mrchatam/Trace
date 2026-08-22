# P22-S08-07 VERIFY-NOTES

## Summary
**PASS** — 0 matrix rows open, 0 checklist `[ ]` (141/141 `[x]`). E01–E04 evidenced. One **FYI (low):** committed `bin/trace-mcp -h` lists **14** tools (pre-S09-07 binary); code registration **15** incl. `trace_agents` via `TestToolNamesRegistered` PASS.

## Preconditions
- **S09 complete:** S09-00…S09-08 all `done` on board — **PASS**
- **Schema:** 27 SQL files (`027_harness_agents.sql`); compat ceiling **27** — **PASS**

## Command floor
| Block | Result | Notes |
|-------|--------|-------|
| Compat + schema | **PASS** | `evidence/01-compat.log`, `02-store.log`; 27 SQL files |
| S01 graph | **PASS** | `evidence/03-s01-graph.log` |
| S02 sync/change | **PASS** | `evidence/04-s02-sync.log` |
| S03 cycle/verify | **PASS** | `evidence/05-s03-cycle.log` |
| S04 impact/regression | **PASS** | `evidence/06-s04-impact.log` |
| S05 query | **PASS** | `evidence/07-s05-query.log` |
| S06 knowledge | **PASS** | `evidence/08-s06-knowledge.log` |
| S07 eval | **PASS** | `evidence/09-s07-eval.log` |
| S08 workflow | **PASS** | `evidence/10-s08-workflow.log` |
| S09 harness agents | **PASS** | `evidence/11-s09-harness.log` (`TestHarnessAgent*` in store/domain; `internal/agents` lib-only) |
| P21 deliberation | **PASS** | `evidence/12-p21-deliberation.log` |
| P21 loop apply | **PASS** | `evidence/13-p21-loop-apply.log` |
| P21 seed/loop | **PASS** | `evidence/14-p21-seed-loop.log` |

## CLI smoke
| # | Check | Result | Evidence |
|---|-------|--------|----------|
| 1 | `install detect` includes `git-hook` | **PASS** | `evidence/16-install-detect.log` |
| 2 | search/changes/regressions/tests verifying help | **PASS** | `evidence/17-cli-search-help.log` + S05 tests |
| 3 | `trace test run --task` | **PASS** | `TestTestRunRecordsOutcome` (S03 floor) |
| 4 | `loop next` cycle flags / planning_evidence / work_conflicts | **PASS** | `TestLoopNextIncludesWorkConflicts`, `TestBuildPolicyInputs`, `TestLoopNextDeliberationSectionPresent` |
| 5 | `tasks conflicts` | **PASS** | `TestTasksConflictsCLI` (S08 floor) |
| 6 | `seed export` omits index; includes harness_agents | **PASS** | `evidence/23-seed-export.log`; `TestSeedExportOmitsDeniedSurfaces` |
| 7 | `./bin/trace-mcp -h` lists **15** tools incl. `trace_agents` | **FYI** | Live `-h` shows **14** (`evidence/15-mcp-help.log`); `TestToolNamesRegistered` PASS (15 in code). Rebuild blocked (proxy); S09-08 noted rebuild cmd |
| 8 | `install agents` + `agents recommend --phase CRITIQUE` | **PASS** | `evidence/21-22*.log` — recommend-only JSON, no spawn |
| 9 | `trace_version` via MCP | **PASS** | `TestTraceVersion` in `internal/mcp/mcp_test.go` (S05/S08 keeper runs) |

## C01–C43 matrix
| ID | Owner | Evidence | Result |
|----|-------|----------|--------|
| C01 | S01-01/03 | `TestIndexDiscoversGoTestFunctions`, `TestArtifactEdgesFunctionsTypesAPIs`, `TestValidatesEdgeExtractedFromImport` | PASS |
| C02 | S01-05 | `TestArchitecturalBoundaryEdges` | PASS |
| C03 | S01-01 | `TestValidatesEdgeExtractedFromImport`, `TestIndexDiscoversGoTestFunctions` | PASS |
| C04 | S02-01/03 | `TestInstallGitHook`, `TestGraphSyncStaleWhenHeadDiffers` | PASS |
| C05 | S02-05 | `TestPromoteVCSCommit` | PASS |
| C06 | S02-07 | `TestCompareStates` | PASS |
| C07 | S01-07 | `TestImpactWalkIncludesAffectedTests` | PASS |
| C08 | S04-01 | `TestRecordPredictedImpact` | PASS |
| C09 | S03-01 + S09 | `TestBuildPolicyInputs`, `TestLoopNextIncludesHarnessRecommendations` | PASS |
| C10 | S06-03 | `TestSynthesizeKnowledge` | PASS |
| C11 | S03-05 | `TestVerificationCycle`, `TestCoordinateVerification` | PASS |
| C12 | S03-03 | `TestTestRunRecordsOutcome` | PASS |
| C13 | S03-05 | `TestVerificationCycle` | PASS |
| C14 | S03-07 | `TestInvariant` | PASS |
| C15 | S03-07 | `TestCompareIterationOutcomes` | PASS |
| C16 | S04-03 | `TestRegressionLinkedToChange` | PASS |
| C17 | S05-03 | `TestRegressionsList` | PASS |
| C18 | S04-05 | `TestRecordImprovement` | PASS |
| C19 | S06-01 | `TestPatternCounts` | PASS |
| C20 | S06-01 | `TestQuerySimilarChanges` | PASS |
| C21 | S06-03 | `TestSynthesizeKnowledge` | PASS |
| C22 | S06-05 | `TestTendHelpHurt` | PASS |
| C23 | S06-05 | `TestSuccessfulApproaches` | PASS |
| C24 | S06-05 | `TestContextIncludesEvaluations` (S05-05) | PASS |
| C25 | S02-01 | `TestInstallGitHook` | PASS |
| C26 | S06-03 | `TestSynthesizeKnowledge` | PASS |
| C27 | S06-03 | `TestSynthesizeKnowledge` | PASS |
| C28 | S08-03 | `TestDetectOverlapping`, `TestLoopNextIncludesWorkConflicts`, `TestTasksConflictsCLI` | PASS |
| C29 | S05-01 | `TestCLISearch`, `TestContextIncludesEvaluations` | PASS |
| C30 | S05-01 | `TestCLIChangesList` | PASS |
| C31 | S05-03 | `TestTestsVerifying` | PASS |
| C32 | S05-03 | `TestRegressionsList` | PASS |
| C33 | S05-03 | `TestSuccessfulApproaches` | PASS |
| C34 | S05-03 | `TestRegressionsList` | PASS |
| C35 | S05-05 | `TestContextIncludesEvaluations` | PASS |
| C36 | S03-05 | `TestCoordinateVerification` | PASS |
| C37 | S05-01 | `TestCLISearch`, MCP context tools | PASS |
| C38 | S03-03 + S08-01 | `TestTestRunRecordsOutcome`, `TestMCPLoop`, `TestHelpIncludesSearchTestVerify` | PASS |
| C39 | S08-01/05 | `TestMCPLoop`, `TestHelpIncludesSearchTestVerify`, harness workflow | PASS |
| C40 | S07-01 | `TestEvalRegistry` | PASS |
| C41 | S07-03 | `TestProjectEvalRules` | PASS |
| C42 | S07-05 + S05-05 | `TestListEvaluationResults`, `TestContextIncludesEvaluations` | PASS |
| C43 | S07-01 | `TestEvalRegistry` (additive contract) | PASS |

## E01–E04
| ID | Evidence | Result |
|----|----------|--------|
| E01 | `TestRecommendSubagentWhenAvailable`, `TestRecommendSubagentHonestWhenUnavailable`, `TestLoopNextIncludesHarnessRecommendations` | PASS |
| E02 | `TestCLIAgentsRecommend`, `TestMCPAgentsRecommend`, `TestRecommendPerformanceReviewerForPerfTask`, `agents recommend --phase CRITIQUE` smoke | PASS |
| E03 | `TestHarnessAgentCatalogMigrate027`, `TestInstallAgentsSeedsDefaults`, `trace/agents/default.json`, seed export 6 agents | PASS |
| E04 | `trace/agents/README.md`, `registry_*` columns, `TestHarnessAgentSeedExportRoundTrip` | PASS |

## Checklist
- Open `[ ]` count: **0** (141/141 `[x]` in `docs/CAPABILITIES_CHECKLIST.md`)
- Spawns created: **none**

## Artifacts
- Evidence dir: `experiments/runs/2026-08-18-p22-s08-07-verify/evidence/`
- Refreshed `trace/graph.json` (6 `harness_agents` after local `install agents`) — `evidence/24-graph-export.log`

## DR-HANDOFF
Status: **OPEN** (S08-08 closes)

## Residuals (non-blocking)
- Rebuild `bin/trace-mcp` so live `-h` lists 15 tools (code already correct).
- `install detect` git-hook `detected: false` in sandbox (git unavailable) — id present in output.
