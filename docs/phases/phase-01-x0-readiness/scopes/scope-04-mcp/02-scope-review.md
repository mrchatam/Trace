# P01 / S04 / 02 — Scope review (MCP adapter)

## Metadata
- id: P01-S04-02
- todo_ids: [P01-S04-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S04 (MCP adapter). Findings by severity; small fixes or spawn. Forward-only. May thicken **upcoming** S05 VERIFY MCP tool checklist with the live tool list from implementer Notes.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop)
- Sibling [`01-mcp.md`](01-mcp.md) + board Notes
- DR-AGENT / DR-SURFACE / DR-API / G19
- Live: `internal/mcp`, `cmd/trace-mcp`, `cmd/trace` (parity contract)

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → review). Fresh subagent — do not share implementer session.

## Locked review assumptions (from S04 planner)

| Item | Expected |
|------|----------|
| Layout | `internal/mcp` + thin `cmd/trace-mcp` |
| SDK | Official `github.com/modelcontextprotocol/go-sdk` (stdio) |
| Required tools | `trace_why`, `trace_context` → retrieval/compiler |
| Parity tools | `trace_add` / `trace_link` / `trace_transition` / `trace_review` (or documented deferral) |
| G19 | No library import of MCP/cmd; no business-logic fork in MCP |
| Out | Daemon/HTTP-as-primary; embeddings; raw SQL tools; MCP required for X0 |
| Regression | `evals/x0` + `evals/p0x` + `evals/honesty` still PASS without MCP |

## Review focus
- Thin adapter only — handlers open store and call `domain` / `retrieval` / `compiler`; no duplicated DONE/why/context logic?
- Tool names/args match [`01-mcp.md`](01-mcp.md) schemas (or Notes document intentional deltas)?
- Transport is stdio; no HTTP/SSE daemon shipped as product primary?
- Import edges: `rg 'internal/mcp|cmd/trace-mcp' internal/` (and `cmd/trace`) should be empty of reverse deps
- `go build ./cmd/trace-mcp` + `CGO_ENABLED=1 go test ./evals/x0/... ./evals/p0x/... ./evals/honesty/... ./...` PASS
- Tool list present in implementer Notes for S05 VERIFY
- If VERIFY checklist needs the concrete tool list, thicken **upcoming** S05 prompts only

## Exit criteria
- [ ] Findings recorded (prefer `REVIEW-NOTES.md` in this folder)
- [ ] blocker/high fixed inline or spawned `a`/`b` pairs immediately below this row
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] S05 upcoming notes updated if tool checklist was underspecified
- [ ] TODO.md status + Notes updated

## Minimal todos
- [ ] Compare claims (01 + Notes) vs repo evidence
- [ ] Grep import boundaries + run build/tests
- [ ] Fix or spawn
- [ ] Re-verify
- [ ] Board update (+ light S05 if needed)
