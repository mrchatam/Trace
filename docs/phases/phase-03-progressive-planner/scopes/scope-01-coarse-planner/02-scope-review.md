# P03 / S01 / 02 — Scope review (coarse planner)

## Metadata
- id: P03-S01-02
- todo_ids: [P03-S01-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of P03-S01-01. APPROVE with evidence or spawn forward remediations. No silent weaken of honesty/p0x bars.

## Session start
Agent → clarify if needed → Plan → execute (review).

## Review focus
- Claims match repo: **`internal/planner`** (not domain Phase/Scope dump) + mig **`006_plan_hierarchy.sql`**
- Hierarchy Goal→phase→scope; `scope_deep_plans` supersede-not-delete; `goal_plan_state` current pointer
- Progressive: DeepPlan fail-closed unless current; one lookahead shallow only; **no** whole-backlog auto-generation
- S02 hooks present: `SupersedeDeepPlan`, List/GetCurrent scope, `auto_replan_count` column
- Thin `trace plan` CLI (G19); no MCP plan tools required
- Laws: no daemon/HTTP/embeddings primary
- Honesty / p0x / x0 still green
- Light-check upcoming S02 stubs still compatible with shipped S01 surface

## Exit criteria
- [ ] Verdict + confidence + REVIEW-NOTES.md
- [ ] No open blocker/high without spawn
- [ ] Board status + Notes

## Minimal todos
- [ ] Diff claims vs evidence
- [ ] Re-run required tests
- [ ] Write REVIEW-NOTES; mark done or spawn
