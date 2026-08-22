# S02 — Packet honesty residuals — scope todos

**Depends-on:** P13-S01-02 APPROVE (DF-65 requires S01 resolve). Owns **DF-61, DF-62, DF-63, DF-65**.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL 2026-08-17 |
| 2 | 01-packet-honesty-residuals | implement | pending |
| 3 | 02-scope-review | review | pending |

## FINAL locks (summary)

| DF | Lock |
|----|------|
| DF-61 | `stale_total` + `stale_truncated`; keep sort-then-cap 8; MD total when truncated |
| DF-62 | Honesty over **pre-trim** file items (not kept-only); false-fresh on I/O miss |
| DF-63 | `items_total` = admit universe (layer0 + unique L1-admissible full candidate list); not post-cap ≤64 |
| DF-65 | `compileAtDepth` Expand **file** seeds (depth 1) after task Expand+FTS; reuse S01 resolve; preserve `edge_provenance` |

**Home:** `internal/compiler` (+ existing Retriever.Expand). **No** mig / analyzer rewrite / path-align. SchemaVersion stays **`0.2`**.

## Depends (from S01 FINAL)
- **DF-65:** Reuse Expand→`resolveImportedFile`; do **not** re-implement path join in compiler.

## Light note → S03
Additive packet fields under SchemaVersion `0.2` only — S03 owns DF-64/66/67 enum/CHECK; do not invent packet `0.3` here.

## Reminders
- Keep P12 SchemaVersion `0.2` + loud budget fields + named keepers
- Next after APPROVE: **P13-S03-00**
