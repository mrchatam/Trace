# P00 / S05 / 02 — Scope review (Work/causal API)

## Metadata
- id: P00-S05-02
- todo_ids: [P00-S05-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S05 (Work/causal API). Severity-tag findings; small fixes or spawn `a` implement + `b` review pairs with **full** prompts. Forward-only. May thicken **upcoming** prompts if this scope’s surface drifted.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop — write proper spawn prompts)
- Sibling `01-*.md` + board Notes

## Session start
Agent → clarify → Plan → review (fresh subagent).

## Review focus
- Exit criteria honestly met vs thickened `01-causal.md` locks?
- Provenance on creates (G5); Goal→Task / Decision→Task / Discovery→PlanChange tests
- `work_state` separate from provenance `status`; transition graph + DONE policy stub
- Events append-only for create/link/transition; no second DB / no CLI / no Claim promotion engine
- `CGO_ENABLED=0` for domain+store; cross-scope notes for S06/S07 still accurate

## Exit criteria
- [ ] Findings recorded
- [ ] blocker/high fixed or spawned
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated

## Minimal todos
- [ ] Compare claims vs evidence
- [ ] Fix or spawn
- [ ] Re-verify
- [ ] Board update
