# P28-S06-00 — Scope planner (R6 FM residual wave)

## Metadata
- id: P28-S06-00
- todo_ids: [P28-S06-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: automated
- hooks: []

## Objective

Plan the **Phase 28 residual wave**: close actionable **R6 / FM-01/02/04/07/08/09/10** gaps promoted from the superseded forward queue. Lock per-FR implement/review prompts so the orchestrator can run **one fresh subagent per board row**.

S00–S05 history stays closed. This scope is **forward-only** after `P28-S05-02`.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — R6 partial; FM §3
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) §3
- [forward-p28-residuals.md](../../../../TODO/forward-p28-residuals.md) — **superseded index**; FR-P28-01…07 source text
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — S05 CLOSED + Residual wave OPEN
- [Phase 28 README](../../README.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked FR → board map (FINAL)

| FR ID | FM | Board implement | Board review | Theme |
|-------|-----|-----------------|--------------|-------|
| FR-P28-01 | FM-01 | P28-S06-01 | P28-S06-02 | Seed-import roster / promotion surfacing |
| FR-P28-02 | FM-02 | P28-S06-03 | P28-S06-04 | Write discoveries/decisions before export |
| FR-P28-03 | FM-04 | P28-S06-05 | P28-S06-06 | Parent vs worker Trace / task inheritance |
| FR-P28-04 | FM-07 | P28-S06-07 | P28-S06-08 | Git-sparsity warn-only **decision** (doc or gate) |
| FR-P28-05 | FM-08 | P28-S06-09 | P28-S06-10 | Prefer task/promotion path after `trace_add` |
| FR-P28-06 | FM-09 | P28-S06-11 | P28-S06-12 | Mode-collapse dual-lane proof beyond Session-B |
| FR-P28-07 | FM-10 | P28-S06-13 | P28-S06-14 | Live promotion API usage / measurement |

**Not board rows** (explicit defers — track only): FR-P28-D1, D2, D3, D4, X1 — see `SCOPE-TODOS.md`.

## Locked defaults

| Item | Value |
|------|-------|
| Product code | Allowed on implement rows **only** when FR acceptance requires it; prefer harness/docs/dogfood when sufficient |
| Auto-spawn discoveries→tasks without human gate | **Out** (FR-P28-D1) |
| Daemon / HTTP / hosted MCP | **Out** |
| Rewrite S00–S05 `done` prompts / Notes | **Out** — forward only |
| FM-07 | Default stay **warn-only** unless implementer records explicit human decision to ship plan-before-edit |
| S07 | Residual-wave VERIFY + DR-HANDOFF close after S06 reviews green |

## Scope boundary

- One FR per implement/review pair; no bundling FM closures across rows
- Reviewers are fresh sessions (not the implementer)
- No Phase 29 scaffold unless human promotes after S07

## Planner gate

- [x] All FR-P28-01…07 implement + review prompts exist and are runnable
- [x] `SCOPE-TODOS.md` lists locks + defers
- [x] S07 stub prompts exist for residual-wave VERIFY

## Exit criteria

- [x] Board rows P28-S06-01…14 + S07-* linked from `docs/TODO/phase-28.md`
- [x] Next runnable **P28-S06-01**
- [x] This row Notes cite FR map + supersession of forward queue

## Todo updates

Status + notes on **P28-S06-00** only.

## Next

`P28-S06-01`
