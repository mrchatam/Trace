# P05 / S01 / 01 — Impact classes

## Metadata
- id: P05-S01-01
- todo_ids: [P05-S01-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: [Shell, Read, Write, Grep, Glob]
- agents: []
- verification: automated

## Objective
Implement **manual/planted decision impact classes + alternatives + thin impact report** against the live Decision surface (`CreateDecision` / `GetDecision` / `LinkDecisionTask` only today). Extend `internal/domain` + store mig **`009_decision_impact`** — no second impact stack, no planner fork, no commercial auto-engine (DR-NOIMP). Expose APIs S02 Gate F can plant/score. Keep honesty / Gate G / Gate E / p0x / x0 / Gate C bars green. No daemon/HTTP/embeddings. No `plan simulate`.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) locks (this scope)
- [phase README](../../README.md)
- [docs/ROADMAP.md](../../../../ROADMAP.md) P12 — impact / knowledge bands
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 5
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) DR-NOIMP
- Live: `internal/domain` CreateDecision/GetDecision/LinkDecisionTask (`decision_affects_task`); `internal/store` decisions + entity_links; CLI `add decision` / `link decision-task`; migs `001`…`008`

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | Keep `go.mod` floor (currently 1.24.0); do not downgrade |
| Package | **`internal/domain`** on `*store.Store` (+ store helpers). **Do not** invent `internal/impact` / second stack or put impact under `internal/planner` |
| Reuse | Keep `CreateDecision` / `GetDecision` / `LinkDecisionTask` (`decision_affects_task`) unchanged in behavior |
| New entity_links rels | **None this scope** — affected work stays `decision_affects_task`; findings may optionally name `related_type`/`related_id` without new rels |
| Migration | Additive embed **`009_decision_impact.sql`** only (do not rewrite `001`–`008`) |
| Impact bands | **`SAFE` \| `CAUTION` \| `HIGH` \| `DESTRUCTIVE` \| `REVERSAL`** (UPPER; reject unknown fail-closed; empty rejected) |
| Uncertainty | **`KNOWN` \| `LIKELY` \| `POSSIBLE` \| `UNKNOWN`** (UPPER; reject unknown; empty on create → **`UNKNOWN`**) |
| Finding kinds | **`AFFECTED_WORK` \| `INVALIDATED_ASSUMPTION` \| `WORK_AT_RISK` \| `NEW_WORK` \| `UNRESOLVED`** (UPPER; reject unknown/empty) |
| Alternatives | Thin rows on a decision; **at most one** `is_recommended=1` (domain enforces) |
| Impact report | Library `ImpactReport(decisionID)` (+ thin CLI); aggregates decision + `decision_affects_task` links + findings + alternatives; **fail-closed** on unknown (see Policy) |
| Why / Expand | Prefer consume existing retrieval/Why over links already present; **do not** invent embeddings or auto-expand cascades |
| CLI | Thin G19: new top-level **`trace impact`** (findings / alternatives / report). Keep `add decision` + `link decision-task` |
| MCP | **Not** required. Do not add MCP impact tools this scope |
| CGO | Domain + store APIs must pass `CGO_ENABLED=0` |
| Carry-forward bars | Honesty A/B/C; Gate G (`TestHonestyEscapeRateGateGPrelim`); Gate E (`TestPlantedDiscoveryReplan`); p0x 7/7; x0; Gate C `dry_run:false` artifacts intact |
| Out | Gate F harness (S02); phase VERIFY (S03); commercial auto-engine; embeddings; daemon/HTTP; VerifiedFact; `plan simulate`; new entity_links rels; planner fork |

### Schema (locked)

```sql
-- 009_decision_impact.sql (additive)

CREATE TABLE IF NOT EXISTS decision_impact_findings (
    id TEXT PRIMARY KEY,
    decision_id TEXT NOT NULL,
    impact_class TEXT NOT NULL,
    uncertainty TEXT NOT NULL DEFAULT 'UNKNOWN',
    kind TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    related_type TEXT NOT NULL DEFAULT '',
    related_id TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_decision_impact_findings_decision
    ON decision_impact_findings(decision_id);

CREATE TABLE IF NOT EXISTS decision_alternatives (
    id TEXT PRIMARY KEY,
    decision_id TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL DEFAULT '',
    is_recommended INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_decision_alternatives_decision
    ON decision_alternatives(decision_id);
```

No ALTER on `decisions` — impact lives in the new tables (manual/planted; DR-NOIMP).

### Vocabulary (locked)

```text
# impact_class (reject unknown / empty)
ImpactClassSAFE         = "SAFE"
ImpactClassCaution      = "CAUTION"
ImpactClassHigh         = "HIGH"
ImpactClassDestructive  = "DESTRUCTIVE"
ImpactClassReversal     = "REVERSAL"

# uncertainty (empty → UNKNOWN; reject other unknowns)
UncertaintyKNOWN    = "KNOWN"
UncertaintyLIKELY   = "LIKELY"
UncertaintyPOSSIBLE = "POSSIBLE"
UncertaintyUNKNOWN  = "UNKNOWN"

# finding kind (reject unknown / empty)
FindingKindAffectedWork            = "AFFECTED_WORK"
FindingKindInvalidatedAssumption   = "INVALIDATED_ASSUMPTION"
FindingKindWorkAtRisk              = "WORK_AT_RISK"
FindingKindNewWork                 = "NEW_WORK"
FindingKindUnresolved              = "UNRESOLVED"
```

### Minimum public API (behavior locked; names may vary slightly)

```text
# Domain — findings
AddImpactFinding(ctx, decisionID, ImpactFindingInput) (store.DecisionImpactFinding, error)
  // decision must exist (GetDecision); validate class/kind; uncertainty default UNKNOWN
  // related_type/related_id optional strings (no FK enforcement this scope)
ListImpactFindings(ctx, decisionID) ([]store.DecisionImpactFinding, error)

# Domain — alternatives
AddDecisionAlternative(ctx, decisionID, AlternativeInput) (store.DecisionAlternative, error)
  // If input.Recommended true → clear other recommended for this decision, set this one
SetRecommendedAlternative(ctx, decisionID, alternativeID) error
  // alternative must belong to decision; clear siblings; set is_recommended=1
ListDecisionAlternatives(ctx, decisionID) ([]store.DecisionAlternative, error)

# Domain — report (S02 Gate F plant/score hook)
ImpactReport(ctx, decisionID) (ImpactReportResult, error)
  // Loads decision; lists decision_affects_task links (ListLinksFrom); findings; alternatives
  // Never mutates plan / tasks / decisions.status
  // Fail-closed fields (locked):
  //   HasUnknown     = any finding.uncertainty == UNKNOWN OR findings empty while ≥1 decision_affects_task link
  //   Incomplete     = HasUnknown OR (alternatives non-empty AND no recommended)
  //   OverallClass   = max severity among findings using order:
  //                    SAFE < CAUTION < HIGH < DESTRUCTIVE < REVERSAL
  //                    (empty findings → OverallClass "" and HasUnknown true when links exist;
  //                     empty findings + no links → OverallClass "" + Incomplete true)
  //   OverallUncertainty = worst among findings (UNKNOWN > POSSIBLE > LIKELY > KNOWN);
  //                        empty findings → UNKNOWN
  // Must include UNKNOWN findings in Findings slice — never silently drop

# Store helpers (as needed)
InsertDecisionImpactFinding / ListDecisionImpactFindingsByDecisionID
InsertDecisionAlternative / UpdateDecisionAlternativeRecommended /
ClearRecommendedAlternatives / ListDecisionAlternativesByDecisionID
```

### Policy (locked)

```text
CreateDecision / LinkDecisionTask: UNCHANGED
Impact APIs are recording + report only — do NOT auto-supersede plans, mutate tasks, or invent links
Fail-closed: reject bad enums; report surfaces Incomplete/HasUnknown; never claim SAFE overall when UNKNOWN present
DR-NOIMP: no auto-classifier / embeddings / “engine” UX theater
plan simulate: OUT
```

### Target tree

```text
internal/store/
  schema/009_decision_impact.sql
  impact.go                 # DecisionImpactFinding + DecisionAlternative + CRUD/list helpers

internal/domain/
  service.go                # impact/uncertainty/kind consts OK here or impact.go
  impact.go                 # AddImpactFinding / List* / AddDecisionAlternative /
                            # SetRecommendedAlternative / ImpactReport
  create.go / link.go       # untouched behavior for Decision CRUD + LinkDecisionTask
  domain_test.go            # policy tests (see Exit criteria)

cmd/trace/
  impact.go                 # thin G19 subcommands
  help.go                   # usage lines
  # add.go / link.go        # keep decision + decision-task; no business impact logic
```

### Tests (required)

- Mig `009_decision_impact.sql` applied on Open (both tables present)
- `AddImpactFinding` + `ListImpactFindings`; defaults uncertainty→UNKNOWN; reject bad class/kind/uncertainty
- `AddDecisionAlternative` + recommend exclusivity (`SetRecommendedAlternative` / Recommended-on-create)
- `ImpactReport` includes `decision_affects_task` linked task IDs + findings + alternatives
- Fail-closed: findings empty + ≥1 task link → `HasUnknown`/`Incomplete`; UNKNOWN finding never omitted; recommended missing when alternatives exist → `Incomplete`
- OverallClass rollup order SAFE < CAUTION < HIGH < DESTRUCTIVE < REVERSAL
- Decision CRUD + `LinkDecisionTask` regression still green
- Do **not** require Gate F package this scope (`evals/impact` is S02)

### CLI (thin G19)

```text
trace impact finding add --decision <id> --class SAFE|CAUTION|HIGH|DESTRUCTIVE|REVERSAL \
  --kind AFFECTED_WORK|INVALIDATED_ASSUMPTION|WORK_AT_RISK|NEW_WORK|UNRESOLVED \
  [--uncertainty KNOWN|LIKELY|POSSIBLE|UNKNOWN] [--body <text>] \
  [--related-type <t>] [--related-id <id>] [--id <uuid>]
trace impact finding list --decision <id>
trace impact alternative add --decision <id> --title <t> [--body <text>] [--recommended] [--id <uuid>]
trace impact alternative list --decision <id>
trace impact alternative recommend --decision <id> --id <alternative_id>
trace impact report --decision <id>
```

Stdout machine-friendly JSON lines (match existing `add`/`review` style). Exit 0/1/2 per CLI norms.

### S02 Gate F hooks (must expose — do not implement harness)

S02 planner/implement will plant via domain APIs (not CLI-only):

- `AddImpactFinding` with known vs conflicting classes/uncertainty
- `LinkDecisionTask` + `ImpactReport` for precision/recall-style tallies
- `HasUnknown` / `Incomplete` / `OverallClass` fields as scoreable signals

## Board rights
Implementer: **status + notes only**. No spawning; no rewriting upcoming prompts.

## Exit criteria
- [ ] Mig `009_decision_impact.sql` + store finding/alternative helpers live
- [ ] Domain APIs: Add/List findings; Add/List/SetRecommended alternatives; `ImpactReport` fail-closed as locked
- [ ] No new entity_links rels; Decision CRUD + `decision_affects_task` unchanged
- [ ] Thin CLI `trace impact` wired (G19); no MCP impact tools; no planner fork / `internal/impact`
- [ ] Domain tests cover cases above
- [ ] `CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/...` green
- [ ] Carry-forward: `CGO_ENABLED=0 go test ./evals/honesty/... ./evals/replan/...` + `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./...` green
- [ ] Gate C artifacts under `docs/verification/gate-c-x0/` untouched
- [ ] TODO.md status + Notes updated (this row only)

## Minimal todos
- [ ] Store: mig 009 + finding/alternative CRUD/list helpers
- [ ] Domain: impact APIs + ImpactReport + consts; Decision link regression tests
- [ ] Thin CLI `trace impact` + help
- [ ] Full carry-forward bars; board Notes

## Out of scope
- `evals/impact` / Gate F named test (S02)
- Phase VERIFY / Phase 06 handoff (S03)
- Commercial impact engine; embeddings; daemon/HTTP; VerifiedFact; `plan simulate`
- Auto-mutating plans on report; new entity_links relationship types
