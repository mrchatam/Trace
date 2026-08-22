# P13 / S03 / 01 — Provenance schema residuals (FINAL)

## Metadata
- id: P13-S03-01
- todo_ids: [P13-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-64** (write validation + mig **012** CHECK + read normalize) and close residuals per sibling **00-PLANNER** FINAL: **DF-66** documented wontfix (docs + keep Law 5 fixtures), **DF-67** explicit out-of-bar (no code). **Stop if 00-PLANNER is still DRAFT.** Do not regress P12 provenance round-trip / Expand/Why/compiler tests.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A5
- DOGFOOD-FINDINGS DF-64/66/67; `_bughunt/post-p12/{prov,symstale}/`
- Live: `internal/store/{file_graph.go,schema/011_import_edge_provenance.sql}`; `evals/compat`
- Docs target: [`docs/ANALYZER_CONTRIBUTION.md`](../../../../../ANALYZER_CONTRIBUTION.md)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do not re-debate FINAL locks.

## Locked defaults (FINAL — do not renegotiate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Home | `internal/store` primary; thin `docs/ANALYZER_CONTRIBUTION.md` for DF-66 |
| DF-64 write | Empty → `EXTRACTED`; only `EXTRACTED`\|`INFERRED`\|`AMBIGUOUS`; garbage → **error** |
| DF-64 read | Empty → `EXTRACTED` on list/read normalize |
| DF-64 mig | **`012_import_provenance_enum.sql`**: rebuild `imports` + `CHECK (...)`; heal `''`/unknown → `EXTRACTED` on copy |
| Compat | Ceiling **12** (011 present; 012 OK; no 013+) |
| DF-66 | **wontfix** analyzer/CLI setter; keep store-fixture Law 5 tests; docs paragraph only |
| DF-67 | **No implement** — residual for VERIFY Notes only |
| Packet SchemaVersion | Keep **`0.2`** |
| Analyzer / retrieval / compiler | Prefer **zero** edits |
| MCP / CLI | No provenance command (G19) |
| Forbidden | Silent write coerce of garbage; call-graph INFERRED; symbol honesty; daemon; full-rebuild; board spawn |
| Keep green | `TestImportProvenanceRoundTrip`, `TestAnalyzerImportProvenanceExtracted`, `TestExpandImportEdgeProvenance`, `TestWhySurfacesEdgeProvenance`, `TestContextWhyTraceEdgeProvenance`, S01/S02 import/honesty tests |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Store schema | `internal/store/schema/012_import_provenance_enum.sql` | **New** — rebuild + CHECK + heal |
| Store | `internal/store/file_graph.go` | Validate on `ReplaceFileImports`; normalize empty on read |
| Store tests | `internal/store/store_test.go` (+ optional focused `*_test.go`) | Garbage reject; empty→EXTRACTED; round-trip keep |
| Compat | `evals/compat/compat_test.go` (+ doc comments if needed) | Ceiling 11→12; allow 012; forbid 013+ |
| Docs | `docs/ANALYZER_CONTRIBUTION.md` | DF-66: no analyzer INFERRED; fixture-only; no CLI setter |
| Analyzers / retrieval / compiler | — | **Prefer zero** |

## Role work
1. TDD: garbage write fails; empty write/read → EXTRACTED; three-enum round-trip still green.
2. Ship mig 012 + write validate/normalize; bump compat ceiling.
3. DF-66: docs paragraph only; prove Law 5 fixture tests still green (no CLI).
4. DF-67: **no code** — Note residual for reviewer/VERIFY.
5. Run locked verify; board **status + Notes only** (no prompt/board spawn).

## Minimal todos
- [ ] Mig `012_import_provenance_enum.sql` + embed applies to 12
- [ ] `ReplaceFileImports` reject garbage; empty→EXTRACTED
- [ ] Read path normalize empty→EXTRACTED
- [ ] Named tests for reject + empty honesty + keep P12 round-trip/Expand/Why
- [ ] Compat ceiling 12 (no 013+)
- [ ] DF-66 docs paragraph in `ANALYZER_CONTRIBUTION.md`
- [ ] DF-67: no symbol honesty code; residual noted on board
- [ ] Carry-forward verify cmds below
- [ ] Board row P13-S03-01 Notes; next **P13-S03-02**

## Named tests (intent locked)

| Test | Intent |
|------|--------|
| `TestImportProvenanceRoundTrip` | **Keep** |
| `TestAnalyzerImportProvenanceExtracted` | **Keep** |
| `TestExpandImportEdgeProvenance` | **Keep** — INFERRED fixture |
| `TestWhySurfacesEdgeProvenance` | **Keep** |
| `TestContextWhyTraceEdgeProvenance` | **Keep** |
| New (name free) | Garbage provenance → `ReplaceFileImports` error |
| New (name free) | Empty provenance write → stored/read `EXTRACTED`; Expand surfaces `EXTRACTED` (not hidden) |

## Verify commands

```bash
# DF-64 store + provenance surface (CGO0)
CGO_ENABLED=0 go test ./internal/store/... ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1

# Named focus (adjust -run to match new names)
CGO_ENABLED=0 go test ./internal/store/... -count=1 -run 'TestImportProvenance|TestReplace.*[Pp]rov|TestProvenance'
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... -count=1 -run 'TestExpandImportEdgeProvenance|TestWhySurfacesEdgeProvenance|TestContextWhyTraceEdgeProvenance'

# Compat ceiling 12
CGO_ENABLED=1 go test ./evals/compat/... -count=1

# Carry-forward (CGO1)
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1

# Product packages
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Do **not** treat `go test ./...` space-path FAIL (known graphify) as S03 failure. Gate C `dry_run:false` untouched. Dry-run ≠ Gate C / ≠ H / ≠ checklist.

## Exit criteria
- [ ] FINAL locks met; DF-64 write+CHECK+read harden shipped
- [ ] DF-66 docs + fixture path retained (wontfix product setter)
- [ ] DF-67 no symbol honesty code; residual explicit in Notes
- [ ] Compat ceiling 12; P12/S01/S02 provenance/honesty keepers green
- [ ] Carry-forward gates green; Gate C `dry_run:false` untouched
- [ ] Board Notes; next **P13-S03-02**

## Next
**P13-S03-02**
