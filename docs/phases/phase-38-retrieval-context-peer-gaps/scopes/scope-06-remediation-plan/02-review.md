# P38-S06-02 — Review remediation plan

## Metadata
- id: P38-S06-02
- todo_ids: [P38-S06-02]
- role: reviewer
- skills: [code-review-and-quality, analyst, doubt-driven-development]
- mcps: [user-trace]
- verification: manual

## Objective

Independent review of **`REMEDIATION-PLAN.md`**. Confirm plan-only artifact — no implement scope smuggled into P38. **Fresh subagent** — do not share S06-01 session.

**APPROVE** or **REQUEST CHANGES**.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — § REMEDIATION-PLAN shape
- [01-plan.md](01-plan.md) — locked defaults + theme registry G1–G9
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- **Upstream (read-only):**
  - [GAP-REGISTRY.md](../scope-04-gap-registry/GAP-REGISTRY.md)
  - [SATURATION-NOTES.md](../scope-05-saturation-gate/SATURATION-NOTES.md) — S05 APPROVE required
  - [h7-compose-desk-check.md](../../../../../experiments/runs/2026-08-22-p38-s05-663/evidence/h7-compose-desk-check.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

---

## Checklist A — Saturation gate + plan gate

- [ ] S05-02 **APPROVE (saturated)** referenced in REMEDIATION-PLAN §1 or board Notes
- [ ] `ready_for_REMEDIATION_PLAN: true` acknowledged (SATURATION-NOTES §7)
- [ ] Artifact path: `scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md`
- [ ] No product Go/TS/web diff in S06-01 row

---

## Checklist B — GAP coverage (mandatory)

Every registry row accounted for:

| GAP / item | Must appear in |
|------------|----------------|
| G-001 … G-011 | Exactly one G1–G9 theme **or** §4 reject with rationale |
| G-004a (vector) | §4 **defer/reject** — not a build theme |
| G-004b (label/graph) | G6 or explicit theme |
| M-001 moat | §1 executive — strengths preserved, not ranked as gap |
| Harness 9/16 fold | G3 or §4 defer |

- [ ] Coverage table in plan or reviewer-verified mapping
- [ ] All **high** severity gaps (10×) appear in ranked themes — not silent omit
- [ ] G-010 (medium) included (G3)

---

## Checklist C — Ranking rubric + theme quality

Locked rubric from P38-S06-00: **impact × law fit ÷ effort sketch**.

- [ ] §2 includes numeric or ordinal rubric columns per theme
- [ ] Final rank order documented with tie-breaker rationale (G1/G3 co-wave note acceptable)
- [ ] Each theme detail has: **problem, GAP ids, peer pattern, phase sketch, risks, not P38**
- [ ] Peer patterns cited as **reference** — not copy-paste mandates (esp. MP 44-tool, CG daemon)
- [ ] Phase sketches are **titles/bullets only** — no detailed implement board in P38

**Spot-check (minimum 3 themes):** G1 (G-001/002), G2 (G-007/H7), G4 (G-011/H11).

---

## Checklist D — H7 + H11 locked decisions (planner gate)

### H7 — compose-first vs unified `trace_explore`

SATURATION-NOTES §4 owner = S06. Desk-check: **not equivalent** (7/7).

- [ ] REMEDIATION-PLAN ranks **compose-first UX before unified `trace_explore`** (Phase 39 vs 40+)
- [ ] Does **not** claim multi-tool compose ≈ CG explore
- [ ] G-007 remains **gap** — not downgraded to non-gap
- [ ] Unified `trace_explore` sketch includes task-aware + law review + optional spike gate

### H11 — doc vs product dual-stack

Investigation conclusion (planner lock): **doc-only**.

- [ ] G4 / §1 states **doc-only** recommendation — not product dual-index integration
- [ ] Law 19 adapter boundaries mentioned for future doc
- [ ] Rejects product default dual-stack as remediation theme

---

## Checklist E — Reject / defer + law hygiene

- [ ] §4 reject/defer registry has **≥12** explicit items
- [ ] Includes PEER-CG §4 anti-patterns: daemon, MCP-only loop, dump defaults, graph-only product, query-only replaces task, CG benchmarks
- [ ] Includes SATURATION-class rejects: vector semantic (G-004a), implement in P38, MP 44-tool copy
- [ ] Laws 6/7/19 / DR-NOSSEM / local-first respected — no silent violation
- [ ] No ranked remediation language in SATURATION-NOTES or GAP-REGISTRY (upstream unchanged)
- [ ] §6 successor recommendation clear for S07 DR-HANDOFF (Phase 39 human-promoted entry)

**Reject spot-check (minimum 3):** daemon anti-pattern, G-004a vector defer, "implement in P38".

---

## Checklist F — Forward-only + spawn rules

- [ ] No rewrite of S01–S05 investigation artifacts
- [ ] No implement scope smuggled into P38
- [ ] If REQUEST CHANGES: spawn `P38-S06-01a` / `01b` below row 667 — do not edit `done` S06-01 prompt body
- [ ] Evidence under `experiments/runs/…-p38-s06-666/evidence/` if S06-01 created synthesis files

---

## Verdict

Document in board Notes:

```text
Verdict: APPROVE | REQUEST CHANGES
Confidence: high | medium | low
Findings: blocker N · high N · medium N · low N · nit N
```

### APPROVE when

- Checklists A–F pass
- No open blocker/high without pending follow-up row
- Confidence **medium** or **high** with explicit residual risks if medium

### REQUEST CHANGES when

- Missing GAP id coverage
- H7 or H11 decisions contradict planner lock / desk-check
- Implement language in P38
- Rubric missing or ranks vector semantic as theme
- Fewer than 12 rejects

---

## Exit criteria

- [ ] Verdict + confidence in board Notes (row 667)
- [ ] On APPROVE: next **P38-S07-00**
- [ ] On REQUEST CHANGES: spawn implement+review pair or notes for S06-01 fix

## Next

**P38-S07-00** on APPROVE.
