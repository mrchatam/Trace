# P07-S02-02 — Scope review notes (2026-08-16)

Independent review of S02 Go language adapter against `00-PLANNER.md` / `01-language-plugins.md` + TODO Notes for `P07-S02-01`. Fresh session; claims verified in-repo (reviewer did not implement S02-01).

## Plan (executed)

1. Diff claims vs `detect.go` / `index.go` `extract` / `extract_go.go` / `analyzers_test.go` + `testdata/sample.go`
2. Confirm `go.mod` `github.com/tree-sitter/tree-sitter-go` **v0.25.0**; go floor 1.24.0; adapter-shaped switches only (no plugin registry)
3. Confirm symbols `function`/`method`/`type`; imports path + named alias; no `011_*`; no Gate H / `evals/perf` invent
4. Confirm S01 T0 walk order + `vendor` SkipDir untouched; CGO analyzers-only
5. Re-run locked verify suite; Gate C `docs/verification/gate-c-x0/` spot-check
6. Severity-tag findings; spawn only for blocker/high (none)
7. Write these notes; mark board; unlock **P07-S03-00**; light-note S03

## Verdict

**APPROVE** — no blocker/high. Confidence: **high**. Spawns: **none**. Next board row: **P07-S03-00**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| Language = Go only (one new adapter) | Pass (`LangGo` + `.go` only; no other langs added) |
| Module `tree-sitter-go` **v0.25.0** | Pass (`go.mod` require + `go.sum`; import `bindings/go`) |
| Extension points = DetectLanguage + `extract` switch + `extract_go.go` | Pass |
| Adapter-shaped (no plugin registry / universal theater) | Pass (no new iface packages; switch only) |
| Symbols `function` / `method` / `type` | Pass (golden: Helper/Main funcs, Run method, Counter/ID/Worker types) |
| Imports path + named alias | Pass (golden: `fmt|`, `os|alias`, `path/filepath|`) |
| Golden + DetectLanguage `.go` | Pass (`TestIndexFileGoGolden`; `TestDetectLanguage` `a.go`→`LangGo`) |
| IndexFile path consistent; unsupported still SkipError | Pass (`TestSkipUnsupported` retained) |
| CGO analyzers-only | Pass (`CGO_ENABLED=0` store/vcs/gitcli/domain + honesty/replan/impact/capability; `tree-sitter-go` import only in `extract_go.go`) |
| No `011_*` | Pass (schema still through `010_capability_surface.sql`) |
| S01 T0 not regressed; `vendor` still SkipDir | Pass (`t0SkipDirs` + walk order T0→lang→T0 file→gitignore; `TestWalkIndexableT0AlwaysSkip`) |
| No Gate H threshold invent | Pass (no `evals/perf/`; no pass numbers / Gate H claim) |
| Carry-forward bars | Pass (honesty / replan / impact / capability / p0x / x0 / `./...`) |
| Gate C untouched | Pass (`dry_run:false` N=3; mean G1 0.8 > B0 0.0; packs not rewritten) |
| go 1.24.0 floor respected | Pass (`go.mod` `go 1.24.0`) |

## Re-verification commands (2026-08-16)

```text
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1
# ok

CGO_ENABLED=1 go test ./cmd/trace/... -count=1
# ok

CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/... ./internal/domain/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
# ok

CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
# ok

CGO_ENABLED=1 go test ./... -count=1
# ok
```

Named proofs: `TestIndexFileGoGolden`, `TestDetectLanguage` (`.go`), analyzers `TestIncrementalIsolation`, cmd/trace `TestWalkIndexableT0AlwaysSkip` / `TestIndexIncrementalIsolation`.

## Findings

### Blocker / high

None.

### Medium

None requiring spawn.

### Low / nit (residual)

- **No CLI walk golden for `.go`:** End-to-end `trace index` walk→Go file is implied by `DetectLanguage` + existing walk (step 2 uses DetectLanguage). Covered at IndexFile golden; optional future CLI plant if S03 wants ladder fixtures.
- **Blank/dot imports:** Spec allows named alias when cheap; `_` / `.` imports not in golden (paths still captured when present as `import_spec`). Acceptable for minimal lock.
- **Interface method elems:** `method` = receiver funcs only (locked); interface `method_elem` (e.g. `Counter.Inc`) not indexed — by design.

## Spawns

None.

## S03 note (upcoming)

Unlock **P07-S03-00**. S02 shipped Go adapter (`tree-sitter-go` v0.25.0) with goldens; **did not** seed `evals/perf` or invent Gate H thresholds. Optional tiny `.go` ladder fixtures remain OK for Gate H; **thresholds still after measurements**. Re-prove Go DetectLanguage/golden + S01 T0 + carry-forward at VERIFY.
