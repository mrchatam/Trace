# P33-S00-02 — Research review

## Metadata
- id: P33-S00-02
- todo_ids: [P33-S00-02]
- role: reviewer
- skills: [code-review-and-quality, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of `RESEARCH.md` against DESIGN-LOCKS, INTAKE, and live repo facts. Fix trivial RESEARCH gaps inline; thicken **upcoming** S01/S02 prompts if structural. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [01-research.md](01-research.md) — locked defaults, research rejects, peer cite paths
- Artifact: [RESEARCH.md](RESEARCH.md)
- Spot-check live: `web/src/App.tsx`, `web/src/layout/Nav.tsx`, `web/src/screens/Graph.tsx`, `cmd/trace/root.go`, `cmd/trace/install.go`, `docs/gui-quickstart.md`
- Spot-check peers (if citing): Graphify `worked/rsl-siege-manager/README.md` + `graph.html`; UA `packages/viewer/README.md` + `bin/viewer.mjs`

## Session start

Follow agent-loop-protocol Session start. Fresh context — do not share the S00-01 session.

## Checklist

### Locks & laws
- [ ] Themes A–C honored (research informs; does not reopen)
- [ ] Laws 6–7: no unbounded full-graph dump recommended as default
- [ ] Law 19 / adapter-only noted where API discussed
- [ ] PATH ≠ `trace install` agents/MCP conflation
- [ ] CLI lean still compatible with Theme C (`trace gui` + open browser + PATH story)

### Artifact quality
- [ ] All RESEARCH template headings present and non-empty
- [ ] Live Explore described as center-first Graph home (not “missing Graph route”)
- [ ] Explore ≠ Nav **Overview** (`/overview`) — Theme B is `/` Graph upgrade
- [ ] Peer cites present (paths or README) for Graphify + UA — not INTAKE-only handwave
- [ ] Peer matrix includes borrow/reject; port-on-conflict row rejects UA auto-increment vs P32-PORT
- [ ] Rejected alternatives cover: full dump default; UA auto-port; PATH via `trace install`; `/overview` as Explore; Three.js / always-on daemon
- [ ] Overview options + **one** S01 recommendation with named budgets
- [ ] API lean actionable (`reuse` | `thin_gap` | `reuse_then_gap_later`); no unbounded dump / seed-export-as-graph-body
- [ ] CLI lean = **`trace gui`** primary (or justified secondary with DESIGN-LOCKS still met)
- [ ] PATH options ranked + teaching note (user vs contributor)

### Forward board rights
- [ ] PASS → next **P33-S01-00**; optionally thicken S01 stubs if IA-shaped gaps
- [ ] blocker/high → inline fix or spawn `P33-S00-02a` / `02b`

## Findings format

`blocker` | `high` | `medium` | `low` | `nit` — every blocker/high needs fix or spawn before PASS.

## Exit criteria

- [ ] Confidence medium or high with evidence in Notes
- [ ] No open blocker/high without pending follow-up
- [ ] Checklist items above addressed or residual risk listed
- [ ] Next **P33-S01-00**

## Todo updates

Status + notes on **P33-S00-02** only.
