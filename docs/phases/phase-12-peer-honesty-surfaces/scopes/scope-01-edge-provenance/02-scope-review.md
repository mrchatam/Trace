# P12 / S01 / 02 — Scope review (Edge provenance) FINAL

## Metadata
- id: P12-S01-02
- todo_ids: [P12-S01-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S01 edge provenance. Fresh subagent. Compare claims + **00-PLANNER FINAL locks** to live code/tests. Spawn `02a`/`02b` for blocker/high.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-edge-provenance.md](01-edge-provenance.md)
- [phase README](../../README.md)
- Research rank 1 — graphify
- Law 5 / A4 — structural enum ≠ causal confidence

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (reviewer).

## Checklist (must all pass for APPROVE)

| # | Check |
|---|--------|
| 1 | Mig **011** + `imports.provenance` persists `EXTRACTED\|INFERRED\|AMBIGUOUS`; default EXTRACTED for pre-mig rows |
| 2 | Analyzers: AST → EXTRACTED; wildcard → AMBIGUOUS; **no** fake call edges / no analyzer INFERRED required |
| 3 | Why + context surface JSON/`Markdown` **`edge_provenance`** on structural import hops; INFERRED fixture not silent-as-EXTRACTED |
| 4 | Causal `confidence` / `Item.Provenance` **not** overloaded with the enum |
| 5 | G19 — no domain fork in CLI/MCP adapters; MCP parity only via library fields |
| 6 | No forbidden architecture (daemon/HTTP/embeddings/Neo4j/full-rebuild) |
| 7 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + Gate C `dry_run:false` |
| 8 | Board Notes accurate; planner row had no product Go |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./internal/analyzers/... ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Exit criteria
- [ ] Checklist evidenced; confidence high (or medium with residuals listed)
- [ ] Board status + Notes; next **P12-S02-00** (unless spawn)
