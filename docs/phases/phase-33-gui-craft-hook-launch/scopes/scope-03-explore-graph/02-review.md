# P33-S03-02 — Explore graph review

## Metadata
- id: P33-S03-02
- todo_ids: [P33-S03-02]
- role: reviewer
- skills: [frontend-design, impeccable, ui-ux-pro-max, code-review-and-quality]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Independent review of Explore hook against Theme B, S01 UX-IA, Laws **6–7** / **19**, and S02 CLI land (`/`). Skills required. Spawn remediation (`P33-S03-02a`/`02b`) if open experience is still center-gate-only, dumps unbounded, or relocates Explore off `/`.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [`../scope-00-research/RESEARCH.md`](../scope-00-research/RESEARCH.md)
- [`../scope-01-design-ux/UX-IA.md`](../scope-01-design-ux/UX-IA.md) + [`DESIGN.md`](../scope-01-design-ux/DESIGN.md)
- [01-implement.md](01-implement.md) — Must-answer locks
- Live: `web/src/screens/Graph.tsx`, optional overview helper, `web/src/App.tsx`, `web/e2e/s03-depth.spec.ts`, `cmd/trace` land URL (read-only — do not reopen S02)

## Session start

**Fresh subagent** (not the implementer). Follow agent-loop-protocol Session start. Apply review skills: frontend-design, impeccable craft-floor, ui-ux-pro-max (Operate / a11y), code-review-and-quality. Note skills in board Notes / `REVIEW.md`.

## Checklist

### Theme B / open experience
- [ ] Open Explore (`/`) = interactive overview hook — **not** EmptyState “Pick center” / “No center selected” as default happy path
- [ ] Loading copy intent matches UX-IA (“Building project overview…” or equivalent)
- [ ] Manual center / search remains **secondary**, not the hero gate

### Laws 6–7 (budgets) + progressive expand
- [ ] Seed pipeline: `getProject` (no fake node) → `listTasks` priority → `search` fill (`q` non-empty) → dedupe; **target 6, ≤8**
- [ ] Parallel `getGraph(…, max_nodes=40, depth=2)` per seed
- [ ] Merge honors **`UI_CAP=100`** (trim keeps seeds first)
- [ ] Expand / re-center `max_nodes≤50`; user-driven only; **no** load-all / expand-all / unbounded dump
- [ ] No Leiden; no seed-export-as-graph-body; API **reuse** only (`ops.ts`)

### Interaction + inspector
- [ ] Pan/zoom/click; click → select + inspector usable
- [ ] Expand does not silently equal select (select ≠ re-center)
- [ ] Kind filter (if present) is client-side only

### Empty / error
- [ ] No seeds / partial / hard-fail: cause + recovery per UX-IA; partial keeps subgraph

### Keyboard / a11y
- [ ] Chrome + node list + inspector keyboard path; visible focus; no trap
- [ ] Canvas keyboard present **or** residual risk explicitly accepted in implementer Notes / REVIEW

### Routing / S02 contract
- [ ] Explore stays at **`/`**; `/overview` still ops Overview — not conflated
- [ ] **S02 CLI land:** SPA root remains `/` so `trace gui` → `http://{addr}/` still hits Explore (no silent relocate)

### S04 hooks + craft boundary
- [ ] `data-kind` / `data-state` (or equiv.) on nodes; kind text visible (color-not-only)
- [ ] No invented palette; shell-wide colorize **not** falsely claimed done (S04 owns)
- [ ] Skills evidence on implement Notes

### Law 19
- [ ] `web/` adapters only — no parallel SQLite / business-logic fork

### Evidence
- [ ] E2E and/or screenshot cited; `npm run build` (or equiv.) green
- [ ] Author [`REVIEW.md`](REVIEW.md) in this folder (PASS/FAIL + confidence + findings)

## Spawn rules

- **Blocker / high** (still center-gate default, unbounded fetch, Explore moved off `/`, Law 19 fork): insert `P33-S03-02a` implement + `P33-S03-02b` review immediately below this row; thicken prompts; do not rewrite `done` history.
- Low/nit: fix in-place only if tiny; else note for S04/S05.

## Exit criteria

- [ ] Confidence medium/high; checklist complete; [`REVIEW.md`](REVIEW.md) written
- [ ] Next **P33-S04-00** (unless spawn remediation)

## Todo updates

Status + notes on **P33-S03-02** only (plus spawn rows if created).

## Next

`P33-S04-00` (or spawned `P33-S03-02a` if FAIL)
