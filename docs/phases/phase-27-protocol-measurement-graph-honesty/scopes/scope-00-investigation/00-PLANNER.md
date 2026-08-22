# P27-S00-00 — Scope planner (investigation)

## Metadata
- id: P27-S00-00
- todo_ids: [P27-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, ask-questions-if-underspecified, code-explorer]
- mcps: [user-codegraph]
- verification: automated
- hooks: []

## Objective

Lock investigation scope for INT-07/08/10 against live repo state and Phase 26 verify residuals. Finalize `01-investigation.md` + `02-review.md` so a fresh subagent produces `AUDIT.md` mapping P25-3 build-only failure to concrete file targets for S01/S02. No product code in this scope.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [Phase 27 README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [Phase 26 VERIFY-NOTES](../../../phase-26-loop-implementation/scopes/scope-05-verify/VERIFY-NOTES.md)
- [Phase 26 DR-HANDOFF](../../../phase-26-loop-implementation/DR-HANDOFF.md)
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Output | `scopes/scope-00-investigation/AUDIT.md` |
| Product Go | **No** on S00-01 |
| Phase 26 residual | P25-3 FAIL on build-only G1 (`discoveries=0 decisions=0`); P25-1/2 PASS; not a P25-C regression |
| INT themes | INT-07 (export `--strict`), INT-08 (protocol v2 / score.sh), INT-10 (two-session rubric) |
| Harness root | `experiments/ab-p25-gap-pass-validation/` (E02/E03 verify arm) |
| Sequence | S00 → S01 → S02 → S03 (serial default) |
| Threshold numbers | Document options only in AUDIT; S01/S02 planners pick |

## Files to audit (at minimum)

| Path | What to find |
|------|--------------|
| `experiments/ab-p25-gap-pass-validation/PROTOCOL.md` | Session modes (build-only vs directed gap); arm isolation |
| `experiments/ab-p25-gap-pass-validation/RUBRIC.md` | P25-3 pass criteria; expected build-only FAIL |
| `experiments/ab-p25-gap-pass-validation/score.sh` | Graph entity counts; `--p25` checks; export shape handling |
| `experiments/ab-p25-gap-pass-validation/prepare.sh` | Whether export is operator step vs automated |
| `experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-BUILD.md` | Build-only arm wording |
| `cmd/trace/seed.go` | `--strict` / `--enforce`; `GateForExport` wiring |
| `internal/domain/seed_export.go` | Export document shape; deliberation_states |
| `internal/loop/gate.go` | `GateForExport` violations |
| `internal/domain/seed_eval_rules_test.go` | Existing strict/export test coverage |
| `experiments/runs/2026-08-20-p26-s05-01-verify/evidence/` | Phase 26 verify score + export snippets |

## Planner gate

- [ ] `01-investigation.md` runnable (exit criteria + no-product-code rule)
- [ ] `02-review.md` includes verify checklist and spawn policy
- [ ] `SCOPE-TODOS.md` lists S00 board rows
- [ ] Live paths above still exist (adjust `01-investigation.md` if renamed)

## Exit criteria

- [ ] Audit implementer prompt locked enough for a fresh subagent
- [ ] Board row P27-S00-00 Notes cite what was verified/thickened
- [ ] Next runnable remains **P27-S00-01** (do not start S01)

## Todo updates

Status + notes on **P27-S00-00** only.

## Next

`P27-S00-01`
