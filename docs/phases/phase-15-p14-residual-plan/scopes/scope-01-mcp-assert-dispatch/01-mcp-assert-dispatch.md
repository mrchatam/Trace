# P15 / S01 / 01 — MCP Assert dispatch (FINAL locks from 00-PLANNER)

## Metadata
- id: P15-S01-01
- todo_ids: [P15-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Wire `domain.Service.AssertToolAllowed` into every MCP CallTool path (residual **R1**) per sibling **00-PLANNER FINAL**. Keep exactly nine tools + `trace_version`. No new mig. No domain Assert semantics change. Board **status + Notes only**.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- Live: `internal/mcp/{server,project,tools_*.go}`, `internal/domain/{capability_decision,capability}.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do **not** re-debate FINAL locks. Grill only if go-sdk forces a registration shape that cannot call a shared helper (unlikely — handlers are already `s.tool*`).

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Home | `internal/mcp` thin adapter (G19) |
| Assert API | `domain.New(st).AssertToolAllowed(ctx, slug)` |
| Slug | `"mcp:" + toolName` where `toolName` ∈ `RegisteredToolNames()` / MCP `Tool.Name` — must match `BuiltinMCPCapabilitySpecs` |
| Helper | Shared unexported helper e.g. `assertMCPToolAllowed(ctx, st *store.Store, toolName string) error` → `AssertToolAllowed(ctx, "mcp:"+toolName)` |
| Call sites | **Entry** of `toolWhy`, `toolContext`, `toolAdd`, `toolLink`, `toolTransition`, `toolReview`, `toolTasks`, `toolCapability`, `toolVersion` — once per CallTool (not again inside review/capability sub-actions) |
| `trace_version` | Must `openStore` (empty project override → server defaultRoot/cwd), Assert, then return `{ok,name,version}` as today |
| Fail-closed | PENDING/DENIED → handler `error` (CallTool fails); no silent success |
| Builtins | AUTO_ALLOWED on first resolve — default path stays green |
| No new tools / no new mig | Hard |
| Forbidden | R2/R3/R4; S05; plan simulate; D21+; YOLO/AllowAll; ImpactWalk edits; install/decide MCP tools; daemon/HTTP/embeddings; rewriting Phase 00–14 history |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| MCP helper | `internal/mcp/*.go` (e.g. `project.go` or small `assert.go`) | Shared `assertMCPToolAllowed` |
| MCP handlers | `tools_why.go`, `tools_context.go`, `tools_write.go`, `tools_parity.go`, `server.go` only if needed | Call helper after `openStore` (version: add openStore) |
| MCP tests | `internal/mcp/mcp_test.go` (+ helpers) | Named Assert regressions |
| Domain / store | Prefer **zero** edits | Reuse Assert + mig 013 |

## Named tests (required)

| Test | Intent |
|------|--------|
| `TestMCPAssertDeniedBlocksCallTool` | `DecideTool` DENIED on a builtin slug (e.g. `mcp:trace_why`) → that tool’s CallTool returns error |
| `TestMCPAssertBuiltinAutoAllowedSucceeds` | Fresh project DB; CallTool on a builtin succeeds (AUTO_ALLOWED path) |
| `TestToolNamesRegistered` | Still exactly the locked nine names (keeper) |

## Role work
1. TDD: add named MCP Assert tests first (red).
2. Add shared helper; call from every tool entry including `toolVersion` (openStore for version).
3. Prove green: named tests + locked verify cmds.
4. Board **status + Notes only** → next **P15-S01-02**.

## Locked verify commands

```text
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered|TestMCPAssert|TestBuiltinMCPCapabilitySpecs|TestCapabilityDecision'

CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Product bar = `./cmd|internal|evals`. Do **not** fail this row for R3 graphify space-in-path on full-module `./...` if present outside product pkgs.

## Exit criteria
- [ ] Every registered tool path Asserts with `mcp:<Name>` (incl. `trace_version`)
- [ ] Named tests pass; DENIED blocks; builtin AUTO_ALLOWED succeeds
- [ ] No new MCP tools; no new mig; G19 intact
- [ ] Locked verify cmds PASS
- [ ] Board Notes → **P15-S01-02**

## Minimal todos
- [ ] Named tests (red → green)
- [ ] Shared helper + wire all nine handlers
- [ ] Locked verify suite
- [ ] Board status + Notes only
