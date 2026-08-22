# P10 / S04 / 02 — Scope review (operator / capability gates)

## Metadata
- id: P10-S04-02
- todo_ids: [P10-S04-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S04 (**DF-17 / DF-18 / DF-24 / DF-26 / DF-31**). Fresh subagent. Compare claims + locks to live code/tests. Honesty escape-rate (Gate G) and DONE=PASS\|hatch must stay coherent with the **operator-flag** supersession. Small inline fix **or** spawn `02a`/`02b` for blocker/high. Do not rewrite S01–S03/`done` history.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) FINAL locks + [01-operator-capability-gates.md](01-operator-capability-gates.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- [experiments/ab-operator-gate/results/G1/PROBE-AGENT-DONE.md](../../../../../experiments/ab-operator-gate/results/G1/PROBE-AGENT-DONE.md)
- Live: `internal/domain/task_state.go`, `cmd/trace/transition.go`, `internal/mcp/tools_write.go`, `evals/honesty/`

## Session start
Agent → clarify → Plan → execute (reviewer).

## Checklist (must all pass for APPROVE)

| # | Check | Evidence |
|---|--------|----------|
| 1 | **DF-17** — →DONE needs hatch **or** (PASS ∧ `AllowOperatorDone`); Actor alone insufficient | domain test |
| 2 | **DF-17** — CLI `--as-operator` + MCP `as_operator` wired (G19) | CLI/MCP tests |
| 3 | **DF-17** — honesty Path C uses operator flag; A/B still reject | honesty test |
| 4 | **DF-18** — leave DONE invalidates linked PASS → `UNCERTAIN`; sticky PASS gone | domain test |
| 5 | **DF-24** — missing caps block transitions; `--allow-missing-caps` / `allow_missing_caps` override | domain + adapter tests |
| 6 | **DF-26** — hatch success emits CLI stderr WARNING + MCP `"warning"`; hatch still works for Gate G | CLI/MCP + honesty G |
| 7 | **DF-31** — `capability missing` without task is usable (hint / clear MCP error) | CLI/MCP |
| 8 | No new mig; no Actor-name allowlist; nine MCP tools retained | grep / schema |
| 9 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + S01–S03 + Gate C `dry_run:false` | locked verify cmds |
| 10 | Board Notes accurate; planner row had no product Go | TODO.md |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Severity / spawn policy
- **blocker/high:** inline fix if tiny; else spawn `P10-S04-02a` (implement) + `P10-S04-02b` (review) immediately below this row
- **medium:** prefer spawn unless trivial
- **Residual OK:** ab-operator-gate experiment re-run is optional evidence (not required if unit/MCP tests cover probe shape); Cursor MCP reload still manual (S02 residual)

## Exit criteria
- [x] Checklist 1–10 evidenced
- [x] Confidence **high** (or **medium** with residuals listed — never silent)
- [x] No open blocker/high without pending follow-up
- [x] Board status + Notes; next **P10-S05-00** (unless spawn)

## Todo updates
Reviewer: status + notes; may spawn forward; may thicken **upcoming** S05 prompts if blast radius requires. Do not edit `done` history.
