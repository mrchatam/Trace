# P20-S07-00 — Phase 20 verify planner (FINAL)

## Metadata
- id: P20-S07-00
- todo_ids: [P20-S07-00]
- role: planner
- verification: automated

## Objective
Lock VERIFY bar covering controller, artifacts, change/effects, gates, regression honesty, protocol, failure modes (doc §29O–Q, 31). Own DR-HANDOFF close policy. No product Go this row.

## Verify focus (must appear in S07-01)

- SelectNext table + hop budget tests (`internal/deliberation/...`)
- blocking uncertainty prevents EXECUTE (integration test via loop or domain)
- test-pass ≠ verified (S04 gate test cited)
- expected vs actual contradiction path (S03 + S05)
- correlated ≠ caused (S05)
- `trace loop next|apply|status` ordinary CLI — archived evidence under `experiments/runs/`
- P19 Loop keepers still green
- Compat checklist ceiling = **live embed max** (`CGO_ENABLED=1 go test ./evals/compat/... -run TestCompatibilitySecurityChecklist`) — S01 landed `015`; S02 **016**; S03 **017**; S04 **018**; S05-00 locked **019** (ceiling **19** after S05-01). Re-lock to embed max at VERIFY time. Seed export of P20 entities: **omit extension** — P17 keepers + documented residual.
- fixture-scale §31 mini-eval: investigate → decide → plan → change → verify/eval
- Failure modes §29O: hallucinated findings untrusted; incomplete evidence cannot promote

## DR-HANDOFF close policy (FINAL)

- S07-01 gathers evidence; **does not** close DR-HANDOFF
- **S07-02** must set successor decision explicitly: **`no successor`**
- Rationale: §16/§18 Future; hosted MCP off-board; human may promote forward phase later
- **Never TBD after S07-02 close**

## Locked defaults (FINAL — 2026-08-18)

| Item | Value |
|------|-------|
| Schema | 015–019 embedded; ceiling **19** |
| Loop | additive v1; P19 keepers + 14 S06 named tests |
| Seed export | P17 keepers PASS; P20 tables **omitted** (residual) |
| Successor | **`no successor`** (S07-02 closes) |
| Evidence | `experiments/runs/YYYY-MM-DD-p20-s07-01-verify/` |
| Next runnable | **P20-S07-01** |

## Planner work

1. [x] Aggregate verify command floor in `01-verify.md`
2. [x] Must coverage checklist from COVERAGE.md in `01-verify.md`
3. [x] §31 mini-eval spec in `01-verify.md`
4. [x] DR-HANDOFF close policy in `02-scope-review.md`
5. [x] Thicken S07-01 + S07-02 with live test names + seed export bar
6. [x] Update DR-HANDOFF successor intent (not closed until S07-02)
7. [x] No product Go

## Exit criteria

- [x] 01/02 thickened with evidence paths + test commands
- [x] Close policy documented in 02-scope-review
- [x] No product Go

## Next

Orchestrator: **P20-S07-01** after this row is `done`.
