# Phase 34 — GUI packaging, embed SPA, auto multi-project ports

**Phase planner.** Row `P34-00`.

## Metadata
- id: P34-00
- todo_ids: [P34-00]
- role: planner
- skills: [planning-and-task-breakdown, diagnosing-bugs]
- verification: automated

## Mission

Fix post–P33 dogfood: (1) **real SPA embedded** in Trace binary so consumer repos never need `web/`, (2) **auto free-port** for concurrent `trace gui`, (3) enforce **consumer = `.trace/` only** for Trace artifacts.

Read [`INTAKE.md`](INTAKE.md) + [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md).

## Gate

Phase 33 is **done** (DR-HANDOFF CLOSED). Proceed. Do not rewrite P33 history. L3 **supersedes** P32/P33 “no auto-port” for default-bind happy path only.

## Scope sequence

```
S00 research (embed vs install-sidecar; auto-port peers; static-dir defaults)
 → S01 plan (build/release embed pipeline + bind policy)
 → S02 implement embed + default static resolution (no consumer web/)
 → S03 implement auto free-port + gui/serve open URL
 → S04 docs/help/quickstart + tests
 → S05 VERIFY + DR-HANDOFF
```

| Scope | Board rows | Artifact |
|-------|------------|----------|
| S00 | P34-S00-00 → S00-01 → S00-02 | `RESEARCH.md` |
| S01 | P34-S01-00 → S01-01 | `PLAN.md` (planner + author; review optional via S01-00 gate) |
| S02 | P34-S02-00 → S02-01 → S02-02 | embed + StaticDir resolution |
| S03 | P34-S03-00 → S03-01 → S03-02 | auto free-port |
| S04 | P34-S04-00 → S04-01 → S04-02 | docs + residual tests |
| S05 | P34-S05-00 → S05-01 → S05-02 | VERIFY-NOTES + DR-HANDOFF CLOSED |

## Hard constraints

- Never require consumer `web/`
- Never put SPA under consumer project as the primary path
- Stub only as last-resort if embed empty (dev mistake) — fail loudly in release VERIFY if embed is stub
- Loopback defaults; `--addr` still honored when set (strict: fail if busy)
- Law 19 / no public bind defaults

## Planner gate (P34-00)

- [x] Phase folder README + light locks + live baseline
- [x] Each scope has runnable `00-PLANNER` / implement / review (or VERIFY) stubs
- [x] `SCOPE-TODOS.md` per scope
- [x] Board scope-sequence prose matches folders S00–S05
- [x] `DR-HANDOFF.md` OPEN; close owner `P34-S05-02`
- [x] No product / GUI / serve implementation in this row

## P34-00 outcome (2026-08-21)

Gate **PASS**. Thickened README (locks, live baseline, serial scopes, ownership, blast radius). DESIGN-LOCKS clarifications (L2 Trace-checkout vs consumer; L3 supersedes P32/P33 no-auto-port; reject SPA under `.trace/`). Protocol-thickened S00–S05 prompts + `SCOPE-TODOS.md`. Updated DR-HANDOFF checklist notes. No product code. Next: **P34-S00-00**.

## Next

`P34-S00-00`
