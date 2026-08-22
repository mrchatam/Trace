# P08-S03-02 — Scope review notes (2026-08-16)

Independent review of S03 production-hardening deliverables vs `00-PLANNER.md` locks + `P08-S03-01` Notes. Fresh session; claims re-verified in-repo (no implementer session shared).

## Verdict

**APPROVE** — no blocker / high findings. Confidence: **high**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| Embed `schema/NNN_*.sql` + `schema_migrations` apply on Open | Pass — `migrate.go` embed + `openStore` → `migrate(db)`; forward-only; no downgrade CLI |
| `MigrationStatus` / `trace migrate status` applied vs embed max | Pass — MaxApplied=EmbedExpected=10; CLI prints `max_applied`/`embed_expected`/`applied`/`pending` |
| Prefer no `011_*` | Pass — schema `001`…`010` only; no justified `011_*` needed |
| Backup = `trace.db` only (`VACUUM INTO`); not repo/source trees | Pass — `BackupTo` / `Backup`; lock held for in-process snapshot |
| Restore path-local `<abs>/.trace/`; Abs `root_path` rebind | Pass — `Restore` → `openStore(..., rebindRoot=true)` → `rebindProjectRoot` |
| Round-trip + `HasBlobLikeColumns` false; token excluded by default | Pass — `TestBackupRestoreRoundTrip` |
| Lock fail-closed on concurrent backup/restore | Pass — `TestBackupFailsWhenLocked` / `TestRestoreFailsWhenLocked` |
| Optional `.trace/access.token` + `TRACE_ACCESS_TOKEN`; fail-closed | Pass — `checkAccessToken` + `ErrUnauthorized`; `TestLocalAccessTokenFailClosed` |
| `trace auth set\|clear\|status`; status does not leak secret | Pass — prints `enabled`/`disabled` only |
| No cloud OAuth / hosted IdP / daemon auth | Pass — FS token + env only |
| MCP: no new tools; Open inherits gate (G19) | Pass — still six tools; `project.go` maps `ErrUnauthorized` like `ErrLocked` |
| Path-local bind + `trace.lock` unchanged (S02) | Pass — Open order mkdir→lock→auth→db→migrate |
| Gate C `dry_run:false` N=3 | Pass — `docs/verification/gate-c-x0/metrics-{b0,g1}.json` |
| Carry-forward Gate H + honesty A/B/C + E/F/G + ablation + p0x + x0 | Pass (independent re-run) |
| S04 VERIFY stubs cover migrate/backup/auth | Pass — `00-PLANNER` / `01-verify` Depends already name surfaces; light S03-02 note added |

## Verify (independent re-run)

```text
CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/... \
  ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... \
  ./evals/impact/... ./evals/capability/... -count=1                           PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... -count=1         PASS
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... \
  ./evals/replan/... ./evals/impact/... ./evals/capability/... ./evals/perf/... \
  -count=1                                                                     PASS
CGO_ENABLED=1 go test ./... -count=1                                           PASS
Gate C dry_run:false N=3                                                       OK
```

## Findings

None at blocker / high / medium (spawn threshold).

### Residuals (low — no spawn)

1. **`auth set <token>` argv exposure** — token may appear in process listings; local-first FS gate by design. Prefer env/`TRACE_ACCESS_TOKEN` for sensitive hosts; S04 checklist can note.
2. **Restore lock release → `openStore` window** — after install, flock is released then re-acquired for migrate/rebind. If another Open wins the window, restore fails fail-closed with `ErrLocked` (DB already installed; `--force` retry). Acceptable under single-writer model.
3. **No dedicated `--include-token` unit test** — default exclude covered by round-trip; opt-in path is thin CLI→`BackupOptions.IncludeToken`.
4. **Carry-forward S02:** lock/unauthorized CLI exit **2** (`exitFail`), not planner-literal **1** — live taxonomy.

## Spawns

None.

## Next

**P08-S04-00** (Phase 08 VERIFY / compat+security checklist planner). Do not start until orchestrator launches that row.
