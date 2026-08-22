# Phase 31 — Stray `trace.db` residual testing

**Active** after Phase 30 close (`P30-S03-02`, 2026-08-21). Human-promoted: Phase 30 product (warn / gitignore / docs) stays; this phase adds **tests + repro evidence** only. Does **not** redesign store paths or add silent delete.

## Goal

Close the “Phase 30 needs more testing” gap. Expand automated coverage and a durable repro script around the shipped warn / path / gitignore behavior so regressions are cheap to catch before Phase 32 GUI work.

## Light locks (do not reopen)

| Lock | Value |
|------|-------|
| Canonical store | `<projectRoot>/.trace/trace.db` via `store.Open` / `OpenExisting` |
| Path redesign | **Forbidden** — Phase 30 S00 already confirmed agent hygiene |
| Delete policy | No silent delete of root stubs; no delete flag in this phase |
| GUI / HTTP features | **Out** — Phase 32 owns graph-first GUI |
| Product scope | Tests + optional dogfood script + docs notes only |
| Successor on green | **Phase 32** — first runnable `P32-00` |

## Live repo baseline (P31-00, 2026-08-21)

| Area | State |
|------|-------|
| Store warn | `internal/store/open.go` — `warnIfStrayRootTraceDB` in `openStore`; Stat-only regular-file; once per open; stderr / `warnWriter` |
| Unit tests | `internal/store/stray_trace_db_test.go` — 4 cases (Open warn, OpenExisting warn, quiet no stub, stub untouched) |
| Known gaps (candidates) | Dir-named root quiet; CLI `trace tasks -C` stderr; `trace serve` startup warn; extra ignore scaffolds; durable repro script; multi-open documented as intentional |
| Gitignore | `/trace.db` in `.gitignore` + `fixtures/x0/.gitignore` |
| P30 VERIFY | PASS — residuals non-blocking; see `phase-30-…/scopes/scope-03-verify/VERIFY-NOTES.md` |

## Scope index (serial)

```
S00 inventory → S01 tests (+ script) + review → S02 VERIFY + DR-HANDOFF → Phase 32
```

| Scope | Title | Primary artifact |
|-------|-------|------------------|
| S00 | Gap inventory vs P30 residuals + live tests | `GAPS.md` |
| S01 | Implement must-add tests (+ optional repro) + review | tests / script |
| S02 | VERIFY + close handoff | `VERIFY-NOTES.md`, successor Phase 32 |

## In scope

- Confirm / trim coverage gaps into `GAPS.md`
- Ship must-add `go test` coverage (prefer `internal/store`; CLI/`serve` only if gap requires)
- Optional durable dogfood script under `scripts/` or `experiments/`
- Document intentional multi-open warn (once per `openStore`) if still only in VERIFY residuals
- Phase VERIFY + DR-HANDOFF → Phase 32

## Out of scope

- Changing canonical store path or dual-store behavior
- Silent or flagged delete of root stubs
- Phase 32 GUI / explorer work
- Rewriting Phase 30 `done` history
- Cloud / hosted SaaS

## Handoff

[`DR-HANDOFF.md`](DR-HANDOFF.md) — **OPEN** until `P31-S02-02`. Successor lean: **Phase 32**.

## Phase planner

[`00-PHASE-PLANNER.md`](00-PHASE-PLANNER.md) — row `P31-00`.
