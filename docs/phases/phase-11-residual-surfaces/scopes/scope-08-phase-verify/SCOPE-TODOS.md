# S08 — Phase VERIFY / residual closeout — scope todos

**Depends-on:** P11-S07-02 done. Owns all **18** P11 DF regressions + carry-forward. (**S07 FINAL forward note:** DF-33 seed `from_id`/`to_id` aliases; DF-30 `phases:[]` + goal `tasks` on plan show; DF-46 snake_case plan JSON; DF-45 review get/show/list CLI; DF-28 thin help handoff SoT — no entity/mig. VERIFY evidence table covers those five plus S01–S06.)

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks + DR-HANDOFF |
| 2 | 01-verify | verify | **done** — PASS; VERIFY-NOTES; DR-HANDOFF start = `no successor` |
| 3 | 02-scope-review | review | **done** — APPROVE; DR-HANDOFF = `no successor`; Phase 11 complete |

## Locked evidence imports (P11-S08-00)

| Scope | Named regressions (must re-prove) |
|-------|-----------------------------------|
| S01 | DF-40 `TestIndexPartialArgvGCAfterRename` (+ `TestListFilePathsByContentHash`); DF-20 retain trio |
| S02 | DF-43 `TestSiblingFailBlocksDone` + PASS-alone / PASS+UNCERTAIN / hatch; DF-44 identity docs + `TestOperatorDoneRequiresFlag` |
| S03 | DF-47 `TestOpenRetrySucceedsWhenLockReleasedSoon` / exclusivity / serialize guidance / init fail-closed |
| S04 | DF-41 `TestUpsertCapabilityBySlugUpdatesExisting` + clash retain; DF-51 `TestAllowDoneDoesNotBypassMissingCaps` + WARNING tests |
| S05 | DF-49/35/48/42 symbol Exact/Why; depth-2 no sibling body; Law 9+4 MD; discovery→task |
| S06 | DF-50/22/37 tip parity print+write; `TestToolNamesRegistered` / `TestTraceVersion` |
| S07 | DF-33/30/46/45/28 seed aliases; plan show; review CLI; help handoff SoT |

## Carry-forward
Honesty A/B/C+G; Gate E; Gate F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` N=3; product `./...` (graphify space FAIL OK; CGO0 analyzers OK residual).

## Reminders
- Independent re-prove S01–S07 + carry-forward gates; close DR-HANDOFF (default **no successor**)
- Carry-forward gates stay green; Gate C `dry_run:false` untouched
- Dry-run ≠ Gate C / F / G / ablation / H / checklist
- Forward-only board; implementers: status + Notes only
- Spawn on fail: `P11-S08-01a` / `01b` / (`01c`) immediately below
- **DR-HANDOFF:** default **`no successor`** — S08-01 starts Notes; **S08-02 owns completion**
- Findings: on PASS, flip 18 DF rows closed (S08-01 or S08-02)
- Next after APPROVE: **(DR-HANDOFF close)** — no Phase 12 unless promoted
