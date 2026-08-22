# FM-02 / FR-P28-02 — notes (P28-S06-03)

**Date:** 2026-08-20  
**Status:** implemented (enforce fail-closed preserved)

## What shipped

1. **`GapPassPrompt` write-before-export** — steps 5–6 require ≥1 discovery OR ≥1 decision linked to `$TRACE_TASK_ID` **before** `seed export --strict --enforce`; thin graphs called out as enforce failures.
2. **Early thin-graph warn** — plain `seed export` (no `--strict`) emits stderr `seed export: warn: … thin graph …; write discoveries/decisions before seed export --strict --enforce` when discoveries=0 and decisions=0. Still writes; exit 0. Does **not** weaken `--strict --enforce`.
3. **Harness nudges** — `PROTOCOL.md` FM-02 paragraph; `PROMPT-G1-BUILD.md` / `PROMPT-G1-DIRECTED-GAP.md` ordered write-then-export.
4. **Help** — `seed export` documents early thin warn.

## Enforce regression (fail-closed)

`--strict --enforce` on a P26 thin fixture still exits blocked and does **not** write the output file (`TestSeedExportStrictEnforceBlocksP26ThinGraph`).

## Evidence

```bash
GOPROXY=direct go test ./cmd/trace/... -count=1 -run 'SeedExport|Enforce|Thin'
GOPROXY=direct go test ./internal/install/... -count=1 -run 'GapPass'
GOPROXY=direct go test ./internal/... ./cmd/trace/... -count=1
```

Primary new coverage: `TestSeedExportPlainThinGraphEarlyWarnWrites` (warn + write); GapPass assertions for `Write-before-export` + `--strict --enforce`.
