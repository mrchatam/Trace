# P33-S01-02 — Design/UX review

## Metadata
- id: P33-S01-02
- todo_ids: [P33-S01-02]
- role: reviewer
- skills: [impeccable, ui-ux-pro-max, frontend-design, code-review-and-quality]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of S01 artifacts (`DESIGN.md`, `UX-IA.md`) for DESIGN-LOCKS compliance, craft floor, contrast/a11y, and Laws 6–7 Explore IA. Skills required (critique / audit mindset). Thicken upcoming **S03/S04** prompts if gaps found. **No product code** unless trivial markdown artifact fix; prefer spawn for material gaps.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [01-design-ux.md](01-design-ux.md)
- [`DESIGN.md`](DESIGN.md) + [`UX-IA.md`](UX-IA.md) (must exist)
- [`../scope-00-research/RESEARCH.md`](../scope-00-research/RESEARCH.md)
- Live baseline for drift checks only: `web/src/styles/tokens.css`, `web/src/screens/Graph.tsx`

## Session start

Fresh subagent. Follow agent-loop-protocol Session start. Load listed design skills before scoring craft/palette.

## Checklist

### Process / skills
- [ ] S01-01 Notes cite impeccable + ui-ux-pro-max + frontend-design with concrete evidence (queries / playbooks)
- [ ] Both artifacts present (not a single merged file unless planner later re-locked — default = two files)

### Locks / Laws
- [ ] Explore open = seed-composed project overview hook — **not** empty “Pick center” as the sole/default story
- [ ] Model is **(D)+(B)+(C)**; Explore = `/` Graph — **≠** `/overview`
- [ ] Budgets explicit within planner bounds (seeds ≤8; per-seed ≤50; merged ≤100–120; depth 2; user-driven expand)
- [ ] API lean = **reuse** (or justified `reuse_then_gap_later`); rejects unbounded dump / seed-export-as-graph-body
- [ ] No silent reopen of Three.js / second SPA / always-on daemon / hosted SaaS

### Craft / palette (impeccable + frontend-design)
- [ ] Visual authority = refine forest-moss Trace (or explicit justified supersession) — not AI purple/cream/broadsheet
- [ ] Rejected directions documented (glow/HUD/neon, glassmorphism-as-identity, generic slate+#22C55E wholesale replace)
- [ ] One clear POV / signature; Operate mode; no landing-page hero clutter for this surface
- [ ] Kind + state encoding without glow-slop; **color-not-only** (labels/borders/patterns)

### A11y / contrast (ui-ux-pro-max UX floor)
- [ ] Text contrast floor ≥4.5:1 called out for light **and** dark
- [ ] Meaningful strokes/icons ≥3:1 where color encodes kind/state
- [ ] Focus ring token/role present; keyboard path for graph selection mentioned or deferred explicitly to S03 with acceptance criteria
- [ ] `prefers-reduced-motion` respected in craft/motion notes
- [ ] Empty + error copy gives cause + recovery (not vibe-only)

### Ownership / handoff
- [ ] S03 vs S04 split matches 01 locked table (IA ship vs shell colorize)
- [ ] S03 handoff checklist in UX-IA is implementable without re-planning
- [ ] If gaps: thicken upcoming S03/S04 prompts or spawn `P33-S01-02a`/`02b` — do **not** rewrite S00 `done` history

## Findings format

Severity: `blocker` | `high` | `medium` | `low` | `nit`. Blocker/high → inline artifact fix **or** spawn implement+review pair. Exit only when no open blocker/high without a pending follow-up.

## Exit criteria

- [ ] Checklist complete with evidence citations (artifact section anchors)
- [ ] Confidence **medium** or **high**; residual risks listed if medium
- [ ] Board Notes on **P33-S01-02**; next runnable **P33-S02-00** (unless spawn inserted)

## Todo updates

Status + notes on **P33-S01-02** only (plus spawn rows if created).
