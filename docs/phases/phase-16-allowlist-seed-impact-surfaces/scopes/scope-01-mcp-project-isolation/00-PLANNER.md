# P16-S01-00 — MCP project isolation (stub — thicken vs live)

## Metadata
- id: P16-S01-00
- todo_ids: [P16-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live `internal/mcp` + `store.Open`, lock **FINAL** defaults for **DF-76**: MCP `project=` must not auto-init a virgin `.trace/` that AUTO_ALLOWs builtins and bypasses DENIED on another root. Thicken sibling `01`/`02`/SCOPE-TODOS. **No product Go.**

## References
- [phase README](../../README.md) — P16-00 FINAL
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- Live: `internal/mcp/project.go` (`openStore` → `store.Open`), `internal/store/open.go` (`MkdirAll`)
- Hunt: `experiments/_bughunt/post-p15/` `mcp-deny/` vs `mcp-noinit/` / `mcp-fresh/`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Phase lock DF-76=fix holds — grill only if exist-check vs workspace-bind trades conflict with Cursor `project=` workspaceFolder.

## Inherited locks (phase — do not reopen)
- MCP must **not** auto-mkdir `.trace/` for a path with no existing store
- Switching `project=` among **already initialized** roots remains valid (per-store SoT)
- Do **not** inherit DENIED across stores
- CLI `store.Open` auto-init **unchanged** (operator explicit)
- No new MCP tools this scope; keep nine + `trace_version`
- Named tests (phase hint; confirm names here): `TestMCPVirginProjectDoesNotCreateStore`; `TestMCPProjectOverrideDeniedDoesNotEscapeViaFreshRoot` (or equivalent)

## Planner work
1. [ ] Confirm live `openStore` → `store.Open` mkdir path
2. [ ] Lock FINAL API (OpenExisting / exist-check in MCP only — prefer **not** changing CLI Open)
3. [ ] Thicken 01/02/SCOPE-TODOS; mark this prompt **FINAL**; next **P16-S01-01**

## Exit criteria
- [ ] 00-PLANNER **FINAL**
- [ ] 01/02 runnable with named tests
- [ ] No product Go
