# P16-S02-00 — Allowlist CHECK + CLI parity (stub — thicken vs live)

## Metadata
- id: P16-S02-00
- todo_ids: [P16-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock **FINAL** defaults for **DF-75** (CHECK + YOLO fail-closed), **DF-77** (CLI Assert parity), **DF-78** (unprefixed slug). Thicken 01/02. **No product Go.**

## References
- [phase README](../../README.md)
- Live: `internal/domain/capability_decision.go` (`ResolveToolDecision` fall-through), store mig 013, `cmd/trace` add/why, `capability decide`
- Hunt: `cap-decisions/`, `mcp-yolo/`, `mcp-cli/`, `mcp-footgun/`

## Inherited locks (phase)
- CHECK `decision IN (AUTO_ALLOWED, PENDING, ALLOWED, DENIED)`; expect mig **014**; bump compat ceiling
- Resolve: unknown/garbage **must not** upsert AUTO_ALLOWED (fail-closed)
- CLI MCP-equivalent verbs call `AssertToolAllowed` with `mcp:<tool>` (G19). **Ungated:** `capability decide` / `decisions`
- DF-78: normalize exact registered names to `mcp:` **or** reject with hint; DENIED unprefixed must gate MCP `trace_why`
- No YOLO/AllowAll; no install/decide MCP tools; no DF-76 reopen (S01)

## Named tests (phase hint)
- `TestCapabilityDecisionCHECKRejectsYOLO`
- `TestResolveGarbageDecisionDoesNotAutoAllow`
- `TestCLIHonorsDeniedMCPSlug`
- `TestDecideUnprefixedBuiltinSlugGatesMCP`

## Planner work
1. [ ] Confirm live CHECK absence + Resolve fall-through
2. [ ] Lock FINAL mig number + CLI call sites
3. [ ] Thicken 01/02; **FINAL**; next **P16-S02-01**
