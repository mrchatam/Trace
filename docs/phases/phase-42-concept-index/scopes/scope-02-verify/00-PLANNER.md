# P42-S02-00 — Scope planner (VERIFY)

## Metadata
- id: P42-S02-00
- todo_ids: [P42-S02-00]
- role: planner
- skills: [planning-and-task-breakdown, qa-lead, shipping-and-launch]
- verification: automated

## Objective

Lock S02 VERIFY blocks for G6+G7 deliverables + Phase 43+ successor policy. Thicken `01-verify.md` + `02-dr-handoff.md`. **No product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [INTAKE.md](../../INTAKE.md)
- Pattern: [P41 S02-00 verify planner](../../../phase-41-layers-intent/scopes/scope-02-verify/00-PLANNER.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)

## Session start

Follow agent-loop-protocol Session start. Unattended: P42-00 locks are authority.

## Locked defaults (FINAL — P42-00)

| Item | Value |
|------|-------|
| Verify scope | G6 product (S00) + G7 product (S01) |
| Precondition | P42-S00-02, P42-S01-02 both **APPROVE** |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p42-s02-01-verify/evidence/` |
| Notes artifact | `scopes/scope-02-verify/VERIFY-NOTES.md` (**required** at S02-01) |
| DR-HANDOFF | Stays **OPEN** until S02-02 |
| Successor default | **`no successor`** — G1–G9 remediation complete after P42 |
| Successor optional | Phase 43+ residuals wave (HTTP index, Tier-2 lang) — human promotion only |
| Close owner | P42-S02-02 |
| Product boundary | S02-00/01 verify only — no feature work |

## Verify blocks (for 01-verify.md)

| Block | Check |
|-------|-------|
| 0 | G6 G6-C1–C7 + S00-02 APPROVE + LAW-REVIEW-NOTES PASS |
| 1 | G7 G7-F1–F6 + S01-02 APPROVE |
| 2 | M-001 moat preserved (concept/index merge; no query-only; no dump) |
| 3 | Laws 6–7 caps honest; Law 19 library-first |
| 4 | G-004a vector absent; DR-NOSSEM honored |
| 5 | Phase 43+ successor named — **`no successor` default** — never TBD |
| 6 | `trace seed export` if entities changed during P42 |

## G6 accept map (Block 0)

| ID | Criterion |
|----|-----------|
| G6-C1 | Discovery/assumption body match → `graph_label_match` |
| G6-C2 | Concept entity filter excludes file/symbol |
| G6-C3 | Compile packet includes concept hits |
| G6-C4 | Limit cap honored |
| G6-C5 | No semantic/vector channel |
| G6-C6 | Deterministic output |
| G6-C7 | Fail-open on concept path error |

## G7 accept map (Block 1)

| ID | Criterion |
|----|-----------|
| G7-F1 | Index status lists supported_languages |
| G7-F2 | List matches adapter table |
| G7-F3 | Unsupported ext honest error |
| G7-F4 | Watch debounced index |
| G7-F5 | Watch foreground exit |
| G7-F6 | HTTP mirror (optional — residual ok) |

## Exit criteria

- [ ] `01-verify.md` + `02-dr-handoff.md` runnable with blocks 0–6
- [ ] Successor never TBD at close template
- [ ] Board row → `done` with Notes

## Next

`P42-S02-01`
