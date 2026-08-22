# P01 / S04 / 01 — MCP adapter

## Metadata
- id: P01-S04-01
- todo_ids: [P01-S04-01]
- role: implementer
- skills: [incremental-implementation, mcp-builder, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Ship a **thin MCP server** that exposes a bounded subset of live CLI semantics (`why` / `context` required; `add` / `link` / `transition` / `review` parity) over the Go library. **No raw SQL. No forked business logic** in the MCP package (G19 / DR-API). X0 remains runnable via CLI without MCP (DR-AGENT).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-AGENT, DR-SURFACE, DR-API
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19
- Sibling [00-PLANNER.md](00-PLANNER.md) locks (this scope)
- Live: `cmd/trace` (why/context/add/link/transition/review); `internal/{store,domain,retrieval,compiler}`
- Official SDK: `github.com/modelcontextprotocol/go-sdk/mcp` (stdio `StdioTransport` + `mcp.AddTool`)

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Live substrate (do not re-guess)

| Fact | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| CLI | Thin `cmd/trace` — G19 library-only; commands include why/context/add/link/transition/review |
| Why | `retrieval.New(st).Why(ctx, entityType, entityID)` → JSON on stdout |
| Context | `compiler.New(st).WithRetrieval(eng)` → `TaskContext` (depth 1) / `ExpandContext` (depth 2); format json\|markdown\|both |
| Domain writes | `domain.New(st)` Create*/Link*/TransitionTask/CreateReview/SetReviewResult/LinkReviewTask |
| Store | `store.Open(abs)` → `.trace/trace.db`; mig 001–005 live |
| S03 X0 | `evals/x0` B0/G1 dry-run via **CLI** why/context — MCP must not gate that path |
| Forbidden today | No MCP package yet; no daemon/HTTP product surface; no embeddings |

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | Keep `go.mod` floor (currently 1.24.0); do not downgrade |
| Layout | **`internal/mcp`** (server + tool handlers) + thin **`cmd/trace-mcp`** `main` only |
| SDK | Official **`github.com/modelcontextprotocol/go-sdk`** — import `…/mcp`. **Do not** use mark3labs/mcp-go or other forks |
| Transport | **stdio only** via `&mcp.StdioTransport{}`. **Not** always-on HTTP/SSE/daemon-as-product |
| Server identity | `Implementation{Name: "trace", Version: "0.0.0-dev"}` (match CLI version string) |
| Project root | Binary flags `-C` / `--project` (same Abs resolve as CLI) **and** optional per-tool `project` string; default cwd |
| G19 | MCP imports library only. `internal/{store,vcs,gitcli,analyzers,domain,retrieval,compiler}` and `cmd/trace` must **not** import `internal/mcp` or `cmd/trace-mcp` |
| Logic | **No** domain/retrieval/compiler reimplementation inside MCP — open store, call library, format MCP tool results |
| Required tools | `why`, `context` (see schemas below) |
| Parity tools (same scope) | `add`, `link`, `transition`, `review` — thin domain wrappers mirroring CLI semantics |
| Out of S04 tools | `init`, `index`/`reindex`, `seed` (CLI/harness); raw SQL; GraphQL; embeddings; dump-all |
| Out surfaces | Daemon-as-product; HTTP primary; embeddings |
| X0 / DR-AGENT | MCP **not** required for X0; regression keeps `evals/x0` + `evals/p0x` + `evals/honesty` green via CLI |
| CGO | MCP package itself should not need CGO; full `./...` may still need `CGO_ENABLED=1` because analyzers/CLI tests link CGO |
| Module dep | Add official go-sdk via `go get`; pin whatever current stable the SDK resolves |

### Tool schemas (locked names + args)

Tool names use `trace_` prefix for discoverability.

#### Required

```text
trace_why
  project?: string          # override project root (else server -C / cwd)
  entity_type: string       # same vocabulary as `trace why <type> <id>`
  id: string
  → call retrieval.Why; return JSON text content (WhyResult)

trace_context
  project?: string
  task_id: string
  depth?: number            # default 1; allow 1|2 only (max 2)
  format?: string           # json|markdown|both; default json
  include_why?: boolean     # default false
  → depth 1: compiler.TaskContext; depth 2: ExpandContext
  → return packet JSON and/or markdown text per format (mirror CLI)
```

#### Parity (same implement row)

```text
trace_add          # kind + fields → domain Create* (goal|task|decision|assumption|discovery|plan-change|claim|evidence)
trace_link         # rel + ids → domain Link* (goal-task|decision-task|discovery-plan-change|claim-evidence)
trace_transition   # task_id + to_state + optional allow_done / evidence — TransitionTask DONE gate unchanged
trace_review       # create | set result PASS|FAIL|UNCERTAIN; optional task link — CreateReview/SetReviewResult/LinkReviewTask
```

Parity tools must accept the same semantic fields the CLI documents (help text is the contract). Exact JSON field names may use snake_case; document them in tool descriptions. **No** new DONE/promotion policy in MCP.

#### Annotations (hints)

- `trace_why`, `trace_context`: `readOnlyHint: true`
- write tools (`add`/`link`/`transition`/`review`): `readOnlyHint: false`; set `destructiveHint` only if the SDK path is clear — prefer honest idempotent/destructive hints without inventing behavior

### Target tree

```text
internal/mcp/
  server.go          # NewServer / RegisterTools / RunStdio
  tools_why.go       # trace_why handler → retrieval
  tools_context.go   # trace_context handler → compiler
  tools_write.go     # add/link/transition/review → domain (or split files)
  project.go         # resolve project root Abs; Open store helper
  mcp_test.go        # handler tests: library calls, no domain fork

cmd/trace-mcp/
  main.go            # parse -C/--project; call internal/mcp RunStdio
```

Do **not** put business rules in `cmd/trace-mcp`. Do **not** create a second SQLite access path.

### Tests (required)

- Unit/handler: temp project + `store.Open`; `trace_why` / `trace_context` return library results for seeded entities (reuse fixture patterns or minimal seed via domain APIs)
- Prove MCP package does not reimplement Why/context (call path is retrieval/compiler)
- Import boundary: no `internal/*` library package imports `internal/mcp` (grep or test)
- `go build -o bin/trace-mcp ./cmd/trace-mcp` succeeds
- Regression: `CGO_ENABLED=1 go test ./evals/x0/... ./evals/p0x/... ./evals/honesty/...` + `./...` PASS
- Board Notes must list the **final tool names** for S05 VERIFY checklist

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] `internal/mcp` + `cmd/trace-mcp` exist; stdio server runs with official go-sdk
- [ ] `trace_why` + `trace_context` call library only (no domain/retrieval fork)
- [ ] Parity tools `trace_add` / `trace_link` / `trace_transition` / `trace_review` wired as thin domain adapters **or** explicitly deferred in Notes with reason (prefer ship all four in this row)
- [ ] Smoke/handler tests + import-boundary check green
- [ ] `go build ./cmd/trace-mcp` ok; `CGO_ENABLED=1 go test ./...` green (incl. x0/p0x/honesty)
- [ ] TODO.md status + Notes updated with **tool list** for VERIFY
- [ ] No daemon/HTTP primary; no embeddings; no raw SQL tools

## Minimal todos
- [ ] `go get github.com/modelcontextprotocol/go-sdk` (or current module path from pkg.go.dev)
- [ ] Scaffold `internal/mcp` + thin `cmd/trace-mcp`
- [ ] Implement `trace_why` + `trace_context`
- [ ] Implement parity write tools (add/link/transition/review)
- [ ] Tests + build + regression evals
- [ ] Board Notes: tool list + binary path
