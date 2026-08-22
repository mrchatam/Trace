| Order | ID | Role | Status note |
|------:|----|------|-------------|
| 1 | 00-PLANNER | scope planner | done — FINAL locks `trace install cursor` print/`--write` |
| 2 | 01-install-wire | implement | pending — CLI + README + ab-simple PROTOCOL |
| 3 | 02-scope-review | review | pending — after 01 |

## Locks (summary)

- **DF-03:** `trace install cursor` → stdout fragment or `--write` upsert `~/.cursor/mcp.json` (`mcpServers.trace`, backup, fail-closed)
- **DF-05:** docs — open **run folder** as workspace so `${workspaceFolder}` is not the Trace monorepo
- **Not** this scope: MCP list-tasks tool (S02 stays CLI `trace tasks`)
