# P29-S02-02 — UX IA review

## Metadata
- id: P29-S02-02
- todo_ids: [P29-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Independent review of `UX-IA.md` against S02-00 locks, ADR, OpenAPI, and Laws 6–7 / 19. Confirm a fresh agent could implement S04/S05 from the IA without inventing screens or second-SoT patterns.

**Fresh subagent.** Docs-only fixes for medium+ gaps preferred; spawn `02a`/`02b` only if structural rewrite needed.

## References

- [01-ux-ia.md](01-ux-ia.md) — implementer contract + locked defaults
- [00-PLANNER.md](00-PLANNER.md)
- Deliverable: [UX-IA.md](UX-IA.md) (must exist)
- [../scope-00-research/RESEARCH.md](../scope-00-research/RESEARCH.md)
- [docs/adr/ADR-HTTP-API-GUI.md](../../../../../adr/ADR-HTTP-API-GUI.md)
- [api/openapi.yaml](../../../../../api/openapi.yaml)
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)

## Session start

Follow agent-loop-protocol Session start. Do not share the implementer’s session.

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-29-http-api-browser-gui/scopes/scope-02-ux-ia/UX-IA.md
test -f docs/adr/ADR-HTTP-API-GUI.md
test -f api/openapi.yaml
# no product UI from this scope
test ! -d web || true   # web/ may appear only after S04; if present early, note blast radius — do not implement GUI here
```

## Checklist

### Completeness

- [ ] Required sections present (framing, nav IA, ×8 screens, empty/error/honesty, ship matrix, OpenAPI map, a11y, out-of-scope)
- [ ] All **eight** must-cover screens specified with purpose / data / actions / empty / error / `gui_ship`
- [ ] Recommended nav (or justified alternative) covers Overview, Tasks, Loop, Graph, Discoveries, Seed, Settings

### Contract vs GUI ship

- [ ] `api_wave` (OpenAPI `x-trace-wave`) distinguished from `gui_ship` (S04/S05/defer)
- [ ] Rich graph explorer = **S05**; S04 stub only — no full-dump viz
- [ ] Reviews **not** promoted into S04 MVP
- [ ] Loop: all five ops (`status` / `gate` / `next` / `apply` / `reset`) mapped; S04 read vs S05 interactive split honored
- [ ] Discoveries/decisions map to entities + search + links/transitions (capability = S05 enrichment)

### Laws / honesty

- [ ] Law 6–7: no full-graph dump as default GUI/API body; seed honesty warnings present
- [ ] Seed paths: project-root confinement called out (ADR)
- [ ] Law 19: no second-SoT UI patterns (no browser SQLite / parallel business rules)
- [ ] Operator chrome (not marketing landing)

### OpenAPI fidelity

- [ ] Operation map uses real `operationId`s / paths from `api/openapi.yaml`
- [ ] No invented endpoints; deferred waves not silently treated as S04
- [ ] Error states reference ADR envelope codes where useful

### Quality / spawn policy

- [ ] Findings by severity: blocker | high | medium | low | nit
- [ ] blocker/high: inline fix **or** spawn `P29-S02-02a` (implement) + `P29-S02-02b` (review) immediately below this row
- [ ] medium: prefer spawn unless trivial doc fix
- [ ] Confidence **medium** or **high** with evidence; residuals listed if medium

## Exit criteria

- [ ] Checklist complete; no open blocker/high without pending follow-up
- [ ] Confidence medium+; next **P29-S03-00**
- [ ] Board Notes cite evidence (paths, any inline fixes / spawns)

## Todo updates

Status + notes on **P29-S02-02** only (plus spawned rows if any). May thicken **upcoming** S03/S04 prompts if IA ship split needs reflection — do not rewrite S00/S01 `done` history.

## Next

**P29-S03-00** (after this row `done`)
