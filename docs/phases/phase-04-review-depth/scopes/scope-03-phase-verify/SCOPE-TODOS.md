# Scope S03 — Phase 04 VERIFY

**Depends-on:** S01 + S02 done (`P04-S01-02`, `P04-S02-02` APPROVE high).

**Bar:** Independent Gate G prelim + honesty A/B/C / p0x / x0 / replan / Gate E / S01 surfaces / Gate C `dry_run:false` integrity. **Phase 01 dry-run ≠ Gate C pass.**

**Gate G (locked by P04-S03-00):**
| Item | Value |
|------|-------|
| Definition | Review reduces false completion (H5 / Gate G prelim) |
| Path / harness | `evals/honesty` **`TestHonestyEscapeRateGateGPrelim`** |
| Schema / metrics | Committed **`schema-gate-g.json`** v1; temp **`metrics-gate-g.json`** (schema-valid) |
| Tallies | escapes=1 / caught=2 / attempts=3 (escape_rate≈1/3); hatch=escape only |
| S01 hooks | `LinkReviewScope` / `review_judges_scope`; OPEN **`POLICY_EXCEPTION`**; **`CountOpenResidualsByScope`** |
| Keep | Paths A/B/C (`TestHonestyFailClosedPlantedClaim`); Gate E; Gate C `dry_run:false` |

**Locked verify commands (summary):** see [01-verify.md](01-verify.md) — Gate G named test; honesty package; Gate E replan; domain/store/planner; full `CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./...`.

**DR-HANDOFF:** On pass — VERIFY creates **`phase-05-decision-impact`** + board `P05-00`; final review owns completion. On fail — spawn 01a/01b/01c; no silent handoff.

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P04-S03-00 | planner | done | 2026-08-16: locked Gate G=`evals/honesty`/`TestHonestyEscapeRateGateGPrelim` + schema/metrics + tallies + S01 hooks + full VERIFY command list + Phase 05=`phase-05-decision-impact`; thickened 01-verify + 02-scope-review; no product Go |
| P04-S03-01 | verify | done | 2026-08-16: **VERIFY PASS / Gate G prelim green**; evidence in [VERIFY-NOTES.md](VERIFY-NOTES.md); Phase 05 scaffold + `P05-00` started; pending P04-S03-02 |
| P04-S03-02 | review | pending | Owns DR-HANDOFF completion |

## Checklist

- [x] P04-S03-00 planner
- [x] P04-S03-01 verify
- [ ] P04-S03-02 review
