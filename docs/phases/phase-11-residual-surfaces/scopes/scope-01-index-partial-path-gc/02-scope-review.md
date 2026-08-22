# P11 / S01 / 02 — Scope review (Index partial-path GC)

## Metadata
- id: P11-S01-02
- todo_ids: [P11-S01-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S01 (**DF-40**). Fresh subagent. Compare claims + FINAL locks to live code/tests. Spawn `02a`/`02b` for blocker/high. Do not rewrite prior `done` history. Small inline fixes OK for medium/low when cheaper than spawn.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-index-partial-path-gc.md](01-index-partial-path-gc.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-40
- Live: `cmd/trace/index.go`; `internal/store/file_graph.go`; `cmd/trace/cli_test.go`
- P10 S03 DF-20 baseline must remain green

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (reviewer).

## Checklist (must all pass for APPROVE)

| # | Check |
|---|--------|
| 1 | DF-40: rename + `index <new-path>` removes old path/symbols/FTS (test + code path) |
| 2 | Mechanism is content-hash orphan GC (or equivalent meeting locks) — **not** project-wide argv set-diff |
| 3 | Isolation: on-disk siblings survive partial argv (`TestIndexIncrementalIsolation`) |
| 4 | P10 DF-20 retained: full-tree GC + missing-argv single delete |
| 5 | G19 — no domain fork in adapters; analyzers stay upsert-only |
| 6 | No forbidden architecture (daemon/HTTP/embeddings/full-rebuild/new mig/MCP index) |
| 7 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + Gate C `dry_run:false` |
| 8 | Board Notes accurate; planner row had no product Go |

## Focus questions
- Does partial GC only delete **missing-on-disk** hash matches (not live duplicates)?
- Is `removed` count observable for the rename case?
- Any accidental full-tree GC when `len(args)>0`?

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Reviewer loop
1. Diff claims (01 Notes + locks) vs repo.
2. Findings by severity; blocker/high → inline fix or spawn `P11-S01-02a`/`02b` immediately below this row.
3. Re-verify until no open blocker/high without pending follow-up.
4. Write [REVIEW-NOTES.md](REVIEW-NOTES.md) with confidence + residuals.
5. Board status + Notes; next **P11-S02-00** unless spawn.

## Exit criteria
- [ ] Checklist evidenced; confidence **high** (or **medium** with residuals listed)
- [ ] Board status + Notes; next **P11-S02-00** (unless spawn)
