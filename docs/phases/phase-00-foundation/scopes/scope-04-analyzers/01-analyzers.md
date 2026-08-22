# P00 / S04 / 01 — tree-sitter analyzers

## Metadata
- id: P00-S04-01
- todo_ids: [P00-S04-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob, user-context7]
- agents: []
- verification: automated

## Objective
Emit **File + minimal Symbol/Import** edges via **tree-sitter** for **TS/JS and Python**, with **file-level incremental** updates (P0-X #2 and #7 / DR-INCREMENTAL). Persist only through S02 store per-path APIs. No full-project rebuild path.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [C_FIRST_SCOPE.md](../../../../init/C_FIRST_SCOPE.md) — structural graph + P0-X #2/#7
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G1 (no blobs), G12 (incremental)
- [D_DECISION_REGISTER.md](../../../../init/D_DECISION_REGISTER.md) — DR-PARSE, DR-ANLANG, DR-INCREMENTAL, DR-P0X, DR-RISK, DR-SURFACE
- [STORAGE_AND_PERFORMANCE.md](../../../../STORAGE_AND_PERFORMANCE.md) §5 incremental indexing
- [ARCHITECTURE.md](../../../../ARCHITECTURE.md) — Code analyzer adapter
- Historical T004: [B_INITIAL_BOARD.md](../../../../init/B_INITIAL_BOARD.md)
- Prior live state: S02 `UpsertFile` / `ReplaceFileSymbols` / `ReplaceFileImports` / `ListSymbolsByPath` / `ListImportsByPath`; S03 `vcs.Repository` (`ShowFile`, `Head`, `Fake`); `internal/analyzers/doc.go` stub; `go.mod` **go 1.24.0**

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | `go 1.24.0` in `go.mod` (sqlite floor; do **not** downgrade) |
| Package path | `internal/analyzers` (package `analyzers`) — **only** product package for this scope |
| Parser backend | **tree-sitter** (DR-PARSE) via official Go bindings |
| Binding | `github.com/tree-sitter/go-tree-sitter` |
| Grammars (separate modules) | `github.com/tree-sitter/tree-sitter-javascript`; `github.com/tree-sitter/tree-sitter-typescript` (typescript **and** tsx languages); `github.com/tree-sitter/tree-sitter-python` |
| Forbidden parsers | Regex-only “temporary” graph; SCIP/Graphify as primary; other languages |
| Languages | **TS/JS + Python only** (DR-ANLANG). Extensions: `.ts`/`.tsx` → typescript/tsx; `.js`/`.jsx`/`.mjs`/`.cjs` → javascript; `.py` → python |
| Persistence | **Only** S02 APIs: `UpsertFile(path, contentHash, gitOID)` + `ReplaceFileSymbols` + `ReplaceFileImports`. **Do not** invent parallel File/Symbol/Import tables |
| Content hash | Hex SHA-256 of raw file bytes (`crypto/sha256`); never store bodies in SQLite |
| Language column | After upsert, set `files.language` via a thin store helper (e.g. `SetFileLanguage(path, lang)`) — **allowed** additive store API; no new migration. Do not leave language forever NULL for indexed supported files |
| Content source | Primary: caller-supplied `[]byte`. Secondary: `vcs.Repository.ShowFile(ctx, rev, path)` wrapper. Analyzers import `vcs` iface + use `vcs.Fake` in tests — **never** import `gitcli` |
| Incremental rule (chosen) | **File-local only:** reindexing path A replaces symbols/imports for A only; **no** cascade reindex of importers/dependents in S04. Document in package comment. (T004 allowed either rule; this is the locked choice for P0-X) |
| Skip policy | Unsupported extension → typed skip/error (batch walkers skip; single `IndexFile` returns clear error). Binary heuristic: NUL byte in first 8KiB → skip/error. Explicit `IndexFile` does **not** need to re-check `.gitignore` (caller/CLI owns walk); if a walk helper is added, respect gitignore |
| CGO | Official bindings use CGO. **Allowed for `internal/analyzers` only.** Store/vcs/gitcli must remain `CGO_ENABLED=0`-clean. Exit tests: analyzers with CGO on; non-analyzer packages still pass with `CGO_ENABLED=0` |
| Surface | Library only — **no** new `cmd/trace` subcommands (S07); **no** MCP/daemon/HTTP |
| Out of scope | Perfect call graphs; cross-file name resolution; other languages; semantic LLM summaries; full-rebuild-as-default; CLI `index` wiring |

### Minimum public API (names may vary slightly; behavior locked)

```text
// Index one path from bytes (unit of incremental update + golden tests):
IndexFile(ctx, st *store.Store, path string, content []byte, opts IndexOptions) error

// Optional convenience (preferred if cheap): load via VCS then IndexFile
IndexFileAtRev(ctx, st *store.Store, repo vcs.Repository, rev, path string, opts IndexOptions) error

IndexOptions at minimum:
  GitOID *string   // passed through to UpsertFile when known
  // Force / DryRun not required

// Language detection helper (export or internal):
DetectLanguage(path string) (lang string, ok bool)
```

Flow for `IndexFile` (locked):

1. Detect language by extension; reject/skip unsupported / binary.
2. `contentHash := sha256hex(content)`.
3. `UpsertFile(path, contentHash, opts.GitOID)`.
4. Set file language string (`javascript` | `typescript` | `tsx` | `python`).
5. Parse with tree-sitter; extract **minimal** symbols + imports (queries or small walkers).
6. `ReplaceFileSymbols(path, symbols)` then `ReplaceFileImports(path, imports)` (order free; both required even if empty slices — empty means clear that file’s edges).

### Minimal symbol / import semantics

**Symbols** (store `Symbol.Kind` vocabulary — use these strings):

| Kind | Extract at least |
|------|------------------|
| `function` | Top-level functions (JS/TS function decl / exported function; Python `FunctionDefinition` at module level) |
| `class` | Class declarations (JS/TS / Python) |
| `method` | Methods inside classes (both languages) when cheap via query |

Optional if queries stay small: TS `interface`, `type`. Do **not** require call-graph edges, references, or overload sets.

**Lines:** `StartLine` / `EndLine` are **1-based** inclusive line numbers from the tree-sitter node range (convert from 0-based rows).

**Imports:**

| Field | Rule |
|-------|------|
| `ImportedPath` | Module specifier as written (`"./foo"`, `"react"`, `"os"`, `"pkg.sub"`) |
| `Symbol` | Optional named binding (`from x import y` → path `x`, symbol `y`; JS named import when obvious). Default import / side-effect import may leave `Symbol` nil |

Cover at least: ES `import … from '…'`; Python `import …` / `from … import …`. `require()` optional.

### Target tree

```text
internal/analyzers/
  doc.go                 # package contract: tree-sitter; file-local incremental; langs
  detect.go              # extension → language
  index.go               # IndexFile / IndexFileAtRev orchestration + hash
  extract_js.go          # JS (+ shared with TS where sensible)
  extract_ts.go          # TS/TSX
  extract_py.go          # Python
  queries_*.scm or inline query strings
  testdata/              # tiny .js/.ts/.py fixtures for golden edges
  analyzers_test.go      # golden + incremental isolation

internal/store/          # additive only if needed
  file_graph.go          # SetFileLanguage (or equivalent) — no new schema migration
```

Always `Close()` tree-sitter `Parser` / `Tree` / `Query` / `QueryCursor` (CGO finalizer caveat per upstream docs).

### Out of scope (this row)

- Perfect / whole-program call graphs
- Go, Rust, or other language grammars
- MCP, daemon, HTTP
- `cmd/trace index` (S07)
- Full `fixtures/x0` corpus (S08) — use `internal/analyzers/testdata` + temp dirs here
- Cascading reindex of dependent files

## Board rights
Implementer: update **status + notes only** on `P00-S04-01`. Do not spawn rows or rewrite later prompts.

## Exit criteria
- [ ] `IndexFile` (and preferably `IndexFileAtRev`) exists under `internal/analyzers` and writes via store Upsert + ReplaceSymbols + ReplaceImports only
- [ ] Golden tests: ≥1 JS or TS fixture **and** ≥1 Python fixture with asserted symbol names/kinds and import paths (stable expected lists)
- [ ] **Incremental isolation test:** index files A and B → mutate/reindex only A → B’s symbol/import rows unchanged (IDs or content equality); A’s edges match new parse; **no** API that rebuilds the whole DB graph as the default update path
- [ ] Unsupported / binary paths handled without corrupting other files’ rows
- [ ] Analyzers do not import `gitcli`; VCS usage is via `vcs.Repository` / `vcs.Fake`
- [ ] `go test ./internal/analyzers/...` passes (CGO enabled as required by bindings)
- [ ] `CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/...` still passes
- [ ] `go test ./...` passes with default toolchain CGO (typical Linux `CGO_ENABLED=1`)
- [ ] No MCP/daemon/HTTP; no new product CLI commands; no source BLOBs in SQLite
- [ ] TODO.md Notes for `P00-S04-01` updated; status `done`

## Minimal todos
- [ ] Add tree-sitter binding + JS/TS/Python grammar modules to `go.mod`
- [ ] `DetectLanguage` + `IndexFile` orchestration (hash → Upsert → SetLanguage → extract → Replace*)
- [ ] TS/JS extractors (symbols + imports) with golden testdata
- [ ] Python extractor with golden testdata
- [ ] Optional `IndexFileAtRev` using `vcs.Fake` / ShowFile
- [ ] Incremental isolation test (A/B paths)
- [ ] Store `SetFileLanguage` (or equivalent) if missing
- [ ] Board status + notes (deps versions, CGO note, commands run)
