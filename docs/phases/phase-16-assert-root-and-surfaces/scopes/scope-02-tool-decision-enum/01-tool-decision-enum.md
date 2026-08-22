# P16 / S02 / 01 — Tool-decision enum + slug prefix (FINAL locks from 00-PLANNER)

## Metadata
- id: P16-S02-01
- todo_ids: [P16-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-75** + **DF-78** per sibling **00-PLANNER FINAL**. CHECK + YOLO fail-closed (heal→PENDING, Resolve never AUTO_ALLOWs garbage). Canonicalize registered MCP Names to `mcp:`+Name. Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- Live: `internal/store/schema/013_capability_tool_decisions.sql`, `internal/store/capability_decision.go`, `internal/domain/capability_decision.go`, `cmd/trace/capability.go`, `internal/mcp/assert.go`
- Pattern: `internal/store/schema/012_import_provenance_enum.sql`
- Hunt: `experiments/_bughunt/post-p15/{cap-decisions,mcp-yolo,mcp-footgun}/`
- Compat: `evals/compat/compat_test.go` (ceiling 13 → **14**)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do **not** re-debate FINAL locks (mig 014 CHECK; YOLO heal→PENDING; Resolve fail-closed; exact `mcp:` canonicalize; `cli:` reserved; no CLI Assert).

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Home | Store mig 014 + domain Resolve/Decide. G19: no slug policy in CLI/MCP adapters |
| Mig | **`014_capability_tool_decision_enum.sql`** — rebuild + `CHECK (decision IN ('AUTO_ALLOWED','PENDING','ALLOWED','DENIED'))`; restore decision index; **do not** rewrite 001–013 |
| Heal (copy) | Invalid/empty/YOLO → **PENDING**. Never AUTO_ALLOWED. Never DROP |
| Resolve | Unknown persisted status → fail-closed (PENDING); **no** fall-through AUTO_ALLOWED upsert |
| Canonicalize | Decide + Resolve: exact builtin **Name** or `mcp:`+Name → `mcp:`+Name. Source: `BuiltinMCPCapabilitySpecs()`. No globs |
| `cli:` | Leave unchanged (S03) |
| Footgun fold | Unprefixed Name + existing `mcp:`+Name → one canonical row; DENIED > PENDING > ALLOWED > AUTO_ALLOWED |
| MCP Assert | **Unchanged** helper `mcp:`+Name |
| CLI | Prefer **zero** `cmd/trace` edits |
| Compat | Ceiling **14**; no **015+** |
| Forbidden | YOLO/AllowAll; new MCP tools; DF-77 CLI Assert; changing P15 Assert helper; daemon/HTTP; R2/R3/R4; S05; plan simulate; D21+ |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Store schema | `internal/store/schema/014_capability_tool_decision_enum.sql` | **New** — rebuild + CHECK + heal + slug fold |
| Store CRUD | `internal/store/capability_decision.go` | Reject non-enum on Upsert (defense) |
| Store tests | `internal/store/store_test.go` (or focused `*_test.go`) | CHECK reject; migrate heal; footgun fold; keeper versions include 14 |
| Domain | `internal/domain/capability_decision.go` | Canonicalize helper; Resolve `default` fail-closed |
| Domain tests | `internal/domain/capability_decision_test.go` | Named canonicalize + YOLO Resolve |
| MCP tests | `internal/mcp/mcp_test.go` | `TestMCPUnprefixedDecideGatesCallTool` + Assert keepers |
| MCP production | **Prefer zero** (`assert.go` unchanged) |
| CLI | **Prefer zero** |
| Compat | `evals/compat/compat_test.go` + `doc.go` | 13→14; allow 014; forbid 015+ |

## Named tests (required)

| Test | Intent |
|------|--------|
| `TestCapabilityToolDecisionCheckRejectsYOLO` | Post-014 INSERT/Upsert YOLO (and empty) errors; four enums OK |
| `TestCapabilityToolDecisionMigrateHealsYOLOToPending` | Pre-CHECK YOLO on builtin → 014 → PENDING, not dropped, not AUTO_ALLOWED |
| `TestResolveYOLOBuiltinDoesNotAutoAllow` | After heal, Resolve/Assert fail-closed; no AUTO_ALLOWED upsert |
| `TestDecideUnprefixedMCPNameCanonicalizes` | `trace_why` DENIED persists `mcp:trace_why`; Assert that slug fails |
| `TestCanonicalizeCustomAndCLISlugsUnchanged` | `tool:…` / `cli:add` unchanged; globs do not match |
| `TestMigrateUnprefixedDeniedFoldsOverAutoAllowed` | Dual rows → canonical DENIED |
| `TestMCPUnprefixedDecideGatesCallTool` | Unprefixed decide DENIED blocks CallTool `trace_why` |
| `TestMCPAssertDeniedBlocksCallTool` | Keeper |
| `TestMCPAssertBuiltinAutoAllowedSucceeds` | Keeper |
| `TestToolNamesRegistered` | Keeper — nine tools |
| `TestCapabilityDecisionAutoAllowBuiltinMCP` | Keeper |
| `TestCapabilityDecisionUnknownPendingFailClosed` | Keeper |
| `TestOpenCreatesDBAndMigratesIdempotent` | Keeper — version **14** |
| `TestMCPVirginProjectDoesNotMkdir` | S01 keeper |

TDD: add named tests first (red), then mig 014 + domain canonicalize/fail-closed (green).

## Role work
1. TDD named store CHECK/heal/fold + domain canonicalize + MCP unprefixed-decide tests (red on live 013).
2. Ship `014_capability_tool_decision_enum.sql`; Resolve default; canonicalize in Decide+Resolve; store Upsert reject; bump compat 14.
3. Prove green: named tests + locked verify cmds. Do **not** implement `cli:` Assert.
4. Board **status + Notes only** → next **P16-S02-02**.

## Locked verify commands

```text
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered|TestMCPAssert|TestMCPUnprefixed|TestMCPVirgin|TestOpenExisting|TestOpenCreates|TestBuiltinMCPCapabilitySpecs|TestCapabilityDecision|TestCapabilityToolDecision|TestResolveYOLO|TestDecideUnprefixed|TestCanonicalize|TestMigrateUnprefixed'

CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Product bar = `./cmd|internal|evals`. Compat ceiling **14**. Do **not** fail this row for R3 graphify space-in-path on full-module `./...` if present outside product pkgs. Do **not** fail for R4 CGO0 analyzer/CLI tree-sitter (keepers proven CGO1).

## Exit criteria
- [ ] DF-75: YOLO cannot persist; existing YOLO heals to PENDING; Resolve does not AUTO_ALLOW garbage
- [ ] DF-78: `decide --slug trace_why DENIED` gates `mcp:trace_why` / CallTool `trace_why`; custom/`cli:` unchanged; no globs
- [ ] Compat 14; no 015+; P15 Assert helper unchanged; S01 virgin keeper green
- [ ] Named tests pass; locked verify cmds PASS
- [ ] Board Notes → **P16-S02-02**

## Minimal todos
- [ ] Named tests (red → green)
- [ ] Mig 014 + Resolve fail-closed + canonicalize + fold
- [ ] Compat ceiling 14
- [ ] Locked verify suite
- [ ] Board status + Notes only
