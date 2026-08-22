# P16 / S01 / 01 — MCP project root / auto-init (FINAL locks from 00-PLANNER)

## Metadata
- id: P16-S01-01
- todo_ids: [P16-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-76** per sibling **00-PLANNER FINAL**. MCP CallTool must not auto-create `.trace/` / `trace.db` on a virgin (or db-less) project root. Bound-root DENIED stays fail-closed. Initialized other roots stay isolated. CLI `store.Open` mkdir **stays**. Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- Live: `internal/mcp/project.go` `openStore`, `internal/store/open.go`, `internal/mcp/mcp_test.go`, `cmd/trace/init.go`
- Hunt: `experiments/_bughunt/post-p15/{mcp-deny,mcp-noinit,mcp-fresh}/`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do **not** re-debate FINAL locks (OpenExisting; fail-closed missing store; CLI Open mkdir stays; per-store SoT HOLD; bind-to-defaultRoot rejected).

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Home | `store.OpenExisting` + MCP `openStore` (G19). No domain fork |
| API | `func OpenExisting(projectRoot string) (*Store, error)` — Abs only, no parent `.trace` walk-up (same as `Open`) |
| Exists | Regular file `<abs>/.trace/trace.db` must exist. Missing file **or** `.trace/` without db → `ErrNotInitialized` |
| Sentinel | `var ErrNotInitialized` in `internal/store` (next to `ErrLocked`). Wrap `%w` |
| Create | Must **not** create `.trace/` or `trace.db`. Never call `Open` on a Stat miss. After exists-check, reuse Open internals (lock/token/migrate/ensureProject). `trace.lock` on an existing store is OK |
| MCP | `openStore` → `store.OpenExisting`. Map `ErrNotInitialized` like `ErrLocked`: `fmt.Errorf("mcp: %w", err)` → CallTool **error** |
| Fail-closed | Virgin / missing db → error; no AUTO_ALLOWED DB; no entity writes |
| Isolation | DENIED on initialized A does not deny initialized B |
| Fallback | **None** — do not bind `project=` miss to `defaultRoot` |
| Assert | Unchanged helper + slug `mcp:`+Name; still after successful open |
| CLI | `store.Open` mkdir **unchanged** (`trace init`, add, why, …) |
| Forbidden | New MCP tools; session-global DENY; daemon/HTTP; YOLO; changing Assert slug; new mig; editing CLI Open mkdir; R2/R3/R4; S05; plan simulate; D21+ |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Store API | `internal/store/open.go` (+ sentinel next to lock/auth errors) | `OpenExisting` + `ErrNotInitialized`; skip mkdir/create when db missing |
| Store docs | `internal/store/doc.go` | Document `OpenExisting` vs `Open` |
| Store tests | `internal/store/store_test.go` | Named `TestOpenExisting*` |
| MCP adapter | `internal/mcp/project.go` | `openStore` uses `OpenExisting`; wrap sentinel |
| MCP tests | `internal/mcp/mcp_test.go` | Named virgin + isolation tests (use existing `callAdd` / `callVersion` + `AddInput.Project`) |
| CLI | **Prefer zero edits** | Keepers `TestInitCreatesDB` / `TestOpenCreatesDBAndMigratesIdempotent` |
| Domain | **Zero edits** | Assert unchanged |

## Named tests (required)

| Test | Intent |
|------|--------|
| `TestMCPVirginProjectDoesNotMkdir` | Virgin temp dir as `ProjectRoot`: `callAdd` + `callVersion` error; **no** `.trace/` created. Also: initialized defaultRoot + `AddInput.Project` = second virgin dir → error, second dir has no `.trace/`. Also: empty `.trace/` (no db) → error, still no `trace.db` |
| `TestMCPInitializedOtherRootIsolated` | A DENIED `mcp:trace_add`; B initialized; server bound to A; `callAdd` `Project=B` succeeds; unbound add (A) still DENIED |
| `TestMCPAssertDeniedBlocksCallTool` | Keeper — bound-root DENIED |
| `TestMCPAssertBuiltinAutoAllowedSucceeds` | Keeper — initialized AUTO_ALLOWED |
| `TestToolNamesRegistered` | Keeper — nine tools |
| `TestOpenExistingMissingReturnsErrNotInitialized` | Virgin dir: `errors.Is(err, ErrNotInitialized)`; no mkdir |
| `TestOpenExistingEmptyTraceDir` | `.trace/` without db: error; no `trace.db` |
| `TestOpenCreatesDBAndMigratesIdempotent` | Keeper — `Open` still creates |
| `TestInitCreatesDB` | Keeper — CLI init still creates |

TDD: add named tests first (red), then `OpenExisting` + `openStore` switch (green).

## Role work
1. TDD named store + MCP tests (red on live `Open` mkdir).
2. Add `ErrNotInitialized` + `OpenExisting`; switch MCP `openStore`; wrap sentinel.
3. Prove green: named tests + locked verify cmds.
4. Board **status + Notes only** → next **P16-S01-02**.

## Locked verify commands

```text
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered|TestMCPAssert|TestMCPVirgin|TestMCPInitialized|TestOpenExisting|TestOpenCreates|TestBuiltinMCPCapabilitySpecs|TestCapabilityDecision'

CGO_ENABLED=0 go test ./cmd/trace/... -count=1 -run 'TestInitCreatesDB|TestInitFailClosedWhenStoreLocked'

CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Product bar = `./cmd|internal|evals`. Compat ceiling remains **13** until S02. Do **not** fail this row for R3 graphify space-in-path on full-module `./...` if present outside product pkgs.

## Exit criteria
- [ ] DF-76: MCP CallTool on virgin / db-less dir fails closed (no `.trace/` or `trace.db` created)
- [ ] Named tests pass; DENIED keeper green; isolation HOLD proven
- [ ] CLI `Open` / `trace init` still mkdir; no new MCP tools; no new mig; G19 intact
- [ ] Locked verify cmds PASS
- [ ] Board Notes → **P16-S01-02**

## Minimal todos
- [ ] Named tests (red → green)
- [ ] `OpenExisting` + MCP `openStore` switch
- [ ] Locked verify suite
- [ ] Board status + Notes only
