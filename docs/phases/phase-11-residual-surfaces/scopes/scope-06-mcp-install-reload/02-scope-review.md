# P11 / S06 / 02 — Scope review (MCP / install reload UX)

## Metadata
- id: P11-S06-02
- todo_ids: [P11-S06-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S06 (**DF-22, DF-37, DF-50**). Fresh subagent. Compare claims + locks to live code/tests. Spawn `02a`/`02b` for blocker/high. Do not rewrite prior `done` history.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- Sibling [01-mcp-install-reload.md](01-mcp-install-reload.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-22, DF-37, DF-50
- Live: `cmd/trace/{install,help,install_test}.go`; README Install; `internal/mcp` tool list / `trace_version`
- Prior: P10 S02 DF-21/22/32; P11-S05 no install coupling

## Checklist (must all pass for APPROVE)

| # | Check |
|---|--------|
| 1 | DF-50: print-only success emits same stderr tip as `--write`; stdout remains JSON-only |
| 2 | DF-22: help + README tip cover print and write; `trace_version` still registered; tip text includes `reload` + `trace-mcp` |
| 3 | DF-37: ops closed via tip/docs only — no PID kill, daemon, HTTP, or new MCP tools |
| 4 | P10 nine-tool set + mcp.json merge/backup semantics retained; G19 — no domain fork in adapters |
| 5 | No forbidden architecture (daemon/HTTP/full-rebuild) |
| 6 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + P11-S01–S05 + Gate C `dry_run:false` |
| 7 | Board Notes accurate; planner row had no product Go |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Prefer named asserts: `TestInstallCursorPrintReloadTip` (or extended print test), `TestInstallCursorWriteMergeBackup`, `TestInstallCursorWriteCreateMissing`, `TestToolNamesRegistered` / `TestTraceVersion` still green (or equiv).

## Exit criteria
- [x] Checklist evidenced; confidence high (or medium with residuals)
- [x] Board status + Notes; next **P11-S07-00** (unless spawn)
- [x] Write [REVIEW-NOTES.md](REVIEW-NOTES.md) on APPROVE / spawn
