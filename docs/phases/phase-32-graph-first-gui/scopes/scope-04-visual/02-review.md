# P32-S04-02 — Visual review

## Metadata
- id: P32-S04-02
- todo_ids: [P32-S04-02]
- role: reviewer
- skills: [code-review-and-quality, frontend-design, accessibility-auditing]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of S04 visual craft vs [`01-implement.md`](01-implement.md) craft locks A/B/C, [`DESIGN-LOCKS.md`](../../DESIGN-LOCKS.md), and S03 depth invariants. **Block** if depth stripped, IA/routes/API changed, Three.js introduced, or Laws 6–7/19 broken. **Block** if craft is absent (ops-shell feel unchanged) or a11y regresses. Fresh subagent — do not share implementer session.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [`UX-IA.md`](../scope-01-ux-ia/UX-IA.md)
- [01-implement.md](01-implement.md) — craft locks A/B/C
- S03 residuals: PacketView density / layout craft; e2e list-vs-canvas nit
- Live: `Inspector.tsx` (`PacketView`), `Graph.tsx`, `tokens.css`, `app.css`, e2e

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Checklist

### Depth / locks (must not regress)

- [ ] Inspector order summary→why→context→impact→reviews→links (+loop tasks); task-scoped omit
- [ ] Select ≠ expand; budgets 50 / UI_CAP 100 + truncation honesty
- [ ] Impact = `getImpact` only; no `/v1/path` in `web/src`
- [ ] Graph home `/` + Explore-first Nav; Law 19 adapters only
- [ ] No Three.js / 3D default

### Craft (S04 owned)

- [ ] **Layout:** `.graph-shell` canvas-first; inspector width/height improved; not ops dual-card dashboard
- [ ] **Typography:** clear inspector type roles; IBM Plex / forest-sage brand evolved, not replaced
- [ ] **PacketView density:** structured primary path tighter; raw JSON secondary (`details`) — S03 residual addressed
- [ ] **Canvas chrome:** center vs selected calm and distinct; empty inspector intentional
- [ ] **Motion:** 2–3 intentional; `prefers-reduced-motion` honored
- [ ] Explorer feel clearly better than Phase 29 ops shell (subjective but cite evidence)

### A11y

- [ ] `:focus-visible` / keyboard path usable on craft surfaces
- [ ] Named landmark / inspector region still discoverable
- [ ] Contrast acceptable on new/changed surfaces

### Verify evidence

- [ ] Re-run `cd web && npm run build` PASS cited
- [ ] Re-run `npm run test:e2e -- e2e/s03-depth.spec.ts e2e/s05-gates.spec.ts` PASS cited
- [ ] E2e canvas-vs-list select: if still list-only, **nit** only (same `onSelect`) — not a fail

## Findings

Severity: blocker | high | medium | low | nit.

- blocker/high: inline fix **or** spawn `P32-S04-02a` / `02b` immediately below this row
- medium: prefer spawn unless trivial
- Do **not** rewrite `done` S03/S02 history
- Visual polish gaps for S05 docs/screenshots → note forward, not S03 rewrites

## Exit criteria

- [ ] No open blocker/high without pending follow-up spawn
- [ ] Confidence medium or high with evidence
- [ ] If PASS: lightly thicken S05 prompts only if craft paths/docs need callouts (optional)
- [ ] Next: **P32-S05-00** (do not start S05 implement)

## Todo updates

Status + notes on **P32-S04-02**.

## Next

`P32-S05-00`
