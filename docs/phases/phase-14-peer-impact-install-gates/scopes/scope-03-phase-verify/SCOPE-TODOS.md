# S03 — Phase VERIFY — scope todos

**Depends-on:** P14-S01-02 + P14-S02-02 APPROVE.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | VERIFY planner | **done** — FINAL locks + DR-HANDOFF = `no successor` |
| 2 | 01-verify | verifier | **done** — VERIFY PASS; DR-HANDOFF started = `no successor` |
| 3 | 02-scope-review | review | **done** — APPROVE high; DR-HANDOFF closed = `no successor`; Phase 14 complete |

## Depends (from S02 — post P14-S02-02 APPROVE)
VERIFY must include:

- **S01:** `TestImpactWalkMultiSeedExcludeSeeds`, `TestImpactWalkContainsAsymmetryNoSiblings`, `TestImpactWalkIncomingImportHop`, `TestImpactWalkLoudTruncation`, `TestImpactWalkHopRiskIncreases` + `TestPlantedImpactConflictsGateFPrelim`
- **S02:** `TestInstallDetectListsCursorStable`, `TestInstallCursorUninstallIdempotent`, `TestInstallConditionalRefusesWithoutMarker`, `TestInstallConditionalWritesWithMarker`, `TestCapabilityDecision*`, existing `TestInstallCursor*`, capability ablation
- **Carry-forward:** honesty A/B/C+G; Gate E/F; ablation; Gate H; compat **13**; p0x; x0; Gate C `dry_run:false`; product `./...`

Optional spot-check only for S01 residual `allowContainsOut` late-upgrade — non-blocking. Also note: `AssertToolAllowed` is library/CLI only (not MCP dispatch) — record honesty, do not fail VERIFY for that deferred wire-up.

## Reminders
- Default DR-HANDOFF = **`no successor`** — **S03-01 starts** Notes; **S03-02 owns completion**
- Do **not** auto-scaffold Phase 15 / S05 / plan simulate / D21+
- Dry-run ≠ Gate C / F / G / ablation / H / checklist
- Phase 13 historical handoff stays intact
- Spawn on fail: `P14-S03-01a` / `01b` (+`01c`) immediately below
- Next after APPROVE: **none** — DR-HANDOFF closed = `no successor`; roadmap closed again
