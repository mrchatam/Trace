# Similar-projects technique review — agent prompt

Paste-ready prompt for a read-only peer review of clones under `similar projects/`. Findings feed a **future** Trace phase (number TBD — do **not** assume Phase 11). This file is the prompt; fill results using [`SIMILAR-PROJECTS-REVIEW-OUTPUT-TEMPLATE.md`](SIMILAR-PROJECTS-REVIEW-OUTPUT-TEMPLATE.md).

## How to run

1. Open a fresh agent session in the Trace repo root (`/home/ali/Desktop/Trace`).
2. Paste everything under **Prompt (copy below)** into the agent.
3. Ask the agent to write filled findings into a dated copy (e.g. `docs/research/SIMILAR-PROJECTS-REVIEW-YYYY-MM-DD.md`) using the output template structure — do **not** mutate board rows in `docs/TODO.md` unless you explicitly request a future-phase scaffold later.
4. Review the proposed phase outline offline; promote to the board only when ready.

---

## Prompt (copy below)

```text
You are reviewing peer projects cloned under `similar projects/` for techniques that could improve Trace — not for cloning their product identity.

## Mission

Extract reusable **techniques, methods, and patterns** from each peer that might strengthen Trace’s local-first project knowledge/causal graph + progressive planning CLI/MCP for AI coding agents. Produce structured findings and a **FUTURE PHASE** outline the human can paste onto `docs/TODO.md` later. Do **not** create Phase 11 (or any new phase board rows) in this run unless the human explicitly asks.

## Trace baseline (read first)

Work from the Trace repo root. Read at least:

- `AGENTS.md` (stack, hard boundaries, current focus)
- `docs/ARCHITECTURE.md`
- `docs/EVALUATION.md` (hypotheses H1–H7)
- `docs/ROADMAP.md`
- `docs/init/G_PROJECT_LAWS.md`
- Hard boundaries in `AGENTS.md` (repeat below)

Also skim as needed: `docs/RETRIEVAL_AND_CONTEXT.md`, `docs/PLANNING.md`, `docs/REVIEW_AND_VERIFICATION.md`, `docs/AGENT_ENVIRONMENT.md`, `docs/STORAGE_AND_PERFORMANCE.md`, `README.md`.

### Trace hard boundaries (non-negotiable filters)

- No product MCP/daemon/HTTP on the P0-X critical path as architecture
- No full-rebuild-on-any-change indexer architecture
- Trace is not a swiss-army knife — prefer deep modules over feature sprawl
- Laws in `G_PROJECT_LAWS.md` (esp. incremental computation, no major infra without measured need, capability minimization, evidence ≠ claims, bounded retrieval)

### Trace evaluation surfaces (map steals here)

- H1 — Project understanding (persistent graph)
- H2 — Progressive planning
- H3 — Discovery-driven replanning
- H4 — Decision impact
- H5 — Evidence-driven review
- H6 — Progressive context / token budgets
- H7 — Capability-aware planning

Named Trace surfaces may also apply: indexing/analyzers, retrieval/Why/context packets, causal domain (goals/tasks/decisions/discoveries), CLI vs MCP adapters, install/agent UX, honesty/eval harnesses, dogfood under `experiments/`.

## Peer corpus

Directory (exact path): `similar projects/`

| Dir | One-line (from README; refresh if stale) |
|-----|------------------------------------------|
| `agentrq` | Agent–human collaboration platform: shared workspace/tasks via MCP, Go API + Vue UI, real-time SSE |
| `codebase-memory-mcp` | Local tree-sitter knowledge graph + MCP tools for AI coding agents; fast structural indexing |
| `codegraph` | Local semantic code graph / surgical context for agents (Rust kernel); install across many agent clients |
| `graphify` | On-demand project knowledge graph (tree-sitter for code); query/path/explain instead of grepping |
| `graphiti` | Temporal context graphs for agents (Zep): evolving facts, provenance/episodes, hybrid retrieval |
there are also some similar proects that there is no source code for them like obsidian vault https://github.com/orgs/obsidianmd/repositories
or notion from notion.com which is more like plane control but it is in its core a taskmanager 

If the directory gains or loses clones, re-inventory from each README before scoring.

## Operating rules

1. **Read-only on peer clones.** Do not modify files under `similar projects/` unless the human explicitly says otherwise. Trace docs under `docs/research/` may be written/updated as deliverables of this review.
2. Prefer **methods over features.** Good: “incremental AST reindex with content-hash watermark,” “EXTRACTED vs INFERRED edge tagging,” “fail-closed review with evidence IDs,” “install matrix that only activates when client markers exist.” Bad: “add 15 MCP tools like X,” “ship a Vue kanban,” “add Docker Compose because Y has it.”
3. For every candidate steal, score:
   - **Trace relevance:** H1–H7 and/or named surface
   - **Fit vs laws:** accept / accept-with-adaptation / reject — cite the law or hard boundary
   - **Effort:** S / M / L (rough)
   - **Risk:** low / med / high (wrong-product, infra creep, law conflict, eval debt)
   - **Lane:** `product` vs `dogfood-only` vs `reject`
4. **Explicit rejects** (call out when peers tempt them):
   - Cloning another product’s identity or UX wholesale
   - Adding MCP/HTTP/daemon “because peers have it” when Trace laws forbid it on the critical path (MCP as thin adapter later is fine only if laws + roadmap allow; never as P0 critical-path architecture)
   - Full-rebuild-on-any-change indexers
   - Broad “add everything” / swiss-knife expansion
   - Graph DB / embeddings / distributed queues without measured need (Law 13)
5. Do **not** change `docs/TODO.md` phase board rows. Output a paste-ready stub labeled **FUTURE PHASE (run when ready; do not assume Phase 11)**.
6. Do not implement product Go changes in this review run.

## Method (suggested)

1. Re-confirm inventory (README one-liners).
2. Per project: skim README + architecture/docs; sample 2–4 concrete mechanisms in code or docs (indexing, retrieval, planning loops, honesty/review, capability/tool selection, context budgets, agent UX, install/MCP wiring, persistence, eval harnesses).
3. Deduplicate techniques across peers; keep the clearest exemplar.
4. Rank survivors by Trace fit × leverage × law safety.
5. Draft future phase: problem, non-goals, ordered scopes, board stub.
6. Fill the output template structure (see below). Save to a dated file under `docs/research/` if asked, or return in-chat in the same structure.

## Output format (required)

Follow `docs/research/SIMILAR-PROJECTS-REVIEW-OUTPUT-TEMPLATE.md` (same section order). Minimum sections:

### A. Inventory
Table of peer dirs + one-line what each is.

### B. Technique findings
One subsection per peer. Bullets of **techniques/methods/patterns** only (not feature laundry lists). Cite paths/docs lightly so a later agent can re-find them.

### C. Candidate steals (ranked)
Table or list rows with: technique | source peer | Trace surface (H1–H7 / named) | fit vs laws | effort | risk | lane (product / dogfood-only / reject) | notes.

### D. Explicit rejects
Short list of attractive-but-wrong ideas and why (law / boundary / wrong-product).

### E. FUTURE PHASE (run when ready; do not assume Phase 11)

Include:

1. **Suggested name + slug** (e.g. `peer-technique-adoption` — human picks phase number later)
2. **Problem** (2–4 sentences)
3. **Non-goals**
4. **Ordered scopes** (S01… with one-line each)
5. **Board-ready TODO stub** — markdown table rows the human can paste later, IDs like `Pxx-00`, `Pxx-S01-00`, … using placeholder `Pxx` (not a real phase number):

| Order | ID | Status | Prompt | Notes |
|------:|----|--------|--------|-------|
| … | Pxx-00 | pending | phases/phase-xx-…/00-PHASE-PLANNER.md | … |

Label this block clearly:

`FUTURE PHASE (run when ready; do not assume Phase 11)`

### F. Open questions for the human
Anything that needs a product call before scheduling.

## Done criteria

- Every peer in `similar projects/` covered
- Steals are techniques with Trace mapping + law fit + lane
- Rejects are explicit
- Future phase outline + board stub present and labeled
- No edits under `similar projects/`; no Phase board mutation unless human ordered it
```
