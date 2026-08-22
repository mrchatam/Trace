# P00 / S05 / 01 — Work/causal API

## Metadata
- id: P00-S05-01
- todo_ids: [P00-S05-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Implement the **canonical work/causal library API** in `internal/domain` on top of the live S02 store: create/link **Goal, Decision, Assumption, Task, Discovery, PlanChange** with **required provenance**, a **Task work-state machine** (actor/reason + append-only events), and **light Claim/Evidence stubs** (no full honesty promotion — Phase 01 / T006). Enables P0-X #1 (Goal/Task/Decision/Discovery round-trip with provenance) and substrate for later `why` / seed / CLI.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [C_FIRST_SCOPE.md](../../../../init/C_FIRST_SCOPE.md) — work/causal subset + P0-X #1
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G5 provenance, G8 discovery→PlanChange, G18 thin events
- [D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-EVT, DR-CLAIM, DR-P0X, DR-SURFACE
- [PROJECT_MODEL.md](../../../../PROJECT_MODEL.md) — entities, provenance, task states, Decision/Discovery
- [B_INITIAL_BOARD.md](../../../../init/B_INITIAL_BOARD.md) — historical T005 (exit: Goal→Task, Decision→Task, Discovery→PlanChange + illegal transitions)
- Live priors: S02 `UpsertGoal`/`GetGoal`, `UpsertTask`/`GetTask`, `AppendEvent`/`ListEventsByEntity`; schema tables already exist for decisions/assumptions/discoveries/plan_changes/claims/evidence/reviews; `internal/domain/doc.go` stub only; `go.mod` **go 1.24.0**

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | `go 1.24.0` in `go.mod` (do **not** downgrade) |
| Package path | `internal/domain` (package `domain`) — **only** product package for this scope’s public API |
| Persistence | **Only** via `*store.Store` in the same `.trace/trace.db`. **Do not** open a second DB, fork schema under domain, or store causal facts only in memory |
| Construction | `domain.New(st *store.Store) *Service` (name may vary; must take store). Domain does **not** call `store.Open` |
| Additive store OK | Upsert/Get for Decision, Assumption, Discovery, PlanChange (+ thin Claim/Evidence); `entity_links` helpers; `tasks.work_state` column — all in `internal/store` + embedded migration. Prefer mirroring Goal/Task patterns |
| Forbidden store fork | Parallel SQLite files; rewriting S02 Goal/Task/Event semantics; source BLOBs |
| Provenance (G5) | Every create/update of semantic facts must set `source_type`, `confidence`, `status`, `created_at`/`updated_at`; `last_verified_at` optional. Default `source_type` when caller omits: `USER_ASSERTED`. Empty title rejected |
| Provenance `status` | `ACTIVE` \| `STALE` \| `SUPERSEDED` on Goal/Decision/Assumption/Discovery/PlanChange/Claim/Evidence (S02 vocabulary). Use `MarkStale` (or equivalent) rather than silent overwrite of history |
| Task **work** state | **Separate** column `tasks.work_state` (migration). **Do not** overload provenance `status` with work-machine enums |
| `work_state` vocabulary | `PENDING` \| `IN_PROGRESS` \| `AWAITING_REVIEW` \| `BLOCKED` \| `FAILED` \| `DONE` \| `STALE` \| `SKIPPED` (PROJECT_MODEL §5). Default on create: `PENDING` |
| Decision workflow enum | Full `PROPOSED`/`ACCEPTED`/`REJECTED` **deferred**. P0-X Decisions use provenance `status=ACTIVE` (accepted fact). Rationale/body in `body` TEXT |
| Goal→Task link | Primary: `tasks.goal_id` (already in schema). Domain `CreateTask` / `LinkGoalTask` must set it. Do **not** require a duplicate `entity_links` row for Goal→Task |
| Other causal links | Table `entity_links` (migration): `id`, `from_type`, `from_id`, `rel`, `to_type`, `to_id`, provenance-ish `source_type`/`confidence`/`created_at`, **UNIQUE(from_type, from_id, rel, to_type, to_id)** |
| Required `rel` values | `decision_affects_task`, `discovery_causes_plan_change` (exact strings). Optional extras OK if documented in package comment |
| Entity type strings | Lowercase singular: `goal`, `task`, `decision`, `assumption`, `discovery`, `plan_change`, `claim`, `evidence`, `review` |
| Events (DR-EVT) | Every successful **create**, **link**, and **work_state transition** appends an event via `store.AppendEvent`. Transition payload JSON must include `actor`, `from`, `to`, `reason` (and `evidence_ids` array, may be empty) |
| Event type names | `entity.created`, `entity.linked`, `task.transition` (locked). `entity_type`/`entity_id` point at the primary subject (for links: the `from` entity) |
| DONE policy stub | Default: reject `→ DONE` unless `TransitionOptions.AllowDoneWithoutReview == true` **or** `EvidenceIDs` non-empty. Full Review/Claim promotion = Phase 01 — do not implement VerifiedFact |
| Claim/Evidence | Light stubs only: create + optional link (`claim_has_evidence` rel). **No** Review PASS/FAIL promotion engine |
| CGO | Domain + any new store APIs must pass `CGO_ENABLED=0`. Do not import `analyzers` / `gitcli` |
| Surface | Library only — **no** new `cmd/trace` subcommands (S07); **no** MCP/daemon/HTTP |
| Out of scope | Autoplanning; impact engine; Requirement/Constraint/Phase/Scope entities; full honesty path; silent history rewrite |

### Task work_state transition graph (locked)

Legal edges (reject all others with a typed/clear error):

```text
PENDING         → IN_PROGRESS | BLOCKED | SKIPPED
IN_PROGRESS     → AWAITING_REVIEW | BLOCKED | FAILED | DONE† | PENDING
AWAITING_REVIEW → DONE† | IN_PROGRESS | FAILED | BLOCKED
BLOCKED         → PENDING | IN_PROGRESS | SKIPPED | FAILED
FAILED          → PENDING | IN_PROGRESS | SKIPPED
DONE            → STALE | PENDING   # reopen / correction only via explicit transition
SKIPPED         → PENDING
STALE           → PENDING | SUPERSEDED-via-MarkStale on provenance status (work_state may stay STALE)
```

† `DONE` also gated by DONE policy stub above.

Transition must record actor + reason (non-empty strings). Append `task.transition` event; update `work_state` + `updated_at`.

### Minimum public API (names may vary slightly; behavior locked)

```text
New(st *store.Store) *Service

// Create* — allocate UUID if empty; apply provenance defaults; persist; append entity.created
CreateGoal(ctx, GoalInput) (store.Goal /* or domain view */, error)
CreateDecision(ctx, DecisionInput) (...)
CreateAssumption(ctx, AssumptionInput) (...)
CreateTask(ctx, TaskInput) (...)          // TaskInput may include GoalID
CreateDiscovery(ctx, DiscoveryInput) (...)
CreatePlanChange(ctx, PlanChangeInput) (...)

// Links
LinkGoalTask(ctx, goalID, taskID, meta LinkMeta) error
  // sets tasks.goal_id; appends entity.linked (rel may be recorded as "goal_has_task" in event payload only;
  // persistence of Goal→Task is goal_id, not entity_links)
LinkDecisionTask(ctx, decisionID, taskID, meta LinkMeta) error
  // entity_links rel=decision_affects_task
LinkDiscoveryPlanChange(ctx, discoveryID, planChangeID, meta LinkMeta) error
  // entity_links rel=discovery_causes_plan_change

// Task machine
TransitionTask(ctx, taskID, toWorkState string, opts TransitionOptions) error
// TransitionOptions: Actor, Reason string; EvidenceIDs []string; AllowDoneWithoutReview bool

MarkStale(ctx, entityType, entityID, reason string) error  // sets provenance status=STALE + event

// Reads (thin wrappers OK)
GetGoal / GetTask / GetDecision / GetDiscovery / … as needed for tests
ListLinksFrom(ctx, entityType, entityID) ([]Link, error)

// Claim/Evidence stubs (minimal)
CreateClaim / CreateEvidence / LinkClaimEvidence (rel=claim_has_evidence) — enough for compile + one smoke test; no promotion
```

Domain may re-export or wrap `store` types; prefer **not** duplicating large struct trees — wrapping store types is fine.

### Store migration (locked shape)

New embedded migration after `002_vcs_index.sql`, e.g. `003_causal_domain.sql`:

1. `ALTER TABLE tasks ADD COLUMN work_state TEXT NOT NULL DEFAULT 'PENDING';` (SQLite; idempotent apply via version table only — do not re-run ALTER on remigrate)
2. `entity_links` table as above + index on `(to_type, to_id)` and `(from_type, from_id)`

Existing rows: `work_state` default `PENDING`. Domain create path always sets it explicitly.

### Target tree

```text
internal/domain/
  doc.go                 # package contract: causal API; provenance; work_state machine; no second DB
  service.go             # New + Service
  create.go              # Create* helpers
  link.go                # Link*
  task_state.go          # TransitionTask + legal graph + DONE policy stub
  claim_stub.go          # optional thin Claim/Evidence
  domain_test.go         # Goal→Task, Decision→Task, Discovery→PlanChange, illegal transition, provenance

internal/store/
  schema/003_causal_domain.sql   # work_state + entity_links
  entities_causal.go             # Upsert/Get Decision/Assumption/Discovery/PlanChange (+ Claim/Evidence stubs)
  links.go                       # InsertLink / ListLinks*
  # extend Task struct + UpsertTask/GetTask for work_state
```

### Out of scope (this row)

- CLI entity-add commands (S07)
- Retrieval `why` / context compiler (S06)
- Fixture GT seed files (S08) — use temp `store.Open` in domain tests
- Full Claim→Evidence→Review→VerifiedFact (Phase 01)
- Autoplanning / impact simulation
- MCP / daemon / HTTP

## Board rights
Implementer: update **status + notes only** on `P00-S05-01`. Do not spawn rows or rewrite later prompts.

## Exit criteria
- [ ] `domain.New(st)` + Create* for Goal, Decision, Assumption, Task, Discovery, PlanChange persist via store and round-trip with provenance (`source_type`, `confidence`, `status`)
- [ ] Goal→Task via `goal_id`; Decision→Task via `entity_links` (`decision_affects_task`); Discovery→PlanChange via `entity_links` (`discovery_causes_plan_change`) — each covered by a test
- [ ] `TransitionTask` enforces the locked graph; at least one illegal transition test fails closed; legal transition appends `task.transition` event with actor/from/to/reason
- [ ] Default policy rejects `→ DONE` without `AllowDoneWithoutReview` and without evidence IDs; with flag or evidence IDs, DONE succeeds
- [ ] Migration `003` applies on `store.Open`; `tasks.work_state` present; no source-content BLOBs introduced
- [ ] Claim/Evidence stub create (and optional link) compiles and has at least one smoke test — no promotion engine
- [ ] `CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/...` passes; `go test ./...` green (analyzers may need CGO as today)
- [ ] No MCP/daemon/HTTP; no new CLI commands; domain does not import `analyzers`/`gitcli`
- [ ] TODO.md Notes for `P00-S05-01` updated; status `done`

## Minimal todos
- [ ] Store migration `003`: `work_state` + `entity_links` + Upsert/Get for Decision/Assumption/Discovery/PlanChange (+ Claim/Evidence stubs) + link helpers; extend Task for `work_state`
- [ ] `domain.New` + Create* with provenance defaults + `entity.created` events
- [ ] LinkGoalTask / LinkDecisionTask / LinkDiscoveryPlanChange + tests
- [ ] TransitionTask graph + DONE policy stub + event payload tests
- [ ] MarkStale + illegal transition rejection tests
- [ ] Claim/Evidence light stub smoke test
- [ ] Board status + notes (migration version, test commands)
