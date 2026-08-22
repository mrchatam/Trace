# P32-S05-00 — Scope planner (polish)

## Metadata
- id: P32-S05-00
- todo_ids: [P32-S05-00]
- role: planner
- skills: [planning-and-task-breakdown, documentation-and-adrs]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock polish: residual UX bugs from S03/S04, packaging nits, and docs — especially [`docs/gui-quickstart.md`](../../../../gui-quickstart.md) **multi-project / port** pattern from P32-PORT.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [OPEN-PORT-MULTI.md](../../OPEN-PORT-MULTI.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- Live: `docs/gui-quickstart.md`, `web/README.md`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Docs | Must document multi-project distinct `--addr` / free-port story as shipped in S02 |
| Scope | Polish + docs — not new explorer features |
| S04 craft (for docs/screenshots) | Cite live explorer craft: `.graph-shell` canvas-first + inspector `minmax(18–26rem)`; taller `--graph-canvas-height`; denser `PacketView` structured `dl` + raw in `<details>`; calm center/selected node chrome; motions `inspector-settle` / node transitions / sticky chrome (all `prefers-reduced-motion`) — not Phase 29 ops dual-card |

## Planner gate

- [x] `01-implement.md` lists doc + residual targets + optional screenshot surfaces from S04
- [x] `02-review.md` checks port docs

## Exit criteria

- [x] Next **P32-S05-01**

## Todo updates

Status + notes on **P32-S05-00** only.

## Next

`P32-S05-01`
