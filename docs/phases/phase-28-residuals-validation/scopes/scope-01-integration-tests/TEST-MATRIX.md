# P25 regression matrix — Phase 28 S01

**Date:** 2026-08-20  
**Git SHA:** unavailable (workspace has no `.git`)  
**Row:** P28-S01-01

## Summary

| P25 | Theme | Automated | Dogfood-only |
|-----|-------|-----------|--------------|
| P25-A | Promotion (apply E2E, candidates, API, MCP nudge) | 4/4 rows (M-01..M-04) | Agent invoke of promotion path (FM-10) → S02 |
| P25-B | Saturation, discoveries-only, reset, STOP reason, CLI/domain reset | 6/6 rows (M-05..M-10) | — |
| P25-C | Install gap pass, parent orchestrator text, hook script | 3/3 rows (M-13..M-15) | Hook deny without `TRACE_TASK_ID` → S03; live parent behavior |
| P25-D | score.sh P25-3a/3b labels + `--arm` | 1/1 row (M-16 grep smoke) | Session-B richness (P25-3b) → S02; P25-4 attestation → S04 |
| P25-E | Thin export `--strict --enforce`; orphan honesty | 2/2 rows (M-11..M-12) | BLOCKING duplicate honesty msg → S04 |

R7 closed for **automated** coverage of P25-A/B/C/D/E. S00 audit “apply promotion E2E gap” is **pre-closed** by M-01 (`apply_test.go`); do not add `apply_promotion_test.go`.

## Matrix

| ID | P25 | INT | Test file | Test function | PASS criterion |
|----|-----|-----|-----------|---------------|----------------|
| M-01 | P25-A | INT-01 | `internal/loop/apply_test.go` | `TestLoopApplyPromotesBlockingDiscoveryViaSpawnedTask` | Apply with `spawned_tasks[].discovery_id` for a BLOCKING discovery: `NewSpawnedTasks=1`; task roster contains discovery ID; discovery→task `RelDiscoveryMentionsTask` persisted |
| M-02 | P25-A | INT-01 | `internal/loop/next_test.go` | `TestLoopNextPromotionCandidates` | `promotion_candidates[]` length 1 = unlinked BLOCKING discovery; linked BLOCKING and INFO excluded |
| M-03 | P25-A | INT-01 | `internal/domain/promote_test.go` | `TestPromoteBlockingDiscoveryAfterImport` (suite: 5 cases incl. create/link, idempotent, fail-closed, seed work_state) | After seed import, `PromoteBlockingDiscovery` creates/links task without breaking imported work_state |
| M-04 | P25-A | INT-06 | `internal/mcp/mcp_test.go` | `TestTraceAddDescriptionMentionsPromotionPath` | `trace_add` description lists discovery first and contains `BLOCKING discovery`, `trace_add kind=task`, `spawned_tasks`, `discovery_id` |
| M-05 | P25-B | INT-02 | `internal/loop/saturation_reset_test.go` | `TestApplyConsecutiveEmptySaturationThreshold` | 1 consecutive empty apply does not STOP; 2 empty applies STOP (`SaturationEmptyThreshold=2` in `internal/deliberation/types.go`) |
| M-06 | P25-B | INT-02 | `internal/loop/saturation_reset_test.go` | `TestApplyDiscoveriesOnlyDoesNotIncrementEmptyCounter` | Discoveries-only apply does not increment consecutive-empty counter / does not trigger saturation STOP |
| M-07 | P25-B | INT-05 | `internal/loop/saturation_reset_test.go` | `TestResetClearsSaturationAndPreventsImmediateReStop` | After saturation STOP, reset then apply does not immediately re-STOP |
| M-08 | P25-B | INT-09 | `internal/loop/saturation_reset_test.go` | `TestExportStopReasonMatchesGateAfterSaturation` | Persisted `stop_reason` matches gate/export/status after saturation STOP |
| M-09 | P25-B | INT-05 | `cmd/trace/loop_test.go` | `TestLoopResetCLIClearsStop` | `trace loop reset` clears STOP so subsequent gate/status is not blocked by prior saturation |
| M-10 | P25-B | INT-05 | `internal/domain/deliberation_test.go` | `TestResetDeliberationStateClearsStopPreservesCritique` | `ResetDeliberationState` clears STOP fields and preserves critique |
| M-11 | P25-E | INT-07 | `cmd/trace/enforce_test.go` | `TestSeedExportStrictEnforceBlocksP26ThinGraph`; `TestSeedExportStrictCleanAllowsWrite` | `--strict --enforce` on thin P26-style graph: no write + fail; clean/rich graph: write allowed |
| M-12 | P25-E | INT-07 | `internal/domain/seed_export_honesty_test.go` | `TestSeedDocumentHonestyOrphanDiscovery`; `TestSeedDocumentHonestyOrphanDecision` | Orphan discovery and orphan decision each produce honesty violations; linked clean graph does not (`TestSeedDocumentHonestyCleanWithLinkedDecision`) |
| M-13 | P25-C | INT-03 | `internal/install/enforcement_test.go` | `TestGapPassPromptNonEmpty`; `TestInstallCursorIncludesLoopGateRule`; `TestInstallClaudeIncludesLoopGateRule` | `GapPassPrompt` non-empty and mentions gap + `TRACE_TASK_ID`; Cursor and Claude install bodies include mandatory gap pass |
| M-14 | P25-C | INT-04 | `internal/install/enforcement_test.go` | `TestInstallCursorIncludesLoopGateRule`; `TestInstallAgentsMDEnforcementBlock` | Cursor rules contain `Parent orchestrator`. **Documented asymmetry (no product change):** `AgentsEnforcementBlock` / AGENTS.md install includes gap pass + loop gate, **not** `ParentOrchestratorRule` (`enforcement.go`: Cursor/Claude get GapPass+Parent; AGENTS block is GapPass only) |
| M-15 | P25-C | INT-04 | `internal/install/enforcement_test.go` | `TestInstallCursorHookCallsGate` | Installed `trace-loop-gate.sh` contains `loop gate` and `--for edit`. Deny-when-strict / empty `TRACE_TASK_ID` **not** asserted here (S03) |
| M-16 | P25-D | INT-08, INT-10 | `evals/p28-regression/score_arm_labels_test.sh` | bash grep smoke (no harness) | `experiments/ab-p25-gap-pass-validation/score.sh` contains `P25-3a`, `P25-3b`, and `--arm`. Does **not** execute `./score.sh` or `./prepare.sh` |

## Install text asymmetry (P25-C / INT-04)

| Surface | Gap pass | Parent orchestrator |
|---------|----------|---------------------|
| Cursor `.cursor/rules/trace-enforcement.mdc` | yes | yes (`TestInstallCursorIncludesLoopGateRule`) |
| Claude install body | yes | yes (`TestInstallClaudeIncludesLoopGateRule`) |
| AGENTS.md `AgentsEnforcementBlock` | yes | **no** — by current product; S01 documents only |

## INT-01..11 mapping (S01 automated vs later scopes)

| INT | Theme | S01 rows | Residual owner |
|-----|-------|----------|----------------|
| INT-01 | Promotion | M-01, M-02, M-03 | Agent invoke / dogfood → S02 |
| INT-02 | Saturation threshold | M-05, M-06 | — |
| INT-03 | Gap-pass install | M-13 | Behavior in live session → S02 |
| INT-04 | Parent orchestrator + hook | M-14, M-15 (text/script install) | Deny without task + drift golden → S03 |
| INT-05 | Deliberation reset | M-07, M-09, M-10 | — |
| INT-06 | MCP trace_add nudge | M-04 | Agent compliance → S02 |
| INT-07 | Seed export honesty | M-11, M-12 | BLOCKING dup msg → S04 |
| INT-08 | score.sh protocol / `--arm` | M-16 | Live harness → S05 |
| INT-09 | STOP reason parity | M-08 | — |
| INT-10 | P25-3a/3b labels | M-16 (labels present) | Directed richness → S02 |
| INT-11 | Hook drift | — | S03 (`hook_drift_test.go` seed) |

## Explicit deferrals (not S01)

| Item | Owner | Rationale |
|------|-------|-----------|
| Hook deny without TRACE_TASK_ID | S03 | FM-05 / R2/R3 — install script still allows empty task id |
| BLOCKING dup honesty msg | S04 | R4 — overlapping `seed.go` + `seed_export_honesty.go` messages |
| Session-B P25-3b richness | S02 | R1 dogfood — M-16 only greps labels |
| P25-4 attestation | S04 | R5 — `score.sh` still `skip` |
| AGENTS vs Cursor orchestrator text | document only | M-14; no product change in S01 |
| `apply_promotion_test.go` | n/a | Pre-closed by M-01 in `apply_test.go` |

## Preflight / regression record (P28-S01-01)

- Path tests: `00-PLANNER.md`, `RESIDUAL-AUDIT.md` present.
- Anchors: `apply_test.go:325` promotion test; `SaturationEmptyThreshold=2`; `score.sh` P25-3a/3b at L197–215.
- `go test ./internal/... -count=1` PASS (module fetch via `GOPROXY=direct` when sandbox proxy 403s `segmentio/encoding`).
- `go test ./cmd/trace/... -count=1` PASS.
- Targeted `-run` slices: see implementer Notes on board row P28-S01-01.
- `bash evals/p28-regression/score_arm_labels_test.sh` PASS.
