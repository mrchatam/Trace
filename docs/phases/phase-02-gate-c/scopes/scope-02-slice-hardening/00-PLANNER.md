# P02 / S02 / 00-PLANNER — Slice hardening

## Metadata
- id: P02-S02-00
- todo_ids: [P02-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize hardening work from **S01 Gate C issue list**. Bound fixes to measurement-driven gaps; do not feature-factory. No product code in planner row.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- S01 GATE-C / REVIEW notes (after S01-02)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner).

## Depends-on
- **Hard:** S01 review done with issue list (or explicit empty list + No-Go skip path).
- **Input artifact:** `../scope-01-x0-gate-c/GATE-C-NOTES.md` issue list (entries shaped `GC-NN` — severity, metric, evidence, proposed_fix_surface, defer). Post–S01-02: GC-01/GC-02 actionable; GC-03/GC-04 deferred. Gate C verdict **Go** (see REVIEW-NOTES).
- Import concrete issues only — no speculative roadmap items from Phase 03+.
- After Gate C **No-Go** with empty/defer-all list: this scope may be `skipped` with reason (phase README).

## Bars (do not weaken)
| Item | Value |
|------|-------|
| Honesty | `evals/honesty` Paths A/B/C stay green |
| P0-X | `evals/p0x` 7/7 regression-keep |
| X0 | Keep dry-run + Gate C artifacts coherent per S01 locks (`evals/x0`, `docs/verification/gate-c-x0/`) |
| Scope | No daemon/HTTP/embeddings; no progressive planner product; no reopening unfair scoring debate — fix surfaces named in GC-NN only |

## Planner work
- Import S01 issue list; prioritize falsification-critical gaps.
- Thicken `01-slice-hardening.md` with ordered fix set + deferrals.
- Light note on S03 if VERIFY must re-check specific hardened surfaces.

## Exit criteria
- [x] `01-slice-hardening.md` runnable alone from S01 issue list
- [x] Honesty + p0x bars not weakened
- [x] Board + SCOPE-TODOS synced

## Minimal todos
- [x] Import S01 issue list (GC-01/02 implement; GC-03/04 defer)
- [x] Thicken 01 + 02 review focus
- [x] Board update
