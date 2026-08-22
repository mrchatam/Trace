# P12-S01-00 — Edge provenance (FINAL)

## Metadata
- id: P12-S01-00
- todo_ids: [P12-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S01 implement/review prompts for **structural edge provenance** (`EXTRACTED|INFERRED|AMBIGUOUS`) surfaced in Why/context. Live inventory confirmed; APIs/tests/migration locked. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law 5
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A1–A7; research rank 1
- [SIMILAR-PROJECTS-REVIEW-2026-08-16.md](../../../../research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md) — graphify provenance
- Live: `internal/analyzers`, `internal/store` (`file_graph.go`, `schema/001_init.sql`…`010_*`), `internal/retrieval` (`expand.go`, `why.go`, `types.go`), `internal/compiler` (`packet.go`, `compiler.go`), `cmd/trace/why.go`, `cmd/trace/context.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material phase locks A1–A7 hold; live schema conflict resolved below (column on `imports`, not causal `confidence`).

## Live inventory (2026-08-16)

| Area | Finding |
|------|---------|
| Structural edges | **`imports` only** — `id, file_id, imported_path, symbol`; **no** provenance column; **no** `calls`/`usages` table |
| Symbols | File-local stubs; not edges — **out of S01 enum home** |
| Causal confidence | `confidence REAL` on goals/tasks/… + `entity_links` + `files` — **do not overwrite / reuse for edge enum** (A4) |
| Analyzers | tree-sitter Go/JS/TS/Py → `ReplaceFileImports`; AST literals only today |
| Expand | `file` → symbols + resolved import→file neighbors (`ReasonGraphNeighbor`); unresolved imports **skipped** |
| Why | Builds steps from Expand; steps have `reason_code`, **no** edge provenance field |
| Context packet | `Item.Provenance` = causal `{status,source_type,confidence}` — **must not** host the edge enum; `WhyTraceStep` has no edge field yet |
| Migrations | Embedded max **010** → next file **`011_*.sql`** |
| MCP | Calls `retrieval.Why` / compiler library — parity free if library emits fields (G19/A6) |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Enum | `EXTRACTED` \| `INFERRED` \| `AMBIGUOUS` (exact strings) |
| Enum home | **`imports.provenance` TEXT NOT NULL** + `store.Import.Provenance string` |
| Migration | **Yes** — `internal/store/schema/011_import_edge_provenance.sql`: `ALTER TABLE imports ADD COLUMN provenance TEXT NOT NULL DEFAULT 'EXTRACTED'`; bump via existing embed migrate (version **011**) |
| Side table | **No** — column on `imports` (ReplaceFileImports already file-local / DR-INCREMENTAL) |
| Packages | `internal/store` → `internal/analyzers` → `internal/retrieval` → `internal/compiler` → CLI thin (`cmd/trace/why.go` / `context.go` unchanged if JSON marshal of library types) |
| JSON field | **`edge_provenance`** on `retrieval.Hit`, `retrieval.WhyStep`, `compiler.WhyTraceStep`, and on `compiler.Item` **only when** the item is a structural import hop (`reason_code=graph_neighbor` carrying import provenance). **Never** write the enum into `Item.Provenance.Confidence` / causal `confidence` |
| Analyzer production (S01) | Concrete AST imports → **`EXTRACTED`**; Python wildcard (`*`) / empty-path extracts → **`AMBIGUOUS`**; **do not produce `INFERRED`** in analyzers this scope (research rank 8 two-pass calls **deferred**) |
| Store contract | `ReplaceFileImports` / `ListImportsByPath` round-trip all three values; empty Provenance on write → default **`EXTRACTED`** |
| Retrieval surfacing | When Expand emits import→file neighbor, copy `Import.Provenance` onto Hit.`EdgeProvenance` (JSON `edge_provenance`); Why copies onto WhyStep; omit field when empty / non-structural |
| Context / MD | Pass-through on WhyTrace (+ Item when structural); Markdown why_trace / item lines show ``edge_provenance: `…` `` when set |
| MCP | Library-only parity — **no** new MCP tool / menu (A6) |
| Forbidden | Fake precise `calls` edges; inventing call/usage tables; embeddings; Neo4j; daemon/HTTP P0; full-rebuild indexer; conflating causal confidence; silent inferred-as-extracted; board spawn by implementer |
| Carry-forward | honesty A/B/C+G; E/F/ablation/H/compat; p0x; x0; Gate C `dry_run:false`; Phase 11 DF regressions |
| Named tests (min) | `TestImportProvenanceRoundTrip` (store); `TestAnalyzerImportProvenanceExtracted` (+ AMBIGUOUS wildcard if py fixture easy); `TestExpandImportEdgeProvenance` / `TestWhySurfacesEdgeProvenance`; `TestContextWhyTraceEdgeProvenance` (compiler); store fixture proving **INFERRED** surfaces when inserted via ReplaceFileImports |
| Verify | See `01-edge-provenance.md` |

## Owns
| Item | Intent |
|------|--------|
| Edge provenance enum | Persist + emit on structural **import** edges |
| Why/context display | Agents see hop/edge confidence; inferred ≠ extracted |

## Explicit deferrals (not S01)
- Two-pass call resolve / `INFERRED` call edges (research rank 8)
- Report mix-% / “verify INFERRED” UX (rank 17) beyond field presence
- Emitting Why steps for **unresolved** imports (remain Expand-skipped)
- Research S03+ / ranks 4+

## Planner work (this row)
1. [x] Inventory live edge/import schema + Why/context paths
2. [x] Lock FINAL defaults (mig 011; column; tests; JSON name)
3. [x] Thicken `01-edge-provenance.md` + `02-scope-review.md` + SCOPE-TODOS
4. [x] Light Depends note for S02 (packets may cite `edge_provenance` if present)

## Exit
- [x] Thicken 01 + 02 + SCOPE-TODOS to FINAL
- [x] Board Notes; next **P12-S01-01**
- [x] Product Go — **not** this row
