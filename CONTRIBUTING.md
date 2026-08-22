# Contributing

## Project direction

The project exists to explore a persistent project knowledge graph and progressive planning system for AI-assisted software engineering.

Good contributions should improve one or more of:

- project understanding;
- graph quality;
- provenance;
- retrieval;
- context efficiency;
- progressive planning;
- discovery propagation;
- decision impact analysis;
- verification;
- environment/capability modeling;
- large-repository performance.

## Before contributing

Read:

- `README.md`
- `docs/ROADMAP.md`
- `docs/PROJECT_MODEL.md`
- `docs/ARCHITECTURE.md`

## Design preference

Prefer:

- deterministic facts over inferred facts;
- incremental computation over full recomputation;
- explicit provenance;
- small composable interfaces;
- adapters over hard-coded integrations;
- measurable performance;
- benchmark evidence over intuition.

Avoid:

- giant generic ontologies;
- unnecessary LLM calls;
- full-project context dumps;
- vendor-specific core assumptions;
- silently mutating project state from an agent claim.

## Changes to the data model

Data-model changes must include:

- migration strategy;
- provenance impact;
- query impact;
- performance considerations;
- benchmark/fixture updates.

## Agent integrations

New agent integrations should use the canonical API and must not duplicate planner state logic.

## Agent workflow (local)

On a Trace-enabled repo, a typical engineering loop is: **`trace index`** (or `trace install git-hook --write` for incremental updates) → **`trace loop next --task <id>`** for a bounded packet → implement → **`trace test run --task <id>`** → **`trace loop apply`** to persist discoveries → **`trace search`** / **`trace why`** for history and causality. Follow the session protocol in [`docs/rules/agent-loop-protocol.md`](docs/rules/agent-loop-protocol.md).

**Moat-first orient (MCP + CLI):** Pick work with **`trace_tasks`**, load bounded context with **`trace_context`** (optional **`query`** for agent FTS merge), run the deliberation loop via **`trace_loop`** (`next`/`status`/`apply`), **`trace_loop action=gate`** (or `trace loop gate`) before product edits, **`trace_review`** before DONE, and **`trace_plan`** when a goal lacks a plan tree. For task-scoped discovery, compose reads in order: **`trace_search`** → **`trace_why`** → **`trace_impact`** → **`trace_capability`** — progressive caps, never a full graph dump. The Trace MCP server exposes this playbook in its **`Instructions`** field at connect time.

**Graph-first GUI (local):** From **any project with `.trace/`** (no Trace checkout required): **`cd your-project && trace gui`** (or **`trace serve`**). The Explore SPA is **embedded in the `trace` binary**; consumer repos do not need `web/`. Contributors serving from the Trace repo may use disk `web/dist` when present. Optional **`-C DIR`** before or after `gui`/`serve` to pick a project root. Run **`trace serve`** / **`trace gui`** and open **`/`** (nav **Explore**) for a seed-composed project graph. A first-visit orient panel explains the moat-first loop (Tasks → Loop → gate → review) and Laws 6–7 budget caps — dismissible via `localStorage` (`trace.orient.dismissed`). The browser GUI is a thin adapter (Law 19) over the same HTTP API as CLI/MCP; no parallel graph logic in `web/`.

**Harness vs product enforcement:** Trace product gates (`trace loop gate`, `--enforce` on DONE/export, status `violations[]`) read SQLite evidence and are authoritative. Harness install (`trace install cursor|claude|cursor-hook`) writes Cursor/Claude rules and optional pre-edit hooks that **call** those CLIs — best-effort reminders, not a second policy engine. Default `.trace/config.json` enforce mode is off; teams opt in explicitly. See `docs/phases/phase-23-enforcement-choke-points/ENFORCEMENT.md`.

Query, history, and loop tools are also available via the **stdio MCP** server (`trace-mcp`). After rebuilding `trace-mcp`, **reload Cursor MCP** (or restart the window) and confirm the live process with **`trace_version`**. A partial tool list (e.g. **9/17**) means a **stale stdio process** — not an intentional tool reduction; rebuild, reload, then call **`trace_version`**. Offline rebuild: `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp` (use `CGO_ENABLED=1` when tree-sitter deps require CGO in your environment). For symbol-level code exploration alongside Trace, see [Trace + Codegraph (optional dual-stack)](#trace--codegraph-optional-dual-stack).

## Trace + Codegraph (optional dual-stack)

Trace and Codegraph are **complementary local-first tools** with different jobs — not a merged product. Trace owns task loop, evidence, and progressive planning; Codegraph (when installed) owns symbol-level code graph exploration. Use one, both, or neither per repo.

### Storage (separate indexes)

| Stack | Local store | Contents |
|-------|-------------|----------|
| **Trace** | `<projectRoot>/.trace/trace.db` | Task backlog, plan tree, evidence, gates, reviews, transitions |
| **Codegraph** | `<projectRoot>/.codegraph/` | Symbol graph, call paths, blast-radius summaries |

Both directories are **gitignored**. They use **separate SQLite indexes** — no shared database, no cross-index reads.

### When to use Trace

Trace is **moat-first** for directed engineering work. See **Moat-first orient** above for the full playbook. Typical surfaces:

- **`trace_tasks`** / **`trace loop next`** — pick and run bounded work packets
- **`trace_context`** (+ optional **`query`**) — progressive task-scoped context
- **`trace_loop`** — deliberation loop; **`trace_loop action=gate`** before product edits
- **`trace_review`** / **`trace_transition`** — evidence-backed completion
- **`trace_plan`** — bootstrap or extend plan tree when a goal lacks structure
- **Portable graph** — `trace/graph.json` for causal entities and plan tree (see **Portable graph** below)

### When to use Codegraph (optional)

Symbol-level exploration is **optional — only when a per-repo symbol graph helps**:

- Run **`codegraph init`** in a project to build `.codegraph/` (user choice; Trace does not require it)
- Use **`codegraph_explore`** MCP (separate harness config) for call paths, symbol source, and blast-radius summaries
- Index is **per project** — pass `projectPath` when querying repos other than the harness cwd

### Setup

Order is independent — neither stack requires the other:

1. **Trace path:** `trace init` → `trace index [paths…]` → configure **`trace-mcp`** stdio server in your harness (see **Agent workflow (local)** above).
2. **Optional Codegraph path:** `codegraph init` in target repo → register **Codegraph MCP** separately in harness config. Reload MCP after rebuilds (same hygiene as Trace **9/17** note above).

### Law 19 — adapter boundaries

Each stack is an **adapter/MCP over its own store** (Law 19). Trace CLI and **`trace-mcp`** never open, read, or write **`.codegraph/`**. Codegraph MCP is a **separate** harness registration — not bundled inside Trace product code. HTTP/GUI adapters follow the same rule: canonical Go API only, no second source of truth.

### Not shipping (product rejects)

Trace **does not** ship:

- **Mandatory dual-index** — both `.trace/` and `.codegraph/` are opt-in per repo
- **Bundled trace+Codegraph MCP** — one stdio server exposing both stacks
- **Trace core indexing Codegraph data** — no symbol graph ingestion into `.trace/trace.db`

### Investigation context

Phase 38 peer investigation documents the complement model:

- [PEER-CG §5 — Trace strengths peers lack](docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-02-codegraph-peer/PEER-CG.md) — moat vs read-only graph
- [PEER-FIXTURES](docs/phases/phase-38-retrieval-context-peer-gaps/PEER-FIXTURES.md) — fixture scenarios used in peer comparison

## Review

Contributors should include:

- tests;
- evidence for behavior;
- benchmark impact where relevant;
- notes on compatibility.

## Portable graph (git)

Trace’s **live** SQLite store is **`<projectRoot>/.trace/trace.db`** only (`.trace/` gitignored). A file named `trace.db` at the project root is **not** the Trace store — do not open or create it (agents: use CLI/MCP). Prefer ignoring root-only `/trace.db` in consumer `.gitignore` without un-ignoring `.trace/`. When `<root>/trace.db` is a regular file, Trace prints a non-fatal warn **once per `openStore`**; each CLI/MCP/HTTP open may re-emit that warn (intentional — no persistent suppress flag). The **portable semantic graph** — causal entities, links, and plan tree — lives at **`trace/graph.json`** at the repo root.

1. **Path:** commit **`trace/graph.json`** at repo root (not `.trace/`).
2. **Before PR:** run `trace seed export -o trace/graph.json` when entity/plan graph changed.
3. **Clone recipe:** `trace init` → `trace seed import trace/graph.json` → `trace index [paths…]` → use `trace why`, `trace context`, `trace plan show` offline.
4. **Evidence:** git **author + commit SHA** and JSON **`exported_at_commit`** are snapshot **evidence**, not entity identity (UUIDs unchanged). **`transition.actor` / review actor / `as_operator` are not authentication** (same as DF-44).
5. **Merge:** parallel PRs may conflict on `trace/graph.json` — resolve in git manually; **no** custom merge driver. When merging parallel PRs, resolve `trace/graph.json` conflicts manually in git (**no** merge driver). Combine arrays by **UUID union**: keep every distinct `id` (or `goal_id` for `goal_plan_state`). If the same UUID appears twice after your edit, **keep one object — the later entry in the array wins** on the next `trace seed import`. Re-import after merge; importer applies **last-import-wins** upsert for entities, links (duplicate no-op), and plan-tree rows.
6. **Hook (DF-86):** optional **`trace install git-hook --write`** installs local post-commit/pre-push fragments that run `trace index` on changed paths and may export `trace/graph.json` — **never wraps `git commit`**. Manual `trace index` + `trace seed export` remain valid without a hook.
7. **Clone honesty (DF-88):** Default `seed export` **omits** reviews, transitions, and task `work_state`. After `seed import`, tasks are **PENDING** until the clone operator transitions them. Live DONE/SKIPPED in the exporter’s `.trace/` is local process, not portable identity. `why` / `plan show` work from links + plan tree without reviews/`work_state`.

## Project eval rules

Projects may commit **`trace/eval-rules.json`** at the repo root (not under `.trace/`). Schema v1: `{ "version": 1, "mechanisms": ["stored_test", …], "invariants": [{"id":"internal_must_not_import_cmd","enabled":true}] }`. Missing file → all four built-in mechanisms enabled with the default invariant on. Unknown mechanism ids in the file are skipped; invalid JSON or unsupported version fails closed. `trace eval rules` prints the effective config; seed export adds `eval_rules_path` when the file exists (pointer only — rules body stays git-canonical).

## License

The project core is intended to use Apache-2.0.

Before accepting large contributions, maintainers may use a contributor agreement or other contribution mechanism appropriate for future commercial/hosted offerings.
