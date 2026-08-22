# S03 — Provenance schema — scope todos

**Depends-on:** P13-S02-02 APPROVE. Owns **DF-64, DF-66, DF-67**.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** (FINAL) |
| 2 | 01-provenance-schema | implement | **done** |
| 3 | 02-scope-review | review | **done** — APPROVE high; [REVIEW-NOTES.md](./REVIEW-NOTES.md) |

## Depends (from S02 FINAL — light)
- S02 keeps packet SchemaVersion **`0.2`** with additive honesty/budget fields. S03 must **not** invent a packet schema bump for enum work; DF-64 lives on **store/import provenance** write validation + mig **012** CHECK.

## FINAL disposition (planner) — closed by S03-02

| DF | Disposition |
|----|-------------|
| DF-64 | **Fix shipped** — write reject garbage; empty→EXTRACTED; read normalize; mig 012 CHECK; compat ceiling 12 |
| DF-66 | **Documented wontfix** — no analyzer/CLI INFERRED setter; Law 5 store-fixture + `ANALYZER_CONTRIBUTION.md` paragraph |
| DF-67 | **Out-of-bar residual** — file-hash honesty only; no symbol-entity staleness in S03; VERIFY must Note |

## Reminders
- DF-66/67 residuals carried to **P13-S04** VERIFY-NOTES (not silently dropped)
- Next after APPROVE: **P13-S04-00**
