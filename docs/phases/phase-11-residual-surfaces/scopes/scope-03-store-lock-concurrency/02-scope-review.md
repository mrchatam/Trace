# P11 / S03 / 02 — Scope review (Store lock / concurrency)

## Metadata
- id: P11-S03-02
- todo_ids: [P11-S03-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S03 (**DF-47**). Fresh subagent. Compare claims + FINAL locks to live code/tests. Spawn `02a`/`02b` for blocker/high. Do not rewrite prior `done` history. Small inline fixes OK for medium/low when cheaper than spawn.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-store-lock-concurrency.md](01-store-lock-concurrency.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-47
- Live: `internal/store/{lock,open}.go`; `cmd/trace/help.go`; `internal/mcp/project.go`; `evals/compat`
- P08 S02 exclusivity + path-local bind must remain green

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (reviewer).

## Checklist (must all pass for APPROVE)

| # | Check |
|---|--------|
| 1 | **DF-47** — short bounded retry: release-during-wait → second Open succeeds |
| 2 | **DF-47** — held lock past budget still `ErrLocked`; exclusivity / compat `trace_lock_ok` green |
| 3 | **DF-47** — ErrLocked and/or help/MCP guide serialize CLI↔MCP or worktrees; exit **2** unchanged |
| 4 | Exclusive flock **not** dropped; no multi-writer / daemon / long-lived MCP store redesign |
| 5 | G19 — no lock-policy fork in adapters |
| 6 | No forbidden architecture (daemon/HTTP/full-rebuild/mig) |
| 7 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + P11-S01/S02 + Gate C `dry_run:false` |
| 8 | Board Notes accurate; planner row had no product Go |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-checks: `TestOpenRetrySucceedsWhenLockReleasedSoon`; `TestConcurrentStoreOpenFailClosed`; ErrLocked/help serialize wording; `TestInitFailClosedWhenStoreLocked`.

## Exit criteria
- [x] Checklist evidenced; confidence high (or medium with residuals listed)
- [x] Board status + Notes; next **P11-S04-00** (unless spawn)

## Todo updates
Reviewer: status + Notes; spawn `02a`/`02b` only for open blocker/high. Forward-only.
