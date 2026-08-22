# P22-S08-07 — Verify Phase 22 (capability completion)

## Metadata
- id: P22-S08-07
- todo_ids: [P22-S08-07]
- role: verify
- skills: [writing-for-agents, planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective

Re-read [`docs/CAPABILITIES_CHECKLIST.md`](../../../../CAPABILITIES_CHECKLIST.md) and README **C01–C43** matrix. Prove every capability is `[x]` with evidence, **or** spawn in-phase remediation rows. Archive evidence. **Does not** close DR-HANDOFF (S08-08 owns). **Does not** start Phase 23 for leftover capabilities.

## Session start

Follow [agent-loop-protocol.md](../../../../../rules/agent-loop-protocol.md). **Unattended:** execute command floor + matrix — do not stop after planning.

## Preconditions (hard)

- [ ] **S09 complete** on board (S09-00…S09-08 `done`) — if not, **FAIL** this row immediately with Notes “blocked: S09 incomplete”; do not mark `done`.
- [ ] Re-read active board after any spawn.

## References

- [00-PLANNER.md](00-PLANNER.md)
- [VERIFY.md](../../VERIFY.md)
- [README.md](../../README.md) coverage matrix (C01–C43)
- [DECISION-LOG.md](../../DECISION-LOG.md)
- Checklist SoT: [`docs/CAPABILITIES_CHECKLIST.md`](../../../../CAPABILITIES_CHECKLIST.md)

## Locked defaults

| Item | Value |
|------|-------|
| Schema max | **027** after S09 (`027_harness_agents.sql`); **no 028+** |
| Compat ceiling | **27** via `TestCompatibilitySecurityChecklist` |
| Checklist state | Expect **141/141 `[x]`** if S08+S09 closed cleanly; **3** were open at S08-00 (C28, C38, C39) |
| Evidence dir | `experiments/runs/2026-08-18-p22-s08-07-verify/evidence/` |
| Notes | `VERIFY-NOTES.md` in this folder (required) |
| Checklist edits | May check boxes **only** with evidence citations; `[ ]` without spawn = **FAIL** |
| E01–E04 | Must appear in VERIFY-NOTES with PASS + evidence (harness recommendations, catalog, CLI/MCP agents) |
| Out-of-scope | Hosted MCP, daemon, HTTP, wrap-git-commit, ML — **out**, not leftovers |
| DR-HANDOFF | Leave **OPEN** |

## Checklist gate (mandatory)

1. Open `docs/CAPABILITIES_CHECKLIST.md` — count `[ ]` (expect **0** on pass).
2. Walk README matrix **C01–C43** — each row: PASS + test/evidence path **or** spawn id `pending`/`in_progress`.
3. Walk **E01–E04** — PASS with S09 evidence or spawn on board.
4. **FAIL** if any matrix row lacks evidence and lacks runnable in-phase spawn.
5. On FAIL: spawn `P22-S08-07a` implement + `07b` review (or capability-specific Na/Nb) **on this board**; leave this row `in_progress` or `failed`.

## Locked command floor (report PASS/FAIL each block)

```bash
# Compat + schema (expect 27 after S09)
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/store/... -count=1 -run 'TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestNoSourceContentColumns'
ls internal/store/schema/*.sql | wc -l  # expect 27

# S01 graph
go test ./internal/analyzers/... ./internal/store/... ./internal/retrieval/... -count=1 -run 'TestIndexDiscoversGoTestFunctions|TestValidatesEdgeExtractedFromImport|TestArtifactEdgesFunctionsTypesAPIs|TestArchitecturalBoundaryEdges|TestImpactWalkIncludesAffectedTests'

# S02 sync/change
go test ./internal/install/... ./internal/domain/... ./internal/compiler/... -count=1 -run 'TestInstallGitHook|TestGraphSyncStaleWhenHeadDiffers|TestPromoteVCSCommit|TestCompareStates'

# S03 cycle/verify
go test ./internal/loop/... ./internal/testrun/... ./internal/domain/... -count=1 -run 'TestBuildPolicyInputs|TestTestRunRecordsOutcome|TestVerificationCycle|TestCoordinateVerification|TestInvariant|TestCompareIterationOutcomes'

# S04 impact/regression/improvements
go test ./internal/domain/... -count=1 -run 'TestRecordPredictedImpact|TestRegressionLinkedToChange|TestRecordImprovement'

# S05 query
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/mcp/... -count=1 -run 'TestCLISearch|TestCLIChangesList|TestTestsVerifying|TestRegressionsList|TestContextIncludesEvaluations|TestToolNamesRegistered'

# S06 knowledge
go test ./internal/domain/... ./internal/loop/... -count=1 -run 'TestPatternCounts|TestQuerySimilarChanges|TestSynthesizeKnowledge|TestTendHelpHurt|TestSuccessfulApproaches'

# S07 eval
go test ./internal/eval/... -count=1 -run 'TestEvalRegistry|TestProjectEvalRules|TestListEvaluationResults'

# S08 workflow
go test ./internal/mcp/... ./internal/domain/... ./cmd/trace/... -count=1 -run 'TestMCPLoop|TestDetectOverlapping|TestHelpIncludesSearchTestVerify|TestTasksConflicts|TestLoopNextIncludesWorkConflicts'

# S09 harness agents (must be done before this row)
go test ./internal/agents/... ./internal/install/... ./internal/loop/... ./cmd/trace/... ./internal/mcp/... -count=1 -run 'TestHarnessAgent|TestInstallAgents|TestLoopNextIncludesHarness|TestCLIAgents|TestMCPAgents'

# P21 keepers (must stay green)
go test ./internal/deliberation/... -count=1
go test ./internal/loop/... -count=1 -run 'TestLoopApply'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedExportRoundTrip|TestSeedExportOmitsDeniedSurfaces|TestLoopNextDeliberationSectionPresent'
```

If a `-run` matches **no tests**, that is a **FAIL** for the owning capability (not a skip). Spawn remediation.

## CLI smoke (capture to evidence dir)

1. `trace install detect` includes `git-hook`
2. `trace search` / `trace changes list` / `trace regressions list` / `trace tests verifying` (fixture)
3. `trace test run --task` on fixture (or documented recorded invoke)
4. `trace loop next` shows cycle flags / planning_evidence / **work_conflicts** when seeded
5. `trace tasks conflicts` when overlap fixture seeded
6. `trace seed export` omits index blobs; includes knowledge/improvements/eval-rules pointer when populated
7. `./bin/trace-mcp -h` lists **15** tools: loop, search/history, **`trace_agents`**
8. `trace install agents` + `trace agents recommend --phase CRITIQUE` — recommendations only, no spawn
9. `trace_version` via MCP after rebuild

## VERIFY-NOTES template

```markdown
# P22-S08-07 VERIFY-NOTES

## Summary
PASS | FAIL — N matrix rows open, M checklist `[ ]`

## Command floor
| Block | Result | Notes |

## C01–C43 matrix
| ID | Owner | Evidence | Result |

## E01–E04
| ID | Evidence | Result |

## Checklist
- Open `[ ]` count: 
- Spawns created: 

## DR-HANDOFF
Status: OPEN (S08-08 closes)
```

## Exit criteria

- [ ] VERIFY-NOTES.md complete
- [ ] Command floor recorded (all blocks)
- [ ] Checklist **141/141 `[x]`** OR in-phase spawns exist and are runnable
- [ ] E01–E04 evidenced
- [ ] DR-HANDOFF still OPEN
- [ ] Board Notes: evidence dir path + FAIL list + next spawn ids

## Minimal todos

- [ ] Confirm S09 done
- [ ] Run command floor
- [ ] Matrix C01–C43 + E01–E04
- [ ] Update checklist boxes with evidence **or** spawn
- [ ] VERIFY-NOTES.md + board notes
