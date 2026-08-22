# P08 / S01 / 01 — Plugin APIs / analyzer contribution surface

## Metadata
- id: P08-S01-01
- todo_ids: [P08-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement the **versioned analyzer contribution surface** locked by `00-PLANNER.md`: `LanguageAdapterAPIVersion`, `LanguageAdapter` iface, compile-time builtin adapter table, DetectLanguage/extract dispatch via adapters, contributor doc, and contribution-path tests. Keep IndexFile orchestration and file-local incremental semantics. Do **not** invent a megastore, dynamic `.so` loader, or weaken Gate H / Gate C.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks finalized 2026-08-16
- [phase README](../../README.md)
- Live: `internal/analyzers/{detect,extract,index,extract_*.go,doc}.go`
- Prior pattern: P07-S02 Go adapter (switch extension) — this row **stabilizes** that pattern behind a versioned iface + static table
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute.

## Locked defaults (FINAL — P08-S01-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Package | **`internal/analyzers` only** (product Go) |
| API version | Exported const **`LanguageAdapterAPIVersion = 1`** |
| Iface | **`LanguageAdapter`**: `ID() string`; `Extensions() []string` (lowercase, leading `.`); `Extract(content []byte) ([]store.Symbol, []store.Import, error)` |
| Registration | Compile-time **static** builtin adapter table — **no** public mutable Register; **no** `.so` / plugin megastore |
| DetectLanguage | Exported signature unchanged; resolve via adapter `Extensions()` → `ID()` |
| extract | Unexported dispatch via adapter `ID()` → `Extract` |
| Existing langs | javascript / typescript / tsx / python / go as adapters; preserve `Lang*` ids + golden behavior |
| IndexFile / IndexFileAtRev | Public signatures + Upsert→SetFileLanguage→extract→Replace* flow **unchanged** |
| Contributor doc | Create **`docs/ANALYZER_CONTRIBUTION.md`**; refresh **`doc.go`** (include Go; link contribution doc) |
| Migration | **No** `011_*` |
| CGO | Analyzers-only |
| CLI / MCP | No new commands/tools |
| Proof | Goldens stay green; `TestLanguageAdapterAPIVersion`; `TestBuiltinLanguageAdaptersContributionPath` |
| Carry-forward | Gate H + honesty A/B/C + Gate G/E/F + ablation + p0x + x0 + Gate C `dry_run:false` |
| Forbidden | Dynamic loader; megastore; full-rebuild; daemon/HTTP/embeddings primary; Gate H threshold invent; Gate C pack rewrite |

## Extension points (exact files)

| File | Work |
|------|------|
| Prefer new `adapter.go` (or `language_adapter.go`) | `LanguageAdapterAPIVersion`; `LanguageAdapter` iface; static `builtinAdapters` table; helpers `adapterByExt` / `adapterByID` if useful |
| `detect.go` | Keep `Lang*` consts; implement `DetectLanguage` via adapter table (may thin-wrap helpers from adapter.go) |
| `index.go` | Keep `IndexFile` / `IndexFileAtRev`; change `extract` switch to adapter lookup |
| `extract_javascript.go` / `extract_ts.go` / `extract_py.go` / `extract_go.go` | Keep extractors; wrap as `LanguageAdapter` (small structs or funcs-as-adapters). Avoid GOOS traps (`*_js.go`) |
| `doc.go` | Refresh supported-lang list; point to `docs/ANALYZER_CONTRIBUTION.md` |
| `docs/ANALYZER_CONTRIBUTION.md` | **New** — version, iface contract, add-language steps, forbidden patterns |
| `analyzers_test.go` (+ testdata as needed) | Version + contribution-path tests; keep existing goldens |

## Role work

1. TDD: add failing `TestLanguageAdapterAPIVersion` + `TestBuiltinLanguageAdaptersContributionPath`.
2. Introduce iface + version + static table; migrate five languages; wire DetectLanguage + extract.
3. Write `docs/ANALYZER_CONTRIBUTION.md`; refresh `doc.go`.
4. Confirm IndexFile isolation / goldens / DetectLanguage still pass.
5. Run locked verify suite; board **status + Notes only**.

### Contribution-path test requirements

`TestBuiltinLanguageAdaptersContributionPath` must assert at least:

- Builtin table length ≥ 5 and covers `LangJavaScript`, `LangTypeScript`, `LangTSX`, `LangPython`, `LangGo`.
- For each adapter: every `Extensions()` entry → `DetectLanguage("x"+ext)` returns `(adapter.ID(), true)`.
- For each language id: `extract(id, []byte("\n"))` does not return “unsupported language” (empty/minimal content may yield empty symbols — OK; hard fail only on missing adapter).
- No duplicate extension claims across adapters.

`TestLanguageAdapterAPIVersion`: `LanguageAdapterAPIVersion == 1`.

## Verify commands (locked)

```bash
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... -count=1
CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/... ./internal/domain/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./evals/perf/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Optional: confirm Gate C artifacts under `docs/verification/gate-c-x0/` remain `dry_run:false` N=3 (do **not** rewrite packs). Spot-check Gate H still green via `./evals/perf/...` above.

## Exit criteria
- [ ] `LanguageAdapterAPIVersion == 1` exported; `LanguageAdapter` iface present
- [ ] Builtin static adapter table drives DetectLanguage + extract; five langs preserved
- [ ] `docs/ANALYZER_CONTRIBUTION.md` exists with versioned contribution steps
- [ ] `doc.go` refreshed; IndexFile orchestration unchanged; no mig; no `.so`/megastore
- [ ] Contribution-path + version tests + goldens green
- [ ] Carry-forward suite green (incl. Gate H / Gate C intact)
- [ ] Board Notes ready for **P08-S01-02**

## Out of scope
- Worktrees (S02); backup/auth (S03); `evals/compat` (S04)
- Adding a sixth natural language beyond wrapping the existing five
- Dynamic plugin loader / public Register API / hosted marketplace
- Daemon / HTTP / embeddings / new MCP tools / CLI subcommands
- Store schema migrations
- Inventing Gate H thresholds or rewriting Mode-B Gate C packs

## Todo updates
Implementer: own row status + Notes only. Do not rewrite planner locks or spawn board rows.

## Minimal todos
- [ ] Add version + contribution-path tests (fail first)
- [ ] Implement LanguageAdapter + static table; wire DetectLanguage/extract
- [ ] Migrate JS/TS/TSX/Python/Go extractors to adapters; keep goldens
- [ ] Write `docs/ANALYZER_CONTRIBUTION.md`; refresh `doc.go`
- [ ] Run locked verify suite; mark P08-S01-01 done with Notes
