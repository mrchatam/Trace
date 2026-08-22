# P16-S05-00 — Impact packet + thin MCP (stub — thicken vs live)

## Metadata
- id: P16-S05-00
- todo_ids: [P16-S05-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock **FINAL** for **DF-71** (context/why include impact), **DF-74** (snake_case report JSON), **DF-72** (thin MCP `trace_impact`). This scope **supersedes** P14 A3 “no MCP impact” for **one** G19 adapter. **No product Go.**

## Inherited locks (phase)
- **DF-71:** Task `context` / `why` packets include linked-decision impact findings + `overall_class` (bounded; Law 6 — no full dump). Prefer compiler/retrieval, not CLI fork
- **DF-74:** `trace impact report` JSON snake_case (`id`, `decision_id`, `impact_class`, `is_recommended`, …) matching `tasks`
- **DF-72:** MCP tool `trace_impact` with `action` (or equivalent) covering CLI subset: `finding` / `alternative` / `report` / `walk`. Builtin slug `mcp:trace_impact`. Assert at tool entry (P15 helper)
- **Supersede** `TestImportBoundaryMCPNoPlanImpactIndexTools`: **allow** `trace_impact`; still **forbid** `trace_plan` / `trace_index` / install / decide MCP
- Update `TestToolNamesRegistered` + `BuiltinMCPCapabilitySpecs` to **ten** tools + `trace_version`
- No `trace_install` / `trace_capability_decide` MCP
- Walk semantics unchanged (no R2 fix)

## Named tests (phase hint)
- `TestContextIncludesImpactOverallClass` (and why if 00 includes it)
- `TestImpactReportJSONSnakeCase`
- `TestMCPTraceImpactReport`
- Keepers: `TestToolNamesRegistered`; updated import-boundary test

## Planner work
1. [ ] Confirm compiler/why packet shape + impact CLI JSON encoding
2. [ ] Lock MCP input schema (action enum) vs CLI
3. [ ] Thicken 01/02; **FINAL**; next **P16-S05-01**
