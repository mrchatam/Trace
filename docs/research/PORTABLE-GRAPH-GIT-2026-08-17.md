# Portable graph via git — 2026-08-17

Investigation: should Trace support **clone-readable semantic collaboration** (decisions, relations, why-chains, notes for out-of-scope work) **without** committing live `.trace/` SQLite?

**Verdict: GO.** Thin **Phase 17** queued **after** Phase 16. Next **runnable** board row remains **P16-S02-01** until Phase 16 VERIFY closes.

Sources: README, `docs/ARCHITECTURE.md`, `docs/SECURITY.md`, `docs/init/G_PROJECT_LAWS.md`, `docs/rules/project-rules.md`, `cmd/trace/seed.go` + `help.go`, `fixtures/x0/seed/gt.json`, P16 S05 planner (seed completeness), `.gitignore`, store `Upsert*` vs `InsertLink`, `experiments/DOGFOOD-FINDINGS.md`, `docs/research/TRACE-GOALS-PROGRESS-2026-08-17.md`, `docs/research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md`, D30/D30b handoff + D21–D23 combo.

Phase 16 product Go and `done` rows were not modified. P16 DFs **68, 70–78** stay on Phase 16. New IDs start at **DF-80** (DF-79 left unused, same class as DF-69).

---

## 0. Go / no-go

| Question | Answer |
|----------|--------|
| Do this? | **Yes** — queued thin phase after P16 |
| Commit `.trace/`? | **No** (law + `.gitignore`) |
| Encryption-as-git? | **No** — not required for policy notes |
| Daemon / sync service? | **No** |
| Duplicate P16 S05 seed rels / impact? | **No** — Depends-on P16 S05 |
| Next runnable today? | **Still P16-S02-01** |

**Why go.** README’s first validation question is whether an agent can understand an **unfamiliar** repository using the graph. Today a clone gets code + `docs/` protocol, then `trace init` creates an **empty** `.trace/`. Causal objects live only in a gitignored SQLite file. Seed **import** already exists (`cmd/trace/seed.go`; fixture `fixtures/x0/seed/gt.json`). The missing half is **export + commit convention + idempotent re-import**. That is a projection of the project DB onto git, not a second SoT and not a hosted sync product.

**Why not no-go.** Remaining roadmap items (`plan simulate`, research S05 supersession, D21+ as a product phase) are real, but a clone that cannot `why` a decision cannot use them. Portable graph is the **clone analog** of D30b Trace-pull. Sequence: finish P16 surfaces → then this thin cut → then deeper planner/simulate work.

---

## Loop 1 — Product / laws / current surface

### Should we?

Yes, as a **git-friendly seed JSON snapshot**, not as pushing the live DB.

| Law / doc | Implication |
|-----------|-------------|
| Law 1 / project-rules store row | Git owns source/history; `.trace/` stays local and gitignored |
| Architecture §3 | Project DB owns goals/decisions/tasks; code index is **derived** (never commit) |
| Law 6 | No full-graph **API dumps**; a committed **seed file** is an explicit, bounded snapshot (same shape as today’s import) |
| Law 4 / SECURITY §8 | Committed graph text is **retrieved data**, not policy authority |
| Law 9 | User decisions stay first-class objects — they are exactly what clones miss |
| Law 18 | After clone `index`, semantic facts can be marked STALE vs code; portable JSON is not a substitute for reindex |
| Law 19 | Export/import stay thin CLI over domain/store — no fork |
| SECURITY §6 | Do not ingest secrets; export must not copy `.trace/access.token` |

### What already exists

- **Seed JSON v1 import only:** `trace seed import <file>` (`cmd/trace/seed.go`). Usage line rejects anything other than `import`. Help: import through domain APIs; relative path under `-C`.
- **Schema (live, pre–P16 S05):** `version`, `goals`, `tasks` (`goal_id`), `decisions`, `assumptions`, `discoveries`, `plan_changes`, `claims`, `evidence`, `links`, `transitions`. Unknown top-level keys **rejected**.
- **Link rels (live):** `goal_has_task` / `goal-task`, `decision_affects_task` / `decision-task`, `discovery_causes_plan_change` / `discovery-plan-change`, `claim_has_evidence` / `claim-evidence`. DF-33 aliases `from_id`/`to_id`.
- **P16 S05 (pending, do not duplicate):** DF-70 `discovery_mentions_task`, DF-73 `findings`/`alternatives`, DF-71 compiler impact on why/context. P17 export of those keys **Depends-on P16 S05**.
- **Stable UUIDs in fixtures:** `fixtures/x0/seed/gt.json` pins goal `1111…` / task `2222…` / decision `3333…` — the collaboration pattern already used for evals.
- **Store upsert:** `UpsertGoal` (and sibling entity upserts) `ON CONFLICT(id) DO UPDATE`. Domain `Create*` calls upsert, then **always** appends `entity.created`.
- **Links are not upsert:** `InsertLink` is plain INSERT; `entity_links` has `UNIQUE(from_type, from_id, rel, to_type, to_id)`. Re-import of the same link **fails**.
- **Handoff SoT (DF-28, closed P11):** help text — predecessor task body + linked decisions via retrieval; successors `trace context` / `why`. No first-class handoff entity.

### What a clone misses today

After `git clone` + `trace init` + `trace index`:

| Present | Missing |
|---------|---------|
| Source tree, `docs/`, AGENTS protocol | Goals, decisions, assumptions, discoveries, plan_changes, claims, evidence |
| Derived code graph (after `index`) | Entity links → empty `why` causal neighborhood |
| Empty task list | Why-chains, out-of-scope notes stored as decisions/discoveries |
| Local capability/tool-decision tables (fresh AUTO_ALLOWED builtins) | Prior reviews, DONE history, events log, impact findings (until P16 S05 + export) |

Recipe that **should** work after P17: clone → `init` → `seed import trace/graph.json` → `index` → `why` / `context`.

---

## Loop 2 — Git / merge / what to commit

### Allow (commit)

| Artifact | Why |
|----------|-----|
| **`trace/graph.json`** (seed JSON v1) | Portable semantic snapshot: entities + links (+ P16 S05 findings/alternatives when present) |
| Optional extra seed files under `trace/` (e.g. fixture overlays) | Same schema; import is path-explicit |
| Docs describing the convention | S02 |

Stable **UUIDs** in the JSON are the merge key. Callers must pass `id` on `trace add` when the object will be exported (fixtures already do this).

### Deny (never commit)

| Artifact | Why |
|----------|-----|
| `.trace/` (`trace.db`, WAL, `trace.lock`) | Local store; gitignored; Law 1 |
| Derived index (files/symbols/imports, FTS, VCS watermark) | Rebuild with `trace index` |
| `.trace/access.token` | Local auth; `backup` already excludes token by default |
| `capability_tool_decisions` / AUTO_ALLOWED | Per-machine / per-store allowlist (P14–P16); clone must not inherit YOLO or DENY |
| Capability catalog as environment dump | Skills/MCP availability differs by machine |
| Events log | Derived audit; re-import would duplicate `entity.created` |
| Reviews / DONE / work_state (default) | Local process; DF-28 remains Trace-pull after import. Exporting DONE would require reviews or `allow_done` (dishonest) |
| Planner coarse plan tables | Not seed v1; out of this thin cut |

**Default export omit:** `transitions`, `reviews`, capability rows, tool decisions, index, tokens, lock. Tasks export as entities (`id`, `title`, `body`, `goal_id`) without replaying work_state. Clone tasks start **PENDING**. Why-chains survive via **links**.

### Merge of two seed JSON PRs

Git will conflict on a single `trace/graph.json` if two branches edit it. **Do not** ship a custom merge driver in P17.

| Strategy | Decision |
|----------|----------|
| Append-only NDJSON / one-file-per-UUID | **Deferred** — more files, new format, not needed for thin cut |
| Single `graph.json` | **Lock** — same seed v1 import already understands |
| Semantic merge | **UUID upsert on `seed import`** after humans resolve or concatenate entity arrays |
| Conflict behavior (S03) | Document: prefer union of objects by `id`; same-id different title/body → last-import-wins (upsert), loud JSON summary; duplicate links → idempotent no-op (not UNIQUE fail) |

### Path convention

| Candidate | Verdict |
|-----------|---------|
| **`trace/graph.json`** | **Lock (recommended).** Distinct from gitignored `.trace/`; not mixed with `docs/` (protocol/design SoT) |
| `.trace-export/` | Reject as default — looks like local derived data; easy to gitignore by mistake |
| `docs/graph.json` | Reject as default — `docs/` is agent/human design, not live graph |
| `seed/gt.json` | Keep for **fixtures/evals** only (`fixtures/x0/seed/gt.json`) |

CLI: `trace seed export [-o <file>]` (stdout if no `-o`, matching machine stdout). Docs + AGENTS: write **`trace/graph.json`** before a PR that changes semantic graph objects.

`.gitignore` **stays** `.trace/`. Do not add `trace/` to gitignore.

---

## Loop 3 — Gaps, bugs, improvements, risks

### New DF IDs (DF-80+; no collision with P16 68–78; DF-79 unused)

| ID | Sev | Kind | Finding | Home |
|----|-----|------|---------|------|
| **DF-80** | med | gap | No `seed export`. Clone cannot reconstruct the causal graph from git. Import-only CLI. | P17 S01 |
| **DF-81** | med | bug/gap | Re-import is not idempotent: entity upsert exists, but `InsertLink` UNIQUE-fails; `Create*` re-appends `entity.created`; transitions replay can fail the work machine. | P17 S03 |
| **DF-82** | low | gap | No commit path / agent habit. Agents will `add`/`link` into `.trace/` and open PRs without exporting. | P17 S02 |
| **DF-83** | low | gap | Two PRs editing `trace/graph.json`: git conflict + undefined semantic merge. | P17 S03 (behavior + docs) |

### Risks (do not board as extra product)

| Risk | Disposition |
|------|-------------|
| Secrets in decision/task bodies | Docs: never put secrets in graph bodies (SECURITY §6). Optional secret scan later — **not** P17. Export does not read env or token file. |
| Prompt injection in committed graph | Law 4 already: retrieved graph text is data. S02 docs: same untrusted_data labeling as source. No new trust elevation. |
| Stale graph vs code (Law 18) | Clone recipe always `index` after import. Export does not include derived index. STALE stays a live-DB concern after code moves. |
| Seed import incomplete rels | **P16 S05** owns DF-70/73. P17 **depends**, does not re-implement. |
| Agent forgets to export | S02: AGENTS / CONTRIBUTING / help one-liner. No git hook required this phase (hooks not in skills map for P0). |
| Merge conflicts | S03 + docs; no merge driver. |
| Encryption | **Not required** for policy notes. Local token stays in `.trace/`. Reject encryption-as-git. |
| DF-28 overlap | Complementary: DF-28 is in-process Trace-pull. Portable JSON is what a **clone** Trace-pulls from. Do not add a handoff entity. Pointer files (`PLANNER-DONE.md`) stay optional; UUIDs in committed graph are the SoT. |

---

## Loop 4 — Cross-check (leverage)

### Goals progress (`TRACE-GOALS-PROGRESS-2026-08-17.md`)

H1 (understanding) is **done for a seeded tree**, not for a **bare clone**. Pillar “local-first graph + `.trace/` SQLite” is done; git-facing **projection** was never a ranked peer technique (ranks 4–20 were install/impact/honesty). That is a product hole relative to the README clone story, not a contradiction of remaining backlog.

Recommended sequence in that memo: (1) P14 impact+install **boarded**, (2) research S05 supersession, (3) `plan simulate`, (4) D21–D23. D21–D23 **already scored** and fed DF-70–74 into **P16 S05**. Portable graph is **not** a substitute for (2)/(3); it is **higher leverage as the next queued thin phase after P16** because:

1. It unblocks H1 for collaborators who never ran the original `init`.
2. It is small (export + docs + idempotent import) vs simulate/adopt-discard.
3. It reuses seed v1 instead of a new persistence model (Law 13).

### Similar projects

Peers keep a **local** graph (SQLite/C/TS). Obsidian’s “markdown as SoT” is explicitly **not** a Trace product target. Steal **one method**: a git-readable snapshot, **not** an Obsidian vault. Do not commit the DB (codebase-memory / graphify style indexes stay derived).

### D30 / D21

- **D30/D30b:** implementer failed when the greeting lived only in planner prose; retry **Trace-pull** passed. Clone collaboration without export is D30 all over again with an empty DB.
- **D21–D23:** combo seeds proved import is the plant path — and that import is **incomplete** (DF-70/73). P17 must not paper over that; it must export whatever S05 makes importable.

---

## Locked defaults for Phase 17 (planner)

| Item | Value |
|------|-------|
| Format | Seed JSON **v1** (round-trip import↔export) |
| Default path | `trace/graph.json` |
| `.gitignore` | Unchanged (`.trace/` only) |
| Export include | goals, tasks (id/title/body/goal_id), decisions, assumptions, discoveries, plan_changes, claims, evidence, links (all importable rels **after** P16 S05), optional findings/alternatives **if** S05 landed those keys |
| Export exclude | index, tokens, lock, tool decisions, capabilities, events, reviews, transitions/work_state (default) |
| Import | Idempotent UUID upsert; duplicate links no-op |
| MCP | **No** new tools (`seed` stays CLI) |
| Daemon / HTTP / encryption | Forbidden |
| DR-HANDOFF | Default **`no successor`** (P16 VERIFY also stays `no successor`; P17 is a **forward human queue**, not P16’s child) |
| Product Go | Forbidden on P17-00 |

---

## Phase recommendation

- **Slug:** `docs/phases/phase-17-portable-graph-git/`
- **GO** with implementation board (S01–S04).
- **Non-goals:** committing `.trace/`; encryption-as-git; daemon/sync; duplicating P16 S05; reviews/DONE in default export; NDJSON split; git merge driver.

---

## Addendum — 2026-08-17 forward (DF-84…86)

P17-00 FINAL + this investigation **excluded** planner coarse-plan tables from seed v1 and required **no** git hook. Upcoming S01–S04 **supersede** those excludes for a complete **local/git** portable graph. Historical loops above stay as the original GO. Live locks: [`docs/phases/phase-17-portable-graph-git/DF-84-FORWARD.md`](../phases/phase-17-portable-graph-git/DF-84-FORWARD.md).

| Item | Upcoming lock |
|------|----------------|
| Plan tree | Export/import `plan_phases` / `plan_scopes` / `scope_deep_plans` / `goal_plan_state` (mig 006). UUIDs stay identity |
| `exported_at_commit` | Git SHA is snapshot **evidence**, not identity |
| Attribution | Git author + SHA = evidence. `transition.actor` is not auth |
| Hook | `trace install git-hook` CONDITIONAL/deferred; must not wrap `git commit`; **not blocking VERIFY** |
| Completion | Two clones, no shared `.trace/`, offline: `init` + `seed import` + `index` + `why`/`context` + plan |
| Non-goals (tightened) | No HTTP; no daemon; no hosted MCP; no OAuth; do not point this repo’s `trace-mcp` at the internet. Hosted product = TODO Later developments, separate repo |

Next free DF after this addendum: **DF-87**.
