# P42-S01-02 — Review (G7 index freshness & langs)

## Metadata
- id: P42-S01-02
- todo_ids: [P42-S01-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter]
- mcps: [user-trace, user-codegraph]
- verification: mixed

## Objective

Fresh independent review of S01-01 G7 implementation vs REMEDIATION-PLAN G7, GAP-REGISTRY G-005, M-001 moat, Laws 6–7/19, and P42-00 lang/watch locks.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md) — G7-F1–F6 acceptance map + touch-list
- [00-PLANNER.md](00-PLANNER.md) — locks + tier policy
- [docs/INDEX_LANG_POLICY.md](../../../../INDEX_LANG_POLICY.md) — produced by S01-01
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [REMEDIATION-PLAN G7](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-005](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- Evidence: [h5-index-watch-contrast.md](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/h5-index-watch-contrast.md)

## Session start

Follow agent-loop-protocol Session start. **Fresh subagent** — do not share implementer session.

## Locked defaults

| Item | Value |
|------|-------|
| APPROVE bar | Medium+ confidence; zero open blocker/high |
| Spawn trigger | Blocker/high → spawn 02a/02b below this row |
| Tier-1 | Exactly 5 langs — no surprise Tier-2 adapter in `builtinAdapters` |
| Daemon reject | Watch must be foreground — no detached process, no `nohup`/fork pattern |
| Git-hook | Primary freshness unchanged — watch does not replace hook fragment |

## Preflight evidence commands

Run fresh subagent — do not trust implementer Notes alone.

```bash
# Policy doc shipped
test -f docs/INDEX_LANG_POLICY.md
grep -q 'Tier-1\|Tier-2\|Tier-3' docs/INDEX_LANG_POLICY.md
grep -q 'git-hook\|install git-hook' docs/INDEX_LANG_POLICY.md

# SupportedLanguages exported
grep -n 'func SupportedLanguages' internal/analyzers/*.go
grep -n 'supported_languages' cmd/trace/index_status.go

# Watch CLI exists + foreground (no daemon keywords in watch path)
test -f cmd/trace/index_watch.go
grep -n 'watch' cmd/trace/index.go cmd/trace/help.go
rg -n 'daemon|Detach|background service|nohup' cmd/trace/index_watch.go || true

# Error UX improved
grep -n 'unsupported' internal/analyzers/errors.go

# Git-hook primary unchanged
grep -n 'trace index' internal/install/githook.go
grep -n 'setHookInstalledFlag' cmd/trace/install.go

# Tests green
CGO_ENABLED=1 go test ./internal/analyzers/... ./cmd/trace/... -count=1 \
  -run 'TestIndexStatusSupportedLanguages|TestSupportedLanguages|TestIndexUnsupported|TestIndexWatch'

# M-001 regression
CGO_ENABLED=1 go test ./internal/compiler/... -count=1 -run 'IndexHonesty|GraphSync'
```

## Review checklist

### A — G7 gap closure

- [ ] `INDEX_LANG_POLICY.md` tier table matches shipped adapters (`language_adapter.go:19–25`)
- [ ] Tier-2 list documented as defer; Tier-3 path-only (.md/.json/.yaml/.yml/.toml)
- [ ] `SupportedLanguages()` exported; stable sorted order
- [ ] `supported_languages` on index status JSON (G7-F1)
- [ ] `trace index watch [--debounce …] [paths…]` foreground + debounce (G7-F4–F5)
- [ ] Unsupported ext honest error cites policy (G7-F3)
- [ ] `ANALYZER_CONTRIBUTION.md` cross-refs INDEX_LANG_POLICY tier-2 gate
- [ ] G7-F1–F6 evidenced green (F6 N/A documented if HTTP deferred)

### B — Freshness paths

- [ ] Git-hook path unchanged (`githook.go:102–113`) and documented as primary in policy
- [ ] Watch does not wrap `git commit` or block VCS
- [ ] File-local incremental index preserved — watch calls `IndexFile`, not full-tree rebuild
- [ ] Watch does not auto-update graph sync watermark on every file (unless explicitly documented otherwise)

### C — M-001 / honesty

- [ ] IndexHonesty / GraphSyncHonesty regression green (`packet.go:58–74`)
- [ ] Status surfaces stale/manual index honestly (`stale`, `hook_installed`, `last_indexed_commit`)
- [ ] `supported_languages` is informational — does not imply unindexed langs are searchable

### D — Rejects

- [ ] No always-on daemon / background service / detached watcher
- [ ] No Tier-2 language adapter added without board spawn
- [ ] No CG detached watcher port (study debounce only — fsnotify OK)
- [ ] No full-rebuild-on-change architecture
- [ ] No new MCP tool solely for lang policy

### E — Law 19

- [ ] Watch/status logic in `cmd/` + `internal/analyzers/` — HTTP mirror thin if present
- [ ] No business logic fork in `web/`

### F — Dependency hygiene

- [ ] `fsnotify` added only for watch — no unrelated deps
- [ ] CGO still required for analyzers only; watch/store paths CGO-free OK

## Spawn criteria (02a/02b)

| Severity | Example | Action |
|----------|---------|--------|
| blocker | Watch forks/detaches; Tier-2 adapter snuck in | Spawn 02a implement + 02b review |
| high | Missing policy doc; status JSON wrong lang count; full rebuild on watch | Spawn or inline fix if trivial |
| medium | HTTP mirror missing when promised; debounce default wrong | Prefer spawn unless trivial |
| low/nit | Help wording; sort order nits | Note only |

## Exit criteria

- APPROVE with medium+ confidence; zero open blocker/high without pending spawn
- Board row → `done` with verdict + confidence in Notes

## Next

`P42-S02-00`
