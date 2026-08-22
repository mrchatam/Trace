# Scope 07 — board map

**VERIFY + DR-HANDOFF** — closes Phase 38; scaffolds Phase 39.

| Order | Board ID | Prompt | Artifact / outcome |
|------:|----------|--------|-------------------|
| 668 | P38-S07-00 | [00-PLANNER.md](00-PLANNER.md) | Thickens 01/02 (FINAL locks) |
| 669 | P38-S07-01 | [01-verify.md](01-verify.md) | `VERIFY-NOTES.md` + evidence archive |
| 670 | P38-S07-02 | [02-dr-handoff.md](02-dr-handoff.md) | CLOSED `DR-HANDOFF.md` + Phase 39 scaffold |

---

## Locked verify blocks (0–6 — S07-00 FINAL)

| Block | Check | Primary evidence |
|-------|-------|------------------|
| **0** | No product Go/TS commits in P38 window | `git log --since=2026-08-22` on `internal/`, `cmd/`, `web/`, `trace/` |
| **1** | All H1–H11 in GAP-REGISTRY + 7 artifacts exist | GAP-REGISTRY §2; artifact manifest |
| **2** | S05 saturation APPROVE on record | SATURATION-NOTES §6 `ready_for_REMEDIATION_PLAN: true`; board 664 |
| **3** | REMEDIATION-PLAN ranked G1–G9 + ≥12 rejects | REMEDIATION-PLAN §2, §4 |
| **4** | Peer cites **CG, UA, GF, MP** | PEER-CG §3; PEER-UA-GF §2–§3; GAP-REGISTRY MP column |
| **5** | Moat row **M-001** | GAP-REGISTRY §3 |
| **6** | Successor Phase 39 in VERIFY-NOTES prep | REMEDIATION-PLAN §6 G1+G3+G4 |

---

## Required artifacts (S07-01 must confirm)

| Artifact | Scope | Path |
|----------|-------|------|
| INVESTIGATION-INDEX.md | S00 | `../scope-00-investigation-index/` |
| TRACE-AUDIT.md | S01 | `../scope-01-trace-audit/` |
| PEER-CG.md | S02 | `../scope-02-codegraph-peer/` |
| PEER-UA-GF.md | S03 | `../scope-03-ua-graphify-peer/` (+ **MP §3**) |
| GAP-REGISTRY.md | S04 | `../scope-04-gap-registry/` (+ **MP column**) |
| SATURATION-NOTES.md | S05 | `../scope-05-saturation-gate/` (APPROVE) |
| REMEDIATION-PLAN.md | S06 | `../scope-06-remediation-plan/` |

---

## Block 4 — Peer cite floor (MP extension)

| Peer | Doc | Required mechanism cites |
|------|-----|-------------------------|
| **CG** | PEER-CG.md §3 | `tools.ts` explore schema/handler; watcher debounce |
| **UA** | PEER-UA-GF.md §2 | `context-builder.ts` L25–79; `search.ts` L14–58 |
| **GF** | PEER-UA-GF.md §2 | EXTRACTED/INFERRED; `graph.html` orient |
| **MP** | PEER-UA-GF.md **§3** | `searcher.py` L276–329; `layers.py` L404–431; `service.py` L60–82; `fact_checker.py` L55–78 |

S03 evidence: `experiments/runs/2026-08-22-p38-s03-657/evidence/h*-mp-*.md`

---

## Phase 39 scaffold checklist (S07-02 — mandatory on APPROVE)

| Item | Path |
|------|------|
| Phase folder | `docs/phases/phase-39-context-orient-harness/` |
| README | goal G1+G3+G4; M-001 moat charter |
| Phase planner | `00-PHASE-PLANNER.md` — runnable **P39-00** |
| INTAKE stub | G1/G3/G4 in/out; links to P38 artifacts |
| DR-HANDOFF | OPEN |
| Scope 00 G1 | `scopes/scope-00-context-orient-merge/` stubs |
| Scope 01 G3 | `scopes/scope-01-harness-orient/` stubs |
| Scope 02 G4 | `scopes/scope-02-dual-stack-docs/` stubs |
| Scope 03 VERIFY | `scopes/scope-03-verify/` stubs |
| Board | `docs/TODO/phase-39.md` — P39-00 first pending |
| Index | `docs/TODO.md` Phase 39 row |

**Entry co-wave lock:** G1 + G3 + G4 (human promotes implement).

**Secondary queue (document only):** G5, G2 compose → Phase 39–40; G2 explore → Phase 40+.

---

## Successor decision (never TBD)

| Outcome | Successor | First runnable |
|---------|-----------|----------------|
| VERIFY green (default) | Phase 39 — Context orient & harness | `P39-00` after human promotion |
| Human idle | `no successor` | — |
| Repair needed | `pending repair spawn` | `P38-S07-02a` |

---

## Hypothesis → verify mapping

| H* | Block | Primary check |
|----|-------|---------------|
| H1–H11 | 1 | GAP-REGISTRY §2.1–§2.2 |
| H6, H8, H9 (MP) | 4 | PEER-UA-GF §3 + matrix MP column |
| H11 | 3, 6 | REMEDIATION-PLAN G4 doc-only; Phase 39 G4 scope |
| moat | 5 | M-001 §3 |

---

## Tool matrix (S07)

| Tool | S07-01 | S07-02 |
|------|--------|--------|
| Read S00–S06 artifacts | verify | review |
| git log boundary | Block 0 | spot-check |
| Archive + manifest.sha256 | Block 1 | confirm |
| VERIFY-NOTES.md | author | review |
| REVIEW-NOTES.md | — | author |
| Phase 39 scaffold | **no** | **yes** |
| Product code | **Forbidden** | **Forbidden** |

---

## Planner gate (S07-00)

- [x] `01-verify.md` thickened — blocks 0–6, MP in block 4, artifact manifest
- [x] `02-dr-handoff.md` thickened — Phase 39 scaffold spec, DR-HANDOFF template
- [x] SCOPE-TODOS 668–670 + verify block table
- [x] Successor lock: Phase 39 G1+G3+G4 entry co-wave

## Next

`P38-S07-01`
