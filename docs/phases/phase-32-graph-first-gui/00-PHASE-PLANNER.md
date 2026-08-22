# Phase 32 — Graph-first GUI (explorer ambition)

**Phase planner.** Row `P32-00`. Phase 31 **CLOSED** (2026-08-21) — this row is the first runnable.

## Metadata
- id: P32-00
- todo_ids: [P32-00]
- role: planner
- skills: [brainstorming, frontend-design, planning-and-task-breakdown]
- verification: automated

## Mission

Raise the browser GUI from Phase 29’s ops shell to a **graph-first explorer** between Graphify and Understand-Anything: rich node detail first, then visual craft.

Read [`DESIGN-LOCKS.md`](DESIGN-LOCKS.md) before thickening scopes.

## Gate

Phase 31 is **done** (DR-HANDOFF CLOSED at `P31-S02-02`). Proceed.

## Scope sequence

```
S00 research (+ P32-PORT note) → S01 UX IA → S02 API gaps + P32-PORT ship → S03 depth → S04 visual → S05 polish (docs ports) → S06 VERIFY
```

| Scope | Board rows | Artifact |
|-------|------------|----------|
| S00 | P32-S00-00 → S00-01 → S00-02 | `RESEARCH.md` |
| S01 | P32-S01-00 → S01-01 → S01-02 | `UX-IA.md` |
| S02 | P32-S02-00 → S02-01 → S02-02 | API gaps / `NO-GAPS.md` **+ P32-PORT** |
| S03 | P32-S03-00 → S03-01 → S03-02 | Graph-home + inspector depth |
| S04 | P32-S04-00 → S04-01 → S04-02 | Visual craft on depth shell |
| S05 | P32-S05-00 → S05-01 → S05-02 | Docs + residual polish |
| S06 | P32-S06-00 → S06-01 → S06-02 | VERIFY-NOTES + DR-HANDOFF CLOSED |

## Hard constraints

- Law 19 / Laws 6–7
- Loopback `trace serve` defaults unchanged (port story may add free-port / clearer errors — see P32-PORT)
- No 3D default; no second SoT in the browser
- Depth scopes must land before visual-only polish is marked done
- **P32-PORT** must be investigated and addressed before VERIFY (see [`OPEN-PORT-MULTI.md`](OPEN-PORT-MULTI.md))
- **S02 always owns P32-PORT** even when API work is `NO-GAPS.md`

## Planner gate (P32-00)

- [x] Phase folder README + light locks + live baseline
- [x] Each scope has runnable `00-PLANNER` / implement / review (or VERIFY) stubs
- [x] `SCOPE-TODOS.md` per scope (S02 notes P32-PORT always)
- [x] Board scope-sequence prose matches folders S00–S06
- [x] `DR-HANDOFF.md` OPEN; close owner `P32-S06-02`
- [x] No product / GUI / serve implementation in this row

## P32-00 outcome (2026-08-21)

Gate **PASS**. Thickened README (locks, in/out, live baseline, serial scopes, P32-PORT ownership). Clarified DESIGN-LOCKS + OPEN-PORT-MULTI: S02 ships port story even with `NO-GAPS.md`. Protocol-thickened S00–S06 prompts + `SCOPE-TODOS.md`. Updated DR-HANDOFF checklist. No product code. Next: **P32-S00-00**.

## Next

`P32-S00-00`
