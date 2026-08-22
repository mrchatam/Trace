# P28-S01-02 — Integration test review

## Metadata
- id: P28-S01-02
- todo_ids: [P28-S01-02]
- role: reviewer
- skills: [code-review-and-quality, qa-lead]
- mcps: [user-codegraph]
- verification: automated
- hooks: []

## Objective

Independent review of S01 deliverables — confirm R7 closure via `TEST-MATRIX.md` and green regression commands before Session-B (S02). **Fresh subagent** — do not reuse S01-01 session.

## References

- [01-implement.md](01-implement.md) — exit criteria + test commands
- [00-PLANNER.md](00-PLANNER.md) — locked defaults
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — R7 seeds
- [TEST-MATRIX.md](TEST-MATRIX.md) (produced by S01-01)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Verify checklist

### Deliverable presence

- [ ] `TEST-MATRIX.md` exists in this scope folder
- [ ] Summary table covers P25-A, P25-B, P25-C, P25-D, P25-E
- [ ] Each matrix row has: test file, function name, PASS criterion
- [ ] Deferrals table lists S02/S03/S04 items (hook deny, dogfood, honesty dup, attestation)

### P25 theme coverage (automated)

- [ ] **P25-A Promotion:** cites `TestLoopApplyPromotesBlockingDiscoveryViaSpawnedTask` in `apply_test.go` (not a duplicate new file)
- [ ] **P25-A** also references `promote_test.go`, `next_test.go` promotion_candidates, and/or `mcp_test.go` promotion description
- [ ] **P25-B Saturation:** cites `TestApplyConsecutiveEmptySaturationThreshold` (1 empty ≠ STOP; 2 empty = STOP)
- [ ] **P25-B Reset:** cites reset sequence in `saturation_reset_test.go` and/or `TestLoopResetCLIClearsStop`
- [ ] **P25-E Honesty:** cites `TestSeedExportStrictEnforceBlocksP26ThinGraph` + rich export pass case
- [ ] **P25-C Install:** cites gap pass tests in `enforcement_test.go`; documents AGENTS vs cursor orchestrator asymmetry if applicable
- [ ] **P25-D Protocol:** automated row present **or** explicitly marked optional with `evals/p28-regression/` smoke added

### Live spot-checks (reviewer runs)

```bash
cd /home/ali/Desktop/Trace

# Promotion E2E exists (audit gap pre-closed)
grep -n 'func TestLoopApplyPromotesBlockingDiscoveryViaSpawnedTask' internal/loop/apply_test.go

# Saturation threshold anchor (no stale saturation.go)
grep -n 'SaturationEmptyThreshold' internal/deliberation/types.go

# score.sh arm labels (P25-D)
grep -n 'P25-3a\|P25-3b' experiments/ab-p25-gap-pass-validation/score.sh

# Full regression (must PASS)
go test ./internal/... -count=1
go test ./cmd/trace/... -run 'SeedExport|Enforce|ThinGraph|LoopReset' -count=1
go test ./internal/install/... -count=1
```

### Process

- [ ] No product behavior change beyond tests/fixtures/docs
- [ ] No Session-B dogfood claimed in S01
- [ ] No hook deny implementation (deferred S03)
- [ ] `go test ./internal/...` PASS

## Findings severity

| Level | Action |
|-------|--------|
| blocker | Missing matrix for any P25 theme; `go test ./internal/...` FAIL; product code changed |
| high | Matrix cites nonexistent tests; duplicate `apply_promotion_test.go` added unnecessarily |
| medium | P25-D row missing with no deferral note; matrix INT mapping incomplete |
| low / nit | Wording/formatting only |

## Spawn policy

- **HIGH/blocker** → insert `P28-S01-02a` (implement fix) + `P28-S01-02b` (review) below this row
- **No HIGH/blocker** → APPROVE → next runnable **P28-S02-00** (Session-B may parallelize if S01 green)

## Exit criteria

- [ ] Verdict: APPROVE or spawn with pending follow-up
- [ ] Confidence: **high** (or **medium** with explicit residual risks)
- [ ] Board row P28-S01-02 status + notes updated

## Todo updates

Status + notes on **P28-S01-02** only.

## Next

**P28-S02-00** (unless spawn pending)
