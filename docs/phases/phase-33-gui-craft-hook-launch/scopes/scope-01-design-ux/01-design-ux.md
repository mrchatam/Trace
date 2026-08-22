# P33-S01-01 — Design + UX

## Metadata
- id: P33-S01-01
- todo_ids: [P33-S01-01]
- role: implementer
- skills: [impeccable, ui-ux-pro-max, frontend-design]
- mcps: []
- verification: automated
- hooks: []

## Objective

Author **markdown-only** design artifacts for Phase 33: a deliberate color/token brief + Explore-as-graph information architecture. Skills required. **No** product CSS/TSX into `web/` or `cmd/` (S03 implements IA; S04 lands tokens into the shell).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [00-PLANNER.md](00-PLANNER.md)
- [`../scope-00-research/RESEARCH.md`](../scope-00-research/RESEARCH.md) — binding **(D)+(B)+(C)** + budget leans
- Live baseline (read-only): `web/src/styles/tokens.css`, `web/src/styles/app.css`, `web/src/screens/Graph.tsx`, `web/src/layout/Nav.tsx`, `web/src/api/ops.ts`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify if blocked → Plan → execute).

**Skills gate (required before locking palette/IA):**

1. **impeccable** — Operate mode; treat this as **shape** + token/craft brief (not craft-floor UI edit). Visual authority = **refine** incumbent forest-moss / IBM Plex world, not greenfield cinema-dark or HUD. Load `shape` / colorize guidance as fits; do **not** ship CSS.
2. **ui-ux-pro-max** — Run at least one `--design-system` and one `--domain color` (and optionally `ux` / `style`) search; record query strings + what you **accepted vs rejected** in `DESIGN.md`. Prefer dense Operate / developer-tool signals; reject purple/cream/broadsheet clusters, HUD neon/glow, and wholesale slate+#22C55E “AI SaaS dark” replacement of Trace’s forest identity.
3. **frontend-design** — One opinionated POV for Trace (local-first knowledge graph for coding agents); one justified aesthetic risk that is **not** glow-slop; keep signature disciplined around the Explore canvas + kind chroma.

Evidence of skills use goes in board Notes (query strings / playbook names).

## Locked defaults (planner — do not re-debate)

| Item | Value |
|------|-------|
| Artifacts | **Two files** in this folder: [`DESIGN.md`](DESIGN.md) + [`UX-IA.md`](UX-IA.md) — do **not** merge into one file |
| Mode | **Operate** (impeccable) — scanability + graph hook; brand in precise details |
| Visual authority | **Refine** live `tokens.css` forest-moss / IBM Plex Sans+Mono — bolder kind/state chroma + clearer surfaces. Keep light+dark. Do **not** replace with purple/cream/broadsheet, glassmorphism/glow, or generic code-dark+#22C55E |
| Inputs | S00 RESEARCH **(D)+(B)+(C)** binding; Explore = `/` Graph — **≠** Nav `/overview` |
| Product CSS/TSX | **Out** — no edits under `web/` or `cmd/` |
| Color ownership | S01 = token **names + roles + hex candidates + craft floor**; **S04** applies to shell |
| IA ownership | S01 = UX-IA; **S03** implements composition in `Graph.tsx` / ops client |
| Clusters (A) | Visual grouping by kind/state only — **no** Leiden/community API in Phase 33 |
| API | Prefer **`reuse`** (`getProject` / `listTasks` / `search` / `getGraph`); `reuse_then_gap_later` only if seeds prove inadequate — never unbounded dump / seed-export-as-graph-body |
| Anti-slop | No glow-first nodes; no emoji-as-icons; color **not** sole encoder (kind label + optional glyph) |

### Budget locks (refine numbers only inside artifact rationale; model stays)

| Lean | Locked default | Hard bound |
|------|----------------|------------|
| Seed count | **Target 6**; acceptable **4–8** | ≤ **8** |
| Per-seed `max_nodes` | **40** | ≤ **50** (today’s `DEFAULT_MAX`) |
| Merged UI visible nodes | Honor live **`UI_CAP=100`** | Argue ≤ **120** only with written rationale in UX-IA; never >150 |
| Depth | **2** | OpenAPI default |
| Expand | User-driven re-center / expand only | No “load all” / expand-all default |

### Seed composition (D) — priority order for first paint

1. **`getProject`** — project/root identity if graph-addressable; else metadata only (do not invent fake nodes).
2. **`listTasks`** — prefer `IN_PROGRESS`, then non-terminal `PENDING` / other active; skip DONE/SKIPPED for seed slots unless needed to fill minimum.
3. **`search`** — fill remaining slots with high-signal entities (goals, capabilities, decisions) under search caps.
4. Dedupe by id; hard-stop at seed cap; then parallel budgeted **`getGraph(center, max_nodes, depth=2)`** per seed (**B**); merge/dedupe client-side under UI cap; **progressive expand** on user action (**C**).

### First paint / empty / error

| State | Required behavior (document in UX-IA) |
|-------|----------------------------------------|
| Happy | Explore open → interactive overview canvas with seed-composed nodes — **not** empty “Pick center” as the default hook |
| No seeds | Helpful empty: explain project has nothing to seed + path to Tasks / search; optional manual center remains **secondary**, not the hero |
| Partial API fail | Show loaded subgraph + non-blocking error (cause + retry); do not blank the canvas if any seeds succeeded |
| Hard fail | Error banner + retry; no fake full dump |

## Must answer (write into artifacts)

1. **Tokens** — CSS variable set (brand/surface/text/accent/focus + **`--kind-*`** + **`--state-*`**) with light+dark hex candidates; contrast floor **WCAG AA** (text ≥4.5:1; meaningful UI chrome/node stroke ≥3:1).
2. **Explore open** — D+B+C flow, seed sources, first-paint node set, expand, empty/error (table above).
3. **Kind/state color** — how nodes encode kind + task `work_state` without glow-slop and without color-only meaning (label/border/pattern).
4. **S03 vs S04** — explicit ownership split (below).

### S03 vs S04 ownership (binding)

| Owner | Owns | Does not own |
|-------|------|--------------|
| **S03** | Explore open composition (seeds → `getGraph` → merge → progressive expand); Graph UX (pan/zoom/click); keep inspector; budget chrome; empty/error states per UX-IA | Full shell recolor; redesigning Nav Overview (`/overview`); new SPA; Three.js |
| **S04** | Land `DESIGN.md` tokens into `tokens.css` / `app.css`; colorize chrome, nav, pills, Explore nodes (kind/state), empty/error surfaces; bolder/polish craft-floor | Re-opening IA model; unbounded graph API; PATH/`trace gui` |

## Artifact templates (author both)

### `DESIGN.md` (minimum sections)

```markdown
# Phase 33 S01 — DESIGN (tokens + craft)

## POV / mode
Operate; refine forest-moss Trace (cite live tokens.css).

## Skills evidence
- ui-ux-pro-max queries + accepted/rejected
- impeccable playbooks consulted
- frontend-design signature risk (one)

## Palette
Named roles + light/dark hex; map to CSS variables (extend existing names where possible).

## Token table
| Token | Role | Light | Dark | Contrast notes |

Must include: surfaces, text, accent/focus, danger/warn/ok, **kind-*** (goal/task/decision/discovery/claim/evidence/…), **state-*** (at least IN_PROGRESS / PENDING / DONE / FAIL-ish).

## Kind + state encoding
Rules for nodes/pills: fill/border/label; color-not-only.

## Craft floor
Anti-slop bans; motion tokens (keep --duration-*); reduced-motion; focus rings.

## Rejected directions
Purple/cream/broadsheet; HUD/neon glow; cinema glass; treating /overview as Explore.
```

### `UX-IA.md` (minimum sections)

```markdown
# Phase 33 S01 — UX-IA (Explore overview)

## Job
Open Explore (/) → interactive project overview graph hook (Graphify energy, Laws 6–7).

## Route clarity
Explore = `/` (Graph). Nav Overview = `/overview` (ops) — not this scope’s target.

## Open sequence (D+B+C)
Step-by-step seed → parallel getGraph → merge → first paint.

## Budgets
Locked numbers + any refined rationale.

## Progressive expand
User actions; what re-center means; no load-all.

## Inspector
Retained; selection behavior on overview vs expanded neighborhood.

## Empty / loading / error
Copy intent + recovery (not CSS).

## S03 handoff checklist
Bullet list implementer can execute without re-planning IA.
```

## Role work

1. Read RESEARCH + DESIGN-LOCKS + live tokens/Graph (baseline only).
2. Run skills gate; record evidence.
3. Author `DESIGN.md` and `UX-IA.md` per templates; answer Must-answer items.
4. Self-check exit criteria; board Notes only on **P33-S01-01**.

## Exit criteria

- [ ] `DESIGN.md` and `UX-IA.md` exist in this scope folder and cite S00 **(D)+(B)+(C)**
- [ ] Skills used with evidence in Notes (impeccable + ui-ux-pro-max + frontend-design)
- [ ] Token set + contrast floor + kind/state encoding specified
- [ ] Explore open IA specified (seeds, budgets, expand, empty/error); Explore ≠ `/overview`
- [ ] S03 vs S04 split explicit and consistent with table above
- [ ] No product code under `web/` or `cmd/`

## Minimal todos

- [ ] Read RESEARCH + locks + live baseline
- [ ] Skills pass (record queries / playbooks)
- [ ] Author `DESIGN.md`
- [ ] Author `UX-IA.md`
- [ ] Board Notes on **P33-S01-01**

## Todo updates

Status + notes on **P33-S01-01** only.

## Next

`P33-S01-02`
