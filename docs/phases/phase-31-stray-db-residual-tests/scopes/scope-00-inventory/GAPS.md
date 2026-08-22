# GAPS — Phase 31 stray residual testing

**Date:** 2026-08-21
**Author row:** P31-S00-01
**Baseline tests:** `internal/store/stray_trace_db_test.go` —
`TestOpenWarnsWhenRootStubPresent`, `TestOpenExistingWarnsWhenRootStubPresent`, `TestOpenQuietWhenNoRootStub`, `TestOpenLeavesRootStubUntouched`
(**no** dir-named stub quiet case)
**open.go cite:** warn choke `warnIfStrayRootTraceDB(absRoot)` at `open.go:85` inside `openStore`; Stat-only regular-file gate `open.go:144–149` (`os.Stat` on `<absRoot>/trace.db`; `err != nil || !fi.Mode().IsRegular()` → quiet return). Live DB join: `traceDirName=".trace"` + `dbFileName="trace.db"` (`open.go:16–17`) → `filepath.Join(traceDir, dbFileName)` at `open.go:87` + `open.go:102`. Message constant notes once-per-`openStore` at `open.go:19–21`.

## Must-add (S01 must ship)

| ID | Gap | Preferred home | Notes |
|----|-----|----------------|-------|
| G1 | Dir-named root `trace.db` → quiet (no warn; open OK; live DB still `.trace/trace.db`) | `internal/store` unit in `stray_trace_db_test.go` (or sibling) | Code already quiet via `!IsRegular()` (`open.go:146–147`); **untested**. P30-S02-02 nit. Pattern: `os.Mkdir` root `trace.db`, capture `warnWriter`, `Open`, assert empty warn + `DBPath` under `.trace/`. |
| G5 | Durable dogfood repro (init → python/sqlite root stub → stderr warn → `.trace/trace.db` live; stub untouched) | `scripts/repro-stray-trace-db.sh` (create `scripts/` if needed) | P30 VERIFY one-shot only under `experiments/runs/2026-08-21-p30-s03-01-verify/evidence/`; **no** checked-in script; **no** `scripts/` dir yet. Prefer script over re-prosing VERIFY. |
| G6 | Document multi-open warn as intentional (once per `openStore`; multiple CLI opens may re-emit; no suppress flag) | Docs-only: short note in `CONTRIBUTING.md` and/or `AGENTS.md` | Code comment at `open.go:19` exists; operator-facing docs do **not**. P30 VERIFY residual listed multi-open as acceptable. **No** product suppress flag. |

## Nice-to-have

| ID | Gap | Notes |
|----|-----|-------|
| G2 | CLI `trace tasks -C <root>` stderr assertion when root stub is a regular file | P30 VERIFY one-shot PASS; no durable `go test`. Store unit already locks the warn choke; G5 script covers the same CLI stderr path. If added: reuse `cmd/trace` `runCapture` / `*_test.go` patterns — do not duplicate G1/G5 unless S01 wants an explicit CLI surface lock. |

## Out of scope (this phase)

| Item | Reason |
|------|--------|
| Store-path redesign | Locked — canonical remains `<root>/.trace/trace.db` |
| Silent / flagged delete of root stub | Locked — future-only (P30 residual) |
| GUI | Phase 32 |
| Persistent warn suppress flag | Locked residual; intentional once-per-`openStore` |
| Extra `/trace.db` in `web/.gitignore` or experiment `.gitignore`s | Not product store scaffolds; frontend-only / experiment harness. CONTRIBUTING already recommends consumer `/trace.db`. No `internal/install` scaffold emits `.gitignore` with `/trace.db`. |

## Defer-with-reason

| ID | Gap | Why defer |
|----|-----|-----------|
| G3 | Dedicated `trace serve` “startup” warn test | HTTP opens via `store.Open(s.root)` **per handler request** (`internal/httpapi/server.go:193–194`), not a process-startup open. `Open` choke already covered by store units (and G1). A serve-only harness adds little unless a serve-exclusive path appears. **Do not invent** a startup open. |
| G4 | Propagate `/trace.db` into non-product ignore files | Same as out-of-scope audit: only product scaffolds today are repo `.gitignore` + `fixtures/x0/.gitignore`. Revisit only if an install/scaffold writer starts emitting consumer `.gitignore`. |

## Candidate disposition (from 00-PHASE-PLANNER)

1. Dir-named root `trace.db` quiet — **must-add (G1)** — store unit; P30-S02-02 nit; code correct, test missing.
2. CLI `trace tasks -C` stderr — **nice-to-have (G2)** — covered operationally by G5 script + store unit; optional CLI harness.
3. `trace serve` startup warn — **defer (G3)** — no startup open; request-scoped `store.Open` only.
4. Extra ignore scaffolds — **out-of-scope / defer (G4)** — no additional product scaffolds; keep repo + `fixtures/x0` only.
5. Dogfood repro script — **must-add (G5)** — durable `scripts/repro-stray-trace-db.sh`.
6. Document multi-open warn — **must-add (G6)** — docs-only; no suppress flag.

## Ignore-scaffold audit

| Path | Has `/trace.db`? | Action |
|------|------------------|--------|
| `.gitignore` | yes (L3–4) | keep |
| `fixtures/x0/.gitignore` | yes (L2) | keep |
| `web/.gitignore` | no (frontend logs/node/dist only) | out-of-scope — not a store scaffold |
| `experiments/**/.gitignore` (ab-incident-tracker, ab-p25-gap-pass-validation, ab-library-hold-desk project/runs) | no | out-of-scope — experiment harness; not install templates |
| `internal/install` | N/A — no `.gitignore` emitter with `/trace.db` | keep as-is; no must-add scaffold write |

## Handoff to S01

- Must-add count: **3** (G1, G5, G6)
- Script required: **yes** (`scripts/repro-stray-trace-db.sh`)
- Docs-only IDs: **G6**
- Nice-to-have (optional): **G2**
- Deferred: **G3**, **G4**
- S01 must not: path change, delete, GUI, invent serve startup open, add warn suppress flag
