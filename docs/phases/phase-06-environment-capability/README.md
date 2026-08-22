# Phase 06 — Environment / capability graph

## Goal

Model **skills / rules / MCP / tool selection** so task packets include only required capabilities, with missing-capability warnings (`A_PROJECT_PLAN` Phase 6 / H7; `docs/AGENT_ENVIRONMENT.md`). Validation toward a **capability-selection ablation** under `evals/` — not ontology bloat or commercial multi-model theater.

**Status:** **complete** (2026-08-16) — S01 capability surface + S02 capability-selection ablation + S03 VERIFY; DR-HANDOFF closed on `P06-S03-02`. Next phase: **`P07-00`** (`phase-07-performance-ladder`).

**Depends on:** Phase 05 decision impact **complete** (Gate F prelim green; DR-HANDOFF closed on `P05-S03-02`).

## Prior phase outcomes (carry forward — live at P06-00)

| Item | Live value |
|------|------------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| Gate F prelim | **Green** — `evals/impact` `TestPlantedImpactConflictsGateFPrelim` (TP=3/FN=0/FP=0/TN=1; P/R=1.0; `schema-gate-f.json` v1) |
| Gate G prelim | **Green** — `evals/honesty` `TestHonestyEscapeRateGateGPrelim` (escapes=1/caught=2/attempts=3) |
| Gate E | **Green** — `evals/replan` `TestPlantedDiscoveryReplan` |
| Gate C | **Go** — mean G1 0.800 > B0 0.000; `docs/verification/gate-c-x0/` (`dry_run:false`, N=3) |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G | Phase 01 dry-run remains regression-only |
| Honesty / P0-X / X0 / replan / impact | Keep green |
| Decision impact | mig `009_decision_impact`; `ImpactReport`; CLI `trace impact` |
| GC-03/04 | Deferred unless promoted with measurement |
| Daemon / HTTP / embeddings | Still forbidden as primary |
| VerifiedFact | Still **out** unless a Phase 06 scope planner explicitly promotes with Notes |
| `plan simulate` | Still out (P13) unless promoted |

## Live inventory vs capability gaps (P06-00)

| Surface | Today (post–Phase 05) | Phase 06 need |
|---------|------------------------|---------------|
| MCP (`internal/mcp` + `cmd/trace-mcp`) | Six tools: `trace_why` / `trace_context` / `trace_add` / `trace_link` / `trace_transition` / `trace_review` (stdio go-sdk; G19) | Treat as **capability ids** in inventory — not a new daemon/HTTP surface |
| CLI (`cmd/trace`) | `init`/`index`/`reindex`/`add`/`link`/`transition`/`review`/`impact`/`plan`/`seed`/`why`/`context` | Thin declare/list capability cmds (S01-00 finalizes names) |
| Task packets (`internal/compiler`) | Layer 0–1 `Packet` — items/why/budget only; **no** required-capabilities / missing warnings | Attach required + missing capability fields (S01; selection in S02) |
| Skills / rules / hooks metadata | Absent as first-class graph entities | Minimal capability catalog — reject ontology megastore |
| Store migrations | `001`…`009_decision_impact` | Additive **`010_capability_surface.sql`** (locked S01-00) |
| Ablation harness | No `evals/capability` (pre-S02) | **FINAL S02-00:** `evals/capability` / `TestPlantedCapabilitySelectionAblation` + `schema-capability.json` v1 (P/R TP=3/FN=0/FP=0/TN=1) |

## Locked phase defaults (P06-00 — 2026-08-16)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Goal | Environment/capability graph (`A_PROJECT_PLAN` Phase 6 / H7) |
| Scope order | **S01** capability surface → **S02** selection / ablation → **S03** VERIFY + Phase 07 handoff |
| Validation gate | Capability-selection ablation — **`evals/capability`** / **`TestPlantedCapabilitySelectionAblation`** + `schema-capability.json` v1 + temp `metrics-capability.json` (TP=3/FN=0/FP=0/TN=1; P/R=1.0) — **FINAL P06-S02-00** |
| Package / mig hint | Prefer extend **`internal/domain`** + store + **`internal/compiler`** packet attach; optional thin MCP id mirror — **avoid** ontology megastore / second stack |
| Migration hint | Next additive **`010_*`** (S01-00 locks filename) |
| CLI vs MCP | CLI primary for mutations; MCP remains thin G19 adapter — capabilities modeled as ids, not new write tools required in S01 |
| Ablation policy | Planted eval only — **no** commercial multi-model capability theater |
| Phase 07 folder | **`phase-07-performance-ladder`** (`A_PROJECT_PLAN` Phase 7 — Performance ladder & language plugins) |
| Review policy | Every scope: `00-PLANNER` → `01` → `02-review` before next scope implement |
| Carry-forward bars | Honesty A/B/C; Gate G; Gate E; Gate F; p0x 7/7; x0; Gate C `dry_run:false` intact |
| Daemon / HTTP / embeddings | Forbidden as primary |
| VerifiedFact / `plan simulate` | Out unless explicitly promoted with Notes |
| DR-HANDOFF | Closing VERIFY scaffolds Phase 07 = **`phase-07-performance-ladder`** (or records `no successor`) |

## Scope run order (locked)

| Scope | Theme | Board IDs | Folder |
|-------|--------|-----------|--------|
| S01 | Capability / env surface | P06-S01-00/01/02 | `scopes/scope-01-capability-surface/` |
| S02 | Selection + ablation harness | P06-S02-00/01/02 | `scopes/scope-02-capability-selection/` |
| S03 | Phase VERIFY + Phase 07 handoff | P06-S03-00/01/02 | `scopes/scope-03-phase-verify/` |

## Out of scope (until planners promote)

- Ontology bloat / universal skill taxonomies without ablation evidence
- Daemon / always-on HTTP / embeddings as primary
- Reopening Gate C / inventing Gate F/G without named harnesses
- Commercial A1 / multi-model capability theater
- `plan simulate` / VerifiedFact promotion engine
- Starting Phase 07 before S03 VERIFY
- Deep-finalizing every implement prompt in `P06-00` (scope planners own that)

## References

- [`docs/init/A_PROJECT_PLAN.md`](../../init/A_PROJECT_PLAN.md) Phase 6 / Phase 7
- [`docs/AGENT_ENVIRONMENT.md`](../../AGENT_ENVIRONMENT.md)
- [`docs/PROJECT_MODEL.md`](../../PROJECT_MODEL.md) §10 Task capability requirements
- [`docs/EVALUATION.md`](../../EVALUATION.md) H7
- [`docs/init/D_DECISION_REGISTER.md`](../../init/D_DECISION_REGISTER.md) DR-HANDOFF
- Phase 05 VERIFY: [`../phase-05-decision-impact/scopes/scope-03-phase-verify/VERIFY-NOTES.md`](../phase-05-decision-impact/scopes/scope-03-phase-verify/VERIFY-NOTES.md)
