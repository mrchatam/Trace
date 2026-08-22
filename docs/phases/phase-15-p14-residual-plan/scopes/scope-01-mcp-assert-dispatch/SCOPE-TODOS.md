# S01 — MCP Assert dispatch — scope todos

**Depends-on:** P15-00 FINAL (disposition R1=fix). Owns residual **R1** only.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** (FINAL 2026-08-17) |
| 2 | 01-mcp-assert-dispatch | implement | pending — runnable alone |
| 3 | 02-scope-review | review | pending |

## Phase locks (FINAL — do not reopen in implement)
- Wire `AssertToolAllowed` into every MCP CallTool entry via shared helper
- Slug = `mcp:` + registered tool `Name` (exact match `BuiltinMCPCapabilitySpecs`)
- Named tests: `TestMCPAssertDeniedBlocksCallTool`, `TestMCPAssertBuiltinAutoAllowedSucceeds`, keep `TestToolNamesRegistered`
- `trace_version` must open store for Assert (same project resolution as peers)
- **No** new MCP tools; keep nine + `trace_version`
- **No** new migration; reuse mig 013 + domain APIs
- Do **not** implement R2/R3/R4
- Do **not** board S05 / plan simulate / D21+

## Depends (to S02 — light)
After S01 APPROVE: S02 VERIFY imports named MCP Assert regressions (`TestMCPAssert*`) + `TestToolNamesRegistered` + carry-forward gates; DR-HANDOFF default `no successor`.

## Reminders
- G19: MCP remains thin adapter; Assert lives in domain
- Forward-only board; implementers: status + Notes only
- Next after APPROVE: **P15-S02-00**
