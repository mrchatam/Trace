# P15-S01-02 REVIEW-NOTES — MCP Assert dispatch

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none  
**Next board:** P15-S02-00

Independent review (fresh subagent). Claims from P15-S01-01 Notes re-verified against live code + locked verify cmds — not trusted alone.

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | Every registered tool Asserts `mcp:<Name>` (incl. `toolVersion`) | PASS | Shared `assertMCPToolAllowed` in `internal/mcp/assert.go` → `AssertToolAllowed(ctx, "mcp:"+toolName)`. Call sites after `openStore`: `toolWhy`, `toolContext`, `toolAdd`, `toolLink`, `toolTransition`, `toolReview`, `toolTasks`, `toolCapability`, `toolVersion` (openStore `""` then Assert). Nine `assertMCPToolAllowed` call sites; nine `AddTool` registrations. |
| 2 | Slugs match `BuiltinMCPCapabilitySpecs` | PASS | Specs + `RegisteredToolNames` identical nine-name order; slugs `mcp:`+name. `TestBuiltinMCPCapabilitySpecs` + `TestToolNamesRegistered` PASS. |
| 3 | DENIED fail-closed | PASS | `TestMCPAssertDeniedBlocksCallTool`: `DecideTool` DENIED on `mcp:trace_why` → CallTool error containing `DENIED`. |
| 4 | Builtin AUTO_ALLOWED succeeds | PASS | `TestMCPAssertBuiltinAutoAllowedSucceeds`: `trace_version` CallTool OK; durable `AUTO_ALLOWED` row for `mcp:trace_version`. |
| 5 | Nine tools; no install/decide MCP | PASS | `TestToolNamesRegistered` wants exactly 9; `registerTools` has only those nine; no `trace_install` / `trace_decide`. |
| 6 | No new mig; G19; no ImpactWalk edits | PASS | Schema still through `013_capability_tool_decisions.sql` only. Helper is thin domain call (G19). Assert logic remains in `domain.AssertToolAllowed`. No MCP edits to `internal/retrieval/impact_walk.go` (R2 `allowContainsOut` late-upgrade still present — expected defer). |
| 7 | Parent-only Assert for review/capability sub-actions | PASS | `toolReview` / `toolCapability` Assert once at entry, then dispatch to `reviewCreate`/`reviewSet` / `capability*` (no second Assert; not separately registered → not ungated CallTool paths). |
| 8 | Carry-forward verify suite | PASS | All `01` locked cmds re-run green (see below). |
| 9 | R2/R3/R4 not falsely claimed fixed | PASS | Implementer Notes claim R1 only. Live ImpactWalk still has R2 pattern; R3/R4 remain disposition wontfix — not “fixed” in S01. |

## Locked verify (re-run 2026-08-17)

```text
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered|TestMCPAssert|TestBuiltinMCPCapabilitySpecs|TestCapabilityDecision'
→ ok mcp, domain; store [no tests to run]

CGO_ENABLED=0 honesty A/B/C+G, replan E, impact F, ablation → ok
CGO_ENABLED=1 Gate H, compat 13, p0x, x0 → ok
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
→ all product pkgs ok
```

## Findings by severity

| Severity | Finding |
|----------|---------|
| blocker | *(none)* |
| high | *(none)* |
| medium | *(none)* |
| low | Optional: no dedicated grep-keeper test that every handler calls `assertMCPToolAllowed` (planner marked optional; code review covers it). |
| nit | `toolReview`/`toolCapability` close store after Assert then re-open in sub-handlers (TOCTOU within one CallTool). Acceptable under FINAL parent-once lock. |

## Residuals for S02

- Import named tests: `TestMCPAssertDeniedBlocksCallTool`, `TestMCPAssertBuiltinAutoAllowedSucceeds`, `TestToolNamesRegistered`.
- Do **not** fail VERIFY for R2 defer / R3–R4 wontfix.
- DR-HANDOFF default remains `no successor`.

## Board

- P15-S01-02 → `done`
- Next runnable → **P15-S02-00**
- No `02a`/`02b` spawn
