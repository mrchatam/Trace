# P09 / S04 / 00-PLANNER — Phase 09 VERIFY / dogfood closeout

## Metadata
- id: P09-S04-00
- todo_ids: [P09-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Lock Phase 09 VERIFY commands: **DF-01 regression** + **S02/S03 spot-checks** + **carry-forward gates** + `./...`. Decide **DR-HANDOFF** from live dogfood ladder gaps (or `no successor`). No product Go.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md)
- [experiments/LADDER.md](../../../../../experiments/LADDER.md)
- Pattern: Phase 08 VERIFY [`../../../phase-08-ecosystem-hardening/scopes/scope-04-phase-verify/`](../../../phase-08-ecosystem-hardening/scopes/scope-04-phase-verify/)
- S01–S03 REVIEW-NOTES under sibling scopes
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner).

## Depends-on (S03 — light)

S03 **shipped + reviewed** (`P09-S03-02` APPROVE high): **`trace install cursor`** print / `--write` MCP config only. VERIFY spot-checks install snippet shape + DF-05 docs; **do not** require a new MCP list-tasks tool (discoverability remains CLI `trace tasks`).

## Live dogfood ladder → DR-HANDOFF decision (2026-08-16)

| Bucket | Items | Phase implication |
|--------|-------|-------------------|
| Product gaps scheduled in Phase 09 | DF-01 → S01; DF-02/04 → S02; DF-03/05 → S03 | Closed by S01–S03 APPROVE high |
| Experiment / rubric only | DF-06, DF-07, DF-13 | Stay in `experiments/` rubrics — **not** a product phase |
| Parallel dogfood (not board-blocking) | DF-08 (D08/D09/…), DF-11 (D04 tighten), DF-12 (D11 tighten) | Continue ladder batch — **not** Phase 10 unless explicitly promoted |
| Known residuals (low) | `plan_scope` ExactLookup out; scope-only review expand untested | Forward notes only — insufficient alone for a successor phase |

**DR-HANDOFF = `no successor`.** Remaining ladder work is the parallel dogfood track (`experiments/ab-*`), same posture as Phase 08 closeout + Phase 09 reopen rule: reopen only with explicit promotion + scaffold.

## Planner work
1. Lock VERIFY command set (DF-01 + S02/S03 + carry-forward + `./...`).
2. Thicken `01-verify.md` evidence table + spawn 01a/b/c + handoff start.
3. Thicken `02-scope-review.md` owns DR-HANDOFF completion (`no successor`).
4. SCOPE-TODOS + board sync.

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Phase gate | Phase 09 dogfood hardening closeout — **not** a new eval planted gate |
| Primary regression | **`TestWhyAndContextWithLinkedReview`** (`internal/retrieval`) — DF-01 |
| S02 spot-checks | `TestTasksListAfterSeed` + `TestSeedImportRelativePathAgainstC` (+ store `TestListTasks` optional) |
| S03 spot-checks | `TestInstallCursorPrintSnippet` (+ write/merge/invalid JSON tests preferred); docs README + `experiments/ab-simple/PROTOCOL.md` DF-05 |
| MCP surface | Still **six** tools — **no** list-tasks / `trace_tasks` requirement |
| Carry-forward | Honesty A/B/C + Gate G; Gate E; Gate F; capability ablation; Gate H; compat checklist; p0x 7/7; x0; Gate C `dry_run:false` N=3 |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist |
| Full bar | `CGO_ENABLED=1 go test ./... -count=1` |
| Allowed Go on VERIFY | **None** for features — re-run + evidence docs only; spawn remediation if fail |
| Spawn | On fail: `01a` implement / `01b` review (+`01c` re-VERIFY if needed) immediately below |
| DR-HANDOFF | **`no successor`** — S04-01 starts Notes; **S04-02 owns completion**. Do **not** scaffold Phase 10 unless Notes explicitly promote. |

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable (thickened)
- [x] VERIFY commands + DR-HANDOFF locked
- [x] Board Notes; next `P09-S04-01`

## Out of scope
- Running VERIFY (S04-01)
- Product Go / new MCP tools / daemon
- Scaffolding Phase 10 without explicit promotion
- Replacing parallel dogfood ladder with harness-only
