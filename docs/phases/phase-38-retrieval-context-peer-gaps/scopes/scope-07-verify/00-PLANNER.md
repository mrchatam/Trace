# P38-S07-00 — Scope planner (verify)

## Metadata
- id: P38-S07-00
- role: planner
- skills: [planning-and-task-breakdown, qa-lead]

## Objective

Lock S07 VERIFY: confirm all P38 artifacts exist, saturation was honored, no product code landed. Prepare **`01-verify.md`** + **`02-dr-handoff.md`** for S07-01/02. Output **`VERIFY-NOTES.md`** is **S07-01 only** — not this row.

## Required artifacts

| Artifact | Scope |
|----------|-------|
| INVESTIGATION-INDEX.md | S00 |
| TRACE-AUDIT.md | S01 |
| PEER-CG.md | S02 |
| PEER-UA-GF.md | S03 (incl. **Mempalace §3**) |
| GAP-REGISTRY.md | S04 (incl. **MP column**) |
| SATURATION-NOTES.md | S05 (APPROVE) |
| REMEDIATION-PLAN.md | S06 |

## Verify blocks (lock for 01)

| Block | Check |
|-------|-------|
| 0 | No product Go/TS commits in P38 scope |
| 1 | All H1–H11 addressed in GAP-REGISTRY |
| 2 | S05 saturation APPROVE on record |
| 3 | REMEDIATION-PLAN has ranked G* + reject list |
| 4 | Peer cites exist for **CG, UA, GF, MP** |
| 5 | Moat row **M-001** present |
| 6 | Successor Phase 39 documented in VERIFY-NOTES prep → DR-HANDOFF at S07-02 |

## Planner gate (S07-00)

- [x] `01-verify.md` thickened — blocks 0–6, commands, VERIFY-NOTES template
- [x] `02-dr-handoff.md` thickened — Phase 39 scaffold (G1+G3+G4), DR-HANDOFF template, REVIEW-NOTES
- [x] SCOPE-TODOS 668–670 + block table + Phase 39 checklist
- [x] Block 4 extended for **MP** (PEER-UA-GF §3 + GAP-REGISTRY MP column)

## Exit criteria

- [x] S07-01/02 prompts runnable alone
- [x] Board `P38-S07-00` → `done`

## Next

`P38-S07-01`
