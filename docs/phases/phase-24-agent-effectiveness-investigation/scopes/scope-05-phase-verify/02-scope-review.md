# P24-S05-02 — Phase 24 verify review + DR-HANDOFF close

## Metadata
- id: P24-S05-02
- todo_ids: [P24-S05-02]
- role: reviewer
- skills: [code-review-and-quality, documentation-and-adrs, writing-for-agents, analyst]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer, tech-lead]
- verification: manual (checklist + spot-check)
- hooks: none

## Objective

Independent **fresh-session** review of S05-01 verify evidence and all Phase 24 scope artifacts. Re-run locked spot-checks (do **not** trust Notes alone). **Close DR-HANDOFF** with explicit Phase 25 successor decision (**never TBD**). Deliver **runnable Phase 25 scaffold** per agent-loop-protocol Phase handoff. Phase 24 complete when this row is `done`.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Reviewer loop + **Phase handoff (mandatory)**
- S05-00 locks: [00-PLANNER.md](00-PLANNER.md)
- S05-01 prompt: [01-verify.md](01-verify.md)
- SoT: [INVESTIGATION.md](../../INVESTIGATION.md), [README.md](../../README.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [FINDINGS.md](../../FINDINGS.md)
- S01–S04 deliverables (see Artifacts under review)
- Pattern: [P23 S06-02 review](../../../phase-23-enforcement-choke-points/scopes/scope-06-phase-verify/02-scope-review.md) (DR-HANDOFF close)
- Pattern: [P07 S03-02 review](../../../phase-07-performance-ladder/scopes/scope-03-phase-verify/02-scope-review.md) (Phase scaffold)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — must not be S05-01 implementer. Board edits: **status + notes only** on prior rows. Unattended: execute review loop until blocker/high clear or spawned forward.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Verify notes | `scopes/scope-05-phase-verify/VERIFY-NOTES.md` |
| Evidence archive | `experiments/runs/…-p24-s05-01-verify/evidence/` |
| Living findings | [FINDINGS.md](../../FINDINGS.md) |
| Post-mortem | [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md) |
| Codebase audit | [CODEBASE-AUDIT.md](../scope-02-codebase-loop-audit/CODEBASE-AUDIT.md) |
| External research | [EXTERNAL-RESEARCH.md](../scope-03-external-research/EXTERNAL-RESEARCH.md) |
| Intervention matrix | [INTERVENTION-MATRIX.md](../scope-04-intervention-matrix/INTERVENTION-MATRIX.md) |
| Phase handoff | [DR-HANDOFF.md](../../DR-HANDOFF.md) |

## Locked DR-HANDOFF close policy (FINAL — S05-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S05-01** — archive + VERIFY-NOTES; DR-HANDOFF stays **OPEN** |
| Who closes | **S05-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Default successor | **Phase 25 — P25-C orchestrator + default gap pass** (first human-promoted implementation theme) |
| Alternative successor | **`no successor`** only if human explicitly wants more dogfood before any Phase 25 build (document reason in board Notes) |
| Theme queue (not one mega-phase) | Human promotes **one theme per phase**: **P25-C** → **P25-A** → **P25-B** (order locked in DR-HANDOFF) |
| Deferred themes | **P25-D** (experiment protocol v2), **P25-E** (graph honesty) — remain deferred unless human overrides |
| Must not | Leave `Successor decision: TBD`, `later`, or empty; rewrite Phase 23 history; implement P25 product Go in this review row |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF **CLOSED** + all Phase 24 board rows `done` |

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S01: Dogfood post-mortem + failure taxonomy → `POSTMORTEM.md`, `FINDINGS.md` draft
- [ ] S02: Codebase loop/policy/task-creation audit → `CODEBASE-AUDIT.md`
- [ ] S03: External + similar-project research → `EXTERNAL-RESEARCH.md`
- [ ] S04: Intervention matrix → `INTERVENTION-MATRIX.md`, consolidated `FINDINGS.md`
- [ ] S05: VERIFY evidence + successor recommendation

### Successor decision table (locked default on APPROVE)

| Decision | When | Phase 25 folder focus |
|----------|------|------------------------|
| **Phase 25 — P25-C orchestrator + default gap pass** | Default — matrix top rank INT-03, INT-04, INT-11 | Harness: gap-pass install bundle, orchestrator Trace-first hook, hook drift checks |
| **`no successor`** | Human explicitly deferred implementation (Notes cite human decision) | No Phase 25 folder required; VERIFY Notes say `no successor` |
| Never | TBD / empty | — |

**P25-A** (discovery→task promotion) and **P25-B** (loop recalibration) are **documented in DR-HANDOFF** as **2nd/3rd** human promotions — not scaffolded as active board phases until P25-C completes or human skips C.

### S04→S05 residuals (list on close — non-blocking)

| Residual | Disposition for Phase 25 planning |
|----------|-------------------------------------|
| Auto-spawn human gate | P25-A planner owns confirm-before-spawn vs guided INT-01 |
| P19 threshold dogfood validation | P25-B owns live recalibration dogfood; not blocking P25-C close |
| Hook API drift (INT-11) | Include in P25-C scope as maintenance spike row |
| Live gate `reason_code` env-dependency | P25-B INT-09; document in DR-HANDOFF residuals |

### DR-HANDOFF.md update template (on APPROVE — default P25-C)

```markdown
**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Closed | YYYY-MM-DD |
| Successor decision | **Phase 25 — P25-C orchestrator + default gap pass** |
| Phase 24 outcome | Two-mode investigation complete; 11 ranked interventions; FM-01..FM-10 mapped; external comparables; FINDINGS consolidated |
| Recommended promotion order | P25-C → P25-A → P25-B (one theme per phase) |
| Top interventions | INT-03 gap pass, INT-04 orchestrator Trace-first, INT-11 hook drift |
| Residuals (non-blocking) | Auto-spawn human gate (P25-A); P19 dogfood validation (P25-B); sticky STOP UX (INT-09); live gate env-dependency |
| Forward | Human promotes Phase 25 board when ready — first runnable **P25-00** |
```

If human chose **`no successor`**: replace successor line and omit Phase 25 scaffold requirement (VERIFY must document human evidence).

## Phase handoff scaffold (mandatory on APPROVE with P25-C successor)

Per agent-loop-protocol, before marking this row `done` with P25-C successor:

### Required Phase 25 artifacts

| Artifact | Path (locked default) |
|----------|----------------------|
| Phase README | `docs/phases/phase-25-orchestrator-gap-pass/README.md` |
| Phase planner | `docs/phases/phase-25-orchestrator-gap-pass/00-PHASE-PLANNER.md` |
| Design SoT stub | `docs/phases/phase-25-orchestrator-gap-pass/GAP-PASS.md` (goal, INT-03/04/11 scope, in/out) |
| DR-HANDOFF open | `docs/phases/phase-25-orchestrator-gap-pass/DR-HANDOFF.md` (**OPEN**) |
| Scope stubs (≥1) | e.g. `scopes/scope-01-gap-pass-install/` with `00-PLANNER.md`, `01-*.md`, `02-scope-review.md`, `SCOPE-TODOS.md` (minimal OK) |
| Board file | `docs/TODO/phase-25.md` — rows **P25-00** (phase planner) first **pending** after Phase 24 last `done` |
| Index link | `docs/TODO.md` — Phase 25 row in phase boards table |

### Phase 25 README minimum content

- Goal: collapse E01 Mode A→B via **default gap pass** + **orchestrator Trace-first** (INT-03, INT-04, INT-11)
- Evidence links: Phase 24 INTERVENTION-MATRIX §1 top-3, POSTMORTEM §2 two-mode table
- In scope: `internal/install/`, cursor hook, orchestrator prompt bundle, hook drift spike
- Out of scope: daemon; hosted MCP; full loop recalibration (P25-B); task promotion product (P25-A); rewriting P24 history
- Human promotes **one theme** — this phase is **P25-C only**

### 00-PHASE-PLANNER minimum content

- Metadata + objective (lock scopes against live repo + Phase 24 artifacts)
- References to Phase 24 DR-HANDOFF + INTERVENTION-MATRIX
- Session start + exit criteria + **Next: P25-S01-00** (or first scope planner)
- **No product Go** on P25-00 row

### Board registration template (`docs/TODO/phase-25.md`)

```markdown
Index: [`docs/TODO.md`](../TODO.md)

## Phase 25 — Orchestrator + default gap pass (P25-C)

Human-promoted after Phase 24 close. Design SoT: [`phases/phase-25-orchestrator-gap-pass/GAP-PASS.md`](../phases/phase-25-orchestrator-gap-pass/GAP-PASS.md).

| Order | ID | Status | Prompt | Notes |
|------:|----|--------|--------|-------|
| 429 | P25-00 | pending | [phases/phase-25-orchestrator-gap-pass/00-PHASE-PLANNER.md](…) | Scaffold from P24-S05-02 close. |
```

(Adjust order IDs to next free integers in master sequence.)

### AGENTS.md update (reviewer may note for orchestrator)

After Phase 24 close, orchestrator paste should point to Phase 25 / P25-00 — **do not** edit AGENTS.md unless reviewer rights explicitly include it; at minimum cite in board Notes.

## Review focus

Confirm independently:

1. **S05-01 evidence** — evidence dir exists; 6 archived files; `manifest.sha256`; VERIFY-NOTES maps Bars 1–6
2. **Completion bar** — all INVESTIGATION/README criteria PASS (spot-check, not trust alone)
3. **Two-mode integrity** — Session A/B not conflated across FINDINGS + POSTMORTEM + matrix
4. **Cross-artifact consistency** — matrix top-3 ↔ DR-HANDOFF ↔ FINDINGS executive summary
5. **Trace law** — top-3 interventions do not require daemon / hosted MCP / full-rebuild indexer
6. **S04→S05 residuals** — documented in VERIFY-NOTES with forward owner (P25-A/B/C)
7. **DR-HANDOFF close** — CLOSED with explicit successor (default P25-C)
8. **Phase 25 scaffold** — runnable if P25-C successor (not README-only)

## Evidence to re-verify (reviewer spot-checks)

| Check | Source | Pass criterion |
|-------|--------|----------------|
| Archive complete | S05-01 evidence dir | 6 files + manifest + metadata |
| FM count | POSTMORTEM §3 | 10 FM rows |
| INT count | INTERVENTION-MATRIX §2 | ≥8 rows (expect 11) |
| URL count | EXTERNAL-RESEARCH | ≥3 http(s) URLs |
| FINDINGS status | FINDINGS status table | no `pending` |
| Recommended themes | DR-HANDOFF | 1–3 with INT links |
| Session B task count | POSTMORTEM | 5 seed / 7 discoveries / 0 new tasks |
| Top-3 ranks | INTERVENTION-MATRIX §1 | INT-03, INT-04, INT-01 |

### Locked re-verify commands (minimum)

```bash
EVID=$(ls -d experiments/runs/*-p24-s05-01-verify/evidence 2>/dev/null | tail -1)
test -d "$EVID" && ls -la "$EVID"
test -f "$EVID/manifest.sha256" && wc -l "$EVID/manifest.sha256"

grep -c '| FM-' docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-01-dogfood-postmortem/POSTMORTEM.md
grep -c '| INT-' docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md
grep -cE 'https?://' docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-03-external-research/EXTERNAL-RESEARCH.md
grep -i pending docs/phases/phase-24-agent-effectiveness-investigation/FINDINGS.md || echo "no pending"
grep 'Recommended' docs/phases/phase-24-agent-effectiveness-investigation/DR-HANDOFF.md
```

Broader walk of VERIFY-NOTES completion bar table encouraged; **minimum** commands above mandatory for APPROVE.

## Review checklist — VERIFY-NOTES + evidence

### Blockers

- [ ] **Blocker:** Missing evidence dir or fewer than 6 archived artifacts
- [ ] **Blocker:** Missing `manifest.sha256` or VERIFY-NOTES
- [ ] **Blocker:** Any completion bar row FAIL without pending spawn
- [ ] **Blocker:** FINDINGS status table has `pending` sections
- [ ] **Blocker:** INTERVENTION-MATRIX <8 rows or missing §1/§3/§5
- [ ] **Blocker:** DR-HANDOFF missing Recommended themes or successor still TBD at close
- [ ] **Blocker:** Session A/B conflated in FINDINGS executive summary

### High

- [ ] **High:** VERIFY-NOTES verdict PASS but spot-check counts below thresholds
- [ ] **High:** Top-3 matrix ranks disagree with DR-HANDOFF without note
- [ ] **High:** S04→S05 residuals missing from VERIFY-NOTES
- [ ] **High:** Phase 25 scaffold incomplete when successor is P25-C
- [ ] **High:** Top-3 recommends daemon/hosted MCP on P0-X core path

### Medium / low

- [ ] **Medium:** Evidence copies stale vs live source (hash mismatch) — re-archive or note
- [ ] **Medium:** EXTERNAL-RESEARCH URL spot-check failed (≥1 broken link)
- [ ] **Low:** VERIFY-NOTES missing git SHA in metadata
- [ ] **Nit:** Typo-only issues in VERIFY-NOTES

## Review checklist — scope deliverables (S01–S04 regression)

- [ ] POSTMORTEM §1–§4 intact; E01 A+B evidence paths cited
- [ ] CODEBASE-AUDIT §2 FM rows have file:line cites
- [ ] EXTERNAL-RESEARCH §2 ≥3 comparables with actionable deltas
- [ ] INTERVENTION-MATRIX §4 human-gate (auto-spawn) documented
- [ ] FINDINGS links INTERVENTION-MATRIX; does not duplicate full §2 table

## Review checklist — DR-HANDOFF + Phase handoff

- [ ] **Blocker:** DR-HANDOFF marked CLOSED before independent re-verify
- [ ] **Blocker:** Successor left TBD when row marked `done`
- [ ] **Blocker:** P25-C successor chosen but Phase 25 folder/board missing
- [ ] **High:** Phase 25 README missing in/out scope for INT-03/04/11
- [ ] **High:** `docs/TODO/phase-25.md` missing or P25-00 not first pending row
- [ ] **High:** `docs/TODO.md` index missing Phase 25 link
- [ ] **Medium:** P25-A/B documented as future promotions (not merged into P25-C mega-scope)

## Spawn policy

| Severity | Action |
|----------|--------|
| blocker / high | Small inline doc fix **or** spawn `P24-S05-02a` (implement) + `02b` (re-review) immediately below |
| medium | Prefer spawn unless typo-only (≤5 lines) |
| low / nit | List in REVIEW-NOTES; do not block close |

Do not rewrite S05-00 / S05-01 `done` prompt bodies.

## Evidence artifacts (reviewer output)

- Read S05-01 `experiments/runs/…/evidence/*` + `VERIFY-NOTES.md`
- Write **`REVIEW-NOTES.md`** in this scope folder (recommended)
- Update [DR-HANDOFF.md](../../DR-HANDOFF.md) on APPROVE
- Create Phase 25 scaffold files if P25-C successor (reviewer/planner rights on **upcoming** artifacts)
- Register `docs/TODO/phase-25.md` + index link
- Update board Notes: verdict, confidence, successor, residuals, next runnable **P25-00**

## Verdict

`APPROVE` | `REQUEST_CHANGES` — confidence **high** | **medium** | **low**

## Exit criteria

- [ ] Independent re-verify spot-checks PASS (minimum commands above)
- [ ] S05-01 evidence reviewed (archive + VERIFY-NOTES bar map)
- [ ] Completion bar Bars 1–6 satisfied
- [ ] S04→S05 residuals listed in REVIEW-NOTES / DR-HANDOFF
- [ ] No open blocker/high without pending spawn
- [ ] **`DR-HANDOFF.md` CLOSED** with successor **Phase 25 P25-C** (or documented **`no successor`**)
- [ ] Phase 25 scaffold **runnable** if P25-C (README, 00-PHASE-PLANNER, ≥1 scope stub, board P25-00)
- [ ] Phase 24 all rows `done` in board Notes
- [ ] Confidence **high** (or **medium** with explicit residuals)
- [ ] `docs/TODO.md` index updated if phase complete

## Forbidden

- Leaving successor **TBD** when row is `done`
- Closing DR-HANDOFF without independent spot-check re-run
- Rewriting Phase 23 or S01–S04 `done` history
- Implementing P25 product Go in this review row
- Scaffolding P25-A and P25-B as simultaneous active phases (queue only)

## Minimal todos

- [ ] Read VERIFY-NOTES + evidence manifest
- [ ] Re-run locked spot-check commands
- [ ] Execute review checklists (VERIFY, deliverables, handoff)
- [ ] Close DR-HANDOFF with locked successor template
- [ ] Scaffold Phase 25 (if P25-C) + register board
- [ ] Write REVIEW-NOTES.md; set row `done` with verdict

## Next

**P25-00** (after human/orchestrator promotes Phase 25) — or **none** if `no successor`
