# P38-S05-02 — Saturation gate review

## Metadata
- id: P38-S05-02
- todo_ids: [P38-S05-02]
- role: reviewer
- skills: [doubt-driven-development, code-review-and-quality, analyst]
- mcps: [user-trace]
- verification: mixed

## Objective

**Gate owner.** Sign saturation or block S06. **Fresh subagent** — do not share S05-01 session.

**Blocks S06 until APPROVE (saturated).** This row is the only authorized path from investigation loops to REMEDIATION-PLAN.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — § Saturation exit criteria
- [GAP-REGISTRY.md](../scope-04-gap-registry/GAP-REGISTRY.md) — upstream SoT (row 661 APPROVE)
- [01-investigate.md](01-investigate.md) — investigation todos T0–T8 + locked defaults
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) — §5 spawn rules

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

---

## Checklist A — DESIGN-LOCKS saturation exit (mandatory)

Map each DESIGN-LOCKS criterion to SATURATION-NOTES §1 pass/fail + evidence:

| # | Criterion | SATURATION-NOTES §1 | Upstream verify |
|---|-----------|---------------------|-----------------|
| 1 | H1–H11 verified / rejected / deferred | | GAP-REGISTRY §2 G-001…G-011 + §5 |
| 2 | Live Trace command per major gap | | GAP-REGISTRY §7 → S01 evidence |
| 3 | Peer mechanism cites CG, UA, GF | | GAP-REGISTRY §2 dual-side rows |
| 4 | Moat row documented | | GAP-REGISTRY §3 M-001 |
| 5 | Spawn empty or deferred with trigger | | SATURATION-NOTES §3–§4 |
| 6 | High confidence — further rows duplicate | | SATURATION-NOTES §2 + §5 |

- [ ] All six criteria explicitly PASS or FAIL (no silent skip)
- [ ] Any FAIL → verdict **SPAWN** or **REQUEST CHANGES** — not APPROVE

---

## Checklist B — Planner must-answer items (P38-S05-00)

- [ ] **Q1** Checklist vs DESIGN-LOCKS — all boxes in §1
- [ ] **Q2** Confidence high/med/low + rationale in §2
- [ ] **Q3** Rejected spawn ideas in §5 (duplicate / out of scope / law conflict)
- [ ] **Q4** Explicit `ready_for_REMEDIATION_PLAN` boolean in §7

**Gate rule:** §7 `true` only valid when §6 = `PROCEED_TO_S06` and Checklist A all PASS.

---

## Checklist C — H7 compose-equivalence (planner lock)

GAP-REGISTRY §6 lists **H7 compose-equivalence** as open trigger. S05-00 locked **defer-with-trigger**, not spawn.

- [ ] SATURATION-NOTES §4 documents H7 defer with **owner** (S06) and **trigger** (remediation sketch / optional Phase 39 spike)
- [ ] `$EV/h7-compose-desk-check.md` exists with dimension table (T4)
- [ ] Desk-check conclusion consistent with G-007 gap verdict (expected: **not equivalent**)
- [ ] If desk-check claims equivalence → **SPAWN** S01-01a or **REQUEST CHANGES** (do not APPROVE saturated)
- [ ] H7 defer does **not** silently downgrade G-007 to non-gap without INVESTIGATION-INDEX “rejected if” evidence

---

## Checklist D — Spawn / defer hygiene

- [ ] §3 spawn list: if non-empty, each item has scope owner + row sketch (S01–S04)
- [ ] §4 deferred items: each has trigger + owner (not silent backlog)
- [ ] §5 reject list: ≥8 items covering duplicate, out-of-scope, law-conflict classes
- [ ] No REMEDIATION-PLAN / G1–Gn ranked themes in SATURATION-NOTES
- [ ] No product Go/TS/web diff in investigate row
- [ ] Evidence under `experiments/runs/…-p38-s05-663/evidence/` with date + row id headers

**SPAWN triggers (reviewer may insert board rows):**

| Condition | Action |
|-----------|--------|
| Checklist A any FAIL without documented defer | **SPAWN** back to owning S01–S04 scope |
| T4 reveals uncovered structural dimension | **SPAWN** S01-01a or S02-01a below row 664 |
| §6 = SPAWN but board rows missing | **REQUEST CHANGES** — implementer adds sketches; reviewer inserts rows |
| Premature S06 language in SATURATION-NOTES | **REQUEST CHANGES** — remove ranked remediation |

---

## Checklist E — Recommendation consistency + S06 gate

- [ ] §6 recommendation ∈ {`PROCEED_TO_S06`, `SPAWN`}
- [ ] `PROCEED_TO_S06` only when Checklist A all PASS
- [ ] `SPAWN` → §3 has actionable row sketches; **S06 blocked**
- [ ] GAP-REGISTRY still authoritative for gap facts — SATURATION-NOTES does not contradict G-001…G-011 without evidence
- [ ] Upstream APPROVED artifacts unchanged (forward-only — no rewrite of S01–S04 history)

**Spot-check (minimum 3 — pick checklist item 2, H7 T4, one reject item):**

- [ ] Open cited S01 evidence path from §1 row 2
- [ ] Open `h7-compose-desk-check.md` — dimension table complete
- [ ] One §5 reject item verified against GAP-REGISTRY §6 closed triggers

---

## Verdict

Document in board Notes:

```text
Verdict: APPROVE (saturated) | SPAWN | REQUEST CHANGES
Confidence: high | medium | low
Findings: blocker N · high N · medium N · low N · nit N
ready_for_REMEDIATION_PLAN: true | false
```

### APPROVE (saturated) when

- Checklists A–E PASS
- All six DESIGN-LOCKS criteria PASS in §1
- H7 defer-with-trigger documented (Checklist C)
- Confidence **high**, or **medium** with explicit residual risks in Notes (H7 live proof deferred is acceptable residual)
- §7 `ready_for_REMEDIATION_PLAN: true`

**Next row:** **P38-S06-00** only on APPROVE.

### SPAWN when

- Any DESIGN-LOCKS criterion FAIL without valid defer
- T4 desk-check reveals high-value investigation room
- §6 recommends SPAWN — insert rows **immediately below** row 664 per INVESTIGATION-INDEX §5
- Re-run S05 after spawn cycle `done`

**Do not start S06.**

### REQUEST CHANGES when

- Missing §1–§7 sections or §7 boolean inconsistent with §6
- H7 trigger missing from §4
- §5 reject list thin (<5 items) or missing law-conflict class
- Checklist walk lacks evidence pointers

Insert **P38-S05-01a** below row 664 if material rework needed; else fix in place if trivial (≤10 lines).

---

## Critical rule (re-state)

**Do not APPROVE saturation** unless investigator can confidently exit with **no high-value investigation room left** OR spawn list is empty/deferred with explicit triggers.

If not saturated → **SPAWN**, not S06.

---

## Next

`P38-S06-00` **only** on APPROVE (saturated).
