# P08-S01-02 — Scope review notes (2026-08-16)

Independent review of S01 plugin-API deliverables vs `00-PLANNER.md` locks + `P08-S01-01` Notes. Fresh session; claims re-verified in-repo (no implementer session shared).

## Verdict

**APPROVE** — no blocker / high / medium findings. Confidence: **high**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| `LanguageAdapterAPIVersion == 1` exported in `internal/analyzers` | Pass (`language_adapter.go`) |
| `LanguageAdapter` iface: `ID` / `Extensions` / `Extract` | Pass |
| Compile-time static `builtinAdapters` (JS/TS/TSX/Py/Go) | Pass — no public `Register`; no `plugin.Open` / `.so` |
| `DetectLanguage` + `extract` dispatch via adapters | Pass (`detect.go`, `index.go` `extract`) |
| `IndexFile` / `IndexFileAtRev` orchestration unchanged | Pass — Upsert → SetFileLanguage → extract → Replace* |
| Lang* ids + goldens + DetectLanguage + isolation | Pass (`TestIndexFile*Golden`, `TestDetectLanguage`, `TestIncrementalIsolation`) |
| `TestLanguageAdapterAPIVersion` + `TestBuiltinLanguageAdaptersContributionPath` | Pass |
| `docs/ANALYZER_CONTRIBUTION.md` + `doc.go` refresh | Pass (Go listed; versioned steps; forbidden patterns) |
| No store migration `011_*` | Pass — schema through `010_capability_surface.sql` only |
| File-local incremental / no full-rebuild | Pass — IndexFile path-local; doc.go DR-INCREMENTAL |
| S02 stubs: no adapter↔worktree coupling | Pass — S02 Depends notes orthogonal; no thickening needed |
| Gate C `dry_run:false` N=3 | Pass (`docs/verification/gate-c-x0/metrics-{b0,g1}.json`) |
| Carry-forward: Gate H + honesty A/B/C + E/F/G + ablation + p0x + x0 | Pass (independent re-run) |

## Verify (independent re-run)

```text
CGO_ENABLED=1 go test ./internal/analyzers/... -count=1                         PASS
CGO_ENABLED=1 go test ./evals/perf/... ./evals/p0x/... ./evals/x0/... \
  ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1  PASS
CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/... \
  ./internal/domain/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... \
  ./evals/capability/... -count=1                                              PASS
CGO_ENABLED=1 go test ./... -count=1                                           PASS
```

## Findings

None at blocker / high / medium.

### Nit (no spawn)

- `init()` in `language_adapter.go` last-wins on duplicate extensions without panic; contribution-path test already fails on duplicates — acceptable for v1 static table.

## Spawns

None.

## Residuals

- Package-local `builtinAdapters` slice is a `var` (not frozen); product path correctly has no exported Register — contributors still append in-tree only.
- S04 `evals/compat` may cite `LanguageAdapterAPIVersion` + contribution doc (S04-00 owns finalization).

## Next

**P08-S02-00** (worktrees scope planner). Do not start until orchestrator launches that row.
