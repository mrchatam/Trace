# Scope 00 — board map

**S00 research** — peers + Laws 6–7 overview + launch/PATH only. Serial: **P33-S00-00 → P33-S00-01 → P33-S00-02**. Primary artifact: `RESEARCH.md` (written in **S00-01**, reviewed in **S00-02**). Do **not** start S01 until S00-02 PASS. Do **not** write product code. Planner (**S00-00**) does **not** author `RESEARCH.md`.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 571 | P33-S00-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 + this file |
| 572 | P33-S00-01 | [01-research.md](01-research.md) | Implementer | Author `RESEARCH.md` |
| 573 | P33-S00-02 | [02-review.md](02-review.md) | Reviewer | Checklist vs DESIGN-LOCKS + INTAKE + peer cites |

## Planner-locked live facts (for 01)

- Explore: `App.tsx` index → `Graph.tsx`; nav “Explore” (`/`); still center-first (`DEFAULT_MAX=50`, `UI_CAP=100`).
- Explore ≠ Overview: nav `/overview` is a separate ops screen — Theme B upgrades `/` Graph, not Overview.
- CLI: `serve` exists; **no** `gui` subcommand; no auto-open browser.
- `trace install` = agents/MCP/hooks — not PATH for binary.
- Docs: `docs/gui-quickstart.md` leads with `./bin/trace serve`.
- P32-PORT: friendly bind conflict + explicit `--addr` — reuse for `gui`; **reject** UA-style auto-increment port.
- API: `/v1/graph` requires `center` + `max_nodes`; client `ops.ts` has `search` / `getGraph` / `listTasks` / `getProject` for composition research.
- Module: `github.com/mrchatam/Trace` (`go.mod`) for PATH/`go install` ranking.

## Planner-locked peer cite paths (for 01)

| Peer | Cite |
|------|------|
| Graphify open-file | `similar projects/graphify/worked/rsl-siege-manager/README.md`, `…/graph.html` |
| Graphify cluster energy | `similar projects/graphify/worked/mixed-corpus/raw/cluster.py` (optional); `…/rsl-siege-manager/review.md` |
| UA launch | `similar projects/Understand-Anything/understand-anything-plugin/packages/viewer/README.md` |
| UA open-browser | `…/packages/viewer/bin/viewer.mjs` (`--no-open`; open/start/xdg-open; auto-port → **reject**) |

## Research rejects (01 must document)

1. Graphify unbounded whole-graph as Trace default (Laws 6–7).
2. UA auto-increment listen port (vs P32-PORT).
3. PATH via `trace install` (agents/MCP only).
4. `/overview` as the Explore hook.
5. Three.js-first / always-on daemon / hosted SaaS.

## Out of this scope

- Writing DESIGN/UX-IA (S01), shipping `trace gui` (S02), Explore UI (S03), colorize (S04).
- Authoring `RESEARCH.md` in the planner row.
