# P07 / S02 / 01 — Language plugins / adapters (Go)

## Metadata
- id: P07-S02-01
- todo_ids: [P07-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Add **one** additional language adapter — **Go** — via the existing `internal/analyzers` boundary (DetectLanguage + extract + IndexFile). Official tree-sitter grammar only. Keep CGO analyzers-only. Do **not** regress S01 T0 walk/ignore or file-local incremental. Golden tests required. No Gate H threshold invent.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks finalized 2026-08-16
- [phase README](../../README.md)
- Live: `internal/analyzers/{detect,extract,index,extract_*.go}.go`; `cmd/trace` T0 walk (`isT0SkipDir`/`isT0SkipPath`/`walkIndexable`)
- Prior pattern: P00-S04 `01-analyzers.md` (JS/TS/Python)

## Session start
Agent → clarify if needed → Plan → execute.

## Locked defaults (FINAL — P07-S02-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go floor | Keep `go 1.24.0` in `go.mod` (do **not** downgrade) |
| Package | **`internal/analyzers` only** for product changes (package `analyzers`) |
| First additional language | **Go** (self-host Trace) — **exactly one** new language in this row |
| Official grammar module | **`github.com/tree-sitter/tree-sitter-go` `v0.25.0`** |
| Binding import | `github.com/tree-sitter/tree-sitter-go/bindings/go` → `Language()` (same shape as Python) |
| Runtime binding | Keep existing `github.com/tree-sitter/go-tree-sitter` `v0.25.0` (MVS may pull grammar’s older require; Trace pin wins) |
| Language id | `LangGo = "go"` (persist on `files.language`) |
| Extension | **`.go` only** (case-insensitive via existing `DetectLanguage` Ext lowercasing) |
| Adapter shape | Extend **existing** switches — `DetectLanguage` + `extract(lang)` — **not** a plugin registry, dynamic loader, or universal language theater |
| Extract file | Prefer `extract_go.go` (safe; avoid GOOS traps like `*_js.go`) |
| IndexFile / IndexFileAtRev | Unchanged orchestration; new lang flows through existing Upsert → SetFileLanguage → extract → ReplaceSymbols/Imports |
| Symbols (minimal) | `function` = top-level funcs; `method` = funcs with receiver; `type` = type declarations when cheap via query (struct/interface/alias all as `type`) |
| Imports (minimal) | Go `import` / import-block path strings → `ImportedPath`; `Symbol` nil unless a named alias is cheap |
| Lines | 1-based inclusive via existing `nodeLines` helpers |
| Skip policy | Unsupported ext still SkipError / walk skip; binary NUL heuristic unchanged |
| CGO | Analyzers-only; store/domain/vcs/gitcli remain `CGO_ENABLED=0`-clean |
| Migration | **No** `011_*` — analyzer extension only |
| S01 T0 | **Do not regress.** Walk order stays T0 dir → DetectLanguage → T0 file/path → gitignore; `vendor` remains T0 SkipDir (Go-friendly). Do not rewrite T0 lists unless a bug blocks Go |
| CLI | No new subcommands; walk already uses `DetectLanguage` — `.go` appears automatically once DetectLanguage accepts it |
| Golden tests | **Required:** DetectLanguage `.go`; IndexFile Go golden symbols+imports; keep existing JS/TS/Py goldens green |
| Gate H | May add tiny Go fixtures later under `evals/perf` **only if** S03 needs them — **no** pass thresholds / Gate H claim in this row |
| Carry-forward | Honesty A/B/C; Gates E/F/G; capability ablation; p0x; x0; Gate C `dry_run:false` intact |
| Surface | No daemon/HTTP/MCP/embeddings as primary |

## Extension points (exact)

1. **`detect.go`** — add `LangGo`; map `.go` → `LangGo`.
2. **`index.go` `extract` switch** — `case LangGo: return extractGo(content)`.
3. **`extract_go.go`** — tree-sitter parse + queries for symbols/imports (mirror `extract_py.go` style).
4. **`analyzers_test.go` (+ optional `testdata/*.go`)** — golden + flip DetectLanguage case that currently expects `.go` unsupported.
5. **`go.mod` / `go.sum`** — `require github.com/tree-sitter/tree-sitter-go v0.25.0`.

Do **not** invent `LanguagePlugin` interfaces, registries, or per-language packages outside `internal/analyzers` in this row.

## Role work

1. TDD: DetectLanguage + Go golden failing → implement extractGo + wiring.
2. Preserve file-local incremental (reindex path A clears A only) — existing isolation tests must stay green.
3. Run locked verify suite; update board Notes only.

## Verify commands (locked)

```bash
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... -count=1
CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/... ./internal/domain/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Optional spot-check: Gate C artifacts under `docs/verification/gate-c-x0/` remain `dry_run:false` N=3 (do not rewrite packs).

## Exit criteria
- [ ] `.go` files indexable via `IndexFile` with language `"go"`
- [ ] Official `tree-sitter-go` **v0.25.0** in `go.mod`; CGO confined to analyzers path
- [ ] Golden tests cover Go symbols + imports; DetectLanguage accepts `.go`
- [ ] Adapter-shaped switch extension only (no universal plugin theater)
- [ ] S01 T0 walk/order/tests not regressed
- [ ] Carry-forward suite green; no Gate H pass/threshold invent
- [ ] Board Notes ready for **P07-S02-02**

## Out of scope
- Second additional language in this row
- Gate H thresholds / declaring Gate H pass
- Rewriting JS/TS/Python extractors without need
- Plugin registry / dynamic grammar loading
- Daemon/HTTP/MCP/embeddings
- Store schema migrations
- Changing T0 ignore lists (unless blocked — then Note + minimal fix)

## Minimal todos
- [ ] Add `tree-sitter-go` v0.25.0 + `LangGo` / DetectLanguage `.go`
- [ ] Implement `extractGo` + `extract` switch case
- [ ] Golden + DetectLanguage tests
- [ ] Run verify commands; board Notes → next **P07-S02-02**
