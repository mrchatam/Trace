# P23-S06-01 — Phase 23 verify

## Metadata
- id: P23-S06-01
- todo_ids: [P23-S06-01]
- role: verify
- skills: [incremental-implementation, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective

Run the **locked verify floor** (S01–S05 deltas + P19/P20 loop keepers + compat **027**), archive CLI smoke evidence under `experiments/runs/YYYY-MM-DD-p23-s06-01-verify/`, and map results to [ENFORCEMENT.md](../../ENFORCEMENT.md) evidence bar. Fix only **blocking** test gaps found in verify (minimal). **Does not** close DR-HANDOFF (S06-02 owns).

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [00-PLANNER.md](00-PLANNER.md) — FINAL evidence bar (this scope)
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S06-02

## Session start

Follow agent-loop-protocol. Unattended: do not stop after planning. This row verifies and records evidence; it does not open new product direction beyond Phase 23 locks.

## Locked defaults (FINAL — S06-00)

| Item | Value |
|------|-------|
| Schema max | **027** (`027_harness_agents.sql`; **no 028+**) |
| Compat ceiling | **27** — re-lock at VERIFY via `TestCompatibilitySecurityChecklist` |
| MCP tools | **15** — no new gate MCP |
| Gate envelope | `trace.loop.gate.v1` — exit **0** allowed / **1** blocked / **2** usage-internal |
| Status schema | `trace.loop.status.v1` — additive `violations[]` only |
| Config | `.trace/config.json` `{ "enforce": "off"\|"warn"\|"strict" }`; default **off**; no auto-enforce on transition/export |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p23-s06-01-verify/evidence/` |
| Notes artifact | `VERIFY-NOTES.md` in this scope folder (recommended) |
| Out-of-scope | No product Go beyond blocking fixes; no DR-HANDOFF close; no Phase 24 scaffold |

## Must checklist (ENFORCEMENT.md → evidence)

S06-01 Notes must cite **PASS** evidence for each row:

- [ ] Gate blocks edit with blocking uncertainty (exit **1** + JSON) — `TestLoopGateBlockedExitOne`, CLI #1
- [ ] Gate allows when policy clear (exit **0**) — `TestLoopGateAllowedExitZero`
- [ ] Status `violations[]` matches gate for same task — `TestLoopStatusViolationsMatchGateEdit`, CLI #2
- [ ] `transition … DONE --enforce` fails with verification debt — `TestTransitionDoneEnforceBlocksVerificationDebt`, CLI #3
- [ ] `seed export --strict --enforce` fails on violations (no write) — `TestSeedExportStrictEnforceNoWriteOnViolation`, CLI #4
- [ ] Default config enforce **off** — `TestTraceConfigEnforceDefaultOff`, `TestLoadEnforceModeMissingFile`
- [ ] Config `strict` does **not** auto-enforce transition/export — plain DONE/export unchanged with config strict
- [ ] `trace install cursor --write` includes enforcement rules — `TestInstallCursorIncludesLoopGateRule`, CLI #5
- [ ] git-hook unchanged (post-commit still works) — `TestInstallGitHookUnchanged`, CLI #6
- [ ] P19/P20 loop keepers PASS — Block F
- [ ] Compat ceiling **027** PASS — Block G
- [ ] No daemon / hosted MCP introduced — grep + review handoff

## Locked verify command floor (FINAL)

Run every block; capture PASS/FAIL in `99-run-metadata.txt`. `-count=1` recommended.

### Block A — S01 gate library (17 + 1 named)

```bash
go test ./internal/loop/... ./internal/domain/... -count=1 -run 'TestEvaluateGate|TestPrematureImplementation_Code'
```

Named tests (all must exist):

| Test | Proves |
|------|--------|
| `TestEvaluateGate_Orient_UnknownTask` | orient fail |
| `TestEvaluateGate_Orient_MissingGoal` | `missing_goal_id` |
| `TestEvaluateGate_Orient_MissingPlan` | `missing_plan_context` |
| `TestEvaluateGate_Orient_OK` | orient allow |
| `TestEvaluateGate_Edit_BlockingUncertainty` | edit block → INVESTIGATE |
| `TestEvaluateGate_Edit_OpenRegression` | edit block → regression |
| `TestEvaluateGate_Edit_PlanMissing` | edit block → PLAN |
| `TestEvaluateGate_Edit_PlanUncritiqued` | edit block → CRITIQUE |
| `TestEvaluateGate_Edit_ExecuteReady` | edit allow |
| `TestEvaluateGate_Edit_HopBudgetStopped` | STOP |
| `TestEvaluateGate_Execute_NotExecutePending` | execute stricter |
| `TestEvaluateGate_Execute_ExecutePendingClear` | execute allow |
| `TestEvaluateGate_Done_VerificationDebt` | done block |
| `TestEvaluateGate_Done_OpenRegression` | done block |
| `TestEvaluateGate_Done_DeliberationIncomplete` | done block |
| `TestEvaluateGate_Done_Clean` | done allow |
| `TestEvaluateGate_Export_SameAsDone` | export parity |
| `TestPrematureImplementation_Code` | stable code string |

Policy reuse keeper:

```bash
go test ./internal/deliberation/... -count=1 -run 'TestSelectNext'
```

### Block B — S02 gate CLI (14 named)

```bash
go test ./cmd/trace -count=1 -run 'TestLoopGate|TestHelpIncludesLoopGate'
```

| Test | Proves |
|------|--------|
| `TestLoopGateAllowedExitZero` | exit 0 + `allowed=true` |
| `TestLoopGateBlockedExitOne` | exit 1 (not usage) |
| `TestLoopGateJSONSchemaVersion` | `trace.loop.gate.v1` |
| `TestLoopGateTopLevelLiftFromViolation` | top-level lift from `violations[0]` |
| `TestLoopGateAllowedEmptyViolations` | empty array not null |
| `TestLoopGateBlockedStderrHint` | stderr hint on block |
| `TestLoopGateDefaultForEdit` | default `--for edit` |
| `TestLoopGateInvalidForFailClosed` | exit 2 invalid `--for` |
| `TestLoopGateMissingTaskFlag` | exit 2 missing `--task` |
| `TestLoopGateUnknownTaskOrientBlocked` | orient block unknown task |
| `TestLoopGateOrientAllowed` | orient allow |
| `TestLoopGateDoneBlockedVerificationDebt` | `--for done` block |
| `TestLoopGateExecuteAllowedWhenPending` | `--for execute` allow |
| `TestHelpIncludesLoopGate` | help documents subcommand |

### Block C — S03 enforce DONE + export strict (17 named)

```bash
go test ./cmd/trace -count=1 -run 'TestTransitionDoneEnforce|TestSeedExportStrict|TestReviewCreateSetDone|TestAllowDoneWarnsOnStderr|TestHelpIncludesTransitionEnforce|TestHelpIncludesSeedExportStrict'
```

Transition (8):

| Test | Proves |
|------|--------|
| `TestTransitionDoneEnforceBlocksVerificationDebt` | enforce block |
| `TestTransitionDoneWithoutEnforceUnchanged` | default unchanged |
| `TestTransitionDoneEnforceAllowsClean` | enforce allow clean |
| `TestTransitionDoneEnforcePreservesAllowDone` | `--allow-done` without enforce |
| `TestTransitionDoneEnforceBlocksDespiteAllowDone` | gate-first with both flags |
| `TestTransitionDoneEnforceIgnoredForNonDone` | enforce no-op non-DONE |
| `TestTransitionDoneEnforcePreservesDomainReviewGate` | domain review unchanged |
| `TestTransitionDoneEnforceStderrHint` | stderr on block |

Export (9):

| Test | Proves |
|------|--------|
| `TestSeedExportStrictEnforceNoWriteOnViolation` | no write on block |
| `TestSeedExportStrictWithoutEnforceExitZero` | strict alone exit 0 |
| `TestSeedExportStrictEnforceBlocksOpenRegression` | regression block |
| `TestSeedExportStrictCleanAllowsWrite` | clean write |
| `TestSeedExportStrictTaskFilter` | `--task` filter |
| (+ round-trip keepers) | `TestSeedExportRoundTrip`, `TestSeedExportOmitsDeniedSurfaces`, `TestSeedExportWritesExportedAtCommit` |

Two-layer DONE walkthrough (document in VERIFY-NOTES):

1. Verification debt + review PASS + `--as-operator` **without** `--enforce` → may succeed (domain escape).
2. Same + `--enforce` → exit **1**, task not DONE.

### Block D — S04 status violations + config (18 named)

```bash
go test ./internal/config/... -count=1
go test ./cmd/trace -count=1 -run 'TestLoopStatusViolations|TestLoopStatusSchemaVersion|TestLoopStatusDeliberationFields|TestLoopStatusBlockedWhenBlockingUncertainty|TestTraceConfig|TestHelpIncludesTraceConfig'
go test ./internal/loop/... -count=1 -run 'TestLoopStatus|Gate'
```

Status (7 CLI):

| Test | Proves |
|------|--------|
| `TestLoopStatusIncludesViolationsWhenBlocked` | violations populated |
| `TestLoopStatusViolationsEmptyWhenClean` | empty when clean |
| `TestLoopStatusViolationsMatchGateEdit` | **parity with gate** |
| `TestLoopStatusViolationsAlwaysArray` | never null |
| `TestLoopStatusSchemaVersionUnchanged` | v1 unchanged |
| `TestLoopStatusDeliberationFields` | P20 fields preserved |
| `TestLoopStatusBlockedWhenBlockingUncertainty` | deliberation.blocked |

Config loader (4):

| Test | Proves |
|------|--------|
| `TestLoadEnforceModeMissingFile` | default off |
| `TestLoadEnforceModeValidValues` | off/warn/strict parse |
| `TestLoadEnforceModeMalformedJSON` | fail-closed off |
| `TestLoadEnforceModeUnknownValue` | fail-closed off |

Config CLI (7):

| Test | Proves |
|------|--------|
| `TestTraceConfigEnforceDefaultOff` | no stderr when off |
| `TestTraceConfigEnforceMalformedFailClosedOff` | malformed → off |
| `TestTraceConfigEnforceInvalidValueFailClosedOff` | invalid → off |
| `TestTraceConfigEnforceWarnSurfacesStderr` | warn hints exit 0 |
| `TestTraceConfigEnforceStrictSurfacesStderr` | strict = warn behavior |
| `TestTraceConfigEnforceOffNoStderrOnViolation` | off silent |
| `TestHelpIncludesTraceConfig` | help documents config |

No auto-enforce spot-check (manual or test):

- Config `strict` + plain `transition … DONE` without `--enforce` → behavior unchanged from S03.

### Block E — S05 harness install (16+ named)

```bash
go test ./internal/install/... -count=1 -run 'TestInstall|Enforcement|CursorHook|GitHook'
go test ./cmd/trace -count=1 -run 'TestInstall|TestHelpIncludesCursorHook|TestHelpIncludesLoop'
```

Enforcement tests (minimum):

| Test | Proves |
|------|--------|
| `TestInstallCursorIncludesLoopGateRule` | `.mdc` has gate + `TRACE_TASK_ID` |
| `TestInstallAgentsMDEnforcementBlock` | AGENTS.md markers |
| `TestInstallAgentsMDMarkersIdempotent` | no duplicate markers |
| `TestInstallClaudeIncludesLoopGateRule` | claude surface |
| `TestInstallCursorHookCallsGate` | hook script gate call |
| `TestInstallCursorHookPreToolUseMatcher` | hooks.json entry |
| `TestInstallDetectIncludesCursorHook` | detect lists target |
| `TestInstallGitHookUnchanged` | P22 fragment preserved |
| `TestInstallCursorMCPUnchanged` | MCP merge intact |
| `TestInstallEnforcementIdempotent` | idempotent writes |

Further rules content (grep in test fixtures):

- Rules mention status violations, `--enforce` DONE, export `--strict`.

### Block F — P19/P20 loop keepers (must stay green)

```bash
go test ./cmd/trace -count=1 -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext|TestLoopNext|TestLoopApply'
```

Combined gate regression (recommended once):

```bash
go test ./cmd/trace -count=1 -run 'TestLoopGate|TestLoopNext|TestLoopApply|TestLoopStatus|TestTransitionDoneEnforce|TestSeedExportStrict|TestInstall|TestHelpIncludesLoop'
```

### Block G — Compat ceiling 027

```bash
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
go test ./internal/store/... -count=1 -run TestMigrationStatusReportsEmbedMax
```

Confirm: **027_harness_agents.sql** present; **no 028+**.

## Required CLI evidence (archive)

Capture stdout/stderr from **ordinary CLI** (build `trace` from checkout). Minimum files:

| # | Artifact | Command sketch | Proves |
|---|----------|----------------|--------|
| 1 | `01-gate-blocked-edit.json` | init + blocking uncertainty → `trace loop gate --task <id> --for edit` | exit **1**, JSON envelope |
| 2 | `02-status-gate-parity.json` | same task → `trace loop status --task <id>` | violations match gate |
| 3 | `03-transition-done-enforce-block.txt` | debt task → `trace transition --task <id> --to DONE --enforce --as-operator` | exit **1** |
| 4 | `04-export-strict-enforce-no-write.txt` | debt → `trace seed export -o /tmp/out.json --strict --enforce` | exit **1**; file absent |
| 5 | `05-install-cursor-rules.txt` | temp dir + `.cursor/` → `trace install cursor --write` | rules + AGENTS block |
| 6 | `06-git-hook-unchanged.txt` | temp git repo → `trace install git-hook --write`; grep post-commit | `# begin-trace` preserved |
| 7 | `99-run-metadata.txt` | — | git SHA, go version, Blocks A–G summary |

Reference setup: reuse loop test helpers pattern (goal → task → plan → apply uncertainty).

## Do not

- Close [DR-HANDOFF.md](../../DR-HANDOFF.md) — S06-02 owns
- Mark criteria `done` without evidence files or test PASS
- Add migration 028+ or new MCP tools
- Auto-write `.trace/config.json` on install (verify absent)
- Scaffold Phase 24

## Exit criteria

- [ ] Evidence directory populated (CLI + test command log)
- [ ] Blocks A–G reported PASS (or row `failed` with reason)
- [ ] Must checklist mapped in board Notes with evidence pointers
- [ ] VERIFY-NOTES.md written (recommended)
- [ ] Residual risks listed (MCP wrapper, non-Cursor hooks, multi-violation lift)
- [ ] DR-HANDOFF remains **OPEN**

## Minimal todos

- [ ] Preflight: confirm all named tests exist (grep)
- [ ] Run Blocks A–G; capture in `99-run-metadata.txt`
- [ ] Script CLI smoke #1–#6; archive under evidence dir
- [ ] Write VERIFY-NOTES.md + board Notes with Must map
- [ ] Set row `done` or `failed` — **do not** close DR-HANDOFF
