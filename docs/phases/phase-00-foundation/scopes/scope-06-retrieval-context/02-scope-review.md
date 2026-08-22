# P00 / S06 / 02 — Scope review (Retrieval + context compiler)

## Metadata
- id: P00-S06-02
- todo_ids: [P00-S06-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S06 (Retrieval + context compiler). Severity-tag findings; small fixes or spawn `a` implement + `b` review pairs with **full** prompts. Forward-only. May thicken **upcoming** prompts if this scope’s surface drifted.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop — write proper spawn prompts)
- Sibling `01-*.md` + board Notes

## Session start
Agent → clarify → Plan → review (fresh subagent).

## Review focus
- Exit criteria honestly met?
- P0-X / incremental / provenance laws
- Silent failures, missing tests
- Cross-scope breakage for later scopes

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
