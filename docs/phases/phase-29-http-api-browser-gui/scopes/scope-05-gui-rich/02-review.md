# P29-S05-02 — Feature-rich GUI review

## Metadata
- id: P29-S05-02
- todo_ids: [P29-S05-02]
- role: reviewer
- skills: [code-review-and-quality, webapp-testing]
- mcps: [cursor-ide-browser]
- verification: mixed
- hooks: []

## Objective

Independent review of GUI-P1 vs [`FEATURE-MATRIX.md`](FEATURE-MATRIX.md), UX-IA `gui_ship: S05`, Laws 6–7/19, and promote/export honesty. Fresh subagent — do not share S05-01 session.

## References

- [00-PLANNER.md](00-PLANNER.md) · [01-implement.md](01-implement.md) · [FEATURE-MATRIX.md](FEATURE-MATRIX.md)
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [UX-IA.md](../scope-02-ux-ia/UX-IA.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) §6–7, §19
- Live: `web/`, `internal/httpapi`, board Notes on P29-S05-01

## Session start

Follow agent-loop-protocol Session start. Fresh subagent.

## Preflight

1. Re-read FEATURE-MATRIX + P29-S05-01 Notes.
2. Confirm `web/` still builds: `cd web && npm run build`.
3. Spot-check `trace serve` serves SPA (not placeholder) when `web/dist` present.

## Checklist

### Matrix / ship

- [ ] Every `M01–M07` row `done` or `deferred` with explicit reason
- [ ] O01–O06 settled (optional done or deferred with reason)
- [ ] G-promote + G-export evidence present and re-runnable
- [ ] No silent promotion of deferred API (`x-trace-wave: defer`) into GUI

### Product surfaces

- [ ] Graph: `@xyflow/react` (or documented xyflow-class); budgeted `center`+`max_nodes`; truncated banner; **no** Three.js / unbounded dump
- [ ] Loop: next/apply/reset usable; apply/reset confirm dialogs; gate still visible when status errors (dogfood UUID / `INTERNAL_ERROR` residual OK at SPA)
- [ ] Discoveries: create + promote paths; S04 search not regressed; SearchResponse `items`
- [ ] Seed: export/import + path-confinement / 501 honesty (no silent ignore of `strict`/`task_id`)
- [ ] Tasks: fuller transitions; gate/warning display without SPA policy engine
- [ ] Reviews: nav + list/detail present

### Laws / architecture

- [ ] Law 19: browser calls `/v1` only; ops wrappers by operationId; no SQLite / invented SoT
- [ ] Law 6–7: budgeted graph; seed status/summary honesty (no full-graph body default)
- [ ] S04 MVP not rewritten (Shell/stack/Context patterns preserved)
- [ ] No agents/reviews leftovers incorrectly left as S04-only gaps

### Quality

- [ ] A11y floor: focusable nav; named icon controls; confirm dialogs keyboard-usable; gate text not icon-only
- [ ] Findings triaged: blocker/high → inline fix or spawn `02a`/`02b` immediately below this row

## Spawn policy

- **blocker/high:** small inline fix **or** insert `P29-S05-02a` (implement) + `P29-S05-02b` (review) under this row with full prompts
- **medium:** prefer spawn unless trivial
- Do **not** rewrite `done` history; thicken **upcoming** S06 prompts only if residuals belong there (e.g. `mapDomainErr` for loop UUID → already S06)

## Exit criteria

- [ ] No open blocker/high without pending follow-up
- [ ] Confidence **medium+** with evidence in Notes (build, matrix rollup, smoke/e2e)
- [ ] Next runnable **P29-S06-00** (unless spawns inserted)

## Todo updates

Status + notes on **P29-S05-02** only; spawn rights per protocol.

## Next

**P29-S06-00** (or spawned `P29-S05-02a` if inserted)
