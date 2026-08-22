# Scope 03 todos — VERIFY + handoff

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P27-S03-00 | planner | done | 2026-08-20: Locked verify floor + enforce policy (both arms `--strict --enforce`); successor table (`no successor` default / Phase 28 alt). Thickened 01-verify + 02-dr-handoff. |
| P27-S03-01 | verify | pending | VERIFY floor + own `score.sh` T02 `--strict --enforce` upgrade; evidence under `experiments/runs/YYYY-MM-DD-p27-s03-01-verify/`. |
| P27-S03-02 | reviewer | pending | Independent re-verify; DR-HANDOFF close; successor never TBD. |

## Locked defaults (S03-00)

| Item | Value |
|------|-------|
| Verify floor | Rebuild; `go test ./internal/...`; cmd honesty tests; product thin enforce demo; score `--arm build` + `--arm directed` |
| Harness enforce | S03-01: `score.sh` T02 → `--strict --enforce` **both** arms; thin → G2 FAIL (WARN→FAIL) |
| Build P25-3a | FAIL on thin = expected residual |
| Directed P25-3b | Label required; PASS only with Session-B rich graph (optional dogfood) |
| Evidence | `experiments/runs/YYYY-MM-DD-p27-s03-01-verify/evidence/` |
| Default successor | `no successor` |
| Alt successor | Phase 28 if regression needs new phase / human promotes |
| Human-gated | P25-4 + optional Session-B — non-blocking |
| Residual | S02-02 BLOCKING duplicate orphan msg — note only |
