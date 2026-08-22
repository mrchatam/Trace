# S04 — Phase VERIFY — scope todos

**Depends-on:** P17-S03-02 APPROVE. Owns Phase 17 close + DR-HANDOFF.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | VERIFY planner | **FINAL** — next P17-S04-01 |
| 2 | 01-verify | verifier | pending — stop if 00 DRAFT |
| 3 | 02-scope-review | review / handoff close | pending — owns DR-HANDOFF completion |

## Locked depends (S01–S03 — imported at S04-00 FINAL)

- **S01 DF-80/84/85** `TestSeedExportRoundTrip` (plan tree) + `TestSeedExportOmitsDeniedSurfaces` + `TestSeedExportWritesExportedAtCommit` + P16 import keepers
- **S02 DF-82/85** `TestHelpSeedExportPath` + `TestHelpHandoffSoT` + `TestAsOperatorFlagIdentityDocs`; path `trace/graph.json`; `.gitignore` `.trace/` only
- **S03 DF-81/83/84** `TestSeedImportIdempotent` + `TestSeedImportDuplicateLinksNoOp` + `TestSeedImportSameIdLastWins` + `TestSeedImportPlanTreeIdempotent`
- **Two-clone recipe** `TestPortableGraphTwoCloneWhyContextPlan` — two dirs, no shared `.trace/`; init + import + index + why + context + plan; offline
- **Carry-forward:** honesty A/B/C+G; E/F; ablation; H; compat; p0x; x0; Gate C `dry_run:false`; product `./cmd\|internal\|evals`

## Residuals (non-blocking — do **not** fail VERIFY)

- Encryption-as-git **wontfix**
- Reviews/DONE omitted from default export
- **DF-86** `trace install git-hook` absent (CONDITIONAL / deferred; not wrapping `git commit`)
- CGO=0 `cmd/trace` tree-sitter FAIL (CGO=1 authoritative)
- S03 `work_state` preservation SQL-only (no dedicated named test)
- P16 DFs are P16’s (do not re-claim)
- Hosted authenticated MCP (TODO Later developments — not this phase)

## Evidence

| Artifact | Location |
|----------|----------|
| VERIFY notes | `scopes/scope-04-phase-verify/VERIFY-NOTES.md` (S04-01 writes; S04-02 confirms) |
| Gate C inspect | `docs/verification/gate-c-x0/` (carry-forward; do not re-score) |
| Scope review | `scopes/scope-04-phase-verify/REVIEW-NOTES.md` (S04-02) |

## Reminders

- Default DR-HANDOFF = **`no successor`** — **S04-01 starts**; **S04-02 owns completion**
- Do **not** auto-scaffold Phase 18 / research S05 / plan simulate / D21+ / hosted server
- Dry-run ≠ Gate C / F / G / ablation / H / checklist
- Phase 16 historical/default handoff stays intact
- Spawn on fail: `P17-S04-01a` / `01b` (+`01c`) immediately below
- SoT: [DF-84-FORWARD.md](../../DF-84-FORWARD.md)
