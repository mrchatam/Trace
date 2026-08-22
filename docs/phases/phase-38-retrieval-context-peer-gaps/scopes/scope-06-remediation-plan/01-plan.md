# P38-S06-01 — Author remediation plan

## Metadata
- id: P38-S06-01
- todo_ids: [P38-S06-01]
- role: implementer
- skills: [planning-and-task-breakdown, analyst, documentation-and-adrs]
- mcps: [user-trace]
- verification: manual
- hooks: none

## Objective

Author **`REMEDIATION-PLAN.md`** — ranked future remediation themes for **human-promoted Phase 39+**. Planning artifact only. **No product code. No implement in P38.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — § REMEDIATION-PLAN shape
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — planner locks from P38-S06-00
- **Required inputs (read-only — do not re-investigate):**
  - [GAP-REGISTRY.md](../scope-04-gap-registry/GAP-REGISTRY.md) — G-001…G-011, M-001, G-004a/b
  - [SATURATION-NOTES.md](../scope-05-saturation-gate/SATURATION-NOTES.md) — S05 APPROVE required; §4 H7 defer owner
  - [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md)
  - [PEER-CG.md](../scope-02-codegraph-peer/PEER-CG.md) — §4 anti-patterns
  - [PEER-UA-GF.md](../scope-03-ua-graphify-peer/PEER-UA-GF.md) — §3 MP
  - [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

**Prerequisite:** S05-02 **APPROVE (saturated)** — `ready_for_REMEDIATION_PLAN: true` in SATURATION-NOTES §7.

## Locked defaults (P38-S06-00 — do not re-debate)

| Item | Value |
|------|-------|
| Output path | `scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md` |
| Product edits | **Forbidden** (Go/TS/web) |
| Method | Synthesis of APPROVED GAP-REGISTRY + SATURATION-NOTES — **not** re-audit S01–S05 |
| Evidence | Link GAP-REGISTRY §7 + S01–S05 `$EV/` paths — no duplicate JSON |
| Theme IDs | **G1–G9** (planner consolidation of G-001…G-011 — see registry below) |
| Ranking rubric | **impact × law fit ÷ effort sketch** (locked formula — § Ranking rubric) |
| H7 owner decision | **Compose-first UX ranks above unified `trace_explore`** for near-term phases — see G2 |
| H11 stack decision | **Doc-only** — user-facing dual-stack recipe; no product integration in first remediation wave |
| G-004a vector | **Reject / defer** — DR-NOSSEM; not a remediation theme |
| Phase sketches | Titles + scope bullets only — **no board spawn in P38** |
| Non-goals | No implement rows; no "fix in P38"; no detailed Phase 39 board |

---

## Ranking rubric (LOCK — apply to every theme)

Score each theme on three axes (1–5 integers). **Rank score = (impact × law_fit) ÷ effort** — higher = earlier remediation.

| Axis | 5 (best for rank) | 3 | 1 (worst for rank) |
|------|-------------------|---|---------------------|
| **Impact** | Blocks agent orient every session (G-001/002/006/007 cluster) | Periodic friction (index, onboarding) | Docs-only / optional polish |
| **Law fit** | Strengthens Laws 6–7 progressive caps + task moat (merge query+task, adapters) | Neutral / adapter-only (G19 GUI) | Requires law review or risks DR-NOSSEM / dump / daemon |
| **Effort sketch** | 1 = docs/harness only · 2 = MCP instructions + compiler param · 3 = new retrieval channel · 4 = watcher/layers · 5 = unified explore + caps redesign |

**Tie-breakers (in order):**

1. More GAP ids covered at same score  
2. Stronger peer pattern with Trace moat preserved (M-001)  
3. Lower law-review risk  
4. Explicit SATURATION-NOTES defer owner (H7 → G2)

**Document in REMEDIATION-PLAN §2:** table with impact, law_fit, effort, computed rank score, final rank order.

---

## Theme registry G1–G9 (planner pre-lock — embed in REMEDIATION-PLAN)

Consolidates 11 gap rows into 9 ranked themes. Every G-001…G-011 must appear in exactly one theme row or §4 reject/defer.

| Theme | GAP ids | Problem (one line) | Peer pattern | Proposed future phase sketch | Risks | Explicit not P38 |
|-------|---------|-------------------|--------------|------------------------------|-------|------------------|
| **G1** Query+task orient merge | G-001, G-002 | Agents cannot merge agent query + task packet in one orient step; compiler FTS uses `task.Title` only | UA `buildChatContext(query)` + SearchEngine; MP `wake_up()` L0+L1 identity | **Phase 39 — Context orient merge:** optional `query` on `trace_context` / compiler path; merge search hits into packet; preserve task UUID + gates | Over-merge → dump risk (Law 6); query-only drift abandons moat | No compiler/MCP code in P38 |
| **G2** Read-surface strategy (H7 owner) | G-007 (+ G-006 partial) | No unified read-orient; 7/7 desk-check dimensions **not equivalent** to CG explore compose | CG `codegraph_explore` single capped call; Trace 16-tool split | **Phase 39 (compose-first):** SERVER_INSTRUCTIONS-style orient recipe, ranked read tools, optional orchestration doc — **before** **Phase 40+ (unified `trace_explore`):** task-aware capped explore after G1 + law spike | Unified tool without task scope → H1 anti-pattern; mega-tool hides write surface | No `trace_explore` implement in P38; live spike optional Phase 39 pre-gate only |
| **G3** MCP discovery & harness orient | G-006, G-010, harness 9/16 | 16 tools without "start here"; moat under-promoted vs CG 1-tool + MP 44-tool extremes | CG `DEFAULT_MCP_TOOLS` + SERVER_INSTRUCTIONS; MP categorized READ sets (contrast — do not copy count) | **Phase 39 — Harness orient:** MCP init playbook, install/bootstrap moat-first messaging, Cursor 9/16 registration hygiene | Copying MP surface; hiding task/write tools | No MCP server.go changes in P38 |
| **G4** Dual-stack documentation (H11) | G-011 | No user doc for Trace + CG complementary workflow | PEER-CG §5 complement; PEER-FIXTURES fixture paths | **Phase 39 — Docs only:** CONTRIBUTING/AGENTS section — when to `codegraph init` vs `trace index`, Law 19 boundaries, storage separation | Product dual-index default → complexity; adapter logic fork | **Doc-only** — investigation conclusion: not product integration |
| **G5** Graph-first onboarding UX | G-008 | GUI `/` → Graph route without committed orient artifact / install hook | GF `graph.html` + worked; UA `onboard-builder.ts`; MP `onboarding.py` | **Phase 39–40 — GUI orient adapter:** graph route content, optional static orient artifact, install hook narrative (Law 19 — adapter calls library) | Business logic in `web/`; full Graphify port | No web/ product code in P38 |
| **G6** Non-semantic concept retrieval | G-004b only | Title-token FTS misses concept/summary graph channels | GF EXTRACTED/INFERRED edges; MP BM25 text leg (not vector) | **Phase 40+ — Graph-label retrieval:** summary/label channel under DR-NOSSEM; law review gate before build | DR-NOSSEM slip into embeddings (G-004a); semantic creep | G-004a vector **reject** — not theme |
| **G7** Index freshness & language coverage | G-005 | 5 langs; manual `trace index`; no watcher | CG watcher debounce + 29 extractors | **Phase 40+ — Index ergonomics:** lang expansion policy; optional watch/hook path (local-first, no daemon) | CG detached daemon anti-pattern; lang sprawl without policy | No analyzer/index code in P38 |
| **G8** Progressive layers L2–L3 | G-003 | Layers 2–3 designed (`doc.go` L7) not in live packet | MP 4-layer stack (L2 on-demand / L3 deep — designed) | **Phase 41+ — Layer expansion:** ship L2–3 or revise spec with alternative | Packet bloat; layer honesty regression | No compiler layer ship in P38 |
| **G9** Intent pipeline | G-009 | RETRIEVAL_AND_CONTEXT §3 intent — zero `internal/retrieval/` code | MP `fact_checker.py` (contrast — offline check pipeline) | **Phase 41+ or doc-revise:** implement intent extraction **or** mark doc aspirational + supersede | Large scope; overlaps G1/G6; law review | No retrieval implement in P38 |

### Pre-computed rank order (planner — verify in §2 table)

| Rank | Theme | impact | law_fit | effort | score (×÷) |
|------|-------|--------|---------|--------|------------|
| 1 | **G1** Query+task merge | 5 | 5 | 3 | 8.33 |
| 2 | **G3** MCP/harness orient | 5 | 5 | 2 | 12.50 |
| 3 | **G4** Dual-stack docs (H11 doc-only) | 4 | 5 | 1 | 20.00 |
| 4 | **G5** Graph onboarding | 4 | 4 | 3 | 5.33 |
| 5 | **G2** Read-surface (compose-first → explore) | 5 | 4 | 4 | 5.00 |
| 6 | **G6** Concept retrieval G-004b | 4 | 3 | 4 | 3.00 |
| 7 | **G7** Index/watch | 3 | 4 | 4 | 3.00 |
| 8 | **G8** Layers L2–L3 | 4 | 4 | 4 | 4.00 |
| 9 | **G9** Intent pipeline | 3 | 3 | 5 | 1.80 |

**Note:** G3 ranks above G1 on raw score (lower effort) — final REMEDIATION-PLAN may swap G1/G3 only with explicit tie-breaker rationale (G1 covers more agent-blocking gaps G-001+G-002). **Planner lock:** publish **G1 → G3 → G4 → G5 → G2 → G6 → G7 → G8 → G9** with note that G1/G3 are Phase 39 co-wave candidates.

### H7 owner decision (LOCK — must appear verbatim in REMEDIATION-PLAN §2 G2 + §5)

**SATURATION-NOTES §4 defer owner = S06.** Desk-check ([`h7-compose-desk-check.md`](../../../../../experiments/runs/2026-08-22-p38-s05-663/evidence/h7-compose-desk-check.md)): **not equivalent** on 7/7 dimensions.

| Option | Rank | Rationale |
|--------|------|-----------|
| **Compose-first UX** (orient recipe + G1 merge + ranked read tools) | **Primary — Phase 39** | Preserves 16-tool write surface + task moat (M-001); lower effort; addresses G-006 discovery without mega-tool |
| **Unified `trace_explore`** (task-aware, capped, P24 transfer) | **Secondary — Phase 40+** | Requires G1 query+task merge first; law review + optional live spike; CG parity without query-only trap |

**Reject for remediation:** Claiming Trace multi-tool compose already equivalent to CG explore (desk-check closed that path).

### H11 stack recommendation (LOCK — investigation conclusion)

| Option | Verdict | Rationale |
|--------|---------|-----------|
| **Doc-only** dual-stack recipe | **Accept — G4 / Phase 39 docs** | [`h11-stack-docs.md`](../../../../../experiments/runs/2026-08-22-p38-s04-660/evidence/h11-stack-docs.md) — zero user-facing workflow; PEER-CG §5 complement exists in investigation artifacts only; Law 19 = adapters, not second product |
| **Product integration** (default dual-index, bundled MCP) | **Reject** | Violates local-first simplicity; blurs moat; no evidence Trace dogfood needs product fork |

---

## Must answer (planner handoff — embed in REMEDIATION-PLAN.md)

1. **G1–G9 themes** from GAP-REGISTRY — ranked with rubric table.
2. **Each theme:** problem, GAP ids, peer pattern, proposed future phase sketch, risks, explicit **not P38**.
3. **Reject list** — minimum 12 items (PEER-CG §4 + SATURATION §5 + plan-specific).
4. **H11:** doc vs product — **doc-only** (investigation conclusion).
5. **H7:** compose-first UX **before** unified `trace_explore` (investigation conclusion).

---

## Planning todos (run in order)

### T0 — Preflight

```bash
EV=experiments/runs/$(date +%Y-%m-%d)-p38-s06-666/evidence
mkdir -p "$EV"
```

- Confirm S05-02 APPROVE (board row 664) and `ready_for_REMEDIATION_PLAN: true`.
- Confirm GAP-REGISTRY + SATURATION-NOTES exist (rows 661, 664).
- Record date + row id in `$EV/t0-preflight.md`.

### T1 — Gap → theme coverage audit

Build `$EV/t1-gap-theme-coverage.md`:

- Table: G-001…G-011 → theme G1–G9 (or §4 reject/defer).
- Confirm **G-004a** → §4 defer only (not theme).
- Confirm **M-001** → REMEDIATION-PLAN §1 executive (moat preserved), not a gap theme.

### T2 — Ranking rubric application

For each G1–G9, document impact/law_fit/effort + score in `$EV/t2-ranking-scores.md`.

Apply tie-breakers; record final order with G1/G3 co-wave note.

### T3 — Reject / defer registry draft

Seed from:

- PEER-CG §4 + [`anti-patterns-not-for-trace.md`](../../../../../experiments/runs/2026-08-22-p38-s02-654/evidence/anti-patterns-not-for-trace.md)
- SATURATION-NOTES §5 (12 items)
- GAP-REGISTRY §4 non-gaps, §5 deferrals

Minimum **12 rejects** in `$EV/t3-reject-registry.md` covering: daemon, MCP-only loop, dump defaults, graph-only product, query-only replaces task, CG benchmarks, MP 44-tool copy, vector semantic, implement in P38, always-on daemon, product dual-stack default, compose-equivalence claim.

### T4 — Phase sketch titles (no board spawn)

`$EV/t4-phase-sketches.md` — map themes to **Phase 39 / 40 / 41+** title bullets only:

| Phase | Themes (planner lock) |
|-------|----------------------|
| **39** | G1, G3, G4 (docs), G5 (start) — co-wave OK |
| **40+** | G2 (unified explore), G6, G7, G5 (complete) |
| **41+** | G8, G9 |

### T5 — H7 + H11 decision sections

`$EV/t5-h7-h11-decisions.md` — copy locked tables from this prompt (compose-first; doc-only H11).

### T6 — Author REMEDIATION-PLAN.md

**Path:** `scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md`

**Required sections:**

```markdown
## §1 Executive summary
Saturation ref; moat M-001 preserved; ranked theme count; Phase 39 recommendation

## §2 Ranked themes G1–G9
Summary table (rank, theme, GAP ids, score, phase sketch)
Detail subsection per theme (problem, GAP ids, peer pattern, phase sketch, risks, not P38)
G2 must include H7 compose-first vs trace_explore ranking
G4 must include H11 doc-only conclusion

## §3 Proposed future phase sketches
Phase 39 / 40 / 41+ titles only — no board rows

## §4 Reject / defer registry
≥12 rejects + G-004a defer + M-001 non-gap pointer

## §5 Open questions for human promotion
Phase 39 scope cut, G2 spike gate, G9 implement vs doc-revise

## §6 Successor recommendation for S07 DR-HANDOFF
Human-promoted Phase 39; entry themes G1+G3+G4; investigation phase close
```

### T7 — Self-check

- [ ] Every G-001…G-011 in theme or §4
- [ ] No "implement in P38" language
- [ ] No Go/TS/web diff
- [ ] Rubric table present
- [ ] H7 + H11 locked decisions documented
- [ ] Links to GAP-REGISTRY evidence paths (not duplicate JSON)

---

## Exit criteria

- [ ] `REMEDIATION-PLAN.md` §§1–6 complete at scope path
- [ ] `$EV/` contains T0–T5 synthesis files (optional but recommended)
- [ ] Board row P38-S06-01 → `done` with confidence in Notes
- [ ] No product code

## Minimal todos

- [ ] T0 preflight + S05 gate confirm
- [ ] T1 gap→theme coverage
- [ ] T2 ranking scores
- [ ] T3 reject registry
- [ ] T4 phase sketches
- [ ] T5 H7/H11 decisions
- [ ] T6 author REMEDIATION-PLAN.md
- [ ] T7 self-check + board update

## Next

`P38-S06-02`
