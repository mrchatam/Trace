# P11 / S02 / 02 — Scope review (Review PASS+FAIL / operator identity)

## Metadata
- id: P11-S02-02
- todo_ids: [P11-S02-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S02 (**DF-43, DF-44**). Fresh subagent. Compare claims + FINAL locks to live code/tests. Spawn `02a`/`02b` for blocker/high. Do not rewrite prior `done` history. Small inline fixes OK for medium/low when cheaper than spawn.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-review-pass-fail-operator.md](01-review-pass-fail-operator.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-43, DF-44
- Live: `internal/domain/task_state.go`; `cmd/trace/{transition,help}.go`; `internal/mcp/{server,tools_write}.go`; `evals/honesty/honesty_test.go`
- P10 S04 DF-17/18/24/26 baseline must remain green

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (reviewer).

## Checklist (must all pass for APPROVE)

| # | Check |
|---|--------|
| 1 | **DF-43** — linked FAIL+PASS + `AllowOperatorDone` rejects →DONE; reason mentions FAIL |
| 2 | **DF-43** — PASS alone / PASS+UNCERTAIN still authorize with flag; hatch bypasses FAIL |
| 3 | **DF-43** — honesty Path C supersedes FAIL before DONE; A/B + Gate G green |
| 4 | **DF-44** — freestanding flag retained; help/MCP state flag≠identity / conscious claim; Actor≠auth; no OAuth |
| 5 | G19 — no domain fork in CLI/MCP adapters |
| 6 | No forbidden architecture (daemon/HTTP/full-rebuild/mig) |
| 7 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + P11-S01 DF-40 + Gate C `dry_run:false` |
| 8 | Board Notes accurate; planner row had no product Go |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-checks: `TestSiblingFailBlocksDone` (or equiv); honesty Path C; `TestOperatorDoneRequiresFlag`; help/MCP identity wording.

## Exit criteria
- [x] Checklist evidenced; confidence high (or medium with residuals listed)
- [x] Board status + Notes; next **P11-S03-00** (unless spawn)

## Todo updates
Reviewer: status + Notes; spawn `02a`/`02b` only for open blocker/high. Forward-only.
