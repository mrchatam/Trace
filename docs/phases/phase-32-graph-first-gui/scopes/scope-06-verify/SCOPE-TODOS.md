# Scope 06 — board map

VERIFY + DR-HANDOFF. Serial: **S06-00 → S06-01 → S06-02**. Close owner **P32-S06-02**.

| Board ID | Row | Prompt | Role |
|----------|-----|--------|------|
| 567 | P32-S06-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner — lock VERIFY + handoff prompts |
| 568 | P32-S06-01 | [01-verify.md](01-verify.md) | Verifier → evidence + `VERIFY-NOTES.md` (DR-HANDOFF stays OPEN) |
| 569 | P32-S06-02 | [02-dr-handoff.md](02-dr-handoff.md) | Close DR-HANDOFF; successor **never TBD** (default **no successor**) |

## Locks (S06-00)

| Item | Value |
|------|-------|
| VERIFY artifact | `VERIFY-NOTES.md` + `experiments/runs/YYYY-MM-DD-p32-s06-01-verify/evidence/` |
| Floor | web build; P32-PORT Go tests; e2e s03+s05; DESIGN-LOCKS/Laws spot-check; port docs + loopback |
| P32-PORT | Must tick: **#1 shipped** + S05 docs; **#2 deferred** |
| Aggregate | Cite S00–S05 artifacts/Notes in VERIFY-NOTES |
| Successor lean | Default **`no successor`**; thin follow-on only if blocker residuals / human promote — then scaffold per protocol |

## Aggregate pointers (for S06-01)

| Scope | Artifact / evidence |
|-------|---------------------|
| S00 | `RESEARCH.md` — peer bar; PORT prefer #1 |
| S01 | `UX-IA.md` — hybrid C; inspector map; budgets |
| S02 | `NO-GAPS.md`; `getImpact`; PORT #1 helpers + serve tests |
| S03 | Graph home + `Inspector.tsx`; e2e depth |
| S04 | Craft A/B/C; no Three.js |
| S05 | `gui-quickstart` Multi-project / ports; `web/README` |

## Non-blocking residuals (do not fail VERIFY alone)

- PORT #2 auto-port deferred
- Serve “listening on” before bind fail (S02 low)
- Chrome box-shadow transition unused
- Keyboard select via list OK
- No screenshots / media pipeline

## Do not

- Start S06-01 from planner row / start S06-02 from verify row
- Product code on planner or verify rows
- Close DR-HANDOFF on S06-01
- Leave successor `TBD` on S06-02
