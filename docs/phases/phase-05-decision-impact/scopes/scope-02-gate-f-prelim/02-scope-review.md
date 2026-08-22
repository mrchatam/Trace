# P05 / S02 / 02 — Scope review (Gate F prelim)

## Metadata
- id: P05-S02-02
- todo_ids: [P05-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of **P05-S02-01** against S02-00 locks and `01-gate-f-prelim.md`. APPROVE with evidence or spawn remediations. Fresh session — do not share the implementer’s context.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-gate-f-prelim.md](01-gate-f-prelim.md)
- S01 APPROVE: [../scope-01-impact-classes/REVIEW-NOTES.md](../scope-01-impact-classes/REVIEW-NOTES.md)
- Live: `evals/impact/` (expect after S02-01); `internal/domain` ImpactReport surface
- Prior review pattern: [../../phase-04-review-depth/scopes/scope-02-honesty-escape-rate/REVIEW-NOTES.md](../../../phase-04-review-depth/scopes/scope-02-honesty-escape-rate/REVIEW-NOTES.md) (Gate G)

## Session start
Agent → clarify if needed → Plan → execute (review).

## Review focus (checklist)

### Gate F harness (blocker if missing)
- [ ] Package **`evals/impact`** exists (not folded into honesty/replan/x0/p0x)
- [ ] Named test **`TestPlantedImpactConflictsGateFPrelim`** present and PASS
- [ ] Committed **`evals/impact/schema-gate-f.json`** with `schema_version` const **1**, `gate` const **F**, `suite` const **impact**
- [ ] Test writes temp **`metrics-gate-f.json`** and validates against schema
- [ ] Planted tallies match locks: TP=3, FN=0, FP=0, TN=1; precision=1.0; recall=1.0; probes=4
- [ ] Probes cover Pos-1 UNKNOWN, Pos-2 DESTRUCTIVE rollup, Pos-3 empty-findings+link, Neg-1 clean SAFE
- [ ] Scoring does **not** trust `OverallClass` alone when `HasUnknown` is required (Pos-1)
- [ ] Evidence is the harness — **not** Notes-only / vibes Gate F

### S01 consume / G19
- [ ] Plants via `AddImpactFinding` + `LinkDecisionTask` + `ImpactReport` (domain library)
- [ ] No new entity_links rels; mig **009** only (no S02 schema fork)
- [ ] Eval package does **not** import `cmd/trace`; no product API invention
- [ ] No `internal/impact` package; no `plan simulate`

### Carry-forward bars (must stay green)
- [ ] Honesty A/B/C `TestHonestyFailClosedPlantedClaim`
- [ ] Gate G `TestHonestyEscapeRateGateGPrelim` (escapes=1/caught=2/attempts=3 untouched)
- [ ] Gate E `TestPlantedDiscoveryReplan`
- [ ] p0x 7/7; x0; `CGO_ENABLED=0` impact+domain+store; `CGO_ENABLED=1` `./...`
- [ ] Gate C `docs/verification/gate-c-x0/` **not rewritten**; `dry_run:false` intact
- [ ] No commercial multi-model Gate F claim; dry-run ≠ Gate C ≠ Gate F

### Re-verify commands

```bash
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./... -count=1
```

## Role work
1. Diff claims (01 + board Notes) vs `evals/impact` tree + S01 APIs
2. Severity-tag findings; blocker/high → small inline fix **or** spawn `02a`/`02b` (full prompts)
3. Write [REVIEW-NOTES.md](REVIEW-NOTES.md); mark board + SCOPE-TODOS
4. Light-thicken upcoming S03 Depends if re-prove command needs a stamp (upcoming only)

## Exit criteria
- [ ] Verdict (APPROVE / spawn) + REVIEW-NOTES.md with evidence table
- [ ] No open blocker/high without pending follow-up
- [ ] Confidence medium or high with residuals listed
- [ ] Board status + Notes; next runnable **P05-S03-00** on APPROVE

## Out of scope
- Implementing Gate F (S02-01); phase VERIFY product; rewriting done S01 history; commercial Gate F
