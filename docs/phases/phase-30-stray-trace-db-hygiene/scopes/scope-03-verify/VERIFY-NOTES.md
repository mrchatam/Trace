# VERIFY-NOTES — P30-S03-01

**Date:** 2026-08-21
**Git SHA:** unavailable (workspace has no `.git` / `git rev-parse` fails; recorded as `unavailable-in-sandbox` in run metadata)
**Overall:** PASS
**Evidence:** experiments/runs/2026-08-21-p30-s03-01-verify/evidence/

| Block | Result | Notes |
|------:|--------|-------|
| 0 Preflight | PASS | `EVID` created; metadata cites `P30-S02-02 PASS; PLAN T1-T4; no store-path change` |
| 1 Store stray tests | PASS | All four: `TestOpenWarnsWhenRootStubPresent`, `TestOpenExistingWarnsWhenRootStubPresent`, `TestOpenQuietWhenNoRootStub`, `TestOpenLeavesRootStubUntouched` — PASS (`01-store-stray-tests.txt`) |
| 2 go test ./internal/... | PASS | Exit 0; all packages ok (`02-internal-test.txt`) |
| 3 Temp repro (init / stub / warn / untouched) | PASS | `bin/trace` rebuilt; init → no root `trace.db`, `.trace/trace.db` present; python 0-byte stub; `tasks` uses `.trace/`; stderr warn with locked substrings; stub size+mtime unchanged (`03-repro-init.txt`, `03-warn-stderr.txt`, `03-warn-verdict.txt`) |
| 4 Docs/gitignore/join | PASS | `/trace.db` in `.gitignore` + `fixtures/x0/.gitignore`; AGENTS/project-rules/CONTRIBUTING state live store `.trace/trace.db`; `warnIfStrayRootTraceDB` Stat-only; join `.trace`+`trace.db` (`04-docs-gitignore-join.txt`) |
| 5 Residuals | listed | See below — non-blocking |

## Residuals (non-blocking)

- Agents can still `sqlite3.connect('trace.db')` and create stubs — mitigated by warn + gitignore; cannot prevent all agent mistakes.
- Optional documented delete of root stub — **future-only**, not Phase 30.
- Warn once per `openStore` (CLI `tasks` emitted warn multiple times because multiple opens) — acceptable; no persistent suppress flag.
- Git SHA unavailable in this workspace (no `.git`) — does not affect behavioral VERIFY.

## Failures (if any)

- None.

## DR-HANDOFF

Remains **OPEN** — close owned by **P30-S03-02**.

## Next

**P30-S03-02**
