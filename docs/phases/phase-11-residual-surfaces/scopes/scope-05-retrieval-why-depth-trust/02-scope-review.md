# P11 / S05 / 02 — Scope review (Retrieval why / depth / trust / DPC attribution)

## Metadata
- id: P11-S05-02
- todo_ids: [P11-S05-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S05 (**DF-49, DF-35, DF-48, DF-42**). Fresh subagent. Compare claims + locks to live code/tests. Spawn `02a`/`02b` for blocker/high. Do not rewrite prior `done` history.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- Sibling [01-retrieval-why-depth-trust.md](01-retrieval-why-depth-trust.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-49, DF-35, DF-48, DF-42
- Live: `internal/retrieval/{exact,expand}.go`; `internal/compiler/packet.go`; `internal/domain/link.go`; `cmd/trace/link.go`; `internal/mcp/tools_write.go`
- Prior: P10 DF-19/27; P11-S04 no coupling

## Checklist (must all pass for APPROVE)

| # | Check |
|---|--------|
| 1 | DF-49: Exact/Why `symbol` by id works; miss OK; no mig |
| 2 | DF-35: depth≥2 TaskContext/ExpandContext omits sibling task **bodies** (titles OK); Expand goal→task no body Excerpt |
| 3 | DF-48: decision/assumption MD Law 9 honor + Law 4 channel; JSON `trust` stays `untrusted_data`; no TrustSystem elevate |
| 4 | DF-42: CLI `discovery-mentions-task` + MCP enum → store `discovery_mentions_task`; multi-goal DPC attribution via link; G19 thin adapters |
| 5 | DF-19 goal-scope + DF-27 title intent retained; no forbidden architecture |
| 6 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + P11-S01–S04 + Gate C `dry_run:false` |
| 7 | Board Notes accurate; planner row had no product Go |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./internal/store/... ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Prefer named asserts: `TestWhySymbolExact`, `TestExpandContextDepth2NoSiblingTaskBody`, updated `TestDecisionMarkdownTrustLabels`, `TestLinkDiscoveryMentionsTask` (+ CLI/MCP smoke), multi-goal DPC attribution (or equiv locked names).

## Exit criteria
- [x] Checklist evidenced; confidence high (or medium with residuals)
- [x] Board status + Notes; next **P11-S06-00** (unless spawn)
- [x] Write [REVIEW-NOTES.md](REVIEW-NOTES.md) on APPROVE / spawn
