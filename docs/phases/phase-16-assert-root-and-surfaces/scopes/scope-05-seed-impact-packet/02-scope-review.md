# P16 / S05 / 02 — Scope review (seed / impact packet + thin MCP) FINAL checklist

## Metadata
- id: P16-S05-02
- todo_ids: [P16-S05-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of S05 **DF-70, DF-71, DF-72, DF-73, DF-74**. Fresh subagent ≠ implementer. Confirm thin `trace_impact` only (P14 A3 superseded for this tool — [`../../DF-72-FORWARD.md`](../../DF-72-FORWARD.md)). Spawn `P16-S05-02a`/`02b` for blocker/high. Prefer `REVIEW-NOTES.md`. Next **P16-S06-00** unless spawn.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-seed-impact-packet.md](01-seed-impact-packet.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [phase README](../../README.md)
- Live: `cmd/trace/seed.go`, `internal/compiler`, `cmd/trace/impact.go`, `cmd/trace/why.go`, `internal/mcp`, `internal/store/impact.go`, `internal/domain/capability.go`
- Dogfood: D22 `impact_report.json` (historical PascalCase); tests are SoT
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do not share the implementer’s session. Re-run verify; do not trust Notes alone. Do not re-open daemon, install/decide MCP, S04 `-C`, P14 walk/R2, or packet schema 0.3 (FINAL).

## Checklist

| # | Check | How to evidence |
|---|--------|-----------------|
| 1 | DF-70 seed link rel + hyphen alias → `LinkDiscoveryMentionsTask` | `TestSeedImportDiscoveryMentionsTask`; grep seed switch for both tokens |
| 2 | DF-71 why/context include findings / `overall_class` when present | `TestContextIncludesImpactOverallClass`; `TestWhyIncludesImpactOverallClass`; Packet `SchemaVersion` still `0.2`; no forked rollup (calls `ImpactReport`) |
| 3 | DF-73 seed `findings`/`alternatives` via existing impact domain | `TestSeedImportImpactFindings`; allowlist; `impact_findings` still rejected; no new rels/tables |
| 4 | DF-74 snake_case impact report JSON | `TestImpactReportJSONSnakeCase`; store json tags; no `ImpactClass` in report stdout |
| 5 | DF-72 thin `trace_impact` MCP; Assert `mcp:trace_impact`; G19 | `TestMCPTraceImpactReport`; `TestMCPImpactDeniedBlocksCallTool`; action covers finding/alternative/report/walk; adapter calls domain/walk only |
| 6 | Boundary keeper **allows** `trace_impact` and still forbids plan/index/install/decide MCP | Updated `TestImportBoundaryMCPNoPlanImpactIndexTools` |
| 7 | Catalog ten including version; specs updated; Gate F keeper green | `TestToolNamesRegistered` len 10; `TestBuiltinMCPCapabilitySpecs`; `TestPlantedImpactConflictsGateFPrelim`; CLI specs still have `cli:impact` only (no `cli:trace_impact`) |
| 8 | Carry-forward + product pkgs PASS; no 015; S04 DashC keeper | Re-run `01` locked verify cmds; compat **14**; `TestInstallClaudeDashCRefuseCitesProjectRoot` |

## Review procedure
1. Read FINAL locks in 00 + 01.
2. Map named tests → code; fresh verify cmds from `01`.
3. APPROVE (high, or medium with residuals listed) or spawn `P16-S05-02a`/`02b` with full prompts.
4. Write `REVIEW-NOTES.md` + board Notes; next **P16-S06-00** unless spawn.
5. If APPROVE: S06 must import the named tests in 00 (already pointed on S06 SCOPE-TODOS).

## Exit criteria
- [ ] Checklist evidenced; confidence high (or medium with residuals listed)
- [ ] REVIEW-NOTES.md written
- [ ] Board status + Notes; next **P16-S06-00** (unless spawn)
- [ ] No rewrite of done P16-S05-00/01 history

## Minimal todos
- [ ] Independent verify + checklist
- [ ] REVIEW-NOTES.md
- [ ] Board sync
