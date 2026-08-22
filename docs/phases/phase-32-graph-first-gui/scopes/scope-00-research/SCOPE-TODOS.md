# Scope 00 — board map

**S00 research** — peer + gap research only. Serial: **P32-S00-00 → P32-S00-01 → P32-S00-02**. Primary artifact: `RESEARCH.md` (written in **S00-01**, reviewed in **S00-02**). Must include **P32-PORT** note + S02 recommendation (prefer #1 friendly bind error). Do **not** start S01 until S00-02 PASS. Do **not** write product code in this scope.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 549 | P32-S00-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 + this file |
| 550 | P32-S00-01 | [01-research.md](01-research.md) | Implementer | Author `RESEARCH.md` (template in 01) |
| 551 | P32-S00-02 | [02-review.md](02-review.md) | Reviewer | Checklist vs DESIGN-LOCKS + P32-PORT |

## Planner-locked live facts (for 01)

- GUI: Overview is index; Graph is `/graph` route; no Graph-side inspector.
- Graph budgets: `DEFAULT_MAX=50`, `UI_CAP=100` in `web/src/screens/Graph.tsx`.
- OpenAPI: why / context / impact / graph / search / reviews present.
- Client: `ops.ts` lacks `getImpact` wrapper (flag in API map).
- Serve: `127.0.0.1:7432`; fail on conflict — P32-PORT → **S02 ships** (even with `NO-GAPS.md`).
- Peers: `similar projects/graphify/`, `similar projects/Understand-Anything/`.

## Out of this scope

- Writing `UX-IA.md` (S01), API/serve implementation (S02), depth UI (S03), visual craft (S04).
