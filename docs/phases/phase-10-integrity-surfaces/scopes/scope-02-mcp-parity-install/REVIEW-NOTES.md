# P10-S02-02 — REVIEW-NOTES (MCP parity / install)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-16  
**Spawns:** none

## Checklist evidence

| # | Check | Result |
|---|--------|--------|
| 1 | DF-21 `trace_tasks` + CLI row keys + optional `goal_id` | **Pass** — `server.go` registers tool; `tools_parity.go` emits `{id,title,work_state,goal_id}`; `TestTraceTasksParity` (incl. goal filter + empty) |
| 2 | DF-21 `trace_capability` five actions via domain (G19) | **Pass** — declare/list/require/unrequire/missing → `domain.New` / Upsert/List/Require/Unrequire/Missing; `TestTraceCapabilityActions` |
| 3 | DF-21 no plan/impact/index MCP | **Pass** — nine `Name:` tools only; `TestImportBoundaryMCPNoPlanImpactIndexTools` |
| 4 | DF-21 BuiltinMCP specs + no auto-seed | **Pass** — nine `mcp:trace_*` incl. tasks/capability/version; `TestBuiltinMCPCapabilitySpecs` asserts Open catalog empty |
| 5 | DF-22 README/help rebuild + abs `--bin` + reload | **Pass** — README Install/Cursor MCP + `help.go` install tip |
| 6 | DF-22 `--write` stderr tip + `trace_version` | **Pass** — `install.go` tip; `TestInstallCursorWriteMergeBackup` / `TestInstallCursorWriteCreateMissing`; `TestTraceVersion` `{ok,name,version}`=`trace`/`0.0.0-dev` |
| 7 | DF-32 CLI list/missing snake_case | **Pass** — `capabilityListRow` DTO; `TestCapabilityListMissingSnakeCase` (no `ID`/`Kind` leak); MCP list/missing same keys |
| 8 | Nine tools; `trace-mcp -h` lists new tools | **Pass** — `RegisteredToolNames` len 9; `cmd/trace-mcp/main.go` help lists tasks/capability/version |
| 9 | S01 inherit intact | **Pass** — `TestExactWhyPlanChangeAlias`, `TestExactWhyCapability`, `TestIncludeWhyFailClosed` green under carry-forward |
| 10 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + Gate C | **Pass** — locked CGO0/CGO1 suites + retrieval/compiler/perf PASS; Gate C `dry_run:false` N=3 mean G1 **0.800** > B0 **0.000** intact |
| 11 | Board Notes / planner no product Go | **Pass** — P10-S02-00 Notes claim no product Go; S02-01 Notes cite tests accurately |

## Findings

| Severity | Finding | Disposition |
|----------|---------|-------------|
| — | No blocker/high | — |
| low | `capabilityListRow` duplicated in CLI + MCP (drift risk) | Residual — G19-thin adapters; both snake_case today |
| low | `TestToolNamesRegistered` asserts hand list, not SDK introspection | Residual — functional handler tests cover live tools |
| nit | `trace_capability` `ReadOnlyHint=false` for whole tool (list/missing inclusive) | Acceptable single-tool pattern (same as `trace_review`) |
| nit | Live Cursor MCP in this session still exposes six tools | Expected DF-22 residual until human reloads MCP |

## Residuals (explicit)

1. **Cursor MCP reload** — Trace cannot force-restart the stdio process; agents must rebuild `trace-mcp`, prefer abs `--bin`, reload MCP / window, then use `trace_version` to confirm.
2. **S04** still owns capability/operator transition gating (DF-17/18/24/26/31).
3. **`go test ./...`** still FAIL only on pre-existing `similar projects/graphify` path space (non-product).
4. Duplicated list-row DTO + static tool-name test gap (low) — no spawn.

## Independent verify (this review)

```text
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
CGO_ENABLED=1 go test ./evals/replan/... ./evals/impact/... ./evals/capability/... ./evals/perf/... ./internal/retrieval/... ./internal/compiler/... -count=1  → PASS
CGO_ENABLED=1 go test ./... -count=1  → product pkgs PASS; FAIL only similar projects/graphify (space)
Gate C artifacts: dry_run:false N=3; G1 mean 0.800 > B0 mean 0.000
```

## Next

**P10-S03-00** (no spawn inserted).
