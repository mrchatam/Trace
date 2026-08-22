# P12 / S01 / 01 — Edge provenance (FINAL)

## Metadata
- id: P12-S01-01
- todo_ids: [P12-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement structural **import** edge provenance (`EXTRACTED|INFERRED|AMBIGUOUS`) and surface it as JSON/`Markdown` field **`edge_provenance`** in Why + context per sibling **00-PLANNER** FINAL locks. **Stop if 00-PLANNER is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A1–A7
- Research rank 1 — graphify edge provenance
- Law 5 — inferred ≠ verified
- Live paths listed in 00-PLANNER inventory
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do not re-debate FINAL locks.

## Locked defaults (FINAL — do not renegotiate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Enum | `EXTRACTED` \| `INFERRED` \| `AMBIGUOUS` |
| Persistence | `imports.provenance` via mig **`011_import_edge_provenance.sql`** (`DEFAULT 'EXTRACTED'`) |
| Go types | `store.Import.Provenance`; `retrieval.Hit` / `WhyStep` field → JSON `edge_provenance`; `compiler.WhyTraceStep` (+ `Item` when structural hop) |
| Analyzer | AST imports → `EXTRACTED`; py wildcard → `AMBIGUOUS`; **no** analyzer-produced `INFERRED`; **no** new call edges |
| Causal fields | Untouched (`confidence` REAL / `Item.Provenance`) |
| Surfacing | Why JSON + context JSON/MD; MCP only if library already returns (no new MCP surface) |
| Forbidden | Fake precise calls; embeddings; Neo4j; daemon/HTTP; full-rebuild; board spawn; product edits outside listed packages without Notes |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Schema | `internal/store/schema/011_import_edge_provenance.sql` | ADD COLUMN |
| Store | `internal/store/file_graph.go` (+ tests) | Import field; Replace/List SQL |
| Analyzers | `internal/analyzers/extract_*.go` (+ tests) | Set Provenance on emit |
| Retrieval | `internal/retrieval/types.go`, `expand.go`, `why.go` (+ tests) | Copy onto Hit/WhyStep |
| Compiler | `internal/compiler/packet.go`, `compiler.go` (+ tests) | Pass-through + MD |
| CLI | Prefer **zero** edits if marshal is library-owned | Thin adapter only if needed (G19) |

## Role work
1. TDD: store round-trip → analyzers defaults → Expand/Why → compiler WhyTrace.
2. Prove `INFERRED` surfaces via **store fixture** (not analyzer) so agents cannot treat it as extracted.
3. Run locked verify suite; board **status + Notes only** (no prompt/board spawn).

## Minimal todos
- [ ] Mig 011 + `Import.Provenance` round-trip (`TestImportProvenanceRoundTrip`)
- [ ] Analyzers set EXTRACTED / AMBIGUOUS per locks
- [ ] Expand import→file + Why emit `edge_provenance`
- [ ] Compiler WhyTrace (+ Item if structural) + Markdown line
- [ ] Named tests green; carry-forward verify cmds below
- [ ] Board row P12-S01-01 Notes; next **P12-S01-02**

## Verify commands

```bash
CGO_ENABLED=0 go test ./internal/store/... ./internal/analyzers/... ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Exit criteria
- [ ] FINAL locks met with named tests above
- [ ] `edge_provenance` visible on Why + context WhyTrace for structural import hops; INFERRED fixture ≠ silent EXTRACTED
- [ ] No causal `confidence` / `Item.Provenance` semantic overwrite
- [ ] Carry-forward gates green; Gate C `dry_run:false` untouched
- [ ] Board Notes; next **P12-S01-02**
