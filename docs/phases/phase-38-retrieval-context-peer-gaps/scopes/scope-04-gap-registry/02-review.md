# P38-S04-02 — Review gap registry

## Metadata
- id: P38-S04-02
- todo_ids: [P38-S04-02]
- role: reviewer
- skills: [code-review-and-quality, analyst, research]
- mcps: [user-trace]
- verification: mixed

## Objective

Independent review of **`GAP-REGISTRY.md`** vs INVESTIGATION-INDEX, planner locks (P38-S04-00), and approved S01–S03 artifacts. **APPROVE**, **REQUEST CHANGES**, or **SPAWN** (back to S01–S03 or S04-01a). **Fresh subagent** — do not share implementer session.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) — §2 verify/reject for H1–H11
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — saturation exit criteria (S05 forward-fit)
- [01-investigate.md](01-investigate.md) — investigation todos T0–T12 + locked defaults
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- **Upstream artifacts (baseline — do not require full re-audit):**
  - [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md)
  - [PEER-CG.md](../scope-02-codegraph-peer/PEER-CG.md)
  - [PEER-UA-GF.md](../scope-03-ua-graphify-peer/PEER-UA-GF.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

---

## Checklist A — Gap ID + hypothesis coverage

For **each** of G-001…G-011 (H1–H11):

- [ ] GAP-REGISTRY §2 row exists with **Gap ID** matching planner lock (G-001 ↔ H1 … G-011 ↔ H11)
- [ ] Verdict ∈ {gap, non-gap, defer} aligns with INVESTIGATION-INDEX **verified if / rejected if**
- [ ] Verdict **consistent** with upstream scope verdicts (TRACE-AUDIT, PEER-CG, PEER-UA-GF) — note justified upgrades (e.g. H8 inconclusive → supported)
- [ ] **Severity** column present — investigation confidence only (`high` \| `medium` \| `low`); **not** build priority / G1–Gn ranking
- [ ] **Law fit** column addresses G6 dump, G19 adapters, local-first, DR-NOSSEM where relevant

**Spot-check (minimum 4 — pick G-001, G-004, G-007, G-008):**

- [ ] Matrix cells populated for locked columns Trace | CG | UA | GF | MP
- [ ] Upstream artifact section cite resolves (§2/§3/§4 in S01–S03)
- [ ] H4 DR-NOSSEM split documented if G-004 is gap+defer

---

## Checklist B — Dual-side evidence gate (mandatory)

Every row with verdict **gap** must have **both**:

1. **Trace side:** `file:line` in Trace repo **or** TRACE-AUDIT evidence path **or** live command cite from S01
2. **Peer side:** ≥1 peer `file:line` **or** PEER-CG / PEER-UA-GF evidence path — **or** explicit `peer N/A` with reason (e.g. H9 UA/GF not applicable)

| Gap class | Minimum dual-side evidence |
|-----------|---------------------------|
| G-001 H1 orient packet | Trace: MCP schema / TRACE-AUDIT H1 + **CG or UA or MP** mechanism cite |
| G-002 H2 FTS | Trace: `compiler.go` L146 + UA SearchEngine contrast |
| G-003 H3 layers | Trace: live packet + `doc.go` L7 — peer optional (Aider doc defer OK) |
| G-004 H4 semantic | Trace: `doc.go` DR-NOSSEM + **GF or MP** mechanism |
| G-005 H5 index | Trace: TRACE-AUDIT H5 + CG watcher/lang cite |
| G-006 H6 MCP surface | Trace: TRACE-AUDIT H6 + **CG and/or MP** contrast cite |
| G-007 H7 explore | Trace: no `trace_explore` grep + CG explore mechanism |
| G-008 H8 onboarding | Trace: App.tsx / TRACE-AUDIT + **GF or UA or MP** hook cite |
| G-009 H9 intent | Trace: RETRIEVAL_AND_CONTEXT + zero grep + **MP** fact_checker contrast |
| G-010 H10 install moat | Trace: install detect + peer README orient-first cite |
| G-011 H11 stack docs | Trace: doc grep/read — CG complementary cite if claiming gap |

- [ ] Reviewer spot-checks ≥3 gap rows by opening cited evidence files or source lines
- [ ] No gap row is **docs-only synthesis** without pointer to S01 live evidence or S02/S03 peer mechanism
- [ ] §7 Evidence index links S01/S02/S03 `$EV/` folders — no orphan claims

**Dual-side failure:** If any gap row lacks Trace **or** peer side → **REQUEST CHANGES** or **SPAWN** S04-01a (not silent defer to S05).

---

## Checklist C — Matrix columns + moat row (planner lock)

- [ ] §2 main table columns include **Trace | CG | UA | GF | MP** (LOCK — not optional peers only)
- [ ] N/A cells justified (e.g. GF column for H5 index — “N/A — lang/index not GF focus”)
- [ ] **M-001** moat row in §3 — distinct from G-010 (under-promotion)
- [ ] Moat row cites merged strengths from TRACE-AUDIT §5 + PEER-CG §5 + PEER-UA-GF §5
- [ ] §4 Non-gaps lists peers-weaker items (CG no task loop, etc.)
- [ ] §5 Deferred includes DR-NOSSEM vector leg and rejected anti-patterns (PEER-CG §4)

---

## Checklist D — Planner must-answer items (P38-S04-00)

Confirm GAP-REGISTRY addresses all four handoff questions:

- [ ] **Q1** Unified G-001…G-011 linked to H* (§2 preamble table)
- [ ] **Q2** Severity = investigation confidence only — no ranked remediation (§2 column + explicit note)
- [ ] **Q3** Moat row M-001 — Trace strengths peers lack (§3)
- [ ] **Q4** Spawn triggers for S05 if incomplete (§6 — empty OK with rationale)

---

## Checklist E — Hygiene, phase law & saturation forward-fit

- [ ] No Go/TS/web product diff in investigate row
- [ ] No **REMEDIATION-PLAN** / G1–Gn ranked themes / “implement in P38” language
- [ ] No orphan claims from S01–S03 — every synthesized gap traceable to upstream §2/§3/§4
- [ ] H8 upgrade from S01 inconclusive documented with S03 evidence
- [ ] H7 compose-equivalence: either evidenced or in §6 spawn (not silent non-gap)
- [ ] G-004: vector leg defer vs label gap split — not collapsed into single vague row
- [ ] Evidence under `experiments/runs/…-p38-s04-660/evidence/` for H11 + synthesis notes
- [ ] Registry complete enough for S05 saturation review (DESIGN-LOCKS checklist forward-fit)

**S05 forward-fit minimum:**

- [ ] All H1–H11 → gap / non-gap / defer with evidence pointer
- [ ] Cross-matrix includes moat row
- [ ] Spawn list empty **or** explicit triggers — reviewer can sign saturation next scope

---

## Verdict

Document in board Notes:

```text
Verdict: APPROVE | REQUEST CHANGES | SPAWN
Confidence: high | medium | low
Findings: blocker N · high N · medium N · low N · nit N
```

### APPROVE when

- Checklists A–E PASS
- Dual-side evidence gate satisfied for all **gap** rows
- Confidence **high**, or **medium** with residual risks listed in Notes

### REQUEST CHANGES when

- Missing G-ID rows, matrix columns, or moat row
- Dual-side evidence missing on any gap row
- Severity confused with build priority
- Upstream verdict contradicted without justification

### SPAWN when

- Structural gap needs scoped re-investigate (insert **P38-S04-01a/01b** below row 661, or back to S01–S03 per INVESTIGATION-INDEX §5)
- H11 requires dedicated doc slice beyond S04-01
- H7 live compose-equivalence test needed before registry complete

**Do not proceed to S05** until registry APPROVE or spawn follow-ups are `done`.

---

## Next

`P38-S05-00` on APPROVE (or spawn first)
