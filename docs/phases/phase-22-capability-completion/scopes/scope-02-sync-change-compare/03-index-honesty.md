# P22-S02-03 — Implement: graph sync honesty + incremental index

## Metadata
- id: P22-S02-03
- todo_ids: [P22-S02-03]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

Keep graph state **synchronized** with the project, with **honesty** when it lags (**C04**). Wire index/hook watermark; surface HEAD ≠ last indexed commit. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## Live baseline

| Present | Absent |
|---------|--------|
| `compiler/buildIndexHonesty` — disk hash drift (`index_honesty.go`) | `graph_sync_state` table |
| `vcs_meta` / `MetaVCSWatermark` — VCS commit index (`gitcli/refresh.go`) | HEAD-vs-indexed-commit notice |
| `cmd/trace/index.go` — Refresh + file-local index | `trace index status` subcommand |
| Compat ceiling **22** | mig **023** |

**Do not conflate:** `MetaVCSWatermark` (VCS history index) vs **`last_indexed_commit`** (symbol/file graph indexed at HEAD).

## Locked defaults

| Item | Value |
|------|-------|
| Mig | **`internal/store/schema/023_graph_sync.sql`** |
| Table | `graph_sync_state` — single row **`id = 1`**: `last_indexed_commit TEXT`, `last_indexed_at TEXT`, `hook_installed INTEGER NOT NULL DEFAULT 0` |
| Compat | Bump **`evals/compat/compat_test.go`** + **`evals/compat/doc.go`** ceiling to **23**; forbid **024+** |
| Watermark update | After **successful** `cmdIndex` when git repo present: set `last_indexed_commit = repo.Head()`, `last_indexed_at = now RFC3339` |
| Partial argv index | Updates watermark only when **all** requested paths index without hard fail (existing stderr summary ok) |
| Hook flag | `hook_installed=1` when `trace install git-hook --write` succeeds (S02-01); cleared on uninstall |
| Honesty — commit lag | New **`GraphSyncHonesty`** (or extend `Packet`) on context compile: when HEAD ≠ watermark → notice string + `stale_commit: true` |
| Honesty — disk lag | **Keep** existing `IndexHonesty` unchanged |
| CLI | **`trace index status`** — JSON: `{ "head", "last_indexed_commit", "stale", "hook_installed" }` |
| Loop | `internal/loop/next.go` `contextSectionFreshness` — treat commit lag as **dirty** (like disk stale) |
| Full rebuild | **Forbidden** (Law 12 / D-22-15) |

## Requirements

1. Store helpers: `GetGraphSyncState`, `UpsertGraphSyncState` in new `internal/store/graph_sync.go`.
2. Wire watermark update at end of `cmdIndex` (`cmd/trace/index.go`).
3. Emit commit-lag honesty from `internal/compiler/` (new small helper; wire in `compiler.go`).
4. Implement `cmdIndexStatus` subcommand (split `cmdIndex` dispatch in `index.go` or `index_status.go`).
5. Named tests + compat bump.

## Touch files

- `internal/store/schema/023_graph_sync.sql`, `graph_sync.go`, `graph_sync_test.go`
- `cmd/trace/index.go` (+ status subcommand)
- `internal/compiler/graph_sync_honesty.go` (or extend `packet.go`), `compiler.go`
- `internal/loop/next.go` (dirty when commit stale)
- `evals/compat/compat_test.go`, `evals/compat/doc.go`

## Named tests

| Test | Proves |
|------|--------|
| `TestGraphSyncStaleWhenHeadDiffers` | HEAD advanced, watermark old → status JSON `stale=true` + packet/compile notice |
| `TestHookIndexUpdatesLastIndexedCommit` | `cmdIndex` after commit → watermark = HEAD |
| `TestMigrationStatusReportsEmbedMax` | embed **23** |
| `TestOpenCreatesDBAndMigratesIdempotent` | 023 applies cleanly |
| `TestCompatibilitySecurityChecklist` | ceiling **23**, no 024+ |

```bash
go test ./internal/store/... ./internal/compiler/... ./cmd/trace/... -count=1 -run 'TestGraphSyncStaleWhenHeadDiffers|TestHookIndexUpdatesLastIndexedCommit|TestMigrationStatusReportsEmbedMax|TestOpenCreatesDBAndMigratesIdempotent'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

## Exit criteria

- [ ] C04 true (sync path + honesty when lagging)
- [ ] Compat **23**; no 024+
- [ ] Disk `IndexHonesty` keepers still PASS
- [ ] Board Notes → **Next `P22-S02-04`**

## Minimal todos

- [ ] Mig 023 + store helpers
- [ ] Index watermark + status subcommand
- [ ] Packet/loop honesty + tests + compat
- [ ] Board notes
