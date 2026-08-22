# P00 / S07 / 02 — Scope review (trace CLI)

## Metadata
- id: P00-S07-02
- todo_ids: [P00-S07-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S07 (trace CLI). Severity-tag findings; small fixes or spawn `a` implement + `b` review pairs with **full** prompts. Forward-only. May thicken **upcoming** prompts if this scope’s surface drifted.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop — write proper spawn prompts)
- Sibling `01-*.md` + board Notes

## Session start
Agent → clarify → Plan → review (fresh subagent).

## Review focus
- Exit criteria honestly met vs thickened `01-cli.md`?
- **G19 / DR-API:** no business logic in CLI (ranking, FTS, provenance machines, parsers)
- Command surface: init / index|reindex / add / link / transition / seed import / why / context + help/version
- Causal path uses `domain.*` only; index uses `analyzers.IndexFile*` + SkipError; why/context call retrieval/compiler
- Seed JSON v1 + TransitionTask for DONE (not CreateTask work_state shortcuts)
- Incremental: index A/B isolation test; no full-rebuild-as-default
- No cobra/MCP/daemon/HTTP/dump; CGO build of `cmd/trace` documented/passing
- Silent failures, missing tests, S08 readiness (seed + CLI walkthrough)
- Cross-scope: S08 can consume locked seed format + commands

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
