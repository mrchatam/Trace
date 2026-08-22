# P16-S06-00 — Phase VERIFY (stub — thicken vs live)

## Metadata
- id: P16-S06-00
- todo_ids: [P16-S06-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock Phase 16 VERIFY evidence: **S01–S05 named DF regressions** + **carry-forward gates** + product pkgs. Decide **DR-HANDOFF** = **`no successor`** (default). **No product Go.**

**Depends-on:** S01–S05 APPROVE (thicken named-test list from REVIEW-NOTES).

## Inherited locks
- Import every S01–S05 named test (DF-76, 75/77/78, 68, 70/73, 71/72/74)
- Carry-forward: honesty A/B/C+G; E/F; ablation; H; compat (**14** if S02 mig landed else **13**); p0x; x0; product `./cmd|internal|evals`
- Residuals **non-fail:** DF-67 defer; R2 defer; R3/R4 wontfix; DF-22/37 manual reload
- Dry-run ≠ Gate C / F / G / ablation / H / checklist
- MCP: ten tools + `trace_version` after S05; **no** install/decide MCP
- DR-HANDOFF default **`no successor`** — S06-01 starts Notes; S06-02 owns completion
- Forbidden: product features on VERIFY; claiming DF-67/R2 fixed; inventing Phase 17 / S05 supersession / plan simulate / D21+ without promotion

## Planner work
1. [ ] Import S01–S05 APPROVE named tests into 01-verify
2. [ ] Lock verify cmds + DR-HANDOFF = `no successor`
3. [ ] Thicken 01/02/SCOPE-TODOS; **FINAL**; next **P16-S06-01**
