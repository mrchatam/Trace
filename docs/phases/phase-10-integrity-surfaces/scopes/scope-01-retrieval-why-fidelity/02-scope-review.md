# P10 / S01 / 02 — Scope review (retrieval / why fidelity)

## Metadata
- id: P10-S01-02
- todo_ids: [P10-S01-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S01 (**DF-19, DF-23, DF-25, DF-27, DF-29**) against **00-PLANNER** FINAL locks and P10-S01-01 Notes. Fresh subagent ≠ implementer. Quality-first: blocker/high → inline fix or spawn `02a`/`02b`.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-retrieval-why-fidelity.md](01-retrieval-why-fidelity.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) Law 4 + Law 9
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- Implementer board Notes on P10-S01-01

## Session start
Agent → clarify → Plan → execute (reviewer).

## Focus checklist

| # | Check | Evidence |
|---|--------|----------|
| 1 | **DF-19** — task Expand/Why no longer dumps all-project DPC | Code path in `discovery_plan_change.go` / `expand.go`; multi-goal test |
| 2 | **DF-19** — single-goal / pair-completion / foreign-task rules match locks | Tests + algorithm vs 00-PLANNER table |
| 3 | **DF-23** — `plan-change` accepted; emitted `plan_change` | Exact/Why/CLI or MCP why; hit EntityType |
| 4 | **DF-25** — `capability` Exact/Why works via `GetCapability` | lookupEntity case + test |
| 5 | **DF-25 residual** — `plan_scope` **not** silently “fixed” or claimed done | Notes list residual |
| 6 | **DF-27** — MD no longer says decisions are “not project policy”; `trust` still `untrusted_data`; no TrustSystem elevation | packet.go + test |
| 7 | **DF-29** — IncludeWhy=true propagates Why errors | compiler.go + failing-Why test |
| 8 | **No new mig** / no new MCP tools / no daemon | `git` / schema dir / mcp tool list |
| 9 | Carry-forward green | honesty / Gates E–H / ablation / compat / p0x / x0 / `./...`; Gate C `dry_run:false` |
| 10 | Board Notes accurate; S02 Depends coherent | TODO.md |

## Verify (independent)

```bash
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Re-read implementer-cited tests; do not trust Notes alone.

## Severity / spawn rules
- **blocker/high:** inline small fix **or** spawn `P10-S01-02a` (implement) + `P10-S01-02b` (review) immediately below this row; thicken prompts.
- **medium:** prefer spawn unless trivial.
- **Residual OK:** `plan_scope` Exact unknown; Mode-B packs historical; S02+ owns MCP tool parity.

## Confidence bar
Prefer **high**. **medium** only with explicit residuals listed (never silent).

## Exit criteria
- [ ] Checklist 1–10 evidenced
- [ ] No open blocker/high without pending follow-up
- [ ] REVIEW-NOTES.md written (verdict, confidence, residuals)
- [ ] Board status + Notes; next **P10-S02-00** (unless spawn)
- [ ] Confidence medium or high

## Todo updates
Reviewer: status + notes; may spawn forward; may thicken **upcoming** S02 prompts if S01 changes blast radius. Do not edit `done` history.
