# P03 / S02 / 02 — Scope review (discovery replan)

## Metadata
- id: P03-S02-02
- todo_ids: [P03-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of P03-S02-01. APPROVE with evidence or spawn forward remediations.

## Session start
Agent → clarify if needed → Plan → execute (review).

## Review focus
- **Severity policy real:** only `PLAN_AFFECTING`/`BLOCKING` auto-open replan; `INFO` does not supersede or increment `auto_replan_count`
- **Churn real:** N=5 fail-closed; `IncrementAutoReplanCount` + `AckAutoReplan` (reset to 0); no unbounded loops
- **S01 consume:** uses `SupersedeDeepPlan` / GetCurrent|ListScopes|GetPlan; keeps `LinkDiscoveryPlanChange`; no second planner package
- **Demo:** `evals/replan` `TestPlantedDiscoveryReplan` matches claims
- **CLI G19:** `--severity`, `plan apply-discovery`, `plan ack-replan` thin only
- **Out of scope held:** DPC attach still global (not scoped); honesty / p0x / Gate C artifacts intact; Mode-B packs not rewritten
- Light-check S03 VERIFY stubs still list Gate E inputs (`evals/replan` + churn/severity proof)

## Exit criteria
- [ ] Verdict + confidence + REVIEW-NOTES.md
- [ ] No open blocker/high without spawn
- [ ] Board status + Notes

## Minimal todos
- [ ] Diff claims vs evidence
- [ ] Re-run required tests (planner/store/domain/replan + honesty + p0x + x0 + `./...`)
- [ ] Write REVIEW-NOTES; mark done or spawn
