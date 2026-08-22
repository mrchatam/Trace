# P10 / S02 / 01 — MCP parity + install

## Metadata
- id: P10-S02-01
- todo_ids: [P10-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-21, DF-22, DF-32** per sibling **00-PLANNER** FINAL locks (2026-08-16). Thin MCP `trace_tasks` + `trace_capability` + `trace_version`; install/docs freshness; snake_case capability JSON on CLI (+ MCP). Inherit S01 why/Exact alias + IncludeWhy (do **not** re-implement DF-19/23/25/27/29). Keep carry-forward gates green. **No new migration. No plan/impact/index MCP.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-21/22/32
- [phase README](../../README.md)
- Live: `internal/mcp/{server,tools_*}.go`; `cmd/trace/{tasks,capability,install,help}.go`; `cmd/trace-mcp/main.go`; `internal/domain/capability.go` (`BuiltinMCPCapabilitySpecs`)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Locks are FINAL — do not re-debate.

## Locked defaults (FINAL — P10-S02-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | `internal/mcp` + thin `cmd/trace-mcp`; `cmd/trace` (capability encode + install tip/docs); `BuiltinMCPCapabilitySpecs` update |
| Migration | **None** |
| DF-21 `trace_tasks` | MCP list = CLI shape `[id,title,work_state,goal_id]`; optional `goal_id` / `project` |
| DF-21 `trace_capability` | One tool; `action` = declare\|list\|require\|unrequire\|missing; domain APIs only |
| DF-21 non-goals | No plan/impact/index MCP |
| DF-22 | README + `--write` stderr reload tip; **`trace_version`** `{ok,name,version}` |
| DF-32 | list/missing rows snake_case `id,kind,slug,title,status` (CLI + MCP) |
| Builtin specs | Six → include `mcp:trace_tasks`, `mcp:trace_capability`, `mcp:trace_version` (no auto-seed) |
| Tool count | **Nine** tools total |
| Carry-forward | honesty A/B/C + G; E/F; ablation; H; compat; p0x; x0; S01 retrieval tests; Gate C `dry_run:false` |
| Forbidden | New mig; daemon/HTTP/embeddings; process kill; plan/impact/index MCP; S01 re-litigation; Mode-B rewrite |

## Extension points (exact files)

| File | Work |
|------|------|
| `internal/mcp/server.go` | Register `trace_tasks`, `trace_capability`, `trace_version` |
| `internal/mcp/tools_*.go` (new or extend) | Handlers; G19 → domain/store only |
| `internal/mcp/mcp_test.go` + `export_test.go` | Tool registration + parity tests (tasks list; capability action; version; G19 import boundary) |
| `internal/domain/capability.go` | Extend `BuiltinMCPCapabilitySpecs` (+ test expected slugs) |
| `cmd/trace/capability.go` | DF-32 snake_case encode for list/missing |
| `cmd/trace/*_test.go` | Assert snake_case keys on capability list/missing stdout |
| `cmd/trace/install.go` | DF-22 stderr tip on successful `--write` |
| `cmd/trace/install_test.go` | Tip present on `--write` (stderr) |
| `cmd/trace/help.go` / `README.md` | Rebuild + abs `--bin` + reload MCP (DF-22); DF-05 unchanged |
| `cmd/trace-mcp/main.go` | Help text lists nine tools |

## Role work

1. TDD DF-32: capability list/missing CLI JSON must use snake_case (fails under current PascalCase).
2. Fix CLI encode (DTO/map or store json tags) — keep declare/require/unrequire envelopes.
3. Add `trace_tasks` MCP + test vs CLI shape / `ListTasks`.
4. Add `trace_capability` MCP (all five actions) + tests; G19 only.
5. Add `trace_version` + test; wire `BuiltinMCPCapabilitySpecs` (9 slugs).
6. DF-22 docs + install `--write` stderr tip + help/README.
7. Locked verify suite; board **status + Notes only** (cite test names + tool list).

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-checks: MCP tool names include `trace_tasks`/`trace_capability`/`trace_version`; `trace tasks` and `trace_tasks` same row keys; `capability list` JSON has `id` not `ID`; BuiltinMCP specs length/slugs match nine tools; README mentions reload after install.

## Exit criteria
- [ ] DF-21: MCP `trace_tasks` + `trace_capability` (five actions) live; G19; no plan/impact/index MCP
- [ ] DF-21: `BuiltinMCPCapabilitySpecs` includes new `mcp:trace_*` (no auto-seed)
- [ ] DF-22: docs + `--write` stderr tip; `trace_version` returns name+version
- [ ] DF-32: CLI (+ MCP) capability list/missing snake_case rows
- [ ] Nine tools registered; help text updated
- [ ] No new mig; carry-forward suite green; Gate C `dry_run:false` untouched
- [ ] Board Notes ready for **P10-S02-02**

## Out of scope
- S01 fidelity rework (already done)
- S03 index GC (DF-20)
- S04 operator/capability transition gates (DF-17/18/24/26/31)
- DF-28/30/33 deferred; daemon/HTTP; killing Cursor MCP processes

## Todo updates
Implementer: **status + notes only**. Record test names + DF checklist evidence. No spawning; no rewriting upcoming prompts.

## Minimal todos
- [ ] DF-32 CLI capability list/missing snake_case + test
- [ ] DF-21 `trace_tasks` MCP + test
- [ ] DF-21 `trace_capability` MCP (5 actions) + tests
- [ ] DF-22 `trace_version` + BuiltinMCP specs + docs/install tip
- [ ] Locked verify suite; board Notes
