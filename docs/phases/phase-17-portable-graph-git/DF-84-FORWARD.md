# Forward correction — DF-84…86 (local/git portable graph complete)

P17-00 `00-PHASE-PLANNER.md` is **done history**. It locked a **narrower** seed: causal entities + links only; **excluded** planner coarse-plan tables; **no** `exported_at_commit`; **no** git hook this phase. A subsequent human/planner cut (complete **local/git** portable graph, still **no server**) **supersedes those excludes for upcoming S01–S04 only**. Do not rewrite the P17-00 prompt body beyond a pointer; do not claim P17-00 was wrong; do not steal P16 DFs 68, 70–78.

Live SoT for upcoming scopes: this file + phase [README.md](README.md). P17-00 FINAL remains the historical GO for DF-80…83.

## Why forward (not rewrite)

Clone-readable collaboration is incomplete if `trace/graph.json` omits the progressive-planning tree (`plan_phases` / `plan_scopes` / `scope_deep_plans` / `goal_plan_state` — mig **006**). A clone can `why` a decision and still cannot `plan show` without the original `.trace/` DB. Snapshot versioning and attribution were unspecified; agents will treat `transition.actor` as identity. Auto-export via hook is useful but **must not** wrap `git commit` or explode P17 into a fifth implement scope.

## New DF IDs (P17-owned; next free after DF-83)

| ID | Sev | Kind | Finding | Disposition | Home |
|----|-----|------|---------|-------------|------|
| **DF-84** | med | gap | Default seed omits plan hierarchy. P17-00 explicitly left planner tables out of seed v1. Clone cannot read the plan without `.trace/`. | **fix** | **S01** export keys + first import of those keys; **S03** idempotent upsert |
| **DF-85** | low | gap | No snapshot evidence field; SHA vs UUID identity unspecified; `actor` may be mistaken for auth. | **fix** | **S01** `exported_at_commit`; **S02** docs (git author+SHA = evidence; actor ≠ identity) |
| **DF-86** | low | gap | Agents forget to export (DF-82 docs help; still no auto-write). | **CONDITIONAL / deferred inside P17** | **Not blocking VERIFY.** Spec below; no S05 board row this cut |

P16 DFs stay on Phase 16. Encryption / `.trace/` commit remain **wontfix**. Hosted MCP remains **later / separate repo** (not a P17 DF).

## Live lock (upcoming S01–S04)

| Item | Value |
|------|-------|
| Path | **`trace/graph.json`** (unchanged) |
| Entity IDs | **UUIDs** (unchanged). Git commit SHA is **snapshot version**, not identity |
| `exported_at_commit` | Top-level string. Export fills `git rev-parse HEAD` when the bound root is a git repo; empty/omit when unknown. Import **accepts** the key and **must not** use it as an entity id or merge key |
| Export include | P17-00 causal set **plus** plan tree: `plan_phases`, `plan_scopes`, `scope_deep_plans`, `goal_plan_state` (mig 006 names). Plus P16 S05 `findings`/`alternatives` **if** those import keys landed. Links as today |
| Export exclude | **Unchanged:** index, tokens, lock, tool decisions, capabilities, events, reviews, transitions/`work_state` (default). Tasks still `id/title/body/goal_id` without replaying DONE |
| Attribution | Git **author + SHA** (commit + `exported_at_commit`) are **evidence**. `transition.actor` / review actor / `as_operator` are **not** authentication (same class as DF-44 flag≠identity). Document in S02 so clones do not treat actor strings as identity |
| Hook (DF-86) | **`trace install git-hook`** — CONDITIONAL install target, **deferred**. Must **not** wrap/replace `git commit`. Writes `trace/graph.json` via `seed export` (post-commit or pre-push). CI docs as backup. **VERIFY must not fail** if the target is absent |
| Completion | Two independent clones, **no shared `.trace/`**, offline, no account, no server: `init` + `seed import trace/graph.json` + `index` + `why`/`context` + `plan show` (or equivalent) work from git JSON |
| MCP / HTTP | **No** `trace_seed` tool; **no** daemon; **no** hosted MCP; **no** OAuth; **do not** point this repo’s `trace-mcp` at the internet |
| DR-HANDOFF | Default still **`no successor`**. Hosted product is **not** Phase 18 |

### Plan-tree JSON keys (do not invent names)

Match store tables (`internal/store/schema/006_plan_hierarchy.sql`). IDs are TEXT UUIDs (same as `InsertPlanPhase` / `InsertPlanScope`).

| Key | Rows | Required fields on export/import |
|-----|------|----------------------------------|
| `plan_phases` | `plan_phases` | `id`, `goal_id`, `title`, `body`, `ord`, `status` |
| `plan_scopes` | `plan_scopes` | `id`, `phase_id`, `title`, `body`, `ord`, `status`, `auto_replan_count` |
| `scope_deep_plans` | `scope_deep_plans` | `id`, `scope_id`, `content_json` (string; store TEXT), `status` — export **all** revisions (ACTIVE + SUPERSEDED) |
| `goal_plan_state` | `goal_plan_state` | `goal_id`, `current_scope_id` (nullable) — required so clone `plan show` knows the cursor |

Timestamps (`created_at` / `updated_at`) optional on seed v1; import may let store allocate. Unknown **top-level** keys still rejected except the keys listed here plus P17-00/`exported_at_commit`/S05 keys.

### DF-86 — not blocking VERIFY if

All of the following hold **without** installing a git hook:

1. S01 named export/round-trip tests include plan tree + `exported_at_commit` (empty-OK outside git).
2. S03 named idempotent import covers plan-tree UUIDs + last-wins.
3. S04 two-clone recipe: `why` / `context` / plan hierarchy readable from `trace/graph.json` alone.

**Acceptance if later implemented (same phase pack, spawn-forward only):**

- Install target id **`git-hook`**, tier **CONDITIONAL** (opt-in; not STABLE Cursor).
- Hook runs `trace seed export -o trace/graph.json` (or library equivalent). Does **not** wrap `git commit`; does **not** `git add -A`; does **not** skip hooks (`--no-verify` still user’s git).
- Uninstall removes the hook fragment Trace owns.
- CI/docs backup: “export before PR” remains valid if the hook is missing.
- No HTTP, no daemon, no hosted MCP.

## Do not

- Rewrite P17-00 FINAL locks as if they always included the plan tree
- Board a hosted/authenticated MCP phase in this repo
- Steal or reorder Phase 16 rows
- Mark P17-S01 `in_progress` from this planner cut
- Commit `.trace/`; add a git merge driver; NDJSON split
