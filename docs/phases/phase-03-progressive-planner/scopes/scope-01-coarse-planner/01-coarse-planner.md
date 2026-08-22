# P03 / S01 / 01 — Coarse progressive planner

## Metadata
- id: P03-S01-01
- todo_ids: [P03-S01-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Implement a **minimal coarse progressive planner**: persist **goal→phase→scope** hierarchy, deep-plan **only the current scope** (+ one lookahead shallow summary), and expose a library + thin CLI. Caller-supplied structure only — **no** LLM auto-generation of an entire backlog (`J` REJECT). Keep honesty / p0x / x0 green. No daemon/HTTP/embeddings. MCP not required.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) locks (this scope)
- [phase README](../../README.md)
- [docs/init/J_BRAINSTORMING_OUTCOMES.md](../../../../init/J_BRAINSTORMING_OUTCOMES.md) — deep-plan current + one lookahead; supersede not delete
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G16 (churn budget **column only** here; S02 enforces)
- Live: `internal/domain` (Goal/Task/Discovery/PlanChange), `internal/store` (mig ≤005), `cmd/trace` add/link

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | Keep `go.mod` floor (currently 1.24.0); do not downgrade |
| **Package path** | **`internal/planner`** — progressive planning orchestration. **Not** a domain entity dump; do **not** put CreatePhase/DeepPlan into `internal/domain` |
| Domain role | Keep using `domain` for Goal (and existing Task/Discovery/PlanChange). Planner **reads** Goal via store/domain; does not fork Goal CRUD |
| Store / DB | One DB: `projectRoot/.trace/trace.db` via `*store.Store`. **No** second database |
| Migration | Additive embed **`006_plan_hierarchy.sql`** (do not rewrite `001`–`005`) |
| Hierarchy | **Goal → plan_phases → plan_scopes**. Ordering via `ord` INTEGER (0-based, stable within parent) |
| Deep plan | Table **`scope_deep_plans`**: one **ACTIVE** revision per scope; prior revisions → **SUPERSEDED** (never DELETE rows for recovery) |
| Current pointer | Table **`goal_plan_state`**: `goal_id` PK → `current_scope_id` (nullable until set) |
| Lookahead | Next scope after current by `(phase.ord ASC, scope.ord ASC)` within the same goal; deep-plan API may write/update a **shallow** summary on that next scope only — **not** a full deep plan for the whole goal |
| Churn column | `plan_scopes.auto_replan_count INTEGER NOT NULL DEFAULT 0` — S01 initializes/preserves; **S02** owns budget/ack enforcement (DR-CHURN N=5) |
| Tasks | **Do not** require `tasks.scope_id` this scope. Existing Goal↔Task (`goal_id`) unchanged. Optional future link is out |
| CONFLICT edges | **Out** of S01 (J ADOPT deferred to later / S02 if measured) |
| Events | Append thin events via store when useful: e.g. `plan.coarse_created`, `plan.current_set`, `plan.deep_planned`, `plan.deep_superseded` (payload: ids + actor). Prefer existing event table patterns |
| Provenance | Phase/scope/revision rows use status `ACTIVE` \| `STALE` \| `SUPERSEDED` (same vocabulary as domain entities). Prefer supersede over delete |
| CLI | Thin G19: `trace plan` subcommands (stdlib argv, no cobra). No business logic in `cmd/trace` |
| MCP | **Not** required for S01 (CLI primary). Do not add MCP plan tools here |
| CGO | `internal/planner` + store mig APIs must pass `CGO_ENABLED=0` |
| Carry-forward bars | honesty Paths A/B/C; p0x 7/7; x0; Gate C artifacts **untouched** |
| Out | Full replan/churn demo (S02); Gate E harness (S03); daemon/HTTP/embeddings; LLM backlog generation; rewriting Mode-B Gate C packs |

### Schema (locked shape — column names may vary slightly; tables/semantics locked)

```sql
-- 006_plan_hierarchy.sql (additive)

CREATE TABLE IF NOT EXISTS plan_phases (
  id TEXT PRIMARY KEY,
  goal_id TEXT NOT NULL,
  title TEXT NOT NULL,
  body TEXT NOT NULL DEFAULT '',
  ord INTEGER NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_plan_phases_goal ON plan_phases(goal_id, ord);

CREATE TABLE IF NOT EXISTS plan_scopes (
  id TEXT PRIMARY KEY,
  phase_id TEXT NOT NULL,
  title TEXT NOT NULL,
  body TEXT NOT NULL DEFAULT '',           -- coarse summary only
  ord INTEGER NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'ACTIVE',
  auto_replan_count INTEGER NOT NULL DEFAULT 0,  -- S02 hook
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_plan_scopes_phase ON plan_scopes(phase_id, ord);

CREATE TABLE IF NOT EXISTS scope_deep_plans (
  id TEXT PRIMARY KEY,
  scope_id TEXT NOT NULL,
  content_json TEXT NOT NULL DEFAULT '{}',  -- DeepPlanDocument JSON
  status TEXT NOT NULL DEFAULT 'ACTIVE',    -- ACTIVE | SUPERSEDED | STALE
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_scope_deep_plans_scope ON scope_deep_plans(scope_id, status);

CREATE TABLE IF NOT EXISTS goal_plan_state (
  goal_id TEXT PRIMARY KEY,
  current_scope_id TEXT,
  updated_at TEXT NOT NULL
);
```

### DeepPlanDocument JSON (locked fields)

```json
{
  "scope_id": "<uuid>",
  "exit_criteria": ["..."],
  "constraints": ["..."],
  "work_items": [{"title": "...", "notes": ""}],
  "lookahead_scope_id": "<uuid or empty>",
  "lookahead_summary": "shallow one-liner or empty"
}
```

- `work_items` are **caller-supplied** (CLI flags / library input / test fixtures). Implementer must **not** invent a goal-wide backlog via LLM or heuristics that expand past current+lookahead.
- Lookahead summary is shallow text only; do not write a second full deep-plan revision for the lookahead scope in the same call unless explicitly superseding that scope later.

### Minimum public API (`internal/planner`)

```text
New(st *store.Store) *Service
  // Optionally accept *domain.Service for Goal existence checks — either OK if documented.
  // Must not import cmd/trace or internal/mcp.

CreateCoarsePlan(ctx, CoarsePlanInput) (CoarsePlan, error)
  // Input: GoalID (must exist), Phases []{Title, Body?, Scopes []{Title, Body?}}
  // Creates plan_phases + plan_scopes with ord = index order.
  // Rejects empty GoalID / empty phase title / empty scope title.
  // Does NOT deep-plan all scopes. Does NOT create Tasks.
  // Initializes goal_plan_state row if missing (current_scope_id NULL).

SetCurrentScope(ctx, goalID, scopeID string) error
  // Validates scope belongs to goal (via phase.goal_id). Updates goal_plan_state.

DeepPlan(ctx, DeepPlanInput) (DeepPlanResult, error)
  // Input: ScopeID, ExitCriteria, Constraints, WorkItems, optional LookaheadSummary override
  // Requires ScopeID == current scope for that goal (fail closed otherwise) — progressive rule.
  // Writes/replaces ACTIVE scope_deep_plans for ScopeID (supersede prior ACTIVE → SUPERSEDED).
  // Resolves lookahead = next scope in goal order; sets lookahead_* on document; may set
  //   plan_scopes.body (coarse) on lookahead to LookaheadSummary if provided / non-empty.
  // Must NOT deep-plan every scope under the goal.

GetPlan(ctx, goalID string) (PlanView, error)
  // Returns phases/scopes (ACTIVE), current_scope_id, ACTIVE deep plan for current (if any),
  // and shallow lookahead identity/summary.

SupersedeDeepPlan(ctx, SupersedeInput) (DeepPlanResult, error)
  // S02 hook: mark ACTIVE deep plan for ScopeID SUPERSEDED; write new ACTIVE revision
  // from caller-supplied document fields. Does not require scope==current (replan may
  // target a scope after discovery). Does not implement severity/churn budget (S02).

ListScopes(ctx, goalID string) ([]ScopeRef, error)  // ordered; helper for tests/S02
GetCurrentScope(ctx, goalID string) (ScopeRef, error) // ErrNotFound / clear error if unset
```

Names may vary slightly; **behavior** above is locked.

### Progressive rule (fail closed)

```text
CreateCoarsePlan → ordered phases/scopes (coarse titles only)
SetCurrentScope(scope_k)
DeepPlan(scope_k, …)  → full document for k + shallow lookahead for scope_{k+1} only
DeepPlan(scope_j≠current) → error
# Whole-goal auto backlog → forbidden
```

### CLI (thin G19)

```text
trace plan create-coarse --goal <id> --phase <title> [--scope <title> ...]
  # Repeatable --phase; each --phase consumes following --scope flags until next --phase
  # OR accept a small JSON file --from <path> with the CoarsePlanInput shape — pick one UX and document in help
trace plan set-current --goal <id> --scope <id>
trace plan deep --scope <id> --exit <text> [--constraint <text>] [--work <title>] ...
  # Maps flags → DeepPlanInput; require current-scope rule via library
trace plan show --goal <id>
  # JSON PlanView on stdout (machine); progress on stderr
```

Exit codes: inherit CLI convention `0/1/2`. Wire help/version. No cobra.

### Target tree

```text
internal/store/
  schema/006_plan_hierarchy.sql
  plan_hierarchy.go          # Upsert/Get/List phases, scopes, deep plans, goal_plan_state
  # extend migrate embed list

internal/planner/
  doc.go
  service.go                 # New + APIs above
  types.go                   # CoarsePlanInput, DeepPlanDocument, PlanView, …
  planner_test.go            # required tests below

cmd/trace/
  plan.go                    # thin subcommands
  help.go / root.go          # wire `plan`
```

### Tests (required)

- CreateCoarsePlan persists phases/scopes with stable `ord`; rejects missing goal / empty titles
- SetCurrentScope + GetCurrentScope round-trip; rejects scope from another goal
- DeepPlan succeeds only for current scope; writes ACTIVE revision; prior ACTIVE → SUPERSEDED on re-deep
- DeepPlan result includes lookahead_scope_id when a next scope exists; no deep revisions created for non-current scopes in that call
- SupersedeDeepPlan creates new ACTIVE + supersedes old (S02 hook smoke)
- `CGO_ENABLED=0 go test ./internal/planner/... ./internal/store/...`
- Regression: `CGO_ENABLED=0 go test ./evals/honesty/...` and `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/...` and `CGO_ENABLED=1 go test ./...`

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] `internal/planner` exists with locked APIs; package is **not** `internal/domain` extension for Phase/Scope CRUD
- [ ] Mig `006_plan_hierarchy.sql` embedded; tables persist phase/scope/deep-plan/current pointer
- [ ] Deep-plan enforces current-scope + one lookahead only (tested)
- [ ] SupersedeDeepPlan available for S02 (tested smoke)
- [ ] Thin `trace plan` CLI wired (G19); help updated
- [ ] `CGO_ENABLED=0` planner+store PASS; honesty + p0x + x0 + `./...` PASS as above
- [ ] No daemon/HTTP/embeddings/MCP plan tools; no LLM whole-backlog generation
- [ ] Board status + Notes only

## Minimal todos
- [ ] Add store mig 006 + store helpers
- [ ] Implement `internal/planner` + unit tests
- [ ] Thin `trace plan` CLI
- [ ] Self-check exit criteria + bars
- [ ] Mark P03-S01-01 done with Notes
