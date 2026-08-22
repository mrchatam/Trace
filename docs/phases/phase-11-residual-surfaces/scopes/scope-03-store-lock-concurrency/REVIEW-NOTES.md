# P11-S03-02 — REVIEW-NOTES (Store lock / concurrency)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | **DF-47** — short bounded retry: release-during-wait → second Open succeeds | **Pass** — `defaultLockWait=350ms`, `lockRetryStep=25ms` in `acquireTraceLock`; `TestOpenRetrySucceedsWhenLockReleasedSoon` (~50ms Close mid-wait) |
| 2 | **DF-47** — held lock past budget still `ErrLocked`; exclusivity / compat `trace_lock_ok` | **Pass** — `TestConcurrentStoreOpenFailClosed` (~0.36s under budget); compat checklist `lock=true` / `trace_lock_ok`; `TestInitFailClosedWhenStoreLocked` exit **2** |
| 3 | **DF-47** — ErrLocked and/or help/MCP guide serialize CLI↔MCP or worktrees; exit **2** | **Pass** — `ErrLocked` text + `TestErrLockedSerializeGuidance`; help Global + `TestHelpSerializeLockGuidance`; MCP wraps `%w` (no policy fork); `exitFail=2` unchanged |
| 4 | Exclusive flock **not** dropped; no multi-writer / daemon / long-lived MCP store | **Pass** — still `LOCK_EX\|LOCK_NB` for Open→Close; MCP per-tool `openStore`+Close; no daemon/HTTP redesign |
| 5 | G19 — no lock-policy fork in adapters | **Pass** — retry/`ErrLocked` only in `internal/store/lock.go`; CLI/MCP consume sentinel |
| 6 | No forbidden architecture | **Pass** — no new mig/`011_*`; no multi-writer SQLite; no indefinite wait (bounded + optional `TRACE_LOCK_WAIT_MS`) |
| 7 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + P11-S01/S02 + Gate C `dry_run:false` | **Pass** — locked CGO0/CGO1 suites green; Gate C artifacts still `dry_run:false`; P11-S01/S02 untouched |
| 8 | Board Notes accurate; planner row had no product Go | **Pass** — P11-S03-00 Notes claim no product Go; P11-S03-01 Notes match live APIs/tests (350ms + env override + named tests) |

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | `open.go` comment still says second Open fails until Close (omits brief-retry nuance) | Residual OK — behavior matches locks |
| low | `go test ./...` still fails setup on `similar projects/graphify` path space | Pre-existing non-product; product pkgs PASS |

## Residuals (explicit)

1. Sustained parallel N writers on one root still fail-closed after budget (planner residual OK).
2. Full-module `./...` FAIL only on pre-existing `similar projects/graphify` space path.
3. Comment polish on `Open` doc is optional forward nit — not spawn-worthy.

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/store/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1  → PASS (product); similar projects/graphify setup fail pre-existing
Named: TestOpenRetrySucceedsWhenLockReleasedSoon / TestConcurrentStoreOpenFailClosed / TestErrLockedSerializeGuidance / TestHelpSerializeLockGuidance / TestInitFailClosedWhenStoreLocked / TraceLockOK → PASS
```

## Next

**P11-S04-00**
