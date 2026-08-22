# P41-S02-00 — Scope planner (VERIFY)

## Metadata
- id: P41-S02-00
- todo_ids: [P41-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, qa-lead, shipping-and-launch]
- verification: automated

## Objective

Lock S02 VERIFY blocks for G8+G9 deliverables + Phase 42+ successor policy. Thicken `01-verify.md` + `02-dr-handoff.md`. **No product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [INTAKE.md](../../INTAKE.md)
- Pattern: [P40 S02-00 verify planner](../../../phase-40-read-surface-retrieval-depth/scopes/scope-02-verify/00-PLANNER.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)

## Session start

Follow agent-loop-protocol Session start. Unattended: DR-HANDOFF secondary queue is authority.

## Locked defaults (FINAL — P41-00)

| Item | Value |
|------|-------|
| Verify scope | G8 product (S00) + G9 product/doc (S01) |
| Precondition | P41-S00-02, P41-S01-02 both **APPROVE** |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p41-s02-01-verify/evidence/` |
| Notes artifact | `scopes/scope-02-verify/VERIFY-NOTES.md` (**required** at S02-01) |
| DR-HANDOFF | Stays **OPEN** until S02-02 |
| Successor | **Phase 42+** — G6/G7 secondary queue or `no successor` — never TBD at close |
| Close owner | P41-S02-02 |
| Product boundary | S02-00/01 verify only — no feature work |

## Verify blocks (for 01-verify.md)

| Block | Check |
|-------|-------|
| 0 | G8 G8-L1–L7 + S00-02 APPROVE |
| 1 | G9 G9-I1–I6 (or doc-revise supersede) + S01-02 APPROVE |
| 2 | M-001 moat preserved (task loop primary; layer/intent merge; no query-only; no dump) |
| 3 | Laws 6–7 caps honest; Law 19 library-first (CLI/MCP adapters thin) |
| 4 | Secondary queue G6/G7 documented in DR-HANDOFF — **no P41 implement rows** |
| 5 | Phase 42+ successor named — never TBD |
| 6 | `trace seed export` if entities changed during P41 |

## G8 accept map (Block 0)

| ID | Criterion |
|----|-----------|
| G8-L1 | Default context → `packet.layer` ≤ 1 |
| G8-L2 | `max_layer=2` → layer-2 items with reason_codes |
| G8-L3 | `max_layer=3` → L3 items when graph supports |
| G8-L4 | L2/L3 subject to budget caps |
| G8-L5 | Trim priority L0 > L1 > L2 > L3 |
| G8-L6 | `--depth` independent of `max_layer` |
| G8-L7 | No dump; MaxCandidateHits honored |

## G9 accept map (Block 1)

| ID | Criterion |
|----|-----------|
| G9-I1 | ExtractIntent from task fields |
| G9-I2 | Entity hint extraction |
| G9-I3 | Query token merge |
| G9-I4 | Search uses intent enrichment |
| G9-I5 | No semantic/vector channel |
| G9-I6 | Deterministic output |
| G9-DOC | §3 revised — intent shipped or aspirational supersede documented |

## Exit criteria

- [ ] `01-verify.md` + `02-dr-handoff.md` runnable with blocks 0–6
- [ ] Successor never TBD at close template
- [ ] Board row → `done` with Notes

## Next

`P41-S02-01`
