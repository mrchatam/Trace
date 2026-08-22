# P16 / S03 / 01 — CLI vs MCP allowlist parity (FINAL locks from 00-PLANNER)

## Metadata
- id: P16-S03-01
- todo_ids: [P16-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-77** per sibling **00-PLANNER FINAL**: gated CLI Asserts `cli:<command>` independently of MCP `mcp:<tool>`. MCP DENIED ≠ CLI DENIED. Default `trace add` stays AUTO_ALLOWED. Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- Live: `cmd/trace/root.go`; `internal/domain/capability_decision.go`; `internal/mcp/assert.go`
- Hunt: `experiments/_bughunt/post-p15/` results `mcp_cli_add.*` + `cap_cli_why.*` + [`POST-P15-BUGHUNT.md`](../../../../../experiments/_bughunt/post-p15/POST-P15-BUGHUNT.md)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do **not** re-debate FINAL locks (dual-slug; no MCP→CLI DENY; full gated list; no unprefixed `add`→`cli:add`; `cli:reindex`→`cli:index` only; whole `capability` ungated).

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Home | Domain builtin CLI specs + Resolve AUTO_ALLOW; thin CLI helper. G19: no slug-policy fork in `cmd/trace` beyond `"cli:"+command` |
| Ungated | `help`/`version`/`init`/`capability` (all subs)/`migrate`/`backup`/`restore`/`auth`/`install`/unknown — **never** Assert |
| Gated | `add` `link` `transition` `review` `why` `context` `tasks` `seed` `impact` `plan` `index`/`reindex` |
| Slugs | `cli:add` … `cli:index` (table in 00). `reindex` → domain fold `cli:reindex`→`cli:index` |
| Builtins | `BuiltinCLICapabilitySpecs()` kind TOOL; **not** merged into `BuiltinMCPCapabilitySpecs()` |
| Resolve | AUTO_ALLOW exact CLI **or** MCP builtin slugs. Reason for CLI: `"builtin CLI command"` |
| Canonicalize | Keep MCP fold; add **only** `cli:reindex`→`cli:index`. Unprefixed `add`/`why` unchanged |
| Helper | `assertCLICommand(ctx, svc, command)` → `AssertToolAllowed(ctx, "cli:"+command)` after store open, before work |
| MCP | **Unchanged** `assertMCPToolAllowed` / nine tools |
| Compat | **14**; no **015+** |
| Forbidden | Shared slug; TTY skip; YOLO flags; new MCP tools; Assert on `install`/`decide`; daemon |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Domain specs | `internal/domain/capability.go` | `BuiltinCLICapabilitySpecs()` (11 titles) |
| Domain resolve | `internal/domain/capability_decision.go` | `isBuiltinCLISlug`; Resolve AUTO_ALLOW CLI; canonicalize `cli:reindex`→`cli:index` |
| Domain tests | `internal/domain/capability_decision_test.go` | AUTO_ALLOW CLI; unprefixed `add`; reindex fold |
| CLI helper | `cmd/trace` (e.g. `assert.go` next to dispatch) | `assertCLICommand` |
| CLI gated cmds | `add.go` `link.go` `transition.go` `review.go` `why.go` `context.go` `tasks.go` `seed.go` `impact.go` `plan.go` `index.go` | After open/`domain.New`, Assert. why/context/tasks/index may `domain.New(st)` **only** to Assert |
| CLI tests | `cmd/trace/*_test.go` | Named isolation + DENIED + ungated decide + index alias |
| MCP tests | `internal/mcp/mcp_test.go` | `TestCLIAddDeniedDoesNotBlockMCPAdd` + keepers |
| MCP production | **Zero** (`assert.go` unchanged) |
| Store / mig | **Zero** (reuse 014) |

Do **not** require a full `openDomain` refactor of every command. Helper + call after existing open is enough. `plan help` / usage paths with no store stay ungated.

## Named tests (required)

| Test | Intent |
|------|--------|
| `TestCLIAddSucceedsWhenMCPAddDenied` | `mcp:trace_add` DENIED → CLI add exit 0 + entity; `cli:add` AUTO_ALLOWED |
| `TestCLIAddDeniedFailClosed` | `cli:add` DENIED → CLI add fails; no entity |
| `TestCLIWhySucceedsWhenMCPWhyDenied` | `mcp:trace_why` DENIED → CLI why exit 0 |
| `TestCLIWhyDeniedFailClosed` | `cli:why` DENIED → CLI why fails |
| `TestCLIAddDeniedDoesNotBlockMCPAdd` | `cli:add` DENIED → CallTool `trace_add` still OK |
| `TestUngatedCapabilityDecideWhenCLIAddDenied` | `cli:add` DENIED → `capability decide` still works |
| `TestUnprefixedAddDecideDoesNotGateCLI` | slug `add` DENIED does not gate `cli:add` |
| `TestCapabilityDecisionAutoAllowBuiltinCLI` | first Resolve `cli:add` AUTO_ALLOWED; independent of `mcp:trace_add` DENIED |
| `TestCanonicalizeCLIReindexFoldsToIndex` | `cli:reindex` ≡ `cli:index` |
| `TestCLIIndexAliasDenied` | `cli:index` DENIED → `index` **and** `reindex` fail |
| Keepers | `TestMCPAssertDeniedBlocksCallTool`, `TestMCPAssertBuiltinAutoAllowedSucceeds`, `TestMCPUnprefixedDecideGatesCallTool`, `TestToolNamesRegistered` (nine), `TestCanonicalizeCustomAndCLISlugsUnchanged`, `TestCapabilityDecisionAutoAllowBuiltinMCP`, `TestMCPVirginProjectDoesNotMkdir`, `TestOpenCreatesDBAndMigratesIdempotent` (v14) |

TDD: named isolation/DENIED tests first (red), then specs + helper call sites (green).

## Role work
1. TDD named tests (red on live zero CLI Assert).
2. `BuiltinCLICapabilitySpecs` + Resolve AUTO_ALLOW + `cli:reindex` fold; CLI helper + gated call sites; **no** Assert on ungated list.
3. Prove green: named tests + locked verify cmds. Do **not** add MCP tools or mig 015.
4. Board **status + Notes only** → next **P16-S03-02**.

## Locked verify commands

```text
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestCLIAddDeniedDoesNotBlockMCPAdd|TestUnprefixedAddDecideDoesNotGateCLI|TestCapabilityDecisionAutoAllowBuiltinCLI|TestCanonicalizeCLIReindexFoldsToIndex|TestCanonicalizeCustomAndCLISlugsUnchanged|TestCapabilityDecision|TestMCPAssert|TestMCPUnprefixed|TestToolNamesRegistered|TestMCPVirgin|TestOpenCreates'

CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestCLIAddSucceedsWhenMCPAddDenied|TestCLIAddDeniedFailClosed|TestCLIWhySucceedsWhenMCPWhyDenied|TestCLIWhyDeniedFailClosed|TestUngatedCapabilityDecideWhenCLIAddDenied|TestCLIIndexAliasDenied|TestUnprefixedAddDecideDoesNotGateCLI'

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
- [ ] DF-77: gated CLI Asserts `cli:<command>`; `mcp:trace_add` DENIED does **not** block `trace add`; `cli:add` DENIED does
- [ ] `cli:why` gated; `capability decide` ungated; `reindex`≡`cli:index`; unprefixed `add` is not a CLI slug
- [ ] Nine MCP tools; Assert helper unchanged; no 015; default `trace add` AUTO_ALLOWED
- [ ] Named tests pass; locked verify cmds PASS
- [ ] Board Notes → **P16-S03-02**

## Minimal todos
- [ ] Named tests (red → green)
- [ ] Builtin CLI specs + Resolve AUTO_ALLOW + reindex fold
- [ ] CLI helper + gated call sites (ungated stay clean)
- [ ] Locked verify suite
- [ ] Board status + Notes only
