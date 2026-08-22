# P16 / S05 / 01 — Impact packet + thin MCP (DF-71, 72, 74)

## Metadata
- id: P16-S05-01
- todo_ids: [P16-S05-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Surface impact findings on context/why, snake_case report JSON, and a thin MCP `trace_impact` adapter (G19). Board **status + Notes only**.

**Stop if `00-PLANNER.md` is not FINAL.**

## Locked defaults
| Item | Value |
|------|-------|
| DF-71 | Compiler/retrieval packet; bounded; include `overall_class` |
| DF-74 | CLI report DTO snake_case (do not leak Go field names) |
| DF-72 | One MCP tool; Assert `mcp:trace_impact`; no domain fork |
| Boundary test | Allow `trace_impact` only among previously forbidden names |
| Forbidden | install/decide/plan/index MCP; R2 ImpactWalk fix; daemon |

## Named tests
See 00-PLANNER.

## Locked verify (minimum)
```text
CGO_ENABLED=0 go test ./internal/compiler/... ./internal/retrieval/... ./internal/mcp/... ./cmd/trace/... ./evals/impact/... -count=1 -run 'TestContextIncludesImpact|TestImpactReportJSONSnakeCase|TestMCPTraceImpact|TestToolNamesRegistered|TestPlantedImpactConflictsGateFPrelim'
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

## Exit criteria
- [ ] DF-71/72/74 named tests pass; tool catalog = ten + version
- [ ] Gate F still green; no R2 claim
- [ ] Board → **P16-S05-02**
