# P00 / S06 / 01 — Retrieval + context compiler

## Metadata
- id: P00-S06-01
- todo_ids: [P00-S06-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Implement **hybrid retrieval without embeddings** and the **Layer 0–1 context compiler** as library APIs:

1. **`internal/retrieval`** — exact lookup, SQLite **FTS5** lexical search, graph expansion (depth ≤ 2), and **`Why`** causal chains with **reason codes**.
2. **`internal/compiler`** — progressive **Layer 0–1** context packets (**JSON canonical + Markdown render**), **token/candidate budgets**, and **untrusted-data labeling** for retrieved project text.

Enables P0-X #3 (`why` + reason codes) and #4 (bounded task context). **No** CLI wiring (S07), **no** embeddings/vector DB (DR-NOSSEM), **no** MCP/daemon/HTTP.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [C_FIRST_SCOPE.md](../../../../init/C_FIRST_SCOPE.md) — hybrid retrieval + Layer 0–1 compiler; P0-X #3/#4
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G4 retrieved≠authority, G6 no dumps, G7 progressive context, G1 no blobs
- [D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-NOSSEM, DR-PACK, DR-SURFACE, DR-API, DR-P0X, DR-INCREMENTAL
- [RETRIEVAL_AND_CONTEXT.md](../../../../RETRIEVAL_AND_CONTEXT.md) — layers, pipeline, reason provenance
- [SECURITY.md](../../../../SECURITY.md) §8 — prompt injection / retrieved text
- [J_BRAINSTORMING_OUTCOMES.md](../../../../init/J_BRAINSTORMING_OUTCOMES.md) — default depth 1 + expand; always reason codes; untrusted channel
- [B_INITIAL_BOARD.md](../../../../init/B_INITIAL_BOARD.md) — historical T007/T008
- Live priors: S02 store (no FTS yet — **S06 adds it**); S03 `vcs.Repository` + thin commit index; S04 `ListSymbolsByPath` / `ListImportsByPath`; S05 `domain` Create*/Link*/ListLinksFrom + events; stubs `internal/retrieval/doc.go`, `internal/compiler/doc.go`; `go.mod` **go 1.24.0**

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Live substrate (do not re-guess)

| Fact | Value |
|------|-------|
| FTS today | **Absent.** Migrations `001_init` / `002_vcs_index` / `003_causal_domain` only. S02 **forbade** FTS5 — **this scope** owns mig **`004_fts.sql`** + store search helpers |
| Embeddings | Forbidden for P0-X (DR-NOSSEM) |
| Packages | `internal/retrieval` + `internal/compiler` only — **never** `internal/context` / `internal/contextx` |
| Causal reads | Prefer `domain.Service` Get*/ListLinksFrom (+ store `ListLinksTo` as needed). **Do not** reimplement Create*/Link*/Transition CRUD inside retrieval |
| Structural graph | Store symbols/imports only (S04). Query via `ListSymbolsByPath` / `ListImportsByPath` / `GetFileByPath` |
| Events | `entity.created`, `entity.linked`, `task.transition`, plus `entity.stale` from MarkStale |
| Links | `decision_affects_task`, `discovery_causes_plan_change`, optional `claim_has_evidence`; Goal→Task via `tasks.goal_id` only (event rel `goal_has_task` is payload-only) |

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | `go 1.24.0` in `go.mod` (do **not** downgrade) |
| Retrieval package | `internal/retrieval` (package `retrieval`) |
| Compiler package | `internal/compiler` (package `compiler`) — **not** `context` / `contextx` |
| Persistence | Same `.trace/trace.db` via `*store.Store`. **No** second DB |
| Construction | `retrieval.New(st *store.Store) *Engine` (+ optional `WithVCS(vcs.Repository)`); `compiler.New(st *store.Store) *Compiler` (may take/`WithRetrieval(*Engine)`). Neither package calls `store.Open` |
| Channels (P0-X) | **Exact** + **FTS5 lexical** + **graph expand** (+ optional **temporal** via VCS iface / thin index). **No** semantic/embedding channel |
| Graph depth | Auto context expansion default **depth = 1**. Hard cap **depth ≤ 2**. Deeper than default requires explicit `ExpandContext` (G6/G7) |
| Dump API | **Forbidden** — no `DumpGraph` / unbounded list-all-entities export (G6) |
| Packet formats | JSON **canonical**; Markdown **render** (DR-PACK). Both from same structured packet |
| Layers this scope | **Layer 0 + Layer 1 only**. Layers 2–3 deferred (may appear only as explicit Expand toward depth 2 neighbors — still not full L2/L3 catalog) |
| Reason codes | **Required** on every included hit/item (never silent inclusion) |
| Untrusted labeling | Retrieved project text (titles/bodies/paths/symbol names used as content) marked `trust=untrusted_data` (G4 / A14). Structured system fields (ids, enums, reason codes, budgets) may be `trust=system` |
| Token budget default | **4096** estimated tokens (heuristic OK — rune/word approx; **no** mandatory tokenizer dependency) |
| Candidate caps | Pre-compile max hits **64**; packet max items **32** (after budget trim). Tests must prove hard cap never exceeded |
| Ranking | Deterministic, explainable weights OK (exact > fts > graph-distance). **No** learned ranker |
| CGO | retrieval + compiler + new store FTS APIs must pass `CGO_ENABLED=0`. Do **not** import `analyzers` or `gitcli` (VCS only via `vcs.Repository` if needed) |
| Surface | Library only — **no** new `cmd/trace` subcommands (S07); **no** MCP/daemon/HTTP |
| Out of scope | Embeddings/vector DB; Layer 2–3 auto-load; impact engine; full-graph dump; CLI; fixture harness (S08) |

### Hybrid pipeline (locked shape)

```text
query / task_id
  → exact lookup
  + FTS5 lexical
  → candidate Hits (each with reason_code, score, distance)
  → graph expand (default depth 1; max 2)
  → optional temporal enrich (LastChanged / History via vcs or thin index — refs only, no blobs)
  → budget trim
  → compiler Layer 0–1 packet (JSON + MD)
```

### Exact lookup (locked)

Resolve by:

- entity UUID + type (`goal`, `task`, `decision`, `assumption`, `discovery`, `plan_change`, `claim`, `evidence`)
- file `path` (store `GetFileByPath`)
- symbol name within a path (via `ListSymbolsByPath`) when query specifies path+name or unambiguous single hit
- commit OID only as **metadata** via thin VCS index / `vcs` — never fetch blobs into SQLite

Misses return empty/not-found — not silent fabricated hits.

### FTS5 lexical (locked)

- Migration **`internal/store/schema/004_fts.sql`** (next after `003_causal_domain.sql`).
- Index **text fields already in tables** — titles/bodies of goals/tasks/decisions/assumptions/discoveries/plan_changes (+ claims/evidence if trivial); file `path`; symbol `name` (+ kind optional UNINDEXED). **No source-file body/BLOB columns** (G1).
- Store API: e.g. `SearchFTS(query string, limit int) ([]FTSHit, error)` + sync helpers (`SyncEntityFTS` / `RebuildFTS` or FTS5 triggers). Tests must insert entities then find them by distinctive tokens.
- Tokenizer: SQLite FTS5 default or `unicode61` — document choice in store comment. Porter optional, not required.
- Retrieval `Search` wraps store FTS and attaches reason `fts_match`.

### Graph expand (locked)

Two neighbor kinds, both depth-capped:

1. **Causal** — `entity_links` via `ListLinksFrom` / `ListLinksTo`; Goal↔Task via `tasks.goal_id` (and inverse: tasks with `goal_id=G`). Rels include `decision_affects_task`, `discovery_causes_plan_change`, optional `claim_has_evidence`.
2. **Structural** — file↔symbols; file↔imports (`imported_path` as neighbor path when that path exists in `files`).

Each expanded hop increments `distance` (seed = 0). Refuse expand when requested depth > 2. Default TaskContext uses depth 1.

### Why (locked)

`retrieval.Why(ctx, entityType, entityID) (WhyResult, error)`:

- Walk causal neighborhood (links + goal_id) and relevant events (`ListEventsByEntity`) into an ordered chain/explanation.
- Every step carries a **reason_code** (see vocabulary).
- Include linked Decision→Task / Discovery→PlanChange / Goal→Task when present.
- Optional: thin temporal note via VCS `LastChanged` for a path if `WithVCS` set and a file seed exists — **refs only**.
- Must be sufficient for P0-X #3 / CLI `trace why` later (S07 adapter only).

### Reason-code vocabulary (locked strings)

Use these exact codes where applicable (extend only with package-doc note):

| Code | When |
|------|------|
| `exact_id` | UUID exact hit |
| `exact_path` | File path exact hit |
| `exact_symbol` | Symbol exact hit |
| `fts_match` | Lexical FTS |
| `direct_task_scope` | Seed task / its goal / work_state packet core |
| `goal_has_task` | Goal↔Task via `goal_id` |
| `decision_affects_task` | Link rel |
| `discovery_causes_plan_change` | Link rel |
| `claim_has_evidence` | Link rel (stub OK) |
| `graph_neighbor` | Structural or other causal neighbor (include `distance`) |
| `recent_event` | Event-backed inclusion |
| `historical_vcs` | Optional VCS/temporal enrich |
| `budget_dropped` | Only in debug/trace of drops — **not** on included items |

Do **not** emit `semantic_match` in P0-X.

### Context compiler Layer 0–1 (locked)

**Layer 0** (always, subject to tiny size):

- task id/title/body snippet, `work_state`, provenance `status`
- linked goal (objective) when `goal_id` set
- exit criteria: use task `body` (or a conventional section if present) — do not invent a new DB column this scope unless already present

**Layer 1** (budgeted):

- decisions linked to the task (`decision_affects_task`)
- open/active assumptions when exact/FTS/graph justifies (prefer linked or fts_match)
- direct files/symbols when path mentions or structural neighbors at distance ≤ default depth
- recent discoveries / plan_changes when linked or clearly in-budget fts/graph hits
- claim/evidence only if linked stubs exist

**Not auto-included:** full repo file list, all events, Layer 2–3 catalogs, embedding neighbors.

### Packet + untrusted labeling (locked)

```text
Packet {
  schema_version: "0.1"
  layer: 0|1
  task_id, generated_at
  budget: { token_limit, tokens_est, max_items, truncated: bool }
  items: [{
    entity_type, entity_id, title?, excerpt?,
    reason_code, distance?, score?,
    trust: "system" | "untrusted_data",
    provenance?: { status, source_type, confidence }
  }]
  why_trace?: [...]   // optional pointer/summary when Compile follows Why
}
```

- Markdown render must **label** untrusted sections (e.g. heading/`trust: untrusted_data` callout). Untrusted excerpts must not be presented as project policy.
- `truncated: true` when budget enforcement drops items.
- Q-INJECTION remains OPEN in the ledger; this labeling is the **P0-X provisional** mitigation (A14).

### Agent-facing library entrypoints (locked names may vary slightly; behavior locked)

```text
// retrieval
New(st *store.Store) *Engine
WithVCS(repo vcs.Repository) *Engine   // optional chaining

Exact(ctx, ExactQuery) ([]Hit, error)
Search(ctx, q string, opts SearchOptions) ([]Hit, error)   // FTS; opts.Limit default ≤ 32
Expand(ctx, seeds []Hit, depth int) ([]Hit, error)         // depth in 1..2 only
Why(ctx, entityType, entityID string) (WhyResult, error)

// compiler — P0-X context surface (CLI will call these in S07)
New(st *store.Store) *Compiler
WithRetrieval(e *retrieval.Engine) *Compiler

TaskContext(ctx, taskID string, opts ContextOptions) (Packet, error)
  // Layer 0–1; default expand depth 1; applies budgets; JSON+MD available via Packet methods or twin return
ExpandContext(ctx, taskID string, depth int, opts ContextOptions) (Packet, error)
  // explicit expansion; reject depth>2 or depth<1

// ContextOptions: TokenBudget int (default 4096); MaxItems int (default 32); IncludeMarkdown bool
```

`TaskContext` / `ExpandContext` orchestrate retrieval + compile. S07 must **not** reimplement ranking/budgets in the CLI.

### Store migration / APIs (locked shape)

```text
internal/store/schema/004_fts.sql
  — FTS5 virtual table(s) for entity text + paths/symbols
  — indexes/triggers OR documented RebuildFTS path

internal/store/fts.go (name flexible)
  — SearchFTS, Sync*/RebuildFTS as needed
  — wire Sync into Upsert* / ReplaceSymbols paths OR provide Rebuild used by tests + Open post-migrate backfill
```

Additive store helpers for retrieval are OK (e.g. `ListTasksByGoalID`). Prefer thin SQL — no business ranking in store.

### Target tree

```text
internal/store/
  schema/004_fts.sql
  fts.go                 # SearchFTS + sync/rebuild

internal/retrieval/
  doc.go                 # contract: exact/FTS/graph/Why; no embeddings; depth caps
  engine.go              # New / WithVCS
  exact.go
  search.go              # FTS wrapper
  expand.go              # causal + structural, depth ≤2
  why.go
  types.go               # Hit, WhyResult, reason code consts
  retrieval_test.go

internal/compiler/
  doc.go                 # Layer 0–1; DR-PACK; untrusted labeling; budgets
  compiler.go            # New / WithRetrieval / TaskContext / ExpandContext
  packet.go              # Packet struct + JSON + Markdown render
  budget.go              # trim by tokens/items
  compiler_test.go       # budget never exceeded; injection label present; layer contents
```

### Out of scope (this row)

- CLI `why` / `context` commands (S07) — library only here
- Fixture GT + `evals/p0x` (S08)
- Embeddings / vector tables / semantic_match
- Auto Layer 2–3 packing
- MCP / daemon / HTTP
- Full Claim→Review honesty promotion
- Importing `analyzers` or `gitcli` from retrieval/compiler

## Board rights
Implementer: update **status + notes only** on `P00-S06-01`. Do not spawn rows or rewrite later prompts.

## Exit criteria
- [ ] `store` mig `004_fts` applies on Open; FTS finds seeded entity title/body and/or path/symbol tokens; **no** source BLOBs introduced
- [ ] `retrieval.Exact` / `Search` / `Expand` (depth 1 and 2) / `Why` covered by tests with **reason_code** on hits; depth > 2 rejected
- [ ] Graph expand uses causal links + `goal_id` and structural symbols/imports — no parallel analyzer DB
- [ ] `compiler.TaskContext` returns Layer 0–1 packet with JSON + Markdown; default depth 1; `ExpandContext` can request depth 2
- [ ] Budget tests: token and/or item caps enforced; `truncated` set when drops occur; **no** dump API
- [ ] Untrusted labeling present on retrieved text items in JSON and visible in Markdown
- [ ] No embeddings / vector deps; packages are `retrieval` + `compiler` only (no `context`/`contextx`)
- [ ] `CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./internal/store/... ./internal/domain/...` passes; `go test ./...` green (analyzers may need CGO as today)
- [ ] No MCP/daemon/HTTP; no new CLI commands; retrieval/compiler do not import `analyzers`/`gitcli`
- [ ] TODO.md Notes for `P00-S06-01` updated; status `done`

## Minimal todos
- [ ] Store `004_fts.sql` + SearchFTS/sync + tests
- [ ] retrieval Exact + Search + Expand (causal/structural, depth caps) + reason codes
- [ ] retrieval Why chain tests (Goal/Task/Decision/Discovery paths)
- [ ] compiler Packet JSON+MD + untrusted labels + budget trim
- [ ] compiler TaskContext / ExpandContext orchestration tests
- [ ] Prove no dump API; CGO_ENABLED=0 suite; board status + notes
