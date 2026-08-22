# P33-S01-00 — Scope planner (design + UX)

## Metadata
- id: P33-S01-00
- todo_ids: [P33-S01-00]
- role: planner
- skills: [impeccable, ui-ux-pro-max, frontend-design, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Finalize S01 so implementer authors **color token brief** + **Explore-as-graph IA** (`DESIGN.md` and/or `UX-IA.md`) from S00 RESEARCH. Skills required. **No product code** (no shipping tokens into `web/` yet — that is S03/S04).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- `../scope-00-research/RESEARCH.md`
- Live: `web/src/styles/app.css`, `web/src/screens/Graph.tsx`, `web/src/layout/Nav.tsx`

## Session start

Follow agent-loop-protocol Session start. Run ui-ux-pro-max / impeccable guidance when thickening 01.

## Locked defaults

| Item | Value |
|------|-------|
| Artifacts | **Two files** (locked): `DESIGN.md` (tokens/palette/craft) + `UX-IA.md` (Explore IA) — do not merge |
| Skills | impeccable + ui-ux-pro-max + frontend-design **required** on 01 and 02 |
| Mode / authority | **Operate**; **refine** forest-moss / IBM Plex `tokens.css` — not purple/cream/broadsheet, HUD/glow, or wholesale slate+#22C55E |
| Explore job | Open → interactive **project overview** graph hook; inspector remains; Laws 6–7 per S00 recommendation |
| S00 IA lean (binding) | **(D)+(B)+(C)** — seed from `getProject` / `listTasks` / `search` → budgeted parallel `getGraph` → merge client-side → **progressive expand** on user action. **Not** empty “pick center” first paint. Explore = `/` Graph — **≠** Nav `/overview` |
| Seed priority | `getProject` → `listTasks` (IN_PROGRESS then active) → `search` fill → dedupe ≤8 |
| Budget locks | Seeds target **6** (4–8, ≤8); per-seed **40** (≤50); merged honor **`UI_CAP=100`** (argue ≤120 only in UX-IA); depth **2**; expand user-driven only |
| Clusters (A) | Visual grouping inspiration only — **no** Leiden/community API required in Phase 33 |
| API | Prefer **`reuse`** composition; `reuse_then_gap_later` only if seeds inadequate — never unbounded dump / seed-export-as-graph-body |
| Color ownership | S01 = **tokens + Explore overview IA only**; **S04** owns full shell colorize/craft |
| Avoid | AI purple/cream/broadsheet clusters; Three.js; second SPA; treating `/overview` as Explore |
| Hand-off | S03 implements IA; S04 applies full colorize — S01 does not ship CSS |

## Must answer (handoff to 01) — locked into `01-design-ux.md`

1. Token/CSS variable set (brand, surface, kind, state) + contrast floor → **yes**, via `DESIGN.md` template + WCAG AA floors.
2. Explore open under D+B+C → **yes**, seed priority + budgets + empty/error table in 01.
3. Kind/state without glow-slop → **yes**, color-not-only + `--kind-*` / `--state-*` required in DESIGN.
4. S03 vs S04 → **binding ownership table** in 01 (IA ship vs shell colorize).

## Planner gate

- [x] `01-design-ux.md` thick enough (skills, exit criteria, artifact template)
- [x] `02-review.md` includes craft + a11y/contrast checklist
- [x] `SCOPE-TODOS.md` accurate

## Exit criteria

- [x] Design implementer locked; next **P33-S01-01**

## Todo updates

Status + notes on **P33-S01-00** only.

## Next

`P33-S01-01`
