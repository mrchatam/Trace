# P15-S01-00 — MCP Assert dispatch (FINAL)

## Metadata
- id: P15-S01-00
- todo_ids: [P15-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live `internal/mcp` + domain capability decisions, lock **FINAL** defaults for wiring `AssertToolAllowed` into MCP tool dispatch (residual **R1**). Thicken sibling `01`/`02`/SCOPE-TODOS. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md) — disposition matrix
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — FINAL; R1=fix
- Live: `internal/mcp/server.go`, tool handlers, `internal/domain/capability_decision.go`, `BuiltinMCPCapabilitySpecs`
- P14 S02 REVIEW (Assert ≠ MCP honesty)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (planner). Phase disposition R1=fix holds — grill only if slug mapping or G19 conflicts. **Unattended:** no architecture blockers; defaults below are FINAL.

## Live inventory (confirmed 2026-08-17)

| Area | Present? | Gap |
|------|----------|-----|
| `AssertToolAllowed` | Yes — `internal/domain/capability_decision.go` | **Zero** calls from `internal/mcp` (grep empty) |
| Builtin slugs | Yes — `BuiltinMCPCapabilitySpecs()` → `mcp:` + name for all nine | Assert must use **exact** same strings |
| MCP nine + version | Yes — `registerTools` + `RegisteredToolNames` / `TestToolNamesRegistered` | Must stay; no new tools |
| mig 013 decisions | Yes — store CRUD + Resolve/Decide | No new mig |
| install/decide MCP tools | Absent (good) | Do not add |
| `toolVersion` | No store open today | **Must** open project store (at least for Assert) so version is gated like peers |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Home | `internal/mcp` thin adapter only (G19). Call `domain.New(st).AssertToolAllowed` — **no** domain fork, **no** Assert logic in MCP |
| Slug | **`mcp:` + registered tool `Name`** — e.g. `mcp:trace_why`, `mcp:trace_version`. Must equal `BuiltinMCPCapabilitySpecs()[].Slug`. Helper: `"mcp:" + toolName` where `toolName` is the MCP `Tool.Name` / `RegisteredToolNames` entry |
| Where | Shared helper used at **entry of every** registered tool handler (`toolWhy`…`toolVersion`, including `reviewCreate`/`reviewSet`/`capability*` only via their parent entry points — **not** a second Assert per sub-action). Prefer `assertMCPToolAllowed(ctx, st, toolName)` after `openStore`; `toolVersion` must `openStore` (empty project → server defaultRoot/cwd) then Assert then return identity JSON |
| Fail-closed | `PENDING` / `DENIED` → return CallTool **error** (do not return success content). Map domain `ErrValidation` (or Assert error) to handler error |
| Builtins | First resolve still **AUTO_ALLOWED** (existing domain) — default nine-tool path stays green without human decide |
| DENIED proof | Plant via `DecideTool` DENIED on a **builtin** slug (e.g. `mcp:trace_why`), then CallTool that tool → error |
| Tests (named) | `TestMCPAssertDeniedBlocksCallTool`; `TestMCPAssertBuiltinAutoAllowedSucceeds`; keep `TestToolNamesRegistered`. Optional keeper: Assert helper coverage noted in REVIEW if every handler calls helper (grep `assertMCPToolAllowed` / equivalent) |
| Verify cmds | See sibling `01` locked verify block |
| Forbidden | New MCP tools; install/decide MCP; YOLO/AllowAll; changing ImpactWalk; R2/R3/R4 fixes; new mig; daemon/HTTP/embeddings; editing domain Assert semantics |
| Carry-forward | honesty A/B/C+G; Gates E/F/H; ablation; compat 13; p0x; x0; product pkgs `./cmd\|internal\|evals` |

## Owns
| Item | Intent |
|------|--------|
| MCP Assert wire-up | Every registered tool CallTool path fails closed via domain Assert with `mcp:<Name>` |
| Honesty close | S02 APPROVE ≠ “audit exists” alone — MCP path gated; closes P14 VERIFY residual R1 |

## Explicit deferrals
- R2 `allowContainsOut` late-upgrade
- R3 graphify space-in-path
- R4 CGO0 analyzers
- S05 / plan simulate / D21+

## Assumptions (unattended)
1. Slug convention is already product law via `BuiltinMCPCapabilitySpecs` — not renegotiated.
2. Opening store on `trace_version` solely to Assert is acceptable (thin adapter; no new tool).
3. Sub-handlers under review/capability are covered by parent tool entry Assert once per CallTool.
4. S02 VERIFY imports the named MCP Assert tests + carry-forward; DR-HANDOFF still default `no successor` (S02 planner owns exact checklist).

## Planner work
1. [x] Confirm live MCP registration + domain open path for Assert
2. [x] Lock FINAL defaults + named tests
3. [x] Thicken `01-mcp-assert-dispatch.md` + `02-scope-review.md` + SCOPE-TODOS
4. [x] Light S02 Depends note (named tests); board Notes → next **P15-S01-01**; mark this prompt **FINAL**

## Exit criteria
- [x] 00-PLANNER **FINAL**
- [x] 01/02 runnable with locked defaults
- [x] No product Go
- [x] Next board row **P15-S01-01**

## Minimal todos
- [x] Inventory MCP dispatch + Assert call sites
- [x] FINAL locks
- [x] Thicken 01/02
- [x] Board sync
