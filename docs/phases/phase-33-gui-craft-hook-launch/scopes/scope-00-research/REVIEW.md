# P33-S00-02 — RESEARCH review

**Verdict: PASS** · Confidence: **high** · 2026-08-21

Artifact: [`RESEARCH.md`](RESEARCH.md) vs DESIGN-LOCKS, INTAKE, live repo, peer cites.

## Checklist

### Locks & laws
- [x] Themes A–C honored (research informs; does not reopen)
- [x] Laws 6–7: no unbounded full-graph dump recommended as default
- [x] Law 19 / adapter-only noted where API discussed
- [x] PATH ≠ `trace install` agents/MCP conflation
- [x] CLI lean still compatible with Theme C (`trace gui` + open browser + PATH story)

### Artifact quality
- [x] All RESEARCH template headings present and non-empty
- [x] Live Explore described as center-first Graph home (not “missing Graph route”)
- [x] Explore ≠ Nav **Overview** (`/overview`) — Theme B is `/` Graph upgrade
- [x] Peer cites present (paths) for Graphify + UA — verified on disk
- [x] Peer matrix includes borrow/reject; port-on-conflict rejects UA auto-increment vs P32-PORT
- [x] Rejected alternatives cover: full dump; UA auto-port; PATH via `trace install`; `/overview` as Explore; Three.js / always-on daemon (+ hosted SaaS)
- [x] Overview options + **one** S01 recommendation (D+B+C) with named budgets
- [x] API lean `reuse` (with `reuse_then_gap_later` escape); rejects unbounded dump / seed-export-as-graph-body
- [x] CLI lean = **`trace gui`** primary
- [x] PATH options ranked + teaching note (user vs contributor)

## Spot-check evidence

| Claim | Live check |
|-------|------------|
| Index → Graph; `/graph` → `/` | `web/src/App.tsx` |
| Explore `/` ≠ Overview `/overview` | `web/src/layout/Nav.tsx` |
| `DEFAULT_MAX=50`, `UI_CAP=100`, Pick center | `web/src/screens/Graph.tsx` |
| No `gui`; serve + install | `cmd/trace/root.go` |
| Quickstart `./bin/trace serve` | `docs/gui-quickstart.md` |
| `/v1/graph` never full dump | `api/openapi.yaml` |
| Graphify `Start-Process` / `graph.html` | `similar projects/graphify/worked/rsl-siege-manager/README.md` |
| UA `--no-open` + auto-port | `…/packages/viewer/bin/viewer.mjs` |

## Findings

None at blocker/high. Residual (non-blocking): exact `--no-open` flag name and final budget numbers are for S01/S02 to lock — RESEARCH already leans them.

## Forward actions (review rights)

Thickened upcoming (not done) S01/S02 planner + implement prompts and SCOPE-TODOS with concrete RESEARCH leans. Did **not** rewrite `RESEARCH.md`.

## Next

**P33-S01-00**
