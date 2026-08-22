# Scope 07 — Residual-wave VERIFY (locked — S07-00)

**Planner:** P28-S07-00 (2026-08-20)  
**Prerequisite:** P28-S06-02…14 all `done` with APPROVE.

## Board rows

| ID | Role | Prompt | Deliverable |
|----|------|--------|-------------|
| P28-S07-00 | planner | `00-PLANNER.md` | Lock floor + thicken siblings |
| P28-S07-01 | verify | `01-verify.md` | `VERIFY-NOTES-RESIDUAL-WAVE.md` + evidence dir |
| P28-S07-02 | DR-HANDOFF | `02-dr-handoff.md` | Close Residual wave only; REVIEW-NOTES-RESIDUAL-WAVE.md |

## Verify floor (FINAL)

| Block | Check |
|-------|-------|
| Build | `CGO_ENABLED=1 go build -o bin/trace ./cmd/trace` |
| Unit | `GOPROXY=direct go test ./internal/... -count=1` |
| Cmd | `GOPROXY=direct go test ./cmd/trace/... -count=1` |
| Install + hook | `go test ./internal/install/...` + `CursorLoopGateFailClosed\|HookDrift\|CursorLoopGateAllowNonStrict` |
| FR spot-check | FM01/02/04/08/09/10-NOTES + FM07-DECISION vs board APPROVE |
| Optional score | Directed only; **no** `prepare.sh`; or cite FM09 archive |

Evidence: `experiments/runs/YYYY-MM-DD-p28-s07-01-verify/evidence/`

## Dual-lane

- Thin = SESSION-A snapshot / FM09 thin (disc=0/dec=0)
- Rich directed / post-directed build = FM09 (or optional re-score)
- Never wipe G1; never conflate thin FAIL with rich PASS

## FR evidence map

| FR | FM | Artifact |
|----|-----|----------|
| FR-P28-01 | FM-01 | `../scope-06-r6-fm-residuals/FM01-NOTES.md` |
| FR-P28-02 | FM-02 | `FM02-NOTES.md` |
| FR-P28-03 | FM-04 | `FM04-NOTES.md` |
| FR-P28-04 | FM-07 | `FM07-DECISION.md` |
| FR-P28-05 | FM-08 | `FM08-NOTES.md` |
| FR-P28-06 | FM-09 | `FM09-NOTES.md` |
| FR-P28-07 | FM-10 | `FM10-NOTES.md` |

Deferred (not board / non-blocking): D1–D4, X1.

## DR-HANDOFF / successor

| Item | Locked |
|------|--------|
| Close scope | Residual wave section only |
| S05 history | Immutable CLOSED |
| Default successor | **`no successor`** |
| Phase 29 | Human promote only — never TBD |

## After S07

Phase 28 residual wave complete when S07-02 APPROVE closes Residual wave.
