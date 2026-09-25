# Phase 44 / S03 / 02b — Review INFERRED seed honesty fix

## Metadata
- id: P44-S03-02b
- todo_ids: [P44-S03-02b]
- role: reviewer
- skills: [code-review]
- verification: automated

## Objective

Independent review that **P44-S03-02a** closed H1 seed honesty. Re-check C1 (CLI-only) still holds — no creep while fixing export.

## Session start

Fresh subagent. Do **not** share implementer session. Follow agent-loop-protocol.

## References

- [02a-inferred-seed-honesty.md](02a-inferred-seed-honesty.md)
- [02-review.md](02-review.md) — parent H1/C1 rubric
- Board Notes on **P44-S03-02** / **P44-S03-02a**

## Rubric (must all PASS)

| Check | Pass |
|-------|------|
| Export includes `source_type` (+ confidence) on seed links | Yes |
| Import preserves `INFERRED` (no silent `USER_ASSERTED` upgrade) | Yes |
| Named round-trip test exists and PASS on re-run | Yes |
| Legacy seeds without fields still import (defaults OK) | Yes |
| Rel strings not mutated with `inferred_` prefix | Yes |
| C1 still PASS: no index-hook / daemon / MCP infer write | Yes |
| Inference rule library / dry-run / explicit-wins unchanged in intent | Yes |

## Evidence checklist

- [ ] Re-run `TestINFERREDScopeMemberSeedRoundTrip` (or Notes name) + `TestInferScopes_*` subset
- [ ] Spot-check `SeedLink` struct + `ImportSeedLink` meta path
- [ ] `rg` index/watch/MCP: still no `InferScopes` registration

## Verdict protocol

| Outcome | When |
|---------|------|
| **APPROVE** | All rubric PASS; confidence medium/high → next **P44-S04-00** |
| **Spawn** | Still dishonest or C1 broken → `P44-S03-02c` / `02d` |

## Exit criteria

- [ ] APPROVE or spawn
- [ ] SCOPE-TODOS S03-2 / S03-C updated
- [ ] Next **P44-S04-00** only after APPROVE

## Next

`P44-S04-00` (after APPROVE)
