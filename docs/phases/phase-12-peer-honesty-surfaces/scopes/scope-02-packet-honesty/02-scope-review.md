# P12 / S02 / 02 — Scope review (Packet honesty) FINAL

## Metadata
- id: P12-S02-02
- todo_ids: [P12-S02-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S02 packet honesty. Fresh subagent. Compare claims + **00-PLANNER FINAL locks (2026-08-17 retry)** to live code/tests. Expect **partial pre-ship** of types + **named tests** from P12-S02-01. Spawn `02a`/`02b` for blocker/high.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-packet-honesty.md](01-packet-honesty.md) — inventory + gap finish
- [phase README](../../README.md)
- Research ranks 2–3
- S01 [REVIEW-NOTES.md](../scope-01-edge-provenance/REVIEW-NOTES.md) — `edge_provenance` must stay green

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (reviewer). Re-inventory briefly; do not trust implementer Notes alone.

## Checklist (must all pass for APPROVE)

| # | Check |
|---|--------|
| 1 | SchemaVersion **`0.2`**; Budget has `items_total` / `items_kept` / `candidates_capped`; MD shows `items=kept/total` |
| 2 | `truncated=true` when kept&lt;total **or** MaxCandidateHits capped — **no silent caps** |
| 3 | `index_honesty.stale_paths` when kept file disk sha256 ≠ `content_hash`; false-fresh on errors; MD banner when set; prefer **sort-then-cap 8** |
| 4 | Causal `Provenance.Status` / Law 18 STALE **not** mutated from index drift |
| 5 | S01 intact — `edge_provenance` pass-through; not written into causal `confidence` / `Item.Provenance` |
| 6 | G19 — thin adapters only; no new MCP tool menu |
| 7 | No forbidden architecture (daemon FS watcher product; embeddings; Neo4j; full-rebuild; size/mtime migration this scope) |
| 8 | Carry-forward + Gate C `dry_run:false`; **named tests PASS**; board Notes accurate; **P12-S02-00 planner did not author product Go on retry** (pre-existing honesty types from earlier partial ship are OK if locks hold) |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/compiler/... ./internal/retrieval/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Named: `TestBudgetLoudTotals` / `TestCandidateCapSetsTruncated` / `TestIndexStaleBanner` / `TestContextWhyTraceEdgeProvenance`

## Focus answers (reviewer must state)
- Where totals are computed (pre-trim `items_total` vs post-keep `items_kept`) — expect `compiler.go` `compileAtDepth`.
- How false-fresh is enforced on I/O failure — expect `index_honesty.go` omit/continue.
- Confirmation S01 EdgeProvenance fields still wired.
- Whether implementer only filled gaps (tests ± small fixes) vs unnecessary rewrite of shipped types.

## Exit criteria
- [ ] Checklist evidenced; confidence high (or medium with residuals listed)
- [ ] Write `REVIEW-NOTES.md` on APPROVE
- [ ] Board status + Notes; next **P12-S03-00** (unless spawn)
