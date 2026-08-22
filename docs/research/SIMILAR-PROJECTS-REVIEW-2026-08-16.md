# Similar-projects technique review — 2026-08-16

Filled from [`SIMILAR-PROJECTS-REVIEW-OUTPUT-TEMPLATE.md`](SIMILAR-PROJECTS-REVIEW-OUTPUT-TEMPLATE.md) using [`SIMILAR-PROJECTS-REVIEW-PROMPT.md`](SIMILAR-PROJECTS-REVIEW-PROMPT.md). Read-only on `similar projects/`; no `docs/TODO.md` board mutation.

**Meta**

| Field | Value |
|-------|-------|
| Review date | 2026-08-16 |
| Reviewer / agent | Composer (Cursor agent) |
| Trace commit / branch | *(workspace has no `.git` at review time)* |
| Peers path | `similar projects/` |
| Baseline read | `AGENTS.md`, `docs/ARCHITECTURE.md`, `docs/EVALUATION.md` (H1–H7), `docs/ROADMAP.md`, `docs/init/G_PROJECT_LAWS.md`; skim `RETRIEVAL_AND_CONTEXT.md`, `AGENT_ENVIRONMENT.md` |

---

## A. Inventory

| Dir | One-line what it appears to be |
|-----|--------------------------------|
| `agentrq` | Agent–human collaboration platform: shared workspace/tasks via MCP, Go API + Vue UI, real-time SSE |
| `codebase-memory-mcp` | Local tree-sitter knowledge graph + MCP tools for AI coding agents; fast structural indexing (Pure C) |
| `codegraph` | Local semantic code graph / surgical context for agents (TS + Rust kernel); install across many agent clients |
| `graphify` | On-demand project knowledge graph (tree-sitter for code); query/path/explain instead of grepping |
| `graphiti` | Temporal context graphs for agents (Zep OSS): evolving facts, provenance/episodes, hybrid retrieval |

### Non-clone peers (no source under `similar projects/`)

| Peer | Notes for this review |
|------|------------------------|
| Obsidian (obsidianmd) | Vault graph / local notes UX — **methods** only: link-first navigation, local markdown as SoT; not a Trace product target |
| Notion | Plane-control / task manager identity — **reject** as Trace product; optional dogfood metaphor for human task boards only |

---

## B. Technique findings

### `agentrq`

- **Callback-injected protocol adapter** — MCP tools bind to typed function fields (`CreateTaskFunc`, …) so the adapter package does not own domain imports; composition root wires callbacks. Steal the *shape* (Law 19); do not copy their repo-bypass wiring. Paths: `backend/internal/controller/mcp/server.go`, `backend/internal/app/app.go`.
- **Graduated permission gates + separate tool-call audit** — auto-allow (exact MCP / shell base-command globs) → else human pending; decisions persist as `ToolCall` rows (`auto_allowed|pending|allowed|denied`), not chat-as-sole-record. Paths: `autoallow.go`, `persistToolCall` in `server.go`.
- **Single-ongoing concurrency + query-shaped dequeue index** — at most one `ongoing` task per workspace; `GetNextTask` equality filter + composite index (peek, not atomic claim — Trace should txn-claim if adopted). Paths: `controller/crud/task.go`, `repository/base/repository.go`.
- **Audience / ownership scoping** — human vs MCP token audiences; store predicates filter `user_id`/workspace (fail-closed IDOR tests). Map to local project/session scopes; reject OAuth product. Paths: `service/auth/jwt.go`, `handler/mcp/mcp.go`.
- **Typed event → trigger template → task with provenance IDs** — publish fans out to templates; tasks carry `TriggerID`/`EventID`. Adapt as Discovery→PlanChange provenance, not HTTP/SSE buses. Paths: `controller/event/event.go`.
- **Idempotent standing-instructions injection** — append workspace note only if absent. Adapt as bounded standing context packets with provenance; never as evidence. Paths: `appendSelfLearningNote` in `task.go`.
- **Client-identity hash for telemetry** — prefer MCP `_meta` clientInfo; hash into attribution. Accept for local/opt-in eval dogfood. Paths: `controller/mcp/client_identity.go`.
- **SQLite WAL + busy_timeout at DSN boundary** — small persistence hardening. Paths: `repository/sqlite/sqlite.go`.

### `codebase-memory-mcp`

- **Content-hash + mtime/size watermark incremental index** — `file_hashes` store `(sha256, mtime_ns, size)`; classify dirty via mtime+size; reparse only changed; **mode-skipped** files keep coverage (FAST vs FULL); sealed staging publish. Paths: `src/pipeline/pipeline_incremental.c`, `src/store/store.c`.
- **LSP-surface watermark + closure-repair delta** — per-file cross-file resolution surface hashed; body-only edits skip dependent recompute; clone-and-patch publish. Adapt pattern without shipping embedded Hybrid LSP engines (Law 13). Paths: `src/pipeline/lsp_surface.h`, `pipeline_delta.c`.
- **Hybrid LSP as second pass over tree-sitter** — AST always answers; semantic pass upgrades `CALLS`/`RESOLVED_CALLS` vs weak `USAGE`. Defer engines; keep method. Paths: README Hybrid LSP, `pass_lsp_cross.c`, `pass_usages.c`.
- **Conditional install matrix** — `STABLE|CONDITIONAL|OPT_IN`; write only when marker/platform/existing config proves activity; capability caps; no YOLO/plugin flips. Paths: `src/cli/agent_clients.h`, README Multi-Agent Support.
- **Typed edges + multi-seed impact BFS** — precision-graded edge types; one multi-source BFS with seed exclusion, hop risk, loud `truncated`. Paths: `cbm_store_bfs_multi` in `store.c`.
- **Bounded agent responses** — compact tabular encoding + exact totals/`truncated` (no silent caps). Adapt loud truncation into compiler packets. Paths: `src/mcp/compact_out.h`.
- **Scout / Verify / Auditor tier contracts** — narrow discovery vs evidence vs auditor limitations; coverage ≠ completeness. Adapt as skill/policy, not 43-client fleets. Paths: `agent_profiles.h`, README.
- **Eval methodology** — Graph vs Explorer A/B; Sillito-anchored D1–D5; blind LLM judge vs mechanical tokens; N/A dims excluded; dual token ratios; pilot-first. Paths: `docs/EVALUATION_PLAN.md`, `docs/BENCHMARK.md`.

### `codegraph`

- **Flat native extract buffers + ABI version** — one crossing per file; `KERNEL_ABI_VERSION` parity gates. Adapt IR/versioning in Go analyzers; reject Rust sibling without measured need. Paths: `codegraph-kernel/`, `docs/design/native-extraction-kernel.md`.
- **Two-tier FS incremental reconcile** — `(size, mtime)` prefilter → content-hash → scoped reparse; removals resurrect inbound refs. Paths: `src/extraction/index.ts` sync, `src/sync/`, `schema.sql`.
- **Surgical explore packet + budgeted render** — source + relationship map + flow spine + compact blast; adaptive skeletonization; session dedup. Prefer one strong CLI context command. Paths: `src/mcp/tools.ts`, `docs/design/explore-budget-allocation.md`, `adaptive-explore-sizing.md`, `explore-session-dedup.md`.
- **Depth-bounded impact radius (contains asymmetry)** — walk incoming deps; expand containers out via `contains`; never climb contains-up into siblings. Paths: `src/graph/traversal.ts` `getImpactRadius`.
- **Install-target registry** — `detect`/`install`/`uninstall` per client; idempotent reverse. Paths: `src/installer/targets/`.
- **Staleness honesty at emission time** — pending-sync banners; emission-time size/mtime then hash vs index; prefer false-fresh over false-stale. Paths: `isFileStaleOnDisk` / `withStalenessNotice` in `tools.ts`.
- **Optional WAL checkpoint valve / unresolved_refs** — adapt only under measured WAL/resolution pain. Paths: `src/db/wal-valve.ts`, `schema.sql`.

### `graphify`

- **Typed edge provenance `EXTRACTED|INFERRED|AMBIGUOUS`** — required on every edge; schema-validated before build; report mix %. Paths: `ARCHITECTURE.md`, `graphify/validate.py`, `report.py`.
- **Two-pass resolution** — AST/import-guided → `EXTRACTED`; bare cross-file / indirect → `INFERRED` (not fake precise `calls`); skip ambiguous. Paths: `symbol_resolution.py`, `extractors/engine.py`.
- **Query / path / explain instead of grep** — depth + token budget; confidence on each hop; fail-closed on ambiguous labels; hub truncation accounts for cut edges. Paths: `skills/.../query.md`, `tests/test_path_cli.py`, `test_explain_cli.py`.
- **Report-time honesty** — rank surprises `AMBIGUOUS` > `INFERRED` > `EXTRACTED`; “verify INFERRED” prompts. Paths: `analyze.py`, `report.py`.
- **Community detection** — Leiden/Louvain; dogfood-only for Trace (Law 13). Path: `cluster.py`.
- **On-demand build framing** — OSS is agent-triggered map→query; always-on is SaaS. Steal query-first UX, not dump-rebuild default (Law 12).

### `graphiti`

- **Bi-temporal fact windows** — `valid_at`/`invalid_at` (world) vs `created_at`/`expired_at` (store); close facts, don’t delete. Paths: `graphiti_core/edges.py`, `edge_operations.py` `resolve_edge_contradictions`.
- **Episodes as immutable provenance** — raw ingest units; derived edges keep episode UUID lists. Paths: `nodes.py` `EpisodicNode`, `EpisodicEdge`.
- **Two-phase contradiction → supersession** — LLM proposes duplicate vs contradict; **code** sets `invalid_at` / expires; multi-witness append. Paths: `prompts/dedupe_edges.py`, `resolve_extracted_edge`.
- **Deterministic-first entity resolution** — exact → entropy → MinHash/LSH → LLM residual. Prefer path+qualified-name keys in Trace. Paths: `dedup_helpers.py`, `node_operations.py`.
- **Hybrid retrieval + RRF** — parallel channels fused with reciprocal rank fusion. Adapt as SQLite FTS5 + graph walks; reject embedding-default. Paths: `search_config_recipes.py`, `search_utils.py` `rrf()`.

---

## C. Candidate steals (ranked)

| Rank | Technique | Source | Trace surface (H1–H7 / named) | Fit vs laws | Effort | Risk | Lane | Notes |
|-----:|-----------|--------|-------------------------------|-------------|--------|------|------|-------|
| 1 | Edge confidence enum `EXTRACTED\|INFERRED\|AMBIGUOUS` on structural edges + Why/context display | graphify | H1, retrieval/Why, analyzers | accept | S–M | low | product | Law 5; causal entities already have confidence — extend to code edges |
| 2 | Emission-time staleness / pending-index honesty banners | codegraph | H5, H6, compiler packets | accept | S | low | product | Complements Law 18 STALE; false-fresh preferred |
| 3 | Loud truncation + totals in context/impact packets | codebase-memory-mcp (+ codegraph budget docs) | H6 | accept | S | low | product | Laws 6–7; format-agnostic |
| 4 | Conditional / marker-gated install matrix + detect→install→uninstall registry | codebase-memory-mcp, codegraph | install/agent UX, H7 | accept | M | low | product | No YOLO; MCP adapter optional later — not P0 architecture |
| 5 | Graduated capability allowlist + durable tool-decision audit (≠ chat) | agentrq | H7, honesty/eval | accept-with-adaptation | M | med | product | Reject AllowAll as default; map to capability surface (Phase 06 already exists) |
| 6 | Multi-seed impact BFS + depth-bounded radius with contains asymmetry | codebase-memory-mcp, codegraph | H4, indexing | accept-with-adaptation | M | med | product | Wire into existing impact domain; no new MCP tool menu |
| 7 | Surgical context packet composition (source + edges + spine + blast) + skeletonization / session dedup | codegraph | H6, `trace context`/`why` | accept-with-adaptation | M–L | med | product | CLI-first; reject explore-as-brand / Read-deny hooks |
| 8 | Conservative two-pass call resolve tagged INFERRED; skip ambiguous | graphify | H1 analyzers | accept | M | med | product | Prefer over LLM similarity edges |
| 9 | Episode/evidence pointer discipline + invalidate≠delete for semantic facts | graphiti | H3, H5, causal domain | accept-with-adaptation | M | low | product | Map onto Evidence/`MarkStale`/`SUPERSEDED`; no Neo4j |
| 10 | LLM-propose / code-apply contradiction → PlanChange + STALE/SUPERSEDED | graphiti | H3, H4 | accept-with-adaptation | M | med | product | Law 2; no silent overwrite |
| 11 | Surface-hash / API-shape skip for dependents (mtime+hash watermarks; mode-skipped GC) | codebase-memory-mcp | H1 indexing | accept-with-adaptation | M | low | product | Harden existing SHA incremental; no full Hybrid LSP matrix |
| 12 | Scout/Verify/Auditor query-tier skill policy | codebase-memory-mcp | H6, H7, dogfood | accept-with-adaptation | S | low | dogfood-only → product later | Policy text first; not 43-client installer |
| 13 | Graph vs Explorer A/B + literature-anchored dims + blind judge / dual token disclosure | codebase-memory-mcp | honesty/eval, H1/H6 | accept | M | low | dogfood-only | Strengthen gates; don’t import 159-lang sweep |
| 14 | RRF over FTS + graph walks (no embeddings default) | graphiti | H6 retrieval | accept-with-adaptation | M | med | product | Law 13 until measured need for embeddings |
| 15 | Single-active-task concurrency + atomic claim | agentrq | H2, causal tasks | accept-with-adaptation | S | low | product | Fix AgentRQ peek→claim under txn |
| 16 | Thin adapter callback injection over canonical library | agentrq | CLI/MCP adapters | accept | S | low | product | Law 19; already Trace intent |
| 17 | Report mix % + “verify INFERRED” prompts | graphify | H5 review UX | accept | S | low | dogfood-only | Easy Why footer |
| 18 | Deterministic entity identity cascade (exact → fuzzy residual) | graphiti | H1, entity_links | accept-with-adaptation | M | med | product | Path+qualified-name first; reject embedding SoT |
| 19 | Typed Discovery publish → PlanChange templates with provenance IDs | agentrq (reframed) | H3 | accept-with-adaptation | M | med | product | Not SSE/event-bus product |
| 20 | Community / subsystem clustering | graphify | H1 | accept-with-adaptation | L | high | dogfood-only | Law 13; heuristics only |
| 21 | Embedded Hybrid LSP engines / napi Rust extract kernel | codebase-memory-mcp, codegraph | indexing | reject as product default | L | high | reject | Law 13; keep Go tree-sitter |
| 22 | In-process SSE fan-out / always-on FS daemon as architecture | agentrq, peers | — | reject | — | high | reject | AGENTS hard boundary |

---

## D. Explicit rejects

| Idea | Why reject (law / hard boundary / wrong-product) |
|------|--------------------------------------------------|
| MCP / daemon / HTTP / SSE as P0 critical-path architecture “because peers have it” | `AGENTS.md` hard boundary; MCP only as thin adapter later |
| Full-rebuild-on-any-change indexer | Law 12 / DR-INCREMENTAL |
| Ship 15+ MCP tools / Cypher console / graph viz UI (`localhost:9749`) | Swiss-knife / wrong-product; Trace is deep modules + CLI |
| Neo4j / FalkorDB / Neptune / proprietary graph DB as SoT | Law 13; keep `.trace/` SQLite |
| Embeddings / vector store / in-binary embed models as default retrieval | Law 13 until measured need |
| 158-language grammar megastore / Hybrid LSP matrix / Rust napi kernel as Trace rewrite | Law 13 + capability sprawl vs TS/JS+Python depth |
| AgentRQ OAuth + Vue kanban + Docker Compose product identity | Wrong-product (generic AI project manager) |
| YOLO / AllowAllCommands as default capability policy | Laws 9, 17; honesty debt |
| Notion-like plane control / Obsidian clone as Trace core | Wrong-product; optional metaphors only |
| LLM similarity “surprising connections” as verified edges | Law 5; graphify’s own corpus review flags false positives |
| Zep cloud / Context Graph Engine | Product identity, not a method |
| Community detection (Leiden) as critical-path dependency | Law 13; dogfood heuristics only |

---

## E. FUTURE PHASE (run when ready; do not assume Phase 11)

> Paste onto `docs/TODO.md` only when the human schedules it. Phase number is TBD — use `Pxx` placeholders below.

`FUTURE PHASE (run when ready; do not assume Phase 11)`

### Human schedule note (2026-08-16)

**Thin Phase 12 human-scheduled** — scopes **S01 + S02 only** (edge provenance + packet honesty) + VERIFY as phase S03. Boarded as [`docs/phases/phase-12-peer-honesty-surfaces/`](../phases/phase-12-peer-honesty-surfaces/) with `P12-*` rows.

### Human schedule note (2026-08-17)

**Thin Phase 14 human-scheduled** — research FUTURE **S03 + S04 only** (impact walks + install/capability gates) + VERIFY. Boarded as [`docs/phases/phase-14-peer-impact-install-gates/`](../phases/phase-14-peer-impact-install-gates/) with `P14-*` rows (goals-gap #1). **Do not** treat supersession (research S05) or the full S01–S06 outline as boarded; ranks 7+ / S05 / `plan simulate` / D21+ remain **deferred** until a later promotion.

### Suggested name + slug

- **Name:** Peer technique adoption (honesty edges, packets, impact walks, install gates)
- **Slug:** `phase-xx-peer-technique-adoption` — *(thin cut boarded as `phase-12-peer-honesty-surfaces`)*

### Problem

Phases 00–11 closed the core causal graph, progressive planning, review, capability surface, and residual honesty gates — but peer projects still show **method gaps** Trace can close without cloning their products: structural edges lack EXTRACTED/INFERRED tagging; context packets under-signal staleness and truncation; impact walks can learn multi-seed / contains-asymmetric rules; install/capability UX can harden marker-gated activation and durable permission audits. This phase adopts those techniques behind existing CLI/library boundaries, with dogfood evals proving H1/H4/H6/H7 lifts — not a peer-feature checklist.

### Non-goals

- No new MCP/daemon/HTTP product architecture; no graph DB or embedding SoT
- No language-count arms race; no Vue/Notion/Obsidian product surfaces
- No rewriting Phase 00–11 board history; spawn forward only
- No full Hybrid LSP / Rust kernel rewrite without Gate-style measured need

### Ordered scopes

| Scope | One-line |
|-------|----------|
| S01 | Structural edge provenance (`EXTRACTED`/`INFERRED`/`AMBIGUOUS`) + Why/context surfacing |
| S02 | Compiler packet honesty: staleness banners, loud truncation, optional skeletonization/dedup |
| S03 | Impact walk upgrades: multi-seed BFS + depth/contains asymmetry wired to existing impact domain |
| S04 | Install/capability gates: marker-gated client registry + graduated allowlist audit (library/CLI) |
| S05 | Causal supersession polish: episode/evidence pointers + contradict→PlanChange/STALE (SQLite) |
| S06 VERIFY | Ablations/dogfood for H1/H4/H6/H7; DR-HANDOFF (no assumed successor) |

### Board-ready TODO stub

```markdown
## Phase XX — Peer technique adoption (FUTURE; schedule when ready — do not assume Phase 11)

| Order | ID | Status | Prompt | Notes |
|------:|----|--------|--------|-------|
| 1 | Pxx-00 | pending | phases/phase-xx-peer-technique-adoption/00-PHASE-PLANNER.md | phase planner; lock S01→S06; cite docs/research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md |
| 2 | Pxx-S01-00 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-01-edge-provenance/00-PLANNER.md | edge confidence enum |
| 3 | Pxx-S01-01 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-01-edge-provenance/01-edge-provenance.md | implement analyzers+store+Why display |
| 4 | Pxx-S01-02 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-01-edge-provenance/02-scope-review.md | |
| 5 | Pxx-S02-00 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-02-packet-honesty/00-PLANNER.md | staleness + truncation |
| 6 | Pxx-S02-01 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-02-packet-honesty/01-packet-honesty.md | compiler/`trace context` |
| 7 | Pxx-S02-02 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-02-packet-honesty/02-scope-review.md | |
| 8 | Pxx-S03-00 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-03-impact-walks/00-PLANNER.md | multi-seed + contains rules |
| 9 | Pxx-S03-01 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-03-impact-walks/01-impact-walks.md | |
| 10 | Pxx-S03-02 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-03-impact-walks/02-scope-review.md | |
| 11 | Pxx-S04-00 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-04-install-capability-gates/00-PLANNER.md | marker install + allowlist audit |
| 12 | Pxx-S04-01 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-04-install-capability-gates/01-install-capability-gates.md | |
| 13 | Pxx-S04-02 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-04-install-capability-gates/02-scope-review.md | |
| 14 | Pxx-S05-00 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-05-supersession-episodes/00-PLANNER.md | evidence pointers + contradict pipeline |
| 15 | Pxx-S05-01 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-05-supersession-episodes/01-supersession-episodes.md | |
| 16 | Pxx-S05-02 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-05-supersession-episodes/02-scope-review.md | |
| 17 | Pxx-S06-00 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-06-phase-verify/00-PLANNER.md | VERIFY planner |
| 18 | Pxx-S06-01 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-06-phase-verify/01-verify.md | ablations + keep prior gates green |
| 19 | Pxx-S06-02 | pending | phases/phase-xx-peer-technique-adoption/scopes/scope-06-phase-verify/02-scope-review.md | VERIFY review + DR-HANDOFF |
```

---

## F. Open questions for the human

1. **Schedule now vs later?** ~~Roadmap is closed (`no successor` after P11).~~ **Resolved 2026-08-16:** thin Phase 12 (S01+S02) human-scheduled; full S01–S06 remains research-only until further promotion.
2. **Scope cut:** ~~Prefer a thin first slice (S01+S02 only…)?~~ **Resolved:** thin S01+S02 boarded as Phase 12.
3. **MCP adapter timing:** Keep install matrix CLI/rules-only, or allow optional thin MCP install behind the same marker gates (still not P0 architecture)? *(deferred with research S04)*
4. **Eval investment:** Steal CBM’s Graph-vs-Explorer A/B into `experiments/` dogfood before any product Go for research S03/S06?
5. **Non-clone peers:** Any desire to capture Obsidian/Notion methods in a separate note, or explicitly out-of-scope forever?
6. **Workspace git:** Review recorded “no `.git`” — restore VCS before scheduling so DR-HANDOFF/commit discipline returns?

---

## Done criteria checklist

- [x] Every peer dir under `similar projects/` covered (+ non-clone note)
- [x] Steals are techniques with Trace mapping + law fit + lane
- [x] Rejects are explicit
- [x] Future phase outline + board stub present and labeled (not Phase 11)
- [x] No edits under `similar projects/`; no `docs/TODO.md` board mutation
