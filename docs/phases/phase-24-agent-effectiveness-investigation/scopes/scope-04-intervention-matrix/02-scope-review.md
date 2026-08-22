# P24-S04-02 — Intervention matrix review

## Metadata
- id: P24-S04-02
- todo_ids: [P24-S04-02]
- role: reviewer
- skills: [code-review-and-quality, documentation-and-adrs, analyst, writing-for-agents]
- mcps: [Read, Glob, Grep]
- verification: manual (checklist)
- hooks: none

## Objective

Independent review of S04-01 synthesis: **`INTERVENTION-MATRIX.md`**, consolidated **`FINDINGS.md`**, and **`DR-HANDOFF.md`** theme updates. Verify ranking evidence, Phase 25 scope boundaries, and Trace law compliance. Spawn fix rows only via board — do not rewrite S04-01 deliverables unless spawning.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md)
- S04-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- S04-01 prompt: [01-intervention-matrix.md](./01-intervention-matrix.md)
- SoT: [INVESTIGATION.md](../../INVESTIGATION.md), [FINDINGS.md](../../FINDINGS.md)
- Evidence base: [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md), [CODEBASE-AUDIT.md](../scope-02-codebase-loop-audit/CODEBASE-AUDIT.md), [EXTERNAL-RESEARCH.md](../scope-03-external-research/EXTERNAL-RESEARCH.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — do not reuse S04-01 session. Board edits: **status + notes only**.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Intervention matrix | `scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md` |
| Living findings | [FINDINGS.md](../../FINDINGS.md) |
| Phase handoff | [DR-HANDOFF.md](../../DR-HANDOFF.md) |

## Evidence to re-verify (reviewer spot-checks)

| Check | Source | Pass criterion |
|-------|--------|----------------|
| Session B task count | POSTMORTEM Must answer §2; G1 graph.json | 5 seed tasks only — supports FM-10 promotion rows |
| Sticky STOP | CODEBASE-AUDIT §3; S03 §5 | INT row for reason-code UX cites S02+S03 |
| Deliberation reset absence | CODEBASE-AUDIT §3 | Reset intervention cites absence or defers with spike |
| spawned_tasks path | CODEBASE-AUDIT §2 FM-10 | Promotion row distinguishes apply spawn vs trace add |
| Peer auto-spawn | EXTERNAL-RESEARCH §3.3, §5 | FM-10 rows reference AR/SWE pattern with human gate |
| Two-mode model | POSTMORTEM §2; FINDINGS | Not flattened in executive summary |
| Trace law | project-rules, G_PROJECT_LAWS | No daemon-core recommendations in top 3 ranks |

## Review checklist — INTERVENTION-MATRIX.md

### Structure (blockers)

- [ ] **Blocker:** Missing §2 ranked table with locked columns (Rank, ID, Addresses, Intervention, Owner, Impact, Effort, Risk, Evidence, Phase 25 theme)
- [ ] **Blocker:** Fewer than **8** intervention rows
- [ ] **Blocker:** Missing §1 ranking rationale or §3 FM coverage matrix
- [ ] **Blocker:** Missing §5 Phase 25 theme mapping
- [ ] **High:** Missing §4 deferred/human-gate items

### Locked schema compliance

- [ ] **Blocker:** Any row missing Owner, Impact, Effort, Risk, or Evidence
- [ ] **Blocker:** Owner value outside `product|harness|protocol|docs`
- [ ] **Blocker:** Impact/Effort outside `high|med|low`
- [ ] **Blocker:** Risk outside `regression|agent confusion|scope creep`
- [ ] **High:** Intervention column >1 sentence or vague (“improve loop” without concrete delta)
- [ ] **High:** Non-contiguous INT IDs
- [ ] **Medium:** Duplicate interventions without differentiation in Notes

### Ranking evidence (S01 / S02 / S03)

- [ ] **Blocker:** Any **high** Impact row lacks S01/S02/S03 Evidence pointer
- [ ] **Blocker:** Top 3 ranks lack documented rubric rationale in §1
- [ ] **High:** Harness-only doc tweak ranked above product behavior change without dogfood feasibility justification
- [ ] **High:** Row claims fix FM-* without §3 coverage entry
- [ ] **High:** Evidence cites wrong scope (e.g. “S04” only — must be S01/S02/S03)
- [ ] **Medium:** Spot-check **≥3** Evidence cells — artifact section exists and supports claim

### Must-include seeds (from 01-intervention-matrix.md)

- [ ] **Blocker:** Missing row for discovery → task promotion (FM-10)
- [ ] **Blocker:** Missing row for hop budget / verify recovery (FM-03)
- [ ] **High:** Missing default gap pass / Mode A→B collapse (FM-09)
- [ ] **High:** Missing orchestrator Trace-first (FM-04/05)
- [ ] **High:** Missing saturation/hop policy calibration
- [ ] **Medium:** Missing experiment protocol / score.sh fix
- [ ] **Medium:** Missing two-session dogfood rubric

### S03→S04 forwarded residuals

- [ ] **High:** Hop/P19 calibration not ranked or §4 deferred with rationale
- [ ] **High:** Auto-spawn vs loop apply human gate not documented
- [ ] **High:** Sticky STOP UX + deliberation reset addressed (row or §4)
- [ ] **High:** Build-mode vs directed-gap defaults (FM-09) in top half of rank or §1 called out
- [ ] **Medium:** SQLite episode spike bounded (not full Graphiti port)
- [ ] **Medium:** Cursor hook API drift noted for Phase 25 harness spike

### Phase 25 scope hygiene

- [ ] **Blocker:** §5 or DR-HANDOFF recommends **all** P25-A…E as one mega-phase
- [ ] **Blocker:** More than **3** themes marked **Recommended** without human-split note
- [ ] **High:** Recommended themes lack INT-ID evidence links
- [ ] **High:** Theme scope boundaries missing “out of scope” column in §5
- [ ] **Medium:** Single intervention row claims to fix all FMs

### Trace law compliance

- [ ] **Blocker:** Top-3 intervention requires daemon/HTTP on P0-X core path
- [ ] **Blocker:** Recommends full-rebuild-on-change indexer
- [ ] **High:** Graphiti/AgentRQ hosted MCP promoted as Trace core (S03 §4 anti-pattern violated)
- [ ] **High:** Auto-spawn without human approval gate documented
- [ ] **Medium:** Portable graph / seed import discipline omitted from protocol rows

### Scope hygiene

- [ ] **Blocker:** Product Go files in reviewer/implementer diff
- [ ] **High:** Rewrites POSTMORTEM/CODEBASE-AUDIT/EXTERNAL-RESEARCH bodies (should link only)
- [ ] **Medium:** Full S02 mechanism table duplicated in matrix (summary OK)

## Review checklist — FINDINGS.md

- [ ] **Blocker:** Status table — any section still `pending` after S04-01 claims done
- [ ] **Blocker:** Missing executive summary (5–8 sentences)
- [ ] **Blocker:** Intervention matrix row still `pending`
- [ ] **High:** Session A/B conflated in synthesis
- [ ] **High:** Contradicts INTERVENTION-MATRIX top ranks without note
- [ ] **High:** Missing link to INTERVENTION-MATRIX.md
- [ ] **Medium:** Duplicates full §2 matrix table (should be top-3 + link)

## Review checklist — DR-HANDOFF.md

- [ ] **Blocker:** S04 checklist item still open when S04-01 done
- [ ] **High:** Candidate themes table not updated with Recommended/deferred labels
- [ ] **High:** Recommended themes missing INT evidence links
- [ ] **Medium:** DR-HANDOFF marked CLOSED (S05-02 owns close)
- [ ] **Medium:** Successor decision still TBD without intervention summary bullets

## Cross-artifact consistency

- [ ] INTERVENTION-MATRIX §3 FM coverage aligns with POSTMORTEM §3 confirmed statuses
- [ ] FM-03 interventions explain both Session A early STOP and Session B post-fix STOP
- [ ] FM-10 interventions explain 7 discoveries + 0 new tasks
- [ ] FINDINGS executive summary matches §1 matrix rationale
- [ ] DR-HANDOFF recommended themes ⊆ INTERVENTION-MATRIX §5
- [ ] INVESTIGATION intervention categories reflected in Owner column distribution

## Spawn policy

- **blocker/high:** spawn `P24-S04-02a` (implement fix) + `02b` (re-review) immediately below this row; or inline doc fix if ≤15 lines and zero new research claims
- **medium:** prefer spawn unless typo-only
- Do not rewrite S04-00 / S04-01 `done` prompt bodies

## Verdict

`APPROVE` | `REQUEST_CHANGES` — confidence **high** | **medium** | **low**

Record in board Notes: verdict, confidence, residuals forwarded to S05.

## Exit criteria

- [ ] Checklists above executed; blockers resolved or forward row spawned
- [ ] Verdict + confidence in board Notes
- [ ] Residual risks listed (e.g. human product call on auto-spawn, spike scope)
- [ ] No product Go in reviewer diff

## Minimal todos

- [ ] Verify INTERVENTION-MATRIX §2 schema + row count (≥8)
- [ ] Spot-check ≥3 Evidence cells against S01/S02/S03
- [ ] Walk top 3 ranks vs ranking rubric in §1
- [ ] Confirm FINDINGS executive summary + status table
- [ ] Confirm DR-HANDOFF 1–3 Recommended themes with INT links
- [ ] Scan top 3 for Trace law violations
- [ ] Set row done with verdict

## Next

**P24-S05-00**
