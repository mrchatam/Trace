# Scope 03 — board map

**S03 VERIFY + handoff** — live feet-seller + close. Serial: **P35-S03-00 → P35-S03-01 → P35-S03-02**.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 619 | P35-S03-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock VERIFY floor + evidence paths (**done** — thickened 01/02) |
| 620 | P35-S03-01 | [01-verify.md](01-verify.md) | Verify | Blocks 0–5 → `VERIFY-NOTES.md` + `experiments/runs/…-p35-s03-01-verify/evidence/`; live ≠ Step1; URL override evidence |
| 621 | P35-S03-02 | [02-dr-handoff.md](02-dr-handoff.md) | Reviewer | Close `DR-HANDOFF.md`; successor lean default **no successor**; TODO/AGENTS |

## Locked floor (S03-00)

- Unit: `cd web && node --experimental-strip-types --test src/lib/pickActiveTask.test.ts`
- Live: `trace gui -C "/home/ali/Desktop/feet seller telegram app"` — Overview/Loop ≠ `33247e2d-…`; prefer `99d8fb92-…`
- Residuals to document (not fail alone): display vs pick truncation; `listTasksForPick` max-pages; HTTP pagination deferred
- DR close owner = **S03-02** only

Evidence under `experiments/runs/` when recording live runs (follow prior phase VERIFY norms).
