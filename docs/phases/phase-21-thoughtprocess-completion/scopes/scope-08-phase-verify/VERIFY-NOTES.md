# P21-S08-01 VERIFY-NOTES

**Date:** 2026-08-18  
**Evidence:** `experiments/runs/2026-08-18-p21-s08-01-verify/evidence/`  
**Outcome:** PASS (named floor Blocks A–B + compat; optional product bar skipped after sandbox network flake)

## Test floor

| Block | Result | Notes |
|-------|--------|-------|
| A — P20 S07-01 floor | **PASS** | All 10 sub-blocks green; `./internal/loop/...` now has 14 S06 tests (no vacuous pass) |
| B — P21 S01–S07 deltas | **PASS** | 12 sub-blocks; `b-s04-promotion` required `-parallel 1` + store fix for `UpsertOutcomeResult` `created_at` on conflict |
| C — compat ceiling **21** | **PASS** | `TestCompatibilitySecurityChecklist` |
| C — product bar (optional) | **SKIP/FAIL** | Sandbox `go mod download` EOF for MCP deps; not required for close |

## §31 live mini-eval (D-14)

- **Workspace:** live repo; `trace-bin` built from checkout; temp DB under `project/`
- **Distinct phases (5):** CRITIQUE → INVESTIGATE → VERIFY → ORIENT → STOP (`02-full-cycle-next.json`)
- **Blocking uncertainty:** step-02 INVESTIGATE (not EXECUTE); `02-apply-blocking-uncertainty.json`
- **Seed P20 keys:** CLI export grep in `01-seed-export-p20-grep.txt`; import round-trip via `trace seed import` after CLI whitelist fix
- **promotion_blocked:** `04-baseline-promote.json` — `eval_regression` after baseline regression helper
- **risk_hints / historical_relationships:** `06-risk-hints.json`, `07-historical-relationships.json`
- **trace why:** `03-why-uncertainty.json` exit 0 + steps
- **Tx rollback:** `05-transactional-apply-fail.json` (unit test log)

**Residual:** `BuildPolicyInputs` still stubs cycle flags (`execute_pending`, etc.) — EXECUTE/TEST/EVALUATE/EXPLORE phases proven by Block B unit tests, not live CLI. Long-horizon multi-agent §31 still depends on seed clone in the wild.

## Must coverage (COVERAGE.md anchors)

| § | Evidence |
|---|----------|
| 1 | Block A S06 + Block B cycle + CLI `02-full-cycle-next.json` |
| 5 | `TestSelectNextFullCycleOrdering`; 14-row table |
| 13 | `TestPromoteBaselineSupersedesPrior`; CLI `04-baseline-promote.json` |
| 16 | `TestCreateExperimentLinksOutcome`; not bake-off |
| 17 | CLI `07-historical-relationships.json` |
| 18 | CLI `06-risk-hints.json`; `TestRiskHintsManyPaths` |
| 25 | CLI `03-why-uncertainty.json`; `TestCLIWhyUncertainty` |
| 29O | gates + `TestLoopApplyTransactionalRollbackOnFailure` |
| 29Q | named floor Blocks A–C |
| 31 | live CLI ≥5 phases |

## Fixes during verify (minimal)

1. `cmd/trace/seed.go` — add P20 keys to `seedImportAllowedKeys` (CLI import matched domain S01)
2. `internal/store/outcomes.go` — `UpsertOutcomeResult` updates `created_at` on conflict (fixes flaky `TestEvalRegressionGateClearsAfterResolve`)

## Residual risks (unchanged)

- Experiments not in seed export
- D-16: Trace does not run tests autonomously
- D-19: hosted MCP out of scope
- DR-HANDOFF remains **OPEN** (S08-02)
