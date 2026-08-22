# P32-S03-02 — Depth review

## Metadata
- id: P32-S03-02
- todo_ids: [P32-S03-02]
- role: reviewer
- skills: [code-review-and-quality, frontend-design]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of S03 depth shell vs [`UX-IA.md`](../scope-01-ux-ia/UX-IA.md), [`DESIGN-LOCKS.md`](../../DESIGN-LOCKS.md), and [`01-implement.md`](01-implement.md) locks. **Do not fail for missing visual craft** (S04 owns typography/motion/density). **Block** if graph-home, inspector depth order, select≠expand, budgets/Laws 6–7/19, invented `/v1/path`, or missing `getImpact` wiring are wrong.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [`UX-IA.md`](../scope-01-ux-ia/UX-IA.md)
- [01-implement.md](01-implement.md) — Must-answer locks A/B/C
- S02: `getImpact` in `web/src/api/ops.ts`; `NO-GAPS.md`
- Live after S03-01: `web/src/App.tsx`, `Nav.tsx`, `Graph.tsx`, Inspector component, e2e

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Fresh subagent — do not share implementer session.

## Checklist (depth — not visual)

### Shell / hybrid C

- [ ] Index `/` is graph canvas home; Overview demoted (e.g. `/overview`)
- [ ] Nav reweighted: Explore/Graph first; not equal CRUD-first
- [ ] `/graph` still reaches graph (redirect or alias)
- [ ] Ops screens remain secondary routes (not a second SPA)

### Inspector depth

- [ ] Section order: summary → why → context → impact → reviews → links
- [ ] Optional loop strip on **task** selection only
- [ ] Task-scoped omit: Context / Impact / Reviews-filter / Loop hidden or collapsed for non-task
- [ ] Impact calls existing `getImpact(taskId)` — no duplicate wrapper; no invented path client
- [ ] Links from loaded neighborhood (+ deep-links); **no** `/v1/path`
- [ ] Structured presentation primary (not raw JSON-only)

### Selection / expand

- [ ] Select ≠ expand/re-center (single-click does not re-center)
- [ ] Explicit expand affordance (dblclick and/or “Use as center”)
- [ ] Progressive expand = new center + same budget ≤ UI_CAP

### Laws / budgets

- [ ] Laws 6–7: always `center` + `max_nodes`; DEFAULT_MAX=50 / UI_CAP=100; truncation banner; no full-dump CTA
- [ ] Law 19: `/v1` adapters only; no browser SoT / business-logic fork
- [ ] No Three.js / 3D default
- [ ] Kind/search filters retained on Graph

### Verify evidence

- [ ] Build PASS cited in Notes
- [ ] E2E/smoke (or documented env blocker + manual checklist) covers search/pick → center → select → inspector; select does not change center
- [ ] Visual polish gaps → note for S04, **not** blockers unless depth broken

## Findings

Severity: blocker | high | medium | low | nit.

- blocker/high: inline fix **or** spawn `P32-S03-02a` / `02b` immediately below this row
- medium: prefer spawn unless trivial
- Do **not** rewrite `done` S02/S01 history

## Exit criteria

- [ ] No open blocker/high without pending follow-up spawn
- [ ] Confidence medium or high with evidence
- [ ] If PASS: thicken upcoming S04 prompts lightly only if depth shell paths changed (optional)
- [ ] Next: **P32-S04-00** (do not start S04 implement)

## Todo updates

Status + notes on **P32-S03-02**.

## Next

`P32-S04-00`
