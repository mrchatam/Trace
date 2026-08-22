# Intervention matrix — Phase 24 synthesis

Phase 24 / Scope 04 / P24-S04-01 deliverable. Ranks product/harness/protocol/docs interventions from S01–S03 evidence. **Investigation only** — no product Go; human promotes Phase 25 themes.

---

## §1 Executive summary + ranking rationale

E01 dogfood splits into **Mode A (build)** and **Mode B (directed gap)**. Both sessions share FM-01 (seed anchoring), FM-03 (loop saturation), and FM-08 (tool surface gap); Session B adds FM-10 (discovery without task promotion). No single intervention fixes all FMs — Phase 25 should ship **harness defaults that collapse Mode A→B** first, then **product loop/task promotion** to sustain directed-gap richness without custom human prompts.

**Top three (INT-03, INT-04, INT-01):**

| Rank | ID | Why it wins |
|------|-----|-------------|
| 1 | **INT-03** | Default gap-pass harness directly targets FM-09 mode collapse; install-only, no daemon; S01 Session A vs B + S03 OH/CUR peer pattern (score 9 = impact 3 × feasibility 3). |
| 2 | **INT-04** | Orchestrator Trace-first closes FM-04/05 parent bypass that thins Session A graph; harness hook change only; complements INT-03 (score 9; tie-break: secondary to explicit gap-pass entry). |
| 3 | **INT-01** | Discovery→task promotion addresses FM-10 (Session B: 7 discoveries, 0 new tasks) and shared FM-01/FM-08; product policy with tests (score 6; tie-break: FM-10 first among score-6 rows). |

**Ranking rationale:** Sort uses impact × feasibility for live Cursor dogfood. INT-03 and INT-04 score 9 because they change agent-visible behavior via harness without SQLite spikes or human product calls. Among score-6 rows, FM-10/FM-09 tie-breakers elevate INT-01 and INT-02 over MCP nudge and export honesty. INT-07–INT-09 rank lower — they improve measurement or UX but do not alone collapse build vs directed-gap. Deferred items (§4) need human gate or multi-phase spike.

**Phase 25 theme preview (1–3 recommended):** **P25-C** (default gap pass + orchestrator Trace-first), **P25-A** (discovery→task promotion), **P25-B** (loop policy + deliberation reset). See §5 and [DR-HANDOFF.md](../../DR-HANDOFF.md).

---

## §2 Ranked intervention table

| Rank | ID | Addresses | Intervention | Owner | Impact | Effort | Risk | Evidence | Phase 25 theme | Notes |
|------|-----|-----------|--------------|-------|--------|--------|------|----------|----------------|-------|
| 1 | INT-03 | FM-09, FM-04, FM-07 | Add install bundle **default gap-pass prompt** that runs after build completes so agents record discoveries/decisions without a custom human “gap analysis” prompt. | harness | high | low | agent confusion | S01 POSTMORTEM §2, §4; S03 EXTERNAL-RESEARCH §3.1, §5 | P25-C | Mirrors OH plan→build handoff; optional not mega-prompt |
| 2 | INT-04 | FM-04, FM-05, FM-09 | Require Multitask **parent orchestrator** to set `TRACE_TASK_ID` and pass `trace loop gate --for edit` on Write via `preToolUse` deny when parent edits without active task. | harness | high | low | regression | S01 POSTMORTEM §2; S02 CODEBASE-AUDIT §2 FM-04; S03 EXTERNAL-RESEARCH §3.4, §5 | P25-C | Contrast ENFORCEMENT.md L5–7 allow-without-task |
| 3 | INT-01 | FM-10, FM-01, FM-08 | Promote **BLOCKING discoveries** to new tasks via `loop apply` `spawned_tasks[]` or guided `trace add task` when discovery lacks executable backlog row. | product | high | med | agent confusion | S01 POSTMORTEM §3 FM-10, Must answer §1; S02 CODEBASE-AUDIT §2 FM-10 | P25-A | Human gate on fully autonomous spawn — see §4 |
| 4 | INT-02 | FM-03, FM-09 | Recalibrate **P19 saturation and hop budget** so first empty apply on greenfield does not sticky-STOP verify work, and post-build gap passes get a separate threshold profile. | product | high | med | regression | S01 POSTMORTEM §3 FM-03; S02 CODEBASE-AUDIT §2 FM-03, §3; S03 EXTERNAL-RESEARCH §3.5 | P25-B | Trace-specific thresholds; no peer copy-paste |
| 5 | INT-05 | FM-03, FM-09, FM-10 | Add **gap-pass deliberation reset** API that clears `Stopped`, resets `hop_count`, and reopens EXECUTE on verify task after gap closure without re-saturating on empty apply. | product | high | med | regression | S01 POSTMORTEM Must answer §3; S02 CODEBASE-AUDIT §3 deliberation reset absence; S03 EXTERNAL-RESEARCH §5 | P25-B | Episode vocabulary from GT/LIT; not full Graphiti port |
| 6 | INT-10 | FM-09, FM-06 | Codify **two-session dogfood rubric** (build vs directed gap) with separate pass/fail criteria and FM-* scoring in experiment protocol. | protocol | med | low | scope creep | S01 POSTMORTEM §2; INVESTIGATION.md Session modes; S03 EXTERNAL-RESEARCH §4 | P25-D | Prevents conflating Session A≡B0 with Session B divergence |
| 7 | INT-06 | FM-08, FM-10 | Reorder MCP **`trace_add`** descriptions and add post-discovery harness nudge to call `trace add task` or `loop apply` spawn before product edits. | product | med | low | agent confusion | S02 CODEBASE-AUDIT §2 FM-08, §4; S03 EXTERNAL-RESEARCH §3.3, §5 | P25-A | Harness nudge is install text; product owns tool ordering |
| 8 | INT-07 | FM-02, FM-01 | Enforce **`seed export --strict`** rules: fail when discoveries exist without linked tasks or when uncertainties missing on BLOCKING gaps. | product | med | med | regression | S01 POSTMORTEM §3 FM-02; S02 CODEBASE-AUDIT §2 FM-02 | P25-E | CI/scorer companion; does not fix agent behavior alone |
| 9 | INT-08 | FM-06, FM-01, FM-07 | Ship **experiment protocol v2**: arm isolation, git-sparsity scoring fix, `score.sh` graph-entity counts, mandatory `seed import` before gate. | protocol | med | med | scope creep | S01 POSTMORTEM §1 git sparsity; S02 CODEBASE-AUDIT §3 export/DB drift; S03 EXTERNAL-RESEARCH §4 | P25-D | FM-06 is protocol-layer; pairs with INT-10 |
| 10 | INT-09 | FM-03 | Unify **sticky STOP reason UX** so gate JSON and export both report the same primary `reason_code` and recovery hint (`p19_saturated` vs `hop_budget_exceeded`). | product | med | med | agent confusion | S02 CODEBASE-AUDIT §3 SelectNext first-match; S03 EXTERNAL-RESEARCH §5 sticky STOP row | P25-B | UX clarity; does not alone unblock verify |
| 11 | INT-11 | FM-05 | Add harness **hook drift verification** spike that tests G1 `trace-loop-gate.sh` against current Cursor `preToolUse` deny schema on each Cursor upgrade. | harness | low | low | scope creep | S03 EXTERNAL-RESEARCH §6 Cursor hook stability; G1 `trace-loop-gate.sh` L5–7 | P25-C | Maintenance row; pairs with INT-04 |

---

## §3 FM coverage matrix

| FM-ID | Addressed by (INT-IDs) | Residual gap |
|-------|------------------------|--------------|
| FM-01 | INT-01, INT-07, INT-08 | Seed UUID import still pins roster until promotion ships; protocol scoring catches but does not expand backlog |
| FM-02 | INT-07 | Export honesty is gate-at-export; agents may still skip writes until harness/product nudge |
| FM-03 | INT-02, INT-05, INT-09 | Full recovery requires reset + recalibration together; sticky STOP label alone insufficient |
| FM-04 | INT-03, INT-04 | Worker-only Trace may persist if parent complies but delegates graph to workers |
| FM-05 | INT-04, INT-11 | Strict config without hook deny still allows post-STOP coding until failClosed |
| FM-06 | INT-08, INT-10 | Arm isolation is protocol; product cannot enforce experiment arm boundaries |
| FM-07 | INT-03, INT-08 | Post-hoc SPEC commits reduced by gap-pass default, not eliminated without plan-before-edit product mode |
| FM-08 | INT-01, INT-06 | Tool surface improved; agent must still choose task over discovery-only path |
| FM-09 | INT-03, INT-04, INT-02, INT-05, INT-10 | Mode collapse needs harness **and** loop policy; no single row closes gap |
| FM-10 | INT-01, INT-06, INT-05 | Auto-spawn without human gate deferred §4; mentions-task links remain until promotion ships |

---

## §4 Deferred / human-gate items

- **Auto-spawn from discoveries (AR `publishEvent` analog):** AgentRQ shows event→task creation power ([S03 EXTERNAL-RESEARCH §3.3 AR row](../scope-03-external-research/EXTERNAL-RESEARCH.md)); Trace law requires **human-approved** backlog expansion. Defer fully autonomous spawn; INT-01 covers guided promotion via `loop apply` / explicit `trace add task`. **Human gate:** product owner decides auto vs confirm-before-spawn in Phase 25 P25-A planning.

- **SQLite episode model spike:** Graphiti episode boundaries are attractive ([S03 §2 GT row](../scope-03-external-research/EXTERNAL-RESEARCH.md)) but imply temporal invalidation tables and migration — **multi-phase spike**, not P0-X. INT-05 delivers minimal deliberation reset first; episode table is optional P25-B follow-on or Phase 26.

- **Graphiti/AgentRQ HTTP MCP as core path:** Rejected per [S03 §4 anti-patterns](../scope-03-external-research/EXTERNAL-RESEARCH.md#4-anti-patterns-trace-law-violations) — conflicts with local-first / no daemon on P0-X.

- **Full Graphiti port / external graph DB:** Out of Trace law; borrow vocabulary only (episode, invalidation).

- **Advisory-only Cursor rules without hooks:** Insufficient per S03 CUR row and ENFORCEMENT.md — INT-04/INT-11 cover harness path.

---

## §5 Phase 25 theme mapping

| Theme | INT-IDs | One-line scope boundary | Out of scope |
|-------|---------|-------------------------|--------------|
| **P25-A** | INT-01, INT-06 | Discovery→task promotion: spawn path, MCP ordering, harness nudge after `trace_add discovery` | Autonomous `publishEvent`-style swarm; HTTP MCP |
| **P25-B** | INT-02, INT-05, INT-09 | Loop policy recalibration: P19/hop thresholds, gap-pass deliberation reset, unified STOP reason UX | Full Graphiti episode DB; daemon |
| **P25-C** | INT-03, INT-04, INT-11 | Orchestrator + default gap pass: install prompt bundle, parent hook failClosed, hook drift checks | Rewriting Multitask Cursor product; hosted enforcement |
| **P25-D** | INT-08, INT-10 | Experiment protocol v2: two-session rubric, arm isolation, `score.sh` fix, import-before-gate | Changing E01 historical results |
| **P25-E** | INT-07 | Graph honesty: `--strict` export enforcement for discoveries/tasks/uncertainties | Replacing SQLite with export-as-SoT |

**Recommended for human promotion (1 theme per phase):** P25-C → P25-A → P25-B. Defer P25-D and P25-E unless measurement blockers appear — they strengthen evidence but do not alone change live dogfood behavior.
