# Index language policy

Honest language coverage and freshness paths for agents using Trace's file-local indexer.

## Purpose

Trace indexes **symbol-level** structure only for Tier-1 languages. Agents should consult this policy (and `trace index status`) before assuming a path is searchable at symbol depth. Freshness is primarily driven by git hooks; optional foreground watch is a session helper — not a daemon.

## Tier-1 (symbol extraction)

Frozen at **five** languages for Phase 42+. Each has a compile-time `LanguageAdapter` in `internal/analyzers`.

| ID | Extensions | Adapter |
|----|------------|---------|
| `go` | `.go` | `goAdapter` |
| `javascript` | `.js`, `.jsx`, `.mjs`, `.cjs` | `jsAdapter` |
| `typescript` | `.ts` | `tsAdapter` |
| `tsx` | `.tsx` | `tsxAdapter` |
| `python` | `.py` | `pythonAdapter` |

`trace index status` exposes these IDs in `supported_languages`.

## Tier-2 (deferred)

No symbol extraction until a human-promoted board row and [ANALYZER_CONTRIBUTION.md](ANALYZER_CONTRIBUTION.md) checklist per language:

`rust`, `java`, `kotlin`, `ruby`, `c`, `cpp`, `csharp`, `swift`, `php`, `lua`, `shell`

Indexing a Tier-2 path (e.g. `.rs`) returns a skip/error that cites this policy — no silent partial index.

## Tier-3 (path-only)

These extensions are indexed in the **files** table / FTS path channel only when explicitly walked — **no** symbol extraction:

`.md`, `.json`, `.yaml`, `.yml`, `.toml`

Do not expect `ListSymbolsByPath` results for Tier-3 paths unless a future phase adds a dedicated path indexer.

## Freshness — primary (git hook)

Recommended default for teams:

```bash
trace install git-hook [--write]
```

- Installs post-commit and pre-push fragments that run incremental `IndexFile` on changed paths.
- Sets `hook_installed: true` in graph sync state (`trace index status`).
- Does **not** wrap or replace `git commit`.

## Freshness — optional (foreground watch)

```bash
trace index watch [--debounce 300ms] [paths…]
```

- **Foreground only** — process lifetime equals the user session; **SIGINT** exits cleanly (exit 0).
- Watches directories under the given paths (default: project root); debounces per-file events (default **300ms**).
- On fire: read working-tree bytes → `analyzers.IndexFile` (same as `trace index <path>`).
- **Does not** update `last_indexed_commit` / graph sync watermark per file — file-local only.
- Skips T0 paths, binary content, and unsupported extensions (quietly or with a stderr note).
- **Not** a replacement for git-hook freshness; no always-on daemon or detached watcher.

## Manual index

```bash
trace index [paths…]
```

- File-local incremental update; full-tree walk when no paths given.
- Updates graph sync watermark when run against a git repo (full or partial index with repo open).
- **Stale** when `HEAD ≠ last_indexed_commit` — see `trace index status` and task-loop `IndexHonesty` / `GraphSyncHonesty` in context packets.

## Contribution gate

New Tier-1 languages require:

1. Human-promoted board row (no dynamic plugin loader).
2. [ANALYZER_CONTRIBUTION.md](ANALYZER_CONTRIBUTION.md) checklist — tree-sitter adapter, golden test, append to `builtinAdapters`.
3. Update this policy's Tier-1 table and `supported_languages` honesty.

There is **no** public `Register` API and **no** `.so` / hosted marketplace loader.
