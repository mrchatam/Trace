# P24-S02-02 — Codebase audit review

## Metadata
- id: P24-S02-02
- todo_ids: [P24-S02-02]
- role: reviewer
- skills: [code-review-and-quality, debugging-and-error-recovery, documentation-and-adrs, writing-for-agents]
- mcps: [Read, Glob, Grep, Shell, user-codegraph]
- verification: manual (checklist)
- hooks: none

## Objective

Independent review of S02-01 **`CODEBASE-AUDIT.md`** and FINDINGS codebase section. Verify file:line claims against live repo; confirm S01 residuals reconciled. Spawn fix rows only via board — do not rewrite S02-01 deliverable unless spawning.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- S02-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- S02-01 prompt: [01-codebase-loop-audit.md](./01-codebase-loop-audit.md)
- S01 evidence: [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md)
- SoT: [INVESTIGATION.md](../../INVESTIGATION.md), [FINDINGS.md](../../FINDINGS.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — do not reuse S02-01 session. Board edits: **status + notes only**.

## Artifacts under review

| Artifact | Path |
|----------|------|
| Codebase audit | `scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md` |
| Living findings | [FINDINGS.md](../../FINDINGS.md) |

## Evidence to re-verify (reviewer runs selective checks)

| Check | Command / path | Pass criterion |
|-------|----------------|----------------|
| SelectNext order | `internal/deliberation/select.go` L7–12 | Audit matches: hop_budget checked before p19_saturated |
| Hop budget constant | `internal/deliberation/types.go` | HopBudget=12 cited correctly if referenced |
| p19 saturation | `internal/loop/policy.go` L103–109 | Audit matches `p19SaturatedFromLastStep` logic |
| trace_add surface | `internal/mcp/server.go`, `tools_write.go` | FM-08 claims match tool descriptions |
| Spawn path | `internal/loop/apply.go` spawned_tasks | FM-10 distinguishes apply spawn vs trace add |
| Install scope | `internal/install/enforcement.go` | FM-04/05 claims — no false “product supports gap pass” |
| Live gate (optional) | `trace loop gate --task …0050 --for edit` on G1 | Audit acknowledges task_not_found vs export STOP |
| FM count | CODEBASE-AUDIT §2 | ≥8 rows with file:line |

## Review checklist — CODEBASE-AUDIT.md

### Structure (blockers)

- [ ] **Blocker:** Missing §2 FM mechanism table (five columns)
- [ ] **Blocker:** Missing §3 S01 residual reconciliation (reason codes, export/DB, deliberation reset)
- [ ] **Blocker:** Fewer than **8** FM rows with file:line
- [ ] **High:** Missing §1 summary or §4 cross-cutting observations

### Citation accuracy

- [ ] **Blocker:** Spot-check **≥3** file:line citations — line numbers exist and support the claim
- [ ] **Blocker:** Any FM row with empty file:line (except FM-06 labeled protocol with explicit “no product path”)
- [ ] **High:** SelectNext / gate semantics misstated (e.g. claims hop_budget at hop_count=1 when export shows p19_saturated — must explain via policy inputs)
- [ ] **High:** Confuses `cmd/trace-mcp/` with `internal/mcp/` handler location
- [ ] **Medium:** Stale path (file moved/renamed since audit written)
- [ ] **Low:** Typo in FM-ID or UUID

### S01 residual coverage

- [ ] **Blocker:** `hop_budget_exceeded` vs `p19_saturated` not reconciled with E01 hop_count=1 evidence
- [ ] **High:** Export vs SQLite sync not addressed (G1 task_not_found vs graph.json)
- [ ] **High:** Deliberation reset after gap pass — no finding (must cite code path or document absence)
- [ ] **Medium:** POSTMORTEM §4 open questions not reflected in §6 or marked deferred

### Change lever hygiene

- [ ] **Blocker:** Any row with vague lever (“improve UX”, “better docs” without product/harness/protocol/experiment)
- [ ] **High:** Product lever proposed for Multitask orchestrator split (FM-04) without harness/protocol label
- [ ] **High:** Harness-only behavior attributed to Trace library bug without evidence
- [ ] **Medium:** FM-06/07 not labeled protocol/experiment where code path absent

### Scope hygiene

- [ ] **Blocker:** Product Go files in reviewer/implementer diff
- [ ] **High:** INTERVENTION-MATRIX ranking in audit (S04 scope)
- [ ] **Medium:** Rewrites POSTMORTEM Session A/B narrative

## Review checklist — FINDINGS.md

- [ ] **Blocker:** Codebase audit row still `pending` after S02-01 claims done
- [ ] **High:** FINDINGS contradicts CODEBASE-AUDIT without note
- [ ] **High:** Two-mode model flattened or Session A/B conflated
- [ ] **Medium:** No link to CODEBASE-AUDIT.md path
- [ ] **Medium:** Summary bullets duplicate full §2 table (should be brief)

## Cross-artifact consistency

- [ ] CODEBASE-AUDIT FM rows align with POSTMORTEM §3 confirmed/partial statuses (no contradiction without explanation)
- [ ] FM-03 mechanism explains both Session A early STOP and Session B post-fix STOP
- [ ] FM-10 mechanism explains discovery + discovery_mentions_task without new task UUIDs
- [ ] INVESTIGATION.md investigation questions B (product/policy) touched in §4

## Spawn policy

- **blocker/high:** spawn `P24-S02-02a` (implement fix) + `02b` (re-review) immediately below this row; or inline doc fix if ≤10 lines and zero new code claims
- **medium:** prefer spawn unless typo-only
- Do not rewrite S02-00 / S02-01 `done` prompt bodies

## Verdict

`APPROVE` | `REQUEST_CHANGES` — confidence **high** | **medium** | **low**

Record in board Notes: verdict, confidence, residuals forwarded to S03/S04.

## Exit criteria

- [ ] Checklists above executed; blockers resolved or forward row spawned
- [ ] Verdict + confidence in board Notes
- [ ] Residual risks listed (e.g. codegraph index stale, G1 DB not inspected)
- [ ] No product Go in reviewer diff

## Minimal todos

- [ ] Spot-check ≥3 CODEBASE-AUDIT file:line citations against live files
- [ ] Verify SelectNext + p19 saturation claims (`select.go`, `policy.go`)
- [ ] Confirm §3 covers all S01 forwarded residuals
- [ ] Walk FM-03, FM-08, FM-10 rows in detail
- [ ] Set row done with verdict

## Next

**P24-S03-00**
