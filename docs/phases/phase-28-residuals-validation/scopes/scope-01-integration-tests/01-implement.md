# P28-S01-01 — Integration test implementer

## Metadata
- id: P28-S01-01
- todo_ids: [P28-S01-01]
- role: implementer
- skills: [test-driven-development, incremental-implementation]
- mcps: [user-codegraph]
- verification: automated
- hooks: []

## Objective

Close **R7** (consolidated regression matrix) for P25-A/B/C/D/E with automated tests — **not** Session-B dogfood (S02). Primary deliverable: `TEST-MATRIX.md` mapping each P25 theme → existing or new test → INT → PASS criterion. Extend tests only where the live survey below shows a gap.

## References

- [00-PLANNER.md](00-PLANNER.md) — locked defaults
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — § Test coverage matrix (R7 seeds)
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [Phase 28 README](../../README.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (live survey 2026-08-20)

Planner verified repo — **do not recreate** tests that already exist.

| P25 | Theme | Test home | Status | Anchor test(s) |
|-----|-------|-----------|--------|----------------|
| P25-A | Promotion | `internal/loop/apply_test.go` | **covered** | `TestLoopApplyPromotesBlockingDiscoveryViaSpawnedTask` (L325) — BLOCKING discovery → `spawned_tasks[].discovery_id` → task row + link |
| P25-A | Promotion candidates | `internal/loop/next_test.go` | **covered** | `TestLoopNextPromotionCandidates` |
| P25-A | Promote API | `internal/domain/promote_test.go` | **covered** | 5 cases incl. `TestPromoteBlockingDiscoveryAfterImport` |
| P25-A | MCP nudge | `internal/mcp/mcp_test.go` | **covered** | `TestTraceAddDescriptionMentionsPromotionPath` |
| P25-B | Saturation threshold | `internal/loop/saturation_reset_test.go` | **covered** | `TestApplyConsecutiveEmptySaturationThreshold` — 1 empty no STOP; 2 empty STOP (`SaturationEmptyThreshold=2` in `internal/deliberation/types.go`) |
| P25-B | Discoveries-only | `internal/loop/saturation_reset_test.go` | **covered** | `TestApplyDiscoveriesOnlyDoesNotIncrementEmptyCounter` |
| P25-B | Reset + re-apply | `internal/loop/saturation_reset_test.go` | **covered** | `TestResetClearsSaturationAndPreventsImmediateReStop` |
| P25-B | STOP reason parity | `internal/loop/saturation_reset_test.go` | **covered** | `TestExportStopReasonMatchesGateAfterSaturation` |
| P25-B | Reset CLI | `cmd/trace/loop_test.go` | **covered** | `TestLoopResetCLIClearsStop` |
| P25-B | Domain reset | `internal/domain/deliberation_test.go` | **covered** | `TestResetDeliberationStateClearsStopPreservesCritique` |
| P25-E | Thin export enforce | `cmd/trace/enforce_test.go` | **covered** | `TestSeedExportStrictEnforceBlocksP26ThinGraph` — blocks write; rich `TestSeedExportStrictCleanAllowsWrite` |
| P25-E | Orphan honesty | `internal/domain/seed_export_honesty_test.go` | **covered** | orphan discovery/decision cases |
| P25-C | Install gap pass | `internal/install/enforcement_test.go` | **covered** | `TestGapPassPromptNonEmpty`; `TestInstallCursorIncludesLoopGateRule`; `TestInstallClaudeIncludesLoopGateRule` |
| P25-C | Parent orchestrator text | `internal/install/enforcement_test.go` | **partial** | Cursor rules assert `Parent orchestrator` (L88–89); AGENTS block has gap pass only (`AgentsEnforcementBlock` — document in matrix, **no product change**) |
| P25-C | Hook script content | `internal/install/enforcement_test.go` | **covered** | `TestInstallCursorHookCallsGate` — defer deny-when-strict to **S03** |
| P25-D | score.sh P25-3a/3b labels | `evals/p28-regression/` (optional) | **gap** | No automated test; optional bash smoke grepping `experiments/ab-p25-gap-pass-validation/score.sh` L197–218 |

**Path correction (from S00):** No `internal/loop/saturation.go` — policy in `internal/deliberation/types.go`; loop tests in `saturation_reset_test.go`.

## Preflight

Run from repo root before editing:

```bash
cd /home/ali/Desktop/Trace

# Scope artifacts
test -f docs/phases/phase-28-residuals-validation/scopes/scope-01-integration-tests/00-PLANNER.md
test -f docs/phases/phase-28-residuals-validation/scopes/scope-00-residual-audit/RESIDUAL-AUDIT.md

# Product anchors (spot-check audit claims still true)
grep -n 'TestLoopApplyPromotesBlockingDiscoveryViaSpawnedTask' internal/loop/apply_test.go
grep -n 'SaturationEmptyThreshold' internal/deliberation/types.go
grep -n 'P25-3a\|P25-3b' experiments/ab-p25-gap-pass-validation/score.sh

# Baseline green (record output in PR notes if any fail)
go test ./internal/... -count=1
go test ./cmd/trace/... -count=1
```

## Minimal todos

1. **Preflight** — run bash block above; abort if anchors missing.
2. **TEST-MATRIX.md** — create in this scope folder using template below; every P25-A/B/C/D/E row must cite ≥1 test file + function name + PASS criterion.
3. **Optional P25-D smoke** — if matrix row would otherwise be `dogfood-only`, add `evals/p28-regression/score_arm_labels_test.sh` (or Go test wrapper) that greps `score.sh` for `P25-3a` / `P25-3b` arm branches and `--arm` flag. No harness execution required.
4. **Regression run** — after any new test file, re-run commands in § Test commands.
5. **Board** — update **P28-S01-01** status + notes only.

## TEST-MATRIX.md template

Write to `scopes/scope-01-integration-tests/TEST-MATRIX.md`:

```markdown
# P25 regression matrix — Phase 28 S01

**Date:** YYYY-MM-DD  
**Git SHA:** (if available)  
**Row:** P28-S01-01

## Summary

| P25 | Theme | Automated | Dogfood-only |
|-----|-------|-----------|--------------|
| P25-A | Promotion | N/M rows | … |
| … | … | … | … |

## Matrix

| ID | P25 | INT | Test file | Test function | PASS criterion |
|----|-----|-----|-----------|---------------|----------------|
| M-01 | P25-A | INT-01 | internal/loop/apply_test.go | TestLoopApplyPromotesBlockingDiscoveryViaSpawnedTask | Apply with discovery_id spawns task; discovery→task link persisted |
| … | … | … | … | … | … |

## Explicit deferrals (not S01)

| Item | Owner | Rationale |
|------|-------|-----------|
| Hook deny without TRACE_TASK_ID | S03 | FM-05 / R2/R3 |
| BLOCKING dup honesty msg | S04 | R4 |
| Session-B P25-3b richness | S02 | R1 dogfood |
| P25-4 attestation | S04 | R5 |
```

Populate all rows from § Locked defaults; add optional P25-D row if smoke added.

## Test commands

Implementer and reviewer both run:

```bash
cd /home/ali/Desktop/Trace

# Full internal suite (must stay green)
go test ./internal/... -count=1

# Targeted P25 regression slice
go test ./internal/loop/... -run 'Promote|Saturation|Reset|DiscoveriesOnly' -count=1
go test ./internal/domain/... -run 'Promote|Honesty|ResetDeliberation' -count=1
go test ./internal/install/... -run 'GapPass|Orchestrator|Hook' -count=1
go test ./internal/mcp/... -run 'Promotion' -count=1
go test ./cmd/trace/... -run 'SeedExport|Enforce|ThinGraph|LoopReset' -count=1

# Optional evals (if added)
test -f evals/p28-regression/score_arm_labels_test.sh && bash evals/p28-regression/score_arm_labels_test.sh
```

## Exit criteria

- [ ] `TEST-MATRIX.md` complete — every P25-A/B/C/D/E theme has ≥1 automated row (P25-D may use optional smoke)
- [ ] Matrix maps to INT-01..11 where applicable; dogfood-only rows explicitly deferred to S02/S03/S04
- [ ] `go test ./internal/... -count=1` **PASS**
- [ ] No product behavior change (tests/fixtures/docs only)
- [ ] No daemon / HTTP
- [ ] Board row P28-S01-01 Notes cite matrix path + any new test files

## Do not

- Add `internal/loop/apply_promotion_test.go` — promotion E2E already in `apply_test.go`
- Run Session-B dogfood or `./prepare.sh G1` (S02)
- Change hook deny semantics (S03)
- Fix BLOCKING duplicate honesty messages (S04)
- Rewrite S00 `done` rows or `RESIDUAL-AUDIT.md`

## Next

**P28-S01-02**
