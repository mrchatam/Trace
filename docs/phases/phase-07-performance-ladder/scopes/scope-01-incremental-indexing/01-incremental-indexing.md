# P07 / S01 / 01 — Incremental indexing / ignore tiers

## Metadata
- id: P07-S01-01
- todo_ids: [P07-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **T0 always-skip** ignore tiers and **measurable incremental indexing quality** against live Trace CLI walk + analyzers. Keep **file-local** incremental (no full-rebuild-on-any-change). Preserve all carry-forward bars. Do **not** declare Gate H pass. Do **not** add languages (S02).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) (locks finalized 2026-08-16)
- [phase README](../../README.md)
- [docs/STORAGE_AND_PERFORMANCE.md](../../../../STORAGE_AND_PERFORMANCE.md) §5 incremental + §9 tiers
- Live: `cmd/trace/index.go` (`walkIndexable`, `gitIgnored`, `indexOne`); `internal/analyzers` (`IndexFile`/`IndexFileAtRev`, SHA-256, NUL binary `SkipError`); `TestIndexIncrementalIsolation`; store schema through `010_capability_surface.sql`

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Go version | Keep `go.mod` floor (currently 1.24.0); do not downgrade |
| Package | **`cmd/trace`** walk/index helpers only — prefer extend `index.go` (+ tests in `cli_test.go` or `index_test.go`). Keep `internal/analyzers` file-local upsert + binary NUL probe. **Do not** invent `internal/indexer` / rewrite analyzers pipeline |
| Architecture | **File-local incremental only** — reindex path A must not cascade-rebuild siblings (DR-INCREMENTAL) |
| Migration | **None** — prefer path filter; **no** `011_*` unless implementer discovers a hard store need (must Note + stop for review; default = no mig) |
| Ignore model | S01 owns **T0 always-skip**. T1–T3 remain **notes only** (no structural-only / semantic tier engine) |
| gitignore | Keep best-effort `git check-ignore` **after** T0 filters |
| Binary | Keep analyzers `SkipError` (NUL in first 8KiB); do not add broad extension blacklists beyond DetectLanguage + T0 file suffixes below |
| Gate H | **Do not** claim Gate H pass; optional seed fixtures/timing hooks for later `evals/perf` only |
| Carry-forward | Honesty A/B/C; Gate G; Gate E; Gate F; capability ablation; p0x 7/7; x0; Gate C `dry_run:false` artifacts intact |
| Product surface | CLI + analyzers + existing store; no daemon/HTTP/embeddings; no VerifiedFact / `plan simulate` |
| Out | New language grammars (S02); Gate H thresholds (S03); commercial multi-model perf theater |

### T0 always-skip — directory basenames (locked)

Skip the directory entirely (`filepath.SkipDir`) when `d.Name()` matches **exactly** (case-sensitive):

```text
.git          # already live
.trace        # already live
node_modules
vendor
__pycache__
.venv
venv
dist
.next
target
coverage
```

Do **not** add `build` / `bin` in this row (too easy to hide real sources). Future rows may promote with measurement Notes.

### T0 always-skip — file suffixes (locked)

After `DetectLanguage` would accept a path, still skip when the **basename** ends with:

```text
.min.js
.min.mjs
.min.cjs
```

### T0 path-segment rule (locked)

A relative path is T0-skipped if **any** path component equals a T0 directory basename above (including when the user passes an explicit argv path under `node_modules/…`, `vendor/…`, etc.). Walk uses `SkipDir`; explicit `trace index <path>` must also refuse/skip T0 paths (count as skipped, not hard fail — prefer same `skipped++` path as `SkipError` / silent filter; document chosen behavior in Notes).

### Walk order (locked)

```text
filepath.WalkDir:
  1. dirs → if T0 basename → SkipDir
  2. files → DetectLanguage; if unsupported → skip
  3. if T0 file suffix OR T0 path-segment → skip
  4. if useGitIgnore && gitIgnored → skip
  5. append rel
```

### Helper surface (behavior locked; names may vary slightly)

```text
# cmd/trace (unexported OK)
isT0SkipDir(name string) bool
isT0SkipPath(rel string) bool   // path segments + .min.js/.min.mjs/.min.cjs
# wire into walkIndexable; also gate indexOne / normalize path before analyzers when args given
```

### Measurement surface (locked)

| Proof | Location | Required |
|-------|----------|----------|
| Sibling isolation | Keep / extend `TestIndexIncrementalIsolation` in `cmd/trace` | **Yes** — unchanged sibling symbols after indexing only A |
| T0 skip correctness | New test e.g. `TestWalkIndexableT0AlwaysSkip` | **Yes** — planted `node_modules/`, `vendor/`, `__pycache__/`, `dist/`, `.next/`, `target/`, `coverage/`, `.venv/`/`venv/` trees with otherwise-indexable `.js`/`.py` must **not** appear in `walkIndexable`; planted `foo.min.js` skipped; control `src/ok.js` included |
| Explicit T0 argv | Test e.g. `TestIndexSkipsExplicitT0Path` | **Yes** — `trace index node_modules/x.js` does not upsert that path |
| Analyzers isolation | Existing `internal/analyzers` `TestIncrementalIsolation` | Keep green (no rewrite required unless broken) |
| Optional latency hook | `t.Logf` wall times and/or tiny planted tree under `evals/perf/fixtures/` (no threshold asserts) | **Optional** — for S03 Gate H; **must not** invent pass numbers or `TestPlantedPerfLadderGateH` pass claim |

### Target tree

```text
cmd/trace/
  index.go           # T0 helpers + walkIndexable / explicit-path gate
  cli_test.go        # or index_test.go — isolation + T0 proofs
  # optional only:
evals/perf/
  fixtures/…         # planted skip/index trees for later Gate H (no thresholds)
```

No new store schema files by default.

### Policy (locked)

```text
File-local incremental only — no dependent cascade / full-repo rebuild
T0 always-skip before gitignore
T1 structural / T2 semantic / T3 deep = OUT (notes only)
No 011_* unless blocked (default none)
No Gate H pass / no threshold invention
No new languages
Daemon/HTTP/embeddings OUT
```

## Role work

1. Add T0 helpers; wire into `walkIndexable` and explicit-path indexing.
2. Add required automated tests (sibling isolation + T0 walk + explicit T0 argv).
3. Optionally seed `evals/perf/fixtures` or timing logs for S03 — no Gate H claim.
4. Run VERIFY commands below; update board Notes only.

## Exit criteria

- [ ] T0 directory + suffix + path-segment rules implemented per locks
- [ ] `TestIndexIncrementalIsolation` green (sibling isolation)
- [ ] New T0 walk + explicit-path skip tests green
- [ ] No `011_*` (or Notes explain blocker + no silent schema invent)
- [ ] Carry-forward suite green (commands below)
- [ ] No daemon/HTTP/embeddings; Gate C artifacts untouched; no Gate H pass claim
- [ ] Board Notes ready for **P07-S01-02**

## VERIFY commands (implementer must run)

```bash
# Focus
CGO_ENABLED=1 go test ./cmd/trace/... -count=1

# Analyzers isolation (keep)
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1

# Carry-forward bars
CGO_ENABLED=0 go test ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Do **not** rewrite Mode-B Gate C packs under `docs/verification/gate-c-x0/`.

## Minimal todos

- [ ] T0 helpers + wire walk / explicit paths
- [ ] Tests: isolation + T0 walk + explicit T0 skip
- [ ] Optional: evals/perf fixture seed or timing logs (no thresholds)
- [ ] Run VERIFY commands; board Notes for P07-S01-02

## Out of scope

- New language plugins (S02)
- Declaring Gate H pass / inventing Gate H thresholds (S03)
- T1–T3 promotion engines / megastore ignore DB
- Weakening prior gates / rewriting done Phase 06 prompts
