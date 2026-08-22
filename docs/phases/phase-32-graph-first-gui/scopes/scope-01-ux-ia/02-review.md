# P32-S01-02 — UX IA review

## Metadata
- id: P32-S01-02
- todo_ids: [P32-S01-02]
- role: reviewer
- skills: [code-review-and-quality, frontend-design]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independently review `UX-IA.md` against DESIGN-LOCKS, RESEARCH (PASS), and planner locks in `01-ux-ia.md`. Ensure graph-home hybrid C, inspector depth sufficient for S03, Laws 6–7 budgets, and honest S02 gap / P32-PORT notes. Thicken upcoming S02/S03 prompts only if IA proves gaps. **No product code** unless a tiny doc-adjacent fix; prefer spawn for substance gaps.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [00-PLANNER.md](00-PLANNER.md) / [01-ux-ia.md](01-ux-ia.md) locked defaults
- S00 [`RESEARCH.md`](../scope-00-research/RESEARCH.md)
- Artifact: [`UX-IA.md`](UX-IA.md)
- Live (spot-check): `web/src/App.tsx`, `web/src/screens/Graph.tsx`, `web/src/api/ops.ts`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: proceed without waiting for plan confirmation. Fresh context — do not share implementer session.

## Checklist

### Hybrid C / shell
- [ ] Graph is **home** (index / primary); Overview demoted — not nav-first CRUD
- [ ] Shell = evolve `web/` (no second SPA implied)
- [ ] Route / panel map covers Graph, Overview, Tasks(+detail), Loop, Reviews(+detail), Discoveries, Seed, Settings
- [ ] Nav reweight described (explore-first); secondary routes remain for deep links

### Inspector / depth (S03-ready)
- [ ] Sections include (order honored or justified): entity summary, why, context, impact, reviews, links/relations
- [ ] Each section maps to concrete `/v1` op(s)
- [ ] Structured presentation preferred over raw JSON-as-primary
- [ ] Select → inspector vs expand/re-center are **distinct** (today’s click-only-recenter called out as gap)
- [ ] Optional loop strip is panel-depth, not unbounded graph
- [ ] No requirement to clone UA tours / Graphify AST-only home

### Laws 6–7 / 19
- [ ] `center` + `max_nodes` required; DEFAULT_MAX/UI_CAP preserved or justified without unbounded dump
- [ ] Truncation honesty; **no** full-graph dump default / primary CTA
- [ ] Law 19: adapters only; no browser SoT / business fork

### API gaps → S02
- [ ] Gap list present (may be empty for library)
- [ ] `getImpact` client wrapper flagged if impact is in inspector map
- [ ] No invented `/v1/path` (or other non-library ops) without explicit reject/defer
- [ ] Notes: **P32-PORT still required** even if API is `NO-GAPS.md` / client-glue only
- [ ] Port design **not** owned by S01 (defer to S02)

### Process / out-of-scope
- [ ] Artifact is docs-only; no SPA rewrite in S01 Notes claim without evidence of code (reject if implementer shipped product code)
- [ ] S04 visual craft not smuggled in as S01/S03 must-have
- [ ] Confidence medium or high with evidence; blocker/high → fix or spawn

## Findings protocol

Per agent-loop-protocol reviewer loop: severity blocker | high | medium | low | nit. Blocker/high → small inline fix on `UX-IA.md` **or** spawn implement+review pair below this row. Medium: prefer spawn unless trivial. Re-verify until no open blocker/high without pending follow-up.

## Exit criteria

- [ ] Checklist complete with evidence in Notes
- [ ] No open blocker/high without follow-up row
- [ ] Confidence **medium** or **high** (residuals listed if medium)
- [ ] Upcoming S02/S03 thickened only if IA changed blast radius
- [ ] Next: **P32-S02-00**

## Todo updates

Status + notes on **P32-S01-02**; may thicken upcoming prompts only; may insert spawn rows immediately below if needed.

## Next

`P32-S02-00`
