# S01 scope todos

Integration test matrix for P25-A/B/C/D/E (R7). Dogfood is **S02**.

| Order | ID | Status | Prompt | Notes |
|------:|----|--------|--------|-------|
| 1 | P28-S01-00 | done | [00-PLANNER.md](00-PLANNER.md) | Planner locked implement/review prompts + live survey |
| 2 | P28-S01-01 | pending | [01-implement.md](01-implement.md) | Deliver `TEST-MATRIX.md`; optional `evals/p28-regression/` P25-D smoke |
| 3 | P28-S01-02 | pending | [02-review.md](02-review.md) | Verify matrix + `go test ./internal/...` green |

## Locked test targets (S01-01)

| Gap (audit) | Resolution |
|-------------|------------|
| Apply promotion E2E | **Pre-closed** — extend `internal/loop/apply_test.go` only if strengthening assertions; do not add `apply_promotion_test.go` |
| Saturation / reset | Document existing `saturation_reset_test.go` rows |
| Honesty enforce | Document `cmd/trace/enforce_test.go` + `seed_export_honesty_test.go` |
| Install gap pass | Document `enforcement_test.go`; note AGENTS block gap-pass-only vs cursor orchestrator |
| score.sh P25-3a/3b | Optional `evals/p28-regression/score_arm_labels_test.sh` |
| Hook deny / drift | Defer **S03** |
| BLOCKING dup | Defer **S04** (add matrix note only) |

## Serial order

S01-00 → S01-01 → S01-02. S02 may start after S01-01 if matrix draft exists (default: wait for S01-02 APPROVE).
