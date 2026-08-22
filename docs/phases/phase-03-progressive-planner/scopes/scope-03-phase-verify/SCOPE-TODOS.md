# Scope S03 — Phase 03 VERIFY

**Depends-on:** S01 + S02 done (`P03-S01-02`, `P03-S02-02` APPROVE high).

**Bar:** Independent Gate E mini-eval + honesty / p0x / x0 / S01–S02 surfaces + Gate C `dry_run:false` integrity. **Phase 01 dry-run ≠ Gate C pass.**

**Gate E (locked by P03-S03-00):**
| Item | Value |
|------|-------|
| Definition | Discovery propagation reduces downstream rework |
| Path / harness | `evals/replan` **`TestPlantedDiscoveryReplan`** |
| Severity | `PLAN_AFFECTING`+ only auto-replan; INFO does not |
| Churn | N=5 fail-closed + `AckReplan` / `AckAutoReplan` |
| Surfaces | S01 `internal/planner` + mig 006; S02 mig 007 + `ApplyDiscoveryReplan` |

**Locked verify commands (summary):** see [01-verify.md](01-verify.md) — Gate E replan test; planner/store/domain; honesty; full `CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./evals/replan/... ./...`.

**DR-HANDOFF:** On pass — VERIFY creates `phase-04-review-depth` + board `P04-00`; final review owns completion. On fail — spawn 01a/01b/01c; no silent handoff.

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P03-S03-00 | planner | done | 2026-08-16: locked Gate E=`evals/replan`/`TestPlantedDiscoveryReplan` + severity/churn + full VERIFY command list + Phase 04=`phase-04-review-depth`; thickened 01-verify + 02-scope-review; no product Go |
| P03-S03-01 | verify | done | 2026-08-16: **Phase 03 VERIFY PASS / Gate E green**; honesty/p0x/x0/replan/S01–S02/`./...` PASS; Gate C intact; Phase 04 scaffold + `P04-00` started — [VERIFY-NOTES.md](VERIFY-NOTES.md) |
| P03-S03-02 | review | done | 2026-08-16: APPROVE **high**; fresh Gate E + locked suite PASS vs VERIFY-NOTES; Gate C `dry_run:false` intact; dry-run≠Gate C; DR-HANDOFF complete (`phase-04-review-depth` + board); **Phase 03 complete**; next **P04-00** — [REVIEW-NOTES.md](REVIEW-NOTES.md) |

## Checklist

- [x] P03-S03-00 planner
- [x] P03-S03-01 verify
- [x] P03-S03-02 review
