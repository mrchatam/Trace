# P10 / S03 / 02 — Scope review (index GC)

## Metadata
- id: P10-S03-02
- todo_ids: [P10-S03-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S03 (**DF-20**). Fresh subagent. Compare claims + locks to live code/tests. **Reject** full-rebuild-as-default. Small inline fix **or** spawn `02a`/`02b` for blocker/high. Do not rewrite S01/S02/`done` history.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) FINAL locks + [01-index-gc.md](01-index-gc.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-20
- Live: `cmd/trace/index.go`; `internal/store/file_graph.go` (+ FTS); `cmd/trace` index tests

## Session start
Agent → clarify → Plan → execute (reviewer).

## Checklist (must all pass for APPROVE)

| # | Check | Evidence |
|---|--------|----------|
| 1 | **DF-20** — full-tree index removes DB paths absent from walk set | `index.go` + rename test |
| 2 | Delete clears **files + symbols + imports + FTS** (not orphan FTS) | store delete + SyncFileFTS mirror |
| 3 | Explicit argv index **does not** project-wide GC | `TestIndexIncrementalIsolation` |
| 4 | Missing explicit argv path deletes **that** path only (per locks) | code path / test |
| 5 | **No** full-rebuild-on-any-change (set-diff delete only) | architecture read |
| 6 | **No** new migration `011_*` (unless Notes justify blocker) | schema dir |
| 7 | **No** MCP index tools; analyzers stay upsert-only | grep / import boundary |
| 8 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + S01/S02 + Gate C `dry_run:false` | locked verify cmds |
| 9 | Board Notes accurate; planner row had no product Go | TODO.md |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Severity / spawn policy
- **blocker/high:** inline fix if tiny; else spawn `P10-S03-02a` (implement) + `P10-S03-02b` (review) immediately below this row
- **medium:** prefer spawn unless trivial
- **Residual OK:** file-local only (no dependent cascade); experiment ab-index may still need agents to run `trace index` after rename — product must make that reindex **correct**, not automatic watchers

## Exit criteria
- [ ] Checklist 1–9 evidenced
- [ ] Confidence **high** (or **medium** with residuals listed — never silent)
- [ ] No open blocker/high without pending follow-up
- [ ] Board status + Notes; next **P10-S04-00** (unless spawn)

## Todo updates
Reviewer: status + notes; may spawn forward; may thicken **upcoming** S04/S05 prompts if blast radius requires. Do not edit `done` history.
