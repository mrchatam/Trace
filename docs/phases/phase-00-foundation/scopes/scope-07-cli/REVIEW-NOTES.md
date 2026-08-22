# P00-S07-02 — Scope review notes (2026-08-15)

Independent review of S07 against `01-cli.md` + TODO Notes for `P00-S07-01`. Fresh session; claims verified in-repo.

## Plan (executed)

1. Diff claims vs `cmd/trace/*` + wiring to store/gitcli/analyzers/domain/retrieval/compiler
2. Grep for cobra/MCP/daemon/HTTP/dump; confirm G19 (no ranking/FTS/provenance machines/parsers in CLI)
3. Re-run `CGO_ENABLED=1 go test ./...` + `go build -o bin/trace ./cmd/trace` + `version` smoke
4. Severity-tag findings; inline-fix medium (index AtRev fallthrough); thicken S08 seed-path notes
5. Write these notes; mark board + SCOPE-TODOS

## Verdict

**APPROVE** — one **medium** fixed inline (`IndexFileAtRev` error masking). Confidence: **high**. Spawns: **none**. Next board row: **P00-S08-00**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| Thin adapter G19 / DR-API — parse → library → print/exit; no ranking/FTS SQL/work-state graphs | Pass (`add`/`link`/`transition`/`seed` → `domain.*`; `why` → `retrieval.Why`; `context` → `compiler.TaskContext`/`ExpandContext`; `index` → `analyzers.IndexFile*`) |
| Locked commands: help/version/init/index/reindex/add/link/transition/seed import/why/context | Pass (`root.go` switch + `help.go`) |
| No cobra / urfave; stdlib argv + `flag` | Pass (`go.mod` clean; `cmd/trace` only) |
| `-C` / `--project`; Abs before `store.Open` | Pass (`parseGlobal` + `resolveRoot`) |
| Seed JSON v1; reject unknown top-level keys; creates then links then `TransitionTask` | Pass (`seed.go`; tasks omit work_state; transitions via `TransitionTask`) |
| Causal mutations only via `domain.*` (no raw store Upsert for entities) | Pass (grep: no UpsertGoal/Task in CLI) |
| Index file-local; SkipError → skip; reindex alias; no full rebuild | Pass (`index.go` + `TestIndexIncrementalIsolation`) |
| why JSON keeps reason_code; context items + budgets + trust labels | Pass (`TestCausalWhyContextRoundTrip`, `TestSeedImportAndWhy`) |
| Exit 0/1/2; stdout machine / stderr progress | Pass (constants + index progress on stderr) |
| No MCP/daemon/HTTP/dump; `internal/*` does not import `cmd/` | Pass (grep) |
| CGO binary + tests; README Build note | Pass (re-verify below; README §Build) |
| S08 can consume locked seed + commands | Pass (S08 `01-fixture.md` thickened with surface + seed path semantics) |

## Re-verification commands (2026-08-15)

```text
CGO_ENABLED=1 go test ./... -count=1                         # ok (cmd/trace + all packages)
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace              # ok
./bin/trace version                                          # 0.0.0-dev
```

## Findings

### Blocker

None.

### High

None.

### Medium (fixed inline)

- **`indexOne` masked non-missing AtRev errors:** After `IndexFileAtRev` failed for any reason, CLI always fell through to working-tree `IndexFile`. Failure mode: store/git/SkipError paths could be retried or partially obscured instead of failing/skipping correctly. **Fix:** fall through only when `errors.Is(err, vcs.ErrNotFound)`; otherwise return the error. (`cmd/trace/index.go`)

### Medium (residual — no spawn)

- **`seed import <file>` path is cwd-relative**, not rewritten under `-C`. Documented for S08 harness (use abs path or matching cwd). Acceptable CLI convention for P0-X.

### Low / nit

- `add`/`link`/`transition` do not use `flagsFirst` (unlike `context`); kind/rel must precede flags — matches locked usage.
- Walk + best-effort `git check-ignore` lives in CLI — adapter orchestration, not domain logic (G19 OK).
- Prior S06 residual (silent GetGoal/Why omit inside compiler) unchanged; CLI surfaces library errors when the packet call itself fails.

## Spawns

None.

## Residual risks

- Seed re-import idempotency relies on Create* upsert semantics (locked acceptable).
- Index without git always uses working-tree bytes (correct for non-repos).
- Q-INJECTION / untrusted labeling remain library concerns (A14).

## Forward edits this review

- `cmd/trace/index.go` — AtRev fallthrough only on `vcs.ErrNotFound`
- `scopes/scope-08-fixture-p0x/01-fixture.md` — CLI surface + seed path + link rel lock notes
- `SCOPE-TODOS.md` — mark S07-02 done
- `docs/TODO.md` — `P00-S07-02` → done + notes
