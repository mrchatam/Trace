# P31-S01-02 — REVIEW-NOTES

**Reviewer:** independent fresh session (not S01-01 implementer)  
**Date:** 2026-08-21  
**Verdict:** **PASS** — G1, G5, G6 closed; no open blocker/high  
**Confidence:** **high** (re-ran `go test ./internal/...` + `bash scripts/repro-stray-trace-db.sh`)  
**Spawns:** none

## Claimed vs evidence

| ID | Claim (S01-01 Notes) | Evidence | Status |
|----|----------------------|----------|--------|
| G1 | `TestOpenQuietWhenRootStubIsDirectory` | `internal/store/stray_trace_db_test.go:155–188`: `os.Mkdir` root `dbFileName`; empty `warnWriter`; Open OK; `DBPath` = `.trace/trace.db`; dir stub still `IsDir` | **closed** |
| G5 | `scripts/repro-stray-trace-db.sh` | Script present, executable; reviewer run → `ALL PASS` (exit 0); greps `project-root trace.db exists…`, `.trace/trace.db`, `agents: use CLI/MCP`; stub size/mtime unchanged | **closed** |
| G6 | CONTRIBUTING + AGENTS multi-open | `CONTRIBUTING.md:83` + `AGENTS.md:9`: once per `openStore`; multi CLI/MCP/HTTP may re-emit; no suppress flag — matches `open.go:19–21` comment + `warnIfStrayRootTraceDB` Stat/`IsRegular` only | **closed** |

**Out of bar:** G2 absent (OK). G3/G4 not reopened.

## Cross-cutting checklist

| Check | Result |
|-------|--------|
| Join still `.trace` + `trace.db` | `open.go:16–17`, `87`, `102`, `174` |
| Warn non-fatal; Stat-only regular-file; no delete/rename | `open.go:142–149`; no `Remove`/`Rename` in `open.go` |
| Existing four stray tests + G1 | All five present; `go test ./internal/store/ -count=1 -run 'TestOpen…'` PASS |
| `go test ./internal/...` | PASS (exit 0) |
| `/trace.db` gitignore | `.gitignore:4` + `fixtures/x0/.gitignore:2` |
| Hard bans | No path redesign, silent delete, GUI, or suppress flag in product/test/script/docs for this scope |

## Findings

| Severity | Finding |
|----------|---------|
| — | **None** at blocker / high / medium |

### Nits (non-blocking; no spawn)

1. **AGENTS.md orchestrator paste** still says Phase 31 “Next: `P31-S00-00`” under Current focus — stale vs board; unrelated to G6 bullet (L9), which is correct. Optional cleanup in S02 docs pass.
2. **G5 script** uses Linux `stat -c` — fine for this repo’s Linux CI/dev; not portable to macOS BSD `stat` without change.

## Re-run commands (reviewer)

```text
go test ./internal/...          → EXIT 0
bash scripts/repro-stray-trace-db.sh → ALL PASS, EXIT 0
```

## Next

**P31-S02-00** (VERIFY planner)
