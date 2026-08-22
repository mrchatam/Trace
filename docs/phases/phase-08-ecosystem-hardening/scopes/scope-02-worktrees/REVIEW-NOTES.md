# P08-S02-02 — Scope review notes (2026-08-16)

Independent review of S02 worktree / path-local bind deliverables vs `00-PLANNER.md` locks + `P08-S02-01` Notes. Fresh session; claims re-verified in-repo (no implementer session shared).

## Verdict

**APPROVE** — no blocker / high / medium findings. Confidence: **high**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| Root resolve Abs-only (`cmd/trace` + MCP); no walk-up | Pass — `resolveRoot` / `resolveProject` / `store.Open` use `filepath.Abs` only |
| `store.Open` binds `<absRoot>/.trace/trace.db` | Pass — `open.go` Join(absRoot, `.trace`, `trace.db`) |
| Exclusive `trace.lock` on Open → Close | Pass — `lock.go` `unix.Flock(LOCK_EX\|LOCK_NB)`; released in `Close` |
| Second Open → `store.ErrLocked`; CLI/MCP fail-closed | Pass — exported `ErrLocked`; CLI stderr + non-zero; MCP wraps `mcp: %w` |
| Isolation: two abs roots → two DBs; no cross-visibility | Pass — `TestProjectBindPathLocalIsolation` |
| Concurrent Open fail-closed; reopen after Close | Pass — `TestConcurrentStoreOpenFailClosed` |
| CLI lock surface | Pass — `TestInitFailClosedWhenStoreLocked` (exit `exitFail`) |
| `gitcli.OpenWithStore` companion (no double-Open) | Pass — CLI/MCP `tryOpenGit` → `OpenWithStore`; `ownsStore` false |
| No `011_*` | Pass — schema through `010_capability_surface.sql` only |
| No swarm / daemon / HTTP / adapter coupling | Pass — store/lock orthogonal to `LanguageAdapter` |
| Help note `-C` / one-writer | Pass — `cmd/trace/help.go` |
| S03 Depends still accurate | Pass — path-local `.trace` + respect `trace.lock`; S02 no mig |
| Gate C `dry_run:false` N=3 | Pass — `docs/verification/gate-c-x0/metrics-{b0,g1}.json` (3 runs each) |
| Carry-forward Gate H + honesty A/B/C + E/F/G + ablation + p0x + x0 | Pass (independent re-run) |

## Verify (independent re-run)

```text
CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/... \
  ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... \
  ./evals/impact/... ./evals/capability/... -count=1                           PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/perf/... ./evals/p0x/... \
  ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... \
  ./evals/capability/... -count=1                                              PASS
CGO_ENABLED=1 go test ./... -count=1                                           PASS
```

## Findings

None at blocker / high / medium.

### Residual severity judgment — CLI exit 2 vs planner exit 1

Planner lock said lock conflict → CLI exit **1**. Implementation uses existing P00-S07 taxonomy: `exitUsage=1`, `exitFail=2`, and maps Open/`ErrLocked` → **`exitFail` (2)**. Checklist for this review requires **non-zero** + clear stderr — both met. Exit **2** is the correct operational-fail code under the live CLI contract; remapping to 1 would collide with usage errors. **Severity: low / residual** — no spawn.

### Nit (no spawn)

- `OpenWithStore` does not assert store `projectRoot` equals git abs root (caller contract). CLI/MCP always pass the same abs — acceptable.
- No explicit nested-dir “no walk-up” plant; Abs-only Open cannot walk up — isolation tests suffice.

## Spawns

None.

## Residuals

- Lock conflict CLI exit **2** (`exitFail`), not planner-literal **1** — intentional alignment with `exitUsage`/`exitFail` taxonomy (see above).
- `trace.lock` via `golang.org/x/sys/unix` flock — Linux/unix CI path; non-unix portability not claimed this scope.
- Optional linked-git-worktree dual-root plant not required; path-local Abs isolation covers multi-root policy.

## Next

**P08-S03-00** (production-hardening scope planner). Do not start until orchestrator launches that row.
