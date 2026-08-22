# P05 / S01 / 00-PLANNER — Impact classes

## Metadata
- id: P05-S01-00
- todo_ids: [P05-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-impact-classes.md` for **decision impact classes + alternatives** against live Decision / entity_links surface. Lock package paths, APIs, persistence, CLI, and exit criteria. No product code in this planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 5
- [docs/ROADMAP.md](../../../../ROADMAP.md) P12 — impact / knowledge bands
- [docs/init/D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) DR-NOIMP / DR-HANDOFF
- Live: `internal/domain` CreateDecision/GetDecision/LinkDecisionTask; `internal/store` decisions + entity_links; CLI `add decision` / `link decision-task`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (gaps — locked 2026-08-16)
| Item | Today (post–Phase 04 / P05-00) | S01 lock |
|------|--------------------------------|----------|
| Decision CRUD | `CreateDecision` / `GetDecision` / store Upsert; provenance ACTIVE\|STALE\|SUPERSEDED | Keep; do not fork entity |
| Links | `decision_affects_task` only (`LinkDecisionTask`) | **No new rels** — findings may optional `related_type`/`related_id` |
| Impact classes | Absent | Findings table; bands `SAFE`\|`CAUTION`\|`HIGH`\|`DESTRUCTIVE`\|`REVERSAL` |
| Uncertainty | Absent on findings | `KNOWN`\|`LIKELY`\|`POSSIBLE`\|`UNKNOWN` (empty→UNKNOWN) |
| Alternatives | Absent | `decision_alternatives` + single recommended |
| Impact report | Absent | `ImpactReport` + CLI `trace impact report` — fail-closed `HasUnknown`/`Incomplete` |
| Why/retrieval | `decision_affects_task` reason only | Prefer consume existing Expand/Why; no embeddings |
| Gate F | No `evals/impact` yet | Expose `AddImpactFinding` + `ImpactReport` hooks for S02 |
| DR-NOIMP | PROVISIONAL full engine deferred | **Manual/planted classes + report** only |
| Migration | `001`…`008` | Additive **`009_decision_impact.sql`** |

## Phase defaults already locked (respect — P05-00)
| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Honesty / Gate G / Gate E / p0x / x0 / Gate C | Keep green / intact |
| Package hint | Prefer **`internal/domain`** + store — no second impact stack; no planner fork |
| Migration hint | Next additive **`009_*`** (name locked here) |
| Daemon/HTTP/embeddings | Forbidden as primary |
| MCP | Optional; CLI primary |
| `plan simulate` | Out this phase |
| VerifiedFact | Out unless explicitly promoted |

## Planner work
- [x] Inventory live Decision / entity_links / retrieval vs impact-class needs (confirm table above)
- [x] Lock package / schema / CLI / impact + alternatives model
- [x] Thicken `01-impact-classes.md` enough to run alone
- [x] Light-update upcoming S02 stubs with expected Gate F hooks from S01 surface
- [x] Sync SCOPE-TODOS + board Notes; mark this row done

## Locked defaults (set 2026-08-16 — do not re-debate in S01-01)

| Item | Value |
|------|-------|
| Migration | **`009_decision_impact.sql`** — tables `decision_impact_findings` + `decision_alternatives` |
| Package | **`internal/domain`** + store helpers — no `internal/impact`; no planner fork |
| entity_links | **Keep `decision_affects_task` only** — no new rels this scope; findings may optional `related_type`/`related_id` |
| Impact bands | `SAFE` \| `CAUTION` \| `HIGH` \| `DESTRUCTIVE` \| `REVERSAL` |
| Uncertainty | `KNOWN` \| `LIKELY` \| `POSSIBLE` \| `UNKNOWN` (empty→UNKNOWN) |
| Finding kinds | `AFFECTED_WORK` \| `INVALIDATED_ASSUMPTION` \| `WORK_AT_RISK` \| `NEW_WORK` \| `UNRESOLVED` |
| APIs | `AddImpactFinding` / `ListImpactFindings`; `AddDecisionAlternative` / `SetRecommendedAlternative` / `ListDecisionAlternatives`; `ImpactReport` (fail-closed `HasUnknown`/`Incomplete`) |
| CLI | Thin **`trace impact`** (finding / alternative / report); keep `add decision` + `link decision-task` |
| MCP | Out this scope |
| Gate F hooks | S02 plants via domain APIs + `ImpactReport` fields — harness stays S02 |

## Exit criteria
- [x] `01-impact-classes.md` runnable alone
- [x] Paths + impact model locked
- [x] Light S02 Depends note updated
- [x] No product Go in this row

## Minimal todos
- [x] Inventory Decision APIs vs S01 needs
- [x] Thicken 01 + 02 + light S02 Depends
- [x] Mark P05-S01-00 done

## Out of scope
- Product Go; Gate F harness (S02); phase VERIFY (S03); embeddings; daemon/HTTP; `plan simulate`; commercial impact engine
