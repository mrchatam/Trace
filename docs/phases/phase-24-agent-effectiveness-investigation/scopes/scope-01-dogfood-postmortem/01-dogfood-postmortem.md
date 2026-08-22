# P24-S01-01 — Dogfood post-mortem investigation

## Metadata
- id: P24-S01-01
- todo_ids: [P24-S01-01]
- role: implementer
- skills: [research, code-explorer, documentation-and-adrs, grilling, writing-for-agents]
- mcps: [Read, Glob, Grep, Write, user-trace]
- agents: [explore, investigator]
- verification: manual (evidence cites)
- hooks: none

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [INVESTIGATION.md](../../INVESTIGATION.md) — FM-01..FM-10 seed list; two-mode model
- [FINDINGS.md](../../FINDINGS.md) — living doc; do not flatten Session A/B
- [experiments/RESULTS.md](../../../../../experiments/RESULTS.md) — E01 verdict row
- S01-00 locks: [00-PLANNER.md](./00-PLANNER.md)

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md): Agent mode → clarify if Session A/B boundary or evidence paths unclear → Plan mode → execute. **Do not conflate E01 Session A (build) with Session B (directed gap).**

## Locked defaults

| Item | Value |
|------|-------|
| Output file | `scopes/scope-01-dogfood-postmortem/POSTMORTEM.md` |
| Sessions required | E01 Session A + Session B |
| Optional prior runs | D44 (Phase 23 / ab-cms-fullstack); D45 only if RESULTS row exists |
| Session A git anchor | `f70aaea` (init G1 starter) |
| Session B git range | `704e2ff` … `a37e7c0` (gap closure → final export) |
| Verify task UUID | `e0100000-0000-4000-8000-000000000050` |
| Seed task UUIDs | `…0010`, `…0020`, `…0030`, `…0040`, `…0050` |
| G1 workspace | `experiments/ab-incident-tracker/runs/G1/` |
| B0 baseline | `experiments/ab-incident-tracker/runs/B0/` |
| Product code | **Forbidden** — docs only |

## Preflight / Plan

Before writing POSTMORTEM prose:

1. Confirm board row **P24-S01-01** is runnable (P24-S01-00 `done`).
2. Skim [00-PLANNER.md](./00-PLANNER.md) locked template — output must match §1–§4 structure.
3. Run the **Commands to run** block below; paste summaries into POSTMORTEM evidence cells.
4. Plan which FM rows are Mode A only, Mode B only, or both — do not decide from memory.

## Objective

Write **`POSTMORTEM.md`** comparing **E01 Session A (build)** vs **Session B (directed gap)**. Update **`FINDINGS.md`** failure taxonomy — do **not** conflate the two sessions.

## Session boundary (locked — do not re-debate)

| Session | Trigger | Git boundary (G1) | Graph expectation |
|---------|---------|-------------------|-------------------|
| **A — Build** | Multitask parent + `PROMPT-G1-ENFORCE.md` + `SUBAGENT-DELEGATION.md` | `f70aaea` (init) through product ship; **no** gap-closure commits | Seed tasks `…0010`–`…0050` only; 0 decisions; 0 discoveries; `internal/` **≡ B0** |
| **B — Directed gap** | Human: gap analysis + plan + fix **using Trace** | `704e2ff` … `a37e7c0` (gap closure → exports) | 7 discoveries, 2 decisions, 4 evidence; product **≠ B0**; tasks still seed-only |

Session A evidence is **sparse in git** (build may be squashed into init). Use B0 diff, RESULTS row, prompts, and pre-`704e2ff` graph snapshots if present. Session B evidence is **rich** in `trace/graph.json` + git.

## Required reads (mandatory)

Read in order. Record path + line/section in every POSTMORTEM claim.

### Phase 24 SoT

| # | Path | Extract |
|---|------|---------|
| 1 | [FINDINGS.md](../../FINDINGS.md) | Two-mode draft tables; do not overwrite Session B facts |
| 2 | [INVESTIGATION.md](../../INVESTIGATION.md) | FM-01..FM-10 seed list; “working as intended” table; investigation questions A–D |
| 3 | [experiments/RESULTS.md](../../../../../experiments/RESULTS.md) | E01 row — Session A vs B verdict |

### E01 Session A (build)

| # | Path | Extract |
|---|------|---------|
| 4 | [PROMPT-G1-ENFORCE.md](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-G1-ENFORCE.md) | Enforcement steps; Multitask workspace lock |
| 5 | [SUBAGENT-DELEGATION.md](../../../../../experiments/ab-incident-tracker/prompts/SUBAGENT-DELEGATION.md) | Orchestrator → worker split; seed UUID table |
| 6 | [PROMPT-B0.md](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-B0.md) | Feature parity bar (G1≡B0 claim) |
| 7 | `experiments/ab-incident-tracker/runs/B0/` | Baseline product tree for Session A convergence check |
| 8 | [runs/G1/seed/gt.json](../../../../../experiments/ab-incident-tracker/runs/G1/seed/gt.json) | Ground-truth seed shape; fixed task UUIDs (FM-01 anchor) |
| 9 | `git -C runs/G1 log --oneline` | Init `f70aaea`; Session B commits through `a37e7c0` |

### E01 Session B (directed gap)

| # | Path | Extract |
|---|------|---------|
| 10 | `experiments/ab-incident-tracker/runs/G1/trace/graph.json` | `tasks[]`, `discoveries[]`, `decisions[]`, `evidence[]`, `links[]` (`discovery_mentions_task`), `deliberation_states[]`, `uncertainties` (expect empty) |
| 11 | `experiments/ab-incident-tracker/runs/G1/SPEC.md` | Post-gap spec deltas (unassign, filters, redirect) |
| 12 | `experiments/ab-incident-tracker/runs/G1/PLANNING-MATRIX.md` | Task→deliverable map; feature checklist |
| 13 | `experiments/ab-incident-tracker/runs/G1/VERIFY.md` | Manual checklist + **Manual run results** table (Session B) |
| 14 | `experiments/ab-incident-tracker/runs/G1/AGENTS.md` + `.cursor/rules/trace-enforcement.mdc` + `.cursor/hooks/trace-loop-gate.sh` | Harness enforcement surface |
| 15 | `git -C runs/G1 log --oneline 704e2ff..a37e7c0` | Gap-closure commit subjects (7 commits) |
| 16 | `git -C runs/G1 show 704e2ff --stat` | Files touched in first gap-closure commit |

### Optional prior dogfood (cite if used; not required for exit)

| # | Path | Extract |
|---|------|---------|
| 17 | [Phase 23 ENFORCEMENT.md](../../../phase-23-enforcement-choke-points/ENFORCEMENT.md) | **D44** — optional; cites ab-cms-fullstack dogfood learnings |
| 18 | [Phase 23 README](../../../phase-23-enforcement-choke-points/README.md) | D44 planning-quality vs premature-impl split |
| 19 | Agent transcripts / RESULTS if D45 row appears later | Any third dogfood run — mark **not reviewed** if not found |

## Commands to run (paste summaries into POSTMORTEM)

Run from repo root unless noted. Capture stdout or file paths in citations.

```bash
# Session boundary
git -C experiments/ab-incident-tracker/runs/G1 log --oneline

# Session A: G1 vs B0 product convergence (Session A claim)
diff -rq experiments/ab-incident-tracker/runs/B0/internal \
         experiments/ab-incident-tracker/runs/G1/internal

# Session B: divergence after gap closure
git -C experiments/ab-incident-tracker/runs/G1 diff f70aaea..704e2ff --stat

# Graph counts (Session B terminal export)
python3 - <<'PY'
import json
p="experiments/ab-incident-tracker/runs/G1/trace/graph.json"
g=json.load(open(p))
for k in ("tasks","discoveries","decisions","evidence","uncertainties"):
    print(k, len(g.get(k,[])))
for d in g.get("deliberation_states",[]):
    print(d["task_id"][-4:], d["current_phase"], d.get("stop_reason"))
PY

# Live loop/gate on verify task (reconcile hop_budget vs p19_saturated)
export TRACE_BIN=/home/ali/Desktop/Trace/bin/trace
export TRACE_PROJECT_ROOT=$PWD/experiments/ab-incident-tracker/runs/G1
$TRACE_BIN loop status --task e0100000-0000-4000-8000-000000000050
$TRACE_BIN loop gate --task e0100000-0000-4000-8000-000000000050 --for edit
$TRACE_BIN loop gate --task e0100000-0000-4000-8000-000000000050 --for done

# Task creation audit (FM-01 / FM-10)
grep -c '"id":' experiments/ab-incident-tracker/runs/G1/trace/graph.json  # tasks section only — count tasks[] manually
# Optional: sqlite3 runs/G1/.trace/trace.db "SELECT type, COUNT(*) FROM entities GROUP BY type;"
```

## Deliverables

| File | Content |
|------|---------|
| `scopes/scope-01-dogfood-postmortem/POSTMORTEM.md` | Locked template below — **create this file** |
| [FINDINGS.md](../../FINDINGS.md) | Update failure taxonomy + FM confirmed/partial/rejected per session |

## POSTMORTEM.md — locked template (S01-00)

Create `POSTMORTEM.md` in this scope folder with **exactly** these sections.

### §1 Runs reviewed

- **Required:** E01 Session A (build), E01 Session B (directed gap) — dates, prompts, harness (Multitask vs single-agent).
- **Optional:** D44 (Phase 23 / ab-cms-fullstack notes), D45 (if evidence found; else “not reviewed”).
- One paragraph each: what the human/agent did, workspace path, Trace binary version if known.

### §2 Two-mode comparison table

Copy structure from [INVESTIGATION.md](../../INVESTIGATION.md) “Working as intended” table. Minimum rows:

| Dimension | Session A | Session B | Evidence pointer |
|-----------|-----------|-----------|------------------|
| Tasks in graph | | | `graph.json` `tasks[]` or FINDINGS |
| Discoveries / decisions / evidence | | | `graph.json` |
| Uncertainties | | | `graph.json` |
| Loop / gate on verify `…0050` | | | `deliberation_states[]` + live `loop status` |
| Product vs B0 | | | `diff -rq …/internal` |
| Orchestrator pattern | | | SUBAGENT-DELEGATION vs directed prompt |
| Tests / VERIFY | | | VERIFY.md |

Add rows as needed; every cell in Session A/B columns must have ≥1 citation.

### §3 Failure mode matrix (FM-01..FM-10)

| FM | Name | Session A | Session B | Both | Status | Evidence |
|----|------|-----------|-----------|------|--------|----------|
| FM-01 | Seed anchoring | Y/N/partial | Y/N/partial | Y/N | confirmed / partial / rejected / unknown | path:detail |
| FM-02 | Graph thin export | | | | | |
| … | … | | | | | |
| FM-10 | Discovery without task promotion | | | | | |

Rules:
- **Session A / B / Both:** which session(s) exhibit the symptom (not whether FM is “fixed”).
- **Status:** investigator judgment vs INVESTIGATION seed text.
- **Evidence:** repo path, git SHA, or command output — **≥8 distinct paths** across POSTMORTEM.

### §4 Open questions (handoff to S02/S04)

Bullet list only — no codebase audit here. Include at minimum:
- Why 7 discoveries but 0 new tasks?
- Was `trace add` (task type) invoked at all in Session B?
- Why verify gate blocked after green tests (`hop_budget_exceeded` vs `p19_saturated` — cite live gate)?
- Directed gap: prompt hack vs product-sufficient entry point?

## Must answer (exit gate)

1. Why did Session B record **7 discoveries** but **0 new tasks**?
2. Did Session B call **`trace add`** for tasks (not just discoveries)? Cite graph + optional `.trace/` DB.
3. Why does verify task `…0050` remain **gate-blocked** after gap work and `go test ./...` PASS?
4. Which FMs are **Mode A only**, **Mode B only**, or **both**?
5. Is “directed gap” a **prompt hack** or evidence Trace is sufficient with the right entry point?

## FINDINGS.md updates

- Set failure taxonomy row owner **S01** → `draft` with link to POSTMORTEM.md.
- For each FM-01..FM-10: `confirmed` | `partial` | `rejected` | `unknown` **per session** (sub-bullets A / B).
- Do **not** delete the two-mode model section; refine only with cited evidence.

## Forbidden

- Product Go changes
- Rewriting Phase 23 (or earlier) board history
- Codebase audit prose (defer to S02) beyond citing file paths from graph/git

## Todo updates

Per [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md): edit **only** board row **P24-S01-01** — set `in_progress` at start, `done` with evidence summary in Notes when exit criteria met. Do not rewrite P24-S01-00 or future rows.

## Minimal todos

- [ ] Read all **Required reads** tables (Phase 24 SoT + Session A + Session B)
- [ ] Run **Commands to run**; capture outputs for citations
- [ ] Create `POSTMORTEM.md` with §1–§4 per locked template below
- [ ] Complete FM-01..FM-10 matrix (Session A / B / Both + evidence)
- [ ] Update [FINDINGS.md](../../FINDINGS.md) FM section per session
- [ ] Self-check **Must answer** list before marking board row done

## Exit criteria

- [ ] `POSTMORTEM.md` exists with §1–§4 per locked template
- [ ] Session A and Session B **both** covered; not blended into one narrative
- [ ] FM table complete for FM-01..FM-10 with Session A / B / Both columns
- [ ] ≥8 distinct cited evidence paths in POSTMORTEM
- [ ] FINDINGS.md FM section updated per session
- [ ] Board row P24-S01-01: status=done, Notes only

## Next

**P24-S01-02**
