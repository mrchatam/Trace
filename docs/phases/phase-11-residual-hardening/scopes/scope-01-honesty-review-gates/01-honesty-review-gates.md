# P11-S01-01 — Honesty / review-gate tightening (STUB)

## Metadata
- id: P11-S01-01
- todo_ids: [P11-S01-01]
- role: implementer
- skills: [tdd, incremental-implementation]
- verification: automated

## Objective
Implement DF-43 / DF-51 / DF-44 per **P11-S01-00** FINAL locks. Thin CLI/MCP only; G19.

## Depends-on
- P11-S01-00 done (FINAL locks)

## Exit criteria (outline — S01-00 thickens)
- [ ] DF-43: DONE rejected when a linked FAIL review exists alongside PASS (named regression test)
- [ ] DF-51: hatch WARNING mentions missing-caps override; behavior matches S01-00 lock
- [ ] DF-44: help/MCP schema clarifies Actor≠auth / flag≠identity
- [ ] Honesty A/B/C + Gate G + Phase 10 operator tests stay green
- [ ] `CGO_ENABLED=0` domain (+mcp as needed) PASS; `CGO_ENABLED=1` honesty/p0x/x0/`./...` product PASS

## Out of scope
- Real authentication; DF-45 review list CLI; S02–S05 work
