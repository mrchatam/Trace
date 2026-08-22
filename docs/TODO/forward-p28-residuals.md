# Forward queue — Phase 28 residuals (R6 / FM gaps)

**Status:** **SUPERSEDED / CLOSED** (2026-08-20) — no longer a runnable backlog; residual wave closed at `P28-S07-02`.

**Superseded by:** Phase 28 residual wave board rows on [`docs/TODO/phase-28.md`](phase-28.md) (**S06** FR-P28-01…07 → `P28-S06-01`…`14`; **S07** VERIFY) — all `done`.

Do **not** treat this file as the master run board. Orchestrator is **idle** (Phase 28 complete; successor **no successor**).

**Historical source (frozen text below):** Phase 28 closed at `P28-S05-02` with successor `no successor`; human later promoted FR-* onto Phase 28 S06/S07 without rewriting S00–S05 `done` history; residual wave closed at `P28-S07-02`. Evidence SoT: [`RESIDUAL-AUDIT.md`](../phases/phase-28-residuals-validation/scopes/scope-00-residual-audit/RESIDUAL-AUDIT.md), [`DR-HANDOFF.md`](../phases/phase-28-residuals-validation/DR-HANDOFF.md) (S05 CLOSED + Residual wave CLOSED).

---

## Promotion status (2026-08-20)

| ID | Residual | Board status |
|----|----------|--------------|
| FR-P28-01 | R6 / FM-01 | **promoted** → `P28-S06-01` / `P28-S06-02` |
| FR-P28-02 | R6 / FM-02 | **promoted** → `P28-S06-03` / `P28-S06-04` |
| FR-P28-03 | R6 / FM-04 | **promoted** → `P28-S06-05` / `P28-S06-06` |
| FR-P28-04 | R6 / FM-07 | **promoted** → `P28-S06-07` / `P28-S06-08` |
| FR-P28-05 | R6 / FM-08 | **promoted** → `P28-S06-09` / `P28-S06-10` |
| FR-P28-06 | R6 / FM-09 | **promoted** → `P28-S06-11` / `P28-S06-12` |
| FR-P28-07 | R6 / FM-10 | **promoted** → `P28-S06-13` / `P28-S06-14` |
| FR-P28-D1…D4, X1 | explicit defers / wontfix | **not board rows** — see S06 `SCOPE-TODOS.md` |

---

## Register (historical — actionable text preserved)

| ID | Residual | Status | Objective | Suggested scope | Acceptance hint |
|----|----------|--------|-----------|-----------------|-----------------|
| FR-P28-01 | R6 / **FM-01** | promoted | Reduce seed-import roster pin so BLOCKING discoveries expand executable backlog without waiting solely on agent memory | Product: guided promotion UX after seed import / `loop next` promotion_candidates surfacing; optional harness nudge post-import | Dogfood or integration: after seed import with orphan BLOCKING discovery, agent/path yields a new task (or explicit decline) without inventing UUIDs; document remaining human-gate if auto-spawn stays deferred |
| FR-P28-02 | R6 / **FM-02** | promoted | Close pre-export skip: agents write discoveries/decisions before `seed export --strict --enforce`, not only at gate time | Harness: stronger gap-pass / write-before-edit nudges; optional product: warn on thin graph earlier than export | Directed or build arm: consecutive session shows disc/dec writes *before* export; enforce still blocks thin export (regression green) |
| FR-P28-03 | R6 / **FM-04** | promoted | Stop worker-only Trace when parent delegates graph work to subagents without `TRACE_TASK_ID` / loop gate | Harness: parent Multitask / orchestrator rules + Option A hook already deny empty-task under strict; extend worker inheritance or parent-must-set-task docs/tests | Live or scripted: parent cannot complete edit path by offloading graph to workers while parent edits without task; document Cursor Multitask limits if product-unfixable |
| FR-P28-04 | R6 / **FM-07** | promoted | Keep git-sparsity / post-hoc SPEC commits as **warn-only** unless product adds plan-before-edit mode | Protocol: document warn semantics in VERIFY/harness; product spike only if human wants fail-closed commit ordering | Acceptance = explicit decision: remain warn-only (document) **or** ship plan-before-edit gate with tests; do not silently turn FM-07 into hard fail |
| FR-P28-05 | R6 / **FM-08** | promoted | Make agents prefer task / promotion path over discovery-only edits after `trace_add` | Product+harness: reinforce INT-06 MCP ordering + post-discovery nudge; optional apply-path smoke in dogfood | Session evidence: discovery → task (or spawned_tasks) before product edits; MCP description regression stays green |
| FR-P28-06 | R6 / **FM-09** | promoted | Prove mode collapse stays closed beyond single Session-B (build≠directed richness) | Protocol/dogfood: repeat dual-lane score (no `prepare.sh` wipe); optional second directed fixture | Dual-lane: thin build baseline documented; directed P25-3b PASS; rich build labeled post-directed — not conflated with Session-A thin FAIL |
| FR-P28-07 | R6 / **FM-10** | promoted | Ensure promotion API is used in live loops (API shipped; build-only exports still risk 0 discoveries) | Dogfood + optional apply E2E assert already in TEST-MATRIX M-01; measure live promotion rate | Directed/build run shows ≥1 discovery linked to task **or** spawned task from BLOCKING; auto-spawn without human gate stays out of scope (see FR-P28-D1) |

---

## Explicit defers (track, do not invent product scope)

| ID | Topic | Status | Rationale | Acceptance hint |
|----|-------|--------|-----------|-----------------|
| FR-P28-D1 | Autonomous discovery→task spawn (no human gate) | deferred | INTERVENTION-MATRIX §4 + Trace law: human-approved backlog expansion | Human product decision recorded; if approved, separate phase with fail-closed confirm UX — not silent swarm |
| FR-P28-D2 | Full Graphiti episode / temporal invalidation DB | deferred | INT-05 minimal reset sufficient; multi-phase spike | Spike README + migration plan only after human promote; no daemon |
| FR-P28-D3 | RESULTS.md parser for P25-4 | deferred | Env attestation (`P25_ATTEST_*`) closed R5; parser optional | Only if env attestation insufficient for CI; otherwise wontfix |
| FR-P28-D4 | Hook Option B (parent-orchestrator-detected deny only) | deferred | S03 locked Option A; Option B deferred | Human re-open only with Multitask detection design |
| FR-P28-X1 | Daemon / HTTP / hosted MCP on Trace core | **superseded (Phase 29)** | Human promoted 2026-08-21: opt-in local HTTP + browser GUI; hosted SaaS still later-developments | See Phase 29 ADR after S01 |

---

## Closed in Phase 28 S00–S05 (do not re-queue)

R1 Session-B / P25-3b · R2/R3/R8 Option A hook + INT-11 · R4 honesty single source · R5 `P25_ATTEST_*` · R7 TEST-MATRIX + VERIFY · FM-03/05/06 largely closed (monitor only; no FR row).
