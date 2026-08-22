# P15 / S01 / 02 — Scope review (MCP Assert dispatch) FINAL checklist

## Metadata
- id: P15-S01-02
- todo_ids: [P15-S01-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of S01 MCP Assert wire-up. Fresh subagent (≠ implementer). Compare claims + **00-PLANNER FINAL** + `01` to live code/tests. Spawn `P15-S01-02a`/`02b` for blocker/high. Prefer `REVIEW-NOTES.md`. Next **P15-S02-00** unless spawn.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-mcp-assert-dispatch.md](01-mcp-assert-dispatch.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [phase README](../../README.md)
- Live: `internal/mcp`, `internal/domain/capability_decision.go`, `BuiltinMCPCapabilitySpecs`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do not share the implementer’s session. Re-run verify; do not trust Notes alone.

## Checklist

| # | Check | How to evidence |
|---|--------|-----------------|
| 1 | Every registered tool entry calls Assert with slug `mcp:<Name>` (nine handlers incl. `toolVersion`) | Grep helper + each `tool*` entry; no missing CallTool path |
| 2 | Slug strings match `BuiltinMCPCapabilitySpecs` (`mcp:` + `RegisteredToolNames`) | Code + `TestBuiltinMCPCapabilitySpecs` |
| 3 | DENIED fail-closed on MCP CallTool | `TestMCPAssertDeniedBlocksCallTool` |
| 4 | Builtin / AUTO_ALLOWED still succeeds | `TestMCPAssertBuiltinAutoAllowedSucceeds` |
| 5 | Tool count still nine; no install/decide MCP | `TestToolNamesRegistered` |
| 6 | No new mig; G19 thin adapter; Assert logic stays in domain; no ImpactWalk edits | Diff |
| 7 | Sub-actions (review create/set, capability actions) not double-gated incorrectly and not **ungated** (parent entry Assert once) | Code review |
| 8 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + product pkgs | Re-run `01` locked verify cmds |
| 9 | R2/R3/R4 not falsely claimed fixed | Diff + Notes |

## Review procedure
1. Read FINAL locks in 00 + 01.
2. Map named tests → code; fresh verify cmds from `01`.
3. APPROVE (high, or medium with residuals listed) or spawn `P15-S01-02a`/`02b` with full prompts.
4. Write `REVIEW-NOTES.md` + board Notes; next **P15-S02-00** unless spawn.
5. If APPROVE: note S02 must import named Assert tests (already stubbed) — thicken S02 only if gaps found.

## Exit criteria
- [x] Checklist evidenced; confidence high (or medium with residuals listed)
- [x] Board status + Notes; next **P15-S02-00** (unless spawn)
- [x] No rewrite of done P15-S01-00/01 history

## Minimal todos
- [x] Independent verify + checklist
- [x] REVIEW-NOTES.md
- [x] Board sync
