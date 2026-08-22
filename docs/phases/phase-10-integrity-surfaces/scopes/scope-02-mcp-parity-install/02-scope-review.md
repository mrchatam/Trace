# P10 / S02 / 02 — Scope review (MCP parity / install)

## Metadata
- id: P10-S02-02
- todo_ids: [P10-S02-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S02 (**DF-21 / DF-22 / DF-32**). Fresh subagent. Compare claims + locks to live code/tests. Small inline fix **or** spawn `02a`/`02b` for blocker/high. Do not rewrite S01/`done` history.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) FINAL locks + [01-mcp-parity-install.md](01-mcp-parity-install.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- Live: `internal/mcp/`, `cmd/trace/{tasks,capability,install}.go`, `BuiltinMCPCapabilitySpecs`

## Session start
Agent → clarify → Plan → execute (reviewer).

## Checklist (must all pass for APPROVE)

| # | Check | Evidence |
|---|--------|----------|
| 1 | **DF-21** — `trace_tasks` registered; output matches CLI row keys; optional `goal_id` | `server.go` + mcp test |
| 2 | **DF-21** — `trace_capability` supports declare\|list\|require\|unrequire\|missing via domain only (G19) | tools + import-boundary test |
| 3 | **DF-21** — **no** plan/impact/index MCP tools | grep tool names |
| 4 | **DF-21** — `BuiltinMCPCapabilitySpecs` includes tasks/capability/version slugs; **no** auto-seed | domain test |
| 5 | **DF-22** — README/help mention rebuild + abs `--bin` + reload MCP | docs |
| 6 | **DF-22** — `--write` stderr tip present; `trace_version` returns ok/name/version | install test + mcp test |
| 7 | **DF-32** — CLI `capability list`/`missing` JSON uses snake_case (`id` not `ID`) | CLI test / sample stdout |
| 8 | Tool count **nine**; `cmd/trace-mcp` help lists new tools | help + registration |
| 9 | S01 inherit intact (no regression of plan-change alias / capability Exact / IncludeWhy) | spot prior tests |
| 10 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + Gate C `dry_run:false` | locked verify cmds |
| 11 | Board Notes accurate; no product Go claimed on planner row | TODO.md |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Severity / spawn policy
- **blocker/high:** inline fix if tiny; else spawn `P10-S02-02a` (implement) + `P10-S02-02b` (review) immediately below this row
- **medium:** prefer spawn unless trivial
- **Residual OK:** Cursor cannot be force-restarted from Trace; agents must reload MCP manually (document residual). S04 still owns transition gating.

## Exit criteria
- [x] Checklist 1–11 evidenced
- [x] Confidence **high** (or **medium** with residuals listed — never silent)
- [x] No open blocker/high without pending follow-up
- [x] Board status + Notes; next **P10-S03-00** (unless spawn)

## Todo updates
Reviewer: status + notes; may spawn forward; may thicken **upcoming** S03/S05 prompts if blast radius requires. Do not edit `done` history.
