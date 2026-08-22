# P24-S05-01 — Phase 24 verify

## Metadata
- id: P24-S05-01
- todo_ids: [P24-S05-01]
- role: verify
- skills: [documentation-and-adrs, writing-for-agents, planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [analyst]
- verification: manual (artifact checklist + evidence archive)
- hooks: none

## Objective

Verify Phase 24 meets the **INVESTIGATION.md / README completion bar**; archive immutable evidence copies (or hashes) of all scope deliverables + consolidated `FINDINGS.md`; write **`VERIFY-NOTES.md`** mapping each criterion to repo paths. **Does not** close DR-HANDOFF (S05-02 owns).

**Investigation phase** — no product Go unless a blocking doc gap requires a ≤15-line fix (prefer `failed` + spawn).

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- S05-00 locks: [00-PLANNER.md](00-PLANNER.md)
- SoT: [INVESTIGATION.md](../../INVESTIGATION.md) — completion bar, session modes, FM taxonomy
- [README.md](../../README.md) — completion bar (5 numbered items)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S05-02
- S01: [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md)
- S02: [CODEBASE-AUDIT.md](../scope-02-codebase-loop-audit/CODEBASE-AUDIT.md)
- S03: [EXTERNAL-RESEARCH.md](../scope-03-external-research/EXTERNAL-RESEARCH.md)
- S04: [INTERVENTION-MATRIX.md](../scope-04-intervention-matrix/INTERVENTION-MATRIX.md)
- Living: [FINDINGS.md](../../FINDINGS.md)
- Pattern: [P23 S06-01 verify](../../../phase-23-enforcement-choke-points/scopes/scope-06-phase-verify/01-verify.md) (evidence archive shape)

## Session start

Follow agent-loop-protocol. Unattended: do not stop after planning. This row verifies documentation completeness and archives evidence; it does not implement Phase 25 or close DR-HANDOFF.

## Locked defaults (FINAL — S05-00)

| Item | Value |
|------|-------|
| Evidence dir | `experiments/runs/YYYY-MM-DD-p24-s05-01-verify/evidence/` (replace `YYYY-MM-DD` with run date) |
| Archive method | **Copy** each artifact into evidence dir **and** write `manifest.sha256` (hashes of source paths at verify time) |
| Artifacts (6) | `FINDINGS.md`, `POSTMORTEM.md`, `CODEBASE-AUDIT.md`, `EXTERNAL-RESEARCH.md`, `INTERVENTION-MATRIX.md`, `DR-HANDOFF.md` |
| Notes artifact | `scopes/scope-05-phase-verify/VERIFY-NOTES.md` (**required**) |
| Product Go | **Forbidden** (investigation phase) |
| DR-HANDOFF | Stays **OPEN** — S05-02 closes |
| Phase 25 scaffold | **Out of scope** — S05-02 owns per agent-loop-protocol Phase handoff |

### Source paths (locked — copy from these)

| Archive name | Source path (repo root relative) |
|--------------|----------------------------------|
| `FINDINGS.md` | `docs/phases/phase-24-agent-effectiveness-investigation/FINDINGS.md` |
| `POSTMORTEM.md` | `docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-01-dogfood-postmortem/POSTMORTEM.md` |
| `CODEBASE-AUDIT.md` | `docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md` |
| `EXTERNAL-RESEARCH.md` | `docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-03-external-research/EXTERNAL-RESEARCH.md` |
| `INTERVENTION-MATRIX.md` | `docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md` |
| `DR-HANDOFF.md` | `docs/phases/phase-24-agent-effectiveness-investigation/DR-HANDOFF.md` |

## Must checklist — INVESTIGATION / README completion bar

S05-01 Notes and `VERIFY-NOTES.md` must cite **PASS/FAIL** evidence for each row:

### Bar 1 — Two-mode model + failure modes (`FINDINGS.md`)

- [ ] **Bar 1a:** `FINDINGS.md` documents **Session A vs Session B** as separate evidence rows (not conflated) with cited artifacts (graph.json, git log, POSTMORTEM paths)
- [ ] **Bar 1b:** ≥**5** distinct failure modes (`FM-*`) with E01 evidence (Session A and/or B columns or narrative)
- [ ] **Bar 1c:** `FINDINGS.md` status table — all sections **done** (not `pending`)
- [ ] **Bar 1d:** Executive summary present (5–8 sentences); links to INTERVENTION-MATRIX top-3

### Bar 2 — Codebase audit (`CODEBASE-AUDIT.md`)

- [ ] **Bar 2a:** Maps **FM-01..FM-10** (or justified subset with “N/A” rationale) to **file:line** or API symbols
- [ ] **Bar 2b:** §1–§6 structure present (mechanism inventory, FM table, S01 residuals, levers, MCP/CLI, summary)
- [ ] **Bar 2c:** S01 forwarded residuals addressed (SelectNext reason codes, export/DB drift, deliberation reset absence)

### Bar 3 — External research (`EXTERNAL-RESEARCH.md`)

- [ ] **Bar 3a:** ≥**3** comparables with **URLs** (Tier B or arXiv) and actionable deltas for Trace
- [ ] **Bar 3b:** §1–§6 structure; Trace law anti-patterns in §4
- [ ] **Bar 3c:** S02 residuals in §5 (sticky STOP, deliberation reset, discovery→task, orchestrator)

### Bar 4 — Intervention matrix (`INTERVENTION-MATRIX.md`)

- [ ] **Bar 4a:** ≥**8** ranked rows (`INT-01` … contiguous)
- [ ] **Bar 4b:** Locked columns: Rank, ID, Addresses, Intervention, Owner, Impact, Effort, Risk, Evidence, Phase 25 theme
- [ ] **Bar 4c:** §1 ranking rationale; §3 FM coverage; §4 deferred/human-gate; §5 Phase 25 theme mapping
- [ ] **Bar 4d:** FM-01..FM-10 each addressed in §3 (≥1 INT per FM or explicit gap note)

### Bar 5 — DR-HANDOFF Phase 25 themes

- [ ] **Bar 5a:** Lists **1–3** **Recommended** Phase 25 implementation themes (not a flat mega-backlog)
- [ ] **Bar 5b:** Recommended themes have INT-ID evidence links (from matrix)
- [ ] **Bar 5c:** Successor decision **not** TBD — candidate table populated (CLOSED status waits for S05-02)
- [ ] **Bar 5d:** Locked recommendation order: **P25-C → P25-A → P25-B**; P25-D/E deferred with rationale

### Bar 6 — Scope deliverables present (S01–S04)

- [ ] **Bar 6a:** [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md) §1–§4 complete (E01 A+B, FM matrix, must-answer)
- [ ] **Bar 6b:** [CODEBASE-AUDIT.md](../scope-02-codebase-loop-audit/CODEBASE-AUDIT.md) exists with §1–§6
- [ ] **Bar 6c:** [EXTERNAL-RESEARCH.md](../scope-03-external-research/EXTERNAL-RESEARCH.md) exists with §1–§6
- [ ] **Bar 6d:** [INTERVENTION-MATRIX.md](../scope-04-intervention-matrix/INTERVENTION-MATRIX.md) exists with §1–§5

### S04→S05 forwarded residuals (verify documented — not implementation)

Each must appear in `VERIFY-NOTES.md` § Residuals with pointer to matrix §4 or FINDINGS:

| Residual | Verify expectation |
|----------|-------------------|
| **Auto-spawn human gate** | INTERVENTION-MATRIX §4 documents human product call; INT-01 is guided promotion only |
| **P19 threshold dogfood validation** | INT-02/INT-05 ranked or §4 deferred; FINDINGS notes recalibration needs live dogfood in Phase 25 |
| **Hook API drift (INT-11)** | INT-11 row present; P25-C theme tag |
| **Live gate `reason_code` env-dependency** | POSTMORTEM/FINDINGS/CODEBASE-AUDIT reconcile `task_not_found` vs `hop_budget_exceeded` with `-C` context |

## Locked verify spot-check commands

Run from repo root. Capture output snippets in evidence dir (`spot-checks.txt`) or cite in VERIFY-NOTES.

```bash
# --- Archive setup ---
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p24-s05-01-verify/evidence"
mkdir -p "$EVID"

# --- Copy artifacts ---
cp docs/phases/phase-24-agent-effectiveness-investigation/FINDINGS.md "$EVID/"
cp docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-01-dogfood-postmortem/POSTMORTEM.md "$EVID/"
cp docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md "$EVID/"
cp docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-03-external-research/EXTERNAL-RESEARCH.md "$EVID/"
cp docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md "$EVID/"
cp docs/phases/phase-24-agent-effectiveness-investigation/DR-HANDOFF.md "$EVID/"

# --- Manifest (hashes of sources at verify time) ---
(
  cd docs/phases/phase-24-agent-effectiveness-investigation
  sha256sum FINDINGS.md DR-HANDOFF.md
  sha256sum scopes/scope-01-dogfood-postmortem/POSTMORTEM.md
  sha256sum scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md
  sha256sum scopes/scope-03-external-research/EXTERNAL-RESEARCH.md
  sha256sum scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md
) > "$EVID/manifest.sha256"

# --- Metadata ---
{
  echo "verify_id=P24-S05-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
} > "$EVID/99-run-metadata.txt"
```

### Content spot-checks (minimum — append to `spot-checks.txt`)

```bash
# FM count in POSTMORTEM §3
grep -c '| FM-' docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-01-dogfood-postmortem/POSTMORTEM.md

# INT row count in matrix §2
grep -c '| INT-' docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md

# External URLs (http/https)
grep -cE 'https?://' docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-03-external-research/EXTERNAL-RESEARCH.md

# FINDINGS status — no pending sections
grep -i pending docs/phases/phase-24-agent-effectiveness-investigation/FINDINGS.md || true

# DR-HANDOFF recommended themes
grep -E 'P25-[ABC]' docs/phases/phase-24-agent-effectiveness-investigation/DR-HANDOFF.md

# CODEBASE-AUDIT file:line cites (sample)
grep -cE '`[^`]+\.go' docs/phases/phase-24-agent-effectiveness-investigation/scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md
```

**Pass thresholds (locked):**

| Check | Minimum |
|-------|---------|
| FM rows in POSTMORTEM §3 | **10** (`FM-01` … `FM-10`) |
| INT rows in INTERVENTION-MATRIX §2 | **8** (S04 delivered **11**) |
| URLs in EXTERNAL-RESEARCH | **3** (S04 delivered **11** comparables) |
| DR-HANDOFF Recommended themes | **1–3** with INT links |
| FINDINGS failure modes with evidence | **≥5** distinct FM-* |

## VERIFY-NOTES.md template (required output)

Write to `scopes/scope-05-phase-verify/VERIFY-NOTES.md`:

```markdown
# Phase 24 VERIFY notes

**Run:** YYYY-MM-DD  
**Row:** P24-S05-01  
**Evidence dir:** experiments/runs/YYYY-MM-DD-p24-s05-01-verify/evidence/  
**Git SHA:** `<rev-parse HEAD>`

## Verdict

PASS | FAIL — confidence high | medium | low

## Completion bar map

| Bar | Criterion | Result | Evidence |
|-----|-----------|--------|----------|
| 1a | Session A vs B separate | PASS/FAIL | FINDINGS.md §…; POSTMORTEM §1–§2 |
| 1b | ≥5 FMs with E01 evidence | PASS/FAIL | FINDINGS.md FM table; POSTMORTEM §3 |
| 1c | FINDINGS status all done | PASS/FAIL | FINDINGS status table |
| 1d | Executive summary + top-3 link | PASS/FAIL | FINDINGS.md L… |
| 2a | FM→file:line map | PASS/FAIL | CODEBASE-AUDIT §2 |
| 2b | §1–§6 structure | PASS/FAIL | CODEBASE-AUDIT headings |
| 2c | S01 residuals | PASS/FAIL | CODEBASE-AUDIT §3 |
| 3a | ≥3 comparables + URLs | PASS/FAIL | EXTERNAL-RESEARCH §2 |
| 3b | §1–§6 + anti-patterns | PASS/FAIL | EXTERNAL-RESEARCH §4 |
| 3c | S02 residuals | PASS/FAIL | EXTERNAL-RESEARCH §5 |
| 4a | ≥8 INT rows | PASS/FAIL | INTERVENTION-MATRIX §2 count |
| 4b | Locked schema | PASS/FAIL | INTERVENTION-MATRIX §2 header |
| 4c | §1/§3/§4/§5 present | PASS/FAIL | INTERVENTION-MATRIX |
| 4d | FM coverage | PASS/FAIL | INTERVENTION-MATRIX §3 |
| 5a | 1–3 Recommended themes | PASS/FAIL | DR-HANDOFF table |
| 5b | INT evidence links | PASS/FAIL | DR-HANDOFF INT-… |
| 5c | Successor not TBD | PASS/FAIL | DR-HANDOFF (OPEN OK) |
| 5d | P25-C→A→B order | PASS/FAIL | DR-HANDOFF |
| 6a–6d | Scope artifacts | PASS/FAIL | paths |

## Artifact inventory

| File | Archived | SHA256 (prefix) |
|------|----------|-----------------|
| FINDINGS.md | yes/no | … |
| POSTMORTEM.md | yes/no | … |
| CODEBASE-AUDIT.md | yes/no | … |
| EXTERNAL-RESEARCH.md | yes/no | … |
| INTERVENTION-MATRIX.md | yes/no | … |
| DR-HANDOFF.md | yes/no | … |

## S04→S05 residuals (documented)

| Residual | Disposition | Pointer |
|----------|-------------|---------|
| Auto-spawn human gate | … | INTERVENTION-MATRIX §4 |
| P19 threshold dogfood validation | … | INT-02, INT-05 |
| Hook API drift (INT-11) | … | INT-11, P25-C |
| Live gate reason_code env-dependency | … | POSTMORTEM / CODEBASE-AUDIT §3 |

## Spot-check command output

(paste key counts from spot-checks.txt)

## Gaps / spawn recommendation

(none | P24-S05-01a implement + 01b review for …)

## DR-HANDOFF status

**OPEN** — S05-02 closes with Phase 25 scaffold + successor decision.
```

## Do not

- Close [DR-HANDOFF.md](../../DR-HANDOFF.md) — S05-02 owns
- Mark criteria `done` without evidence dir + VERIFY-NOTES map
- Implement Phase 25 product Go or scaffold Phase 25 folder (S05-02)
- Rewrite S01–S04 `done` deliverable bodies (link only)
- Ship product Go under P24 investigate rows

## Exit criteria

- [ ] Evidence directory populated (`manifest.sha256`, 6 copies, `99-run-metadata.txt`)
- [ ] Must checklist **Bars 1–6** mapped in VERIFY-NOTES with PASS/FAIL
- [ ] S04→S05 residuals documented in VERIFY-NOTES § Residuals
- [ ] Spot-check thresholds met (or row `failed` with gap list)
- [ ] Board Notes cite evidence dir path + verdict summary
- [ ] DR-HANDOFF remains **OPEN**

## Minimal todos

- [ ] Preflight: confirm all 6 source artifacts exist
- [ ] Run archive commands; write `manifest.sha256` + metadata
- [ ] Run spot-check commands; append to `spot-checks.txt`
- [ ] Walk Must checklist Bars 1–6; fill VERIFY-NOTES template
- [ ] Document S04→S05 residuals
- [ ] Set row `done` or `failed` — **do not** close DR-HANDOFF

## Next

**P24-S05-02**
