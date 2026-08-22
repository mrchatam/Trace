# P00 / S04 / 02 — Scope review (tree-sitter analyzers)

## Metadata
- id: P00-S04-02
- todo_ids: [P00-S04-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S04 (tree-sitter analyzers). Severity-tag findings; small fixes or spawn `a` implement + `b` review pairs with **full** prompts. Forward-only. May thicken **upcoming** prompts if this scope’s surface drifted.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop — write proper spawn prompts)
- Sibling `01-*.md` + board Notes

## Session start
Agent → clarify → Plan → review (fresh subagent).

## Review focus
- Exit criteria honestly met? (golden JS/TS + Python; incremental A/B isolation)
- DR-PARSE / DR-ANLANG: official tree-sitter binding + only TS/JS/Python
- DR-INCREMENTAL: file-local replace only; no full-rebuild default path
- Persistence only via store UpsertFile + ReplaceFileSymbols + ReplaceFileImports (no parallel schema)
- CGO confined to analyzers; store/vcs/gitcli still `CGO_ENABLED=0`-clean
- No `gitcli` import from analyzers; no MCP/daemon/HTTP/CLI commands
- Silent failures, missing tests, Close() on tree-sitter C objects
- Cross-scope: S06/S07/S08 can consume IndexFile + ListSymbols/Imports

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
