# P29-S00-01 — Peer + surface research

## Metadata
- id: P29-S00-01
- todo_ids: [P29-S00-01]
- role: implementer
- skills: [research, code-explorer]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Produce `RESEARCH.md`: peer GUI patterns, Trace CLI/MCP surface map, recommended browser+HTTP stack, API resource families for S01 OpenAPI, and a draft law carve-out for AGENTS/project-rules. **Investigation only — no product HTTP/GUI code.**

## References

- [00-PLANNER.md](00-PLANNER.md)
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [Phase 29 README](../../README.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) §19
- [docs/TODO.md](../../../../TODO.md) Later developments (hosted product boundary)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Human locks are settled — do not re-grill delivery/transport/stack lean.

## Locked defaults

| Item | Value |
|------|-------|
| Deliverable | `scopes/scope-00-research/RESEARCH.md` |
| Product code | **No** (no `internal/httpapi`, `web/`, `serve` cmd, no Go/TS product edits) |
| Peer floor | ≥3 comparable systems/modules — paths locked below |
| Stack lean | Recommend **TypeScript + Vite + React** unless evidence in RESEARCH overturns |
| Cloud hosting vendor | **Do not** choose |
| Package/path ADR | Recommend only; **S01** locks `internal/httpapi`, OpenAPI path, `web/` |
| Law carve-out | Draft paste-ready text only; **apply in S06** (not this row) |
| Context model | Progressive / bounded — **no** full-graph dump as default HTTP/GUI behavior |

## Locked peer survey targets (minimum)

Survey **all three** (plus skill for #1). Extra peers under `similar projects/` OK if time permits.

| # | Path | Why |
|---|------|-----|
| 1 | `similar projects/Understand-Anything/understand-anything-plugin/packages/dashboard` | Required peer; React+Vite+xyflow graph IA |
| 1b | `similar projects/Understand-Anything/understand-anything-plugin/skills/understand-dashboard` | Agent-facing dashboard skill patterns |
| 2 | `similar projects/codebase-memory-mcp/graph-ui` | React+Vite local graph UI; also skim `src/ui`, `scripts/security-ui.sh`, `scripts/embed-frontend.sh` for daemon/embed security |
| 3 | `similar projects/agentrq/frontend` | Vue+Vite contrast row for stack table |

## Verified Trace inventories (re-check in preflight; update RESEARCH if drift)

### CLI (`cmd/trace/root.go` switch)

`help`, `version`, `init`, `index`, `reindex`, `add`, `link`, `transition`, `review`, `impact`, `capability`, `plan`, `seed`, `tasks`, `why`, `context`, `loop`, `agents`, `migrate`, `backup`, `restore`, `auth`, `install`, `changes`, `patterns`, `knowledge`, `search`, `test`, `tests`, `verify`, `eval`, `outcomes`, `regressions`

**Absent (expected until S03):** `serve`

### MCP (`internal/mcp/mcp_test.go` catalog checklist)

`trace_why`, `trace_context`, `trace_add`, `trace_link`, `trace_transition`, `trace_review`, `trace_tasks`, `trace_capability`, `trace_impact`, `trace_version`, `trace_search`, `trace_changes`, `trace_regressions`, `trace_loop`, `trace_agents`

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -f cmd/trace/root.go
test -d internal/mcp
test -f internal/mcp/mcp_test.go
test -d "similar projects/Understand-Anything/understand-anything-plugin/packages/dashboard"
test -d "similar projects/Understand-Anything/understand-anything-plugin/skills/understand-dashboard"
test -d "similar projects/codebase-memory-mcp/graph-ui"
test -d "similar projects/agentrq/frontend"
test ! -d internal/httpapi
test ! -d web
! grep -q 'case "serve"' cmd/trace/root.go
```

If preflight fails unexpectedly (e.g. `serve` already exists), **stop** and mark `blocked` with Notes — do not invent a second architecture.

## Investigation bounds

**In**

1. Peer UIs: navigation/IA, graph viz, task/loop boards, backend transport (HTTP vs embed vs IPC), security defaults worth copying or avoiding.
2. Trace CLI + MCP inventory → HTTP **resource families** + browser **pages** with P0/P1/defer (parity *intent*, not 1:1 tool mirroring).
3. Stack comparison table (React+Vite lean vs Vue, Svelte, Go templates-only, other) — one recommendation with trade-offs.
4. Law carve-out draft: opt-in `trace serve`, loopback default, Law 19 (adapters → library only), cloud via same OpenAPI / separate hosted product (`docs/TODO.md` Later developments).

**Out**

- Implement server or UI
- Choose cloud hosting vendor / tenancy design
- Point MCP at the internet
- Lock package paths beyond recommendations (S01 owns ADR)
- Rewrite AGENTS.md / project-rules (draft only)

## Role work

1. Run preflight; note any CLI/MCP drift vs inventories above.
2. Survey the three locked peers (read package.json, README/skills, key route/layout files — not a full rewrite).
3. Map Trace surfaces → HTTP families + GUI pages (progressive context).
4. Fill stack options + recommendation.
5. Draft law carve-out bullets.
6. Write `RESEARCH.md` using the template below.
7. Board Notes on **P29-S00-01** only (paths read, peer count, stack verdict).

## RESEARCH.md template

```markdown
# RESEARCH — Phase 29 HTTP API + browser GUI

**Date:** YYYY-MM-DD  
**Row:** P29-S00-01

## Executive summary
(3–5 sentences)

## Peer matrix
| Peer | Nav / IA | Graph | Tasks / loop | Backend transport | Takeaways for Trace |
(≥3 rows; include Understand-Anything dashboard, codebase-memory-mcp graph-ui, agentrq frontend)

## Trace surface map
| Surface | CLI | MCP | HTTP candidate family | GUI page | Priority (P0/P1/defer) |

## API resource families (for S01 OpenAPI)
(List: health, project/meta, tasks, loop, entities, seed, graph-bounded, …)

## Stack options
| Option | Pros | Cons | Verdict |
| Recommended stack | … | rationale |
(Overturn TS+Vite+React only with evidence)

## Law carve-out draft
(Paste-ready bullets for AGENTS.md / project-rules — applied in S06)

## Risks / open decisions for S01–S02
```

## Do not

- Implement server or UI
- Choose cloud hosting vendor
- Point MCP at the internet
- Lock package paths beyond recommendations (S01 owns ADR)
- Start P29-S00-02 or S01 work

## Exit criteria

- [ ] `RESEARCH.md` exists with all template sections
- [ ] Peer matrix ≥3 rows covering locked peer paths
- [ ] Trace surface map covers CLI + MCP parity intent (inventories above)
- [ ] Explicit API resource family list for S01
- [ ] Stack recommendation with trade-offs (lean or evidence overturn)
- [ ] Law carve-out draft present
- [ ] No product code changed
- [ ] Board Notes on P29-S00-01 cite evidence

## Minimal todos

- [ ] Preflight PASS
- [ ] Peer survey (3 locked paths) + Trace CLI/MCP inventory
- [ ] Write RESEARCH.md
- [ ] Board status + Notes on P29-S00-01

## Todo updates

Status + notes on **P29-S00-01** only.

## Next

**P29-S00-02**
