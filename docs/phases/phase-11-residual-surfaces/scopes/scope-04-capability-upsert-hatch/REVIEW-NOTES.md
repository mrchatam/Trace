# P11-S04-02 — REVIEW-NOTES (Capability upsert + hatch vs caps)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | **DF-41** — empty-ID re-declare same slug → same id + updated fields; explicit different-id slug clash still fails | **Pass** — `UpsertCapability` empty-ID → `GetCapabilityBySlug` reuse id; `TestUpsertCapabilityBySlugUpdatesExisting`; clash in `TestUpsertCapabilityGetAndReject` |
| 2 | **DF-51** — `AllowDoneWithoutReview` does **not** bypass missing-caps; caps→DONE order; WARNING/docs mention override | **Pass** — caps gate before →DONE hatch in `task_state.go`; `TestAllowDoneDoesNotBypassMissingCaps`; WARNING phrases in CLI/MCP; `TestAllowDoneWarnsOnStderr` / `TestTransitionAllowDoneEmitsWarning` |
| 3 | Gate G hatch + DF-24 fail-closed retained; no mig; no hatch→auto-missing-caps | **Pass** — hatch still review/operator only; flags independent; no `011_*` / no new mig; honesty/Gate G still green via carry-forward |
| 4 | G19 — no domain fork in CLI/MCP adapters | **Pass** — CLI/MCP declare call `svc.UpsertCapability`; transition flags map to domain opts only; WARNING strings copy-only |
| 5 | No forbidden architecture | **Pass** — no daemon/HTTP/full-rebuild; no multi-writer; no collapsing hatch into missing-caps |
| 6 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + P11-S01–S03 + Gate C `dry_run:false` | **Pass** — locked CGO0/CGO1 suites green; Gate C metrics still `dry_run:false`; prior S01–S03 untouched |
| 7 | Board Notes accurate; planner row had no product Go | **Pass** — P11-S04-00 Notes claim no product Go; P11-S04-01 Notes match live APIs/tests |

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | Full-module `go test ./...` may still fail setup under `similar projects/graphify` space path | Pre-existing non-product; product pkgs PASS |

## Residuals (explicit)

1. Product packages (`./cmd/... ./internal/... ./evals/...`) PASS; research `similar projects/` space-path setup fail remains residual OK.
2. Optional help/MCP schema phrase asserts not required beyond WARNING tests — live help/schema already mention independence.

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1  → PASS (product)
Named: TestUpsertCapabilityBySlugUpdatesExisting / TestUpsertCapabilityGetAndReject (clash) / TestAllowDoneDoesNotBypassMissingCaps / TestAllowDoneWarnsOnStderr / TestTransitionAllowDoneEmitsWarning → PASS
Gate C: docs/verification/gate-c-x0/metrics-{b0,g1}.json dry_run:false intact
```

## Next

**P11-S05-00**
