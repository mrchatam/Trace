# P00 / S01 / 02 — Scope review (Go scaffold)

## Metadata
- id: P00-S01-02
- todo_ids: [P00-S01-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S01 deliverables. Find gaps/bugs; small fixes inline or **spawn** `03a` implement + `03b` review (full prompts). Forward-only. Thicken upcoming S02+ prompts if path assumptions wrong.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop)
- `01-go-module-scaffold.md` + TODO Notes for P00-S01-01

## Session start
Agent → clarify → Plan → review (fresh agent; no implementer memory).

## Review focus
- Module path / layout vs locked S01 taxonomy (`store`, `vcs`, `gitcli`, `analyzers`, `domain`, `retrieval`, `compiler`)
- Forbidden names under `internal/`: `context`, `contextx`, `sqlite`
- Accidental daemon/HTTP/MCP code
- Missing `.trace/` / `bin/` gitignore
- Misleading stubs claimed complete
- Downstream path drift in S02+ prompts

## Board rights
Reviewer may spawn rows + edit **upcoming** prompts. Must not mutate `done` prompt bodies for S01-00/01.

## Exit criteria
- [x] Findings recorded (REVIEW-NOTES.md or TODO Notes)
- [x] Every blocker/high fixed or has pending spawn pair
- [x] Confidence medium or high
- [x] TODO.md updated

## Minimal todos
- [x] Diff claims vs tree
- [x] Findings → fix or spawn
- [x] Re-verify build/test
- [x] Board update
