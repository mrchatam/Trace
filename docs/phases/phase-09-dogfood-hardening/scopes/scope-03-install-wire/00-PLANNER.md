# P09-S03-00 — Install & Cursor wire (DF-03/DF-05)

## Metadata
- id: P09-S03-00
- todo_ids: [P09-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S03 implement/review prompts for a Graphify-like **Cursor MCP install path**: print or merge config for `trace-mcp` with `-C ${workspaceFolder}`, and document the workspace-root footgun (DF-05). **No product Go in this row.**

## Depends-on (S02 APPROVE)

S02 shipped **CLI** `trace tasks` (+ DF-04 seed resolve-vs-`-C`). Install-wire must **not** invent an MCP list-tasks / `trace_tasks` tool — keep discoverability CLI-primary; this scope owns **editor config wire only**.

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `cmd/trace/root.go` | No `install` (or `mcp`) command |
| `cmd/trace/help.go` | No Cursor / MCP install docs |
| `cmd/trace-mcp` | Stdio server; `-C`/`--project`; tools why/context/add/link/transition/review |
| `internal/mcp` | Server name `"trace"` |
| README | Build only — **no** MCP install section |
| Live `~/.cursor/mcp.json` | Server key `trace`; `command` abs path to `trace-mcp`; `args: ["-C","${workspaceFolder}"]` |
| `experiments/ab-simple/PROTOCOL.md` | Manual install notes + DF-05 workspace warning (canonical) |
| Dogfood | DF-03 (no install writer); DF-05 (open run folder, not Trace monorepo) |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Command | **`trace install cursor`** — not `trace mcp snippet` (single Graphify-like verb) |
| Subcommand shape | `trace install <target>`; only **`cursor`** required this scope (`install` alone → usage) |
| Default mode | **Print** merge-ready JSON fragment to stdout: `{"mcpServers":{"trace":{…}}}` (pretty); exit 0; **no** file write |
| Write mode | **`--write`** → upsert into MCP config file (default `$HOME/.cursor/mcp.json`) |
| Backup | Before any overwrite of an existing file: copy to `$HOME/.cursor/mcp.json.bak.<UTC-YYYYMMDDTHHMMSSZ>` (or sibling of `--mcp-json`); print backup path on stderr |
| Merge semantics | Parse existing JSON object; ensure `mcpServers` map; **upsert only** key `trace`; leave other servers untouched; create parent `.cursor/` + file if missing |
| Fail-closed | Existing `--mcp-json` / default path that is non-empty but **invalid JSON** → exit 2, no write (no `--force` this scope) |
| Idempotent | Re-run `--write` with same snippet OK (upsert replaces `trace` entry) |
| Server key | **`trace`** (match `internal/mcp` + live dogfood mcp.json) |
| Entry fields | `"type":"stdio"`, `"command":…`, `"args":["-C","${workspaceFolder}"]` — literal `${workspaceFolder}` string (Cursor expands) |
| Command binary | Default **`trace-mcp`** (PATH); optional **`--bin <path>`** for abs path (dogfood machines) |
| Test override | **`--mcp-json <path>`** — write/print target path for tests (default still `$HOME/.cursor/mcp.json`) |
| Docs | **README** short Install / Cursor MCP section; **`experiments/ab-simple/PROTOCOL.md`** — point at `trace install cursor` + DF-05 footgun (open **run folder** as workspace, not Trace monorepo) |
| Help | `help.go` line for `install cursor [--write] [--bin] [--mcp-json]` |
| Packages | Thin **`cmd/trace` only** (G19). No `internal/mcp` tool changes; no new store/mig |
| MCP tools | **Out of scope** — do not add list-tasks / `trace_tasks` |
| Carry-forward | honesty A/B/C + Gate G; p0x; x0; S01 Why/context-with-review; S02 `trace tasks`; `CGO_ENABLED=1 go test ./...` |
| Forbidden | Daemon/HTTP/embeddings; new MCP tools; new `011_*` mig; rewriting Phase 00–08 / S01–S02 history; weakening S02 CLI discoverability |

## Exit
- [x] Thicken `01-install-wire.md` + `02-scope-review.md`
- [x] SCOPE-TODOS + board Notes; next **P09-S03-01**
- [x] Product Go — **not** this row
