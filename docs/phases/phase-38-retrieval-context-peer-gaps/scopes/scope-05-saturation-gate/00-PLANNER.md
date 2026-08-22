# P38-S05-00 — Scope planner (saturation gate)

## Metadata
- id: P38-S05-00
- todo_ids: [P38-S05-00]
- role: planner
- skills: [planning-and-task-breakdown, analyst]
- verification: automated

## Objective

Lock **saturation gate** — the only authorized exit from investigation loops. Output **`SATURATION-NOTES.md`** (S05-01). See DESIGN-LOCKS saturation exit criteria.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — § Saturation exit criteria
- [GAP-REGISTRY.md](../scope-04-gap-registry/GAP-REGISTRY.md) — §6 spawn triggers (S04 APPROVED)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) — §4.4 saturation forward-fit

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Inputs (required)

- [GAP-REGISTRY.md](../scope-04-gap-registry/GAP-REGISTRY.md) (APPROVED — row 661)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md)
- [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md) — summary §§2–5
- [PEER-CG.md](../scope-02-codegraph-peer/PEER-CG.md) — summary §§3–5
- [PEER-UA-GF.md](../scope-03-ua-graphify-peer/PEER-UA-GF.md) — summary §§3–5

## Critical rule

**Do not plan S05-01 to APPROVE saturation unless investigator can exit with no high-value investigation room left** — or spawn list is empty/deferred with explicit triggers.

If not saturated: plan for **SPAWN** (back to S01–S04), **not S06**.

## Locked defaults

| Item | Value |
|------|-------|
| Artifact | `scopes/scope-05-saturation-gate/SATURATION-NOTES.md` |
| Method | Checklist walk against DESIGN-LOCKS + GAP-REGISTRY — **not** re-audit S01–S03 |
| Evidence | `experiments/runs/YYYY-MM-DD-p38-s05-663/evidence/` (saturation notes + optional T4 desk-check) |
| Upstream SoT | GAP-REGISTRY §2–§7; link S01–S04 `$EV/` paths — no duplicate JSON |
| Product edits | **Forbidden** |
| H7 compose-equivalence | **Defer-with-trigger** (planner lock — **no spawn row**) — see below |
| Expected exit | **PROCEED_TO_S06** if checklist passes + confidence high |

## H7 compose-equivalence (planner decision — lock for 01/02)

| Option | Verdict | Rationale |
|--------|---------|-----------|
| Spawn S01-01a / S02-01a live compose test | **Reject** | Structural gap G-007 already verified (high); S01 live + S02 mechanism cites prove multi-tool split; live side-by-side unlikely to change verdict; duplicates saturation criterion |
| Defer-with-trigger in SATURATION-NOTES §4 | **Accept** | GAP-REGISTRY §6 documents open trigger; satisfies DESIGN-LOCKS “spawn list empty **or** deferred with trigger”; S05-01 T4 desk-check closes trigger with documented non-equivalence dimensions |
| Silent close without trigger | **Reject** | INVESTIGATION-INDEX H7 “rejected if” path requires evidence or explicit defer |

**T4 desk-check (S05-01):** Compare CG explore outputs (from `h7-explore-mechanism.md`) vs Trace compose surface (from `h7-explore-gap.md` + TRACE-AUDIT H1/H2) on **dimensions**: single-call cap, query param, source verbatim, call path, blast radius, task-packet merge. Record `$EV/h7-compose-desk-check.md`. Expected outcome: **not equivalent** — triggers S06 to rank unified tool vs compose-first UX.

**S06 trigger (defer owner):** REMEDIATION-PLAN must address G-007 with peer pattern + compose-vs-unified tradeoff; optional Phase 39 pre-implement live spike — not P38 investigation.

## Must answer for 01 handoff (planner pre-lock)

### 1. Checklist vs DESIGN-LOCKS saturation exit criteria

| # | Criterion | Pre-assessment | Evidence pointer |
|---|-----------|----------------|------------------|
| 1 | Every H1–H11 verified / rejected / deferred | **PASS** | GAP-REGISTRY G-001…G-011; G-004a defer; §5 rejects |
| 2 | Live Trace command per major gap claim | **PASS** | GAP-REGISTRY §7 → S01 `$EV/`; spot-check not required in S05 |
| 3 | Peer mechanism cites (CG, UA, GF) | **PASS** (+ MP) | GAP-REGISTRY §2 dual-side; S02/S03 mechanism files |
| 4 | Cross-matrix moat row | **PASS** | GAP-REGISTRY §3 M-001 |
| 5 | Spawn list empty or deferred with trigger | **PASS (1 defer)** | H7 compose-equivalence → §4 defer + S06 trigger |
| 6 | Confidence high — new rows duplicate | **PASS (expected)** | S05-01 documents reject list §5 |

### 2. Reviewer confidence (expected)

**High** — all H* routed through APPROVED S01–S04; dual-side evidence complete; one residual **medium** slice on H7 compose live proof (explicitly deferred, not blocking).

### 3. Rejected spawn ideas

| Idea | Reject reason |
|------|---------------|
| S01-01a / S02-01a H7 live compose-equivalence | Duplicate — structural non-equivalence already evidenced (`h7-explore-gap.md`, TRACE-AUDIT H1/H2) |
| Optional S01-01a symbol packet richness | Defer — GAP-REGISTRY §6; not saturation-blocking |
| H12+ Mempalace uncovered slice | Duplicate — mapped to H1/H4/H6/H8/H9 in S03 |
| Re-audit S01–S03 live CLI | Duplicate — GAP-REGISTRY synthesis APPROVED |
| Live `codegraph_explore` on Trace repo | Blocked — no `.codegraph/` index; mechanism cite sufficient (S02) |
| Cursor MCP 9/16 tool exposure | Fold — S06 harness/install hygiene (GAP-REGISTRY §6) |
| Semantic embedding spike | Law conflict — DR-NOSSEM (G-004a defer) |
| Implement `trace_explore` in P38 | Phase law — investigation only; S06 plan-only |
| Full peer re-scan (UA/GF/MP) | Duplicate — S03 APPROVED |
| Spawn back to S04 registry | Duplicate — row 661 APPROVE |

### 4. Ready for REMEDIATION-PLAN

**`true` (expected)** — if S05-01 checklist all PASS, T4 desk-check documents H7 defer, recommendation `PROCEED_TO_S06`. S05-02 must APPROVE before S06 starts.

## Planner gate

- [x] `01-investigate.md` thickened (T0–T8, locked H7 defer, deliverable §§1–6)
- [x] `02-review.md` requires Checklists A–E + saturation verdict table
- [x] SCOPE-TODOS IDs 662–664 + checklist registry
- [x] H7: defer-with-trigger, not spawn

## Exit criteria

- [x] S05-01/02 prompts runnable alone
- [x] Board `P38-S05-00` → `done`

## Next

`P38-S05-01`
