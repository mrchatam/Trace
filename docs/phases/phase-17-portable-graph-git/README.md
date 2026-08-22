# Phase 17 — Portable graph via git (thin)

**Status:** **complete** (2026-08-17) — S01–S03 + S04 VERIFY green; two-clone recipe PASS; DR-HANDOFF **`no successor`**. Human-scheduled **queue** after Phase 16 (P16 historical/default DR-HANDOFF remains `no successor`). Phase folder evidence: [`scopes/scope-04-phase-verify/VERIFY-NOTES.md`](scopes/scope-04-phase-verify/VERIFY-NOTES.md) + [`REVIEW-NOTES.md`](scopes/scope-04-phase-verify/REVIEW-NOTES.md).

## Why this phase exists

`.trace/` is gitignored (Law 1). Seed **import** exists; **export** does not. A clone gets code + docs but an empty causal graph **and** an empty plan tree, so `why`/`context`/`plan show` cannot see decisions or the progressive plan. This phase adds a git-friendly **seed JSON v1** snapshot (`trace/graph.json`) — causal entities **plus** plan hierarchy — and an idempotent import path. Not a live-SQLite push, not encryption, not a sync daemon, **not a hosted MCP**.

Research SoT: [`docs/research/PORTABLE-GRAPH-GIT-2026-08-17.md`](../../research/PORTABLE-GRAPH-GIT-2026-08-17.md) (P17-00 investigation; upcoming locks in [`DF-84-FORWARD.md`](DF-84-FORWARD.md)). Findings: [`experiments/DOGFOOD-FINDINGS.md`](../../../experiments/DOGFOOD-FINDINGS.md).

## Disposition matrix (P17-00 FINAL + DF-84-FORWARD)

P17-00 prompt history excluded planner tables and required no git hook. Live lock for **upcoming** S01–S04: [`DF-84-FORWARD.md`](DF-84-FORWARD.md).

| ID | Residual / finding | Disposition | One-line rationale | Home |
|----|--------------------|-------------|--------------------|------|
| DF-80 | No `seed export` | **fix** | Clone cannot reconstruct the graph from git | **S01** |
| DF-81 | Re-import not idempotent (links UNIQUE-fail; extra `entity.created`) | **fix** | UUID upsert + duplicate-link no-op | **S03** |
| DF-82 | No commit path / agents forget to export | **fix** | Convention + AGENTS/help; `.gitignore` unchanged | **S02** |
| DF-83 | Two PRs / merge of `graph.json` undefined | **fix** | Document + last-import-wins upsert; no merge driver | **S03** |
| **DF-84** | Plan tree omitted from seed (P17-00 exclude) | **fix** (forward) | Clone must `plan show` without original `.trace/` | **S01** + **S03** |
| **DF-85** | No `exported_at_commit`; actor vs git evidence unspecified | **fix** (forward) | SHA = snapshot evidence, not identity; actor ≠ auth | **S01** + **S02** |
| **DF-86** | No auto-export hook | **CONDITIONAL / deferred** | `trace install git-hook` later-in-P17; **not blocking VERIFY** | pack residual |
| P16 DF-70/73 | Seed rels / findings incomplete until S05 | **depend** | Do **not** duplicate; export those keys after P16 S05 | — |
| Encryption-as-git | Policy notes in git | **wontfix** | SECURITY §6: don’t put secrets in bodies; token stays local | — |
| Commit `.trace/` | Live SQLite | **wontfix** | Law 1 + gitignore | — |
| Reviews/DONE in default export | Process history | **out** | DF-28 remains Trace-pull; tasks export without work_state | — |
| Hosted MCP / HTTP / OAuth | Cloud product | **out** (later, separate repo) | Not Trace core; see TODO **Later developments** | — |

## Scope order (locked at P17-00; DFs thickened by DF-84-FORWARD)

| Scope | Focus | DFs |
|-------|--------|-----|
| S00 / phase planner | Inventory + disposition + spawn | **done** (P17-00) |
| S01 | `seed export` + plan tree + `exported_at_commit` + round-trip | DF-80, **DF-84**, **DF-85** (export field) |
| S02 | Commit convention + SHA/author evidence + attribution + merge docs | DF-82, **DF-85** (docs) |
| S03 | Idempotent import (UUID upsert incl. plan tree, conflict behavior) | DF-81, DF-83, **DF-84** (import) |
| S04 | Phase VERIFY + **two-clone** recipe + DR-HANDOFF | named S01–S03 + DF-84/85; DF-86 **non-fail** |

No new board rows this cut (no S05). DF-86 stays a labeled residual inside this pack.

## Out of scope unless promoted

- Committing `.trace/`; encryption-as-git; daemon/HTTP/sync service
- Hosted authenticated MCP, OAuth, pointing `trace-mcp` at the internet (separate repo if ever)
- New MCP tools (`seed` stays CLI)
- Duplicating P16 S05 (DF-70/71/72/73/74 / `trace_impact`)
- Git merge driver; NDJSON / one-file-per-UUID split
- Default export of reviews, transitions, tool decisions, tokens, index, allowlists
- Wrapping `git commit` with a Trace wrapper
- Research S05 / `plan simulate` / D21+ ladder / ranks 7+
- Rewriting Phase 00–16 `done` history; hijacking AGENTS current focus from Phase 16

## Assumptions (P17-00 + DF-84-FORWARD)

1. Human cut is **clone-readable semantic collaboration** via seed JSON (causal **and** plan tree), not live DB replication, not a server.
2. P16 S05 will land before any P17 implement row (board order). Export includes S05 keys **if present**; S01 must not re-implement import rels.
3. Single `trace/graph.json` (seed v1) is enough; append-only layouts deferred.
4. VERIFY default DR-HANDOFF = **`no successor`**. Hosted product is not Phase 18.
5. Encryption is not required for policy notes.
6. Phase 16 closed before P17 VERIFY (2026-08-17).
7. P17-00 “planner tables not seed v1” is **superseded for upcoming scopes** — [`DF-84-FORWARD.md`](DF-84-FORWARD.md).
8. DF-86 hook absence does **not** fail VERIFY if the two-clone git-JSON recipe passes.

## Completion bar (VERIFY)

Two clones, **no shared `.trace/`**, offline, no account, no server: `init` + `seed import trace/graph.json` + `index` + `why`/`context` + plan hierarchy readable.

## Parallel track (not board-blocking)

Optional dogfood under `experiments/`; feed new DF-* **forward** only (next free **DF-87**).
