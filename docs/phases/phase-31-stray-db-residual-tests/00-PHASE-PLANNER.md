# Phase 31 — Stray `trace.db` residual testing

**Phase planner.** Row `P31-00`.

## Metadata
- id: P31-00
- todo_ids: [P31-00]
- role: planner
- skills: [test-driven-development, diagnosing-bugs, planning-and-task-breakdown]
- verification: automated

## Mission

Close the “Phase 30 needs more testing” gap without reopening store-path design. Expand automated coverage and a durable repro script around the shipped warn / path / gitignore behavior.

## Locked (human 2026-08-21)

- Canonical store remains `<root>/.trace/trace.db`
- No silent delete of root stubs
- No GUI work in this phase (Phase 32)
- Successor on green: **Phase 32** (graph-first GUI)

## Suggested coverage gaps (S00 confirms / trims)

1. Directory-named root `trace.db` → no warn (quiet)
2. MCP / `OpenExisting` path already covered — add CLI `trace tasks -C` stderr assertion in integration-style test or script
3. `trace serve` open path emits warn once at startup when stub present (if not covered)
4. Install / scaffold templates still carry `/trace.db` if more ignore scaffolds exist
5. Checked-in dogfood script under `experiments/` or `scripts/` reproducing: init → python stub → warn → `.trace/` still live
6. Document multi-open warn (once per `openStore`) as intentional

## Scope sequence

```
S00 inventory → S01 tests (+ script) + review → S02 VERIFY → Phase 32
```

| Scope | Board rows | Artifact |
|-------|------------|----------|
| S00 | P31-S00-00 → S00-01 | `GAPS.md` |
| S01 | P31-S01-00 → S01-01 → S01-02 | tests (+ optional script) |
| S02 | P31-S02-00 → S02-01 → S02-02 | VERIFY-NOTES + DR-HANDOFF CLOSED |

## Planner gate (P31-00)

- [x] Phase folder README + light locks
- [x] Each scope has runnable `00-PLANNER` / implement / review (or VERIFY) stubs
- [x] `SCOPE-TODOS.md` per scope
- [x] Board scope prose matches folders (no S03; S02 = VERIFY)
- [x] `DR-HANDOFF.md` OPEN; successor lean Phase 32
- [x] No product / test implementation in this row

## P31-00 outcome (2026-08-21)

Gate **PASS**. Thickened README (locks, in/out, live baseline, serial scopes). Clarified DR-HANDOFF (close owner `P31-S02-02`, successor Phase 32). Protocol-thickened S00–S02 prompts + `SCOPE-TODOS.md`. Fixed board scope-sequence prose. No product code; no store-path change. Next: **P31-S00-00**.

## Next

`P31-S00-00`
