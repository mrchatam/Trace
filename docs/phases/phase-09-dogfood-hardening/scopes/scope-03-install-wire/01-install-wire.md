# P09 / S03 / 01 — Install & Cursor wire (DF-03/DF-05)

## Metadata
- id: P09-S03-01
- todo_ids: [P09-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Ship **`trace install cursor`** so agents/humans can print or merge Cursor MCP config for `trace-mcp -C ${workspaceFolder}` (DF-03), and document the workspace-root footgun (DF-05). Keep carry-forward gates green. **Do not** add MCP list-tasks tools.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL 2026-08-16
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-03 / DF-05
- Live: `cmd/trace/{root,help}.go`; `cmd/trace-mcp/main.go`; `internal/mcp` server name `trace`; `experiments/ab-simple/PROTOCOL.md`; README Build section
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute.

## Depends-on (S02 APPROVE)

S02 shipped **CLI** `trace tasks` (+ DF-04 seed resolve-vs-`-C`). This row owns Cursor MCP install/wire only — **do not** add an MCP `trace_tasks` / list-tasks tool unless a later board row promotes it.

## Locked defaults (FINAL — P09-S03-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Command | `trace install cursor` |
| Default | Print `{"mcpServers":{"trace":{…}}}` pretty JSON to stdout; no file write |
| Write | `--write` upserts `mcpServers.trace` into `$HOME/.cursor/mcp.json` (or `--mcp-json`) |
| Backup | Existing file → `*.bak.<UTC>` before overwrite; path on stderr |
| Entry | `type=stdio`, `command` = `trace-mcp` or `--bin`, `args=["-C","${workspaceFolder}"]` |
| Fail-closed | Invalid existing JSON → exit 2, no write |
| Docs | README Install/MCP + `experiments/ab-simple/PROTOCOL.md` DF-05 |
| Packages | `cmd/trace` only (+ tests). **No** new MCP tools / daemon / HTTP / mig |
| Carry-forward | honesty A/B/C + Gate G; p0x; x0; S01 review Why/context; S02 tasks; `./...` |
| Forbidden | MCP list-tasks; daemon/HTTP; board spawn rights; S01/S02 regress |

## Extension points (exact files)

| File | Work |
|------|------|
| `cmd/trace/install.go` (new) | `cmdInstall`: parse `cursor`, flags, build snippet, print or merge+backup |
| `cmd/trace/root.go` | Dispatch `case "install":` |
| `cmd/trace/help.go` | Document `install cursor [--write] [--bin path] [--mcp-json path]` |
| `cmd/trace/install_test.go` or `cli_test.go` | Print shape; `--write` to temp `--mcp-json`; backup created; other servers preserved; invalid JSON fails |
| `README.md` | Short **Install / Cursor MCP** section after Build |
| `experiments/ab-simple/PROTOCOL.md` | Replace hand-edit notes with `trace install cursor` (+ `--write`); keep DF-05 run-folder warning |

## Role work

1. TDD: failing test — `trace install cursor` stdout contains `mcpServers.trace` with `-C` / `${workspaceFolder}`.
2. Implement print path; then `--write` merge + backup against `--mcp-json` temp file (preserve sibling servers).
3. Wire help + README + ab-simple PROTOCOL (DF-05).
4. Run locked verify suite; board **status + Notes only**.

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-check:

```bash
trace install cursor
# → pretty JSON fragment with mcpServers.trace.args ["-C","${workspaceFolder}"]

trace install cursor --write --mcp-json /tmp/trace-mcp-test.json
# → file upserted; if prior existed, .bak.* created
```

## Exit criteria
- [ ] `trace install cursor` prints merge-ready `mcpServers.trace` snippet (stdio + `-C ${workspaceFolder}`)
- [ ] `--write` upserts only `trace`; backup before overwrite; invalid JSON fail-closed
- [ ] `--bin` / `--mcp-json` honored; default command `trace-mcp`
- [ ] Help + README + `experiments/ab-simple/PROTOCOL.md` document install + DF-05 footgun
- [ ] **No** new MCP tools / daemon / HTTP / migration
- [ ] Carry-forward green (S01 review Why/context; S02 tasks)
- [ ] Board Notes ready for **P09-S03-02**

## Out of scope
- New MCP `trace_tasks` / list-tasks tool (stay with S02 CLI)
- Installing binaries to PATH (`go install` / copy) — document only if already in README Build
- Other editor targets (VS Code, Claude Desktop) — Cursor only
- Rewriting Gate packs or Phase 00–08 / S01–S02 history

## Todo updates
Implementer: own row status + Notes only. Do not rewrite planner locks or spawn board rows.

## Minimal todos
- [ ] TDD print snippet + merge/backup tests
- [ ] Implement `trace install cursor` (+ flags); help dispatch
- [ ] README + ab-simple PROTOCOL (DF-03/DF-05)
- [ ] Run locked verify suite; mark P09-S03-01 done with Notes
