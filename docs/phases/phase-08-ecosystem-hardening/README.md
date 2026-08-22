# Phase 08 — Ecosystem & hardening

## Goal

Stable **plugin APIs**, multi-agent worktrees, and production concerns (migrations, backup, auth) for adoption beyond research (`A_PROJECT_PLAN` Phase 8). Expected outcome toward **versioned APIs** + contributor analyzer plugin path — not locking bad APIs early.

**Depends on:** Phase 07 **complete** — VERIFY PASS / Gate H green; DR-HANDOFF closed on `P07-S03-02`. **Phase 08 complete** (2026-08-16) — S01–S03 done; S04 VERIFY checklist green; DR-HANDOFF closed on `P08-S04-02` = **`no successor`**. Next runnable: **none**.

## Prior phase outcomes (carry forward)

| Item | Live value |
|------|------------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| Gate H | **Green** — `evals/perf` `TestPlantedPerfLadderGateH` (smoke/~1k/~10k; measure-then-threshold; `schema-gate-h.json` v1) |
| S01 T0 + isolation | Green (`cmd/trace`) |
| S02 Go adapter | Green (`tree-sitter-go` v0.25.0); adapter-shaped switches only |
| Capability ablation | Green — TP=3/FN=0/FP=0/TN=1; P/R=1.0 |
| Gate F / G / E | Green |
| Gate C | **Go** — G1 0.800 > B0 0.000; `dry_run:false`, N=3 |
| Dry-run ≠ Gate C / ≠ Gate F / ≠ Gate G / ≠ ablation / ≠ Gate H | Phase 01 dry-run remains regression-only |
| Store mig | Through `010_capability_surface.sql` |
| 100k / 1M CI ladders | Deferred |
| GC-03/04 | Deferred |
| Daemon / HTTP / embeddings | Still forbidden as primary unless Phase 08 promotes with Notes |
| VerifiedFact / `plan simulate` | Still **out** unless explicitly promoted |

## Locked phase defaults (P08-00 — 2026-08-16)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Goal | Ecosystem & hardening (`A_PROJECT_PLAN` Phase 8) |
| Scope order | **S01** plugin APIs → **S02** multi-agent worktrees → **S03** production hardening → **S04** VERIFY + checklist |
| Validation gate | Compatibility + security checklist — **`evals/compat`** / `TestCompatibilitySecurityChecklist` + `schema-compat.json` v1 + temp `metrics-compat.json` (**S04-00 FINAL**; S04-01 creates harness) |
| Package hint | Versioned analyzer extension points over premature megastore |
| Plugin policy | Adapter-shaped + versioned contribution surface; **do not** lock universal dynamic registry early |
| Worktree policy | Safe multi-root / worktree bind + concurrent-agent fail-closed (not swarm) |
| Production policy | Migrate hygiene + backup/restore + **local** auth/binding (not cloud OAuth) |
| Successor | `A_PROJECT_PLAN` ends at Phase 8 — VERIFY records **`no successor`** unless Notes promote |
| Carry-forward bars | Gate H + honesty A/B/C + Gate G/E/F + ablation + p0x + x0 + Gate C `dry_run:false` |
| Review policy | Every scope: `00-PLANNER` → `01` → `02-review` before next scope implement |
| Full-rebuild-on-any-change | Still forbidden |

## Scope run order (locked)

| Scope | Theme | Board IDs | Folder |
|-------|--------|-----------|--------|
| S01 | Plugin APIs / analyzer contribution surface | P08-S01-00/01/02 | `scopes/scope-01-plugin-apis/` |
| S02 | Multi-agent worktrees / project bind | P08-S02-00/01/02 | `scopes/scope-02-worktrees/` |
| S03 | Production hardening (migrations / backup / auth) | P08-S03-00/01/02 | `scopes/scope-03-production-hardening/` |
| S04 | Phase VERIFY + compat/security checklist | P08-S04-00/01/02 | `scopes/scope-04-phase-verify/` |

## Out of scope (until planners promote)

- Rewriting Gate H / commercial 1M-LOC theater
- Reopening Gate C without contradicting evidence
- Declaring A1 commercially validated
- Daemon / always-on HTTP / embeddings as primary without explicit promotion
- Universal plugin megastore / dynamic .so loader theater
- Cloud SaaS auth / hosted offering shape
- Swarm orchestration frameworks

## References

- [`docs/init/A_PROJECT_PLAN.md`](../../init/A_PROJECT_PLAN.md) Phase 8
- [`docs/init/D_DECISION_REGISTER.md`](../../init/D_DECISION_REGISTER.md) DR-HANDOFF
- Phase 07 VERIFY: [`../phase-07-performance-ladder/scopes/scope-03-phase-verify/VERIFY-NOTES.md`](../phase-07-performance-ladder/scopes/scope-03-phase-verify/VERIFY-NOTES.md)
