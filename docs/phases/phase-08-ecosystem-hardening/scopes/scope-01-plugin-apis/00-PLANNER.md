# P08 / S01 / 00-PLANNER — Plugin APIs

## Metadata
- id: P08-S01-00
- todo_ids: [P08-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling implement/review prompts for **stable plugin / analyzer contribution APIs** after `P08-00` locks phase order. Prefer versioned adapter-shaped surface over premature megastore. No product code in this planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 8
- Live: `internal/analyzers` `DetectLanguage` + `extract` + `extract_{javascript,ts,py,go}.go`
- Phase 07 S02 Go pattern: adapter switch extension (no plugin registry)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Live inventory (locked 2026-08-16 — P08-S01-00)

| Item | Today | S01 lock |
|------|-------|----------|
| Extension points | `detect.go` switch (`DetectLanguage`); `index.go` `extract(lang)` switch; per-lang `extract_*.go`; `IndexFile` orchestration Upsert→SetFileLanguage→extract→Replace* | **Versioned** `LanguageAdapter` iface + compile-time adapter table; `IndexFile` / `IndexFileAtRev` orchestration **unchanged** |
| Languages | JS / TS / TSX / Python / Go (`Lang*` consts) | Same five as first-class `LanguageAdapter` implementations; contributor path = add N+1 adapter + table entry (no IndexFile rewrite) |
| Registry | Absent (dual switches only) | **Static** `[]LanguageAdapter` (or equivalent compile-time table) — **not** dynamic `.so` / plugin megastore / runtime Register mutator for product |
| Version | Absent | `LanguageAdapterAPIVersion = 1` const in `internal/analyzers` |
| Docs | `doc.go` outdated (omits Go) | Package godoc + **`docs/ANALYZER_CONTRIBUTION.md`** contributor steps |
| Tests | Per-lang goldens + DetectLanguage | Keep goldens; add version/compat + adapter-table contribution-path smoke |
| Migration | Through `010_*` | **No** mig this scope |
| CGO | Analyzers-only | Unchanged |

### Exact live call sites (implementer must touch)

1. **`detect.go`** — `Lang*` consts + `DetectLanguage` extension switch.
2. **`index.go`** — `IndexFile` calls `DetectLanguage` then `extract(lang, content)`; `extract` language switch → `extractJS` / `extractTS` / `extractPython` / `extractGo`.
3. **`extract_{javascript,ts,py,go}.go`** — private extractors + tree-sitter queries.
4. **`doc.go`** — package comment (still lists JS/TS/Python only — refresh).
5. **`analyzers_test.go`** + `testdata/` — goldens; DetectLanguage cases.

`cmd/trace` walk already uses `DetectLanguage` — no CLI change required when adapters drive detection.

## Phase defaults already locked (respect — P08-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Prefer | Adapter-shaped language plugins (P07-S02 Go pattern); versioned contribution surface |
| Forbidden | Premature universal plugin megastore; weakening Gate H / Gate C |
| Carry-forward | Gate H + honesty + Gates E/F/G + ablation + p0x + x0 + Gate C `dry_run:false` |
| Later scopes | S02 worktrees; S03 production; S04 VERIFY owns `evals/compat` |

## Planner work
- [x] Inventory live DetectLanguage/extract extension points vs contributor needs
- [x] Lock API shape (version const + `LanguageAdapter` iface + docs) without megastore
- [x] Thicken `01-plugin-apis.md` + `02-scope-review.md`
- [x] Light S02/S03 Depends notes; SCOPE-TODOS sync; board Notes

## Locked defaults (FINAL — 2026-08-16 — do not re-debate in S01-01)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Package | **`internal/analyzers` only** for product changes (package `analyzers`) |
| API version | **`LanguageAdapterAPIVersion = 1`** (exported const). Bump only on breaking iface/contract change; document in `docs/ANALYZER_CONTRIBUTION.md` |
| Iface | **`LanguageAdapter`** with methods: `ID() string`; `Extensions() []string` (lowercase, leading-dot, e.g. `".go"`); `Extract(content []byte) ([]store.Symbol, []store.Import, error)` |
| Registration | Compile-time **static table** of built-in adapters (e.g. `var builtinAdapters []LanguageAdapter` or package-level slice initialized in one file). Product code **must not** expose dynamic `.so` loading or a public mutable global Register API |
| DetectLanguage | Must resolve via the adapter table (extensions → `ID()`). Keep exported `DetectLanguage(path) (lang string, ok bool)` signature. Unsupported ext → `ok=false` |
| extract | Unexported `extract(lang, content)` must dispatch via adapter `ID()` → `Extract`. Keep `IndexFile` body orchestration stable |
| Existing langs | Wrap/migrate current extractors into adapters: javascript, typescript, tsx, python, go — preserve `Lang*` string ids and golden behavior |
| New language (contrib recipe) | (1) `extract_<lang>.go` + official tree-sitter grammar; (2) `LanguageAdapter` impl; (3) append to static table; (4) golden test; (5) update contribution doc — **no** IndexFile rewrite |
| Contributor doc | Create **`docs/ANALYZER_CONTRIBUTION.md`** (version, iface, steps, forbidden patterns). Refresh `doc.go` to list Go + point at that doc |
| Migration | **No** `011_*` — analyzer surface only; version is a code const, not DB metadata |
| IndexFile / IndexFileAtRev | Unchanged public signatures and file-local incremental semantics (DR-INCREMENTAL) |
| CGO | Analyzers-only; store/domain/vcs/gitcli remain `CGO_ENABLED=0`-clean |
| CLI / MCP | No new subcommands or MCP tools this scope |
| Proof | Keep JS/TS/Py/Go goldens + DetectLanguage; add **`TestLanguageAdapterAPIVersion`** (const == 1) + **`TestBuiltinLanguageAdaptersContributionPath`** (every builtin adapter: Extensions→DetectLanguage→ID; Extract wired; table non-empty). Do **not** invent Gate H thresholds or rewrite Gate C packs |
| Carry-forward | Gate H + honesty A/B/C + Gate G/E/F + ablation + p0x + x0 + Gate C `dry_run:false` |
| Forbidden | Dynamic plugin loader / `.so` theater; universal plugin megastore; full-rebuild-on-any-change; daemon/HTTP/embeddings as primary; VerifiedFact / `plan simulate` |

## Cross-scope notes (light)

| Scope | Note from S01-00 |
|-------|------------------|
| S02 worktrees | Independent of `LanguageAdapter`. Keep `-C`/`DetectLanguage` call sites. Do not couple worktree bind to analyzer registration |
| S03 production | S01 adds **no** store mig. S03 may still introduce `011_*` for backup/auth. Do not persist `LanguageAdapterAPIVersion` in SQLite this phase |
| S04 VERIFY | Owns `evals/compat` checklist finalization; may cite `LanguageAdapterAPIVersion` + contribution doc as checklist items — S04-00 finalizes |

## Exit criteria
- [x] `01-*.md` runnable alone with locked paths
- [x] No product Go; board Notes

## Minimal todos
- [x] Inventory live DetectLanguage/extract extension points
- [x] Lock API shape + exit criteria; thicken 01/02
- [x] Light path notes for S02/S03; SCOPE-TODOS sync

## Out of scope
- Implementing adapters (S01-01); independent review (S01-02)
- Worktrees (S02); backup/auth (S03); `evals/compat` harness (S04)
- Second new natural language beyond table migration of existing five
- Dynamic loaders / hosted plugin marketplace
