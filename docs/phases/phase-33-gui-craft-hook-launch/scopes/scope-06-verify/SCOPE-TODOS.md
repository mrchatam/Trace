# Scope 06 — board map

**S06 VERIFY + handoff.** Serial: **P33-S06-00 → P33-S06-01 → P33-S06-02**.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 589 | P33-S06-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock VERIFY floor (Themes A–C; canvas capture\|waive; no successor) |
| 590 | P33-S06-01 | [01-verify.md](01-verify.md) | Verify | `VERIFY-NOTES.md` + evidence dir; Themes A–C; canvas capture\|waive |
| 591 | P33-S06-02 | [02-dr-handoff.md](02-dr-handoff.md) | Close handoff | **done** — DR-HANDOFF CLOSED; successor **no successor** ([REVIEW-NOTES.md](REVIEW-NOTES.md)) |

## Locked by S06-00 (2026-08-21)

| Lock | Value |
|------|-------|
| VERIFY artifact | `scopes/scope-06-verify/VERIFY-NOTES.md` |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p33-s06-01-verify/evidence/` |
| Close owner | **P33-S06-02** (S06-01 leaves DR-HANDOFF **OPEN**) |
| Successor lean | Default **`no successor`** |
| Themes A–C | Must tick in VERIFY-NOTES (craft + Explore hook + `trace gui` docs/launch) |
| Budgets | seeds 6≤8; getGraph 40/depth2; UI_CAP=100; expand ≤50 |
| Docs primary | `trace gui`; PATH=`go install …/cmd/trace@…` ≠ `trace install`; serve secondary |
| Canvas shot | Capture under this scope `evidence/` **or** explicit waive (S05 deferred) |
| Product code in S06 | **Forbidden** (verify/handoff notes only) |

## Carried from S05 review (PASS)

| Item | Disposition for VERIFY |
|------|------------------------|
| Docs primary = `trace gui` | Re-spot-check; already flipped |
| PATH ≠ `trace install` | Re-spot-check quickstart |
| Addr-in-use `gui\|serve` | Landed — spot-check only |
| Explore canvas screenshot | Capture under `evidence/` **or** explicit waive |
| EmptyState CTA / craft literacy | Done in S05 — no rework |
| Canvas keyboard arrow-roving | Accepted residual — out of phase |
| S04 `explore-{light,dark}.png` | Valid Theme A evidence (list-heavy OK) |

## Command floor (S06-01) — summary

1. `npm run build` (web)
2. `go test ./cmd/trace/` + `go test ./internal/httpapi/ -run FormatAddrInUse|IsAddrInUse` + `go run ./cmd/trace gui --help`
3. overviewCompose unit + `e2e/s03-depth.spec.ts`
4. Themes A–C + Laws spot-check greps
5. Skills cites on S01/S03/S04
6. Canvas capture **or** waive
7. Residuals list → handoff

## Out of this scope

- New product features (spawn remediation rows instead).
- Retokenize / compose / route / budget changes.
- Closing DR-HANDOFF from S06-01 (owner is S06-02).
- Inventing Phase 34 unless thin follow-on exception fires.
