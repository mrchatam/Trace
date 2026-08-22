# P00-S04-02 — Scope review notes (2026-08-15)

Independent review of S04 against `01-analyzers.md` + TODO Notes for `P00-S04-01`. Fresh session; claims verified in-repo.

## Plan (executed)

1. Diff claims vs `internal/analyzers` + store surface + `go.mod`
2. Re-run golden / A/B / CGO gate tests
3. Severity-tag findings; fix/spawn only for blocker/high
4. Write these notes; mark board + SCOPE-TODOS; light thicken upcoming if needed

## Verdict

**APPROVE** — no blocker/high. Confidence: **high**. Spawns: **none**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| `IndexFile` / `IndexFileAtRev` / `DetectLanguage` under `internal/analyzers` | Pass (`index.go`, `detect.go`) |
| Official `github.com/tree-sitter/go-tree-sitter` v0.25.0 + JS/TS/Python grammar modules | Pass (`go.mod`: js/py v0.25.0, ts v0.23.2) |
| Languages TS/JS/Python only (DR-ANLANG); extensions match lock table | Pass (`detect.go`) |
| Persist only via `UpsertFile` + `SetFileLanguage` + `ReplaceFileSymbols` + `ReplaceFileImports` | Pass (`index.go`; no raw SQL / parallel schema in analyzers) |
| SHA-256 hex of bytes; no source blobs | Pass (`sha256Hex`; store `files` has `content_hash` only) |
| File-local incremental; no full-rebuild default API | Pass (`doc.go` + `TestIncrementalIsolation`; no IndexAll/Rebuild) |
| No `extract_js.go` GOOS=js trap | Pass (`extract_javascript.go` + comment; file absent) |
| CGO confined: analyzers need CGO; store/vcs/gitcli `CGO_ENABLED=0`-clean | Pass (re-run below) |
| Golden JS + TS + Python; A/B isolation; SkipError for unsupported/binary | Pass (`analyzers_test.go` + `testdata/`) |
| `Close()` on Parser / Tree / Query / QueryCursor | Pass (`extract.go`, extract_* defer Close) |
| No `gitcli` import; VCS via `vcs.Repository` / `vcs.Fake` | Pass (imports + `TestIndexFileAtRevUsesVCS`) |
| No MCP/daemon/HTTP; no new CLI index commands | Pass (`cmd/trace` unchanged) |
| Cross-scope: S06/S07/S08 can consume IndexFile + ListSymbols/Imports | Pass (upcoming prompts already name the surface) |

## Re-verification commands (2026-08-15)

```text
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1   # ok
CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/... -count=1  # ok
CGO_ENABLED=1 go test ./... -count=1                     # ok
CGO_ENABLED=0 go test ./internal/analyzers/...           # build failed (expected; grammar CGO)
```

## Findings

### Blocker / high

None.

### Medium (residual — no spawn)

- **Non-atomic IndexFile:** locked flow upserts + sets language before extract/`Replace*`. If extract fails after upsert, `content_hash`/`language` can advance while prior symbols/imports remain. Rare for well-formed text; no short-circuit-on-hash yet. Acceptable for P0; transactional wrap or extract-first would be a later harden.
- **JS/TS symbol query is not strictly module-top-level:** nested `function_declaration` nodes are captured (Python correctly scopes via `(module …)`). Minimal-graph over-extraction vs lock wording “top-level”; not wrong for P0-X.

### Low / nit

- `TestAnalyzersDoNotImportGitcli` is a compile-time `vcs.Fake` assert only — does not grep import graph; package imports are already clean.
- `IndexFile` ignores `ctx` (`_ = ctx`) — fine until cancellation is wired.
- TS golden asserts symbols only (imports covered by JS + Python goldens) — meets “≥1 JS or TS **and** ≥1 Python” exit criterion.

## Spawns

None.

## Residual risks

- Arrow/`const` function expressions and `require()` not extracted (explicitly optional / out of minimal set).
- CGO required to *build* analyzers; CI must enable CGO for packages that import `internal/analyzers` (S07/S08). Store/vcs/gitcli remain pure Go.
- Grammar version skew (TS module at v0.23.2 vs JS/Python v0.25.0) is upstream-normal; monitor query breakage on bumps.

## Forward edits this review

- `SCOPE-TODOS.md` — mark S04-01/02 done
- `docs/TODO.md` — `P00-S04-02` → done + notes
- `scopes/scope-07-cli/01-cli.md` — light API note: `IndexFile` / `IndexFileAtRev` / `SkipError`; CGO when linking analyzers
