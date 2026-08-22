# P28-S00-00 — Scope planner (residual audit)

## Metadata
- id: P28-S00-00
- todo_ids: [P28-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, ask-questions-if-underspecified, code-explorer]
- mcps: [user-codegraph]
- verification: automated
- hooks: []

## Objective

Lock investigation scope for Phase 28 residuals against live repo state and Phases 24–27 outcomes. Finalize `01-residual-audit.md` + `02-review.md` so a fresh subagent produces `RESIDUAL-AUDIT.md` mapping R1–R8 and INT-01..11 to closure tasks for S01–S04. **No product code** in this scope.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [Phase 28 README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [Phase 27 VERIFY-NOTES](../../../phase-27-protocol-measurement-graph-honesty/scopes/scope-03-verify/VERIFY-NOTES.md)
- [Phase 27 AUDIT.md](../../../phase-27-protocol-measurement-graph-honesty/scopes/scope-00-investigation/AUDIT.md)
- [Phase 26 VERIFY-NOTES](../../../phase-26-loop-implementation/scopes/scope-05-verify/VERIFY-NOTES.md)
- [Phase 26 PLAN.md](../../../phase-26-loop-implementation/scopes/scope-01-planning/PLAN.md)
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Output | `scopes/scope-00-residual-audit/RESIDUAL-AUDIT.md` |
| Product Go | **No** on S00-01 |
| Residual anchor | R1–R8 from Phase 28 README; P27 VERIFY residuals (P25-3b, P25-4, BLOCKING dup, hook allow path) |
| INT coverage | INT-01..11 — each row: **implemented** / **partial** / **residual** + test coverage |
| Harness root | `experiments/ab-p25-gap-pass-validation/` (E02 G1) |
| Hook gap | `CursorLoopGateHookScript()` L106–108: empty `TRACE_TASK_ID` → **allow** (FM-05 / R3) |
| Sequence | S00 → S01 → S02 → S03 → S04 → S05 (serial default; S02 may parallel S01 after audit) |
| Session-B | **Not** in S00 — audit only; dogfood is S02 |
| Threshold numbers | Document in audit; S01/S02 planners pick — do not lock here |

## Files to audit (at minimum)

| Path | What to find | INT / Residual |
|------|--------------|----------------|
| `internal/domain/promote.go`, `internal/loop/apply.go` | BLOCKING discovery → spawned task | INT-01, P25-A |
| `internal/loop/saturation.go`, migration 028 | `SaturationEmptyThreshold`, reset path | INT-02/05/09, P25-B |
| `internal/install/gappass.go`, `enforcement.go` | GapPassPrompt + ParentOrchestratorRule in install output | INT-03/04, P25-C |
| `internal/domain/seed_export_honesty.go`, `cmd/trace/seed.go` | strict/enforce; BLOCKING orphan dedupe (R4) | INT-07, R4 |
| `experiments/ab-p25-gap-pass-validation/score.sh` | `--arm build\|directed`, P25-3a/3b, P25-4 skip | INT-08/10, R5 |
| `experiments/ab-p25-gap-pass-validation/PROTOCOL.md` | two-session + attestation block | INT-10, R5 |
| `experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-DIRECTED-GAP.md` | Session-B prompt exists | R1 |
| `internal/install/enforcement.go` L100–122 | Hook allow when no TRACE_TASK_ID | INT-04, R2/R3 |
| `internal/install/cursorhook.go` | Hook JSON schema; drift doc | INT-11, R8 |
| `internal/loop/*_test.go`, `cmd/trace/*_test.go`, `internal/install/*_test.go` | Existing automated coverage | R6, R7 |
| `docs/phases/phase-27-.../VERIFY-NOTES.md` | Residual table at Phase 27 close | R1–R5 |
| `docs/phases/phase-24-.../INTERVENTION-MATRIX.md` §3 | FM coverage matrix post INT-01..11 | R6 |

## Required reads (S00-01 implementer)

1. Phase 28 README residual register (R1–R8)
2. Phase 27 VERIFY-NOTES residuals table
3. Phase 27 AUDIT.md INT mapping
4. Phase 26 VERIFY-NOTES + PLAN.md
5. INTERVENTION-MATRIX §3 FM matrix
6. Live spot-check: hook script allow path + score.sh P25-3a/3b labels

## Planner gate

- [ ] `01-residual-audit.md` runnable (metadata, preflight, audit table, template, exit criteria)
- [ ] `02-review.md` includes verify checklist, live spot-checks, spawn policy
- [ ] `SCOPE-TODOS.md` lists S00 board rows
- [ ] Live paths above still exist (adjust `01-residual-audit.md` if renamed)

## Exit criteria

- [ ] Audit implementer prompt locked enough for a fresh subagent
- [ ] Board row P28-S00-00 Notes cite what was verified/thickened
- [ ] Next runnable remains **P28-S00-01** (do not start S01)

## Todo updates

Status + notes on **P28-S00-00** only.

## Next

`P28-S00-01`
