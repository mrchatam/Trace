# S02 — Phase VERIFY — scope todos

**Depends-on:** P15-S01-02 APPROVE. Owns Phase 15 close + DR-HANDOFF.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | VERIFY planner | **done** — FINAL locks + DR-HANDOFF = `no successor` |
| 2 | 01-verify | verifier | **done** — VERIFY PASS; VERIFY-NOTES + DR-HANDOFF started = `no successor` |
| 3 | 02-scope-review | review / handoff close | **done** — APPROVE high; DR-HANDOFF **closed** = `no successor`; Phase 15 complete |

## Depends (from S01 — post P15-S01-02 APPROVE)
VERIFY must include:

- **S01:** `TestMCPAssertDeniedBlocksCallTool`, `TestMCPAssertBuiltinAutoAllowedSucceeds`, `TestToolNamesRegistered` (+ keep `TestBuiltinMCPCapabilitySpecs` / `TestCapabilityDecision` in the S01 package run)
- **Carry-forward:** honesty A/B/C+G; Gate E/F; ablation; Gate H; compat **13**; p0x; x0; Gate C `dry_run:false`; product `./cmd|internal|evals`

## Residuals (non-blocking — do **not** fail VERIFY)
- **R2 defer:** `allowContainsOut` late-upgrade — Notes only
- **R3 wontfix:** `similar projects/graphify` space-in-path on full `./...`
- **R4 wontfix:** CGO0 analyzers FAIL (tree-sitter) — product bar is CGO1

## Reminders
- Default DR-HANDOFF = **`no successor`** — **S02-01 starts** Notes; **S02-02 owns completion**
- Do **not** auto-scaffold Phase 16 / S05 / plan simulate / D21+
- Dry-run ≠ Gate C / F / G / ablation / H / checklist
- Phase 14 historical handoff stays intact
- Spawn on fail: `P15-S02-01a` / `01b` (+`01c`) immediately below
- Next after APPROVE: **none** — DR-HANDOFF closed = `no successor`; roadmap closed again
