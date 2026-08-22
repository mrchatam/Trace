# P10 / S03 / 01 — Index GC

## Metadata
- id: P10-S03-01
- todo_ids: [P10-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-20** per sibling **00-PLANNER** FINAL locks (2026-08-16). Full-tree `trace index` must drop missing paths from files/symbols/imports/FTS (delete-on-missing). Keep file-local argv isolation. **No new migration. No MCP index tools.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law 4 / DR-INCREMENTAL
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-20
- [phase README](../../README.md)
- Live: `cmd/trace/index.go`; `internal/store/file_graph.go`; `internal/store/fts.go` `SyncFileFTS`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Locks are FINAL — do not re-debate.

## Locked defaults (FINAL — P10-S03-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | `internal/store` (`ListFilePaths` + `DeleteFileByPath`); `cmd/trace/index.go` GC after full walk |
| Migration | **None** |
| Full-tree GC | `len(args)==0`: index live walk set, then delete every `files.path` ∉ live set |
| Explicit argv | No project-wide GC; missing-on-disk path → delete **that** path only |
| Delete | FTS file+symbol rows for path, then `DELETE FROM files` (CASCADE symbols/imports); idempotent |
| stderr | Include removed count (`indexed N, skipped M, removed K`) |
| Architecture | Set-diff deletes only — **forbidden** full-rebuild-on-any-change |
| Carry-forward | honesty A/B/C + Gate G; Gate E/F; ablation; Gate H; compat; p0x; x0; S01/S02; Gate C `dry_run:false` intact |
| Forbidden | New mig; MCP index; daemon/HTTP/embeddings; analyzer megastore rewrite; Mode-B Gate C rewrite; board spawn |

## Extension points (exact files)

| File | Work |
|------|------|
| `internal/store/file_graph.go` (or sibling) | `ListFilePaths`, `DeleteFileByPath` (+ FTS cleanup) |
| `internal/store/*_test.go` | Unit: delete removes symbols/imports/FTS; list paths; idempotent delete |
| `cmd/trace/index.go` | After full walk index: GC set-diff; argv missing → single-path delete; stderr `removed` |
| `cmd/trace/cli_test.go` (or `index_*_test.go`) | `TestIndexGCAfterPathRename` + keep `TestIndexIncrementalIsolation` |

## Role work

1. TDD store delete: insert file+symbols → `DeleteFileByPath` → GetFile/ListSymbols fail; FTS clean; second delete OK.
2. TDD CLI rename: index `a.js`+`b.js` → rename `a.js`→`c.js` → `trace index` (no args) → old path gone, `c.js` indexed, `b.js` intact.
3. Wire full-tree GC + argv missing-path delete; stderr removed count.
4. Confirm `TestIndexIncrementalIsolation` still passes (argv must not GC sibling).
5. Run locked verify suite; board **status + Notes only** (cite test names).

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-checks: after rename+full index, `GetFileByPath(old)` errors / empty symbols; sibling unchanged; single-file argv index does not remove other files.

## Exit criteria
- [ ] Full-tree `trace index` GC: missing/renamed-away paths removed from files + symbols + imports + FTS
- [ ] Explicit argv index does **not** project-wide GC (isolation green)
- [ ] Missing explicit argv path deletes that path only (when applicable)
- [ ] Regression test `TestIndexGCAfterPathRename` (or equiv) green
- [ ] No new mig; no MCP index tools; no full-rebuild architecture
- [ ] Carry-forward suite green; Gate C `dry_run:false` untouched
- [ ] Board Notes ready for **P10-S03-02**

## Out of scope
- S04 operator/capability transition gates (DF-17/18/24/26/31)
- S02 MCP surface changes
- Dependent-file cascade reindex (still file-local)
- Rewriting Mode-B Gate C packs / Phase 00–09 history

## Todo updates
Implementer: **status + notes only**. Record test names + DF-20 evidence. No spawning; no rewriting upcoming prompts.

## Minimal todos
- [ ] Store `ListFilePaths` + `DeleteFileByPath` (+ FTS) + unit tests
- [ ] Full-tree GC in `cmdIndex` + stderr removed count
- [ ] Argv missing-path single delete
- [ ] `TestIndexGCAfterPathRename` (+ isolation retained)
- [ ] Locked verify cmds green; board Notes
