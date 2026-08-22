# S04 scope todos — product polish

| ID | Status | Prompt | Notes |
|----|--------|--------|-------|
| P28-S04-00 | done | [00-PLANNER.md](00-PLANNER.md) | R4/R5 locked 2026-08-20; R6 deferred |
| P28-S04-01 | pending | [01-implement.md](01-implement.md) | Remove store BLOCKING loop + P25-4 env attestation |
| P28-S04-02 | pending | [02-review.md](02-review.md) | Independent review |

## Locked targets (from P28-S04-00)

| Target | Lock |
|--------|------|
| R4 | Remove store-backed BLOCKING orphan loop in `cmd/trace/seed.go` `collectExportGraphHonestyViolations` |
| R4 sole source | `domain.CollectSeedDocumentHonestyViolations` |
| R4 forbidden | Add `Severity` to `SeedEntity`; keep dual message paths |
| R4 test | BLOCKING orphan → exactly one honesty violation for that ID |
| R5 build | `P25_ATTEST_BUILD=Y` + `--arm build` → pass P25-4 |
| R5 directed | `P25_ATTEST_DIRECTED=Y` + `--arm directed` → pass P25-4 |
| R5 unset | skip (backward compatible) |
| R5 wrong-arm | ignored |
| R5 docs | PROTOCOL.md required; RUBRIC.md optional one-liner |
| R6 | **Deferred** — thin-but-strict stderr hint not in S04 |
| Residuals | Close R4, R5; leave R6 open/partial for S05 |
| Out | Daemon/HTTP; RESULTS.md parser; reopen R1–R3/R8 |

## Implementer acceptance (S04-01)

- [ ] Store BLOCKING loop removed
- [ ] Single orphan message source
- [ ] Regression: BLOCKING orphan → one violation
- [ ] score.sh env attestation wired
- [ ] PROTOCOL.md updated
- [ ] `go test ./internal/... ./cmd/trace/...` PASS
- [ ] R6 not implemented
