# VERIFY-NOTES — P31-S02-01

**Date:** 2026-08-21
**Git SHA:** unavailable-no-.git-in-workspace
**Overall:** PASS
**Evidence:** experiments/runs/2026-08-21-p31-s02-01-verify/evidence/
**Precondition:** P31-S01-02 PASS high (G1+G5+G6)

| Block | Result | Notes |
|------:|--------|-------|
| 0 Preflight | PASS | Evidence dir + metadata; SHA unavailable (no `.git` in workspace, same as P30 VERIFY) |
| 1 Store stray tests (5) | PASS | All five PASS (`-v`): WarnsWhenRootStubPresent, ExistingWarns, QuietWhenNoRootStub, LeavesRootStubUntouched, QuietWhenRootStubIsDirectory (G1); exit 0 |
| 2 go test ./internal/... | PASS | exit 0; all packages ok |
| 3 Repro script G5 | PASS | `scripts/repro-stray-trace-db.sh` executable; ALL PASS (init no root stub; live `.trace/trace.db`; stub untouched; warn on stderr) |
| 4 Docs/gitignore/join/G6 | PASS | `/trace.db` in `.gitignore` + `fixtures/x0/.gitignore`; G6 once-per-`openStore`/multi-open/no suppress in CONTRIBUTING L83 + AGENTS L9; `warnIfStrayRootTraceDB` + `.trace`/`trace.db` join; `IsRegular` gate; no `os.Remove`/`os.Rename` in `open.go` |
| 5 Residuals | listed | Non-blocking only; did not fail VERIFY |

## Residuals (non-blocking)
- G2 CLI stderr unit absent — nice-to-have; G5 script + store units cover CLI path
- G3 serve “startup” warn harness — deferred; request-scoped `store.Open` only
- G4 extra `/trace.db` in web/experiment ignores — out-of-scope / deferred
- multi-open once-per-openStore — intentional; G6 documents; no suppress flag
- agents can still create stubs — mitigated by warn + gitignore; agent hygiene
- optional delete future-only — not this phase
- G5 script uses Linux `stat -c` — acceptable for this repo’s Linux CI/dev (S01 nit)

## Failures (if any)
- None

## DR-HANDOFF
remains OPEN — close owner **P31-S02-02**

## Next
P31-S02-02
