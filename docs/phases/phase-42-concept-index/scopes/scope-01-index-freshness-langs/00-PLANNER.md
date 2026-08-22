# P42-S01-00 — Scope planner (G7 index freshness & langs)

## Metadata
- id: P42-S01-00
- todo_ids: [P42-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, context-engineering]
- mcps: [user-trace, user-codegraph]
- verification: automated

## Objective

Lock S01 **G7** against live repo: index freshness & language coverage (G-005). Thicken `01-implement.md` + `02-review.md`. **No product code in this row.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md) — P42-00 Q2/Q3 resolutions
- [REMEDIATION-PLAN §2 G7](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-005](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [ANALYZER_CONTRIBUTION.md](../../../../ANALYZER_CONTRIBUTION.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors (verified 2026-08-22 P42-00):
  - `internal/analyzers/detect.go:9–15` — 5 lang IDs: go, javascript, typescript, tsx, python
  - `internal/analyzers/language_adapter.go:19–25` — compile-time static table (5 adapters)
  - `internal/analyzers/doc.go:3` — DR-ANLANG supported langs
  - `internal/analyzers/index.go:22–36` — file-local IndexFile; unsupported ext → error
  - `cmd/trace/index_status.go:14–19` — head, last_indexed_commit, stale, hook_installed
  - `internal/install/githook.go:102–113` — post-commit/pre-push incremental index (exists)
  - `cmd/trace/install.go:258` — sets `hook_installed` on git-hook install
  - `internal/compiler/packet.go:58–74` — IndexHonesty + GraphSyncHonesty (Phase 12+)
  - Evidence: [h5-index-langs.txt](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h5-index-langs.txt) · [h5-index-watch-contrast.md](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/h5-index-watch-contrast.md)

## Session start

Follow agent-loop-protocol Session start. Unattended: INTAKE + P42-00 locks are authority.

## Locked defaults (FINAL — P42-00)

| Item | Value |
|------|-------|
| GAP ids | G-005 |
| Verdict | **Accept** per REMEDIATION-PLAN G7 |
| P42-00 Q2 lang policy | **Tier-1 frozen (5 langs)** — no new analyzer in P42 unless S01-01 blocked on policy-only path |
| P42-00 Q3 watch/hook | **Git-hook primary**; optional **foreground** `trace index watch` — no always-on daemon |
| Tier-1 (shipped) | `go`, `javascript`, `typescript`, `tsx`, `python` |
| Tier-2 (defer) | rust, java, kotlin, ruby, c, cpp, csharp, swift, php, lua, shell — human promotion + ANALYZER_CONTRIBUTION per lang |
| Tier-3 (path-only) | `.md`, `.json`, `.yaml`, `.yml`, `.toml` — path FTS via files table only; no symbol extraction |
| Freshness primary | `trace install git-hook [--write]` — post-commit incremental index (already shipped) |
| Freshness optional | `trace index watch [--debounce 300ms] [paths…]` — foreground process; exits on SIGINT |
| Honesty | Extend `trace index status` JSON with `supported_languages[]`; document stale/manual index |
| M-001 | Index honesty surfaces merge into task loop — no silent stale reads |
| Out | Always-on daemon; CG detached watcher stack; full-rebuild-on-change; Tier-2 langs without board spawn |

## Live repo gap (re-verified P42-S01-00)

| Check | Shipped | Gap |
|-------|---------|-----|
| Lang count | 5 IDs (`detect.go:9–15`, `language_adapter.go:19–25`) | Policy doc + status JSON listing |
| `hook_installed` flag | Set on git-hook install (`install.go:258`) | Document primary freshness path |
| File watcher | **None** (`index_watch.go` absent; no fsnotify dep) | Optional foreground watch CLI |
| Unsupported ext UX | `SkipError` "unsupported extension" (`errors.go:16–18`) | Honest message + policy pointer |
| Index status JSON | 4 fields (`index_status.go:14–19`) | Add `supported_languages` |
| M-001 honesty | IndexHonesty + GraphSyncHonesty (`packet.go:58–74`) | Preserve; status JSON complements |

## Accept / reject (G7)

| Decision | Item |
|----------|------|
| **Accept** | `docs/INDEX_LANG_POLICY.md` — tier table + contribution gate |
| **Accept** | `analyzers.SupportedLanguages()` or equivalent + `trace index status` field |
| **Accept** | `trace index watch` foreground watcher (debounced IndexFile) |
| **Accept** | Improve unsupported-language error with policy link |
| **Accept** | Tests G7-F1–F6 (see thickened `01-implement.md`) |
| **Accept** | Update `ANALYZER_CONTRIBUTION.md` cross-ref to policy |
| **Reject** | Always-on background daemon / detached MCP watcher |
| **Reject** | Tier-2 language adapter without human-promoted board row |
| **Reject** | Full-project rebuild on file change |
| **Reject** | HTTP `POST /v1/index` (defer to Phase 43+ residual) |

## Must lock for S01-01 (delivered in thickened 01-implement)

1. Touch-list: policy doc → status JSON → watch CLI → error UX → tests.
2. Six acceptance tests G7-F1–F6.
3. Git-hook documented as default freshness path (not replaced by watch).

## Exit criteria

- [x] `01-implement.md` + `02-review.md` runnable with file targets + G7 accept map
- [x] SCOPE-TODOS G7-0 marked done
- [x] Board row → `done` with evidence in Notes

## Next

`P42-S01-01`
