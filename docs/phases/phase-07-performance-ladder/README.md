# Phase 07 — Performance ladder & language plugins

## Goal

Prove **practical scale** of Trace indexing and analysis: incremental indexing quality, ignore tiers, additional language adapters, and Gate H size ladders (`A_PROJECT_PLAN` Phase 7 / Gate H). Expected outcome toward **benchmark tables** (order-of-magnitude 10k–1M LOC) — not premature optimization theater.

**Status:** Phase 07 **complete** (2026-08-16). S01 T0/incremental **APPROVE**; S02 Go adapter **APPROVE** (`tree-sitter-go` v0.25.0); S03 VERIFY **Gate H green** (`evals/perf` `TestPlantedPerfLadderGateH`; measure-then-threshold); DR-HANDOFF closed on `P07-S03-02`. Next board row: **`P08-00`**.

## Prior phase outcomes (carry forward — live at P07-00)

| Item | Live value |
|------|------------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| Capability-selection ablation | **Green** — `evals/capability` `TestPlantedCapabilitySelectionAblation` (TP=3/FN=0/FP=0/TN=1; P/R=1.0; `schema-capability.json` v1) |
| Gate F prelim | **Green** — `evals/impact` `TestPlantedImpactConflictsGateFPrelim` |
| Gate G prelim | **Green** — `evals/honesty` `TestHonestyEscapeRateGateGPrelim` |
| Gate E | **Green** — `evals/replan` `TestPlantedDiscoveryReplan` |
| Gate C | **Go** — mean G1 0.800 > B0 0.000; `docs/verification/gate-c-x0/` (`dry_run:false`, N=3) |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H | Phase 01 dry-run remains regression-only; Gate H = planted `evals/perf` (S03-01) |
| Honesty / P0-X / X0 / replan / impact / capability | Keep green |
| Capability surface | mig `010_capability_surface`; packet required/missing; thin `trace capability` |
| GC-03/04 | Deferred unless promoted with measurement |
| Daemon / HTTP / embeddings | Still forbidden as primary |
| VerifiedFact / `plan simulate` | Still **out** unless a Phase 07 scope planner explicitly promotes with Notes |

## Live inventory vs performance gaps (locked 2026-08-16 — P07-00)

| Surface | Today (post–Phase 06) | Phase 07 need |
|---------|------------------------|---------------|
| Indexing | `cmd/trace` `walkIndexable`: skip `.git`/`.trace`; `DetectLanguage` filter; best-effort `git check-ignore`; file-local `IndexFile`/`IndexFileAtRev` + SHA-256 upsert; CLI `TestIndexIncrementalIsolation` | **Ignore tiers** (STORAGE T0+) beyond gitignore; measurable incremental quality / hot-path latency |
| Analyzers | tree-sitter **JS/TS/TSX/Python + Go** (`LangGo`; `tree-sitter-go` v0.25.0) | S02 shipped; Gate H may optionally plant tiny `.go` |
| Perf evidence | **No** `evals/perf` yet (S01/S02 left none) | **S03-01 creates** planted smoke/~1k/~10k under **`evals/perf`**; measure-then-threshold; 100k/1M deferred |
| Store | SQLite `.trace/trace.db`; mig through **`010`** | Additive mig only if needed (`011_*` hint); measure SQLite ceiling (A5) |
| Full-rebuild | Forbidden architecture | Keep file-local / incremental — no full-rebuild-on-any-change |

## Locked phase defaults (P07-00 — 2026-08-16)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Goal | Performance ladder & language plugins (`A_PROJECT_PLAN` Phase 7 / Gate H) |
| Scope order | **S01** incremental indexing / ignore tiers → **S02** language plugins → **S03** VERIFY + Gate H + Phase 08 handoff |
| Validation gate | Gate H **FINAL (S03-00):** **`evals/perf`** / **`TestPlantedPerfLadderGateH`** + **`schema-gate-h.json` v1** + temp **`metrics-gate-h.json`**; **S03-01 creates** harness; thresholds **after first measurements** (measure-then-threshold) |
| Package hint | Prefer extend `internal/analyzers` + `cmd/trace` walk/index + optional `evals/perf` — avoid premature infra rewrite / graph DB |
| Ignore-tier hint | Align `STORAGE_AND_PERFORMANCE` T0–T3; S01 owns **T0 always-skip** — **S01-00 locked**: dirs `node_modules`/`vendor`/`__pycache__`/`.venv`/`venv`/`dist`/`.next`/`target`/`coverage` (+`.git`/`.trace`); files `.min.js`/`.min.mjs`/`.min.cjs` + path-segment; no `011_*`; T1–T3 notes only |
| Language (S02) | **FINAL:** **Go** + `github.com/tree-sitter/tree-sitter-go` **v0.25.0**; adapter-shaped `DetectLanguage`/`extract` extension; golden required; no Gate H invent |
| Migration hint | Prefer path-filter / analyzer extension without mig; additive **`011_*`** only if store needs it (S01/S02-00 lock) |
| Phase 08 folder | **`phase-08-ecosystem-hardening`** (`A_PROJECT_PLAN` Phase 8 — Ecosystem & hardening) |
| Review policy | Every scope: `00-PLANNER` → `01` → `02-review` before next scope implement |
| Carry-forward bars | Honesty A/B/C; Gate G; Gate E; Gate F; capability ablation; p0x 7/7; x0; Gate C `dry_run:false` intact |
| Daemon / HTTP / embeddings | Forbidden as primary |
| VerifiedFact / `plan simulate` | Out unless explicitly promoted with Notes |
| Perf policy | **Measure first** — no premature optimization theater; no commercial multi-model perf theater |
| DR-HANDOFF | Closing VERIFY scaffolds Phase 08 = **`phase-08-ecosystem-hardening`** (or records `no successor`) |

## Scope run order (confirmed)

| Scope | Theme | Board IDs | Folder |
|-------|--------|-----------|--------|
| S01 | Incremental indexing / ignore tiers | P07-S01-00/01/02 | `scopes/scope-01-incremental-indexing/` |
| S02 | Language plugins / adapters | P07-S02-00/01/02 | `scopes/scope-02-language-plugins/` |
| S03 | Phase VERIFY + Gate H + Phase 08 handoff | P07-S03-00/01/02 | `scopes/scope-03-phase-verify/` |

## Out of scope (until planners promote)

- Premature optimization theater / rewriting core for unmeasured gains
- Daemon / always-on HTTP / embeddings as primary
- Reopening Gate C / inventing Gate F/G/ablation without named harnesses
- Declaring Gate H **pass** before thresholds exist
- Commercial A1 / multi-model theater
- `plan simulate` / VerifiedFact promotion engine
- Starting Phase 08 before S03 VERIFY
- Deep-finalizing every implement prompt in phase planner (scope planners own that)

## References

- [`docs/init/A_PROJECT_PLAN.md`](../../init/A_PROJECT_PLAN.md) Phase 7 / Phase 8
- [`docs/STORAGE_AND_PERFORMANCE.md`](../../STORAGE_AND_PERFORMANCE.md)
- [`docs/EVALUATION.md`](../../EVALUATION.md) Gate H
- [`docs/init/D_DECISION_REGISTER.md`](../../init/D_DECISION_REGISTER.md) DR-HANDOFF
- [`docs/init/E_ASSUMPTION_REGISTER.md`](../../init/E_ASSUMPTION_REGISTER.md) A5 / A16
- Phase 06 VERIFY: [`../phase-06-environment-capability/scopes/scope-03-phase-verify/VERIFY-NOTES.md`](../phase-06-environment-capability/scopes/scope-03-phase-verify/VERIFY-NOTES.md)
