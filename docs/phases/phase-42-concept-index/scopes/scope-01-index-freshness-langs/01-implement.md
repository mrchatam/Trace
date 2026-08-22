# P42-S01-01 — Implement (G7 index freshness & langs)

## Metadata
- id: P42-S01-01
- todo_ids: [P42-S01-01]
- role: implementer
- skills: [backend-dev, test-driven-development]
- mcps: [user-trace, user-codegraph]
- verification: mixed

## Objective

Implement **G7**: index freshness ergonomics and language coverage policy (G-005). Local-first; no always-on daemon default.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md) — Laws 6–7, 19
- [00-PLANNER.md](00-PLANNER.md) — **SoT** for locks
- [REMEDIATION-PLAN G7](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-005](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [ANALYZER_CONTRIBUTION.md](../../../../ANALYZER_CONTRIBUTION.md)
- Live anchors (P42-S01-00 re-verified 2026-08-22):
  - `internal/analyzers/detect.go:9–15` — 5 lang ID consts
  - `internal/analyzers/language_adapter.go:19–25` — `builtinAdapters` static table (5 adapters)
  - `internal/analyzers/doc.go:3–4` — DR-ANLANG supported langs doc comment
  - `internal/analyzers/index.go:33–36` — `DetectLanguage` fail → `unsupportedExt`
  - `internal/analyzers/errors.go:16–18` — bare `"unsupported extension"` today (G7-F3 gap)
  - `cmd/trace/index_status.go:14–19` — 4 JSON fields; no `supported_languages`
  - `cmd/trace/index.go:19–22` — `index status` dispatch only; no `watch` branch
  - `cmd/trace/help.go:27–29` — index help omits watch + policy
  - `internal/install/githook.go:102–113` — post-commit/pre-push incremental index (primary freshness)
  - `cmd/trace/install.go:258–272` — `setHookInstalledFlagForProject` on git-hook install
  - `internal/compiler/packet.go:58–74` — IndexHonesty + GraphSyncHonesty (M-001; do not regress)
  - `internal/httpapi/handlers_p1.go:338–364` — GET /v1/index status mirror (4 fields)
  - Evidence: [h5-index-langs.txt](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h5-index-langs.txt) · [h5-index-watch-contrast.md](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/h5-index-watch-contrast.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Preflight

- Re-read `00-PLANNER.md` locked defaults before coding.
- Confirm gaps still match live repo (no surprise Tier-2 adapter landed):
  ```bash
  rg -n 'builtinAdapters|SupportedLanguages|index watch|INDEX_LANG_POLICY' \
    internal/analyzers cmd/trace docs/INDEX_LANG_POLICY.md 2>/dev/null || true
  test ! -f cmd/trace/index_watch.go  # watch not shipped yet
  ```
- Git-hook remains primary — watch is opt-in foreground helper only.

## Locked defaults

| Item | Value |
|------|-------|
| GAP ids | G-005 |
| Verdict | **Accept — ship** (P42-00 lock) |
| Lang policy | Tier-1 frozen at 5 langs; Tier-2 defer; Tier-3 path-only documented |
| Watch | Foreground `trace index watch` only — debounce default 300ms; SIGINT exits |
| Watch dep | Add `github.com/fsnotify/fsnotify` — no CG detached watcher stack |
| Git-hook | Primary freshness path — do not remove, wrap, or replace git commit |
| Status JSON | Add `supported_languages: ["go","javascript","typescript","tsx","python"]` |
| HTTP | Optional mirror on GET /v1/index — not blocking (Note if deferred) |
| Graph export | Policy doc only — no entity schema change expected |
| M-001 | IndexHonesty/GraphSyncHonesty in packets unchanged; status JSON complements loop |

## Touch-list (library → CLI → docs → tests)

| Step | File | Action |
|------|------|--------|
| 1 | `docs/INDEX_LANG_POLICY.md` | **New** — tier table, contribution gate, watch/hook guidance (see outline below) |
| 2 | `internal/analyzers/language_adapter.go` | Export `SupportedLanguages() []string` — sorted IDs from `builtinAdapters` |
| 3 | `internal/analyzers/detect_test.go` or `language_adapter_test.go` | G7-F2 unit test |
| 4 | `cmd/trace/index_status.go` | Add `SupportedLanguages []string \`json:"supported_languages"\``; populate from analyzers |
| 5 | `cmd/trace/index_watch.go` | **New** — foreground watch + debounced `IndexFile` per changed path |
| 6 | `cmd/trace/index.go` | Dispatch `index watch` before generic index paths (`args[0] == "watch"`) |
| 7 | `cmd/trace/help.go` | Document `index watch [--debounce 300ms] [paths…]` + policy pointer |
| 8 | `internal/analyzers/errors.go` | Richer `unsupportedExt` message → cite INDEX_LANG_POLICY + tier |
| 9 | `docs/ANALYZER_CONTRIBUTION.md` | Cross-ref INDEX_LANG_POLICY tier-2 gate; keep v1 table in sync |
| 10 | `cmd/trace/index_status_test.go` | **New** or extend — G7-F1 |
| 11 | `cmd/trace/index_watch_test.go` | **New** — G7-F4–F5 |
| 12 | `internal/httpapi/handlers_p1.go` | Optional: mirror `supported_languages` on GET /v1/index |
| 13 | `go.mod` / `go.sum` | `fsnotify` dependency for watch |

**Explicit non-touch:**

- New Tier-2 language adapters (rust/java/…)
- Always-on daemon / background service / detached process
- Full-rebuild indexer architecture
- `web/` GUI index controls
- Replacing git-hook with watch-only workflow
- MCP new tools (status via existing CLI/HTTP paths)

## INDEX_LANG_POLICY.md outline (required sections)

1. **Purpose** — honest lang coverage + freshness paths for agents
2. **Tier-1 (symbol extraction)** — table: id, extensions, adapter file; frozen at 5 langs for P42
3. **Tier-2 (deferred)** — rust, java, kotlin, ruby, c, cpp, csharp, swift, php, lua, shell — human-promoted board row + ANALYZER_CONTRIBUTION per lang
4. **Tier-3 (path-only)** — `.md`, `.json`, `.yaml`, `.yml`, `.toml` — files table / FTS only; no symbol extraction
5. **Freshness primary** — `trace install git-hook [--write]`; post-commit incremental index; `hook_installed` flag
6. **Freshness optional** — `trace index watch [--debounce 300ms] [paths…]` foreground; SIGINT exits; not a daemon
7. **Manual index** — `trace index [paths…]` file-local; stale when HEAD ≠ last_indexed_commit
8. **Contribution gate** — pointer to ANALYZER_CONTRIBUTION.md; no dynamic plugin loader

## Watch CLI sketch (implementer guidance)

```text
trace index watch [--debounce 300ms] [paths…]
  - Default paths: project root (respect existing walkIndexable / T0 skip rules)
  - fsnotify watcher on paths; coalesce events per path with debounce timer
  - On fire: read file → analyzers.IndexFile (same as cmdIndex indexOne)
  - Stderr: "watching …" on start; "indexed <path>" per success; skip binary/unsupported quietly or log
  - SIGINT / context cancel → clean exit 0 (G7-F5)
  - Do NOT update graph sync watermark on every file (file-local only); document in policy
  - Do NOT fork/detach; process lifetime = user session
```

## Implementation order

```text
1. SupportedLanguages() + INDEX_LANG_POLICY.md + ANALYZER_CONTRIBUTION cross-ref
2. index status JSON field + tests G7-F1, G7-F2
3. Unsupported ext error UX G7-F3
4. go get github.com/fsnotify/fsnotify
5. trace index watch foreground CLI + tests G7-F4–F5
6. help.go update
7. Optional HTTP mirror G7-F6
8. CGO_ENABLED=1 go test ./internal/analyzers/... ./cmd/trace/... -count=1
```

## Acceptance tests (must pass)

| ID | Suggested name | Assert |
|----|----------------|--------|
| G7-F1 | `TestIndexStatusSupportedLanguages` | `trace index status` JSON includes `supported_languages` with exactly 5 tier-1 IDs |
| G7-F2 | `TestSupportedLanguagesMatchesAdapters` | `SupportedLanguages()` matches `builtinAdapters` IDs (sorted stable order) |
| G7-F3 | `TestIndexUnsupportedExtMessage` | Indexing `.rs` (or other Tier-2 ext) error/skip mentions policy or "unsupported" + lang tier honesty |
| G7-F4 | `TestIndexWatchDebounced` | Temp dir: write `.go` file → watch indexes after debounce (symbol row or file row in store) |
| G7-F5 | `TestIndexWatchForegroundExit` | Cancel context / SIGINT path exits cleanly (exit 0); no goroutine leak (use `-race` if practical) |
| G7-F6 | `TestHTTPIndexStatusLanguages` | GET /v1/index includes `supported_languages` (if mirrored; else Note + skip) |

### G7-F4 test hints

- Use short debounce (50ms) via flag in test harness
- Init store in temp dir; run watch in goroutine with cancelable context
- Avoid real git dependency unless existing cli_test patterns require it

## Regression tests (must stay green)

```bash
CGO_ENABLED=1 go test ./internal/install/... -count=1 -run GitHook
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestIndex'
CGO_ENABLED=1 go test ./internal/compiler/... -count=1 -run 'IndexHonesty|GraphSync'
```

## Self-check before done

```bash
test -f docs/INDEX_LANG_POLICY.md
grep -q 'SupportedLanguages' internal/analyzers/*.go
grep -q 'supported_languages' cmd/trace/index_status.go
grep -q 'watch' cmd/trace/index.go cmd/trace/help.go
CGO_ENABLED=1 go test ./internal/analyzers/... ./cmd/trace/... -count=1 \
  -run 'TestIndexStatusSupportedLanguages|TestSupportedLanguages|TestIndexUnsupported|TestIndexWatch'
```

## Role work

1. Document lang tiers — ship policy, not sprawl.
2. Git-hook remains primary; watch is opt-in foreground helper.
3. Self-check G7-F1–F6 before marking row done.

## Exit criteria

- [ ] G7-F1–F6 green (F6 N/A ok if HTTP deferred with Note)
- [ ] `docs/INDEX_LANG_POLICY.md` authored per outline
- [ ] Board row → `done` with files + test command in Notes

## Next

`P42-S01-02`
