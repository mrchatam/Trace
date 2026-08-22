# P16 / S05 / 01 — seed / impact packet + thin MCP (FINAL locks from 00-PLANNER)

## Metadata
- id: P16-S05-01
- todo_ids: [P16-S05-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-70, DF-71, DF-72, DF-73, DF-74** per sibling **00-PLANNER FINAL**. Thin MCP `trace_impact` only (P14 A3 superseded for this tool — [`../../DF-72-FORWARD.md`](../../DF-72-FORWARD.md)). Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- Live: `cmd/trace/seed.go`; `internal/compiler`; `cmd/trace/impact.go`; `cmd/trace/why.go`; `internal/mcp`; `internal/store/impact.go`; `internal/domain/impact.go`; `internal/domain/capability.go`
- Dogfood: [`BATCH-D21-D23.md`](../../../../../experiments/BATCH-D21-D23.md); D22 [`impact_report.json`](../../../../../experiments/ab-combo-context-impact/results/_surfaces/impact_report.json)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do **not** re-debate FINAL locks (seed keys `findings`/`alternatives`; packet `impact` omitempty; catalog **10 including version**; thin `trace_impact` only). Do **not** reopen S04 `-C` / install MCP / daemon. **No board spawn.**

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| DF-70 | Seed switch: `discovery_mentions_task` + `discovery-mentions-task` → `LinkDiscoveryMentionsTask` |
| DF-73 | Seed top-level `findings` / `alternatives` → `AddImpactFinding` / `AddDecisionAlternative`; unknown other keys still rejected; stub names `impact_findings` stay unknown |
| DF-71 | Packet `impact []DecisionImpact` omitempty; SchemaVersion **0.2**; attach via `ImpactReport` when `len(Findings)>0`; MD `## Impact`; why JSON inherits (do not change Why BFS) |
| DF-74 | `json` tags on store finding + alternative structs; report nested snake_case (`id`, `impact_class`, `is_recommended`, …) |
| DF-72 | MCP tool `trace_impact`; slug `mcp:trace_impact`; Assert at entry; `action` = finding\|alternative\|report\|walk; `op` = add\|list\|recommend as locked |
| Catalog | Ten names: … `trace_capability`, `trace_impact`, `trace_version`. Update `RegisteredToolNames` + `BuiltinMCPCapabilitySpecs` |
| Boundary | Allow `trace_impact`; forbid `trace_plan` / `trace_index` / install / decide MCP |
| CLI specs | `BuiltinCLICapabilitySpecs` **unchanged** (`cli:impact` already exists) |
| Compat | **14**; no **015+**. Do **not** edit `014_*.sql` Name IN list |
| Forbidden | Daemon; new entity_links rels; rewriting P14 walk (R2); install/decide MCP; `cli:install`; packet 0.3 |

Copy exact seed object keys, MCP input table, and catalog order from **00-PLANNER**. Do not invent `impact_findings` as an accepted alias.

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Seed CLI | `cmd/trace/seed.go` | Rel case + allowlist + `seedFinding`/`seedAlternative` structs + domain calls + summary counts |
| Seed/CLI tests | `cmd/trace/cli_test.go` | Named DF-70/73/74/71-why tests (CGO1) |
| Why CLI | `cmd/trace/why.go` | After `Why()`, merge `impact` from domain helper (G19) |
| Store | `internal/store/impact.go` | snake_case `json` tags on finding + alternative structs |
| Domain | `internal/domain/impact.go` | Optional `ImpactSummariesForTask`; **no** rollup change |
| Domain specs | `internal/domain/capability.go` + `capability_test.go` | Ten MCP slugs including `mcp:trace_impact` |
| Compiler | `internal/compiler/packet.go`, `compiler.go`, `compiler_test.go` | `impact` field + attach + Markdown; may import domain |
| MCP | `internal/mcp/server.go`, **new** `tools_impact.go`, `export_test.go`, `mcp_test.go` | Register tool; action switch; CallImpact seam; named tests + catalog/boundary updates |
| Help | `cmd/trace/help.go` | Only if seed/link usage lists rels — add mentions-task if seed help exists; do not rewrite impact walk help semantics |
| Install / walk eval | **Zero** production | Gate F keeper only |

Do **not** thread `InstallOpts.ProjectRoot` into these files. Do **not** add MCP install/decide/plan/index tools.

## Named tests (required)

| Test | Intent |
|------|--------|
| `TestSeedImportDiscoveryMentionsTask` | Underscore + hyphen rels write `discovery_mentions_task` |
| `TestContextIncludesImpactOverallClass` | Packet JSON+MD include DESTRUCTIVE `overall_class` / `impact_class`; zero-finding neighbor omitted |
| `TestWhyIncludesImpactOverallClass` | `trace why task` JSON includes `impact` + DESTRUCTIVE |
| `TestSeedImportImpactFindings` | Seed findings+alternatives round-trip; `impact_findings` key rejected |
| `TestImpactReportJSONSnakeCase` | Report JSON snake_case nested keys; no `ImpactClass` / `IsRecommended` |
| `TestMCPTraceImpactReport` | MCP `action=report` snake_case + overall_class |
| `TestMCPImpactDeniedBlocksCallTool` | DENIED `mcp:trace_impact` fail-closed |
| Keepers | `TestToolNamesRegistered` (**10**), `TestBuiltinMCPCapabilitySpecs` (**10**), updated `TestImportBoundaryMCPNoPlanImpactIndexTools`, `TestMCPVirginProjectDoesNotMkdir`, `TestCLIAddDeniedFailClosed`, `TestPlantedImpactConflictsGateFPrelim`, `TestOpenCreates*` (v14) |

TDD: named tests red on live gaps, then implement. Do not change `hasClaudeMarker`, Why BFS, or ImpactReport rank tables to pass.

## Role work
1. TDD named tests (red: unknown seed rel; `findings` unknown key; packet has no `impact`; report PascalCase; MCP name missing / boundary forbids `trace_impact`).
2. Seed switch + allowlist + domain calls; store json tags; compiler attach + MD; why merge; MCP `trace_impact` + catalog/specs/boundary/Assert.
3. Prove green: named tests + locked verify cmds. Do **not** add install/decide MCP, mig 015, or walk algorithm changes.
4. Board **status + Notes only** → next **P16-S05-02**.

## Locked verify commands

```text
CGO_ENABLED=0 go test ./internal/compiler/... ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestContextIncludesImpactOverallClass|TestToolNamesRegistered|TestBuiltinMCPCapabilitySpecs|TestImportBoundaryMCPNoPlanImpactIndexTools|TestMCPTraceImpactReport|TestMCPImpactDeniedBlocksCallTool|TestMCPVirgin|TestOpenCreates'

CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedImportDiscoveryMentionsTask|TestSeedImportImpactFindings|TestImpactReportJSONSnakeCase|TestWhyIncludesImpactOverallClass|TestCLIAddDeniedFailClosed|TestLinkDiscoveryMentionsTaskCLI|TestInstallClaudeDashCRefuseCitesProjectRoot'

CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Product bar = `./cmd|internal|evals`. Compat ceiling **14**. Named `cmd/trace` tests are **CGO1** (R4: CGO0 `./cmd/trace/...` tree-sitter). Do **not** fail this row for R3 graphify space-in-path on full-module `./...` if present outside product pkgs.

## Exit criteria
- [ ] DF-70/71/72/73/74 named tests green
- [ ] Catalog = ten tools including `trace_version`; `mcp:trace_impact` Asserted; boundary allows impact only; no install/decide/plan/index MCP
- [ ] Gate F keeper green; walk unrewritten; no 015; compat **14**
- [ ] Named tests pass; locked verify cmds PASS
- [ ] Board Notes → **P16-S05-02**

## Minimal todos
- [ ] Named tests (red → green) for DF-70/71/73/74/72
- [ ] Seed rel + findings/alternatives allowlist
- [ ] Compiler `impact` attach + why inherit + store json tags
- [ ] MCP `trace_impact` + catalog 10 + boundary + Assert
- [ ] Locked verify suite
- [ ] Board status + Notes only
