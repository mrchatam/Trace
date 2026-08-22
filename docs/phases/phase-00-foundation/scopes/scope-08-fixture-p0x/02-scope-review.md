# P00 / S08 / 02 — Scope review (Fixture + P0-X harness)

## Metadata
- id: P00-S08-02
- todo_ids: [P00-S08-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S08 (Fixture + P0-X harness). Severity-tag findings; small fixes or spawn `a` implement + `b` review pairs with **full** prompts. Forward-only. May thicken **upcoming** prompts if this scope’s surface drifted.

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
- **Seed path:** harness uses **absolute** (or cwd-valid) path for `seed import` — not rewritten under `-C`
- **Both langs:** fixture has TS/JS **and** Python; index asserts symbols/imports non-vacuous
- **7/7 map:** each criterion has a real assertion; ≥5 named understanding queries; incremental sibling-file isolation
- **No** committed `.trace/`; no MCP/daemon; no `evals/x0` agent work

## Exit criteria
- [x] Findings recorded
- [x] blocker/high fixed or spawned
- [x] Confidence medium or high (residuals listed if medium)
- [x] TODO.md updated

## Minimal todos
- [x] Compare claims vs evidence
- [x] Fix or spawn
- [x] Re-verify
- [x] Board update
