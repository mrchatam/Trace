# P13-S01-00 — Import path resolve (FINAL)

## Metadata
- id: P13-S01-00
- todo_ids: [P13-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S01 implement/review prompts for **DF-60**: Expand/Why resolve relative/module import strings by joining the importer directory, normalizing `./`, and handling common extensions so live subdir trees emit `edge_provenance`. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A1–A8 (esp. A3 resolve-time)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-60
- [experiments/_bughunt/post-p12/POST-P12-BUGHUNT.md](../../../../../experiments/_bughunt/post-p12/POST-P12-BUGHUNT.md)
- Live: `internal/retrieval/expand.go` (file→imports loop), `store.NormalizePath`, analyzers import emit
- Optional dogfood: [`experiments/ab-import-resolve/`](../../../../../experiments/ab-import-resolve/)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks from phase A3 hold; live gap confirmed below.

## Live inventory (2026-08-17)

| Area | Finding |
|------|---------|
| Expand import loop | `expand.go` `case "file"`: `ListImportsByPath` then **`GetFileByPath(store.NormalizePath(imp.ImportedPath))` only** — unresolved → `continue` (no neighbor, no `edge_provenance`) |
| `NormalizePath` | `entities.go`: `\`→`/`, **`TrimPrefix("./")` only** — no importer-dir join, no `../` clean, no extensions |
| P12 named fixtures | `TestExpandImportEdgeProvenance` / `TestWhySurfacesEdgeProvenance` insert **repo-exact** `ImportedPath` (`src/util.go`, `why/b.go`) — bypass live analyzer strings → green today |
| Analyzer emit | JS/TS: raw AST source (`./util`, `./util.js`) via `extract_javascript.go`; Go: import path literal (`github.com/…` or rare `./…`); Py: module path / wildcard → AMBIGUOUS. **No** project-relative rewrite at index |
| Root positive | `_bughunt/post-p12/rootimp/`: `./util.js` → Normalize → `util.js` matches root file → Why shows `EXTRACTED` |
| Subdir fail | `prov/src/app.ts` → `./util`; `whyfile/src/app.js` → `./util.js`; `ab-import-resolve` `src/app.js` → `./util.js` — Normalize → bare `util`/`util.js` ≠ `src/…` → **no hop** |
| Go module paths | `whyfile/pkg/a.go` → `github.com/example/whyfile/pkg/b` — **out of S01** (no module/GOPATH resolve); remain Expand-skipped |
| Schema | No new migration — resolve is retrieval-time only (A3) |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Home | **`internal/retrieval`** — resolve helper + Expand call site (+ tests). **No** analyzer rewrite of stored `imported_path`; **no** sqlite path-align hooks as product fix |
| API | Unexported helper (name free; e.g. `resolveImportedFile` / `importPathCandidates`) used **only** from Expand file→import loop. Signature intent: given **importer file path** + raw `Import.ImportedPath` → first indexed `*store.File` or not-found |
| Relative detect | On **raw** string: prefix `./` or `../` (also accept `.\\` / `..\\` before slash normalize). Bare modules (`fmt`, `node:fs`, `github.com/…`) → **no join**; try NormalizePath exact only, else skip |
| Candidate order | (1) `NormalizePath(imported)` exact; (2) if relative: `NormalizePath(path.Clean(join(Dir(importer), imported)))`; (3) for each base **without** a known source suffix, append same base + extensions; (4) for extensionless joined base, also try `base/index` + extensions. Dedupe preserving order; **first `GetFileByPath` hit wins** |
| Extensions | `.js`, `.jsx`, `.mjs`, `.cjs`, `.ts`, `.tsx`, `.go`, `.py` |
| Positive control | Root `./util.js`→`util.js` **must stay green**; P12 exact-path fixtures **must stay green** |
| Provenance | Unchanged — copy `imp.Provenance` onto Hit when neighbor resolves (P12 contract) |
| Packages touched | Prefer **`internal/retrieval` only** (`expand.go` + new helper file optional + `retrieval_test.go`). Compiler/CLI/MCP untouched if Expand already feeds Why |
| Migration | **No** |
| Forbidden | Full-rebuild indexer; fake call edges; product path-align / rewriting analyzer strings as the fix; new MCP tools; daemon/HTTP; conflating causal confidence; board spawn by implementer |
| Carry-forward | P12 provenance named tests; honesty A/B/C+G; E/F/ablation/H/compat; p0x; x0; Gate C `dry_run:false`; Law 18 causal STALE untouched |
| Named tests (min) | **Keep** `TestExpandImportEdgeProvenance`, `TestWhySurfacesEdgeProvenance`. **Add** (names free but intent locked): subdir `./util.js` → Expand neighbor + `EXTRACTED`; subdir extensionless `./util` → file with `.ts`/`.js`/…; root positive control Expand/Why still emits `edge_provenance`. Prefer store fixtures (analyzer-shaped strings) under CGO0 |
| Optional dogfood | `experiments/ab-import-resolve/` prepare → why on `src/app.js` shows util hop — **not** board blocker (A8) |
| Verify | See `01-import-path-resolve.md` |

## Owns
| DF | Intent |
|----|--------|
| DF-60 | Subdir `./util` / `./util.js` resolve to indexed file → Why/context show `edge_provenance` without sqlite path-align hooks |

## Explicit deferrals (not S01)
- Node/Go module / `node_modules` / GOPATH resolve for bare specifiers
- Re-keying analyzer `imported_path` to repo-relative at index time
- Emitting Why steps for still-unresolved imports
- DF-65 context hop **admission** (S02 — may rely on this resolve)
- DF-61…64, 66–67

## Depends note (S02)
DF-65 needs Expand to admit import→file neighbors from **analyzer-native** `./` strings. S01 ships that resolve; S02 must not re-implement path join — see light note on S02 SCOPE-TODOS.

## Planner work (this row)
1. [x] Live inventory Expand + NormalizePath + analyzer stored strings vs P12 fixtures
2. [x] Lock FINAL APIs/tests; thicken `01-import-path-resolve.md` + `02-scope-review.md` + SCOPE-TODOS
3. [x] Light Depends note for S02 (context hops may rely on resolve)

## Exit
- [x] FINAL locks; next **P13-S01-01**
- [x] Product Go — **not** this row
