# Project Graph for AI Agents

> A versioned, evolving knowledge graph of a software project that connects goals, decisions, assumptions, plans, tasks, code, history, discoveries, evidence, and agent capabilities.

The project is designed for AI-assisted software development. Its purpose is not to replace Git, an IDE, or a coding agent. It adds the layer those systems generally do not provide:

- why code exists;
- which goal, decision, or task caused it to exist;
- how a change affects other work;
- what the project currently believes to be true;
- how implementation discoveries should change future planning;
- what evidence supports a completion claim;
- which skills, tools, MCP servers, hooks, and rules are needed for a task.

The system treats planning as **progressive refinement**, not as exhaustive up-front prediction.

## Core model

```text
Goal
  ↓
coarse project plan
  ↓
phase
  ↓
scope plan
  ↓
minimal tasks
  ↓
implementation
  ↓
discovery / review / evidence
  ↓
graph update
  ↓
replan affected future work
  ↓
continue
```

The graph combines five logical layers:

1. **Code graph** — files, symbols, imports, calls, dependencies.
2. **History graph** — commits, changes, project states, provenance.
3. **Work graph** — goals, phases, scopes, tasks, reviews.
4. **Causal/decision graph** — decisions, assumptions, discoveries, rationale, impact.
5. **Environment graph** — skills, rules, tools, MCP servers, hooks, agents, models.

The project uses Git as the canonical version-history substrate rather than reimplementing Git.

## Design principles

- **Progressive planning:** detail is created when it becomes useful.
- **Discovery is normal:** implementation gaps are expected and become planning inputs.
- **Evidence over assertions:** an agent saying “done” is never sufficient evidence.
- **Independent review:** implementation and review are separate contexts/identities.
- **Context minimization:** agents receive the smallest high-value context first and can expand it on demand.
- **Hybrid retrieval:** exact lookup, lexical search, semantic search, graph traversal, and temporal history work together.
- **Provenance everywhere:** inferred facts are never silently treated as verified facts.
- **Forward progression:** backward movement must be explicit; reversals are recorded as new states.
- **Decision awareness:** user decisions are first-class objects with impact analysis and alternative routes.
- **Capability-aware planning:** tasks can select the skills, rules, tools, MCPs, and hooks needed for their environment.
- **Git delegation:** source history and content remain in Git; the project graph stores meaning and references.
- **Incremental updates:** changing one area must not trigger a complete project re-analysis.
- **Human authority:** the system warns and advises; it does not silently override user decisions.

## Repository documentation

- `docs/TODO.md` — **execution board index**; row tables in `docs/TODO/phase-NN.md`.
- `docs/gui-quickstart.md` — **opt-in** `trace gui` + browser GUI (local-first; loopback default; `serve` for headless/scripting).
- `docs/rules/` — agent loop protocol, project rules, skills map.
- `docs/phases/` — runnable prompts.
- `docs/ROADMAP.md` — implementation roadmap and milestone gates.
- `docs/ARCHITECTURE.md` — system architecture and component boundaries.
- `docs/PROJECT_MODEL.md` — entities, relations, provenance, and state model.
- `docs/PLANNING.md` — progressive planning algorithm.
- `docs/RETRIEVAL_AND_CONTEXT.md` — retrieval, context compilation, RAG/graph retrieval strategy.
- `docs/REVIEW_AND_VERIFICATION.md` — multi-layer review, evidence, and anti-hallucination design.
- `docs/DECISION_IMPACT.md` — user decisions, impact analysis, alternatives, and hypothetical plans.
- `docs/AGENT_ENVIRONMENT.md` — skills, rules, MCPs, hooks, tools, agents, and capability-aware planning.
- `docs/STORAGE_AND_PERFORMANCE.md` — storage, indexing, caching, scaling, and large-repository strategy.
- `docs/SECURITY.md` — trust boundaries and safety model.
- `docs/EVALUATION.md` — benchmarks and research questions.
- `docs/init/` — initialization planning registers (decisions, P0-X bar, laws).
- `CONTRIBUTING.md` — contribution model and development rules.
- `LICENSE` — Apache-2.0.
- `AGENTS.md` — agent entrypoint.

## Status

This repository is intentionally starting with the foundational knowledge and planning layer before attempting a large multi-agent control plane. The first validation target is not “can it orchestrate many agents?” but:

> Can an agent understand an unfamiliar repository, plan a bounded task, adapt when implementation discovers something new, and make better decisions using the project graph than using raw repository contents alone?

## Build

```bash
CGO_ENABLED=1 go test ./...
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
./bin/trace version
```

Requires Go 1.24+ (`go.mod` / `modernc.org/sqlite`). The full `trace` binary links tree-sitter analyzers and needs **`CGO_ENABLED=1`**. Library packages that do not import `analyzers` remain usable with `CGO_ENABLED=0`.

Also build the MCP stdio server when using Cursor:

```bash
CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp
```

## Install / Cursor MCP

Print a merge-ready MCP snippet (no file write):

```bash
trace install cursor
```

Upsert `mcpServers.trace` into `~/.cursor/mcp.json` (creates a `*.bak.<UTC>` backup of an existing file; path on stderr):

```bash
trace install cursor --write
# optional: --bin /abs/path/to/trace-mcp  --mcp-json /path/to/mcp.json
```

The entry uses `type=stdio`, command `trace-mcp` (or `--bin`), and `args: ["-C", "${workspaceFolder}"]`. Open the **project (or experiment run folder)** as the Cursor workspace so `${workspaceFolder}` points at the seeded tree — not the Trace monorepo root by mistake (DF-05).

After a successful `trace install cursor` (print or `--write`) or after rebuilding `trace-mcp`, prefer an **absolute `--bin` path**, then **reload/restart Cursor MCP** (or reload the window) so the long-lived stdio process is not stale. The live Cursor tool catalog may lag until reload (DF-37). Use MCP `trace_version` to confirm the live process identity.

## Portable graph (clone recipe)

After cloning a repo that commits `trace/graph.json`:

```bash
trace init
trace seed import trace/graph.json
trace index
trace plan show
trace why goal <id>
trace context <task-id>
```

After `seed import`, tasks are **PENDING** (default export omits reviews, transitions, and task `work_state`; live DONE/SKIPPED stays on the exporter’s `.trace/`).

`trace index` rebuilds the derived code graph locally; causal and plan data come from the git-committed JSON. See [CONTRIBUTING.md](CONTRIBUTING.md) for export-before-PR and merge conventions.

## License

The core project is intended to be released under Apache-2.0. The project may later offer hosted, enterprise, support, or other commercial services without restricting the open-source core.
