# P22-S08-05 — Implement: useful in normal agent workflow

## Metadata
- id: P22-S08-05
- todo_ids: [P22-S08-05]
- role: implementer
- skills: [incremental-implementation, writing-for-agents]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Finish **C39**: Trace is usable in an agent’s normal engineering workflow (help, docs, MCP catalog complete vs CLI for P22 tools). Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md)
- Live: `cmd/trace/help.go`, `cmd/trace-mcp/main.go`, `CONTRIBUTING.md`, `README.md`

## Live baseline

| Present | Gap |
|---------|-----|
| Help: search, test run, verify, changes, knowledge, loop, git-hook | `TestHelpIncludesSearchTestVerify` absent |
| MCP **14** after S08-01 (`trace_loop`) | `trace-mcp -h` stale (missing regressions + loop) |
| CLI loop/test/verify/search/changes/regressions | CONTRIBUTING agent workflow paragraph missing |
| `TestHelpIncludesLoopNext` | Full workflow doc path for new agents |

## Locked defaults

| Item | Value |
|------|-------|
| Help strings | Must mention subcommands: **`search`**, **`test run`**, **`verify`**, **`changes`**, **`knowledge`**, **`install git-hook`**, **`loop`** (next/apply/status) |
| CONTRIBUTING | Add short **“Agent workflow (local)”** paragraph: `trace index` → `trace loop next --task <id>` → `trace test run --task <id>` → `trace loop apply` → `trace search` / `trace why` for history; note stdio MCP + rebuild |
| MCP `-h` | List **all** registered tools in order (14 post-S08-01); must match `RegisteredToolNames()` |
| Catalog parity | MCP mirrors P22 CLI query/history/loop tools; tests/outcomes evidence **CLI-only** (S05 lock) |
| Rebuild | Document + run in evidence: `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp`; verify `./bin/trace-mcp -h` |
| `trace_version` | Still required after MCP rebuild (P18 lesson) — smoke in Notes |
| Hosted MCP | **Forbidden** |
| Checklist | Box C39 when help + docs + catalog aligned |

## Requirements

1. **`TestHelpIncludesSearchTestVerify`** — assert help contains workflow subcommand strings (can extend `loop_test.go` or new help test file).
2. Fix **`cmd/trace-mcp/main.go`** `-h` to list all 14 tools (include `trace_regressions`, `trace_loop`).
3. **CONTRIBUTING.md** — agent workflow paragraph (5–8 lines, link to `docs/rules/agent-loop-protocol.md`).
4. Optional **README.md** one-liner under usage if missing — only if CONTRIBUTING alone is insufficient for discoverability.
5. **`TestToolNamesRegistered`** still expects **14**; no new MCP tools this row (S09 adds `trace_agents`).

## Touch files

- `cmd/trace/help.go`, `cmd/trace/loop_test.go` or `help_test.go`
- `cmd/trace-mcp/main.go`
- `CONTRIBUTING.md` (required)
- `README.md` (optional one-liner)

## Named tests

| Test | Proves |
|------|--------|
| `TestHelpIncludesSearchTestVerify` | help strings for agent workflow |
| `TestHelpIncludesLoopNext` | loop subcommands in help (keeper) |
| `TestToolNamesRegistered` | MCP catalog **14** matches registration order |

```bash
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/mcp/... -count=1 -run 'TestHelpIncludesSearchTestVerify|TestToolNamesRegistered|TestHelpIncludesLoopNext'
# Evidence (Notes): ./bin/trace-mcp -h | grep trace_loop
```

## Exit criteria

- [ ] C39 true (together with S08-01 loop half)
- [ ] Named tests PASS
- [ ] Live `trace-mcp -h` matches catalog
- [ ] CONTRIBUTING agent workflow paragraph present
- [ ] Board Notes: rebuild command + `-h` snippet

## Minimal todos

- [ ] Help tests + fix trace-mcp `-h`
- [ ] CONTRIBUTING paragraph
- [ ] Rebuild binary (evidence in Notes)
- [ ] Board notes

## Residual risks (carry to S08-06)

- S09 will add 15th tool — S08-06 must not claim full catalog final until S09-08
- README vs CONTRIBUTING duplication — keep single source in CONTRIBUTING
