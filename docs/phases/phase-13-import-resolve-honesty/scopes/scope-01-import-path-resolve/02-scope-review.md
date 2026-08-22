# P13 / S01 / 02 — Scope review (Import path resolve) FINAL

## Metadata
- id: P13-S01-02
- todo_ids: [P13-S01-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S01 DF-60. Fresh subagent. Compare **00-PLANNER FINAL locks** + implementer claims to live code/tests. Spawn `02a`/`02b` for blocker/high.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-import-path-resolve.md](01-import-path-resolve.md)
- [phase README](../../README.md)
- DF-60 / POST-P12-BUGHUNT / ab-import-resolve
- Phase A3 — resolve-time, not analyzer re-key

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (reviewer).

## Checklist (must all pass for APPROVE)

| # | Check |
|---|--------|
| 1 | Subdir `./` / `../` imports resolve via **importer-dir join** → Expand neighbor + `edge_provenance` (store fixture and/or live-shaped strings) |
| 2 | Extensionless `./util` tries locked extensions (+ `index.*` when applicable) and hits indexed file |
| 3 | Root positive control `./util.js`→`util.js` still green with `edge_provenance` |
| 4 | P12 `TestExpandImportEdgeProvenance` / `TestWhySurfacesEdgeProvenance` green; causal `confidence` / `Item.Provenance` untouched |
| 5 | Home is retrieval resolve — **no** analyzer rewrite of `imported_path`, **no** product path-align hook, **no** new migration |
| 6 | Bare modules (`fmt`, `node:fs`, `github.com/…`) still Expand-skipped (no fake module resolver) |
| 7 | Forbidden architecture absent (full-rebuild / daemon/HTTP / new MCP / fake calls); G19 intact |
| 8 | Carry-forward + Gate C `dry_run:false` intact; board Notes accurate; no implementer board spawn |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run 'TestExpandImportEdgeProvenance|TestWhySurfacesEdgeProvenance|TestExpand.*[Ii]mport|TestWhy.*[Ii]mport|TestResolve'
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Optional (not blocker): `experiments/ab-import-resolve/` prepare + why probe for `EXTRACTED` on `src/util.js`.

## Exit criteria
- [ ] Checklist evidenced; confidence high (or medium with residuals listed)
- [ ] [REVIEW-NOTES.md](./REVIEW-NOTES.md) written
- [ ] Board status + Notes; next **P13-S02-00** (unless spawn)
