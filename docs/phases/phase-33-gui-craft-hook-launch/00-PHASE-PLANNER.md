# Phase 33 — GUI craft, Explore hook graph, `trace gui` launch

**Phase planner.** Row `P33-00`.

## Metadata
- id: P33-00
- todo_ids: [P33-00]
- role: planner
- skills: [impeccable, ui-ux-pro-max, frontend-design, planning-and-task-breakdown]
- verification: automated

## Mission

Close post–Phase 32 satisfaction gaps: **color/craft**, **Graphify-like Explore hook**, **practical `trace gui` launch** (PATH install).

Read [`INTAKE.md`](INTAKE.md) + [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md).

## Gate

Phase 32 is **done** (DR-HANDOFF CLOSED). Proceed. Do not rewrite P32 history.

## Scope sequence

```
S00 research (peers + Laws 6–7 overview graph + launch patterns)
 → S01 design tokens + UX (color system; Explore-as-graph IA)
 → S02 CLI `trace gui` + install/PATH (+ open browser; port reuse)
 → S03 Explore interactive project graph (hook)
 → S04 colorize / craft pass across shell (impeccable + ui-ux-pro-max)
 → S05 polish + docs (primary story = `trace gui`)
 → S06 VERIFY + DR-HANDOFF
```

| Scope | Board rows | Artifact |
|-------|------------|----------|
| S00 | P33-S00-00 → S00-01 → S00-02 | `RESEARCH.md` |
| S01 | P33-S01-00 → S01-01 → S01-02 | `DESIGN.md` / `UX-IA.md` |
| S02 | P33-S02-00 → S02-01 → S02-02 | `trace gui` + PATH install |
| S03 | P33-S03-00 → S03-01 → S03-02 | Explore project-graph hook UI |
| S04 | P33-S04-00 → S04-01 → S04-02 | Colorize / craft shell |
| S05 | P33-S05-00 → S05-01 → S05-02 | Docs + residual polish |
| S06 | P33-S06-00 → S06-01 → S06-02 | VERIFY-NOTES + DR-HANDOFF CLOSED |

## Hard constraints

- Law 19 / Laws 6–7 / loopback defaults
- Skills mandatory on S01/S03/S04 implement + review prompts
- Do not treat `./bin/trace serve` as the primary user-facing launch story after this phase
- `trace install` (agents/MCP) ≠ PATH install for the `trace` binary — keep distinct
- Prefer reuse of `/v1` graph + P32 Explore shell; no full-graph dump API

## Planner gate (P33-00)

- [x] Phase folder README + light locks + live baseline
- [x] Each scope has runnable `00-PLANNER` / implement / review (or VERIFY) stubs
- [x] `SCOPE-TODOS.md` per scope
- [x] Board scope-sequence prose matches folders S00–S06
- [x] `DR-HANDOFF.md` OPEN; close owner `P33-S06-02`
- [x] No product / GUI / serve implementation in this row

## P33-00 outcome (2026-08-21)

Gate **PASS**. Thickened README (locks, live baseline, serial scopes, ownership, blast radius). Light DESIGN-LOCKS clarifications (Explore center-first gap; `gui` vs `serve`; PATH vs `trace install`). Protocol-thickened S00–S06 prompts + `SCOPE-TODOS.md`. Updated DR-HANDOFF checklist notes. No product code. Next: **P33-S00-00**.

## Next

`P33-S00-00`
