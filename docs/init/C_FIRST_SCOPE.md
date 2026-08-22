# C — First Deep Scope Plan

## Scope ID

`S0 — Experimental Vertical Slice`  
Contains **P0-X** (foundation tiny experiment) and post-P0 completion toward full **X0** (Gate C agent comparison).

## Objective

Deliver a Go library + `trace` CLI that can (1) prove the data model, retrieval, human causal seeding, and benchmark harness on a synthetic fixture (**P0-X → close P0**), then (2) support agent baseline-vs-graph evaluation (**X0 / Gate C**).

## Why this scope first

Docs alone do not close P0 (DR-P0). The thesis still ultimately needs Gate C, but the foundation must be empirically real first.

## In scope

1. Go module: canonical library + `cmd/trace` CLI (no daemon/HTTP/MCP in P0-X).
2. Repository binding + **`.trace/`** project SQLite store (gitignored).
3. Git adapter: **CLI-backed VCS interface**; commit metadata index, file-at-ref resolution, changed-files; **no content duplication**.
4. Structural code graph for **TypeScript/JavaScript and Python** via **tree-sitter** (files + minimal symbols/imports) — **required for P0-X**, not deferred.
5. Work/causal subset: Project, Goal, Decision, Assumption, Task, Discovery, Claim, Evidence, Review, PlanChange, Event (+ thin Phase/Scope if needed).
6. Provenance + STALE.
7. Hybrid retrieval **without embeddings**: exact, lexical (FTS5), graph expansion depth ≤ 2, temporal via Git.
8. Context compiler: Layer 0–1 JSON + Markdown; reason codes; untrusted labeling.
9. CLI: `init`, `index`/`reindex`, entity add commands, `context`, `why`, seed import; claim/evidence/review for Phase 1.
10. Synthetic Apache-2.0 fixture + **human-curated** ground-truth seed.
11. **T012a** deterministic P0-X harness; **T012b** agent X0 harness after P0.
12. Todo-level review / honesty path in Phase 1 (not required to close P0).

## Out of scope (explicit)

- Daemon / loopback HTTP in P0-X
- MCP until after query/context API validated via CLI
- Embeddings / vector DB
- Environment/capability graph
- Automated decision impact & plan simulate
- Scope/phase review automation
- Multi-agent / worktrees
- Full progressive planner
- UI/dashboard
- Analyzers beyond TS/JS/Python
- Real OSS benchmark corpus (after synthetic)

## Relevant decisions

DR-LANG, DR-NAME, DR-GOMOD, DR-TRACEDIR, DR-SURFACE, DR-AGENT, DR-API, DR-P0, DR-P0X, DR-PARSE, DR-GIT, DR-INCREMENTAL, DR-RISK, DR-BENCH, DR-SEED, DR-SLICE, DR-NOSSEM, DR-NOENV, DR-NOIMP, DR-EVT, DR-CLAIM; D1–D14.

## Relevant assumptions

A1, A4–A8, A9 (CLI-first), A10, A11, A12 (MCP later), A13, A14.

## Dependencies

| Dependency | Status |
|------------|--------|
| Round 1 answers | Done |
| Round 2 answers | Soft-block (defaults allowed) |
| Go toolchain + git on host | Assumed |
| Fixture design | In-scope |

## Implementation boundaries (Go)

```text
cmd/trace                 — thin CLI
internal/...              — domain, sqlite, vcs (interface), gitcli, analyzers (tree-sitter), retrieval, compiler
fixtures/x0               — synthetic project + ground truth
evals/p0x                 — deterministic foundation harness (incl. incremental test)
evals/x0                  — agent comparison harness (post-P0)
```

Module path: **`github.com/mrchatam/Trace`**.

**Library must not import CLI.** Future MCP imports library only.

**DR-RISK:** Prefer correct foundation over fastest scaffold. Do not build disposable file-only or full-rebuild architectures.

## P0-X validation (closes P0) — SETTLED (DR-P0X)

P0-X is complete when **all** hold:

1. Goal / Task / Decision / Discovery round-trip with provenance.
2. Files + minimal symbols/imports can be represented (tree-sitter).
3. `trace why` returns a causal chain with evidence/reason codes.
4. `trace context` produces bounded task-specific context.
5. Human-seeded graph matches fixture ground truth.
6. Deterministic tests pass **several** understanding queries without an LLM (default ≥5; see Q-UNDERSTAND-N).
7. Incremental update of a changed file works **without rebuilding the entire fixture graph**.

**Forbidden architecture:** “whenever anything changes, rebuild the entire project graph.” Incremental localized update must be designed in from T003/T004, not bolted on later.

## Success vs non-success

- **P0-X fail:** revise schema/retrieval/seed UX; do not expand to planner/MCP/impact.
- **P0-X pass, Gate C fail later:** foundation OK; product thesis still endangered — stop feature factory.
- **Both pass:** proceed to progressive planner phases.

## Verification for full S0

See board T014 + honesty path; laws 1–7, 12–15, 19; no blob duplication; incremental reindex smoke.
