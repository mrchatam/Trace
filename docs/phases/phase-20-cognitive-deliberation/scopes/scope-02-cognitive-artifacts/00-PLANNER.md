# P20-S02-00 — Cognitive artifacts planner

## Metadata
- id: P20-S02-00
- todo_ids: [P20-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock persistence for Uncertainty/Question, Hypothesis, Assumption invalidation, Decision reconsideration. **Reuse** Discovery, Decision, DecisionAlternative, Assumption. **No product Go this row.**

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [COVERAGE.md](../../COVERAGE.md) merge table (unchanged)
- Live: `internal/domain/create.go`, `internal/domain/task_state.go` (`MarkStale`), `internal/domain/impact.go` (`FindingKindInvalidatedAssumption`, `AddImpactFinding`), `internal/domain/deliberation.go` (`ApplyDeliberationTransition`), `internal/store/entities_causal.go`, `internal/store/entities.go` (`StatusStale`/`StatusSuperseded`), `internal/store/impact.go`, `internal/store/schema/015_deliberation_state.sql`, `internal/deliberation/types.go` (`PolicyInputs.BlockingUncertaintyCount`)
- Laws 5, 8, 9, 11, 16, 18, 19

## Doc map
§2, 3A, 7, 8, 9

## Live inventory (2026-08-18)

| Entity | Live API | S02 action |
|--------|----------|------------|
| Assumption | `CreateAssumption` / `UpsertAssumption` / `GetAssumption`; provenance `ACTIVE`\|`STALE`\|`SUPERSEDED`; `MarkStale` → `entity.stale` | **Add** `InvalidateAssumption` (STALE/SUPERSEDED, no delete, optional Discovery + impact finding) |
| Decision | `CreateDecision` / `UpsertDecision` / `GetDecision` | Reuse row; **do not** ALTER `decisions` |
| DecisionAlternative | `InsertDecisionAlternative` / `ListDecisionAlternativesByDecisionID` (mig 009) | Reuse; reconsideration must leave alternatives intact |
| Discovery | `CreateDiscovery` + severity `INFO`\|`PLAN_AFFECTING`\|`BLOCKING` (mig 007) | Reuse for findings — **no Finding table** |
| DecisionImpactFinding | `InsertDecisionImpactFinding`; kind includes `INVALIDATED_ASSUMPTION` | Emit on invalidate when `assumption_supports_decision` links exist |
| `MarkStale` | Generic Law 18 provenance stale on assumption/decision/… | Keep; InvalidateAssumption is the cognitive path (`assumption.invalidated`) |
| Uncertainty/Question | **missing** | New `uncertainties` |
| Hypothesis | **missing** | New `hypotheses` |
| Requirement / Constraint / Risk / Option | n/a | Merge per COVERAGE — **no new tables** |
| `PolicyInputs.BlockingUncertaintyCount` | Caller-populated; SelectNext never EXECUTE when > 0 | S02 **exposes** `CountOpenBlockingUncertaintiesByTaskID`; S06 fills PolicyInputs |
| `ApplyDeliberationTransition` | Persists hop + `deliberation.transition` from complete `PolicyInputs` | S02 does **not** auto-hop; named test **passes the count into** ApplyDeliberationTransition |
| Seed export | No uncertainties/hypotheses keys | **Out of S02** (residual S07); library persistence only |
| Next migration | max **015** | S02-01 adds **`016_cognitive_artifacts.sql`** |

## FINAL locked defaults (S02-01 must not re-debate)

| Item | Value |
|------|-------|
| Migration | **`016_cognitive_artifacts.sql`** — additive; do not rewrite 001–015; **no ALTER** on `assumptions` / `decisions` |
| Compat ceiling | **16** (forbid `017+`); bump `evals/compat`, `production_hardening_test`, `deliberation_test` EmbedExpected, `TestOpenCreatesDBAndMigratesIdempotent` versions |
| New tables | `uncertainties`, `hypotheses`, `decision_reconsiderations` only |
| Forbidden tables | Finding, Requirement, Constraint, Risk, Option, Question-as-separate-from-Uncertainty |
| History | Law 11 — no `Delete*` APIs; reversals are explicit status transitions |
| Raw CoT | Forbidden |
| CLI / MCP / loop | **Library-only** — S06 owns apply write keys + PolicyInputs wiring |
| FTS | Do **not** `SyncEntityFTS` for new types this row (unknown type fail-closed today) |
| Seed JSON | Do **not** extend `SeedDocument` this row |

### `uncertainties` (question text = `title`)

```sql
CREATE TABLE IF NOT EXISTS uncertainties (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    severity TEXT NOT NULL DEFAULT 'INFO'
        CHECK (severity IN ('INFO', 'BLOCKING')),
    status TEXT NOT NULL DEFAULT 'OPEN'
        CHECK (status IN ('OPEN', 'RESOLVED', 'SUPERSEDED')),
    kind TEXT NOT NULL DEFAULT ''
        CHECK (kind IN ('', 'risk', 'gap', 'unknown')),
    confidence REAL NOT NULL DEFAULT 0,
    source_type TEXT NOT NULL DEFAULT '',
    resolution TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    last_verified_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_uncertainties_status_severity
    ON uncertainties(status, severity);
```

| Field | Lock |
|-------|------|
| `title` | Question text (doc §7). Required. |
| `severity` | `INFO` \| `BLOCKING` only. Empty create → `INFO`. `PLAN_AFFECTING` **fail closed** (Discovery-only). |
| `status` | Lifecycle `OPEN` \| `RESOLVED` \| `SUPERSEDED` (UPPER, match residuals). Not provenance ACTIVE/STALE. |
| `kind` | `''` \| `risk` \| `gap` \| `unknown` — Risk merge lives here (`kind=risk`). |
| Links | `entity_links`; **no** task_id/goal_id columns |

**Rels (new):**

| Rel | From → To | Purpose |
|-----|-----------|---------|
| `uncertainty_blocks_task` | uncertainty → task | **SelectNext count** |
| `uncertainty_affects_goal` | uncertainty → goal | Graph/why only — **does not** increment count |

**Create rules:** status on create must be `OPEN` (empty → OPEN). `severity=BLOCKING` **requires** `TaskID` and inserts `uncertainty_blocks_task` (fail closed if missing). `GoalID` optional → `uncertainty_affects_goal`.

**Transitions (from OPEN only):** `ResolveUncertainty(id, resolution)` → `RESOLVED` (resolution required); `SupersedeUncertainty(id, reason)` → `SUPERSEDED` (reason required). RESOLVED/SUPERSEDED are terminal. Events: `entity.created`; `uncertainty.resolved`; `uncertainty.superseded`.

### Blocking count (S01/S06 consume)

```sql
SELECT COUNT(*) FROM uncertainties u
INNER JOIN entity_links l
  ON l.from_type = 'uncertainty' AND l.from_id = u.id
WHERE l.rel = 'uncertainty_blocks_task'
  AND l.to_type = 'task' AND l.to_id = ?
  AND u.severity = 'BLOCKING'
  AND u.status = 'OPEN'
```

- Store: `CountOpenBlockingUncertaintiesByTaskID(taskID string) (int, error)` — empty taskID fail closed.
- Domain: `CountBlockingUncertainties(ctx, taskID) (int, error)`.
- INFO or RESOLVED/SUPERSEDED rows **do not** count.
- S02 **must not** call `SelectNext` / auto-`ApplyDeliberationTransition` on create. S06 passes this int as `PolicyInputs.BlockingUncertaintyCount` into `ApplyDeliberationTransition` every hop.
- Named test: count=1 → `ApplyDeliberationTransition` with that `PolicyInputs` → `INVESTIGATE` / `blocking_uncertainty`.

### `hypotheses` (statement = `title`)

```sql
CREATE TABLE IF NOT EXISTS hypotheses (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'OPEN'
        CHECK (status IN ('OPEN', 'CONFIRMED', 'REJECTED', 'SUPERSEDED')),
    confidence REAL NOT NULL DEFAULT 0,
    source_type TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    last_verified_at TEXT
);
CREATE INDEX IF NOT EXISTS idx_hypotheses_status ON hypotheses(status);
```

| Rel | From → To | Purpose |
|-----|-----------|---------|
| `hypothesis_supported_by` | hypothesis → evidence | Evidence links — **reuse** `evidence`; no second Discovery |
| `hypothesis_addresses_uncertainty` | hypothesis → uncertainty | Optional |

Create: status empty → `OPEN`. Optional `EvidenceIDs` / `UncertaintyID` insert those rels. Transitions from OPEN only: `CONFIRMED` \| `REJECTED` \| `SUPERSEDED` (reason required). Event `entity.created`.

### Assumption invalidate (reuse row)

**Signature (domain):**

```text
InvalidateAssumption(ctx, assumptionID string, in InvalidateAssumptionInput) (Assumption, *Discovery, error)

InvalidateAssumptionInput:
  Status        string   // required STALE | SUPERSEDED (fail closed otherwise, including ACTIVE)
  Reason        string   // required non-empty
  SupersededBy  string   // optional assumption id when SUPERSEDED
  EmitDiscovery bool
  DiscoveryTitle, DiscoveryBody string  // Title required if EmitDiscovery
  TaskIDs       []string // optional RelDiscoveryMentionsTask when discovery emitted
```

| Rule | Lock |
|------|------|
| Row | `UpsertAssumption` status only — **never DELETE** |
| From | `ACTIVE` → `STALE` or `SUPERSEDED`; `STALE` → `SUPERSEDED`; `SUPERSEDED` terminal; **no** revive to ACTIVE (replacement = new `CreateAssumption`) |
| vs `MarkStale` | Keep generic Law 18 `entity.stale`. Invalidate is cognitive: event `assumption.invalidated` |
| Discovery | If `EmitDiscovery`: `CreateDiscovery` severity **`PLAN_AFFECTING`** + rel `discovery_invalidates_assumption` (discovery → assumption) |
| Impact | For each `assumption_supports_decision` link: `InsertDecisionImpactFinding` kind `INVALIDATED_ASSUMPTION`, `related_type=assumption`, `related_id=<id>`, impact_class **`CAUTION`**, uncertainty **`UNKNOWN`** |
| Reconsider | Same linked decisions: auto `RecordDecisionReconsideration` trigger=`invalidated_assumption`, status=`FIRED` |
| Replan | **Do not** auto-replan (Laws 9, 16) |
| New rels | `assumption_supports_decision` (assumption → decision); `assumption_affects_task` (assumption → task) — create-time optional links; invalidate consumes them |

### Decision reconsideration (child table, not JSON on `decisions`)

```sql
CREATE TABLE IF NOT EXISTS decision_reconsiderations (
    id TEXT PRIMARY KEY,
    decision_id TEXT NOT NULL,
    trigger TEXT NOT NULL
        CHECK (trigger IN ('contradicted_effect', 'new_evidence', 'invalidated_assumption')),
    status TEXT NOT NULL DEFAULT 'FIRED'
        CHECK (status IN ('OPEN', 'FIRED')),
    reason TEXT NOT NULL DEFAULT '',
    related_type TEXT NOT NULL DEFAULT '',
    related_id TEXT NOT NULL DEFAULT '',
    reconsider_at TEXT NOT NULL,
    created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_decision_reconsiderations_decision
    ON decision_reconsiderations(decision_id);
```

- `OPEN` = registered watch predicate; `FIRED` = evidence satisfied **now** (default when recording a trigger).
- `reconsider_at` RFC3339; empty input → now.
- Append-only. Does **not** delete Decision or DecisionAlternative. Does **not** auto-`MarkStale` the decision (caller/S06 may).
- Event `decision.reconsider` on `entity_type=decision`.
- Domain: `RecordDecisionReconsideration(ctx, decisionID, ReconsiderationInput)`.

### Files (S02-01)

| Path | Role |
|------|------|
| `internal/store/schema/016_cognitive_artifacts.sql` | Tables + CHECKs + indexes |
| `internal/store/cognitive.go` | Upsert/Get/Count + reconsideration insert/list |
| `internal/store/cognitive_test.go` | Store round-trip + count SQL |
| `internal/domain/cognitive.go` | Create/resolve/invalidate/reconsider APIs + rel constants |
| `internal/domain/cognitive_test.go` | Named domain tests below |
| `internal/domain/service.go` | `EntityUncertainty`, `EntityHypothesis`, rels, event names |
| Compat / embed tests listed above | Ceiling **16** |

Do **not** edit `internal/loop`, `cmd/trace`, MCP, or S01 SelectNext table.

### Named tests (minimum — exact names)

1. `TestCreateUncertaintyDefaultsOpenInfo`
2. `TestBlockingUncertaintyRequiresTaskID`
3. `TestBlockingUncertaintyIncrementsCountForTask`
4. `TestInfoUncertaintyDoesNotIncrementBlockingCount`
5. `TestResolveUncertaintyClearsBlockingCount`
6. `TestSupersedeUncertaintyClearsBlockingCount`
7. `TestCountBlockingUncertaintiesFeedsApplyDeliberationTransition`
8. `TestInvalidateAssumptionSetsStaleAndKeepsRow`
9. `TestInvalidateAssumptionSupersededNoDelete`
10. `TestInvalidateAssumptionEmitsImpactFindingOnLinkedDecision`
11. `TestInvalidateAssumptionOptionalPlanAffectingDiscovery`
12. `TestHypothesisLinksEvidenceWithoutDiscoveryTable`
13. `TestDecisionReconsiderPreservesDecisionAndAlternatives`
14. `TestUnknownUncertaintySeverityFailClosed`

## Merge table

COVERAGE.md entity merge **unchanged**. No gap found.

## Later scopes (upcoming notes only)

- **S06:** call `CountBlockingUncertainties` when building `PolicyInputs`; apply keys `uncertainties` / `hypotheses`; ceiling **16** after this migration.
- **S07:** seed export of new entities is an explicit residual; VERIFY ceiling = live embed max (16 after S02-01).
- **S03:** contradicted effects may link a Hypothesis; do not fork Discovery.

## Planner work

1. [x] Lock `uncertainties` + `hypotheses` schema and migration **016**.
2. [x] Lock Assumption invalidate domain API signature.
3. [x] Lock Decision reconsideration child table.
4. [x] Thicken `01-cognitive-artifacts.md` + `02-scope-review.md`.
5. [x] Confirm COVERAGE merge table unchanged.

## Exit criteria

- [x] 01/02 thickened with named tests
- [x] Merge table unchanged
- [x] Blocking uncertainty query locked for S01/S06
- [x] No product Go

## Next

Orchestrator: **P20-S02-01** after this row is `done`.
