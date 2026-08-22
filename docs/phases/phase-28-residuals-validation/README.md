# Phase 28 — Residuals closure + implementation validation

Human-promoted after Phase 27 close (`no successor`). Phase 27 shipped INT-07/08/10 but left expected residuals: build-only P25-3a FAIL, no Session-B dogfood (P25-3b unvalidated), hook failClosed beyond install text (FM-05), P25-4 attestation manual, BLOCKING duplicate honesty msg, and no consolidated regression matrix for P25-A/B/C/D/E.

## Goal

Close every open residual from Phases 24–27 and **validate shipped implementations** with automated tests plus live Session-B dogfood — without rewriting done phase history.

## Evidence basis

- [Phase 27 DR-HANDOFF](../phase-27-protocol-measurement-graph-honesty/DR-HANDOFF.md) (CLOSED; residuals listed at close)
- [Phase 27 VERIFY-NOTES](../phase-27-protocol-measurement-graph-honesty/scopes/scope-03-verify/VERIFY-NOTES.md)
- [Phase 27 AUDIT.md](../phase-27-protocol-measurement-graph-honesty/scopes/scope-00-investigation/AUDIT.md)
- [Phase 26 VERIFY-NOTES](../phase-26-loop-implementation/scopes/scope-05-verify/VERIFY-NOTES.md)
- [INTERVENTION-MATRIX.md](../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) — FM matrix §3
- E02 harness: [`experiments/ab-p25-gap-pass-validation/`](../../../experiments/ab-p25-gap-pass-validation/)

## Residual register (promotion input)

| ID | Source | Description | Target scope |
|----|--------|-------------|--------------|
| R1 | P27 VERIFY | P25-3b never validated — no Session-B dogfood | S02 |
| R2 | P26 DR-HANDOFF | INT-04 hook enforcement beyond install text | S03 |
| R3 | FM-05 matrix | Strict config without hook deny allows post-STOP edits | S03 |
| R4 | P27 S02-02 | BLOCKING duplicate message in honesty check | S04 |
| R5 | P27 VERIFY | P25-4 operator attestation manual/skip | S04 |
| R6 | FM matrix §3 | FM-01/02/07/08/10 residual gaps post INT-01..11 | S00, S01 |
| R7 | — | No consolidated regression matrix for P25 stack | S01, S05 |
| R8 | INT-11 | Hook drift verification not automated | S03 |

## Scope sequence

```
S00 → S01 → S02 → S03 → S04 → S05 (closed) → S06 R6 FM residuals → S07 residual-wave VERIFY
```

| Scope | Title | Deliverable |
|-------|-------|-------------|
| S00 | Residual audit | `RESIDUAL-AUDIT.md` |
| S01 | Integration tests | `TEST-MATRIX.md` + automated coverage |
| S02 | Session-B dogfood | `SESSION-B-NOTES.md` + P25-3b score |
| S03 | Hook failClosed | strict deny without `TRACE_TASK_ID` |
| S04 | Product polish | honesty dedupe + P25-4 attestation |
| S05 | VERIFY | `VERIFY-NOTES.md` + DR-HANDOFF S05 close |
| S06 | R6 FM residual wave | FR-P28-01…07 implement/review pairs |
| S07 | Residual-wave VERIFY | `VERIFY-NOTES-RESIDUAL-WAVE.md` + Residual wave close |

**S00→S05** serial history is `done` (immutable). **Residual wave** runs S06→S07 forward-only after `P28-S05-02`.

## In scope

- Investigation + audit artifacts under this phase folder
- Integration tests extending existing `internal/`, `cmd/`, `evals/` packages
- Hook failClosed hardening in `internal/install/`
- Product polish in honesty/export + harness attestation wiring
- Session-B dogfood on existing E02 G1 workspace (no prepare wipe)
- Full regression VERIFY at S05
- **Residual wave:** R6 FM-01/02/04/07/08/09/10 closures (harness/product/dogfood as per FR prompts)

## Out of scope

- Daemon / HTTP / hosted service on Trace core
- Rewriting Phase 24–27 `done` board history or delivered code “as if never shipped”
- Rewriting Phase 28 S00–S05 `done` prompts/Notes
- Wiping E02 G1 with `./prepare.sh` (invalidates Session-A evidence)
- Auto-spawn from discoveries without human gate (Phase 24 §4)

## Hard constraints

- `go test ./internal/...` green after product-touching scopes
- Session-B uses `--arm directed`; build arm evidence stays separate
- DR-HANDOFF S05 successor was **no successor**; Residual wave closes separately at S07-02 (**never TBD**)

## Board

[`docs/TODO/phase-28.md`](../../TODO/phase-28.md) — residual wave Next runnable: **P28-S06-01**.

## Related

- Phase 24 [INTERVENTION-MATRIX.md](../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md)
- Phase 26 [PLAN.md](../phase-26-loop-implementation/scopes/scope-01-planning/PLAN.md)
- Phase 27 [AUDIT.md](../phase-27-protocol-measurement-graph-honesty/scopes/scope-00-investigation/AUDIT.md)
