# DR-HANDOFF — Phase 24

**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Opened | 2026-08-20 |
| Predecessor | Phase 23 closed 2026-08-20 (`no successor` at P23 close; human promoted P24 after E01) |
| Trigger | E01 dogfood + **directed gap session** (Session B): Trace mode-dependent — discoveries work when asked; build mode does not |
| Phase 24 outcome | Investigation complete through S04 — ranked interventions + consolidated FINDINGS |
| Closed | 2026-08-20 |
| Successor decision | **Phase 25 — P25-C orchestrator + default gap pass** |
| Recommended promotion order | **P25-C → P25-A → P25-B** (one theme per phase) |
| Top interventions | **INT-03, INT-04, INT-11** |
| Residuals (non-blocking) | Auto-spawn human gate (P25-A); P19 threshold validation + sticky STOP reason UX (P25-B); hook API drift (P25-C INT-11); live gate env-dependency note |
| Forward | Phase 25 scaffold created; next runnable **P25-00** |

## Scope checklist (closed)

- [x] S01: Dogfood post-mortem + failure taxonomy → `POSTMORTEM.md`, `FINDINGS.md` draft
- [x] S02: Codebase loop/policy/task-creation audit → `CODEBASE-AUDIT.md`
- [x] S03: External + similar-project research → `EXTERNAL-RESEARCH.md`
- [x] S04: Intervention matrix → `INTERVENTION-MATRIX.md`, consolidated `FINDINGS.md`
- [x] S05: VERIFY evidence + successor recommendation

## Intervention summary (top 3 ranks)

1. **INT-03** — Default gap-pass install bundle so build sessions end with mandatory Trace gap review (FM-09, FM-04).
2. **INT-04** — Orchestrator Trace-first: parent sets `TRACE_TASK_ID` and failClosed `preToolUse` gate (FM-04, FM-05).
3. **INT-01** — Discovery→task promotion via `loop apply` `spawned_tasks[]` / guided `trace add task` (FM-10, FM-01, FM-08).

Full matrix: [INTERVENTION-MATRIX.md](scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md).

## Recommended Phase 25 themes

Human promotes **one theme per phase**; order below is suggested priority, not a mega-phase.

| Theme | One-line | Status | Evidence |
|-------|----------|--------|----------|
| **P25-C** | **Orchestrator + default gap pass** — install gap-pass prompt bundle; parent orchestrator failClosed hook; hook drift checks | **Recommended (1st).** Collapses Mode A→B without custom human gap prompt. | INT-03, INT-04, INT-11 ([INTERVENTION-MATRIX.md](scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md)) |
| **P25-A** | **Discovery → task promotion** — spawn path, MCP tool ordering, harness nudge after discovery | **Recommended (2nd).** Addresses Session B FM-10 (7 discoveries, 0 new tasks). | INT-01, INT-06 |
| **P25-B** | **Loop policy recalibration** — P19/hop thresholds, gap-pass deliberation reset, unified STOP reason UX | **Recommended (3rd).** Unblocks verify after gap fixes (FM-03). | INT-02, INT-05, INT-09 |
| P25-D | **Experiment protocol v2** — two-session rubric, arm isolation, `score.sh` fix | **Deferred.** Strengthens measurement; does not alone change live dogfood behavior. | INT-08, INT-10 |
| P25-E | **Graph honesty** — `--strict` export when discoveries lack linked tasks | **Deferred.** CI/export gate; pair after P25-A promotion ships. | INT-07 |

## Deferred / human-gate (from matrix §4)

- **Auto-spawn from discoveries:** Human product call — guided promotion (INT-01) vs autonomous AR-style `publishEvent`; not ranked as standalone row.
- **SQLite episode model spike:** Multi-phase; INT-05 minimal reset first; full Graphiti patterns out of Trace law for P0-X.
