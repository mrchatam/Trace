# P24-S03-02 — External research review

## Metadata
- id: P24-S03-02
- todo_ids: [P24-S03-02]
- role: reviewer
- skills: [code-review-and-quality, research, documentation-and-adrs, writing-for-agents]
- mcps: [Read, Glob, Grep, Shell, WebFetch]
- verification: manual (checklist + URL spot-check)
- hooks: none

## Objective

Independent review of S03-01 **`EXTERNAL-RESEARCH.md`** and FINDINGS external section. Verify sources are real; recommendations respect Trace laws; FM linkage and S02 residuals are addressed. Spawn fix rows only via board — do not rewrite S03-01 deliverable unless spawning.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md)
- S03-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- S03-01 prompt: [01-external-research.md](./01-external-research.md)
- SoT: [INVESTIGATION.md](../../INVESTIGATION.md), [FINDINGS.md](../../FINDINGS.md)
- Internal baseline: [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md), [CODEBASE-AUDIT.md](../scope-02-codebase-loop-audit/CODEBASE-AUDIT.md)

## Session start

Follow agent-loop-protocol. **Fresh reviewer context** — do not reuse S03-01 session. Board edits: **status + notes only**.

## Artifacts under review

| Artifact | Path |
|----------|------|
| External research | `scopes/scope-03-external-research/EXTERNAL-RESEARCH.md` |
| Living findings | [FINDINGS.md](../../FINDINGS.md) |

## Evidence to re-verify (reviewer runs selective checks)

| Check | Command / action | Pass criterion |
|-------|------------------|----------------|
| Local paths exist | `ls "similar projects/<name>"` for ≥2 cited local comparables | Paths match §2 table |
| URL spot-check | WebFetch **≥3** external URLs from §2 (SWE, OH, AID, CUR, or LIT) | HTTP 200 or valid redirect; content supports claim |
| Locked ID coverage | Compare §2 to [01-external-research.md](./01-external-research.md) locked list | UA, CG, GF, GT, CM, AR, SWE, OH, AID, CUR each present or §6 skip with reason |
| FM linkage | Scan §2 FM IDs column | Each FM cited ≥1× across doc OR §6 explains omission |
| S02 crosswalk | Read §5 | All 5 S02 residuals addressed |

## Review checklist — EXTERNAL-RESEARCH.md

### Structure (blockers)

- [ ] **Blocker:** Missing §2 comparable systems table (approach | transfer | risk columns)
- [ ] **Blocker:** Missing §3 research question answers (§3.1–§3.5)
- [ ] **Blocker:** Missing §4 anti-patterns
- [ ] **Blocker:** Missing §5 S02 residual crosswalk
- [ ] **Blocker:** Fewer than **6** comparables in §2
- [ ] **High:** Missing §1 executive summary or §6 open gaps

### Source integrity

- [ ] **Blocker:** Any URL appears fabricated (404, wrong domain, placeholder)
- [ ] **Blocker:** Fewer than **3** verified external URLs in §2
- [ ] **Blocker:** Fewer than **2** local `similar projects/` citations with valid paths
- [ ] **High:** Tier C literature missing (no paper/survey row)
- [ ] **High:** Marketing copy without mechanism detail in **approach** column
- [ ] **Medium:** Broken relative path to local clone
- [ ] **Low:** Access date not noted for web sources (optional nit)

### Trace law compliance

- [ ] **Blocker:** Recommends **hosted MCP or daemon on Trace P0-X core path** without “anti-pattern / harness-only” label
- [ ] **Blocker:** Recommends **full-graph dump** as default Trace context strategy
- [ ] **Blocker:** Recommends **full-rebuild-on-change** indexer architecture
- [ ] **High:** Transfer column ignores local-first / `.trace/` gitignored constraint
- [ ] **High:** Transfer column proposes autonomous task creation with no human authority note (violates G_PROJECT_LAWS)
- [ ] **Medium:** Conflates peer MCP (optional harness) with Trace product surface

### FM and investigation coverage

- [ ] **Blocker:** Q-D (forced replanning; plan-before-edit) not answered in §3
- [ ] **High:** FM-03, FM-10 not referenced in any comparable row
- [ ] **High:** FM-04 orchestrator bypass not addressed in §3.4
- [ ] **High:** §2 **tags** omit intervention category on any row
- [ ] **Medium:** FM IDs listed but mechanism mismatch vs POSTMORTEM §3

### S02 residual coverage

- [ ] **Blocker:** Sticky STOP / reason-code UX not in §5
- [ ] **Blocker:** Deliberation reset after gap pass not in §5
- [ ] **High:** Discovery-without-task promotion not compared to peers
- [ ] **High:** Orchestrator vs worker memory not in §3.4
- [ ] **Medium:** Plan-before-edit enforcement not contrasted with Trace install ([ENFORCEMENT.md](../../../phase-23-enforcement-choke-points/ENFORCEMENT.md))

### Scope hygiene

- [ ] **Blocker:** Product Go files in reviewer/implementer diff
- [ ] **High:** INTERVENTION-MATRIX ranking in external doc (S04 scope — §5 hints OK, full ranking not)
- [ ] **High:** Rewrites CODEBASE-AUDIT or POSTMORTEM as if S02/S01 never ran
- [ ] **Medium:** Duplicate full §2 table into FINDINGS (should be brief bullets + link)

## Review checklist — FINDINGS.md

- [ ] **Blocker:** External comparables row still `pending` after S03-01 claims done
- [ ] **High:** FINDINGS contradicts EXTERNAL-RESEARCH without reconciliation note
- [ ] **High:** Two-mode model flattened
- [ ] **Medium:** No link to EXTERNAL-RESEARCH.md path
- [ ] **Medium:** Summary duplicates entire §2 (should be 3–5 bullets)

## Cross-artifact consistency

- [ ] EXTERNAL-RESEARCH transfers align with CODEBASE-AUDIT change levers (no contradiction without explanation)
- [ ] Peer “task promotion” patterns explain E01 Session B (7 discoveries, 0 new tasks)
- [ ] Anti-patterns in §4 consistent with [project-rules.md](../../../../rules/project-rules.md) settled stack
- [ ] INVESTIGATION intervention categories used consistently in tags

## Spawn policy

- **blocker/high:** spawn `P24-S03-02a` (implement fix) + `02b` (re-review) immediately below this row; or inline doc fix if ≤15 lines and zero new research claims
- **medium:** prefer spawn unless typo-only
- Do not rewrite S03-00 / S03-01 `done` prompt bodies

## Verdict

`APPROVE` | `REQUEST_CHANGES` — confidence **high** | **medium** | **low**

Record in board Notes: verdict, confidence, residuals forwarded to S04.

## Exit criteria

- [ ] Checklists above executed; blockers resolved or forward row spawned
- [ ] Verdict + confidence in board Notes
- [ ] Residual risks listed (e.g. paywalled paper, stale peer version)
- [ ] No product Go in reviewer diff

## Minimal todos

- [ ] Verify ≥3 external URLs via WebFetch
- [ ] Confirm ≥2 local `similar projects/` paths exist
- [ ] Walk §5 S02 residual crosswalk
- [ ] Scan §4 for Trace law violations in *transfer* recommendations
- [ ] Check FINDINGS external section (brief + link)
- [ ] Set row done with verdict

## Next

**P24-S04-00**
