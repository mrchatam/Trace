# P03 / S02 / 01 — Discovery→PlanChange replan demo

## Metadata
- id: P03-S02-01
- todo_ids: [P03-S02-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Implement **discovery→PlanChange replan with churn controls** and a planted-discovery demo: severity-gated auto-replan (`PLAN_AFFECTING`+ only), G16 / DR-CHURN budget (N=5) with human ack, consuming live S01 `internal/planner` APIs. Keep honesty / p0x / x0 green. Do not rewrite Mode-B Gate C packs. No daemon/HTTP/embeddings. MCP not required.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) locks (this scope)
- [phase README](../../README.md)
- [docs/init/J_BRAINSTORMING_OUTCOMES.md](../../../../init/J_BRAINSTORMING_OUTCOMES.md) — severity tiers; churn budget; supersede not delete
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G16
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-CHURN (N=5)
- Live S01: `internal/planner` (`SupersedeDeepPlan`, `GetCurrentScope`/`ListScopes`/`GetPlan`); mig `006` `auto_replan_count` column (no mutator yet)
- Live causal: `domain.LinkDiscoveryPlanChange` (`discovery_causes_plan_change`)

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | Keep `go.mod` floor (currently 1.24.0); do not downgrade |
| **Package path** | Stay in **`internal/planner`** (+ store/domain helpers). **Do not** invent a second planner package or put replan orchestration under `internal/domain` |
| S01 consume | Call existing `SupersedeDeepPlan`, `GetCurrentScope`/`ListScopes`/`GetPlan`. Do **not** fork supersede logic |
| Causal link | Keep `domain.LinkDiscoveryPlanChange`; do not invent a parallel link rel |
| Store / DB | One DB via `*store.Store`. Additive mig **`007_discovery_severity.sql`** only (do not rewrite `001`–`006`) |
| Severity | `discoveries.severity TEXT NOT NULL DEFAULT 'INFO'`; allowed values **`INFO` \| `PLAN_AFFECTING` \| `BLOCKING`** (reject others fail-closed) |
| Auto-replan gate | Only **`PLAN_AFFECTING` and `BLOCKING`** auto-open replan (`J` ADOPT). **`INFO`** may create/link a PlanChange but **must not** call `SupersedeDeepPlan` or increment `auto_replan_count` |
| Churn budget | **N=5** default (`DefaultMaxAutoReplans = 5`, overridable in tests). Fail closed when `auto_replan_count >= N` before auto-replan. Human **ack resets count to 0** |
| Store mutators | `IncrementAutoReplanCount(scopeID)` and `AckAutoReplan(scopeID)` (reset to 0). S01 only DEFAULT 0 + preserve on body update — S02 owns enforce |
| Apply API | `planner.ApplyDiscoveryReplan` (name locked) — see Minimum public API |
| Demo harness | **`evals/replan`** package; named test **`TestPlantedDiscoveryReplan`** (honesty-style library demo) |
| DPC retrieval | **Out** — keep GC-01 global attach; do not scope Why/TaskContext DPC unless a later measured residual demands it |
| CLI | Thin G19: `trace add discovery --severity`; `trace plan apply-discovery`; `trace plan ack-replan`. No business logic in `cmd/trace` |
| MCP | **Not** required. Do not add MCP plan/replan tools |
| CGO | planner + store + `evals/replan` must pass `CGO_ENABLED=0` where applicable |
| Carry-forward bars | honesty Paths A/B/C; p0x 7/7; x0; Gate C artifacts **untouched** |
| Out | Gate E final bar / VERIFY (S03); daemon/HTTP/embeddings; LLM backlog generation; Mode-B Gate C pack rewrites; scoping DPC attach |

### Schema (locked)

```sql
-- 007_discovery_severity.sql (additive)

ALTER TABLE discoveries ADD COLUMN severity TEXT NOT NULL DEFAULT 'INFO';
-- Enforce INFO | PLAN_AFFECTING | BLOCKING in store/domain validation (SQLite CHECK optional).
```

Existing `plan_scopes.auto_replan_count` (mig 006) is the churn counter — **no new column** for ack; ack = reset to 0.

### Severity constants (locked names)

```text
SeverityINFO           = "INFO"
SeverityPlanAffecting  = "PLAN_AFFECTING"
SeverityBlocking       = "BLOCKING"
```

Place constants in `internal/domain` and/or `internal/planner` (one canonical set; re-export OK). Default when omitted: `INFO`.

### Minimum public API

#### Store (`internal/store`)

```text
IncrementAutoReplanCount(scopeID string) (newCount int, error)
  // Atomically increment plan_scopes.auto_replan_count by 1; return new value.
  // Fail if scope missing.

AckAutoReplan(scopeID string) error
  // Set auto_replan_count = 0. Fail if scope missing.

// Extend Discovery struct + UpsertDiscovery/GetDiscovery to include Severity.
// Existing rows after mig: DEFAULT 'INFO'.
```

#### Domain (`internal/domain`)

```text
// DiscoveryInput gains Severity string (optional; default INFO; validate enum).
CreateDiscovery(ctx, DiscoveryInput) (store.Discovery, error)
  // Persist severity; reject unknown severity values.

// Optional helper OK:
SetDiscoverySeverity(ctx, discoveryID, severity string) error
```

#### Planner (`internal/planner`)

```text
const DefaultMaxAutoReplans = 5

ApplyDiscoveryReplan(ctx, ApplyDiscoveryReplanInput) (ApplyDiscoveryReplanResult, error)

// ApplyDiscoveryReplanInput (locked fields; names may vary slightly):
//   DiscoveryID string          // required
//   ScopeID string              // required — deep plan target (need not be current)
//   PlanChangeID string         // optional: if empty, create PlanChange from Title/Body overrides
//   PlanChangeTitle, PlanChangeBody string  // used when creating PlanChange
//   ExitCriteria, Constraints, WorkItems    // passed through to SupersedeDeepPlan
//   LookaheadScopeID, LookaheadSummary      // optional pass-through
//   MaxAutoReplans int          // 0 → DefaultMaxAutoReplans
//   Actor string

// Behavior (fail closed):
// 1. Load discovery; validate exists.
// 2. If severity == INFO:
//      - Optionally ensure PlanChange exists + LinkDiscoveryPlanChange when PlanChange fields/ID provided
//      - Do NOT SupersedeDeepPlan; do NOT IncrementAutoReplanCount
//      - Return result with AutoReplanApplied=false, reason severity_info
// 3. If severity is PLAN_AFFECTING or BLOCKING:
//      a. Load scope; if auto_replan_count >= MaxAutoReplans → error ErrReplanBudgetExceeded (or equivalent)
//      b. Ensure PlanChange (create if needed) + LinkDiscoveryPlanChange(discovery, planChange)
//      c. Call SupersedeDeepPlan with caller-supplied document fields
//      d. IncrementAutoReplanCount(scopeID)
//      e. Return AutoReplanApplied=true + revision id + new count
// 4. Unknown severity → validation error

AckReplan(ctx, scopeID string) error
  // Thin wrapper over store.AckAutoReplan; append thin event e.g. plan.replan_acked
```

Consume S01: do **not** reimplement supersede-not-delete; call `SupersedeDeepPlan`.

### CLI (thin G19)

```text
trace add discovery --title <t> [--severity INFO|PLAN_AFFECTING|BLOCKING] …
  # Default severity INFO when flag omitted

trace plan apply-discovery --discovery <id> --scope <id>
                [--plan-change <id>] [--pc-title <t>] [--pc-body <b>]
                [--exit <text>] [--constraint <text>] [--work <title>] …
  # Maps to ApplyDiscoveryReplan; library enforces severity + budget

trace plan ack-replan --scope <id>
  # Maps to AckReplan → reset auto_replan_count to 0

# Existing: create-coarse | set-current | deep | show — keep; extend help
```

Exit codes: inherit `0/1/2`. Stdout machine JSON where applicable; progress on stderr. No cobra.

### Target tree

```text
internal/store/
  schema/007_discovery_severity.sql
  entities_causal.go          # Discovery.Severity + Upsert/Get
  plan_hierarchy.go           # IncrementAutoReplanCount, AckAutoReplan
  # extend migrate embed list

internal/domain/
  create.go / …               # DiscoveryInput.Severity + validation

internal/planner/
  replan.go (or service.go)   # ApplyDiscoveryReplan, AckReplan, DefaultMaxAutoReplans
  types.go                    # ApplyDiscoveryReplanInput/Result, ErrReplanBudgetExceeded
  replan_test.go              # unit: INFO no supersede; PLAN_AFFECTING supersede+count; budget fail; ack reset

evals/replan/
  doc.go
  replan_test.go              # TestPlantedDiscoveryReplan

cmd/trace/
  add.go                      # --severity on discovery
  plan.go                     # apply-discovery + ack-replan + help
```

### Tests (required)

**Unit (`internal/planner` / store / domain):**
- CreateDiscovery with severity; default INFO; reject garbage severity
- INFO ApplyDiscoveryReplan: no new ACTIVE deep-plan revision; count unchanged
- PLAN_AFFECTING ApplyDiscoveryReplan: LinkDiscoveryPlanChange present; SupersedeDeepPlan new ACTIVE; count += 1
- BLOCKING behaves like PLAN_AFFECTING for auto-replan gate
- After N=5 successful auto-replans, 6th fails closed (`ErrReplanBudgetExceeded`); AckReplan resets to 0; then Apply succeeds again
- SupersedeDeepPlan still works when called directly (S01 smoke retained)

**Demo (`evals/replan`):**
- `TestPlantedDiscoveryReplan`: plant Goal → CreateCoarsePlan → SetCurrentScope → DeepPlan; plant PLAN_AFFECTING discovery → ApplyDiscoveryReplan → assert deep plan superseded + count; plant INFO discovery → assert no supersede; exhaust budget → fail → ack → succeed once

**Bars:**
- `CGO_ENABLED=0 go test ./internal/planner/... ./internal/store/... ./internal/domain/... ./evals/replan/...`
- `CGO_ENABLED=0 go test ./evals/honesty/...`
- `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/...`
- `CGO_ENABLED=1 go test ./...`

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] Mig `007_discovery_severity.sql` embedded; Discovery severity persisted + validated
- [ ] Store `IncrementAutoReplanCount` + `AckAutoReplan` (reset to 0); budget enforce on Apply
- [ ] `planner.ApplyDiscoveryReplan` + `AckReplan`; INFO no auto-replan; PLAN_AFFECTING+ auto-replan via `SupersedeDeepPlan` + `LinkDiscoveryPlanChange`
- [ ] `evals/replan` `TestPlantedDiscoveryReplan` PASS
- [ ] Thin CLI: `add discovery --severity`, `plan apply-discovery`, `plan ack-replan` (G19)
- [ ] Required test bars PASS; Gate C packs / honesty / p0x / x0 intact
- [ ] No second planner package; no DPC attach scoping; no daemon/HTTP/embeddings/MCP replan tools
- [ ] Board status + Notes only

## Minimal todos
- [ ] Add mig 007 + Discovery.Severity + store Increment/Ack mutators
- [ ] Implement ApplyDiscoveryReplan + AckReplan + unit tests
- [ ] Add evals/replan TestPlantedDiscoveryReplan
- [ ] Thin CLI flags/subcommands + help
- [ ] Self-check exit criteria + bars
- [ ] Mark P03-S02-01 done with Notes
