# P11 / S01 / 01 — Index partial-path GC

## Metadata
- id: P11-S01-01
- todo_ids: [P11-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-40** per sibling **00-PLANNER** FINAL locks (2026-08-16). After rename, `trace index <new-path>` must remove ghost old path/symbols/FTS via **content-hash orphan GC**, without project-wide argv set-diff and without full-rebuild. **No new migration. No MCP index tools.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law 4 / DR-INCREMENTAL
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-40
- [experiments/_post_p10/BUGHUNT.md](../../../../../experiments/_post_p10/BUGHUNT.md) — DF-40 repro
- [phase README](../../README.md)
- Live: `cmd/trace/index.go`; `internal/store/file_graph.go` (`ListFilePaths`, `DeleteFileByPath`, `GetFileByPath` / `ContentHash`)
- Prior: P10 S03 DF-20 full-tree GC — keep green
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Locks are FINAL — do not re-debate. If 00-PLANNER is still DRAFT, stop and return to planner.

## Locked defaults (FINAL — P11-S01-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | `cmd/trace/index.go` + optional `internal/store` helper `ListFilePathsByContentHash` |
| Migration | **None** |
| Partial GC | After successful argv index of path P with hash H: delete other DB paths with hash H that are **missing on disk** |
| Full-tree / missing-argv | **Unchanged** from P10 S03 |
| Isolation | Never delete on-disk siblings on partial argv |
| stderr | `indexed N, skipped M, removed K` |
| Architecture | Orphan / set-diff deletes only — **forbidden** full-rebuild-on-any-change |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` intact; DF-20 tests green |
| Forbidden | Project-wide GC on every argv index; new mig; MCP index; daemon/HTTP/embeddings; analyzer megastore rewrite; Mode-B Gate C rewrite; board spawn |

## Extension points (exact files)

| File | Work |
|------|------|
| `internal/store/file_graph.go` (or sibling) | Optional `ListFilePathsByContentHash(hash string) ([]string, error)` — ordered paths with that `content_hash` |
| `internal/store/*_test.go` | Unit for hash list helper (if added) |
| `cmd/trace/index.go` | After successful partial `indexOne` (or after argv loop): content-hash orphan GC; bump `removed`; keep full-tree + missing-argv paths intact |
| `cmd/trace/cli_test.go` (or `index_*_test.go`) | **`TestIndexPartialArgvGCAfterRename`**; keep DF-20 + isolation tests |

## Role work

1. TDD CLI: index `a.js`+`b.js` → rename `a.js`→`c.js` → `trace index c.js` only → assert old path/symbols/FTS gone, `c.js` present, `b.js` intact, stderr has `removed`.
2. Optional store helper for hash→paths; unit test.
3. Wire orphan GC in `cmdIndex` for `!fullTree` after successful index (not on skip/fail).
4. Confirm P10 regressions: `TestIndexGCAfterPathRename`, `TestIndexArgvMissingPathDeletesOnlyThatPath`, `TestIndexIncrementalIsolation`.
5. Run locked verify suite; board **status + Notes only** (cite test names + DF-40).

## Algorithm sketch (non-normative — locks above win)

```text
for each successfully indexed argv path P:
  H = GetFileByPath(P).ContentHash
  for each Q in ListFilePathsByContentHash(H):
    if Q == P: continue
    if Q missing on disk: DeleteFileByPath(Q); removed++
```

Do **not** run this as a substitute for full-tree set-diff when `len(args)==0`.

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-checks: rename + partial argv clears ghost; sibling on disk survives; full-tree GC still works; missing argv still deletes only that path.

## Exit criteria
- [ ] `TestIndexPartialArgvGCAfterRename` (or equiv) green — DF-40
- [ ] P10 DF-20 tests still green (full-tree + missing-argv + isolation)
- [ ] No project-wide argv set-diff; no new mig; no MCP index; no full-rebuild architecture
- [ ] Carry-forward suite green; Gate C `dry_run:false` untouched
- [ ] Board Notes ready for **P11-S01-02**

## Out of scope
- S02+ Phase 11 product work (DF-43/44, locks, etc.)
- Rename+edit same move when hash diverges (residual; full-tree still GC’s)
- Dependent-file cascade reindex
- Rewriting Mode-B Gate C packs / Phase 00–10 history

## Todo updates
Implementer: **status + notes only**. Record test names + DF-40 evidence. No spawning; no rewriting upcoming prompts.

## Minimal todos
- [ ] Optional store `ListFilePathsByContentHash` + unit test
- [ ] Partial argv content-hash orphan GC in `cmdIndex` + stderr `removed`
- [ ] `TestIndexPartialArgvGCAfterRename` (+ DF-20/isolation retained)
- [ ] Locked verify cmds green; board Notes
