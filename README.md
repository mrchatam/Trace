# Trace

Local-first project knowledge graph and progressive planning for AI coding agents.

Trace sits beside Git and your coding agent. It records **why** work exists, what the project currently believes, and how discoveries should change the next plan — then serves that back as bounded context (CLI, MCP, optional local GUI).

## Why Trace

Coding agents see files and diffs. They usually do not see:

- which goal or decision caused a change;
- what assumptions are still load-bearing;
- what evidence supports a “done” claim;
- how a discovery should replan affected work;
- which skills/tools/MCP servers a task actually needs.

Trace is that missing layer: a versioned SQLite graph under `.trace/`, with Git remaining the source-of-truth for history and content.

## Install

**Requirements:** Go matching [`go.mod`](go.mod) (currently **1.26+**), plus a C toolchain for the full CLI (`CGO_ENABLED=1` — tree-sitter analyzers).

The module is **not** published to the Go module proxy yet, so `go install …@latest` fails until a version is tagged. Build from a checkout:

```bash
git clone https://github.com/mrchatam/Trace.git
cd Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp
cp -f bin/trace bin/trace-mcp ~/.local/bin/   # or any dir on PATH
./bin/trace version
```

- Full `trace` binary needs **`CGO_ENABLED=1`**.
- `trace-mcp` builds with **`CGO_ENABLED=0`** in typical environments.
- Library packages that do not import analyzers remain usable with `CGO_ENABLED=0`.
- No GitHub Releases yet — use the build above. See also [`docs/gui-quickstart.md`](docs/gui-quickstart.md) for PATH notes.

`trace install …` configures agent/MCP/hook snippets; it does **not** put binaries on PATH.

## 60-second quickstart

From any project directory (Git optional but recommended):

```bash
trace init
# → creates .trace/trace.db (never project-root trace.db)

trace add goal --title "Harden auth" --body "Close session gaps"
# → {"id":"<goal-id>","ok":true,"type":"goal"}

trace add task --title "Audit cookie flags" --goal-id <goal-id>
# → {"id":"<task-id>","ok":true,"type":"task"}

trace index
trace search cookie
trace context <task-id>
trace loop status --task <task-id>
```

Clone of a repo that commits `trace/graph.json`:

```bash
trace init
trace seed import trace/graph.json
trace index
trace plan show
trace why goal <id>
trace context <task-id>
```

After `seed import`, tasks are **PENDING** (default export omits reviews, transitions, and task `work_state`).

## Surfaces

| Surface | Binary / command | Role |
|---------|------------------|------|
| **CLI** | `trace` | Canonical operator interface (`init`, `add`, `index`, `search`, `context`, `loop`, `plan`, …) |
| **MCP** | `trace-mcp` (stdio) | Same graph for Cursor/other MCP hosts — tools `trace_why`, `trace_context`, `trace_search`, `trace_loop`, `trace_plan`, `trace_explore`, … (17 tools; confirm with `trace_version`) |
| **HTTP / GUI** | `trace serve` / `trace gui` | **Opt-in** loopback API + embedded Explore SPA (default `127.0.0.1:7432`; not a daemon) |

Cursor MCP snippet (print only):

```bash
trace install cursor
# merge into ~/.cursor/mcp.json:
trace install cursor --write
# prefer --bin /abs/path/to/trace-mcp, then reload Cursor MCP
```

GUI (consumer repo needs only `.trace/`):

```bash
cd your-project && trace gui
# headless twin: trace serve
```

Details: [`docs/gui-quickstart.md`](docs/gui-quickstart.md).

## What it is / isn’t

**Is**

- Local-first knowledge + planning substrate for one project root (`-C` / cwd; no parent walk-up).
- Hybrid retrieval: exact lookup, lexical **FTS**, graph-label traversal, temporal/history — **not** embedding/vector semantic search (deferred — **DR-NOSSEM**; do not expect `semantic_match`).
- Progressive planning: coarse → deep near execution → replan from discoveries.
- Evidence and independent review before DONE; provenance on inferred vs verified facts.
- Portable causal graph via committed `trace/graph.json` (`.trace/` stays local/gitignored).

**Isn’t**

- A Git replacement, IDE, or coding-agent runtime.
- Cloud SaaS, multi-tenant hosting, or always-on network daemon (HTTP/GUI is opt-in loopback).
- Embedding/RAG semantic search (not shipped).
- A claim that Jev or other external integrations are productized here.

## Design principles (condensed)

- Progressive planning; discovery is normal and feeds the plan.
- Evidence over assertions; implementation ≠ review identity.
- Context minimization — smallest high-value packet first, expand on demand.
- Hybrid FTS/graph/temporal retrieval (no embeddings yet — DR-NOSSEM).
- Provenance everywhere; forward progression; human authority.
- Capability-aware tasks; Git for history; incremental index updates.

## Feature map (shipped vs deeper docs)

| Area | Shipped surface (examples) |
|------|----------------------------|
| Graph CRUD | `trace add`, `link`, `transition`, `review`, `tasks` |
| Retrieval | `search`, `why`, `context`, `explore` / MCP `trace_explore` |
| Loop / gates | `loop next\|apply\|status\|gate`, optional `--enforce` |
| Planning | `plan create-coarse\|deep\|show\|bootstrap\|…` |
| Index | `index` / `index status` / `index watch` (file-local incremental) |
| Portable graph | `seed import\|export` → `trace/graph.json` |
| Impact / caps | `impact`, `capability`, `agents recommend` |
| Install adapters | `install cursor\|claude\|cursor-hook\|git-hook\|agents` |
| Local HTTP/GUI | `serve`, `gui` (embedded SPA) |

For command truth, prefer `trace --help` and MCP `Instructions` / `trace_version` over older blog-style docs.

## Documentation map

| Doc | What |
|-----|------|
| [`docs/TODO.md`](docs/TODO.md) | Execution board index (phase tables under `docs/TODO/`) |
| [`docs/gui-quickstart.md`](docs/gui-quickstart.md) | Opt-in GUI / serve |
| [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) | Component boundaries |
| [`docs/SECURITY.md`](docs/SECURITY.md) | Trust boundaries |
| [`docs/PROJECT_MODEL.md`](docs/PROJECT_MODEL.md) | Entities, relations, provenance |
| [`docs/RETRIEVAL_AND_CONTEXT.md`](docs/RETRIEVAL_AND_CONTEXT.md) | Retrieval strategy (incl. DR-NOSSEM) |
| [`docs/ROADMAP.md`](docs/ROADMAP.md) | Roadmap / milestone gates |
| [`PROJECT_DOCS_INDEX.md`](PROJECT_DOCS_INDEX.md) | Full doc index |
| [`CONTRIBUTING.md`](CONTRIBUTING.md) | Dev rules, portable-graph export, dual-stack notes |
| [`AGENTS.md`](AGENTS.md) | Agent entrypoint / phase focus |
| [`LICENSE`](LICENSE) | Apache-2.0 |

## Status

Active research/engineering codebase: foundational graph, CLI, MCP, and opt-in local GUI are in use and under rapid iteration (multi-phase boards; see `AGENTS.md`). The validation bar is not “orchestrate many agents,” but:

> Can an agent understand an unfamiliar repo, plan a bounded task, adapt when implementation discovers something new, and decide better with the project graph than with raw files alone?

Expect sharp edges, evolving commands, and honest gaps. Read `trace --help` and the active phase board before assuming a feature exists.

## License

Apache-2.0. Hosted/enterprise offerings may appear later without restricting the open-source core.
