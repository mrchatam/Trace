# Scope 02 — Implement todos

| ID | Task | Status | Notes |
|----|------|--------|-------|
| T1 | Docs: AGENTS / project-rules / CONTRIBUTING (+ optional help) | done | P30-S02-01; help.go init one-liner |
| T2 | Warn once in `openStore` (stderr / injectable writer) | done | `warnIfStrayRootTraceDB` + `warnWriter` |
| T3 | `/trace.db` in `.gitignore` + `fixtures/x0/.gitignore` | done | leading slash; `.trace/` kept |
| T4 | Tests: stub present/absent/untouched | done | `stray_trace_db_test.go`; `go test ./internal/...` PASS |

SoT: [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md). Board: `P30-S02-01` implement done; `P30-S02-02` review **PASS**.
