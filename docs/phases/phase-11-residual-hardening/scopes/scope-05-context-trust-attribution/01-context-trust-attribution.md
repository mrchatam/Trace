# P11-S05-01 — Context depth / trust / attribution (STUB)

## Metadata
- id: P11-S05-01
- todo_ids: [P11-S05-01]
- role: implementer
- verification: automated

## Objective
Implement DF-35 / DF-48 / DF-42 per **P11-S05-00** FINAL locks.

## Depends-on
- P11-S05-00 done; Phase 10 DF-19/27 regressions stay green

## Exit criteria (outline)
- [ ] Depth≥2 sibling leak fixed or fail-closed; named test
- [ ] Trust/binding wording consistent (Law 4/9); no `system` elevate
- [ ] DF-42 path per FINAL lock; DF-19 multi-goal still clean
- [ ] CGO0 retrieval/compiler + honesty; CGO1 p0x/x0 PASS
