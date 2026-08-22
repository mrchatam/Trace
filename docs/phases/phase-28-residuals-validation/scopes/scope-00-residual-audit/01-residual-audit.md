# P28-S00-01 — Residual audit implementer

## Metadata
- id: P28-S00-01
- todo_ids: [P28-S00-01]
- role: implementer
- skills: [code-explorer, spec-miner, investigator]
- mcps: [user-codegraph]
- verification: mixed
- hooks: []

## Objective

Audit live codebase and Phase 24–27 artifacts for every open residual (R1–R8) and INT-01..11 implementation vs test-coverage gap. Produce `RESIDUAL-AUDIT.md` with closure seeds for S01–S04. **No product code** on this row.

## References

- [00-PLANNER.md](00-PLANNER.md) — locked defaults + file audit table
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [Phase 28 README](../../README.md)
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
| Deliverable | `scopes/scope-00-residual-audit/RESIDUAL-AUDIT.md` |
| Product code | **No** — audit artifact only |
| Residual register | R1–R8 from Phase 28 README |
| INT themes | INT-01..11 — implementation + test row each |
| Harness root | `experiments/ab-p25-gap-pass-validation/` |
| Hook gap anchor | `internal/install/enforcement.go` `CursorLoopGateHookScript()` L106–108 allow-without-task |
| Session-B | Document R1 status only — do **not** run dogfood |

## Preflight

Run from repo root. All paths must exist (planner verified 2026-08-20, P28-S00-00):

```bash
cd /home/ali/Desktop/Trace

# Phase artifacts
test -f docs/phases/phase-27-protocol-measurement-graph-honesty/scopes/scope-03-verify/VERIFY-NOTES.md
test -f docs/phases/phase-27-protocol-measurement-graph-honesty/scopes/scope-00-investigation/AUDIT.md
test -f docs/phases/phase-26-loop-implementation/scopes/scope-05-verify/VERIFY-NOTES.md
test -f docs/phases/phase-26-loop-implementation/scopes/scope-01-planning/PLAN.md
test -f docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md
test -f docs/phases/phase-28-residuals-validation/README.md

# Harness + prompts
test -f experiments/ab-p25-gap-pass-validation/PROTOCOL.md
test -f experiments/ab-p25-gap-pass-validation/RUBRIC.md
test -f experiments/ab-p25-gap-pass-validation/score.sh
test -f experiments/ab-p25-gap-pass-validation/prepare.sh
test -f experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-BUILD.md
test -f experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-DIRECTED-GAP.md

# Product / install paths
test -f internal/domain/promote.go
test -f internal/loop/apply.go
test -f internal/loop/gate.go
test -f internal/deliberation/types.go
test -f internal/store/schema/028_deliberation_consecutive_empty.sql
test -f internal/domain/seed_export_honesty.go
test -f cmd/trace/seed.go
test -f internal/install/gappass.go
test -f internal/install/enforcement.go
test -f internal/install/cursorhook.go
```

If any path is missing, **stop** and mark row `blocked` with the missing path in Notes.

**Path note:** There is no `internal/loop/saturation.go`. Saturation policy lives in `internal/deliberation/types.go` (`SaturationEmptyThreshold`); reset/apply behavior is tested in `internal/loop/saturation_reset_test.go`.

## Live anchors (spot-check before writing)

| Topic | Path | Lines (2026-08-20) |
|-------|------|---------------------|
| Hook allow without task | `internal/install/enforcement.go` | L106–108: empty `TRACE_TASK_ID` → `permission: allow` |
| P25-3a/3b + P25-4 skip | `experiments/ab-p25-gap-pass-validation/score.sh` | L197–218 |
| BLOCKING dup (R4) | `cmd/trace/seed.go` | L158–170 store-backed BLOCKING check (may duplicate orphan link in `seed_export_honesty.go` L43–48) |
| Thin-graph honesty | `internal/domain/seed_export_honesty.go` | L21–25 min count; L43–48 orphan links |
| Saturation threshold | `internal/deliberation/types.go` | L11–13 `SaturationEmptyThreshold = 2` |
| INT-11 drift doc | `internal/install/cursorhook.go` | `HookDriftNote` const (manual upgrade checklist) |
| INT-01 promotion | `internal/domain/promote.go` | `PromoteBlockingDiscovery` |
| INT-03/04 install text | `internal/install/gappass.go`, `enforcement.go` | `GapPassPrompt`, `ParentOrchestratorRule`, `CursorLoopGateHookScript` |

## Files to audit (minimum)

| Path | What to find | INT / Residual |
|------|--------------|----------------|
| `internal/domain/promote.go`, `internal/loop/apply.go` | BLOCKING discovery → spawned task | INT-01, P25-A |
| `internal/deliberation/types.go`, `internal/loop/saturation_reset_test.go`, migration `028_*` | `SaturationEmptyThreshold`, reset path | INT-02/05/09, P25-B |
| `internal/install/gappass.go`, `enforcement.go` | Gap pass + orchestrator in install text | INT-03/04, P25-C |
| `internal/domain/seed_export_honesty.go`, `cmd/trace/seed.go` | strict/enforce; BLOCKING orphan dedupe (R4) | INT-07, R4 |
| `experiments/ab-p25-gap-pass-validation/score.sh`, `PROTOCOL.md` | `--arm build\|directed`, P25-3a/3b, P25-4 skip, attestation | INT-08/10, R5 |
| `experiments/.../prompts/PROMPT-G1-DIRECTED-GAP.md` | Session-B prompt exists | R1 |
| `internal/install/enforcement.go` L100–122 | allow when `TRACE_TASK_ID` empty | INT-04, R2/R3 |
| `internal/install/cursorhook.go` | Hook schema drift doc | INT-11, R8 |
| `internal/loop/*_test.go`, `internal/deliberation/*_test.go`, `cmd/trace/*_test.go`, `internal/install/*_test.go` | Automated vs dogfood-only coverage | R6, R7 |
| `INTERVENTION-MATRIX.md` §3 | FM-01..10 post INT-01..11 residual gaps | R6 |
| Phase 27 `VERIFY-NOTES.md` residuals table | P25-3b, P25-4, BLOCKING dup, hook | R1–R5 |

### INT-01..11 theme reference

| INT | Theme | Primary paths |
|-----|-------|---------------|
| INT-01 | BLOCKING discovery → task promotion | `promote.go`, `apply.go` spawned_tasks |
| INT-02 | P19 saturation / hop budget recalibration | `deliberation/types.go`, `loop/apply.go` |
| INT-03 | Default gap-pass prompt in install | `gappass.go`, `enforcement.go` |
| INT-04 | Parent orchestrator + hook gate | `enforcement.go` hook script |
| INT-05 | Gap-pass deliberation reset | `loop` reset API, `saturation_reset_test.go` |
| INT-06 | MCP trace_add ordering / nudge | `internal/mcp/` (if present) |
| INT-07 | seed export graph honesty | `seed_export_honesty.go`, `seed.go` strict path |
| INT-08 | Protocol v2 / score.sh enforce | `score.sh`, `PROTOCOL.md`, `prepare.sh` |
| INT-09 | Sticky STOP reason UX parity | `deliberation/select.go`, gate JSON |
| INT-10 | Two-session rubric P25-3a/3b | `score.sh`, `RUBRIC.md`, prompts |
| INT-11 | Hook drift verification | `cursorhook.go` HookDriftNote |

For each INT-01..11, record:

1. **Implementation status:** shipped / partial / missing
2. **Test coverage:** unit / integration / dogfood-only / none
3. **Residual:** open / closed / defer with rationale

## RESIDUAL-AUDIT.md template

Write to `scopes/scope-00-residual-audit/RESIDUAL-AUDIT.md`:

```markdown
# Residual audit — Phase 28

**Date:** YYYY-MM-DD  
**Git SHA:** (if available)  
**Auditor row:** P28-S00-01

## Executive summary

(3–5 sentences: what is closed from P26/P27, what remains, recommended scope order)

## Residual register disposition (R1–R8)

| ID | Status | Evidence | Target scope | Notes |
|----|--------|----------|--------------|-------|
| R1 | close in P28 / defer / already closed | … | S02 | … |
| … | … | … | … | … |

## INT-01..11 implementation matrix

| INT | Theme | Implementation | Test coverage | Residual |
|-----|-------|----------------|---------------|----------|
| INT-01 | … | shipped/partial/missing | unit/dogfood/none | … |
| … | … | … | … | … |

## Dogfood residuals

(R1 Session-B; P25-3b; arm isolation; G1 workspace state)

## Harness residuals

(R2/R3/R8 hook allow path; score attestation; FM-05)

## Product residuals

(R4 BLOCKING dup; R5 P25-4; R6 FM gaps)

## Test coverage matrix (R7)

| Area | P25 theme | Existing tests | Gap | S01 seed |
|------|-----------|----------------|-----|----------|
| Promotion | P25-A | … | … | file + case |
| … | … | … | … | … |

## Closure plan → S01..S04

| Scope | Tasks (actionable) | Files |
|-------|-------------------|-------|
| S01 | … | … |
| S02 | … | … |
| S03 | … | … |
| S04 | … | … |

## Explicit defers (with rationale)

(Items not closing in Phase 28 — must have owner + reason)

## Risks / open decisions

(For scope planners — threshold picks, parallel S01/S02, etc.)
```

## Role work

1. Run preflight bash block — all paths must PASS.
2. Read Phase 28 README residual register (R1–R8).
3. Read Phase 27 VERIFY-NOTES residuals table (L53–63) + Phase 27 AUDIT INT mapping.
4. Read Phase 26 VERIFY-NOTES + PLAN.md for INT-01..05 baseline.
5. Read INTERVENTION-MATRIX §3 FM coverage matrix — map open FM gaps to R6.
6. Spot-check live anchors table (grep/read — cite line numbers in RESIDUAL-AUDIT).
7. Inventory existing tests: `internal/loop/*_test.go` (8 files), `internal/deliberation/*_test.go`, `cmd/trace/enforce_test.go`, `internal/install/enforcement_test.go` — note dogfood-only vs automated.
8. Write `RESIDUAL-AUDIT.md` using template above.
9. Self-check: every R1–R8 row + every INT-01..11 row populated.

## Do not

- Implement fixes or change product code
- Run live Session-B dogfood (S02)
- Wipe E02 G1 with `./prepare.sh`
- Mark residuals closed without evidence citation

## Exit criteria

- [ ] `RESIDUAL-AUDIT.md` exists with all template sections
- [ ] R1–R8 each have status: **close in P28** / **defer** / **already closed**
- [ ] INT-01..11 each have implementation + test coverage row
- [ ] S01–S04 closure plan has actionable file + gap seeds
- [ ] No product code changed

## Minimal todos

- [ ] Preflight paths PASS
- [ ] Read Phase 27 VERIFY-NOTES + AUDIT
- [ ] Audit INT table with live line citations
- [ ] Write RESIDUAL-AUDIT.md
- [ ] Update board row P28-S00-01 status + notes

## Todo updates

Status + notes on **P28-S00-01** only.

## Next

**P28-S00-02**
