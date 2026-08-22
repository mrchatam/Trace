# Scope S03 — Phase 05 VERIFY

**Depends-on:** S01 + S02 APPROVE (**P05-S02-02 APPROVE high** — 2026-08-16).

**DR-HANDOFF:** On pass — VERIFY starts next-phase scaffold **`phase-06-environment-capability`** + board `P06-00`; S03-02 owns completion (or records `no successor`).

**Gate F re-prove (S03-00 locked):**
`CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim`
(schema `evals/impact/schema-gate-f.json` v1; temp `metrics-gate-f.json`; tallies TP=3/FN=0/FP=0/TN=1; precision=1.0; recall=1.0)

**Carry-forward (locked in `01-verify.md`):** honesty A/B/C; Gate G; Gate E; p0x 7/7; x0; domain/store/planner; Gate C `dry_run:false`; full `./...` with `evals/impact`.

**Spawn on fail:** `P05-S03-01a` / `01b` / `01c` (forward-only).

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P05-S03-00 | planner | done | 2026-08-16: locked Gate F re-prove + carry-forward command table + spawn 01a/b/c + Phase 06 checklist; thickened 01/02; no product Go |
| P05-S03-01 | verify | done | 2026-08-16: **VERIFY PASS / Gate F prelim green**; Phase 06 scaffold started; see [VERIFY-NOTES.md](VERIFY-NOTES.md); pending S03-02 |
| P05-S03-02 | review | pending | Owns DR-HANDOFF completion |

## Checklist

- [x] P05-S03-00 planner
- [x] P05-S03-01 verify
- [ ] P05-S03-02 review
