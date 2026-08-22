# P34-S00-00 — Scope planner (research)

## Metadata
- id: P34-S00-00
- todo_ids: [P34-S00-00]
- role: planner
- skills: [planning-and-task-breakdown, research]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock research scope so a fresh subagent can produce `RESEARCH.md`: embed vs install-sidecar, StaticDir consumer vs Trace-checkout policy, auto-port algorithm (L3 supersession), peer patterns. **No product code.** Do **not** author `RESEARCH.md` in this row.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [Phase 34 README](../../README.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — L1–L4
- [INTAKE.md](../../INTAKE.md)
- Live: `internal/httpapi/static.go`, `internal/httpapi/embed.go`, `internal/httpapi/embeddist/`, `internal/httpapi/addr_in_use.go`, `internal/httpapi/server.go`, `cmd/trace/local_http.go`, `cmd/trace/help.go`, `docs/gui-quickstart.md`
- Peer (auto-port borrow under L3): `similar projects/Understand-Anything/understand-anything-plugin/packages/viewer/bin/viewer.mjs` (if present)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Do not reopen L1–L4.

## Locked defaults

| Item | Value |
|------|-------|
| Output of S00-01 | `scopes/scope-00-research/RESEARCH.md` |
| Product / CLI / web edits | **No** on S00 |
| L2 lean | Prefer **`go:embed`** of real SPA; install-sidecar only if RESEARCH proves embed unblockable — still must not require consumer `web/` |
| L3 | Default bind busy → auto next free **loopback** port; `--addr` strict fail-if-busy |
| Supersession | Document that P33 rejected UA auto-port; **L3 overturns** for default bind only |
| Sequence | S00 → S01 serial; do not start S01 until S00-02 PASS |

## Must answer (handoff to 01)

1. Embed pipeline options (manual cp, make target, CI step, `go:embed` of built dist) + recommend one for S01/S02?
2. When (if ever) prefer disk StaticDir — Trace checkout only vs any `-C` with `web/dist`?
3. Auto-port algorithm: start port, scan range, `:0` then advertise, or peer-style increment — recommend one + loopback constraint?
4. `--addr` strict vs default auto — confirm L3 split?
5. Consumer layout audit: what paths today violate L1 in docs/help?

## Planner gate

- [x] `01-research.md` runnable (metadata, exit criteria, RESEARCH template)
- [x] `02-review.md` checklist vs DESIGN-LOCKS + INTAKE
- [x] `SCOPE-TODOS.md` accurate
- [x] Do **not** write `RESEARCH.md` in this planner row

## Exit criteria

- [x] Research implementer prompt locked for a fresh subagent
- [x] Board Notes cite what was thickened
- [x] Next remains **P34-S00-01**

## Todo updates

Status + notes on **P34-S00-00** only.

## Next

`P34-S00-01`
