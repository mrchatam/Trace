# P24-S01-02 — Dogfood post-mortem review

## Metadata
- id: P24-S01-02
- todo_ids: [P24-S01-02]
- role: reviewer
- skills: [code-review-and-quality, documentation-and-adrs, writing-for-agents]
- mcps: [Read, Glob, Grep, Shell]
- verification: manual (checklist)

## Objective

Independent review of S01-01 artifacts (`POSTMORTEM.md`, `FINDINGS.md`). Spawn S02 gaps only via board notes — do not rewrite S01-01 deliverables unless spawning a fix row.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- S01-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- S01-01 prompt: [01-dogfood-postmortem.md](./01-dogfood-postmortem.md)
- SoT: [INVESTIGATION.md](../../INVESTIGATION.md), [FINDINGS.md](../../FINDINGS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — do not reuse S01-01 session. Board edits: **status + notes only**.

## Review checklist (summary)

Execute every item in the three checklists below before verdict:

1. **POSTMORTEM.md** — structure blockers, evidence quality, must-answer coverage, scope hygiene
2. **FINDINGS.md** — two-mode model preserved; FM per-session statuses; link to POSTMORTEM
3. **Cross-artifact** — POSTMORTEM §2 ↔ FINDINGS tables; FM-09/FM-10 explicit in both files

## Artifacts under review

| Artifact | Path |
|----------|------|
| Post-mortem | `scopes/scope-01-dogfood-postmortem/POSTMORTEM.md` |
| Living findings | [FINDINGS.md](../../FINDINGS.md) |

## Evidence to re-verify (reviewer runs selective checks)

| Check | Command / path | Pass criterion |
|-------|----------------|----------------|
| G1≡B0 (Session A) | `diff -rq runs/B0/internal runs/G1/internal` | POSTMORTEM claim matches diff (identical or lists deltas) |
| Session B divergence | `git -C runs/G1 diff f70aaea..704e2ff --stat` | POSTMORTEM cites gap-closure files |
| Graph counts | `runs/G1/trace/graph.json` | 5 tasks, 7 discoveries, 2 decisions, 4 evidence (or corrected with cite) |
| Verify gate | `trace loop status/gate --task …0050` | POSTMORTEM reconciles `hop_budget_exceeded` vs `p19_saturated` |
| Git range | `git -C runs/G1 log --oneline` | Session B commits `704e2ff`..`a37e7c0` referenced |

## Review checklist — POSTMORTEM.md

### Structure (blockers)

- [ ] **Blocker:** Missing §1 Runs reviewed (Session A + B required)
- [ ] **Blocker:** Missing §2 two-mode comparison table
- [ ] **Blocker:** Missing §3 FM-01..FM-10 matrix with Session A / B / Both columns
- [ ] **Blocker:** Missing §4 open questions
- [ ] **Blocker:** Sessions A and B conflated into single narrative

### Evidence quality

- [ ] **Blocker:** Fewer than **8 distinct** cited evidence paths
- [ ] **Blocker:** Any FM row with Status ≠ unknown but empty Evidence cell
- [ ] **High:** FM claim contradicts `graph.json` without explanation
- [ ] **High:** Session A claims lack B0 diff or RESULTS/prompt cite
- [ ] **High:** Session B claims lack git or graph.json cite
- [ ] **Medium:** Optional D44/D45 mentioned but uncited when marked “reviewed”
- [ ] **Medium:** VERIFY.md manual table not cited for Session B test pass
- [ ] **Low:** Typo in task UUIDs (`…0010`–`…0050` seed set)

### Must-answer coverage (from 01 prompt)

- [ ] **Blocker:** Discovery count vs zero new tasks not explained (even if “unknown” with evidence gap noted)
- [ ] **High:** `trace add` task vs discovery usage not addressed
- [ ] **High:** Verify task gate block after green tests not addressed
- [ ] **High:** FM Mode A / B / both classification missing or inconsistent with §3 table
- [ ] **High:** Prompt hack vs product entry point not addressed

### Scope hygiene

- [ ] **Blocker:** Product Go files in diff
- [ ] **High:** Codebase audit implementation detail (S02 scope) presented as fact without “hypothesis” label
- [ ] **Medium:** Intervention recommendations (S04 scope) in POSTMORTEM

## Review checklist — FINDINGS.md

- [ ] **Blocker:** Two-mode model section removed or flattened to single session
- [ ] **Blocker:** FM taxonomy not updated (still `pending` with no per-session status)
- [ ] **High:** FINDINGS contradicts POSTMORTEM without cross-reference
- [ ] **High:** FM status uses only global confirmed/partial — missing Session A / B split
- [ ] **Medium:** Preliminary conclusion unchanged despite new POSTMORTEM evidence (should refine or note “unchanged”)
- [ ] **Medium:** No link from FINDINGS to POSTMORTEM.md path

## Cross-artifact consistency

- [ ] POSTMORTEM §2 aligns with FINDINGS Session A / B tables
- [ ] FM-09 (mode-dependent effectiveness) and FM-10 (discovery without task promotion) explicitly addressed in both files
- [ ] INVESTIGATION.md seed FM symptoms referenced — deviations documented

## Spawn policy

- **blocker/high:** spawn `P24-S01-02a` (implement fix) + `02b` (re-review) immediately below this row; or inline doc fix if ≤10 lines and zero new claims
- **medium:** prefer spawn unless typo-only
- Do not rewrite S01-00 / S01-01 `done` prompt bodies

## Verdict

`APPROVE` | `REQUEST_CHANGES` — confidence **high** | **medium** | **low**

Record in board Notes: verdict, confidence, open residuals forwarded to S02/S04.

## Exit criteria

- [ ] Checklist above executed; blockers resolved or forward row spawned
- [ ] Verdict + confidence in board Notes
- [ ] Residual risks listed (e.g. Session A sparse git → S02 may need transcript)
- [ ] No product Go in reviewer diff

## Minimal todos

- [ ] Spot-check ≥3 POSTMORTEM citations (graph, git, diff)
- [ ] Walk FM-01..FM-10 row-by-row vs INVESTIGATION seed
- [ ] Confirm FINDINGS FM section has per-session statuses
- [ ] Set row done with verdict

## Next

**P24-S02-00**
