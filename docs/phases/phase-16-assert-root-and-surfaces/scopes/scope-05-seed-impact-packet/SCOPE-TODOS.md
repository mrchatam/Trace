# S05 — seed / impact packet + thin MCP — scope todos

**Depends-on:** P16-S04-02 APPROVE (board). Owns **DF-70, DF-71, DF-72, DF-73, DF-74**.

## Depends (from S04 — live)
- **P16-S04-02 APPROVE.** DF-68 install `-C` is **S04**. Do not reopen `cmdInstall` / Cursor STABLE home / CONDITIONAL markers / `cli:install`.
- No install MCP; `install` stays ungated. S05 does not thread `InstallOpts.ProjectRoot` into seed/impact/compiler.
- Catalog was nine including `trace_version` after S04; this scope adds `trace_impact` → **ten including version**.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **FINAL** (P16-S05-00) |
| 2 | 01-seed-impact-packet | implement | **done** (P16-S05-01) |
| 3 | 02-scope-review | review | **done APPROVE** (P16-S05-02) — next **P16-S06-00** |

## Phase locks (P16-S05-00 FINAL)
- Seed `discovery_mentions_task` + hyphen `discovery-mentions-task` → `LinkDiscoveryMentionsTask` (DF-70)
- Compiler packet `impact` + MD `## Impact`; why JSON inherits; SchemaVersion **0.2** (DF-71)
- Thin `trace_impact` MCP (DF-72) — P14 A3 **superseded for this tool only** ([`../../DF-72-FORWARD.md`](../../DF-72-FORWARD.md)); slug `mcp:trace_impact`; Assert at entry
- Seed top-level `findings` / `alternatives` (DF-73); unknown other keys rejected
- Impact report nested JSON snake_case (DF-74)
- Catalog **10** including `trace_version` (`trace_impact` before version)
- **No** install/decide/plan/index MCP; no new entity_links rels; no P14 walk rewrite; compat **14**

## Named tests (S06 import)
`TestSeedImportDiscoveryMentionsTask`, `TestContextIncludesImpactOverallClass`, `TestWhyIncludesImpactOverallClass`, `TestSeedImportImpactFindings`, `TestImpactReportJSONSnakeCase`, `TestMCPTraceImpactReport`, `TestMCPImpactDeniedBlocksCallTool`, updated `TestToolNamesRegistered` / `TestBuiltinMCPCapabilitySpecs` / `TestImportBoundaryMCPNoPlanImpactIndexTools`

## Reminders
- Next after APPROVE: **P16-S06-00**
